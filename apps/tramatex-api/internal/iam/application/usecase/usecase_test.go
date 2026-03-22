package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/joran-cortez/tramatex/internal/iam/application/usecase"
	"github.com/joran-cortez/tramatex/internal/iam/domain/model"
	"github.com/joran-cortez/tramatex/internal/shared/infrastructure/security"
)

const testUserIDStr = "00000000-0000-0000-0000-000000000001"

type fakeUserRepo struct {
	byIDFunc    func(ctx context.Context, id uuid.UUID) (*model.User, error)
	byEmailFunc func(ctx context.Context, email *model.Email) (*model.User, error)
	saveFunc    func(ctx context.Context, user *model.User) error
	deleteFunc  func(ctx context.Context, id uuid.UUID) error
	listFunc    func(ctx context.Context) ([]*model.User, error)
	savedUsers  []*model.User
}

func (f *fakeUserRepo) ByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	if f.byIDFunc != nil {
		return f.byIDFunc(ctx, id)
	}
	return nil, model.ErrUserNotFound
}

func (f *fakeUserRepo) ByEmail(ctx context.Context, email *model.Email) (*model.User, error) {
	if f.byEmailFunc != nil {
		return f.byEmailFunc(ctx, email)
	}
	return nil, model.ErrUserNotFound
}

func (f *fakeUserRepo) Save(ctx context.Context, user *model.User) error {
	if f.saveFunc != nil {
		return f.saveFunc(ctx, user)
	}
	f.savedUsers = append(f.savedUsers, user)
	return nil
}

func (f *fakeUserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if f.deleteFunc != nil {
		return f.deleteFunc(ctx, id)
	}
	return nil
}

func (f *fakeUserRepo) List(ctx context.Context) ([]*model.User, error) {
	if f.listFunc != nil {
		return f.listFunc(ctx)
	}
	return f.savedUsers, nil
}

type fakeJWTService struct {
	accessToken  string
	refreshToken string
	validateFunc func(ctx context.Context, token string) (*security.TokenClaims, error)
	accessErr    error
	refreshErr   error
}

func (f *fakeJWTService) GenerateAccessToken(ctx context.Context, claims *security.TokenClaims) (string, error) {
	if f.accessErr != nil {
		return "", f.accessErr
	}
	if f.accessToken == "" {
		return "access-token", nil
	}
	return f.accessToken, nil
}

func (f *fakeJWTService) GenerateRefreshToken(ctx context.Context, claims *security.TokenClaims) (string, error) {
	if f.refreshErr != nil {
		return "", f.refreshErr
	}
	if f.refreshToken == "" {
		return "refresh-token", nil
	}
	return f.refreshToken, nil
}

func (f *fakeJWTService) ValidateToken(ctx context.Context, token string) (*security.TokenClaims, error) {
	if f.validateFunc != nil {
		return f.validateFunc(ctx, token)
	}
	return nil, security.ErrInvalidToken
}

func newUser(t *testing.T, email string, password string, role model.Role) *model.User {
	t.Helper()

	emailVO, err := model.NewEmail(email)
	if err != nil {
		t.Fatalf("invalid email: %v", err)
	}
	passwordVO, err := model.NewPassword(password)
	if err != nil {
		t.Fatalf("invalid password: %v", err)
	}
	user, err := model.NewUserWithUUID(emailVO, passwordVO, role)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	return user
}

