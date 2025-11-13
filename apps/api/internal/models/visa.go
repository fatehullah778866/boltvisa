package models

import (
	"time"

	"gorm.io/gorm"
)

type VisaStatus string

const (
	StatusDraft     VisaStatus = "draft"
	StatusSubmitted VisaStatus = "submitted"
	StatusInReview  VisaStatus = "in_review"
	StatusApproved  VisaStatus = "approved"
	StatusRejected  VisaStatus = "rejected"
	StatusCancelled VisaStatus = "cancelled"
)

type VisaCategory struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	Country     string    `gorm:"not null" json:"country"`
	Duration    string    `json:"duration"` // e.g., "90 days", "1 year"
	Price       float64   `gorm:"default:0" json:"price"`
	Active      bool      `gorm:"default:true" json:"active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (VisaCategory) TableName() string {
	return "visa_categories"
}

type VisaApplication struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	UserID         uint           `gorm:"not null;index" json:"user_id"`
	ConsultantID   *uint          `json:"consultant_id,omitempty"` // Optional consultant
	CategoryID     uint           `gorm:"not null;index" json:"category_id"`
	Status         VisaStatus     `gorm:"type:varchar(20);default:'draft'" json:"status"`
	PassportNumber string         `json:"passport_number"`
	DateOfBirth    *time.Time     `json:"date_of_birth"`
	Nationality    string         `json:"nationality"`
	TravelDate     *time.Time     `json:"travel_date"`
	Notes          string         `gorm:"type:text" json:"notes"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User       User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Consultant *User        `gorm:"foreignKey:ConsultantID" json:"consultant,omitempty"`
	Category   VisaCategory `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Documents  []Document   `gorm:"foreignKey:ApplicationID" json:"documents,omitempty"`
}

func (VisaApplication) TableName() string {
	return "visa_applications"
}
