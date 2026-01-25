package security

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTServiceImpl is the concrete implementation of JWTService.
type JWTServiceImpl struct {
	secret          string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

// NewJWTService creates a new JWT service implementation.
func NewJWTService(secret string, accessTokenTTL, refreshTokenTTL string) (JWTService, error) {
	// Parse TTL durations
	accessDuration, err := time.ParseDuration(accessTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("invalid accessTokenTTL: %w", err)
	}

	refreshDuration, err := time.ParseDuration(refreshTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("invalid refreshTokenTTL: %w", err)
	}

	if secret == "" {
		return nil, fmt.Errorf("JWT secret cannot be empty")
	}

	return &JWTServiceImpl{
		secret:          secret,
		accessTokenTTL:  accessDuration,
		refreshTokenTTL: refreshDuration,
	}, nil
}

// GenerateAccessToken creates a new access token.
func (j *JWTServiceImpl) GenerateAccessToken(ctx context.Context, claims *TokenClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  claims.Subject(),
		"email": claims.Email(),
		"role": claims.Role(),
		"iat":  claims.IssuedAt().Unix(),
		"exp":  time.Now().Add(j.accessTokenTTL).Unix(),
	})

	tokenString, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign access token: %w", err)
	}

	return tokenString, nil
}

// GenerateRefreshToken creates a new refresh token.
func (j *JWTServiceImpl) GenerateRefreshToken(ctx context.Context, claims *TokenClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": claims.Subject(),
		"iat": claims.IssuedAt().Unix(),
		"exp": time.Now().Add(j.refreshTokenTTL).Unix(),
	})

	tokenString, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken verifies and parses a signed token.
func (j *JWTServiceImpl) ValidateToken(ctx context.Context, tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method: %v", token.Header["alg"])
		}
		return []byte(j.secret), nil
	})

	if err != nil {
		if err.Error() == "token is expired" {
			return nil, ErrTokenExpired
		}
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidClaims
	}

	// Extract claims
	subject, ok := (*claims)["sub"].(string)
	if !ok || subject == "" {
		return nil, fmt.Errorf("invalid subject claim")
	}

	email, ok := (*claims)["email"].(string)
	if !ok || email == "" {
		return nil, fmt.Errorf("invalid email claim")
	}

	role, ok := (*claims)["role"].(string)
	if !ok || role == "" {
		return nil, fmt.Errorf("invalid role claim")
	}

	// Create TokenClaims value object
	tokenClaims, err := NewTokenClaims(subject, email, role, time.Now(), time.Unix(int64((*claims)["exp"].(float64)), 0))
	if err != nil {
		return nil, fmt.Errorf("failed to create token claims: %w", err)
	}

	return tokenClaims, nil
}
