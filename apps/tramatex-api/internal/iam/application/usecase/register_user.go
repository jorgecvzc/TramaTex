package usecase

import (
	"context"
	"errors"
	"fmt"

	domain_model "github.com/joran-cortez/tramatex/internal/iam/domain/model"
	domain_repo "github.com/joran-cortez/tramatex/internal/iam/domain/repository"
)

// RegisterUserUseCase implements user registration logic.
type RegisterUserUseCase struct {
	userRepo domain_repo.Repository
}

// NewRegisterUserUseCase creates a new register user use case.
func NewRegisterUserUseCase(userRepo domain_repo.Repository) *RegisterUserUseCase {
	return &RegisterUserUseCase{userRepo: userRepo}
}

// Execute registers a new user with email/password and default role.
func (uc *RegisterUserUseCase) Execute(ctx context.Context, input RegisterInput) (*RegisterOutput, error) {
	if err := validateRegisterInput(input); err != nil {
		return nil, err
	}

	email, err := domain_model.NewEmail(input.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email: %w", err)
	}

	// Check if user already exists
	if _, err := uc.userRepo.ByEmail(ctx, email); err == nil {
		return nil, domain_model.ErrUserAlreadyExists
	} else if !errors.Is(err, domain_model.ErrUserNotFound) {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	password, err := domain_model.NewPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("invalid password: %w", err)
	}

	user, err := domain_model.NewUserWithUUID(email, password, domain_model.RoleCommercial)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	if err := uc.userRepo.Save(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	return &RegisterOutput{
		ID:    user.ID().String(),
		Email: user.Email().Value(),
		Role:  string(user.Role()),
	}, nil
}

func validateRegisterInput(input RegisterInput) error {
	if input.Email == "" {
		return fmt.Errorf("email is required")
	}
	if input.Password == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}
