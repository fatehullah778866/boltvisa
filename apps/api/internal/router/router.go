package router

import (
    "github.com/boltvisa/api/internal/config"
    "github.com/boltvisa/api/internal/handlers"
    "github.com/boltvisa/api/internal/middleware"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

func Setup(db *gorm.DB, cfg *config.Config) *gin.Engine {
    if cfg.Environment == "production" {
        gin.SetMode(gin.ReleaseMode)
    }

    r := gin.New()
    r.Use(gin.Recovery())

    // Middleware (order matters)
    r.Use(middleware.Logging())
    r.Use(middleware.Metrics())
    r.Use(middleware.RateLimit()) // Rate limiting for all routes
    r.Use(middleware.CORSMiddleware()) // CORS before routes

	// Initialize handlers
	h := handlers.New(db, cfg)

	// Health check (public)
	r.GET("/health", handlers.Health)

	// Metrics endpoint (for monitoring)
	r.GET("/metrics", h.GetMetrics)

	// OpenAPI specification
	r.GET("/openapi.json", h.GetOpenAPISpec)
	r.GET("/swagger.json", h.GetOpenAPISpec) // Alias for compatibility

	// Webhook routes (public, signature verified)
	webhooks := r.Group("/webhooks")
	{
		webhooks.POST("/stripe", h.StripeWebhook)
		webhooks.POST("/razorpay", h.RazorpayWebhook)
	}

	// API routes
	api := r.Group("/api/v1")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", h.Register)
			auth.POST("/login", h.Login)
			auth.POST("/refresh", h.RefreshToken)
			auth.POST("/forgot-password", h.ForgotPassword)
			auth.POST("/reset-password", h.ResetPassword)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.Auth(cfg.JWTSecret))
		protected.Use(middleware.AuthRateLimit()) // Stricter rate limit for authenticated routes
		{
			// Dashboard
			protected.GET("/dashboard", h.GetDashboard)

			// User routes
			protected.GET("/users/me", h.GetCurrentUser)
			protected.PUT("/users/me", h.UpdateCurrentUser)

			// Visa categories (public but may need auth for some operations)
			protected.GET("/visa-categories", h.GetVisaCategories)
			protected.GET("/visa-categories/:id", h.GetVisaCategory)

			// Visa applications
			protected.GET("/applications", h.GetApplications)
			protected.POST("/applications", h.CreateApplication)
			protected.GET("/applications/:id", h.GetApplication)
			protected.PUT("/applications/:id", h.UpdateApplication)
			protected.DELETE("/applications/:id", h.DeleteApplication)

			// Documents
			protected.POST("/applications/:id/documents", h.UploadDocument)
			protected.GET("/applications/:id/documents", h.GetDocuments)
			protected.DELETE("/documents/:id", h.DeleteDocument)

			// Notifications
			protected.GET("/notifications", h.GetNotifications)
			protected.PUT("/notifications/:id/read", h.MarkNotificationRead)
			protected.PUT("/notifications/read-all", h.MarkAllNotificationsRead)

			// Payments
			protected.POST("/payments", h.CreatePayment)
			protected.GET("/payments", h.GetPayments)
			protected.POST("/payments/:id/confirm", h.ConfirmPayment)

			// Admin routes
			admin := protected.Group("/admin")
			admin.Use(middleware.RequireRole("admin"))
			{
				admin.GET("/users", h.ListUsers)
				admin.PUT("/users/:id", h.UpdateUser)
				admin.POST("/visa-categories", h.CreateVisaCategory)
				admin.PUT("/visa-categories/:id", h.UpdateVisaCategory)
				admin.GET("/audit-logs", h.GetAuditLogs)
				admin.GET("/audit-logs/:id", h.GetAuditLog)
			}

			// Consultant routes
			consultant := protected.Group("/consultant")
			consultant.Use(middleware.RequireRole("consultant", "admin"))
			{
				consultant.GET("/clients", h.GetConsultantClients)
				consultant.GET("/applications", h.GetConsultantApplications)
			}

			// Assign consultant (admin/consultant only)
			protected.PUT("/applications/:id/assign-consultant", middleware.RequireRole("admin", "consultant"), h.AssignConsultant)
		}
	}

	return r
}
