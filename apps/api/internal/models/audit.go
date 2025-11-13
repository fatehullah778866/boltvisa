package models

import (
	"time"
)

type AuditAction string

const (
	ActionCreate     AuditAction = "create"
	ActionUpdate     AuditAction = "update"
	ActionDelete     AuditAction = "delete"
	ActionLogin      AuditAction = "login"
	ActionLogout     AuditAction = "logout"
	ActionPayment    AuditAction = "payment"
	ActionDocument   AuditAction = "document"
	ActionRoleChange AuditAction = "role_change"
)

type AuditLog struct {
	ID          uint        `gorm:"primaryKey" json:"id"`
	UserID      *uint       `json:"user_id,omitempty"` // Nullable for system actions
	Action      AuditAction `gorm:"type:varchar(50);not null" json:"action"`
	Resource    string      `gorm:"type:varchar(100);not null" json:"resource"` // e.g., "user", "application", "payment"
	ResourceID  *uint       `json:"resource_id,omitempty"`
	Description string      `gorm:"type:text" json:"description"`
	IPAddress   string      `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent   string      `gorm:"type:text" json:"user_agent"`
	Metadata    string      `gorm:"type:text" json:"metadata,omitempty"` // JSON metadata
	CreatedAt   time.Time   `json:"created_at"`

	// Relations
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}
