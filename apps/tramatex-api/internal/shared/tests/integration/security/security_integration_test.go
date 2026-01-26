package security_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joran-cortez/tramatex/internal/shared/infrastructure/logging"
	infra_middleware "github.com/joran-cortez/tramatex/internal/shared/infrastructure/middleware"
	"github.com/joran-cortez/tramatex/internal/shared/infrastructure/security"
	http_middleware "github.com/joran-cortez/tramatex/internal/shared/interfaces/http/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	// Initialize logger for tests
	logging.InitLogger("test")
}

// generateTestToken creates an access token for testing
func generateTestToken(t *testing.T, jwtService security.JWTService, userID, email, role string) string {
	now := time.Now()
	expiresAt := now.Add(15 * time.Minute)

	claims, err := security.NewTokenClaims(userID, email, role, now, expiresAt)
	require.NoError(t, err)

	token, err := jwtService.GenerateAccessToken(context.Background(), claims)
	require.NoError(t, err)

	return token
}

// setupTestRouter creates a Gin router with all security middlewares applied
func setupTestRouter(t *testing.T) (*gin.Engine, security.JWTService) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Create JWT service for token generation (use hours: 168h = 7 days)
	jwtService, err := security.NewJWTService("test-secret-key-min-32-chars-long!", "15m", "168h")
	require.NoError(t, err)

	// Apply security middlewares in correct order
	router.Use(infra_middleware.SecurityHeadersMiddleware())
	router.Use(infra_middleware.RequestIDMiddleware())
	router.Use(infra_middleware.SecurityLoggerMiddleware())
	router.Use(infra_middleware.CORSMiddleware("http://localhost:3000,http://localhost:5173"))
	router.Use(infra_middleware.ErrorHandlerMiddleware("production"))

	// Create protected routes
	protected := router.Group("/api/protected")
	protected.Use(http_middleware.AuthMiddleware(jwtService))
	{
		// Admin/Manager only endpoint
		protected.POST("/admin-only", infra_middleware.RequireRole("admin", "manager"), func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "admin access granted"})
		})

		// All authenticated users
		protected.GET("/read-only", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "read access granted"})
		})
	}

	return router, jwtService
}

// TestRBAC_AdminCanAccess verifies admin role can access write endpoints
func TestRBAC_AdminCanAccess(t *testing.T) {
	router, jwtService := setupTestRouter(t)

	// Generate admin token
	token := generateTestToken(t, jwtService, "admin-user-id", "admin@tramatex.com", "admin")

	// Make request with admin token
	req := httptest.NewRequest("POST", "/api/protected/admin-only", bytes.NewBuffer([]byte("{}")))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "admin access granted", response["message"])
}

