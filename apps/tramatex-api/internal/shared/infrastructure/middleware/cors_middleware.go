package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware configures CORS with whitelist-based origin validation.
// Allowed origins are passed as a comma-separated string (e.g., "http://localhost:3000,https://app.tramatex.com")
// Use "*" to allow all origins (NOT recommended for production).
func CORSMiddleware(allowedOrigins string) gin.HandlerFunc {
	origins := parseOrigins(allowedOrigins)
	allowAll := allowedOrigins == "*"

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Allow all origins if configured
		if allowAll {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		} else if origin != "" && isOriginAllowed(origin, origins) {
			// Allow whitelisted origin
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Vary", "Origin")
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Request-ID, X-User-ID")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400") // 24 hours

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// parseOrigins splits comma-separated origins into a slice
func parseOrigins(originsStr string) []string {
	if originsStr == "" || originsStr == "*" {
		return []string{}
	}

	parts := strings.Split(originsStr, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// isOriginAllowed checks if the origin is in the whitelist
func isOriginAllowed(origin string, allowed []string) bool {
	for _, o := range allowed {
		if o == origin {
			return true
		}
	}
	return false
}
