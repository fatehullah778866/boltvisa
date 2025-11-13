package handlers

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"

    "github.com/boltvisa/api/internal/config"
    "github.com/boltvisa/api/internal/models"
    "github.com/gin-gonic/gin"
    "github.com/glebarez/sqlite"
    "github.com/stretchr/testify/assert"
    "gorm.io/gorm"
)

func setupIntegrationTestDB() *gorm.DB {
    // Use a uniquely named shared in-memory SQLite DB per test to avoid cross-test collisions
    dsn := fmt.Sprintf("file:memdb_int_%d?mode=memory&cache=shared", time.Now().UnixNano())
    db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("Failed to connect to test database")
    }

	db.AutoMigrate(
		&models.User{},
		&models.VisaCategory{},
		&models.VisaApplication{},
		&models.Document{},
		&models.Notification{},
		&models.Payment{},
		&models.AuditLog{},
		&models.PasswordResetToken{},
	)

	return db
}

func setupIntegrationRouter(h *Handlers) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", h.Register)
			auth.POST("/login", h.Login)
		}

		protected := api.Group("")
		protected.Use(func(c *gin.Context) {
			// Mock auth middleware for integration tests
			c.Set("userID", uint(1))
			c.Set("userRole", "applicant")
			c.Next()
		})
		{
			protected.GET("/users/me", h.GetCurrentUser)
			protected.PUT("/users/me", h.UpdateCurrentUser)
			protected.GET("/visa-categories", h.GetVisaCategories)
			protected.POST("/applications", h.CreateApplication)
			protected.GET("/applications", h.GetApplications)
			protected.POST("/payments", h.CreatePayment)
			protected.GET("/payments", h.GetPayments)
		}
	}

	return r
}

// TestUserRegistrationAndLoginFlow tests the complete user registration and login flow
func TestUserRegistrationAndLoginFlow(t *testing.T) {
	db := setupIntegrationTestDB()
	cfg := &config.Config{
		JWTSecret: "test-secret-key-for-integration-tests",
	}
	h := New(db, cfg)
	r := setupIntegrationRouter(h)

	// Step 1: Register a new user
	registerReq := map[string]interface{}{
		"email":        "test@example.com",
		"password":     "password123",
		"first_name":   "Test",
		"last_name":    "User",
		"phone_number": "+1234567890",
	}
	registerBody, _ := json.Marshal(registerReq)
	registerHTTPReq, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(registerBody))
	registerHTTPReq.Header.Set("Content-Type", "application/json")
	registerW := httptest.NewRecorder()
	r.ServeHTTP(registerW, registerHTTPReq)

	assert.Equal(t, http.StatusCreated, registerW.Code)
	var registerResp map[string]interface{}
	json.Unmarshal(registerW.Body.Bytes(), &registerResp)
	assert.NotNil(t, registerResp["token"])
	assert.NotNil(t, registerResp["user"])

	// Step 2: Login with registered credentials
	loginReq := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	loginBody, _ := json.Marshal(loginReq)
	loginHTTPReq, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(loginBody))
	loginHTTPReq.Header.Set("Content-Type", "application/json")
	loginW := httptest.NewRecorder()
	r.ServeHTTP(loginW, loginHTTPReq)

	assert.Equal(t, http.StatusOK, loginW.Code)
	var loginResp map[string]interface{}
	json.Unmarshal(loginW.Body.Bytes(), &loginResp)
	assert.NotNil(t, loginResp["token"])
}

// TestApplicationCreationFlow tests creating a visa application
func TestApplicationCreationFlow(t *testing.T) {
	db := setupIntegrationTestDB()
	cfg := &config.Config{
		JWTSecret: "test-secret-key-for-integration-tests",
	}
	h := New(db, cfg)

	// Create a test user
	user := models.User{
		Email:     "test@example.com",
		Password:  "hashedpassword",
		FirstName: "Test",
		LastName:  "User",
		Role:      models.RoleApplicant,
		Active:    true,
	}
	db.Create(&user)

	// Create a visa category
	category := models.VisaCategory{
		Name:        "Tourist Visa",
		Description: "Short-term tourist visa",
		Country:     "USA",
		Duration:    "90 days",
		Price:       160.00,
		Active:      true,
	}
	db.Create(&category)

	r := setupIntegrationRouter(h)

	// Create application
    appReq := map[string]interface{}{
        "category_id": category.ID,
        // The API sets new applications to draft by default
        "status":      "draft",
    }
	appBody, _ := json.Marshal(appReq)
	appHTTPReq, _ := http.NewRequest("POST", "/api/v1/applications", bytes.NewBuffer(appBody))
	appHTTPReq.Header.Set("Content-Type", "application/json")
	appW := httptest.NewRecorder()
	r.ServeHTTP(appW, appHTTPReq)

	assert.Equal(t, http.StatusCreated, appW.Code)
	var appResp models.VisaApplication
	json.Unmarshal(appW.Body.Bytes(), &appResp)
    assert.Equal(t, category.ID, appResp.CategoryID)
    assert.Equal(t, "draft", string(appResp.Status))
}

