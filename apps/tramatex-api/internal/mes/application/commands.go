package application

import "github.com/google/uuid"

type CreateTaskCommand struct {
	ActorID     string  `json:"-"`
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
	IsActive    *bool   `json:"is_active"`
}

type UpdateTaskCommand struct {
	ActorID     string `json:"-"`
	ID          uuid.UUID
	Name        *string `json:"name"`
	Description *string `json:"description"`
	IsActive    *bool   `json:"is_active"`
}

type DeleteTaskCommand struct {
	ActorID string
	ID      uuid.UUID
}

type CreatePositionCommand struct {
	ActorID     string  `json:"-"`
	Name        string  `json:"name" binding:"required"`
	Code        string  `json:"code" binding:"required"`
	Description *string `json:"description"`
	IsActive    *bool   `json:"is_active"`
}

type UpdatePositionCommand struct {
	ActorID     string `json:"-"`
	ID          uuid.UUID
	Name        *string `json:"name"`
	Code        *string `json:"code"`
	Description *string `json:"description"`
	IsActive    *bool   `json:"is_active"`
}

type DeletePositionCommand struct {
	ActorID string
	ID      uuid.UUID
}

type ServiceGroupTaskInput struct {
	TaskID   uuid.UUID `json:"task_id" binding:"required"`
	Sequence int       `json:"sequence" binding:"required,min=1"`
}

type CreateServiceGroupCommand struct {
	ActorID         string                  `json:"-"`
	Name            string                  `json:"name" binding:"required"`
	Description     *string                 `json:"description"`
	ProductGroupID  *uuid.UUID              `json:"product_group_id"`
	IsActive        *bool                   `json:"is_active"`
	TaskAssignments []ServiceGroupTaskInput `json:"task_assignments"`
}

type UpdateServiceGroupCommand struct {
	ActorID         string `json:"-"`
	ID              uuid.UUID
	Name            *string                 `json:"name"`
	Description     *string                 `json:"description"`
	ProductGroupID  *uuid.UUID              `json:"product_group_id"`
	IsActive        *bool                   `json:"is_active"`
	TaskAssignments []ServiceGroupTaskInput `json:"task_assignments"`
}

type DeleteServiceGroupCommand struct {
	ActorID string
	ID      uuid.UUID
}

type CreateMESWorkServiceGroupInput struct {
	ServiceGroupID uuid.UUID `json:"service_group_id" binding:"required"`
	PositionID     uuid.UUID `json:"position_id" binding:"required"`
	DesignFilePath *string   `json:"design_file_path"`
	Notes          *string   `json:"notes"`
	Sequence       int       `json:"sequence" binding:"required,min=1"`
}

type CreateMESWorkCommand struct {
	ActorID                 string                           `json:"-"`
	WorkName                string                           `json:"work_name" binding:"required"`
	PartyID                 string                           `json:"party_id" binding:"required"`
	TangibleGroupID         uuid.UUID                        `json:"tangible_group_id" binding:"required"`
	GarmentNotes            *string                          `json:"garment_notes"`
	Status                  *string                          `json:"status"`
	Priority                *string                          `json:"priority"`
	ServiceGroupAssignments []CreateMESWorkServiceGroupInput `json:"service_group_assignments" binding:"required,min=1"`
}

type UpdateMESWorkTaskStatusCommand struct {
	ActorID string `json:"-"`
	WorkID  uuid.UUID
	TaskID  uuid.UUID
	Action  string  `json:"action" binding:"required"`
	Notes   *string `json:"notes"`
}
