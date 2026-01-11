package auth_test

import (
	"errors"
	"testing"
	"time"

	authapp "github.com/joran-cortez/tramatex/internal/application/auth"
	"github.com/joran-cortez/tramatex/internal/domain/security"
	"github.com/joran-cortez/tramatex/internal/domain/user"
)

// MockUserRepository mocks the UserRepository interface
type MockUserRepository struct {
	user   *user.User
	found  bool
	err    error
}

func (m *MockUserRepository) FindByEmail(email string) (*user.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	if !m.found {
		return nil, errors.New("user not found")
	}
	return m.user, nil
}

// MockJWTService mocks the JWTService interface
type MockJWTService struct {
	token string
	err   error
}

func (m *MockJWTService) Generate(claims *security.TokenClaims) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.token, nil
}

func (m *MockJWTService) Verify(tokenString string) (*security.TokenClaims, error) {
	// Not used in LoginUseCase, but required for interface
	return nil, errors.New("not implemented")
}

func TestLoginUseCaseWithValidCredentials(t *testing.T) {
	// Arrange
	email, _ := user.NewEmail("user@example.com")
	password, _ := user.NewPassword("SecurePassword123!")
	testUser, _ := user.NewUser(email, password, user.RoleUser)

	mockRepo := &MockUserRepository{
		user:  testUser,
		found: true,
	}

	mockJWT := &MockJWTService{
		token: "valid-jwt-token",
	}

	usecase := authapp.NewLoginUseCase(mockRepo, mockJWT)

	// Act
	result, err := usecase.Execute(&authapp.LoginRequest{
		Email:    "user@example.com",
		Password: "SecurePassword123!",
	})

	// Assert
	if err != nil {
		t.Errorf("LoginUseCase.Execute with valid credentials should not fail: %v", err)
	}

	if result == nil {
		t.Error("LoginUseCase.Execute should return LoginResponse")
	}

	if result.Token != "valid-jwt-token" {
		t.Errorf("Token = %q, want %q", result.Token, "valid-jwt-token")
	}
}

func TestLoginUseCaseWithInvalidEmail(t *testing.T) {
	// Arrange
	mockRepo := &MockUserRepository{}
	mockJWT := &MockJWTService{token: "token"}

	usecase := authapp.NewLoginUseCase(mockRepo, mockJWT)

	// Act
	result, err := usecase.Execute(&authapp.LoginRequest{
		Email:    "invalid-email-format",
		Password: "SomePassword123!",
	})

	// Assert
	if err == nil {
		t.Error("LoginUseCase.Execute with invalid email should fail")
	}

	if result != nil {
		t.Error("LoginUseCase.Execute should return nil LoginResponse on error")
	}
}

func TestLoginUseCaseWithUserNotFound(t *testing.T) {
	// Arrange
	mockRepo := &MockUserRepository{
		found: false,
		err:   errors.New("user not found"),
	}
	mockJWT := &MockJWTService{token: "token"}

	usecase := authapp.NewLoginUseCase(mockRepo, mockJWT)

	// Act
	result, err := usecase.Execute(&authapp.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "SomePassword123!",
	})

	// Assert
	if err == nil {
		t.Error("LoginUseCase.Execute with nonexistent user should fail")
	}

	if result != nil {
		t.Error("LoginUseCase.Execute should return nil LoginResponse on error")
	}
}

func TestLoginUseCaseWithWrongPassword(t *testing.T) {
	// Arrange
	email, _ := user.NewEmail("user@example.com")
	password, _ := user.NewPassword("CorrectPassword123!")
	testUser, _ := user.NewUser(email, password, user.RoleUser)

	mockRepo := &MockUserRepository{
		user:  testUser,
		found: true,
	}
	mockJWT := &MockJWTService{token: "token"}

	usecase := authapp.NewLoginUseCase(mockRepo, mockJWT)

	// Act
	result, err := usecase.Execute(&authapp.LoginRequest{
		Email:    "user@example.com",
		Password: "WrongPassword123!",
	})

	// Assert
	if err == nil {
		t.Error("LoginUseCase.Execute with wrong password should fail")
	}

	if result != nil {
		t.Error("LoginUseCase.Execute should return nil LoginResponse on error")
	}
}

