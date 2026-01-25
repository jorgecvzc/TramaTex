package security_test

import (
	"testing"
	"time"

	"github.com/joran-cortez/tramatex/internal/shared/infrastructure/security"
)

func TestTokenClaimsCreation(t *testing.T) {
	now := time.Now()
	expiresIn := 15 * time.Minute
	expiresAt := now.Add(expiresIn)

	claims, err := security.NewTokenClaims("user-123", "user@example.com", "admin", now, expiresAt)

	if err != nil {
		t.Errorf("NewTokenClaims with valid data should not fail: %v", err)
	}

	if claims == nil {
		t.Error("NewTokenClaims should return claims instance")
	}

	if claims.Subject() != "user-123" {
		t.Errorf("Subject() = %q, want %q", claims.Subject(), "user-123")
	}

	if claims.Email() != "user@example.com" {
		t.Errorf("Email() = %q, want %q", claims.Email(), "user@example.com")
	}

	if claims.Role() != "admin" {
		t.Errorf("Role() = %q, want %q", claims.Role(), "admin")
	}
}

func TestTokenClaimsIsExpired(t *testing.T) {
	now := time.Now()

	// Create expired token (expiresAt is in the past)
	expiredAt := now.Add(-1 * time.Hour)

	claims, err := security.NewTokenClaims("user-123", "user@example.com", "admin", now.Add(-2*time.Hour), expiredAt)

	// Should fail because expiration is in the past
	if err == nil {
		t.Error("NewTokenClaims with past expiration should fail")
	}

	if claims != nil {
		t.Error("NewTokenClaims should return nil on error")
	}
}

func TestTokenClaimsNotExpired(t *testing.T) {
	now := time.Now()
	expiresAt := now.Add(15 * time.Minute) // 15 minutes in future

	claims, err := security.NewTokenClaims("user-123", "user@example.com", "admin", now, expiresAt)

	if err != nil {
		t.Errorf("NewTokenClaims with future expiration should succeed: %v", err)
	}

	if claims.IsExpired() {
		t.Error("Token should not be expired when expiresAt is in future")
	}
}

func TestTokenClaimsJustExpired(t *testing.T) {
	now := time.Now()
	expiresAt := now.Add(-1 * time.Millisecond) // Just expired

	// Should fail creation because expiration is in past
	claims, err := security.NewTokenClaims("user-123", "user@example.com", "admin", now, expiresAt)

	if err == nil {
		t.Error("NewTokenClaims with just-expired time should fail")
	}

	if claims != nil {
		t.Error("NewTokenClaims should return nil")
	}
}

func TestTokenClaimsAllFields(t *testing.T) {
	subject := "user-123"
	email := "user@example.com"
	role := "manager"
	now := time.Now()
	expiresAt := now.Add(7 * 24 * time.Hour) // 7 days

	claims, _ := security.NewTokenClaims(subject, email, role, now, expiresAt)

	// Verify all fields
	if claims.Subject() != subject {
		t.Errorf("Subject mismatch")
	}

	if claims.Email() != email {
		t.Errorf("Email mismatch")
	}

	if claims.Role() != role {
		t.Errorf("Role mismatch")
	}

	if !claims.IssuedAt().Equal(now) {
		t.Errorf("IssuedAt mismatch")
	}

	if !claims.ExpiresAt().Equal(expiresAt) {
		t.Errorf("ExpiresAt mismatch")
	}
}

func TestTokenClaimsWithEmptySubject(t *testing.T) {
	now := time.Now()
	expiresAt := now.Add(15 * time.Minute)

	claims, err := security.NewTokenClaims("", "user@example.com", "admin", now, expiresAt)

	if err == nil {
		t.Error("NewTokenClaims with empty subject should fail")
	}

	if claims != nil {
		t.Error("NewTokenClaims should return nil on error")
	}
}

func TestTokenClaimsWithEmptyEmail(t *testing.T) {
	now := time.Now()
	expiresAt := now.Add(15 * time.Minute)

	claims, err := security.NewTokenClaims("user-123", "", "admin", now, expiresAt)

	if err == nil {
		t.Error("NewTokenClaims with empty email should fail")
	}

	if claims != nil {
		t.Error("NewTokenClaims should return nil on error")
	}
}

func TestTokenClaimsWithEmptyRole(t *testing.T) {
	now := time.Now()
	expiresAt := now.Add(15 * time.Minute)

	claims, err := security.NewTokenClaims("user-123", "user@example.com", "", now, expiresAt)

	if err == nil {
		t.Error("NewTokenClaims with empty role should fail")
	}

	if claims != nil {
		t.Error("NewTokenClaims should return nil on error")
	}
}

func TestTokenClaimsExpirationAfterIssuance(t *testing.T) {
	now := time.Now()
	expiresAt := now.Add(-1 * time.Hour) // Expires before issuance

	claims, err := security.NewTokenClaims("user-123", "user@example.com", "admin", now, expiresAt)

	if err == nil {
		t.Error("NewTokenClaims with expiration before issuance should fail")
	}

	if claims != nil {
		t.Error("NewTokenClaims should return nil on error")
	}
}

func TestTokenClaimsString(t *testing.T) {
	now := time.Now()
	expiresAt := now.Add(15 * time.Minute)

	claims, _ := security.NewTokenClaims("user-123", "user@example.com", "admin", now, expiresAt)

	str := claims.String()

	// Should contain key information (not redacted, unlike password)
	if len(str) == 0 {
		t.Error("String() returned empty")
	}

	if !contains(str, "user-123") {
		t.Errorf("String() should contain subject")
	}

	if !contains(str, "user@example.com") {
		t.Errorf("String() should contain email")
	}
}

// Helper function for testing
func contains(s, substr string) bool {
	for i := 0; i < len(s); i++ {
		if len(s[i:]) < len(substr) {
			return false
		}
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
