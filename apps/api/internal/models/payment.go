package models

import (
	"time"
)

type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "pending"
	PaymentStatusProcessing PaymentStatus = "processing"
	PaymentStatusCompleted  PaymentStatus = "completed"
	PaymentStatusFailed     PaymentStatus = "failed"
	PaymentStatusRefunded   PaymentStatus = "refunded"
)

type PaymentMethod string

const (
	PaymentMethodStripe   PaymentMethod = "stripe"
	PaymentMethodRazorpay PaymentMethod = "razorpay"
)

type Payment struct {
	ID              uint          `gorm:"primaryKey" json:"id"`
	ApplicationID   uint          `gorm:"not null;index" json:"application_id"`
	UserID          uint          `gorm:"not null;index" json:"user_id"`
	Amount          float64       `gorm:"not null" json:"amount"`
	Currency        string        `gorm:"default:'USD'" json:"currency"`
	Status          PaymentStatus `gorm:"type:varchar(20);default:'pending'" json:"status"`
	Method          PaymentMethod `gorm:"type:varchar(20)" json:"method"`
	TransactionID   string        `gorm:"index" json:"transaction_id"`
	PaymentIntentID string        `json:"payment_intent_id,omitempty"`         // Stripe payment intent ID
	RazorpayOrderID string        `json:"razorpay_order_id,omitempty"`         // Razorpay order ID
	Metadata        string        `gorm:"type:text" json:"metadata,omitempty"` // JSON metadata
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`

	// Relations
	Application VisaApplication `gorm:"foreignKey:ApplicationID" json:"application,omitempty"`
	User        User            `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (Payment) TableName() string {
	return "payments"
}
