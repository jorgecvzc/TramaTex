package iam_test

import (
	"context"
	"errors"
	"testing"

	iam_usecase "github.com/joran-cortez/tramatex/internal/iam/application/usecase"
	iam_model "github.com/joran-cortez/tramatex/internal/iam/domain/model"
	"github.com/joran-cortez/tramatex/internal/shared/infrastructure/security"
)

// MockUserRepository mocks the UserRepository interface
type MockUserRepository struct {
	user *iam_model.User
	err  error
}

func (m *MockUserRepository) ByID(ctx context.Context, id string) (*iam_model.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.user != nil && m.user.ID() == id {
		return m.user, nil
	}
	return nil, iam_model.ErrUserNotFound
}

func (m *MockUserRepository) ByEmail(ctx context.Context, email *iam_model.Email) (*iam_model.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.user != nil && m.user.Email().Equals(email) {
		return m.user, nil
	}
	return nil, iam_model.ErrUserNotFound
}

func (m *MockUserRepository) Save(ctx context.Context, user *iam_model.User) error {
	if m.err != nil {
		return m.err
	}
	m.user = user
	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	if m.err != nil {
		return m.err
	}
	if m.user != nil && m.user.ID() == id {
		m.user = nil
	}
	return nil
}

// MockJWTService mocks the JWTService interface
type MockJWTService struct {
	token string
	err   error
}

func (m *MockJWTService) GenerateAccessToken(ctx context.Context, claims *security.TokenClaims) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return "valid-access-token", nil
}

func (m *MockJWTService) GenerateRefreshToken(ctx context.Context, claims *security.TokenClaims) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return "valid-refresh-token", nil
}

func (m *MockJWTService) ValidateToken(ctx context.Context, tokenString string) (*security.TokenClaims, error) {
	if m.err != nil {
		return nil, m.err
	}
	return nil, errors.New("not implemented")
}

func TestLoginUseCaseWithValidCredentials(t *testing.T) {
	// Arrange
	email, _ := iam_model.NewEmail("user@example.com")
	password, _ := iam_model.NewPassword("SecurePassword123!")
	testUser, _ := iam_model.NewUser("user-123", email, password, iam_model.RoleOperator)

	mockRepo := &MockUserRepository{
		user: testUser,
	}

	mockJWT := &MockJWTService{}

	usecase := iam_usecase.NewLoginUseCase(mockRepo, mockJWT)

	// Act
	result, err := usecase.Execute(context.Background(), iam_usecase.LoginInput{
		Email:    "user@example.com",
		Password: "SecurePassword123!",
	})

	// Assert
	if err != nil {
		t.Errorf("LoginUseCase.Execute with valid credentials should not fail: %v", err)
	}

	if result == nil {
		t.Fatal("LoginUseCase.Execute should return LoginOutput")
	}

	if result.AccessToken != "valid-access-token" {
		t.Errorf("Token = %q, want %q", result.AccessToken, "valid-access-token")
	}
	if result.User.Email != "user@example.com" {
		t.Errorf("User Email = %q, want %q", result.User.Email, "user@example.com")
	}
}

func TestLoginUseCaseWithUserNotFound(t *testing.T) {
	// Arrange
	mockRepo := &MockUserRepository{
		user: nil, // User does not exist
	}
	mockJWT := &MockJWTService{}

	usecase := iam_usecase.NewLoginUseCase(mockRepo, mockJWT)

	// Act
	_, err := usecase.Execute(context.Background(), iam_usecase.LoginInput{
		Email:    "nonexistent@example.com",
		Password: "SomePassword123!",
	})

	// Assert
	if err == nil {
		t.Fatal("LoginUseCase.Execute with nonexistent user should fail")
	}
	if !errors.Is(err, iam_model.ErrUserNotFound) {
		t.Errorf("Expected error %v, got %v", iam_model.ErrUserNotFound, err)
	}
}

func TestLoginUseCaseWithWrongPassword(t *testing.T) {
	// Arrange
	email, _ := iam_model.NewEmail("user@example.com")
	password, _ := iam_model.NewPassword("CorrectPassword123!")
	testUser, _ := iam_model.NewUser("user-123", email, password, iam_model.RoleOperator)

	mockRepo := &MockUserRepository{
		user: testUser,
	}
	mockJWT := &MockJWTService{}

	usecase := iam_usecase.NewLoginUseCase(mockRepo, mockJWT)

	// Act
	_, err := usecase.Execute(context.Background(), iam_usecase.LoginInput{
		Email:    "user@example.com",
		Password: "WrongPassword123!",
	})

	// Assert
	if err == nil {
		t.Fatal("LoginUseCase.Execute with wrong password should fail")
	}
	if !errors.Is(err, iam_model.ErrInvalidPassword) {
		t.Errorf("Expected error %v, got %v", iam_model.ErrInvalidPassword, err)
	}
}
