package usecase

import (
	"context"
	"fmt"

	"github.com/joran-cortez/tramatex/internal/shared/infrastructure/security"
)

// LogoutUserUseCase implements logout logic by revoking tokens.
type LogoutUserUseCase struct {
	jwtService     security.JWTService
	tokenBlacklist security.TokenBlacklist
}

// NewLogoutUserUseCase creates a new logout use case.
func NewLogoutUserUseCase(jwtService security.JWTService, tokenBlacklist security.TokenBlacklist) *LogoutUserUseCase {
	return &LogoutUserUseCase{jwtService: jwtService, tokenBlacklist: tokenBlacklist}
}

// Execute revokes the access token and optional refresh token.
func (uc *LogoutUserUseCase) Execute(ctx context.Context, accessToken string, input LogoutInput) error {
	if accessToken == "" {
		return fmt.Errorf("access token is required")
	}

	claims, err := uc.jwtService.ValidateToken(ctx, accessToken)
	if err != nil {
		return err
	}

	if uc.tokenBlacklist != nil {
		uc.tokenBlacklist.Revoke(accessToken, claims.ExpiresAt())
	}

	if input.RefreshToken != "" && uc.tokenBlacklist != nil {
		refreshClaims, err := uc.jwtService.ValidateToken(ctx, input.RefreshToken)
		if err != nil {
			return err
		}
		uc.tokenBlacklist.Revoke(input.RefreshToken, refreshClaims.ExpiresAt())
	}

	return nil
}
