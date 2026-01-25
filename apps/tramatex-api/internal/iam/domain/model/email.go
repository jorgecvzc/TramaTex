package model

import (
	"fmt"
	"regexp"
	"strings"
)

// Email represents an email address as a Value Object.
// Immutable and validated at construction time.
type Email struct {
	value string
}

// emailRegex is a simplified RFC 5322 validation pattern.
// Matches: local@domain with basic character classes.
// Includes support for plus-addressing (test+tag@domain.com).
// Prevents consecutive dots in domain part.
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._+\-]+@([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`)

// NewEmail creates a new Email value object with validation.
// Returns error if email format is invalid.
func NewEmail(value string) (*Email, error) {
	if err := validateEmail(value); err != nil {
		return nil, err
	}

	return &Email{value: value}, nil
}

// validateEmail performs email format validation.
func validateEmail(value string) error {
	// Check empty
	if value == "" {
		return fmt.Errorf("email cannot be empty")
	}

	// Trim whitespace
	value = strings.TrimSpace(value)

	// Check length (RFC 5321: max 254 chars)
	if len(value) > 254 {
		return fmt.Errorf("email is too long (max 254 characters)")
	}

	// Check format with regex (simplified RFC 5322)
	if !emailRegex.MatchString(value) {
		return fmt.Errorf("invalid email format")
	}

	// Check local part length (before @)
	parts := strings.Split(value, "@")
	if len(parts[0]) == 0 || len(parts[0]) > 64 {
		return fmt.Errorf("email local part must be 1-64 characters")
	}

	return nil
}

// Value returns the email string representation.
func (e *Email) Value() string {
	return e.value
}

// Equals compares two Email value objects for equality.
func (e *Email) Equals(other *Email) bool {
	if e == nil || other == nil {
		return e == other
	}
	return e.value == other.value
}

// String implements the Stringer interface.
func (e *Email) String() string {
	return e.value
}
