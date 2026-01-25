package model_test

import (
	"strings"
	"testing"

	"github.com/joran-cortez/tramatex/internal/iam/domain/model"
)

func TestPasswordNewWithValidLength(t *testing.T) {
	tests := []string{
		"password1",                      // 9 chars
		"a_secure_password",              // 17 chars
		"12345678",                       // Exactly 8 chars (min)
		"aVeryLongPassword123!@#$%^&*()", // 31 chars
	}

	for _, pwd := range tests {
		t.Run(pwd, func(t *testing.T) {
			p, err := model.NewPassword(pwd)
			if err != nil {
				t.Errorf("NewPassword(%q) failed: %v", pwd, err)
			}
			if p == nil {
				t.Error("NewPassword returned nil")
			}
		})
	}
}

func TestPasswordNewWithTooShort(t *testing.T) {
	tests := []string{
		"",        // Empty
		"pass",    // 4 chars
		"1234567", // 7 chars (below minimum)
	}

	for _, pwd := range tests {
		t.Run(pwd, func(t *testing.T) {
			p, err := model.NewPassword(pwd)
			if err == nil {
				t.Errorf("NewPassword(%q) should fail for short password", pwd)
			}
			if p != nil {
				t.Error("NewPassword should return nil on error")
			}
		})
	}
}

func TestPasswordNewWithTooLong(t *testing.T) {
	// Create password longer than 72 chars
	longPwd := ""
	for i := 0; i < 73; i++ {
		longPwd += "a"
	}

	p, err := model.NewPassword(longPwd)
	if err == nil {
		t.Error("NewPassword with >72 chars should fail")
	}
	if p != nil {
		t.Error("NewPassword should return nil on error")
	}
}

func TestPasswordNewWithEmptyString(t *testing.T) {
	p, err := model.NewPassword("")
	if err == nil {
		t.Error("NewPassword(\"\") should fail for empty string")
	}
	if p != nil {
		t.Error("NewPassword(\"\") should return nil")
	}
}

func TestPasswordMatchesWithCorrectPassword(t *testing.T) {
	plaintext := "mySecurePassword123"
	p, _ := model.NewPassword(plaintext)

	if !p.Matches(plaintext) {
		t.Error("Matches() should return true for correct password")
	}
}

func TestPasswordMatchesWithWrongPassword(t *testing.T) {
	plaintext := "mySecurePassword123"
	p, _ := model.NewPassword(plaintext)

	if p.Matches("wrongPassword") {
		t.Error("Matches() should return false for wrong password")
	}
}

func TestPasswordMatchesWithEmpty(t *testing.T) {
	p, _ := model.NewPassword("password123")

	tests := []struct {
		name     string
		pwd      *model.Password
		input    string
		expected bool
	}{
		{"Correct password", p, "password123", true},
		{"Wrong password", p, "wrong", false},
		{"Empty input", p, "", false},
		{"Nil password", nil, "anything", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pwd.Matches(tt.input)
			if result != tt.expected {
				t.Errorf("Matches(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestPasswordNeverStoredPlaintext(t *testing.T) {
	plaintext := "mySecurePassword"
	p, _ := model.NewPassword(plaintext)

	hash := p.Hash()

	// Hash should not be equal to plaintext (obviously, it's hashed)
	if hash == plaintext {
		t.Error("Password hash should never be plaintext")
	}

	// Hash should be bcrypt format (starts with $2a$, $2b$, or $2x$)
	if len(hash) < 20 || !isBcryptHash(hash) {
		t.Errorf("Hash format invalid: %v", hash)
	}
}

// isBcryptHash checks if a string is a valid bcrypt hash
func isBcryptHash(hash string) bool {
	return len(hash) >= 20 && (strings.HasPrefix(hash, "$2a$") || strings.HasPrefix(hash, "$2b$") || strings.HasPrefix(hash, "$2x$"))
}

func TestPasswordBcryptCostAtLeast10(t *testing.T) {
	if model.BcryptCost < 10 {
		t.Errorf("BcryptCost = %d, must be >= 10 for security", model.BcryptCost)
	}
}

func TestPasswordString(t *testing.T) {
	p, _ := model.NewPassword("password123")

	str := p.String()
	if str != "[REDACTED]" {
		t.Errorf("String() = %q, want [REDACTED]", str)
	}
}

