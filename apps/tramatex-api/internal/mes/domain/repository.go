package domain

import (
	"context"

	"github.com/google/uuid"
)

// ProductionRecipeRepository defines methods for ProductionRecipe aggregate persistence.
type ProductionRecipeRepository interface {
	Save(ctx context.Context, recipe *ProductionRecipe, createdBy string, modifiedBy string) error
	FindByID(ctx context.Context, id uuid.UUID) (*ProductionRecipe, error)
	FindByClientIDAndProductID(ctx context.Context, clientID, productID uuid.UUID) (*ProductionRecipe, error)
	FindAll(ctx context.Context) ([]*ProductionRecipe, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// ProductionOrderRepository defines methods for ProductionOrder aggregate persistence.
type ProductionOrderRepository interface {
	Save(ctx context.Context, order *ProductionOrder, createdBy string, modifiedBy string) error
	FindByID(ctx context.Context, id uuid.UUID) (*ProductionOrder, error)
	FindBySalesOrderID(ctx context.Context, salesOrderID uuid.UUID) (*ProductionOrder, error)
	FindAll(ctx context.Context) ([]*ProductionOrder, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// TaskInstanceRepository defines methods for TaskInstance entity persistence.
type TaskInstanceRepository interface {
	Save(ctx context.Context, task *TaskInstance, createdBy string, modifiedBy string) error
	FindByID(ctx context.Context, id uuid.UUID) (*TaskInstance, error)
	FindByProductionOrderID(ctx context.Context, productionOrderID uuid.UUID) ([]*TaskInstance, error)
	FindAll(ctx context.Context) ([]*TaskInstance, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// WorkCenterRepository defines methods for WorkCenter entity persistence.
type WorkCenterRepository interface {
	Save(ctx context.Context, workCenter *WorkCenter, createdBy string, modifiedBy string) error
	FindByID(ctx context.Context, id uuid.UUID) (*WorkCenter, error)
	FindAll(ctx context.Context) ([]*WorkCenter, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// QualityCheckRepository defines methods for QualityCheck entity persistence.
type QualityCheckRepository interface {
	Save(ctx context.Context, check *QualityCheck, createdBy string, modifiedBy string) error
	FindByID(ctx context.Context, id uuid.UUID) (*QualityCheck, error)
	FindByProductionOrderID(ctx context.Context, productionOrderID uuid.UUID) ([]*QualityCheck, error)
	FindAll(ctx context.Context) ([]*QualityCheck, error)
	Delete(ctx context.Context, id uuid.UUID) error
}