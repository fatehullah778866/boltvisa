package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/boltvisa/api/internal/models"
	"github.com/boltvisa/api/internal/services"
	"github.com/boltvisa/api/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

func (h *Handlers) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user
	var user models.User
	if err := h.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Don't reveal if user exists or not (security best practice)
			c.JSON(http.StatusOK, gin.H{"message": "If the email exists, a password reset link has been sent"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Generate reset token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	token := hex.EncodeToString(tokenBytes)

	// Create password reset token (expires in 1 hour)
	resetToken := models.PasswordResetToken{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour),
		Used:      false,
	}

	if err := h.db.Create(&resetToken).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reset token"})
		return
	}

	// Send email with reset link
	emailService := services.NewEmailService(h.cfg.SendGridAPIKey, h.cfg.SendGridFromEmail, h.cfg.SendGridFromName)
	if emailService != nil {
		resetURL := h.cfg.FrontendURL + "/reset-password?token=" + token
		emailContent := `
		<!DOCTYPE html>
		<html>
		<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
			<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
				<h2>Password Reset Request</h2>
				<p>Hello ` + user.FirstName + `,</p>
				<p>You requested to reset your password. Click the link below to reset it:</p>
				<p><a href="` + resetURL + `" style="background-color: #0ea5e9; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px; display: inline-block;">Reset Password</a></p>
				<p>Or copy and paste this URL into your browser:</p>
				<p style="word-break: break-all;">` + resetURL + `</p>
				<p>This link will expire in 1 hour.</p>
				<p>If you didn't request this, please ignore this email.</p>
				<p>Best regards,<br>Visa Help Center Team</p>
			</div>
		</body>
		</html>
		`

		// Send email (non-blocking)
		go func() {
			_ = emailService.SendEmail(c.Request.Context(), user.Email, user.FirstName+" "+user.LastName,
				"Password Reset Request", emailContent)
		}()
	}

	// Log audit event (non-blocking)
	if h.auditService != nil {
		go func() {
			_ = h.auditService.Log(c.Request.Context(), &user.ID, models.ActionLogin, "user", &user.ID,
				"Password reset requested", c.ClientIP(), c.Request.UserAgent(), nil)
		}()
	}

	c.JSON(http.StatusOK, gin.H{"message": "If the email exists, a password reset link has been sent"})
}

func (h *Handlers) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find reset token
	var resetToken models.PasswordResetToken
	if err := h.db.Where("token = ? AND used = ?", req.Token, false).First(&resetToken).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired reset token"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Check if token expired
	if time.Now().After(resetToken.ExpiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Reset token has expired"})
		return
	}

	// Find user
	var user models.User
	if err := h.db.First(&user, resetToken.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Update password
	user.Password = hashedPassword
	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	// Mark token as used
	resetToken.Used = true
	h.db.Save(&resetToken)

	// Log audit event (non-blocking)
	if h.auditService != nil {
		go func() {
			_ = h.auditService.Log(c.Request.Context(), &user.ID, models.ActionUpdate, "user", &user.ID,
				"Password reset completed", c.ClientIP(), c.Request.UserAgent(), nil)
		}()
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password has been reset successfully"})
}
