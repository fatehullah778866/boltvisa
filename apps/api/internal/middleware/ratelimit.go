package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// In-memory rate limiter (can be replaced with Redis for distributed systems)
type rateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	rate     int           // requests per duration
	duration time.Duration // time window
}

type visitor struct {
	windowStart time.Time
	count       int
}

var limiter = &rateLimiter{
	visitors: make(map[string]*visitor),
	rate:     100,              // 100 requests
	duration: time.Minute * 15, // per 15 minutes
}

// Cleanup old visitors every hour
func init() {
	go func() {
		for {
			time.Sleep(time.Hour)
			limiter.mu.Lock()
			now := time.Now()
			for ip, v := range limiter.visitors {
				if now.Sub(v.windowStart) > 3*time.Hour {
					delete(limiter.visitors, ip)
				}
			}
			limiter.mu.Unlock()
		}
	}()
}

func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter.mu.Lock()
		defer limiter.mu.Unlock()

		now := time.Now()
		v, exists := limiter.visitors[ip]
		if !exists {
			limiter.visitors[ip] = &visitor{
				windowStart: now,
				count:       1,
			}
			c.Next()
			return
		}

		// Reset counter if time window has passed
		if now.Sub(v.windowStart) > limiter.duration {
			v.count = 1
			v.windowStart = now
			c.Next()
			return
		}

		// Check if limit exceeded
		if v.count >= limiter.rate {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": fmt.Sprintf("Rate limit exceeded. Maximum %d requests per %v", limiter.rate, limiter.duration),
			})
			c.Abort()
			return
		}

		v.count++
		c.Next()
	}
}

// AuthRateLimit - Stricter rate limit for authenticated endpoints
func AuthRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.Next()
			return
		}

		key := fmt.Sprintf("user:%v", userID)
		limiter.mu.Lock()
		defer limiter.mu.Unlock()

		now := time.Now()
		v, exists := limiter.visitors[key]
		if !exists {
			limiter.visitors[key] = &visitor{
				windowStart: now,
				count:       1,
			}
			c.Next()
			return
		}

		// Stricter limit for authenticated users: 200 requests per 15 minutes
		authRate := 200
		if now.Sub(v.windowStart) > limiter.duration {
			v.count = 1
			v.windowStart = now
			c.Next()
			return
		}

		if v.count >= authRate {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": fmt.Sprintf("Rate limit exceeded. Maximum %d requests per %v", authRate, limiter.duration),
			})
			c.Abort()
			return
		}

		v.count++
		c.Next()
	}
}
