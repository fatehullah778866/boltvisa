package database

import (
	"log"
	"strings"
	"time"

	"github.com/boltvisa/api/internal/models"
	"github.com/glebarez/sqlite" // Pure Go SQLite driver (no CGO required)
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(databaseURL string) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	// Use SQLite for local development if database URL is empty or contains "sqlite"
	if databaseURL == "" || strings.Contains(databaseURL, "sqlite") {
		dbPath := "boltvisa.db"
		if databaseURL != "" && strings.Contains(databaseURL, "sqlite://") {
			// Extract path from sqlite://path format
			dbPath = strings.TrimPrefix(databaseURL, "sqlite://")
		}
		log.Printf("Using SQLite database: %s", dbPath)
		// Use pure Go SQLite driver (no CGO required)
		// Configure SQLite with connection pool and timeout settings
		db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
			PrepareStmt: true, // Enable prepared statements for better performance
		})
	} else {
		log.Printf("Using PostgreSQL database")
		db, err = gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	}

	if err != nil {
		log.Printf("ERROR: Failed to connect to database: %v", err)
		return nil, err
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	log.Println("✅ Database connection established")
	return db, nil
}

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.User{},
		&models.VisaCategory{},
		&models.VisaApplication{},
		&models.Document{},
		&models.Notification{},
		&models.Payment{},
		&models.AuditLog{},
		&models.PasswordResetToken{},
	)
	if err != nil {
		log.Printf("ERROR: Failed to run migrations: %v", err)
		return err
	}

	log.Println("✅ Database migrations completed")
	return nil
}
