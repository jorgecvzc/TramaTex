package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecurityHeadersMiddleware adds HTTP security headers to all responses.
// Implements defense-in-depth by hardening HTTP responses against common attacks.
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent MIME-type sniffing
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")

		// Prevent clickjacking attacks
		c.Writer.Header().Set("X-Frame-Options", "DENY")

		// Enable XSS protection (legacy, but still useful for older browsers)
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")

		// Enforce HTTPS (only if not in development)
		// HSTS header tells browsers to only access site via HTTPS for 1 year
		if gin.Mode() == gin.ReleaseMode {
			c.Writer.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}

		// Content Security Policy (basic - can be extended)
		c.Writer.Header().Set("Content-Security-Policy", "default-src 'self'")

		// Referrer Policy (avoid leaking sensitive URLs)
		c.Writer.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Remove server identification headers
		c.Writer.Header().Del("Server")
		c.Writer.Header().Del("X-Powered-By")

		c.Next()
	}
}
