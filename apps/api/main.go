package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

/*
=========================

	Models
	=========================
*/
type user struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
	Active    bool   `json:"active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type visaCategory struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Country     string `json:"country"`
	Duration    string `json:"duration"`
	Price       int64  `json:"price"`
	Active      bool   `json:"active"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type visaStatus string

const (
	statusDraft     visaStatus = "draft"
	statusSubmitted visaStatus = "submitted"
	statusInReview  visaStatus = "in_review"
	statusApproved  visaStatus = "approved"
	statusRejected  visaStatus = "rejected"
	statusCancelled visaStatus = "cancelled"
)

type visaApplication struct {
	ID           int64         `json:"id"`
	UserID       int64         `json:"user_id"`
	ConsultantID *int64        `json:"consultant_id,omitempty"`
	CategoryID   int64         `json:"category_id"`
	Status       visaStatus    `json:"status"`
	PassportNo   string        `json:"passport_number,omitempty"`
	DateOfBirth  string        `json:"date_of_birth,omitempty"`
	Nationality  string        `json:"nationality,omitempty"`
	TravelDate   string        `json:"travel_date,omitempty"`
	Notes        string        `json:"notes,omitempty"`
	CreatedAt    string        `json:"created_at"`
	UpdatedAt    string        `json:"updated_at"`
	User         *user         `json:"user,omitempty"`
	Consultant   *user         `json:"consultant,omitempty"`
	Category     *visaCategory `json:"category,omitempty"`
}

/*
=========================

	In-memory stores
	=========================
*/
var (
	users            = map[string]user{} // key: email
	nextUserID int64 = 1

	categories       = []visaCategory{}
	nextCatID  int64 = 1

	applications       = []visaApplication{}
	nextAppID    int64 = 1
)

/*
=========================

	DTOs
	=========================
*/
type registerRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}
type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
type authResponse struct {
	Token     string `json:"token"`
	User      user   `json:"user"`
	ExpiresAt string `json:"expires_at"`
}

// NEW: create-application request (dev)
type createApplicationRequest struct {
	CategoryID  int64  `json:"category_id" binding:"required"`
	PassportNo  string `json:"passport_number"`
	DateOfBirth string `json:"date_of_birth"` // ISO-8601 (optional)
	Nationality string `json:"nationality"`
	TravelDate  string `json:"travel_date"` // ISO-8601 (optional)
	Notes       string `json:"notes"`
	Submit      bool   `json:"submit"` // if true -> status=submitted, else draft
}

/*
=========================

	Helpers / seeders
	=========================
*/
func nowISO() string { return time.Now().UTC().Format(time.RFC3339) }

func seedCategoriesOnce() {
	if len(categories) > 0 {
		return
	}
	add := func(name, desc, country, duration string, price int64) {
		categories = append(categories, visaCategory{
			ID:          nextCatID,
			Name:        name,
			Description: desc,
			Country:     country,
			Duration:    duration,
			Price:       price,
			Active:      true,
			CreatedAt:   nowISO(),
			UpdatedAt:   nowISO(),
		})
		nextCatID++
	}
	add("Tourist Visa", "Short stay tourism", "United States", "6 months", 100)
	add("Student Visa", "University enrollment", "United Kingdom", "2 years", 200)
	add("Work Visa", "Skilled worker", "Canada", "3 years", 300)
}

func seedDefaultUserOnce() {
	if len(users) > 0 {
		return
	}
	email := "demo@boltvisa.local"
	users[email] = user{
		ID:        nextUserID,
		Email:     email,
		FirstName: "Demo",
		LastName:  "User",
		Role:      "applicant",
		Active:    true,
		CreatedAt: nowISO(),
		UpdatedAt: nowISO(),
	}
	nextUserID++
}

/*
=========================

	Auth + user helpers
	=========================
*/
func tokenFor(email string) string { return "devtoken-" + strings.ToLower(email) }

func userFromContext(c *gin.Context) (user, bool) {
	val, ok := c.Get("user")
	if !ok {
		return user{}, false
	}
	u, ok := val.(user)
	return u, ok
}

/*
=========================

	Handlers
	=========================
*/
func registerHandler(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	email := strings.ToLower(req.Email)
	if _, exists := users[email]; exists {
		c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
		return
	}
	u := user{
		ID:        nextUserID,
		Email:     email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      "applicant",
		Active:    true,
		CreatedAt: nowISO(),
		UpdatedAt: nowISO(),
	}
	nextUserID++
	users[email] = u

	c.JSON(http.StatusOK, gin.H{"data": authResponse{
		Token:     tokenFor(email),
		User:      u,
		ExpiresAt: time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339),
	}})
}

