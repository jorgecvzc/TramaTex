package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	sharedomain "github.com/joran-cortez/tramatex/internal/shared/domain"
	"github.com/joran-cortez/tramatex/internal/shared/infrastructure/logging"
)

// ErrorHandlerMiddleware sanitizes error responses based on environment.
// Production: Generic error messages (avoid leaking implementation details)
// Development: Detailed error messages (for debugging)
//
// Domain errors implementing shared/domain.HTTPStatuser automatically determine
// the HTTP status code. Untyped errors fall back to 500.
func ErrorHandlerMiddleware(environment string) gin.HandlerFunc {
	isProduction := environment == "production"

	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		ginErr := c.Errors.Last()

		// Get request ID for log correlation
		requestID, _ := c.Get("requestID")
		reqID, _ := requestID.(string)

		logging.WithRequestID(reqID).WithFields(map[string]interface{}{
			"path":   c.Request.URL.Path,
			"method": c.Request.Method,
			"error":  ginErr.Error(),
		}).Error("Request processing error")

		// Determine HTTP status from the error type when available
		status := http.StatusInternalServerError
		var statuser sharedomain.HTTPStatuser
		if errors.As(ginErr.Err, &statuser) {
			status = statuser.HTTPStatus()
		}

		if isProduction {
			c.JSON(status, gin.H{
				"error":      "An error occurred while processing your request",
				"request_id": reqID,
			})
		} else {
			c.JSON(status, gin.H{
				"error":      ginErr.Error(),
				"request_id": reqID,
			})
		}
	}
}
