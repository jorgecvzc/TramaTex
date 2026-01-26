package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joran-cortez/tramatex/internal/shared/infrastructure/middleware"
	"github.com/stretchr/testify/assert"
)

func TestRequireRole(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	t.Run("allows access when user has allowed role", func(t *testing.T) {
		// Setup
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userRole", "admin")

		// Create middleware
		roleMiddleware := middleware.RequireRole("admin", "manager")

		// Execute
		roleMiddleware(c)

		// Assert
		assert.False(t, c.IsAborted())
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("allows access when user has one of multiple allowed roles", func(t *testing.T) {
		// Setup
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userRole", "manager")

		// Create middleware
		roleMiddleware := middleware.RequireRole("admin", "manager", "operator")

		// Execute
		roleMiddleware(c)

		// Assert
		assert.False(t, c.IsAborted())
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("denies access when user role is not in allowed roles", func(t *testing.T) {
		// Setup
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userRole", "operator")

		// Create middleware
		roleMiddleware := middleware.RequireRole("admin", "manager")

		// Execute
		roleMiddleware(c)

		// Assert
		assert.True(t, c.IsAborted())
		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Contains(t, w.Body.String(), "forbidden")
	})

	t.Run("denies access when userRole is not in context", func(t *testing.T) {
		// Setup
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		// Note: NOT setting userRole in context

		// Create middleware
		roleMiddleware := middleware.RequireRole("admin")

		// Execute
		roleMiddleware(c)

		// Assert
		assert.True(t, c.IsAborted())
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "user role not found")
	})

	t.Run("handles invalid role type in context", func(t *testing.T) {
		// Setup
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userRole", 123) // Invalid type (should be string)

		// Create middleware
		roleMiddleware := middleware.RequireRole("admin")

		// Execute
		roleMiddleware(c)

		// Assert
		assert.True(t, c.IsAborted())
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "invalid role type")
	})

	t.Run("works with single allowed role", func(t *testing.T) {
		// Setup
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userRole", "admin")

		// Create middleware with single role
		roleMiddleware := middleware.RequireRole("admin")

		// Execute
		roleMiddleware(c)

		// Assert
		assert.False(t, c.IsAborted())
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("case-sensitive role matching", func(t *testing.T) {
		// Setup
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userRole", "Admin") // Capital A

		// Create middleware
		roleMiddleware := middleware.RequireRole("admin") // lowercase

		// Execute
		roleMiddleware(c)

		// Assert (should deny because of case mismatch)
		assert.True(t, c.IsAborted())
		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}
