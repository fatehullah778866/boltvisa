package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log request details
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if raw != "" {
			path = path + "?" + raw
		}

		log.Printf("[%s] %s %s %d %v %s %s",
			clientIP,
			method,
			path,
			statusCode,
			latency,
			errorMessage,
			c.Request.UserAgent(),
		)

		// In production, send to Stackdriver/Cloud Logging
		// logEntry := &logging.Entry{
		// 	Severity: logging.Info,
		// 	Payload: map[string]interface{}{
		// 		"method":      method,
		// 		"path":        path,
		// 		"status_code": statusCode,
		// 		"latency_ms":  latency.Milliseconds(),
		// 		"client_ip":   clientIP,
		// 	},
		// }
	}
}
