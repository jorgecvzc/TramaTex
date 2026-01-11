package security

import (
	"context"
)

// JWTService defines the contract for JWT token generation and validation.
// Implemented by infrastructure layer, consumed by application layer.
// Domain layer defines the interface, infrastructure provides the implementation.
type JWTService interface {
	// GenerateAccessToken creates a new access token from claims.
	// TTL: 15 minutes (short-lived, for API requests).
	// Returns the signed token string or error if generation fails.
	GenerateAccessToken(ctx context.Context, claims *TokenClaims) (string, error)

	// GenerateRefreshToken creates a new refresh token from claims.
	// TTL: 7 days (long-lived, for obtaining new access tokens).
	// Returns the signed token string or error if generation fails.
	GenerateRefreshToken(ctx context.Context, claims *TokenClaims) (string, error)

	// ValidateToken verifies and parses a signed token.
	// Returns the extracted TokenClaims if valid, or error if invalid/expired.
	// Errors:
	// - ErrInvalidToken: token format or signature invalid
	// - ErrExpiredToken: token is expired
	ValidateToken(ctx context.Context, token string) (*TokenClaims, error)
}
