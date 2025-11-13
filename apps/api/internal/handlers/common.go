package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIError struct {
	Error   string      `json:"error"`
	Details interface{} `json:"details,omitempty"`
}

func WriteErr(c *gin.Context, code int, msg string, details interface{}) {
	c.AbortWithStatusJSON(code, APIError{Error: msg, Details: details})
}

func WriteOK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"data": data})
}
