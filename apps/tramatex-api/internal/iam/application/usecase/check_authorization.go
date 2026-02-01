package usecase

import (
	"context"
	"errors"
	"fmt"

	domain_model "github.com/joran-cortez/tramatex/internal/iam/domain/model"
	domain_repo "github.com/joran-cortez/tramatex/internal/iam/domain/repository"
)

// CheckAuthorizationUseCase verifies if a user has any of the required roles.
type CheckAuthorizationUseCase struct {
	userRepo domain_repo.Repository
}

// NewCheckAuthorizationUseCase creates a new authorization check use case.
func NewCheckAuthorizationUseCase(userRepo domain_repo.Repository) *CheckAuthorizationUseCase {
	return &CheckAuthorizationUseCase{userRepo: userRepo}
}

// Execute checks if the user has any of the required roles.
func (uc *CheckAuthorizationUseCase) Execute(ctx context.Context, input CheckAuthorizationInput) (*CheckAuthorizationOutput, error) {
	if len(input.RequiredRoles) == 0 {
		return nil, fmt.Errorf("required_roles is required")
	}
	if input.UserID == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	user, err := uc.userRepo.ByID(ctx, input.UserID)
	if err != nil {
		if errors.Is(err, domain_model.ErrUserNotFound) {
			return nil, domain_model.ErrUserNotFound
		}
		return nil, fmt.Errorf("user not found: %w", err)
	}

	userRole := string(user.Role())
	allowed := false
	for _, role := range input.RequiredRoles {
		if role == userRole {
			allowed = true
			break
		}
	}

	return &CheckAuthorizationOutput{
		Allowed: allowed,
		Role:    userRole,
	}, nil
}