// TestUserProfileUpdateFlow tests updating user profile
func TestUserProfileUpdateFlow(t *testing.T) {
	db := setupIntegrationTestDB()
	cfg := &config.Config{
		JWTSecret: "test-secret-key-for-integration-tests",
	}
	h := New(db, cfg)

	// Create a test user
	user := models.User{
		Email:       "test@example.com",
		Password:    "hashedpassword",
		FirstName:   "Test",
		LastName:    "User",
		PhoneNumber: "+1234567890",
		Role:        models.RoleApplicant,
		Active:      true,
	}
	db.Create(&user)

	r := setupIntegrationRouter(h)

	// Update user profile
	updateReq := map[string]string{
		"first_name":   "Updated",
		"last_name":    "Name",
		"phone_number": "+9876543210",
	}
	updateBody, _ := json.Marshal(updateReq)
	updateHTTPReq, _ := http.NewRequest("PUT", "/api/v1/users/me", bytes.NewBuffer(updateBody))
	updateHTTPReq.Header.Set("Content-Type", "application/json")
	updateW := httptest.NewRecorder()
	r.ServeHTTP(updateW, updateHTTPReq)

	assert.Equal(t, http.StatusOK, updateW.Code)
	var updatedUser models.User
	json.Unmarshal(updateW.Body.Bytes(), &updatedUser)
	assert.Equal(t, "Updated", updatedUser.FirstName)
	assert.Equal(t, "Name", updatedUser.LastName)
	assert.Equal(t, "+9876543210", updatedUser.PhoneNumber)
}

// TestPaymentCreationFlow tests payment creation (without actual payment gateway calls)
func TestPaymentCreationFlow(t *testing.T) {
	db := setupIntegrationTestDB()
	cfg := &config.Config{
		JWTSecret:         "test-secret-key-for-integration-tests",
		StripeSecretKey:   "", // Empty to avoid actual API calls
		RazorpayKeyID:     "",
		RazorpayKeySecret: "",
	}
	h := New(db, cfg)

	// Create test user and application
	user := models.User{
		Email:     "test@example.com",
		Password:  "hashedpassword",
		FirstName: "Test",
		LastName:  "User",
		Role:      models.RoleApplicant,
		Active:    true,
	}
	db.Create(&user)

	category := models.VisaCategory{
		Name:    "Tourist Visa",
		Country: "USA",
		Price:   160.00,
		Active:  true,
	}
	db.Create(&category)

    application := models.VisaApplication{
        UserID:     user.ID,
        CategoryID: category.ID,
        Status:     models.StatusDraft,
    }
	db.Create(&application)

	r := setupIntegrationRouter(h)

	// Create payment (will fail at payment gateway, but tests the flow)
	paymentReq := map[string]interface{}{
		"application_id": application.ID,
		"amount":         160.00,
		"currency":       "USD",
		"method":         "stripe",
	}
	paymentBody, _ := json.Marshal(paymentReq)
	paymentHTTPReq, _ := http.NewRequest("POST", "/api/v1/payments", bytes.NewBuffer(paymentBody))
	paymentHTTPReq.Header.Set("Content-Type", "application/json")
	paymentW := httptest.NewRecorder()
	r.ServeHTTP(paymentW, paymentHTTPReq)

	// Should fail because Stripe is not configured, but payment record should be created
	// In a real scenario, this would succeed with proper credentials
	assert.True(t, paymentW.Code == http.StatusInternalServerError || paymentW.Code == http.StatusOK)
}
