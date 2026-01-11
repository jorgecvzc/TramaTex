package auth

import (
	"context"
	"testing"
	"time"

	"github.com/joran-cortez/tramatex/internal/domain/security"
	"github.com/joran-cortez/tramatex/internal/domain/user"
)

// MockUserRepository is a test double for user.Repository.
type MockUserRepository struct {
	ByEmailFunc func(ctx context.Context, email *user.Email) (*user.User, error)
	ByIDFunc    func(ctx context.Context, id string) (*user.User, error)
	SaveFunc    func(ctx context.Context, u *user.User) error
	DeleteFunc  func(ctx context.Context, id string) error
}

func (m *MockUserRepository) ByEmail(ctx context.Context, email *user.Email) (*user.User, error) {
	return m.ByEmailFunc(ctx, email)
}

func (m *MockUserRepository) ByID(ctx context.Context, id string) (*user.User, error) {
	return m.ByIDFunc(ctx, id)
}

func (m *MockUserRepository) Save(ctx context.Context, u *user.User) error {
	return m.SaveFunc(ctx, u)
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	return m.DeleteFunc(ctx, id)
}

// MockJWTService is a test double for security.JWTService.
type MockJWTService struct {
	GenerateAccessTokenFunc  func(ctx context.Context, claims *security.TokenClaims) (string, error)
	GenerateRefreshTokenFunc func(ctx context.Context, claims *security.TokenClaims) (string, error)
	ValidateTokenFunc        func(ctx context.Context, token string) (*security.TokenClaims, error)
}

func (m *MockJWTService) GenerateAccessToken(ctx context.Context, claims *security.TokenClaims) (string, error) {
	return m.GenerateAccessTokenFunc(ctx, claims)
}

func (m *MockJWTService) GenerateRefreshToken(ctx context.Context, claims *security.TokenClaims) (string, error) {
	return m.GenerateRefreshTokenFunc(ctx, claims)
}

func (m *MockJWTService) ValidateToken(ctx context.Context, token string) (*security.TokenClaims, error) {
	return m.ValidateTokenFunc(ctx, token)
}

// Test: Successful login with valid credentials
func TestLoginWithValidCredentials(t *testing.T) {
	ctx := context.Background()

	// Setup: Create test user
	email, _ := user.NewEmail("user@example.com")
	password, _ := user.NewPassword("validPassword123")
	testUser, _ := user.NewUser("user-123", email, password, user.RoleOperator)

	// Mock repository
	mockRepo := &MockUserRepository{
		ByEmailFunc: func(ctx context.Context, e *user.Email) (*user.User, error) {
			if e.Equals(email) {
				return testUser, nil
			}
			return nil, nil
		},
	}

	// Mock JWT service
	mockJWT := &MockJWTService{
		GenerateAccessTokenFunc: func(ctx context.Context, claims *security.TokenClaims) (string, error) {
			return "access_token_xyz", nil
		},
		GenerateRefreshTokenFunc: func(ctx context.Context, claims *security.TokenClaims) (string, error) {
			return "refresh_token_abc", nil
		},
	}

	useCase := NewLoginUseCase(mockRepo, mockJWT)
	input := LoginInput{
		Email:    "user@example.com",
		Password: "validPassword123",
	}

	output, err := useCase.Execute(ctx, input)

	if err != nil {
		t.Errorf("Execute with valid credentials failed: %v", err)
	}

	if output == nil {
		t.Error("Execute should return output")
	}

	if output.User.ID != "user-123" {
		t.Errorf("User ID mismatch: %q", output.User.ID)
	}

	if output.AccessToken != "access_token_xyz" {
		t.Errorf("Access token mismatch: %q", output.AccessToken)
	}

	if output.RefreshToken != "refresh_token_abc" {
		t.Errorf("Refresh token mismatch: %q", output.RefreshToken)
	}

	if output.ExpiresIn != 900 {
		t.Errorf("ExpiresIn should be 900, got %d", output.ExpiresIn)
	}
}

// Test: Login with user not found
func TestLoginWithUserNotFound(t *testing.T) {
	ctx := context.Background()

	// Mock repository that returns error
	mockRepo := &MockUserRepository{
		ByEmailFunc: func(ctx context.Context, e *user.Email) (*user.User, error) {
			return nil, userNotFoundError()
		},
	}

	mockJWT := &MockJWTService{}

	useCase := NewLoginUseCase(mockRepo, mockJWT)
	input := LoginInput{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}

	output, err := useCase.Execute(ctx, input)

	if err == nil {
		t.Error("Execute should fail when user not found")
	}

	if output != nil {
		t.Error("Execute should return nil output on error")
	}
}

