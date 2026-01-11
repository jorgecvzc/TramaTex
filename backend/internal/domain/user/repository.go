package user

import (
	"context"
)

// Repository defines the contract for user persistence.
// Implemented by the infrastructure layer (database adapters).
// Domain layer depends on this interface, infrastructure implements it.
type Repository interface {
	// ByID retrieves a user by their unique identifier.
	// Returns ErrUserNotFound if user does not exist.
	ByID(ctx context.Context, id string) (*User, error)

	// ByEmail retrieves a user by their email address.
	// Returns ErrUserNotFound if user does not exist.
	ByEmail(ctx context.Context, email *Email) (*User, error)

	// Save persists a user to storage.
	// Creates new user or updates existing one.
	Save(ctx context.Context, user *User) error

	// Delete removes a user from storage.
	// Should be soft delete (mark as inactive) in production.
	Delete(ctx context.Context, id string) error
}
