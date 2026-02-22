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

type ServiceGroupFilters struct {
	IsActive *bool
	Search   string
}

type MESWorkFilters struct {
	Status  *ProductionStatus
	Search  string
	PartyID string
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

type ServiceGroupRepository interface {
	Save(ctx context.Context, serviceGroup *ServiceGroup) error
	FindByID(ctx context.Context, id uuid.UUID) (*ServiceGroup, error)
	FindAll(ctx context.Context, filters *ServiceGroupFilters) ([]*ServiceGroup, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type MESWorkRepository interface {
	Save(ctx context.Context, work *MESWork) error
	FindByID(ctx context.Context, id uuid.UUID) (*MESWork, error)
	FindAll(ctx context.Context, filters *MESWorkFilters) ([]*MESWork, error)
	CountByYear(ctx context.Context, year int) (int64, error)
}
