package models

import (
	"time"

	"gorm.io/gorm"
)

type DocumentType string

const (
	DocTypePassport      DocumentType = "passport"
	DocTypePhoto         DocumentType = "photo"
	DocTypeBankStatement DocumentType = "bank_statement"
	DocTypeInvitation    DocumentType = "invitation"
	DocTypeOther         DocumentType = "other"
)

type Document struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	ApplicationID uint           `gorm:"not null;index" json:"application_id"`
	Type          DocumentType   `gorm:"type:varchar(50);not null" json:"type"`
	Name          string         `gorm:"not null" json:"name"`
	FileName      string         `gorm:"not null" json:"file_name"`
	GCSURL        string         `gorm:"not null" json:"gcs_url"`
	Size          int64          `json:"size"` // Size in bytes
	MimeType      string         `json:"mime_type"`
	UploadedBy    uint           `gorm:"not null" json:"uploaded_by"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Application VisaApplication `gorm:"foreignKey:ApplicationID" json:"application,omitempty"`
	Uploader    User            `gorm:"foreignKey:UploadedBy" json:"uploader,omitempty"`
}

func (Document) TableName() string {
	return "documents"
}
