package handlers

import (
	"fmt"
	"net/http"

	"github.com/boltvisa/api/internal/models"
	"github.com/boltvisa/api/internal/services"
	"github.com/gin-gonic/gin"
)

type CreatePaymentRequest struct {
	ApplicationID uint    `json:"application_id" binding:"required"`
	Amount        float64 `json:"amount" binding:"required"`
	Currency      string  `json:"currency"`
	Method        string  `json:"method" binding:"required"` // "stripe" or "razorpay"
}

func (h *Handlers) CreatePayment(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify application exists and belongs to user
	var application models.VisaApplication
	if err := h.db.First(&application, req.ApplicationID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	userRole, _ := c.Get("userRole")
	if userRole == "applicant" && application.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if req.Currency == "" {
		req.Currency = "USD"
	}

	// Create payment record
	payment := models.Payment{
		ApplicationID: req.ApplicationID,
		UserID:        userID.(uint),
		Amount:        req.Amount,
		Currency:      req.Currency,
		Status:        models.PaymentStatusPending,
		Method:        models.PaymentMethod(req.Method),
	}

	if err := h.db.Create(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment"})
		return
	}

	// Create payment intent based on method
	var paymentIntentID string
	var clientSecret string
	var razorpayOrderID string

	if req.Method == "stripe" {
		paymentService := services.NewPaymentService(h.cfg.StripeSecretKey, "", "")
		pi, err := paymentService.CreateStripePaymentIntent(c.Request.Context(), req.Amount, req.Currency, map[string]string{
			"payment_id":     fmt.Sprintf("%d", payment.ID),
			"application_id": fmt.Sprintf("%d", req.ApplicationID),
			"user_id":        fmt.Sprintf("%d", userID.(uint)),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment intent"})
			return
		}
		paymentIntentID = pi.ID
		clientSecret = pi.ClientSecret
		payment.PaymentIntentID = pi.ID
		payment.Status = models.PaymentStatusProcessing
		h.db.Save(&payment)
	} else if req.Method == "razorpay" {
		paymentService := services.NewPaymentService("", h.cfg.RazorpayKeyID, h.cfg.RazorpayKeySecret)
		order, err := paymentService.CreateRazorpayOrder(c.Request.Context(), req.Amount, req.Currency, map[string]string{
			"payment_id":     fmt.Sprintf("%d", payment.ID),
			"application_id": fmt.Sprintf("%d", req.ApplicationID),
			"user_id":        fmt.Sprintf("%d", userID.(uint)),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Razorpay order"})
			return
		}
		razorpayOrderID = order.ID
		payment.RazorpayOrderID = order.ID
		payment.Status = models.PaymentStatusProcessing
		h.db.Save(&payment)
	}

	response := gin.H{
		"payment_id": payment.ID,
		"amount":     payment.Amount,
		"currency":   payment.Currency,
	}

	if paymentIntentID != "" {
		response["payment_intent_id"] = paymentIntentID
		response["client_secret"] = clientSecret
	}
	if razorpayOrderID != "" {
		response["razorpay_order_id"] = razorpayOrderID
		response["razorpay_key_id"] = h.cfg.RazorpayKeyID
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handlers) ConfirmPayment(c *gin.Context) {
	paymentID := c.Param("id")
	userID, _ := c.Get("userID")

	var payment models.Payment
	if err := h.db.First(&payment, paymentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	if payment.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var req struct {
		PaymentIntentID string `json:"payment_intent_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify payment based on method
	paymentService := services.NewPaymentService(h.cfg.StripeSecretKey, h.cfg.RazorpayKeyID, h.cfg.RazorpayKeySecret)

	if payment.Method == models.PaymentMethodStripe && req.PaymentIntentID != "" {
		pi, err := paymentService.ConfirmStripePayment(c.Request.Context(), req.PaymentIntentID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Payment verification failed"})
			return
		}

		if pi.Status == "succeeded" {
			payment.Status = models.PaymentStatusCompleted
			payment.TransactionID = pi.ID
		} else {
			payment.Status = models.PaymentStatusFailed
		}
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
	} else if payment.Method == models.PaymentMethodRazorpay {
		var razorpayReq struct {
			RazorpayOrderID   string `json:"razorpay_order_id"`
			RazorpayPaymentID string `json:"razorpay_payment_id"`
			RazorpaySignature string `json:"razorpay_signature"`
		}
		if err := c.ShouldBindJSON(&razorpayReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		valid, err := paymentService.VerifyRazorpayPayment(c.Request.Context(),
			razorpayReq.RazorpayOrderID, razorpayReq.RazorpayPaymentID, razorpayReq.RazorpaySignature)
		if err != nil || !valid {
			payment.Status = models.PaymentStatusFailed
			h.db.Save(&payment)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Payment verification failed"})
			return
		}

		payment.Status = models.PaymentStatusCompleted
		payment.TransactionID = razorpayReq.RazorpayPaymentID
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
	}

	c.JSON(http.StatusOK, payment)
}

func (h *Handlers) GetPayments(c *gin.Context) {
	userID, _ := c.Get("userID")
	userRole, _ := c.Get("userRole")

	query := h.db.Model(&models.Payment{}).Preload("Application")

	if userRole == "applicant" {
		query = query.Where("user_id = ?", userID)
	}

	var payments []models.Payment
	if err := query.Find(&payments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch payments"})
		return
	}

	c.JSON(http.StatusOK, payments)
}
