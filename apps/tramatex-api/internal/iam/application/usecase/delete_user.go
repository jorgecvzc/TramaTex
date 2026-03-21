package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	domain_repo "github.com/joran-cortez/tramatex/internal/iam/domain/repository"
)

// DeleteUserUseCase allows admins to soft delete a user.
type DeleteUserUseCase struct {
	userRepo domain_repo.Repository
}

// NewDeleteUserUseCase creates a new delete user use case.
func NewDeleteUserUseCase(userRepo domain_repo.Repository) *DeleteUserUseCase {
	return &DeleteUserUseCase{userRepo: userRepo}
}

// Execute deletes a user by id.
func (uc *DeleteUserUseCase) Execute(ctx context.Context, userID string) error {
	if userID == "" {
		return fmt.Errorf("user id is required")
	}

	parsedID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user id: %w", err)
	}

	if err := uc.userRepo.Delete(ctx, parsedID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
