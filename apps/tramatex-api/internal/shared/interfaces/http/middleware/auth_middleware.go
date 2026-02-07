package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joran-cortez/tramatex/internal/shared/infrastructure/security"
)

func AuthMiddleware(jwtService security.JWTService, blacklist security.TokenBlacklist) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Set("authError", "missing_authorization_header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Set("authError", "invalid_authorization_format")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			return
		}

		tokenString := parts[1]
		if blacklist != nil && blacklist.IsRevoked(tokenString) {
			c.Set("authError", "revoked_token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims, err := jwtService.ValidateToken(c.Request.Context(), tokenString)
		if err != nil {
			c.Set("authError", "invalid_token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Set user context
		c.Set("userID", claims.Subject())
		c.Set("userEmail", claims.Email())
		c.Set("userRole", claims.Role())

		ctx := context.WithValue(c.Request.Context(), "userID", claims.Subject())
		ctx = context.WithValue(ctx, "actorID", claims.Subject())
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
