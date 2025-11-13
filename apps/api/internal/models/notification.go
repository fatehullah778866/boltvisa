package models

import (
	"time"
)

type NotificationType string

const (
	NotifTypeApplicationUpdate NotificationType = "application_update"
	NotifTypeDocumentRequest   NotificationType = "document_request"
	NotifTypePayment           NotificationType = "payment"
	NotifTypeSystem            NotificationType = "system"
)

type Notification struct {
	ID        uint             `gorm:"primaryKey" json:"id"`
	UserID    uint             `gorm:"not null;index" json:"user_id"`
	Type      NotificationType `gorm:"type:varchar(50);not null" json:"type"`
	Title     string           `gorm:"not null" json:"title"`
	Message   string           `gorm:"type:text;not null" json:"message"`
	Read      bool             `gorm:"default:false" json:"read"`
	ReadAt    *time.Time       `json:"read_at,omitempty"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`

	// Relations
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (Notification) TableName() string {
	return "notifications"
}
