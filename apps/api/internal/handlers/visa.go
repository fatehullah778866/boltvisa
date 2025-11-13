package handlers

import (
	"net/http"
	"strings"

	"github.com/boltvisa/api/internal/models"
	"github.com/boltvisa/api/internal/utils"
	"github.com/gin-gonic/gin"
)

func (h *Handlers) GetVisaCategories(c *gin.Context) {
	var categories []models.VisaCategory
	query := h.db.Where("active = ?", true)

	if err := query.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch visa categories"})
		return
	}

	c.JSON(http.StatusOK, categories)
}

func (h *Handlers) GetVisaCategory(c *gin.Context) {
	id := c.Param("id")

	var category models.VisaCategory
	if err := h.db.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Visa category not found"})
		return
	}

	c.JSON(http.StatusOK, category)
}

func (h *Handlers) CreateVisaCategory(c *gin.Context) {
	var category models.VisaCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create visa category"})
		return
	}

	c.JSON(http.StatusCreated, category)
}

func (h *Handlers) UpdateVisaCategory(c *gin.Context) {
	id := c.Param("id")

	var category models.VisaCategory
	if err := h.db.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Visa category not found"})
		return
	}

	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update visa category"})
		return
	}

	c.JSON(http.StatusOK, category)
}

func (h *Handlers) GetApplications(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userRole, exists := c.Get("userRole")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return
	}

	role, ok := userRole.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user role type"})
		return
	}

	query := h.db.Model(&models.VisaApplication{}).Preload("Category").Preload("User")

	// Applicants see only their applications, consultants see their clients', admins see all
	if role == "applicant" {
		query = query.Where("user_id = ?", userID)
	} else if role == "consultant" {
		query = query.Where("consultant_id = ?", userID)
	}

	// Filters
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if categoryID := c.Query("category_id"); categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if search := c.Query("search"); search != "" {
		searchPattern := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(passport_number) LIKE ? OR LOWER(nationality) LIKE ?", searchPattern, searchPattern)
	}

	// Pagination
	params := utils.GetPaginationParams(c)
	var applications []models.VisaApplication

	paginated, err := utils.Paginate(query, params, &applications)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch applications"})
		return
	}

	c.JSON(http.StatusOK, paginated)
}

func (h *Handlers) GetApplication(c *gin.Context) {
	id := c.Param("id")
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userRole, exists := c.Get("userRole")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return
	}

	var application models.VisaApplication
	query := h.db.Preload("Category").Preload("User").Preload("Documents")

	if err := query.First(&application, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	// Check permissions
	if role, ok := userRole.(string); ok && role == "applicant" {
		if userIDVal, ok := userID.(uint); ok && application.UserID != userIDVal {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}

	c.JSON(http.StatusOK, application)
}

func (h *Handlers) CreateApplication(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDVal, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	var application models.VisaApplication
	if err := c.ShouldBindJSON(&application); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	application.UserID = userIDVal
	application.Status = models.StatusDraft

	if err := h.db.Create(&application).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create application"})
		return
	}

	if err := h.db.Preload("Category").First(&application, application.ID).Error; err != nil {
		// Log error but continue - application was created
	}

	// Send notification (non-blocking)
	if h.notificationService != nil {
		go func() {
			_ = h.notificationService.CreateNotification(c.Request.Context(), userIDVal,
				models.NotifTypeApplicationUpdate, "Application Created",
				"Your visa application has been created successfully.", nil)
		}()
	}

	c.JSON(http.StatusCreated, application)
}

func (h *Handlers) UpdateApplication(c *gin.Context) {
	id := c.Param("id")
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userRole, exists := c.Get("userRole")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return
	}

	var application models.VisaApplication
	if err := h.db.First(&application, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	// Check permissions
	if role, ok := userRole.(string); ok && role == "applicant" {
		if userIDVal, ok := userID.(uint); ok && application.UserID != userIDVal {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}

	var oldStatus models.VisaStatus
	oldStatus = application.Status

	if err := c.ShouldBindJSON(&application); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Save(&application).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update application"})
		return
	}

	if err := h.db.Preload("Category").First(&application, application.ID).Error; err != nil {
		// Log error but continue - application was updated
	}

	// Send notification if status changed (non-blocking)
	if h.notificationService != nil && oldStatus != application.Status {
		go func() {
			_ = h.notificationService.SendApplicationUpdate(c.Request.Context(), application.UserID, application.ID, string(application.Status))
		}()
	}

	c.JSON(http.StatusOK, application)
}

func (h *Handlers) DeleteApplication(c *gin.Context) {
	id := c.Param("id")
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userRole, exists := c.Get("userRole")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return
	}

	var application models.VisaApplication
	if err := h.db.First(&application, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	// Check permissions
	if role, ok := userRole.(string); ok && role == "applicant" {
		if userIDVal, ok := userID.(uint); ok && application.UserID != userIDVal {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}

	// Only allow deletion of draft applications
	if application.Status != models.StatusDraft {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Can only delete draft applications"})
		return
	}

	if err := h.db.Delete(&application).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete application"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Application deleted"})
}

func (h *Handlers) GetConsultantApplications(c *gin.Context) {
	consultantID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	query := h.db.Model(&models.VisaApplication{}).Preload("Category").Preload("User").
		Where("consultant_id = ?", consultantID)

	// Filters
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if userID := c.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	// Pagination
	params := utils.GetPaginationParams(c)
	var applications []models.VisaApplication

	paginated, err := utils.Paginate(query, params, &applications)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch applications"})
		return
	}

	c.JSON(http.StatusOK, paginated)
}

func (h *Handlers) GetConsultantClients(c *gin.Context) {
	consultantID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var applications []models.VisaApplication
	if err := h.db.Where("consultant_id = ?", consultantID).
		Select("DISTINCT user_id").
		Find(&applications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch clients"})
		return
	}

	// Extract unique user IDs
	userIDs := make([]uint, 0)
	seen := make(map[uint]bool)
	for _, app := range applications {
		if !seen[app.UserID] {
			userIDs = append(userIDs, app.UserID)
			seen[app.UserID] = true
		}
	}

	var users []models.User
	if err := h.db.Where("id IN ?", userIDs).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch clients"})
		return
	}

	// Remove passwords
	for i := range users {
		users[i].Password = ""
	}

	c.JSON(http.StatusOK, users)
}

func (h *Handlers) AssignConsultant(c *gin.Context) {
	id := c.Param("id")
	userRole, exists := c.Get("userRole")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return
	}

	role, ok := userRole.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user role type"})
		return
	}

	var application models.VisaApplication
	if err := h.db.First(&application, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	// Only admins and consultants can assign consultants
	if role != "admin" && role != "consultant" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var req struct {
		ConsultantID *uint `json:"consultant_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify consultant exists and has consultant role if provided
	if req.ConsultantID != nil {
		var consultant models.User
		if err := h.db.Where("id = ? AND role = ?", *req.ConsultantID, "consultant").First(&consultant).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid consultant"})
			return
		}
		application.ConsultantID = req.ConsultantID
	} else {
		// Unassign consultant
		application.ConsultantID = nil
	}

	if err := h.db.Save(&application).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign consultant"})
		return
	}

	if err := h.db.Preload("Category").Preload("Consultant").Preload("User").First(&application, application.ID).Error; err != nil {
		// Log error but continue - consultant was assigned
	}
	c.JSON(http.StatusOK, application)
}
