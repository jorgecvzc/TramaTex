package usecase

import (
	"context"
	"errors"
	"fmt"

	domain_model "github.com/joran-cortez/tramatex/internal/iam/domain/model"
	domain_repo "github.com/joran-cortez/tramatex/internal/iam/domain/repository"
)

// CreateUserUseCase allows admins to create users with a specific role.
type CreateUserUseCase struct {
	userRepo domain_repo.Repository
}

// NewCreateUserUseCase creates a new create user use case.
func NewCreateUserUseCase(userRepo domain_repo.Repository) *CreateUserUseCase {
	return &CreateUserUseCase{userRepo: userRepo}
}

// Execute creates a new user with the provided role.
func (uc *CreateUserUseCase) Execute(ctx context.Context, input CreateUserInput) (*CreateUserOutput, error) {
	if err := validateCreateUserInput(input); err != nil {
		return nil, err
	}

	email, err := domain_model.NewEmail(input.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email: %w", err)
	}

	if _, err := uc.userRepo.ByEmail(ctx, email); err == nil {
		return nil, domain_model.ErrUserAlreadyExists
	} else if !errors.Is(err, domain_model.ErrUserNotFound) {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	password, err := domain_model.NewPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("invalid password: %w", err)
	}

	role, err := domain_model.NewRole(input.Role)
	if err != nil {
		return nil, err
	}

	user, err := domain_model.NewUserWithUUID(email, password, role)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	if err := uc.userRepo.Save(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	return &CreateUserOutput{
		ID:    user.ID().String(),
		Email: user.Email().Value(),
		Role:  string(user.Role()),
	}, nil
}

func validateCreateUserInput(input CreateUserInput) error {
	if input.Email == "" {
		return fmt.Errorf("email is required")
	}
	if input.Password == "" {
		return fmt.Errorf("password is required")
	}
	if input.Role == "" {
		return fmt.Errorf("role is required")
	}
	return nil
}
