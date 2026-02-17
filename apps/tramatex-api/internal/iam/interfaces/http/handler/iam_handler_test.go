package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joran-cortez/tramatex/internal/iam/application/usecase"
	"github.com/joran-cortez/tramatex/internal/iam/domain/model"
	"github.com/joran-cortez/tramatex/internal/shared/infrastructure/logging"
	"github.com/joran-cortez/tramatex/internal/shared/infrastructure/security"
)

func init() {
	logging.InitLogger("test")
}

type fakeUserRepo struct {
	byIDFunc    func(ctx context.Context, id string) (*model.User, error)
	byEmailFunc func(ctx context.Context, email *model.Email) (*model.User, error)
	saveFunc    func(ctx context.Context, user *model.User) error
	deleteFunc  func(ctx context.Context, id string) error
	listFunc    func(ctx context.Context) ([]*model.User, error)
}

func (f *fakeUserRepo) ByID(ctx context.Context, id string) (*model.User, error) {
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
	return nil
}

func (f *fakeUserRepo) Delete(ctx context.Context, id string) error {
	if f.deleteFunc != nil {
		return f.deleteFunc(ctx, id)
	}
	return nil
}

func (f *fakeUserRepo) List(ctx context.Context) ([]*model.User, error) {
	if f.listFunc != nil {
		return f.listFunc(ctx)
	}
	return nil, nil
}

type fakeJWTService struct {
	validateFunc func(ctx context.Context, token string) (*security.TokenClaims, error)
}

func (f *fakeJWTService) GenerateAccessToken(ctx context.Context, claims *security.TokenClaims) (string, error) {
	return "access-token", nil
}

func (f *fakeJWTService) GenerateRefreshToken(ctx context.Context, claims *security.TokenClaims) (string, error) {
	return "refresh-token", nil
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

func newIAMHandler(repo *fakeUserRepo, jwt *fakeJWTService, blacklist security.TokenBlacklist) *IAMHandler {
	return NewIAMHandler(
		usecase.NewLoginUseCase(repo, jwt),
		usecase.NewRegisterUserUseCase(repo),
		usecase.NewCreateUserUseCase(repo),
		usecase.NewRefreshTokenUseCase(repo, jwt),
		usecase.NewLogoutUserUseCase(jwt, blacklist),
		usecase.NewAssignRoleUseCase(repo),
		usecase.NewCheckAuthorizationUseCase(repo),
		usecase.NewListUsersUseCase(repo),
		usecase.NewDeleteUserUseCase(repo),
	)
}

func performRequest(t *testing.T, handlerFunc func(*gin.Context), method, path string, body interface{}, headers map[string]string, setup func(*gin.Context)) *httptest.ResponseRecorder {
	t.Helper()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	var reqBody []byte
	switch v := body.(type) {
	case string:
		reqBody = []byte(v)
	case nil:
		reqBody = nil
	default:
		data, err := json.Marshal(v)
		if err != nil {
			t.Fatalf("failed to marshal body: %v", err)
		}
		reqBody = data
	}

	req := httptest.NewRequest(method, path, bytes.NewReader(reqBody))
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	c.Request = req

	if setup != nil {
		setup(c)
	}

	handlerFunc(c)
	c.Writer.WriteHeaderNow()
	return w
}

func TestIAMHandler_Login(t *testing.T) {
	repo := &fakeUserRepo{}
	jwt := &fakeJWTService{}
	h := newIAMHandler(repo, jwt, nil)

	response := performRequest(t, h.Login, http.MethodPost, "/auth/login", "{", nil, nil)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", response.Code)
	}

	repo.byEmailFunc = func(ctx context.Context, email *model.Email) (*model.User, error) {
		return nil, model.ErrUserNotFound
	}
	response = performRequest(t, h.Login, http.MethodPost, "/auth/login", usecase.LoginInput{Email: "user@example.com", Password: "password"}, nil, nil)
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", response.Code)
	}

	user := newUser(t, "user@example.com", "password", model.RoleCommercial)
	repo.byEmailFunc = func(ctx context.Context, email *model.Email) (*model.User, error) {
		return user, nil
	}
	response = performRequest(t, h.Login, http.MethodPost, "/auth/login", usecase.LoginInput{Email: "user@example.com", Password: "password"}, nil, nil)
	if response.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", response.Code)
	}
}

