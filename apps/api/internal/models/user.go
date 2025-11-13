package models

import (
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	RoleAdmin      UserRole = "admin"
	RoleConsultant UserRole = "consultant"
	RoleApplicant  UserRole = "applicant"
)

type User struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Email       string         `gorm:"uniqueIndex;not null" json:"email"`
	Password    string         `gorm:"not null" json:"-"` // Never return password in JSON
	FirstName   string         `json:"first_name"`
	LastName    string         `json:"last_name"`
	PhoneNumber string         `gorm:"type:varchar(20);index" json:"phone_number,omitempty"`
	Role        UserRole       `gorm:"type:varchar(20);default:'applicant'" json:"role"`
	Active      bool           `gorm:"default:true" json:"active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "users"
}
