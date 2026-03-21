package usecase

import (
	"context"
	"fmt"
	"time"

	domain_model "github.com/joran-cortez/tramatex/internal/iam/domain/model"
	domain_repo "github.com/joran-cortez/tramatex/internal/iam/domain/repository"
	"github.com/joran-cortez/tramatex/internal/shared/infrastructure/security"
)

// LoginUseCase implements the login business logic.
// Orchestrates: user repository + password validation + JWT generation.
type LoginUseCase struct {
	userRepo   domain_repo.Repository
	jwtService security.JWTService
}

// NewLoginUseCase creates a new login use case with dependencies.
func NewLoginUseCase(userRepo domain_repo.Repository, jwtService security.JWTService) *LoginUseCase {
	return &LoginUseCase{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

// Execute runs the login use case.
// Steps:
// 1. Validate input (email, password)
// 2. Find user by email
// 3. Verify password
// 4. Generate tokens
// 5. Return user info + tokens
func (uc *LoginUseCase) Execute(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	// Step 1: Validate input
	if err := validateLoginInput(input); err != nil {
		return nil, err
	}

	// Step 2: Create Email VO (validates format)
	email, err := domain_model.NewEmail(input.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email: %w", err)
	}

	// Step 3: Find user by email
	foundUser, err := uc.userRepo.ByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Check if user is active
	if !foundUser.IsActive() {
		return nil, fmt.Errorf("user account is inactive")
	}

	// Step 4: Verify password
	if !foundUser.Password().Matches(input.Password) {
		return nil, domain_model.ErrInvalidPassword
	}

	// Step 5: Generate tokens
	now := time.Now()
	expiresAt := now.Add(15 * time.Minute) // Access token expires in 15 minutes
	claims, err := security.NewTokenClaims(
		foundUser.ID().String(),
		foundUser.Email().Value(),
		string(foundUser.Role()),
		now,       // issuedAt: current time
		expiresAt, // expiresAt: 15 minutes from now
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create token claims: %w", err)
	}

	// Generate access token (15 min)
	accessToken, err := uc.jwtService.GenerateAccessToken(ctx, claims)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token (7 days)
	refreshToken, err := uc.jwtService.GenerateRefreshToken(ctx, claims)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Map User to UserDTO
	userDTO := UserDTO{
		ID:    foundUser.ID().String(),
		Email: foundUser.Email().Value(),
		Role:  string(foundUser.Role()),
	}

	return &LoginOutput{
		User:         userDTO,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    900, // 15 minutes in seconds
	}, nil
}

// validateLoginInput validates the login input data.
func validateLoginInput(input LoginInput) error {
	if input.Email == "" {
		return fmt.Errorf("email is required")
	}

	if input.Password == "" {
		return fmt.Errorf("password is required")
	}

	return nil
}