func TestIAMHandler_Register(t *testing.T) {
	repo := &fakeUserRepo{}
	jwt := &fakeJWTService{}
	h := newIAMHandler(repo, jwt, nil)

	response := performRequest(t, h.Register, http.MethodPost, "/auth/register", "{", nil, nil)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", response.Code)
	}

	repo.byEmailFunc = func(ctx context.Context, email *model.Email) (*model.User, error) {
		return newUser(t, email.Value(), "password123", model.RoleCommercial), nil
	}
	response = performRequest(t, h.Register, http.MethodPost, "/auth/register", usecase.RegisterInput{Email: "user@example.com", Password: "password123"}, nil, nil)
	if response.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d", response.Code)
	}

	repo.byEmailFunc = func(ctx context.Context, email *model.Email) (*model.User, error) {
		return nil, model.ErrUserNotFound
	}
	response = performRequest(t, h.Register, http.MethodPost, "/auth/register", usecase.RegisterInput{Email: "user@example.com", Password: "password123"}, nil, nil)
	if response.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", response.Code)
	}
}

func TestIAMHandler_CreateUser(t *testing.T) {
	repo := &fakeUserRepo{}
	jwt := &fakeJWTService{}
	h := newIAMHandler(repo, jwt, nil)

	response := performRequest(t, h.CreateUser, http.MethodPost, "/auth/users", "{", nil, nil)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", response.Code)
	}

	repo.byEmailFunc = func(ctx context.Context, email *model.Email) (*model.User, error) {
		return newUser(t, email.Value(), "password123", model.RoleAdmin), nil
	}
	response = performRequest(t, h.CreateUser, http.MethodPost, "/auth/users", usecase.CreateUserInput{Email: "user@example.com", Password: "password123", Role: "admin"}, nil, nil)
	if response.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d", response.Code)
	}

	repo.byEmailFunc = func(ctx context.Context, email *model.Email) (*model.User, error) {
		return nil, model.ErrUserNotFound
	}
	response = performRequest(t, h.CreateUser, http.MethodPost, "/auth/users", usecase.CreateUserInput{Email: "user@example.com", Password: "password123", Role: "admin"}, nil, nil)
	if response.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", response.Code)
	}
}

func TestIAMHandler_Refresh(t *testing.T) {
	repo := &fakeUserRepo{}
	jwt := &fakeJWTService{}
	h := newIAMHandler(repo, jwt, nil)

	response := performRequest(t, h.Refresh, http.MethodPost, "/auth/refresh", "{", nil, nil)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", response.Code)
	}

	jwt.validateFunc = func(ctx context.Context, token string) (*security.TokenClaims, error) {
		return nil, security.ErrInvalidToken
	}
	response = performRequest(t, h.Refresh, http.MethodPost, "/auth/refresh", usecase.RefreshInput{RefreshToken: "bad"}, nil, nil)
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", response.Code)
	}

	claims, _ := security.NewTokenClaims("user-1", "user@example.com", "admin", time.Now(), time.Now().Add(time.Hour))
	jwt.validateFunc = func(ctx context.Context, token string) (*security.TokenClaims, error) {
		return claims, nil
	}
	repo.byIDFunc = func(ctx context.Context, id string) (*model.User, error) {
		return newUser(t, "user@example.com", "password123", model.RoleAdmin), nil
	}
	response = performRequest(t, h.Refresh, http.MethodPost, "/auth/refresh", usecase.RefreshInput{RefreshToken: "valid"}, nil, nil)
	if response.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", response.Code)
	}
}

func TestIAMHandler_Logout(t *testing.T) {
	repo := &fakeUserRepo{}
	jwt := &fakeJWTService{}
	blacklist := security.NewInMemoryTokenBlacklist()
	h := newIAMHandler(repo, jwt, blacklist)

	response := performRequest(t, h.Logout, http.MethodPost, "/auth/logout", usecase.LogoutInput{}, nil, nil)
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", response.Code)
	}

	response = performRequest(t, h.Logout, http.MethodPost, "/auth/logout", usecase.LogoutInput{}, map[string]string{"Authorization": "Token bad"}, nil)
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", response.Code)
	}

	jwt.validateFunc = func(ctx context.Context, token string) (*security.TokenClaims, error) {
		return nil, security.ErrInvalidToken
	}
	response = performRequest(t, h.Logout, http.MethodPost, "/auth/logout", usecase.LogoutInput{}, map[string]string{"Authorization": "Bearer bad"}, nil)
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", response.Code)
	}

	claims, _ := security.NewTokenClaims("user-1", "user@example.com", "admin", time.Now(), time.Now().Add(time.Hour))
	jwt.validateFunc = func(ctx context.Context, token string) (*security.TokenClaims, error) {
		return claims, nil
	}
	response = performRequest(t, h.Logout, http.MethodPost, "/auth/logout", usecase.LogoutInput{RefreshToken: "refresh"}, map[string]string{"Authorization": "Bearer access"}, nil)
	if response.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", response.Code)
	}
}

