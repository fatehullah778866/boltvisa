package services

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentintent"
)

type PaymentService struct {
	stripeSecretKey   string
	razorpayKeyID     string
	razorpayKeySecret string
}

type RazorpayOrder struct {
	ID       string `json:"id"`
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
	Status   string `json:"status"`
}

type RazorpayOrderResponse struct {
	ID       string `json:"id"`
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
	Status   string `json:"status"`
}

func NewPaymentService(stripeSecretKey, razorpayKeyID, razorpayKeySecret string) *PaymentService {
	if stripeSecretKey != "" {
		stripe.Key = stripeSecretKey
	}
	return &PaymentService{
		stripeSecretKey:   stripeSecretKey,
		razorpayKeyID:     razorpayKeyID,
		razorpayKeySecret: razorpayKeySecret,
	}
}

func (s *PaymentService) CreateStripePaymentIntent(ctx context.Context, amount float64, currency string, metadata map[string]string) (*stripe.PaymentIntent, error) {
	if s.stripeSecretKey == "" {
		return nil, fmt.Errorf("Stripe not configured")
	}

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(amount * 100)), // Convert to cents
		Currency: stripe.String(currency),
		Metadata: metadata,
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment intent: %w", err)
	}

	return pi, nil
}

func (s *PaymentService) ConfirmStripePayment(ctx context.Context, paymentIntentID string) (*stripe.PaymentIntent, error) {
	if s.stripeSecretKey == "" {
		return nil, fmt.Errorf("Stripe not configured")
	}

	pi, err := paymentintent.Get(paymentIntentID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment intent: %w", err)
	}

	return pi, nil
}

func (s *PaymentService) RefundStripePayment(ctx context.Context, paymentIntentID string) error {
	if s.stripeSecretKey == "" {
		return fmt.Errorf("Stripe not configured")
	}

	// Implementation would use Stripe refund API
	// For now, return nil as placeholder
	return nil
}

func (s *PaymentService) CreateRazorpayOrder(ctx context.Context, amount float64, currency string, metadata map[string]string) (*RazorpayOrder, error) {
	if s.razorpayKeyID == "" || s.razorpayKeySecret == "" {
		return nil, fmt.Errorf("Razorpay not configured")
	}

	// Convert amount to paise (smallest currency unit for INR) or cents for other currencies
	amountInSmallestUnit := int64(amount * 100)
	if currency == "INR" {
		amountInSmallestUnit = int64(amount * 100) // Already in paise
	}

	orderData := map[string]interface{}{
		"amount":   amountInSmallestUnit,
		"currency": currency,
		"notes":    metadata,
	}

	jsonData, err := json.Marshal(orderData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal order data: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.razorpay.com/v1/orders", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Basic auth with Razorpay credentials
	auth := base64.StdEncoding.EncodeToString([]byte(s.razorpayKeyID + ":" + s.razorpayKeySecret))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create Razorpay order: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Razorpay API error: %s", string(body))
	}

	var orderResponse RazorpayOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&orderResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &RazorpayOrder{
		ID:       orderResponse.ID,
		Amount:   orderResponse.Amount,
		Currency: orderResponse.Currency,
		Status:   orderResponse.Status,
	}, nil
}

func (s *PaymentService) VerifyRazorpayPayment(ctx context.Context, orderID, paymentID, signature string) (bool, error) {
	if s.razorpayKeySecret == "" {
		return false, fmt.Errorf("Razorpay not configured")
	}

	// Verify signature using Razorpay's signature verification
	// Signature format: HMAC SHA256(order_id + "|" + payment_id, secret)
	message := orderID + "|" + paymentID
	expectedSignature := s.GenerateHMAC(message, s.razorpayKeySecret)

	return expectedSignature == signature, nil
}

func (s *PaymentService) GenerateHMAC(message, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}
