package services

import (
	"context"
	"fmt"

	"github.com/boltvisa/api/internal/models"
	"gorm.io/gorm"
)

type NotificationService struct {
	db     *gorm.DB
	pubsub *PubSubClient
	email  *EmailService
	sms    *SMSService
}

func NewNotificationService(db *gorm.DB, pubsub *PubSubClient, email *EmailService, sms *SMSService) *NotificationService {
	return &NotificationService{
		db:     db,
		pubsub: pubsub,
		email:  email,
		sms:    sms,
	}
}

func (s *NotificationService) CreateNotification(ctx context.Context, userID uint, notifType models.NotificationType, title, message string, metadata map[string]interface{}) error {
	// Create notification in database
	notification := models.Notification{
		UserID:  userID,
		Type:    notifType,
		Title:   title,
		Message: message,
		Read:    false,
	}

	if err := s.db.Create(&notification).Error; err != nil {
		return err
	}

	// Get user for email/phone
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return err
	}

	// Publish to Pub/Sub if available
	if s.pubsub != nil {
		msg := NotificationMessage{
			UserID:   userID,
			Type:     string(notifType),
			Title:    title,
			Message:  message,
			Email:    user.Email,
			Metadata: metadata,
		}
		if err := s.pubsub.PublishNotification(ctx, msg); err != nil {
			// Log error but don't fail the operation
			// In production, you might want to retry or use a dead letter queue
		}
	}

	return nil
}

func (s *NotificationService) SendApplicationUpdate(ctx context.Context, userID uint, applicationID uint, status string) error {
	title := "Application Status Updated"
	message := fmt.Sprintf("Your visa application status has been updated to: %s", status)

	metadata := map[string]interface{}{
		"application_id": applicationID,
		"status":         status,
	}

	return s.CreateNotification(ctx, userID, models.NotifTypeApplicationUpdate, title, message, metadata)
}

func (s *NotificationService) SendDocumentRequest(ctx context.Context, userID uint, applicationID uint, docType string) error {
	title := "Document Request"
	message := fmt.Sprintf("Please upload the following document: %s", docType)

	metadata := map[string]interface{}{
		"application_id": applicationID,
		"document_type":  docType,
	}

	return s.CreateNotification(ctx, userID, models.NotifTypeDocumentRequest, title, message, metadata)
}

func (s *NotificationService) SendPaymentNotification(ctx context.Context, userID uint, amount float64, status string) error {
	title := "Payment Update"
	message := fmt.Sprintf("Payment of $%.2f has been %s", amount, status)

	metadata := map[string]interface{}{
		"amount": amount,
		"status": status,
	}

	return s.CreateNotification(ctx, userID, models.NotifTypePayment, title, message, metadata)
}
