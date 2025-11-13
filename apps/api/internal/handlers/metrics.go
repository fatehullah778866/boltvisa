package handlers

import (
	"net/http"

	"github.com/boltvisa/api/internal/middleware"
	"github.com/gin-gonic/gin"
)

func (h *Handlers) GetMetrics(c *gin.Context) {
	metrics := middleware.GetMetrics()
	c.JSON(http.StatusOK, gin.H{
		"metrics": metrics,
		"timestamp": gin.H{
			"unix": gin.H{},
		},
	})
}
