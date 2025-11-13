package config

import (
	"log"
	"os"
)

type Config struct {
	DatabaseURL   string
	JWTSecret     string
	Port          string
	GCPProjectID  string
	GCSBucketName string
	PubSubTopic   string
	Environment   string
	FrontendURL   string

	// Email (SendGrid)
	SendGridAPIKey    string
	SendGridFromEmail string
	SendGridFromName  string

	// SMS (Twilio)
	TwilioAccountSID string
	TwilioAuthToken  string
	TwilioFromNumber string

	// Payments
	StripeSecretKey       string
	StripePublishableKey  string
	StripeWebhookSecret   string
	RazorpayKeyID         string
	RazorpayKeySecret     string
	RazorpayWebhookSecret string
}

func Load() *Config {
	jwtSecret := getEnv("JWT_SECRET", "your-secret-key-change-in-production")
	env := getEnv("ENVIRONMENT", "development")

	// Warn if using default JWT secret in production
	if env == "production" && jwtSecret == "your-secret-key-change-in-production" {
		log.Println("⚠️  WARNING: Using default JWT secret in production! This is a security risk!")
		log.Println("⚠️  Please set JWT_SECRET environment variable to a secure random string")
	}

	return &Config{
		DatabaseURL:   getEnv("DATABASE_URL", "sqlite://boltvisa.db"),
		JWTSecret:     jwtSecret,
		Port:          getEnv("PORT", "8080"),
		GCPProjectID:  getEnv("GCP_PROJECT_ID", ""),
		GCSBucketName: getEnv("GCS_BUCKET_NAME", ""),
		PubSubTopic:   getEnv("PUBSUB_TOPIC", "visa-notifications"),
		Environment:   getEnv("ENVIRONMENT", "development"),
		FrontendURL:   getEnv("FRONTEND_URL", "http://localhost:3000"),

		// Email
		SendGridAPIKey:    getEnv("SENDGRID_API_KEY", ""),
		SendGridFromEmail: getEnv("SENDGRID_FROM_EMAIL", "noreply@boltvisa.com"),
		SendGridFromName:  getEnv("SENDGRID_FROM_NAME", "Visa Help Center"),

		// SMS
		TwilioAccountSID: getEnv("TWILIO_ACCOUNT_SID", ""),
		TwilioAuthToken:  getEnv("TWILIO_AUTH_TOKEN", ""),
		TwilioFromNumber: getEnv("TWILIO_FROM_NUMBER", ""),

		// Payments
		StripeSecretKey:       getEnv("STRIPE_SECRET_KEY", ""),
		StripePublishableKey:  getEnv("STRIPE_PUBLISHABLE_KEY", ""),
		StripeWebhookSecret:   getEnv("STRIPE_WEBHOOK_SECRET", ""),
		RazorpayKeyID:         getEnv("RAZORPAY_KEY_ID", ""),
		RazorpayKeySecret:     getEnv("RAZORPAY_KEY_SECRET", ""),
		RazorpayWebhookSecret: getEnv("RAZORPAY_WEBHOOK_SECRET", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