func TestCreateUserUseCase(t *testing.T) {
	repo := &fakeUserRepo{}
	uc := usecase.NewCreateUserUseCase(repo)

	_, err := uc.Execute(context.Background(), usecase.CreateUserInput{})
	if err == nil {
		t.Fatalf("expected validation error")
	}

	repo.byEmailFunc = func(ctx context.Context, email *model.Email) (*model.User, error) {
		return newUser(t, email.Value(), "password123", model.RoleAdmin), nil
	}
	_, err = uc.Execute(context.Background(), usecase.CreateUserInput{Email: "admin@example.com", Password: "password123", Role: "admin"})
	if !errors.Is(err, model.ErrUserAlreadyExists) {
		t.Fatalf("expected user already exists error")
	}

	repo.byEmailFunc = func(ctx context.Context, email *model.Email) (*model.User, error) {
		return nil, model.ErrUserNotFound
	}
	output, err := uc.Execute(context.Background(), usecase.CreateUserInput{Email: "new@example.com", Password: "password123", Role: "admin"})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if output.Role != "admin" {
		t.Fatalf("unexpected role: %s", output.Role)
	}
}

func TestRegisterUserUseCase(t *testing.T) {
	repo := &fakeUserRepo{}
	uc := usecase.NewRegisterUserUseCase(repo)

	_, err := uc.Execute(context.Background(), usecase.RegisterInput{})
	if err == nil {
		t.Fatalf("expected validation error")
	}

	repo.byEmailFunc = func(ctx context.Context, email *model.Email) (*model.User, error) {
		return newUser(t, email.Value(), "password123", model.RoleCommercial), nil
	}
	_, err = uc.Execute(context.Background(), usecase.RegisterInput{Email: "user@example.com", Password: "password123"})
	if !errors.Is(err, model.ErrUserAlreadyExists) {
		t.Fatalf("expected user already exists error")
	}

	repo.byEmailFunc = func(ctx context.Context, email *model.Email) (*model.User, error) {
		return nil, model.ErrUserNotFound
	}
	output, err := uc.Execute(context.Background(), usecase.RegisterInput{Email: "user@example.com", Password: "password123"})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if output.Role != string(model.RoleCommercial) {
		t.Fatalf("expected default role commercial")
	}
}

func TestLoginUseCase(t *testing.T) {
	repo := &fakeUserRepo{}
	jwtService := &fakeJWTService{}
	uc := usecase.NewLoginUseCase(repo, jwtService)

	_, err := uc.Execute(context.Background(), usecase.LoginInput{})
	if err == nil {
		t.Fatalf("expected validation error")
	}

	repo.byEmailFunc = func(ctx context.Context, email *model.Email) (*model.User, error) {
		return nil, model.ErrUserNotFound
	}
	_, err = uc.Execute(context.Background(), usecase.LoginInput{Email: "user@example.com", Password: "password123"})
	if err == nil {
		t.Fatalf("expected user not found error")
	}

	user := newUser(t, "user@example.com", "password123", model.RoleCommercial)
	user.Deactivate()
	repo.byEmailFunc = func(ctx context.Context, email *model.Email) (*model.User, error) {
		return user, nil
	}
	_, err = uc.Execute(context.Background(), usecase.LoginInput{Email: "user@example.com", Password: "password123"})
	if err == nil {
		t.Fatalf("expected inactive user error")
	}

	user.Activate()
	repo.byEmailFunc = func(ctx context.Context, email *model.Email) (*model.User, error) {
		return user, nil
	}
	_, err = uc.Execute(context.Background(), usecase.LoginInput{Email: "user@example.com", Password: "wrong"})
	if !errors.Is(err, model.ErrInvalidPassword) {
		t.Fatalf("expected invalid password error")
	}

	output, err := uc.Execute(context.Background(), usecase.LoginInput{Email: "user@example.com", Password: "password123"})
	if err != nil {
		t.Fatalf("expected login success, got error: %v", err)
	}
	if output.AccessToken == "" || output.RefreshToken == "" {
		t.Fatalf("expected tokens")
	}
}

