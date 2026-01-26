package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestIDMiddleware generates a unique ID for each request and adds it to the context.
// This ID is used for log correlation across the application.
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if request ID was provided in header (for distributed tracing)
		requestID := c.GetHeader("X-Request-ID")

		// Generate new UUID if not provided
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Add to context for handler access
		c.Set("requestID", requestID)

		// Add to response headers for client-side correlation
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}
