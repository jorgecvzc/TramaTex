package usecase

import (
	"context"
	"errors"
	"fmt"

	domain_model "github.com/joran-cortez/tramatex/internal/iam/domain/model"
	domain_repo "github.com/joran-cortez/tramatex/internal/iam/domain/repository"
)

// AssignRoleUseCase assigns a role to an existing user.
type AssignRoleUseCase struct {
	userRepo domain_repo.Repository
}

// NewAssignRoleUseCase creates a new assign role use case.
func NewAssignRoleUseCase(userRepo domain_repo.Repository) *AssignRoleUseCase {
	return &AssignRoleUseCase{userRepo: userRepo}
}

// Execute assigns a role to a user and persists it.
func (uc *AssignRoleUseCase) Execute(ctx context.Context, input AssignRoleInput) (*AssignRoleOutput, error) {
	if input.UserID == "" {
		return nil, fmt.Errorf("user_id is required")
	}
	if input.Role == "" {
		return nil, fmt.Errorf("role is required")
	}

	role, err := domain_model.NewRole(input.Role)
	if err != nil {
		return nil, err
	}

	user, err := uc.userRepo.ByID(ctx, input.UserID)
	if err != nil {
		if errors.Is(err, domain_model.ErrUserNotFound) {
			return nil, domain_model.ErrUserNotFound
		}
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if err := user.ChangeRole(role); err != nil {
		return nil, err
	}

	if err := uc.userRepo.Save(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	return &AssignRoleOutput{
		UserID: user.ID(),
		Role:   string(user.Role()),
	}, nil
}