func TestRefreshTokenUseCase(t *testing.T) {
	claims, _ := security.NewTokenClaims(testUserIDStr, "user@example.com", "commercial", time.Now(), time.Now().Add(time.Hour))
	repo := &fakeUserRepo{}
	jwtService := &fakeJWTService{
		validateFunc: func(ctx context.Context, token string) (*security.TokenClaims, error) {
			return claims, nil
		},
	}
	uc := usecase.NewRefreshTokenUseCase(repo, jwtService)

	_, err := uc.Execute(context.Background(), usecase.RefreshInput{})
	if err == nil {
		t.Fatalf("expected validation error")
	}

	jwtService.validateFunc = func(ctx context.Context, token string) (*security.TokenClaims, error) {
		return nil, security.ErrInvalidToken
	}
	_, err = uc.Execute(context.Background(), usecase.RefreshInput{RefreshToken: "bad"})
	if err == nil {
		t.Fatalf("expected invalid token error")
	}

	jwtService.validateFunc = func(ctx context.Context, token string) (*security.TokenClaims, error) {
		return claims, nil
	}
	repo.byIDFunc = func(ctx context.Context, id uuid.UUID) (*model.User, error) {
		return nil, model.ErrUserNotFound
	}
	_, err = uc.Execute(context.Background(), usecase.RefreshInput{RefreshToken: "refresh"})
	if !errors.Is(err, model.ErrUserNotFound) {
		t.Fatalf("expected user not found error")
	}

	inactive := newUser(t, "user@example.com", "password123", model.RoleCommercial)
	inactive.Deactivate()
	repo.byIDFunc = func(ctx context.Context, id uuid.UUID) (*model.User, error) {
		return inactive, nil
	}
	_, err = uc.Execute(context.Background(), usecase.RefreshInput{RefreshToken: "refresh"})
	if err == nil {
		t.Fatalf("expected inactive user error")
	}

	active := newUser(t, "user@example.com", "password123", model.RoleCommercial)
	repo.byIDFunc = func(ctx context.Context, id uuid.UUID) (*model.User, error) {
		return active, nil
	}
	output, err := uc.Execute(context.Background(), usecase.RefreshInput{RefreshToken: "refresh"})
	if err != nil {
		t.Fatalf("expected refresh success, got error: %v", err)
	}
	if output.AccessToken == "" {
		t.Fatalf("expected access token")
	}
}

func TestAssignRoleUseCase(t *testing.T) {
	repo := &fakeUserRepo{}
	uc := usecase.NewAssignRoleUseCase(repo)

	_, err := uc.Execute(context.Background(), usecase.AssignRoleInput{})
	if err == nil {
		t.Fatalf("expected validation error")
	}

	_, err = uc.Execute(context.Background(), usecase.AssignRoleInput{UserID: "user", Role: "invalid"})
	if err == nil {
		t.Fatalf("expected invalid role error")
	}

	repo.byIDFunc = func(ctx context.Context, id uuid.UUID) (*model.User, error) {
		return nil, model.ErrUserNotFound
	}
	_, err = uc.Execute(context.Background(), usecase.AssignRoleInput{UserID: testUserIDStr, Role: "admin"})
	if !errors.Is(err, model.ErrUserNotFound) {
		t.Fatalf("expected user not found error")
	}

	repo.byIDFunc = func(ctx context.Context, id uuid.UUID) (*model.User, error) {
		return newUser(t, "user@example.com", "password123", model.RoleCommercial), nil
	}
	output, err := uc.Execute(context.Background(), usecase.AssignRoleInput{UserID: testUserIDStr, Role: "admin"})
	if err != nil {
		t.Fatalf("expected assign role success, got error: %v", err)
	}
	if output.Role != "admin" {
		t.Fatalf("unexpected role: %s", output.Role)
	}
}