func TestLoginUseCaseWithRepositoryError(t *testing.T) {
	// Arrange
	mockRepo := &MockUserRepository{
		err: errors.New("database connection error"),
	}
	mockJWT := &MockJWTService{token: "token"}

	usecase := authapp.NewLoginUseCase(mockRepo, mockJWT)

	// Act
	result, err := usecase.Execute(&authapp.LoginRequest{
		Email:    "user@example.com",
		Password: "Password123!",
	})

	// Assert
	if err == nil {
		t.Error("LoginUseCase.Execute should propagate repository errors")
	}

	if result != nil {
		t.Error("LoginUseCase.Execute should return nil LoginResponse on error")
	}
}

func TestLoginUseCaseWithJWTGenerationError(t *testing.T) {
	// Arrange
	email, _ := user.NewEmail("user@example.com")
	password, _ := user.NewPassword("SecurePassword123!")
	testUser, _ := user.NewUser(email, password, user.RoleUser)

	mockRepo := &MockUserRepository{
		user:  testUser,
		found: true,
	}

	mockJWT := &MockJWTService{
		err: errors.New("JWT generation failed"),
	}

	usecase := authapp.NewLoginUseCase(mockRepo, mockJWT)

	// Act
	result, err := usecase.Execute(&authapp.LoginRequest{
		Email:    "user@example.com",
		Password: "SecurePassword123!",
	})

	// Assert
	if err == nil {
		t.Error("LoginUseCase.Execute should return JWT generation errors")
	}

	if result != nil {
		t.Error("LoginUseCase.Execute should return nil LoginResponse on JWT error")
	}
}

func TestLoginUseCaseTokenGeneration(t *testing.T) {
	// Arrange
	email, _ := user.NewEmail("manager@example.com")
	password, _ := user.NewPassword("ManagerPass123!")
	testUser, _ := user.NewUser(email, password, user.RoleManager)

	mockRepo := &MockUserRepository{
		user:  testUser,
		found: true,
	}

	expectedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test.token"
	mockJWT := &MockJWTService{
		token: expectedToken,
	}

	usecase := authapp.NewLoginUseCase(mockRepo, mockJWT)

	// Act
	result, err := usecase.Execute(&authapp.LoginRequest{
		Email:    "manager@example.com",
		Password: "ManagerPass123!",
	})

	// Assert
	if err != nil {
		t.Errorf("LoginUseCase.Execute should succeed: %v", err)
	}

	if result.Token != expectedToken {
		t.Errorf("Token = %q, want %q", result.Token, expectedToken)
	}
}

func TestLoginUseCaseResponseStructure(t *testing.T) {
	// Arrange
	email, _ := user.NewEmail("user@example.com")
	password, _ := user.NewPassword("SecurePassword123!")
	testUser, _ := user.NewUser(email, password, user.RoleUser)

	mockRepo := &MockUserRepository{
		user:  testUser,
		found: true,
	}

	mockJWT := &MockJWTService{
		token: "generated-token",
	}

	usecase := authapp.NewLoginUseCase(mockRepo, mockJWT)

	// Act
	result, err := usecase.Execute(&authapp.LoginRequest{
		Email:    "user@example.com",
		Password: "SecurePassword123!",
	})

	// Assert
	if err != nil {
		t.Errorf("LoginUseCase.Execute failed: %v", err)
	}

	if result == nil {
		t.Fatal("LoginUseCase.Execute returned nil result")
	}

	// Response should contain token and potentially user info
	if result.Token == "" {
		t.Error("LoginResponse.Token should not be empty")
	}

	// Optional: Response should have metadata (timestamp, expires in, etc.)
	// Verify response is properly populated
	if !isValidTokenFormat(result.Token) {
		t.Errorf("Token format appears invalid: %q", result.Token)
	}
}

// Helper function to validate token format (basic check)
func isValidTokenFormat(token string) bool {
	return len(token) > 0
}