// DEV: auto-create on login to remove 401 friction
func loginHandler(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	email := strings.ToLower(req.Email)
	u, ok := users[email]
	if !ok {
		u = user{
			ID:        nextUserID,
			Email:     email,
			FirstName: strings.Split(email, "@")[0],
			LastName:  "User",
			Role:      "applicant",
			Active:    true,
			CreatedAt: nowISO(),
			UpdatedAt: nowISO(),
		}
		nextUserID++
		users[email] = u
	}

	c.JSON(http.StatusOK, gin.H{"data": authResponse{
		Token:     tokenFor(email),
		User:      u,
		ExpiresAt: time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339),
	}})
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		raw := strings.TrimPrefix(h, "Bearer ")
		parts := strings.SplitN(raw, "devtoken-", 2)
		if len(parts) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		email := parts[1]
		u, ok := users[email]
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		c.Set("user", u)
		c.Next()
	}
}

func meHandler(c *gin.Context) {
	u, ok := userFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": u})
}

/*
=========================

	Applications & Dashboard
	=========================
*/
func listCategoriesHandler(c *gin.Context) {
	seedCategoriesOnce()
	c.JSON(http.StatusOK, gin.H{"data": categories})
}

func listApplicationsHandler(c *gin.Context) {
	u, ok := userFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	seedCategoriesOnce()

	out := make([]visaApplication, 0, len(applications))
	for _, a := range applications {
		if u.Role == "applicant" && a.UserID != u.ID {
			continue
		}
		appCopy := a
		// populate relations
		appCopy.User = &u
		for i := range categories {
			if categories[i].ID == a.CategoryID {
				cat := categories[i]
				appCopy.Category = &cat
				break
			}
		}
		out = append(out, appCopy)
	}
	c.JSON(http.StatusOK, gin.H{"data": out})
}

// NEW: create application
func createApplicationHandler(c *gin.Context) {
	u, ok := userFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req createApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.CategoryID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	// validate category exists
	var cat *visaCategory
	for i := range categories {
		if categories[i].ID == req.CategoryID {
			cat = &categories[i]
			break
		}
	}
	if cat == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category_id"})
		return
	}

	now := nowISO()
	app := visaApplication{
		ID:          nextAppID,
		UserID:      u.ID,
		CategoryID:  req.CategoryID,
		Status:      statusDraft,
		PassportNo:  req.PassportNo,
		DateOfBirth: req.DateOfBirth,
		Nationality: req.Nationality,
		TravelDate:  req.TravelDate,
		Notes:       req.Notes,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if req.Submit {
		app.Status = statusSubmitted
	}
	nextAppID++
	applications = append(applications, app)

	// populate relations for response
	app.User = &u
	catCopy := *cat
	app.Category = &catCopy

	c.JSON(http.StatusOK, gin.H{"data": app})
}

func dashboardHandler(c *gin.Context) {
	u, ok := userFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	counts := map[visaStatus]int{
		statusDraft: 0, statusSubmitted: 0, statusInReview: 0,
		statusApproved: 0, statusRejected: 0, statusCancelled: 0,
	}
	total := 0
	recent := []visaApplication{}
	for _, a := range applications {
		if u.Role == "applicant" && a.UserID != u.ID {
			continue
		}
		total++
		counts[a.Status]++
		if len(recent) < 5 {
			recent = append(recent, a)
		}
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{
		"stats": gin.H{
			"total_applications": total,
			"draft":              counts[statusDraft], "submitted": counts[statusSubmitted],
			"in_review": counts[statusInReview], "approved": counts[statusApproved],
			"rejected": counts[statusRejected], "cancelled": counts[statusCancelled],
		},
		"recent_applications": recent,
	}})
}

/*
=========================

	Server bootstrap
	=========================
*/
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	seedCategoriesOnce()
	seedDefaultUserOnce()

	r := gin.Default()

	// CORS for local dev
	r.Use(func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "http://127.0.0.1:3000" || origin == "http://localhost:3000" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}
		c.Writer.Header().Set("Vary", "Origin")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	// health
	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })

	// API
	api := r.Group("/api/v1")
	{
		api.POST("/auth/register", registerHandler)
		api.POST("/auth/login", loginHandler)
		api.GET("/users/me", authMiddleware(), meHandler)

		api.GET("/visa-categories", authMiddleware(), listCategoriesHandler)
		api.GET("/applications", authMiddleware(), listApplicationsHandler)
		api.POST("/applications", authMiddleware(), createApplicationHandler) // <-- NEW
		api.GET("/dashboard", authMiddleware(), dashboardHandler)
	}

	srv := &http.Server{
		Addr:           ":" + port,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("boltvisa-api listening on :%s", port)
	log.Fatal(srv.ListenAndServe())
}
