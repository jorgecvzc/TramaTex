package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system as a Root Aggregate.
// Immutable after creation (all fields private, only getters exposed).
type User struct {
	id        uuid.UUID
	email     *Email
	password  *Password
	role      Role
	active    bool
	createdAt time.Time
	updatedAt time.Time
}

// NewUser creates a new User with validation of invariants.
// Returns error if any invariant is violated.
func NewUser(id uuid.UUID, email *Email, password *Password, role Role) (*User, error) {
	// Validate ID
	if id == uuid.Nil {
		return nil, fmt.Errorf("user ID cannot be empty")
	}

	// Validate Email
	if email == nil {
		return nil, fmt.Errorf("user email cannot be nil")
	}

	// Validate Password
	if password == nil {
		return nil, fmt.Errorf("user password cannot be nil")
	}

	// Validate Role
	if !role.IsValid() {
		return nil, fmt.Errorf("invalid role: %s", role)
	}

	now := time.Now()

	return &User{
		id:        id,
		email:     email,
		password:  password,
		role:      role,
		active:    true, // Default to active
		createdAt: now,
		updatedAt: now,
	}, nil
}

// ID returns the user's unique identifier.
func (u *User) ID() uuid.UUID {
	return u.id
}

// Email returns the user's email as a value object.
func (u *User) Email() *Email {
	return u.email
}

// Password returns the user's password as a value object.
func (u *User) Password() *Password {
	return u.password
}

// Role returns the user's role.
func (u *User) Role() Role {
	return u.role
}

// IsActive returns whether the user is active.
func (u *User) IsActive() bool {
	return u.active
}

// CreatedAt returns when the user was created.
func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

// UpdatedAt returns when the user was last updated.
func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

// ChangePassword updates user's password (internal only).
// Used by infrastructure layer for password reset flows.
// Returns error if password is invalid.
func (u *User) ChangePassword(newPassword *Password) error {
	if newPassword == nil {
		return fmt.Errorf("new password cannot be nil")
	}
	u.password = newPassword
	u.updatedAt = time.Now()
	return nil
}

// Deactivate marks user as inactive.
func (u *User) Deactivate() {
	u.active = false
	u.updatedAt = time.Now()
}

// Activate marks user as active.
func (u *User) Activate() {
	u.active = true
	u.updatedAt = time.Now()
}

// ChangeRole updates the user's role.
// Returns error if role is invalid.
func (u *User) ChangeRole(newRole Role) error {
	if !newRole.IsValid() {
		return fmt.Errorf("invalid role: %s", newRole)
	}
	u.role = newRole
	u.updatedAt = time.Now()
	return nil
}

// NewUserWithUUID generates a new user with a randomly generated UUID as ID.
func NewUserWithUUID(email *Email, password *Password, role Role) (*User, error) {
	return NewUser(uuid.New(), email, password, role)
}

// Error variables
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrUserAlreadyExists = errors.New("user already exists")
)