// Test: Login with invalid password
func TestLoginWithInvalidPassword(t *testing.T) {
	ctx := context.Background()

	// Setup: Create test user
	email, _ := user.NewEmail("user@example.com")
	password, _ := user.NewPassword("correctPassword123")
	testUser, _ := user.NewUser("user-123", email, password, user.RoleOperator)

	// Mock repository
	mockRepo := &MockUserRepository{
		ByEmailFunc: func(ctx context.Context, e *user.Email) (*user.User, error) {
			if e.Equals(email) {
				return testUser, nil
			}
			return nil, nil
		},
	}

	mockJWT := &MockJWTService{}

	useCase := NewLoginUseCase(mockRepo, mockJWT)
	input := LoginInput{
		Email:    "user@example.com",
		Password: "wrongPassword456",
	}

	output, err := useCase.Execute(ctx, input)

	if err == nil {
		t.Error("Execute should fail with wrong password")
	}

	if output != nil {
		t.Error("Execute should return nil output on error")
	}
}

// Test: Login with empty email
func TestLoginWithEmptyEmail(t *testing.T) {
	ctx := context.Background()

	mockRepo := &MockUserRepository{}
	mockJWT := &MockJWTService{}

	useCase := NewLoginUseCase(mockRepo, mockJWT)
	input := LoginInput{
		Email:    "",
		Password: "password123",
	}

	output, err := useCase.Execute(ctx, input)

	if err == nil {
		t.Error("Execute should fail with empty email")
	}

	if output != nil {
		t.Error("Execute should return nil output on error")
	}
}

// Test: Login with empty password
func TestLoginWithEmptyPassword(t *testing.T) {
	ctx := context.Background()

	mockRepo := &MockUserRepository{}
	mockJWT := &MockJWTService{}

	useCase := NewLoginUseCase(mockRepo, mockJWT)
	input := LoginInput{
		Email:    "user@example.com",
		Password: "",
	}

	output, err := useCase.Execute(ctx, input)

	if err == nil {
		t.Error("Execute should fail with empty password")
	}

	if output != nil {
		t.Error("Execute should return nil output on error")
	}
}

// Test: Login output DTO mapping
func TestLoginOutputDTOMapping(t *testing.T) {
	ctx := context.Background()

	// Setup: Create test user with specific data
	email, _ := user.NewEmail("john@example.com")
	password, _ := user.NewPassword("securePassword")
	testUser, _ := user.NewUser("user-john-123", email, password, user.RoleManager)

	// Mock repository
	mockRepo := &MockUserRepository{
		ByEmailFunc: func(ctx context.Context, e *user.Email) (*user.User, error) {
			return testUser, nil
		},
	}

	// Mock JWT service
	mockJWT := &MockJWTService{
		GenerateAccessTokenFunc: func(ctx context.Context, claims *security.TokenClaims) (string, error) {
			return "token_access", nil
		},
		GenerateRefreshTokenFunc: func(ctx context.Context, claims *security.TokenClaims) (string, error) {
			return "token_refresh", nil
		},
	}

	useCase := NewLoginUseCase(mockRepo, mockJWT)
	input := LoginInput{
		Email:    "john@example.com",
		Password: "securePassword",
	}

	output, err := useCase.Execute(ctx, input)

	if err != nil {
		t.Errorf("Execute failed: %v", err)
	}

	// Verify DTO mapping
	if output.User.ID != "user-john-123" {
		t.Errorf("User ID mapping failed: %q", output.User.ID)
	}

	if output.User.Email != "john@example.com" {
		t.Errorf("User Email mapping failed: %q", output.User.Email)
	}

	if output.User.Role != "manager" {
		t.Errorf("User Role mapping failed: %q", output.User.Role)
	}
}

// Helper: Create error for user not found
func userNotFoundError() error {
	return userNotFoundErr{}
}

type userNotFoundErr struct{}

func (e userNotFoundErr) Error() string {
	return "user not found"
}

// Test: Login with inactive user
func TestLoginWithInactiveUser(t *testing.T) {
	ctx := context.Background()

	// Setup: Create inactive test user
	email, _ := user.NewEmail("inactive@example.com")
	password, _ := user.NewPassword("password123")
	testUser, _ := user.NewUser("user-inactive", email, password, user.RoleOperator)
	testUser.Deactivate()

	// Mock repository
	mockRepo := &MockUserRepository{
		ByEmailFunc: func(ctx context.Context, e *user.Email) (*user.User, error) {
			if e.Equals(email) {
				return testUser, nil
			}
			return nil, nil
		},
	}

	mockJWT := &MockJWTService{}

	useCase := NewLoginUseCase(mockRepo, mockJWT)
	input := LoginInput{
		Email:    "inactive@example.com",
		Password: "password123",
	}

	output, err := useCase.Execute(ctx, input)

	if err == nil {
		t.Error("Execute should fail for inactive user")
	}

	if output != nil {
		t.Error("Execute should return nil output on error")
	}
}
