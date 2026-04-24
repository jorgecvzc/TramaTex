package persistence

import (
	"context"
	"errors"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/joran-cortez/tramatex/internal/iam/domain/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	testID1 = "00000000-0000-0000-0000-000000000001"
	testID2 = "00000000-0000-0000-0000-000000000002"
)

func newMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	t.Helper()

	sqlDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		sqlDB.Close()
		t.Fatalf("failed to open gorm db: %v", err)
	}
	cleanup := func() {
		sqlDB.Close()
	}
	return gormDB, mock, cleanup
}

func TestPostgresUserRepository_ByID_NotFound(t *testing.T) {
	gormDB, mock, cleanup := newMockDB(t)
	defer cleanup()

	mock.ExpectQuery("SELECT .* FROM \"users\"").
		WithArgs(testID1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	repo := NewPostgresUserRepository(gormDB).(*PostgresUserRepository)
	user, err := repo.ByID(context.Background(), uuid.MustParse(testID1))
	if !errors.Is(err, model.ErrUserNotFound) {
		t.Fatalf("expected not found error, got %v", err)
	}
	if user != nil {
		t.Fatalf("expected nil user")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestPostgresUserRepository_ByEmail_NotFound(t *testing.T) {
	gormDB, mock, cleanup := newMockDB(t)
	defer cleanup()

	mock.ExpectQuery("SELECT .* FROM \"users\"").
		WithArgs("user@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	repo := NewPostgresUserRepository(gormDB).(*PostgresUserRepository)
	email, _ := model.NewEmail("user@example.com")
	user, err := repo.ByEmail(context.Background(), email)
	if !errors.Is(err, model.ErrUserNotFound) {
		t.Fatalf("expected not found error, got %v", err)
	}
	if user != nil {
		t.Fatalf("expected nil user")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestPostgresUserRepository_ByID_Success(t *testing.T) {
	gormDB, mock, cleanup := newMockDB(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "email", "password", "role", "is_active", "created_at", "updated_at", "deleted_at"}).
		AddRow(testID1, "user@example.com", "$2a$10$hash", "commercial", true, time.Now(), time.Now(), nil)

	mock.ExpectQuery("SELECT .* FROM \"users\"").
		WithArgs(testID1).
		WillReturnRows(rows)

	repo := NewPostgresUserRepository(gormDB).(*PostgresUserRepository)
	user, err := repo.ByID(context.Background(), uuid.MustParse(testID1))
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if user == nil || user.ID() != uuid.MustParse(testID1) {
		t.Fatalf("unexpected user result")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestPostgresUserRepository_List(t *testing.T) {
	gormDB, mock, cleanup := newMockDB(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "email", "password", "role", "is_active", "created_at", "updated_at", "deleted_at"}).
		AddRow(testID1, "user1@example.com", "$2a$10$hash", "commercial", true, time.Now(), time.Now(), nil).
		AddRow(testID2, "user2@example.com", "$2a$10$hash", "admin", true, time.Now(), time.Now(), nil)

	mock.ExpectQuery("SELECT .* FROM \"users\"").
		WillReturnRows(rows)

	repo := NewPostgresUserRepository(gormDB).(*PostgresUserRepository)
	users, err := repo.List(context.Background())
	if err != nil {
		t.Fatalf("expected list success, got error: %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestPostgresUserRepository_Delete_EmptyID(t *testing.T) {
	repo := &PostgresUserRepository{}
	if err := repo.Delete(context.Background(), uuid.Nil); err == nil {
		t.Fatalf("expected error for empty id")
	}
}

func TestPostgresUserRepository_Save_NilUser(t *testing.T) {
	repo := &PostgresUserRepository{}
	if err := repo.Save(context.Background(), nil); err == nil {
		t.Fatalf("expected error for nil user")
	}
}

func TestPostgresUserRepository_ModelToDomain_InvalidEmail(t *testing.T) {
	repo := &PostgresUserRepository{}
	_, err := repo.modelToDomain(&UserModel{ID: testID1, Email: "invalid", PasswordHash: "hash", Role: "admin"})
	if err == nil {
		t.Fatalf("expected error for invalid email")
	}
}

func TestPostgresUserRepository_ModelToDomain_NilModel(t *testing.T) {
	repo := &PostgresUserRepository{}
	_, err := repo.modelToDomain(nil)
	if err == nil {
		t.Fatalf("expected error for nil model")
	}
}

func TestPostgresUserRepository_Delete_Success(t *testing.T) {
	gormDB, mock, cleanup := newMockDB(t)
	defer cleanup()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE .*\"users\".*deleted_at").
		WithArgs(sqlmock.AnyArg(), testID1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := NewPostgresUserRepository(gormDB).(*PostgresUserRepository)
	if err := repo.Delete(context.Background(), uuid.MustParse(testID1)); err != nil {
		t.Fatalf("expected delete success, got error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
