package repository

import (
	"context"

	"github.com/joran-cortez/tramatex/internal/iam/domain/model"
)

// RoleRepository defines the contract for role persistence.
type RoleRepository interface {
	ByID(ctx context.Context, id string) (model.Role, error)
	ByName(ctx context.Context, name string) (model.Role, error)
	Save(ctx context.Context, role model.Role) error
}
