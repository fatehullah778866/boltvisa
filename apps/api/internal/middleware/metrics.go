package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	requestCount   = make(map[string]int64)
	requestLatency = make(map[string]time.Duration)
	requestErrors  = make(map[string]int64)
	metricsMutex   sync.RWMutex
)

func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		// Process request
		c.Next()

		// Record metrics
		latency := time.Since(start)
		statusCode := c.Writer.Status()

		metricsMutex.Lock()
		requestCount[path]++
		requestLatency[path] = latency
		if statusCode >= 400 {
			requestErrors[path]++
		}
		metricsMutex.Unlock()
	}
}

func GetMetrics() map[string]interface{} {
	metricsMutex.RLock()
	defer metricsMutex.RUnlock()

	result := make(map[string]interface{})
	for path := range requestCount {
		result[path] = map[string]interface{}{
			"count":   requestCount[path],
			"latency": requestLatency[path].Milliseconds(),
			"errors":  requestErrors[path],
		}
	}
	return result
}

func ResetMetrics() {
	metricsMutex.Lock()
	defer metricsMutex.Unlock()
	requestCount = make(map[string]int64)
	requestLatency = make(map[string]time.Duration)
	requestErrors = make(map[string]int64)
}
