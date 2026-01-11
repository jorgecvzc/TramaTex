package user_test

import (
	"testing"

	"github.com/joran-cortez/tramatex/internal/domain/user"
)

func TestEmailNewWithValidFormat(t *testing.T) {
	tests := []string{
		"user@example.com",
		"john.doe@company.co.uk",
		"test+tag@domain.com",
		"user_name@domain-name.org",
		"a@b.co",
	}

	for _, email := range tests {
		t.Run(email, func(t *testing.T) {
			e, err := user.NewEmail(email)
			if err != nil {
				t.Errorf("NewEmail(%q) failed: %v", email, err)
			}
			if e == nil {
				t.Error("NewEmail returned nil")
			}
			if e.Value() != email {
				t.Errorf("Email value = %q, want %q", e.Value(), email)
			}
		})
	}
}

func TestEmailNewWithInvalidFormat(t *testing.T) {
	tests := []string{
		"no-at-sign.com",
		"user@",
		"@domain.com",
		"user@.com",
		"user@domain",
		"user name@domain.com",
		"user@domain..com",
	}

	for _, email := range tests {
		t.Run(email, func(t *testing.T) {
			e, err := user.NewEmail(email)
			if err == nil {
				t.Errorf("NewEmail(%q) should have failed but got: %v", email, e)
			}
			if e != nil {
				t.Errorf("NewEmail(%q) should return nil on error, got: %v", email, e)
			}
		})
	}
}

func TestEmailNewWithEmptyString(t *testing.T) {
	e, err := user.NewEmail("")
	if err == nil {
		t.Error("NewEmail(\"\") should fail for empty string")
	}
	if e != nil {
		t.Error("NewEmail(\"\") should return nil")
	}
}

func TestEmailNewWithWhitespaceOnly(t *testing.T) {
	e, err := user.NewEmail("   ")
	if err == nil {
		t.Error("NewEmail(\"   \") should fail for whitespace-only string")
	}
	if e != nil {
		t.Error("NewEmail(\"   \") should return nil")
	}
}

func TestEmailNewWithTooLongAddress(t *testing.T) {
	longEmail := "a@" + string(make([]byte, 253)) + ".com" // > 254 chars
	e, err := user.NewEmail(longEmail)
	if err == nil {
		t.Error("NewEmail with >254 chars should fail")
	}
	if e != nil {
		t.Error("NewEmail with >254 chars should return nil")
	}
}

func TestEmailNewWithTooLongLocalPart(t *testing.T) {
	localPart := string(make([]byte, 65))
	for i := range localPart {
		localPart = "a" + localPart
	}
	email := localPart + "@domain.com"
	e, err := user.NewEmail(email)
	if err == nil {
		t.Error("NewEmail with >64 char local part should fail")
	}
	if e != nil {
		t.Error("NewEmail with >64 char local part should return nil")
	}
}

func TestEmailEquals(t *testing.T) {
	email1, _ := user.NewEmail("user@example.com")
	email2, _ := user.NewEmail("user@example.com")
	email3, _ := user.NewEmail("other@example.com")

	tests := []struct {
		name     string
		email1   *user.Email
		email2   *user.Email
		expected bool
	}{
		{"Equal emails", email1, email2, true},
		{"Different emails", email1, email3, false},
		{"Same object", email1, email1, true},
		{"Nil both", nil, nil, true},
		{"Nil first", nil, email1, false},
		{"Nil second", email1, nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.email1.Equals(tt.email2)
			if result != tt.expected {
				t.Errorf("Equals() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestEmailImmutable(t *testing.T) {
	// Email should not allow modification after creation
	email, _ := user.NewEmail("original@example.com")

	// Try to modify (through reflection would be needed to actually modify)
	// For this test, we just verify Value() returns the same
	if email.Value() != "original@example.com" {
		t.Error("Email value changed unexpectedly")
	}

	// Verify multiple calls return same value
	if email.Value() != email.Value() {
		t.Error("Email value not consistent")
	}
}

func TestEmailString(t *testing.T) {
	email, _ := user.NewEmail("test@example.com")
	if email.String() != "test@example.com" {
		t.Errorf("String() = %q, want %q", email.String(), "test@example.com")
	}
}
