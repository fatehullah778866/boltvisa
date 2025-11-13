package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/boltvisa/api/internal/models"
	"gorm.io/gorm"
)

type AuditService struct {
	db *gorm.DB
}

func NewAuditService(db *gorm.DB) *AuditService {
	return &AuditService{db: db}
}

func (s *AuditService) Log(ctx context.Context, userID *uint, action models.AuditAction, resource string, resourceID *uint, description, ipAddress, userAgent string, metadata map[string]interface{}) error {
	var metadataJSON string
	if metadata != nil {
		bytes, err := json.Marshal(metadata)
		if err == nil {
			metadataJSON = string(bytes)
		}
	}

	auditLog := models.AuditLog{
		UserID:      userID,
		Action:      action,
		Resource:    resource,
		ResourceID:  resourceID,
		Description: description,
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
		Metadata:    metadataJSON,
	}

	return s.db.Create(&auditLog).Error
}

func (s *AuditService) LogUserAction(ctx context.Context, userID uint, action models.AuditAction, resource string, resourceID *uint, description, ipAddress, userAgent string) error {
	return s.Log(ctx, &userID, action, resource, resourceID, description, ipAddress, userAgent, nil)
}

func (s *AuditService) LogAdminAction(ctx context.Context, adminID uint, action models.AuditAction, resource string, resourceID *uint, description, ipAddress, userAgent string, metadata map[string]interface{}) error {
	return s.Log(ctx, &adminID, action, resource, resourceID, description, ipAddress, userAgent, metadata)
}

func (s *AuditService) LogPayment(ctx context.Context, userID uint, paymentID uint, amount float64, status string, ipAddress, userAgent string) error {
	metadata := map[string]interface{}{
		"payment_id": paymentID,
		"amount":     amount,
		"status":     status,
	}
	return s.Log(ctx, &userID, models.ActionPayment, "payment", &paymentID,
		fmt.Sprintf("Payment %s: $%.2f", status, amount), ipAddress, userAgent, metadata)
}

func (s *AuditService) LogRoleChange(ctx context.Context, adminID uint, targetUserID uint, oldRole, newRole string, ipAddress, userAgent string) error {
	metadata := map[string]interface{}{
		"target_user_id": targetUserID,
		"old_role":       oldRole,
		"new_role":       newRole,
	}
	return s.Log(ctx, &adminID, models.ActionRoleChange, "user", &targetUserID,
		fmt.Sprintf("Role changed from %s to %s", oldRole, newRole), ipAddress, userAgent, metadata)
}
