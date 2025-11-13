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
    "github.com/boltvisa/api/internal/utils"
    "github.com/gin-gonic/gin"
    "github.com/glebarez/sqlite"
    "github.com/stretchr/testify/assert"
    "gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
    // Use a uniquely named shared in-memory SQLite DB per test to avoid cross-test collisions
    dsn := fmt.Sprintf("file:memdb_%d?mode=memory&cache=shared", time.Now().UnixNano())
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

func setupRouter(h *Handlers) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	api := r.Group("/api/v1")
	{
		api.POST("/auth/register", h.Register)
		api.POST("/auth/login", h.Login)
		api.GET("/users/me", h.GetCurrentUser)
		api.GET("/visa-categories", h.GetVisaCategories)
		api.POST("/applications", h.CreateApplication)
		api.GET("/applications", h.GetApplications)
	}

	return r
}

func TestRegister(t *testing.T) {
	db := setupTestDB()
	cfg := &config.Config{JWTSecret: "test-secret"}
	h := New(db, cfg)
	r := setupRouter(h)

	reqBody := map[string]string{
		"email":      "test@example.com",
		"password":   "password123",
		"first_name": "Test",
		"last_name":  "User",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.NotNil(t, response["token"])
	assert.NotNil(t, response["user"])
}

func TestLogin(t *testing.T) {
    db := setupTestDB()
    cfg := &config.Config{JWTSecret: "test-secret"}
    h := New(db, cfg)

    // Create a user first
    hashed, _ := utils.HashPassword("password123")
    user := models.User{
        Email:     "test@example.com",
        Password:  hashed, // bcrypt hash of password123
        FirstName: "Test",
        LastName:  "User",
        Role:      models.RoleApplicant,
        Active:    true,
    }
	db.Create(&user)

	r := setupRouter(h)
	reqBody := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.NotNil(t, response["token"])
}

func TestGetVisaCategories(t *testing.T) {
	db := setupTestDB()
	cfg := &config.Config{JWTSecret: "test-secret"}
	h := New(db, cfg)

	// Create test categories
	category1 := models.VisaCategory{
		Name:        "Tourist Visa",
		Description: "Short-term tourist visa",
		Country:     "USA",
		Duration:    "90 days",
		Price:       160.00,
		Active:      true,
	}
	category2 := models.VisaCategory{
		Name:        "Business Visa",
		Description: "Business travel visa",
		Country:     "USA",
		Duration:    "1 year",
		Price:       200.00,
		Active:      true,
	}
	db.Create(&category1)
	db.Create(&category2)

	r := setupRouter(h)
	req, _ := http.NewRequest("GET", "/api/v1/visa-categories", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var categories []models.VisaCategory
	json.Unmarshal(w.Body.Bytes(), &categories)
	assert.GreaterOrEqual(t, len(categories), 2)
}
