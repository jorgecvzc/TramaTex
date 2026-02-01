package persistence

import (
	"context"
	"fmt"

	"github.com/joran-cortez/tramatex/internal/iam/domain/model"
	"github.com/joran-cortez/tramatex/internal/iam/domain/repository"
	"gorm.io/gorm"
)

// PostgresRoleRepository provides role lookups backed by a static set.
// NOTE: Roles are fixed for MVP and not stored in a separate table.
type PostgresRoleRepository struct {
	db *gorm.DB
}

// NewPostgresRoleRepository creates a new role repository instance.
func NewPostgresRoleRepository(db *gorm.DB) repository.RoleRepository {
	return &PostgresRoleRepository{db: db}
}

// ByID returns a role by id. For MVP, roles are fixed and id == name.
func (r *PostgresRoleRepository) ByID(ctx context.Context, id string) (model.Role, error) {
	return r.ByName(ctx, id)
}

// ByName returns a role by name.
func (r *PostgresRoleRepository) ByName(ctx context.Context, name string) (model.Role, error) {
	role, err := model.NewRole(name)
	if err != nil {
		return "", err
	}
	return role, nil
}

// Save is a no-op for MVP since roles are fixed.
func (r *PostgresRoleRepository) Save(ctx context.Context, role model.Role) error {
	if !role.IsValid() {
		return fmt.Errorf("invalid role: %s", role)
	}
	return nil
}
