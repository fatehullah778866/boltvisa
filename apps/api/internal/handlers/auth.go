package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/boltvisa/api/internal/models"
	"github.com/boltvisa/api/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RegisterRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	PhoneNumber string `json:"phone_number,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token     string      `json:"token"`
	User      models.User `json:"user"`
	ExpiresAt time.Time   `json:"expires_at"`
}

func (h *Handlers) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		WriteErr(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	// Check if user exists
	var existingUser models.User
	if err := h.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		WriteErr(c, http.StatusConflict, "email_already_registered", map[string]interface{}{"email": req.Email})
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		WriteErr(c, http.StatusInternalServerError, "Failed to hash password", nil)
		return
	}

	// Create user
	user := models.User{
		Email:       req.Email,
		Password:    hashedPassword,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		Role:        models.RoleApplicant,
		Active:      true,
	}

	if err := h.db.Create(&user).Error; err != nil {
		WriteErr(c, http.StatusInternalServerError, "Failed to create user", nil)
		return
	}

	// Generate JWT
	token, err := utils.GenerateJWT(user.ID, user.Email, string(user.Role), h.cfg.JWTSecret)
	if err != nil {
		WriteErr(c, http.StatusInternalServerError, "Failed to generate token", nil)
		return
	}

	user.Password = "" // Don't return password

	// Log audit event (non-blocking)
	if h.auditService != nil {
		go func() {
			_ = h.auditService.Log(c.Request.Context(), &user.ID, models.ActionCreate, "user", &user.ID,
				"User registered", c.ClientIP(), c.Request.UserAgent(), nil)
		}()
	}

	c.JSON(http.StatusCreated, AuthResponse{
		Token:     token,
		User:      user,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	})
}

func (h *Handlers) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		WriteErr(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	// Find user
	var user models.User
	if err := h.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			WriteErr(c, http.StatusUnauthorized, "invalid_credentials", nil)
			return
		}
		WriteErr(c, http.StatusInternalServerError, "Database error", nil)
		return
	}

	// Check password
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		WriteErr(c, http.StatusUnauthorized, "invalid_credentials", nil)
		return
	}

	// Check if user is active
	if !user.Active {
		WriteErr(c, http.StatusForbidden, "Account is inactive", nil)
		return
	}

	// Generate JWT
	token, err := utils.GenerateJWT(user.ID, user.Email, string(user.Role), h.cfg.JWTSecret)
	if err != nil {
		WriteErr(c, http.StatusInternalServerError, "Failed to generate token", nil)
		return
	}

	user.Password = "" // Don't return password

	// Log audit event (non-blocking)
	if h.auditService != nil {
		go func() {
			_ = h.auditService.LogUserAction(c.Request.Context(), user.ID, models.ActionLogin, "user", &user.ID,
				"User logged in", c.ClientIP(), c.Request.UserAgent())
		}()
	}

	c.JSON(http.StatusOK, AuthResponse{
		Token:     token,
		User:      user,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	})
}

func (h *Handlers) RefreshToken(c *gin.Context) {
	// Get token from Authorization header (allows expired tokens)
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		WriteErr(c, http.StatusUnauthorized, "authorization_required", "Authorization header required")
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		WriteErr(c, http.StatusUnauthorized, "invalid_authorization", "Invalid authorization header format")
		return
	}

	tokenString := parts[1]

	// Validate token (allowing expired tokens for refresh)
	claims, err := utils.ValidateJWTAllowExpired(tokenString, h.cfg.JWTSecret)
	if err != nil {
		WriteErr(c, http.StatusUnauthorized, "invalid_token", "Invalid or malformed token")
		return
	}

	// Find user
	var user models.User
	if err := h.db.First(&user, claims.UserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			WriteErr(c, http.StatusNotFound, "user_not_found", "User not found")
			return
		}
		WriteErr(c, http.StatusInternalServerError, "database_error", "Failed to find user")
		return
	}

	// Check if user is active
	if !user.Active {
		WriteErr(c, http.StatusForbidden, "account_inactive", "Account is inactive")
		return
	}

	// Generate new token
	token, err := utils.GenerateJWT(user.ID, user.Email, string(user.Role), h.cfg.JWTSecret)
	if err != nil {
		WriteErr(c, http.StatusInternalServerError, "token_generation_failed", "Failed to generate token")
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, AuthResponse{
		Token:     token,
		User:      user,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	})
}
