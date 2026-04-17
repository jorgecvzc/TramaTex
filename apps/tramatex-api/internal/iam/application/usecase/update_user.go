package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	domain_model "github.com/joran-cortez/tramatex/internal/iam/domain/model"
	domain_repo "github.com/joran-cortez/tramatex/internal/iam/domain/repository"
)

// UpdateUserUseCase allows updating a user's basic info.
type UpdateUserUseCase struct {
	userRepo domain_repo.Repository
}

// NewUpdateUserUseCase creates a new update user use case.
func NewUpdateUserUseCase(userRepo domain_repo.Repository) *UpdateUserUseCase {
	return &UpdateUserUseCase{userRepo: userRepo}
}

// Execute updates an existing user.
func (uc *UpdateUserUseCase) Execute(ctx context.Context, input UpdateUserInput) (*UpdateUserOutput, error) {
	id, err := uuid.Parse(input.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	user, err := uc.userRepo.ByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain_model.ErrUserNotFound) {
			return nil, domain_model.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}

	if input.Email != nil && *input.Email != "" {
		email, err := domain_model.NewEmail(*input.Email)
		if err != nil {
			return nil, fmt.Errorf("invalid email: %w", err)
		}
		
		// Check if email already exists for another user
		existing, err := uc.userRepo.ByEmail(ctx, email)
		if err == nil && existing.ID() != user.ID() {
			return nil, domain_model.ErrUserAlreadyExists
		} else if err != nil && !errors.Is(err, domain_model.ErrUserNotFound) {
			return nil, fmt.Errorf("failed to check existing email: %w", err)
		}

		if err := user.ChangeEmail(email); err != nil {
			return nil, err
		}
	}

	if input.Password != nil && *input.Password != "" {
		password, err := domain_model.NewPassword(*input.Password)
		if err != nil {
			return nil, fmt.Errorf("invalid password: %w", err)
		}
		if err := user.ChangePassword(password); err != nil {
			return nil, err
		}
	}

	if err := uc.userRepo.Save(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	return &UpdateUserOutput{
		ID:    user.ID().String(),
		Email: user.Email().Value(),
		Role:  string(user.Role()),
	}, nil
}
