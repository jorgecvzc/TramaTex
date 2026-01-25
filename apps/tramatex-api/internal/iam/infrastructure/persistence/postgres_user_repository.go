package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/joran-cortez/tramatex/internal/iam/domain/model"
	"github.com/joran-cortez/tramatex/internal/iam/domain/repository"
	"gorm.io/gorm"
)

// PostgresUserRepository implementa el contrato repository.Repository usando PostgreSQL
type PostgresUserRepository struct {
	db *gorm.DB
}

// NewPostgresUserRepository crea una nueva instancia del repository
func NewPostgresUserRepository(db *gorm.DB) repository.Repository {
	return &PostgresUserRepository{db: db}
}

// ByID retrieves a user by their unique identifier
// Returns model.ErrUserNotFound if user does not exist or is deleted
func (r *PostgresUserRepository) ByID(ctx context.Context, id string) (*model.User, error) {
	var dbModel UserModel

	if err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&dbModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to fetch user by id: %w", err)
	}

	return r.modelToDomain(&dbModel)
}

// ByEmail retrieves a user by their email address
// Returns model.ErrUserNotFound if user does not exist or is deleted
func (r *PostgresUserRepository) ByEmail(ctx context.Context, email *model.Email) (*model.User, error) {
	var dbModel UserModel

	if err := r.db.WithContext(ctx).
		Where("email = ? AND deleted_at IS NULL", email.Value()).
		First(&dbModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to fetch user by email: %w", err)
	}

	return r.modelToDomain(&dbModel)
}

// Save persists a user to storage
// Creates a new user if id does not exist, updates if it does (upsert)
func (r *PostgresUserRepository) Save(ctx context.Context, u *model.User) error {
	if u == nil {
		return fmt.Errorf("user cannot be nil")
	}

	dbModel := UserModel{
		ID:           u.ID(),
		Email:        u.Email().Value(),
		PasswordHash: u.Password().Hash(),
		Role:         string(u.Role()),
		IsActive:     u.IsActive(),
	}

	// GORM upsert: si ID existe, actualiza; sino, crea
	if err := r.db.WithContext(ctx).Save(&dbModel).Error; err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	return nil
}

// Delete performs a soft delete of a user
// Sets deleted_at timestamp instead of removing the record
func (r *PostgresUserRepository) Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("user id cannot be empty")
	}

	if err := r.db.WithContext(ctx).
		Model(&UserModel{}).
		Where("id = ?", id).
		Update("deleted_at", gorm.Expr("CURRENT_TIMESTAMP")).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// modelToDomain convierte un UserModel (BD) a User (domain)
// Recrea los value objects: Email, Password desde el modelo
func (r *PostgresUserRepository) modelToDomain(dbModel *UserModel) (*model.User, error) {
	if dbModel == nil {
		return nil, fmt.Errorf("user model cannot be nil")
	}

	// Crear email value object
	email, err := model.NewEmail(dbModel.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email in database: %w", err)
	}

	// Crear password value object desde hash (no hashear de nuevo)
	password := model.NewPasswordFromHash(dbModel.PasswordHash)

	// Crear user entity con los value objects
	return model.NewUser(dbModel.ID, email, password, model.Role(dbModel.Role))
}
