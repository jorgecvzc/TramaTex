package repository

import (
	"context"

	"github.com/joran-cortez/tramatex/internal/iam/domain/model"
)

// PermissionRepository defines the contract for permission persistence.
type PermissionRepository interface {
	ByID(ctx context.Context, id string) (*model.Permission, error)
	ByName(ctx context.Context, name string) (*model.Permission, error)
	Save(ctx context.Context, permission *model.Permission) error
}
