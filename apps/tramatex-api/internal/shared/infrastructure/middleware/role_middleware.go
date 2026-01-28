package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireRole creates a middleware that checks if the authenticated user
// has one of the allowed roles before accessing the endpoint.
//
// Prerequisites: Must be used AFTER AuthMiddleware (expects claims in context).
//
// Usage:
//
//	router.POST("/organizations", middleware.RequireRole("admin", "manager"), handler.CreateOrganization)
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract user role from context (set by AuthMiddleware)
		roleValue, exists := c.Get("userRole")
		if !exists {
			c.Set("authError", "role_not_found_in_context")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized: user role not found in context",
			})
			c.Abort()
			return
		}

		userRole, ok := roleValue.(string)
		if !ok {
			c.Set("authError", "invalid_role_type_in_context")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error: invalid role type",
			})
			c.Abort()
			return
		}

		// Check if user role is in the list of allowed roles
		roleAllowed := false
		for _, allowedRole := range allowedRoles {
			if userRole == allowedRole {
				roleAllowed = true
				break
			}
		}

		if !roleAllowed {
			c.Set("authError", "insufficient_permissions")
			c.JSON(http.StatusForbidden, gin.H{
				"error": "forbidden: insufficient permissions",
			})
			c.Abort()
			return
		}

		// Role is allowed, continue to next handler
		c.Next()
	}
}
