package middleware_test

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joran-cortez/tramatex/internal/shared/infrastructure/middleware"
	"github.com/stretchr/testify/assert"
)

func init() {
	// Force test mode for all tests in this package
	gin.SetMode(gin.TestMode)
}

func TestSecurityHeadersMiddleware(t *testing.T) {
	router := gin.New()

	router.Use(middleware.SecurityHeadersMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Verify security headers
	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
	assert.Equal(t, "1; mode=block", w.Header().Get("X-XSS-Protection"))
	assert.Equal(t, "default-src 'self'", w.Header().Get("Content-Security-Policy"))
	assert.Equal(t, "strict-origin-when-cross-origin", w.Header().Get("Referrer-Policy"))

	// Verify server headers are removed
	assert.Empty(t, w.Header().Get("Server"))
	assert.Empty(t, w.Header().Get("X-Powered-By"))

	// HSTS should not be set in test mode (debug mode)
	assert.Empty(t, w.Header().Get("Strict-Transport-Security"))
}

func TestSecurityHeadersMiddleware_ProductionMode(t *testing.T) {
	// Save original mode and restore after test
	originalMode := gin.Mode()
	defer gin.SetMode(originalMode)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(middleware.SecurityHeadersMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// HSTS should be set in release mode
	assert.Equal(t, "max-age=31536000; includeSubDomains", w.Header().Get("Strict-Transport-Security"))
}
