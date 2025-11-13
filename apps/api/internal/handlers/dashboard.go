package handlers

import (
	"net/http"

	"github.com/boltvisa/api/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DashboardData struct {
	Applications  int `json:"applications"`
	Pending       int `json:"pending"`
	Approved      int `json:"approved"`
	Rejected      int `json:"rejected"`
	Notifications int `json:"notifications"`
}

func (h *Handlers) GetDashboard(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		WriteErr(c, http.StatusUnauthorized, "unauthorized", "User not authenticated")
		return
	}

	uid, ok := userID.(uint)
	if !ok {
		WriteErr(c, http.StatusInternalServerError, "invalid_user_id", "Invalid user ID type")
		return
	}

	// Count applications by status
	var totalApps, pending, approved, rejected int64

	if err := h.db.Model(&models.VisaApplication{}).
		Where("user_id = ?", uid).
		Count(&totalApps).Error; err != nil && err != gorm.ErrRecordNotFound {
		WriteErr(c, http.StatusInternalServerError, "database_error", "Failed to count applications")
		return
	}

	if err := h.db.Model(&models.VisaApplication{}).
		Where("user_id = ? AND status = ?", uid, "in_review").
		Count(&pending).Error; err != nil && err != gorm.ErrRecordNotFound {
		WriteErr(c, http.StatusInternalServerError, "database_error", "Failed to count pending applications")
		return
	}

	if err := h.db.Model(&models.VisaApplication{}).
		Where("user_id = ? AND status = ?", uid, "approved").
		Count(&approved).Error; err != nil && err != gorm.ErrRecordNotFound {
		WriteErr(c, http.StatusInternalServerError, "database_error", "Failed to count approved applications")
		return
	}

	if err := h.db.Model(&models.VisaApplication{}).
		Where("user_id = ? AND status = ?", uid, "rejected").
		Count(&rejected).Error; err != nil && err != gorm.ErrRecordNotFound {
		WriteErr(c, http.StatusInternalServerError, "database_error", "Failed to count rejected applications")
		return
	}

	// Count unread notifications
	var unreadNotifications int64
	if err := h.db.Model(&models.Notification{}).
		Where("user_id = ? AND read = ?", uid, false).
		Count(&unreadNotifications).Error; err != nil && err != gorm.ErrRecordNotFound {
		// Don't fail if notifications count fails, just use 0
		unreadNotifications = 0
	}

	data := DashboardData{
		Applications:  int(totalApps),
		Pending:       int(pending),
		Approved:      int(approved),
		Rejected:      int(rejected),
		Notifications: int(unreadNotifications),
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}