func TestIAMHandler_AssignRole(t *testing.T) {
	repo := &fakeUserRepo{}
	jwt := &fakeJWTService{}
	h := newIAMHandler(repo, jwt, nil)

	response := performRequest(t, h.AssignRole, http.MethodPost, "/auth/assign-role", "{", nil, nil)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", response.Code)
	}

	repo.byIDFunc = func(ctx context.Context, id string) (*model.User, error) {
		return nil, model.ErrUserNotFound
	}
	response = performRequest(t, h.AssignRole, http.MethodPost, "/auth/assign-role", usecase.AssignRoleInput{UserID: "user", Role: "admin"}, nil, nil)
	if response.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", response.Code)
	}

	repo.byIDFunc = func(ctx context.Context, id string) (*model.User, error) {
		return newUser(t, "user@example.com", "password123", model.RoleCommercial), nil
	}
	response = performRequest(t, h.AssignRole, http.MethodPost, "/auth/assign-role", usecase.AssignRoleInput{UserID: "user", Role: "admin"}, nil, nil)
	if response.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", response.Code)
	}
}

func TestIAMHandler_CheckAuthorization(t *testing.T) {
	repo := &fakeUserRepo{}
	jwt := &fakeJWTService{}
	h := newIAMHandler(repo, jwt, nil)

	response := performRequest(t, h.CheckAuthorization, http.MethodPost, "/auth/authorize", "{", nil, nil)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", response.Code)
	}

	repo.byIDFunc = func(ctx context.Context, id string) (*model.User, error) {
		return nil, model.ErrUserNotFound
	}
	response = performRequest(t, h.CheckAuthorization, http.MethodPost, "/auth/authorize", usecase.CheckAuthorizationInput{UserID: "user", RequiredRoles: []string{"admin"}}, nil, nil)
	if response.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", response.Code)
	}

	repo.byIDFunc = func(ctx context.Context, id string) (*model.User, error) {
		return newUser(t, "user@example.com", "password123", model.RoleAdmin), nil
	}
	response = performRequest(t, h.CheckAuthorization, http.MethodPost, "/auth/authorize", usecase.CheckAuthorizationInput{RequiredRoles: []string{"admin"}}, nil, func(c *gin.Context) {
		c.Set("userID", "user")
	})
	if response.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", response.Code)
	}
}

func TestIAMHandler_DeleteUser(t *testing.T) {
	repo := &fakeUserRepo{}
	jwt := &fakeJWTService{}
	h := newIAMHandler(repo, jwt, nil)

	response := performRequest(t, h.DeleteUser, http.MethodDelete, "/auth/users", nil, nil, func(c *gin.Context) {
		c.Params = []gin.Param{{Key: "id", Value: ""}}
	})
	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", response.Code)
	}

	repo.deleteFunc = func(ctx context.Context, id string) error {
		return errors.New("delete failed")
	}
	response = performRequest(t, h.DeleteUser, http.MethodDelete, "/auth/users/1", nil, nil, func(c *gin.Context) {
		c.Params = []gin.Param{{Key: "id", Value: "1"}}
	})
	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", response.Code)
	}

	repo.deleteFunc = func(ctx context.Context, id string) error {
		return nil
	}
	response = performRequest(t, h.DeleteUser, http.MethodDelete, "/auth/users/1", nil, nil, func(c *gin.Context) {
		c.Params = []gin.Param{{Key: "id", Value: "1"}}
	})
	if response.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", response.Code)
	}
}

func TestIAMHandler_ListUsers(t *testing.T) {
	repo := &fakeUserRepo{}
	jwt := &fakeJWTService{}
	h := newIAMHandler(repo, jwt, nil)

	repo.listFunc = func(ctx context.Context) ([]*model.User, error) {
		return nil, errors.New("db error")
	}
	response := performRequest(t, h.ListUsers, http.MethodGet, "/auth/users", nil, nil, nil)
	if response.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", response.Code)
	}

	repo.listFunc = func(ctx context.Context) ([]*model.User, error) {
		return []*model.User{newUser(t, "user@example.com", "password123", model.RoleCommercial)}, nil
	}
	response = performRequest(t, h.ListUsers, http.MethodGet, "/auth/users", nil, nil, nil)
	if response.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", response.Code)
	}
}

func TestIAMHandler_Health(t *testing.T) {
	repo := &fakeUserRepo{}
	jwt := &fakeJWTService{}
	h := newIAMHandler(repo, jwt, nil)

	response := performRequest(t, h.Health, http.MethodGet, "/health", nil, nil, nil)
	if response.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", response.Code)
	}
}
