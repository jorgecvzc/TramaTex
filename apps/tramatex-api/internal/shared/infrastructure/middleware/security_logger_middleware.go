package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joran-cortez/tramatex/internal/shared/infrastructure/logging"
)

// SecurityLoggerMiddleware logs HTTP requests with security context.
// It captures request details, user info, and response status for audit trails.
func SecurityLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Record start time
		start := time.Now()

		// Get request ID from context (set by RequestIDMiddleware)
		requestID, _ := c.Get("requestID")
		reqID, _ := requestID.(string)

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get user info from context (set by AuthMiddleware)
		userID, _ := c.Get("userID")
		userEmail, _ := c.Get("userEmail")
		userRole, _ := c.Get("userRole")

		// Build log entry with all context
		logEntry := logging.WithRequestID(reqID)

		// Add user info if authenticated
		if userID != nil {
			logEntry = logEntry.WithFields(map[string]interface{}{
				"userID":    userID,
				"userEmail": logging.MaskEmail(userEmail.(string)),
				"userRole":  userRole,
			})
		}

		// Add HTTP details
		logEntry = logEntry.WithFields(map[string]interface{}{
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"status":     c.Writer.Status(),
			"latency_ms": latency.Milliseconds(),
			"ip":         c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		})

		// Log with appropriate level and message based on status code and context
		status := c.Writer.Status()
		msg := "HTTP request completed"
		if authErr, exists := c.Get("authError"); exists {
			if errStr, ok := authErr.(string); ok {
				msg = errStr
			}
		}

		switch {
		case status >= 500:
			logEntry.Error(msg)
		case status >= 400:
			logEntry.Warn(msg)
		default:
			logEntry.Info(msg)
		}
	}
}
