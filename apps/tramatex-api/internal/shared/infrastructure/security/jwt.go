package security

import (
	"errors"
	"fmt"
	"time"
)

// TokenClaims represents standard JWT claims as a Value Object.
// Immutable after creation.
type TokenClaims struct {
	subject   string    // user ID (sub claim)
	email     string    // user email
	role      string    // user role
	issuedAt  time.Time // iat claim
	expiresAt time.Time // exp claim
}

// NewTokenClaims creates a new TokenClaims value object.
// Returns error if any parameter is invalid.
func NewTokenClaims(subject, email, role string, issuedAt, expiresAt time.Time) (*TokenClaims, error) {
	// Validate subject (user ID)
	if subject == "" {
		return nil, fmt.Errorf("token subject (user ID) cannot be empty")
	}

	// Validate email
	if email == "" {
		return nil, fmt.Errorf("token email cannot be empty")
	}

	// Validate role
	if role == "" {
		return nil, fmt.Errorf("token role cannot be empty")
	}

	// Validate timestamps
	if expiresAt.Before(issuedAt) {
		return nil, fmt.Errorf("token expiration must be after issuance")
	}

	// Token should not be already expired (soft check)
	if expiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("token expiration is in the past")
	}

	return &TokenClaims{
		subject:   subject,
		email:     email,
		role:      role,
		issuedAt:  issuedAt,
		expiresAt: expiresAt,
	}, nil
}

// Subject returns the token subject (user ID).
func (tc *TokenClaims) Subject() string {
	return tc.subject
}

// Email returns the token email.
func (tc *TokenClaims) Email() string {
	return tc.email
}

// Role returns the token role.
func (tc *TokenClaims) Role() string {
	return tc.role
}

// IssuedAt returns when the token was issued.
func (tc *TokenClaims) IssuedAt() time.Time {
	return tc.issuedAt
}

// ExpiresAt returns when the token will expire.
func (tc *TokenClaims) ExpiresAt() time.Time {
	return tc.expiresAt
}

// IsExpired checks if the token has expired.
// Returns true if current time >= expiresAt.
func (tc *TokenClaims) IsExpired() bool {
	return time.Now().After(tc.expiresAt)
}

// String implements Stringer interface.
func (tc *TokenClaims) String() string {
	return fmt.Sprintf("TokenClaims{sub=%s, email=%s, role=%s, iat=%v, exp=%v}",
		tc.subject, tc.email, tc.role, tc.issuedAt, tc.expiresAt)
}

// Error variables
var (
	ErrInvalidToken  = errors.New("invalid token")
	ErrTokenExpired  = errors.New("token expired")
	ErrInvalidClaims = errors.New("invalid token claims")
)