func TestCheckAuthorizationUseCase(t *testing.T) {
	repo := &fakeUserRepo{}
	uc := usecase.NewCheckAuthorizationUseCase(repo)

	_, err := uc.Execute(context.Background(), usecase.CheckAuthorizationInput{})
	if err == nil {
		t.Fatalf("expected validation error")
	}

	_, err = uc.Execute(context.Background(), usecase.CheckAuthorizationInput{RequiredRoles: []string{"admin"}})
	if err == nil {
		t.Fatalf("expected user id required error")
	}

	repo.byIDFunc = func(ctx context.Context, id uuid.UUID) (*model.User, error) {
		return nil, model.ErrUserNotFound
	}
	_, err = uc.Execute(context.Background(), usecase.CheckAuthorizationInput{UserID: testUserIDStr, RequiredRoles: []string{"admin"}})
	if !errors.Is(err, model.ErrUserNotFound) {
		t.Fatalf("expected user not found error")
	}

	user := newUser(t, "user@example.com", "password123", model.RoleDesigner)
	repo.byIDFunc = func(ctx context.Context, id uuid.UUID) (*model.User, error) {
		return user, nil
	}
	output, err := uc.Execute(context.Background(), usecase.CheckAuthorizationInput{UserID: testUserIDStr, RequiredRoles: []string{"admin", "designer"}})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if !output.Allowed {
		t.Fatalf("expected authorization allowed")
	}

	output, err = uc.Execute(context.Background(), usecase.CheckAuthorizationInput{UserID: testUserIDStr, RequiredRoles: []string{"admin"}})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if output.Allowed {
		t.Fatalf("expected authorization denied")
	}
}

func TestListUsersUseCase(t *testing.T) {
	repo := &fakeUserRepo{}
	uc := usecase.NewListUsersUseCase(repo)

	repo.listFunc = func(ctx context.Context) ([]*model.User, error) {
		return nil, errors.New("db error")
	}
	_, err := uc.Execute(context.Background())
	if err == nil {
		t.Fatalf("expected list error")
	}

	repo.listFunc = func(ctx context.Context) ([]*model.User, error) {
		return []*model.User{newUser(t, "user@example.com", "password123", model.RoleCommercial)}, nil
	}
	users, err := uc.Execute(context.Background())
	if err != nil {
		t.Fatalf("expected list success, got error: %v", err)
	}
	if len(users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(users))
	}
}

func TestDeleteUserUseCase(t *testing.T) {
	repo := &fakeUserRepo{}
	uc := usecase.NewDeleteUserUseCase(repo)

	if err := uc.Execute(context.Background(), ""); err == nil {
		t.Fatalf("expected validation error")
	}

	repo.deleteFunc = func(ctx context.Context, id uuid.UUID) error {
		return errors.New("delete failed")
	}
	if err := uc.Execute(context.Background(), testUserIDStr); err == nil {
		t.Fatalf("expected delete error")
	}

	repo.deleteFunc = func(ctx context.Context, id uuid.UUID) error {
		return nil
	}
	if err := uc.Execute(context.Background(), testUserIDStr); err != nil {
		t.Fatalf("expected delete success, got error: %v", err)
	}
}

func TestLogoutUserUseCase(t *testing.T) {
	claims, _ := security.NewTokenClaims("user-1", "user@example.com", "admin", time.Now(), time.Now().Add(time.Hour))
	jwtService := &fakeJWTService{
		validateFunc: func(ctx context.Context, token string) (*security.TokenClaims, error) {
			return claims, nil
		},
	}
	blacklist := security.NewInMemoryTokenBlacklist()
	uc := usecase.NewLogoutUserUseCase(jwtService, blacklist)

	if err := uc.Execute(context.Background(), "", usecase.LogoutInput{}); err == nil {
		t.Fatalf("expected validation error")
	}

	if err := uc.Execute(context.Background(), "access-token", usecase.LogoutInput{RefreshToken: "refresh-token"}); err != nil {
		t.Fatalf("expected logout success, got error: %v", err)
	}
	if !blacklist.IsRevoked("access-token") {
		t.Fatalf("expected access token to be revoked")
	}
	if !blacklist.IsRevoked("refresh-token") {
		t.Fatalf("expected refresh token to be revoked")
	}
}
