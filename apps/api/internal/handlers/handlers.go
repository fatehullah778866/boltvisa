package handlers

import (
	"context"
	"log"

	"github.com/boltvisa/api/internal/config"
	"github.com/boltvisa/api/internal/services"
	"gorm.io/gorm"
)

type Handlers struct {
	db                  *gorm.DB
	cfg                 *config.Config
	notificationService *services.NotificationService
	auditService        *services.AuditService
}

func New(db *gorm.DB, cfg *config.Config) *Handlers {
	ctx := context.Background()

	// Initialize Pub/Sub client
	var pubsubClient *services.PubSubClient
	if cfg.GCPProjectID != "" && cfg.PubSubTopic != "" {
		var err error
		pubsubClient, err = services.NewPubSubClient(ctx, cfg.GCPProjectID, cfg.PubSubTopic, "")
		if err != nil {
			log.Printf("WARNING: Failed to initialize Pub/Sub client: %v. Continuing without Pub/Sub notifications.", err)
		} else {
			log.Println("✅ Pub/Sub client initialized successfully")
		}
	}

	// Initialize email service
	emailService := services.NewEmailService(cfg.SendGridAPIKey, cfg.SendGridFromEmail, cfg.SendGridFromName)

	// Initialize SMS service
	smsService := services.NewSMSService(cfg.TwilioAccountSID, cfg.TwilioAuthToken, cfg.TwilioFromNumber)

	// Initialize notification service
	notificationService := services.NewNotificationService(db, pubsubClient, emailService, smsService)

	// Initialize audit service
	auditService := services.NewAuditService(db)

	return &Handlers{
		db:                  db,
		cfg:                 cfg,
		notificationService: notificationService,
		auditService:        auditService,
	}
}
