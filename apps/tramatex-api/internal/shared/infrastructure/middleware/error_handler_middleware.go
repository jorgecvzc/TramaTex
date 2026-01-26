package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/joran-cortez/tramatex/internal/shared/infrastructure/logging"
)

// ErrorHandlerMiddleware sanitizes error responses based on environment.
// Production: Generic error messages (avoid leaking implementation details)
// Development: Detailed error messages (for debugging)
func ErrorHandlerMiddleware(environment string) gin.HandlerFunc {
	isProduction := environment == "production"

	return func(c *gin.Context) {
		c.Next()

		// Check if there are errors after processing the request
		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			// Get request ID for correlation
			requestID, _ := c.Get("requestID")
			reqID, _ := requestID.(string)

			// Log the actual error with full details
			logging.WithRequestID(reqID).WithFields(map[string]interface{}{
				"path":   c.Request.URL.Path,
				"method": c.Request.Method,
				"error":  err.Error(),
			}).Error("Request processing error")

			// Return sanitized error to client
			if isProduction {
				// Generic error in production
				c.JSON(c.Writer.Status(), gin.H{
					"error":      "An error occurred while processing your request",
					"request_id": reqID,
				})
			} else {
				// Detailed error in development
				c.JSON(c.Writer.Status(), gin.H{
					"error":      err.Error(),
					"request_id": reqID,
				})
			}
		}
	}
}
