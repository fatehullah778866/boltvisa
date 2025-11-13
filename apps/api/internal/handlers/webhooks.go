package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/boltvisa/api/internal/models"
	"github.com/boltvisa/api/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/webhook"
)

// StripeWebhook handles Stripe webhook events
func (h *Handlers) StripeWebhook(c *gin.Context) {
	const MaxBodyBytes = int64(65536)
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error reading request body"})
		return
	}

	// Verify webhook signature
	event, err := webhook.ConstructEvent(body, c.GetHeader("Stripe-Signature"), h.cfg.StripeWebhookSecret)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Webhook signature verification failed"})
		return
	}

	// Handle different event types
	switch event.Type {
	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing payment intent"})
			return
		}

		// Find payment by payment intent ID
		var payment models.Payment
		if err := h.db.Where("payment_intent_id = ?", paymentIntent.ID).First(&payment).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
			return
		}

		// Update payment status
		payment.Status = models.PaymentStatusCompleted
		payment.TransactionID = paymentIntent.ID
		h.db.Save(&payment)

		// Log audit event (non-blocking)
		if h.auditService != nil {
			go func() {
				_ = h.auditService.LogPayment(c.Request.Context(), payment.UserID, payment.ID,
					payment.Amount, string(payment.Status), c.ClientIP(), c.Request.UserAgent())
			}()
		}

		// Send notification (non-blocking)
		if h.notificationService != nil {
			go func() {
				_ = h.notificationService.SendPaymentNotification(c.Request.Context(), payment.UserID, payment.Amount, string(payment.Status))
			}()
		}

	case "payment_intent.payment_failed":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing payment intent"})
			return
		}

		// Find payment by payment intent ID
		var payment models.Payment
		if err := h.db.Where("payment_intent_id = ?", paymentIntent.ID).First(&payment).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
			return
		}

		// Update payment status
		payment.Status = models.PaymentStatusFailed
		h.db.Save(&payment)

		// Log audit event (non-blocking)
		if h.auditService != nil {
			go func() {
				_ = h.auditService.LogPayment(c.Request.Context(), payment.UserID, payment.ID,
					payment.Amount, string(payment.Status), c.ClientIP(), c.Request.UserAgent())
			}()
		}
	}

	c.JSON(http.StatusOK, gin.H{"received": true})
}

// RazorpayWebhook handles Razorpay webhook events
func (h *Handlers) RazorpayWebhook(c *gin.Context) {
	var webhookData map[string]interface{}
	if err := c.ShouldBindJSON(&webhookData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify webhook signature
	razorpaySignature := c.GetHeader("X-Razorpay-Signature")
	if razorpaySignature == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing signature"})
		return
	}

	// Get the raw body for signature verification
	bodyBytes, _ := json.Marshal(webhookData)
	bodyString := string(bodyBytes)

	// Verify signature (Razorpay uses HMAC SHA256)
	// Note: Razorpay webhook signature is calculated differently - using the webhook secret
	paymentService := services.NewPaymentService("", h.cfg.RazorpayKeyID, h.cfg.RazorpayKeySecret)
	expectedSignature := paymentService.GenerateHMAC(bodyString, h.cfg.RazorpayWebhookSecret)

	if expectedSignature != razorpaySignature {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid signature"})
		return
	}

	// Handle different event types
	event, ok := webhookData["event"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event type"})
		return
	}

	payload, ok := webhookData["payload"].(map[string]interface{})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	switch event {
	case "payment.captured":
		paymentEntity, ok := payload["payment"].(map[string]interface{})
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment entity"})
			return
		}
		orderID, _ := paymentEntity["order_id"].(string)
		paymentID, _ := paymentEntity["id"].(string)

		// Find payment by Razorpay order ID
		var payment models.Payment
		if err := h.db.Where("razorpay_order_id = ?", orderID).First(&payment).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
			return
		}

		// Update payment status
		payment.Status = models.PaymentStatusCompleted
		payment.TransactionID = paymentID
		h.db.Save(&payment)

		// Log audit event (non-blocking)
		if h.auditService != nil {
			go func() {
				_ = h.auditService.LogPayment(c.Request.Context(), payment.UserID, payment.ID,
					payment.Amount, string(payment.Status), c.ClientIP(), c.Request.UserAgent())
			}()
		}

		// Send notification (non-blocking)
		if h.notificationService != nil {
			go func() {
				_ = h.notificationService.SendPaymentNotification(c.Request.Context(), payment.UserID, payment.Amount, string(payment.Status))
			}()
		}

	case "payment.failed":
		paymentEntity, ok := payload["payment"].(map[string]interface{})
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment entity"})
			return
		}
		orderID, _ := paymentEntity["order_id"].(string)

		// Find payment by Razorpay order ID
		var payment models.Payment
		if err := h.db.Where("razorpay_order_id = ?", orderID).First(&payment).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
			return
		}

		// Update payment status
		payment.Status = models.PaymentStatusFailed
		h.db.Save(&payment)

		// Log audit event (non-blocking)
		if h.auditService != nil {
			go func() {
				_ = h.auditService.LogPayment(c.Request.Context(), payment.UserID, payment.ID,
					payment.Amount, string(payment.Status), c.ClientIP(), c.Request.UserAgent())
			}()
		}
	}

	c.JSON(http.StatusOK, gin.H{"received": true})
}