// TestRBAC_ManagerCanAccess verifies manager role can access write endpoints
func TestRBAC_ManagerCanAccess(t *testing.T) {
	router, jwtService := setupTestRouter(t)

	// Generate manager token
	token := generateTestToken(t, jwtService, "manager-user-id", "manager@tramatex.com", "manager")

	req := httptest.NewRequest("POST", "/api/protected/admin-only", bytes.NewBuffer([]byte("{}")))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

// TestRBAC_OperatorDeniedAccess verifies operator role cannot access write endpoints
func TestRBAC_OperatorDeniedAccess(t *testing.T) {
	router, jwtService := setupTestRouter(t)

	// Generate operator token
	token := generateTestToken(t, jwtService, "operator-user-id", "operator@tramatex.com", "operator")

	req := httptest.NewRequest("POST", "/api/protected/admin-only", bytes.NewBuffer([]byte("{}")))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should get 403 Forbidden
	assert.Equal(t, 403, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "insufficient permissions")
}

// TestRBAC_OperatorCanRead verifies operator can access read-only endpoints
func TestRBAC_OperatorCanRead(t *testing.T) {
	router, jwtService := setupTestRouter(t)

	token := generateTestToken(t, jwtService, "operator-user-id", "operator@tramatex.com", "operator")

	req := httptest.NewRequest("GET", "/api/protected/read-only", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

// TestRBAC_UnauthenticatedDenied verifies unauthenticated requests are rejected
func TestRBAC_UnauthenticatedDenied(t *testing.T) {
	router, _ := setupTestRouter(t)

	req := httptest.NewRequest("GET", "/api/protected/read-only", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should get 401 Unauthorized
	assert.Equal(t, 401, w.Code)
}

// TestSecurityHeaders_AllPresent verifies all security headers are set
func TestSecurityHeaders_AllPresent(t *testing.T) {
	router, _ := setupTestRouter(t)

	req := httptest.NewRequest("GET", "/api/protected/read-only", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Verify security headers
	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
	assert.Equal(t, "1; mode=block", w.Header().Get("X-XSS-Protection"))
	assert.Equal(t, "default-src 'self'", w.Header().Get("Content-Security-Policy"))
	assert.Equal(t, "strict-origin-when-cross-origin", w.Header().Get("Referrer-Policy"))

	// Server headers should be removed
	assert.Empty(t, w.Header().Get("Server"))
	assert.Empty(t, w.Header().Get("X-Powered-By"))
}

// TestRequestID_Generated verifies request ID is generated and added to headers
func TestRequestID_Generated(t *testing.T) {
	router, jwtService := setupTestRouter(t)

	token := generateTestToken(t, jwtService, "user-id", "user@tramatex.com", "admin")

	req := httptest.NewRequest("GET", "/api/protected/read-only", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Verify X-Request-ID header is present and is a valid UUID format
	requestID := w.Header().Get("X-Request-ID")
	assert.NotEmpty(t, requestID)
	assert.Len(t, requestID, 36) // UUID format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
}

// TestRequestID_PreservedFromClient verifies client-provided request ID is preserved
func TestRequestID_PreservedFromClient(t *testing.T) {
	router, jwtService := setupTestRouter(t)

	token := generateTestToken(t, jwtService, "user-id", "user@tramatex.com", "admin")

	clientRequestID := "client-request-12345"
	req := httptest.NewRequest("GET", "/api/protected/read-only", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Request-ID", clientRequestID)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Verify the client's request ID is preserved
	assert.Equal(t, clientRequestID, w.Header().Get("X-Request-ID"))
}

// TestCORS_AllowedOrigin verifies CORS allows whitelisted origins
func TestCORS_AllowedOrigin(t *testing.T) {
	router, _ := setupTestRouter(t)

	req := httptest.NewRequest("GET", "/api/protected/read-only", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, "http://localhost:3000", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "Origin", w.Header().Get("Vary"))
}

// TestCORS_DisallowedOrigin verifies CORS blocks non-whitelisted origins
func TestCORS_DisallowedOrigin(t *testing.T) {
	router, _ := setupTestRouter(t)

	req := httptest.NewRequest("GET", "/api/protected/read-only", nil)
	req.Header.Set("Origin", "http://evil.com")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// CORS header should not be set for disallowed origin
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
}

// TestCORS_Preflight verifies CORS preflight requests are handled
func TestCORS_Preflight(t *testing.T) {
	router, _ := setupTestRouter(t)

	req := httptest.NewRequest("OPTIONS", "/api/protected/admin-only", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, 204, w.Code)
	assert.Equal(t, "http://localhost:3000", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "POST")
}

// TestIntegration_FullSecurityStack verifies all security controls work together
func TestIntegration_FullSecurityStack(t *testing.T) {
	router, jwtService := setupTestRouter(t)

	// Generate admin token
	token := generateTestToken(t, jwtService, "admin-123", "admin@tramatex.com", "admin")

	// Make request with all security context
	req := httptest.NewRequest("POST", "/api/protected/admin-only", bytes.NewBuffer([]byte("{}")))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Verify successful access (RBAC + Auth work)
	assert.Equal(t, 200, w.Code)

	// Verify CORS headers
	assert.Equal(t, "http://localhost:3000", w.Header().Get("Access-Control-Allow-Origin"))

	// Verify security headers
	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))

	// Verify request ID
	assert.NotEmpty(t, w.Header().Get("X-Request-ID"))

	// Verify response
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "admin access granted", response["message"])
}
