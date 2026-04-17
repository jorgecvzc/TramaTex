package domain

import (
	"context"

	"github.com/google/uuid"
)

type TaskFilters struct {
	IsActive *bool
	Search   string
}

type PositionFilters struct {
	IsActive *bool
	Search   string
}

type WorkTypeFilters struct {
	IsActive *bool
	Search   string
}

type WorkSetupFilters struct {
	IsActive *bool
	Search   string
	PartyID  string
}

type WorkOrderFilters struct {
	Status        *ProductionStatus
	ExcludeStatus *ProductionStatus
	Search        string
	PartyID       string
	WorkSetupID   *uuid.UUID
}

type TaskRepository interface {
	Save(ctx context.Context, task *Task) error
	FindByID(ctx context.Context, id uuid.UUID) (*Task, error)
	FindAll(ctx context.Context, filters *TaskFilters) ([]*Task, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
}

type PositionRepository interface {
	Save(ctx context.Context, position *Position) error
	FindByID(ctx context.Context, id uuid.UUID) (*Position, error)
	FindAll(ctx context.Context, filters *PositionFilters) ([]*Position, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type WorkTypeRepository interface {
	Save(ctx context.Context, workType *WorkType) error
	FindByID(ctx context.Context, id uuid.UUID) (*WorkType, error)
	FindAll(ctx context.Context, filters *WorkTypeFilters) ([]*WorkType, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type WorkSetupRepository interface {
	Save(ctx context.Context, workSetup *WorkSetup) error
	FindByID(ctx context.Context, id uuid.UUID) (*WorkSetup, error)
	FindAll(ctx context.Context, filters *WorkSetupFilters) ([]*WorkSetup, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type WorkOrderRepository interface {
	Save(ctx context.Context, workOrder *WorkOrder) error
	FindByID(ctx context.Context, id uuid.UUID) (*WorkOrder, error)
	FindAll(ctx context.Context, filters *WorkOrderFilters) ([]*WorkOrder, error)
	CountByYear(ctx context.Context, year int) (int64, error)
}
