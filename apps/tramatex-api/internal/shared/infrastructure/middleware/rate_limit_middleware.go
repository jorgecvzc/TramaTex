package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type rateLimitEntry struct {
	count     int
	resetTime time.Time
}

// RateLimitMiddleware applies a fixed-window rate limit per IP.
// Example: RateLimitMiddleware(10, time.Minute)
func RateLimitMiddleware(maxRequests int, window time.Duration) gin.HandlerFunc {
	var mu sync.Mutex
	entries := make(map[string]*rateLimitEntry)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()

		mu.Lock()
		entry, exists := entries[ip]
		if !exists || now.After(entry.resetTime) {
			entries[ip] = &rateLimitEntry{count: 1, resetTime: now.Add(window)}
			mu.Unlock()
			c.Next()
			return
		}

		if entry.count >= maxRequests {
			mu.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			return
		}

		entry.count++
		mu.Unlock()
		c.Next()
	}
}
