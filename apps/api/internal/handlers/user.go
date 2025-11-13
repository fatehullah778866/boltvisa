package handlers

import (
	"net/http"
	"strings"

	"github.com/boltvisa/api/internal/models"
	"github.com/boltvisa/api/internal/utils"
	"github.com/gin-gonic/gin"
)

func (h *Handlers) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, user)
}

func (h *Handlers) UpdateCurrentUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var updateData struct {
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		PhoneNumber string `json:"phone_number"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if updateData.FirstName != "" {
		user.FirstName = updateData.FirstName
	}
	if updateData.LastName != "" {
		user.LastName = updateData.LastName
	}
	if updateData.PhoneNumber != "" {
		user.PhoneNumber = updateData.PhoneNumber
	}

	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, user)
}

func (h *Handlers) ListUsers(c *gin.Context) {
	query := h.db.Model(&models.User{})

	// Filters
	if role := c.Query("role"); role != "" {
		query = query.Where("role = ?", role)
	}
	if active := c.Query("active"); active != "" {
		query = query.Where("active = ?", active == "true")
	}
	if search := c.Query("search"); search != "" {
		searchPattern := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(email) LIKE ? OR LOWER(first_name) LIKE ? OR LOWER(last_name) LIKE ?",
			searchPattern, searchPattern, searchPattern)
	}

	// Pagination
	params := utils.GetPaginationParams(c)
	var users []models.User

	paginated, err := utils.Paginate(query, params, &users)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	// Remove passwords
	if usersData, ok := paginated.Data.(*[]models.User); ok {
		for i := range *usersData {
			(*usersData)[i].Password = ""
		}
	}

	c.JSON(http.StatusOK, paginated)
}

func (h *Handlers) UpdateUser(c *gin.Context) {
	userID := c.Param("id")

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var updateData struct {
		FirstName   string          `json:"first_name"`
		LastName    string          `json:"last_name"`
		PhoneNumber string          `json:"phone_number"`
		Role        models.UserRole `json:"role"`
		Active      *bool           `json:"active"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if updateData.FirstName != "" {
		user.FirstName = updateData.FirstName
	}
	if updateData.LastName != "" {
		user.LastName = updateData.LastName
	}
	if updateData.PhoneNumber != "" {
		user.PhoneNumber = updateData.PhoneNumber
	}
	oldRole := user.Role
	if updateData.Role != "" {
		user.Role = updateData.Role
	}
	if updateData.Active != nil {
		user.Active = *updateData.Active
	}

	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	// Log audit event for role changes (non-blocking)
	if h.auditService != nil && updateData.Role != "" && oldRole != updateData.Role {
		adminID, exists := c.Get("userID")
		if exists {
			if adminIDUint, ok := adminID.(uint); ok {
				go func() {
					_ = h.auditService.LogRoleChange(c.Request.Context(), adminIDUint, user.ID,
						string(oldRole), string(updateData.Role), c.ClientIP(), c.Request.UserAgent())
				}()
			}
		}
	}

	user.Password = ""
	c.JSON(http.StatusOK, user)
}
