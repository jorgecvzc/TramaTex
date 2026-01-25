package model

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Password represents a hashed password as a Value Object.
// Immutable and hashed at construction time using bcrypt.
// Never stores plaintext passwords.
type Password struct {
	hash string
}

// MinPasswordLength defines minimum password length for security.
const MinPasswordLength = 8

// MaxPasswordLength defines maximum password length (bcrypt constraint).
const MaxPasswordLength = 72

// BcryptCost is the cost parameter for bcrypt hashing.
// Must be between 4 and 31. Higher = more secure but slower.
// 10 = ~100ms hashing time (good balance for MVP).
const BcryptCost = 10

// NewPassword creates a new Password value object.
// Hashes the plaintext password using bcrypt.
// Returns error if password validation fails or hashing fails.
func NewPassword(plaintext string) (*Password, error) {
	if err := validatePasswordInput(plaintext); err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(plaintext), BcryptCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	return &Password{hash: string(hash)}, nil
}

// NewPasswordFromHash creates a Password from an existing bcrypt hash.
// Used when reading from database (no hashing needed).
// Infrastructure layer only.
func NewPasswordFromHash(hash string) *Password {
	return &Password{hash: hash}
}

// validatePasswordInput validates plaintext password before hashing.
func validatePasswordInput(plaintext string) error {
	if plaintext == "" {
		return fmt.Errorf("password cannot be empty")
	}

	if len(plaintext) < MinPasswordLength {
		return fmt.Errorf("password must be at least %d characters", MinPasswordLength)
	}

	if len(plaintext) > MaxPasswordLength {
		return fmt.Errorf("password must be at most %d characters", MaxPasswordLength)
	}

	return nil
}

// Matches compares a plaintext password with the stored hash.
// Returns true if password matches, false otherwise.
// Never exposes the hash or plaintext.
func (p *Password) Matches(plaintext string) bool {
	if p == nil || p.hash == "" || plaintext == "" {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(p.hash), []byte(plaintext))
	return err == nil
}

// Hash returns the bcrypt hash (private for domain, infrastructure layer only).
func (p *Password) Hash() string {
	return p.hash
}

// String implements Stringer but returns masked hash for security.
func (p *Password) String() string {
	return "[REDACTED]"
}
