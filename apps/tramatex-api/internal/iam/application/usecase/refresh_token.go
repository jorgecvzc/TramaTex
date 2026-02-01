package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	domain_model "github.com/joran-cortez/tramatex/internal/iam/domain/model"
	domain_repo "github.com/joran-cortez/tramatex/internal/iam/domain/repository"
	"github.com/joran-cortez/tramatex/internal/shared/infrastructure/security"
)

// RefreshTokenUseCase implements token refresh logic.
type RefreshTokenUseCase struct {
	userRepo   domain_repo.Repository
	jwtService security.JWTService
}

// NewRefreshTokenUseCase creates a new refresh token use case.
func NewRefreshTokenUseCase(userRepo domain_repo.Repository, jwtService security.JWTService) *RefreshTokenUseCase {
	return &RefreshTokenUseCase{userRepo: userRepo, jwtService: jwtService}
}

// Execute validates the refresh token and issues a new access token.
func (uc *RefreshTokenUseCase) Execute(ctx context.Context, input RefreshInput) (*RefreshOutput, error) {
	if input.RefreshToken == "" {
		return nil, fmt.Errorf("refresh_token is required")
	}

	claims, err := uc.jwtService.ValidateToken(ctx, input.RefreshToken)
	if err != nil {
		return nil, err
	}

	user, err := uc.userRepo.ByID(ctx, claims.Subject())
	if err != nil {
		if errors.Is(err, domain_model.ErrUserNotFound) {
			return nil, domain_model.ErrUserNotFound
		}
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if !user.IsActive() {
		return nil, fmt.Errorf("user account is inactive")
	}

	now := time.Now()
	expiresAt := now.Add(15 * time.Minute)
	newClaims, err := security.NewTokenClaims(
		user.ID(),
		user.Email().Value(),
		string(user.Role()),
		now,
		expiresAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create token claims: %w", err)
	}

	accessToken, err := uc.jwtService.GenerateAccessToken(ctx, newClaims)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	return &RefreshOutput{
		AccessToken: accessToken,
		ExpiresIn:   900,
	}, nil
}
