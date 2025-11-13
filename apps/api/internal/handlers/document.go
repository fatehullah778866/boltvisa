package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/boltvisa/api/internal/models"
	"github.com/boltvisa/api/internal/storage"
	"github.com/gin-gonic/gin"
)

func (h *Handlers) UploadDocument(c *gin.Context) {
	applicationID := c.Param("id")
	userID, _ := c.Get("userID")

	// Verify application exists and user has access
	var application models.VisaApplication
	if err := h.db.First(&application, applicationID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	userRole, _ := c.Get("userRole")
	if userRole == "applicant" && application.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	docType := c.PostForm("type")
	if docType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Document type is required"})
		return
	}

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer src.Close()

	// Upload to GCS if configured
	var gcsURL string
	if h.cfg.GCSBucketName != "" {
		ctx := context.Background()
		gcsClient, err := storage.NewGCSClient(ctx, h.cfg.GCSBucketName, os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize storage client"})
			return
		}
		defer gcsClient.Close()

		objectPath := storage.GenerateDocumentPath(application.ID, file.Filename)
		gcsURL, err = gcsClient.UploadFile(ctx, src, objectPath, file.Header.Get("Content-Type"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to upload file: %v", err)})
			return
		}
	} else {
		// Fallback URL if GCS not configured
		gcsURL = fmt.Sprintf("https://storage.googleapis.com/%s/documents/%d/%s",
			h.cfg.GCSBucketName, application.ID, file.Filename)
	}

	document := models.Document{
		ApplicationID: application.ID,
		Type:          models.DocumentType(docType),
		Name:          file.Filename,
		FileName:      file.Filename,
		GCSURL:        gcsURL,
		Size:          file.Size,
		MimeType:      file.Header.Get("Content-Type"),
		UploadedBy:    userID.(uint),
	}

	if err := h.db.Create(&document).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save document"})
		return
	}

	c.JSON(http.StatusCreated, document)
}

func (h *Handlers) GetDocuments(c *gin.Context) {
	applicationID := c.Param("id")
	userID, _ := c.Get("userID")

	// Verify application exists and user has access
	var application models.VisaApplication
	if err := h.db.First(&application, applicationID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	userRole, _ := c.Get("userRole")
	if userRole == "applicant" && application.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var documents []models.Document
	if err := h.db.Where("application_id = ?", applicationID).Find(&documents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch documents"})
		return
	}

	c.JSON(http.StatusOK, documents)
}

func (h *Handlers) DeleteDocument(c *gin.Context) {
	documentID := c.Param("id")
	userID, _ := c.Get("userID")

	var document models.Document
	if err := h.db.Preload("Application").First(&document, documentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	userRole, _ := c.Get("userRole")
	if userRole == "applicant" && document.Application.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Delete from GCS if configured
	if h.cfg.GCSBucketName != "" && document.GCSURL != "" {
		ctx := context.Background()
		gcsClient, err := storage.NewGCSClient(ctx, h.cfg.GCSBucketName, os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
		if err == nil {
			defer gcsClient.Close()
			// Extract object path from URL
			objectPath := document.GCSURL[len(fmt.Sprintf("https://storage.googleapis.com/%s/", h.cfg.GCSBucketName)):]
			gcsClient.DeleteFile(ctx, objectPath)
		}
	}

	if err := h.db.Delete(&document).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete document"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document deleted"})
}
