package application

import "github.com/google/uuid"

type CreateTaskCommand struct {
	ActorID     string  `json:"-"`
	Name        string  `json:"name" binding:"required"`
	Reference   *string `json:"reference"`
	Description *string `json:"description"`
	IsActive    *bool   `json:"is_active"`
}

type UpdateTaskCommand struct {
	ActorID     string `json:"-"`
	ID          uuid.UUID
	Name        *string `json:"name"`
	Reference   *string `json:"reference"`
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

type WorkTypeTaskInput struct {
	TaskID   uuid.UUID `json:"task_id" binding:"required"`
	Sequence int       `json:"sequence" binding:"required,min=1"`
}

type CreateWorkTypeCommand struct {
	ActorID         string              `json:"-"`
	Name            string              `json:"name" binding:"required"`
	Reference       *string             `json:"reference"`
	Description     *string             `json:"description"`
	IsActive        *bool               `json:"is_active"`
	TaskAssignments []WorkTypeTaskInput `json:"task_assignments"`
}

type UpdateWorkTypeCommand struct {
	ActorID         string `json:"-"`
	ID              uuid.UUID
	Name            *string             `json:"name"`
	Reference       *string             `json:"reference"`
	Description     *string             `json:"description"`
	IsActive        *bool               `json:"is_active"`
	TaskAssignments []WorkTypeTaskInput `json:"task_assignments"`
}

type DeleteWorkTypeCommand struct {
	ActorID string
	ID      uuid.UUID
}

type CreateWorkOrderCommand struct {
	ActorID          string     `json:"-"`
	WorkName         string     `json:"work_name" binding:"required"`
	PartyID          string     `json:"party_id" binding:"required"`
	WorkSetupID      *uuid.UUID `json:"work_setup_id"`
	Notes            *string    `json:"notes"`
	Priority         *string    `json:"priority"`
	Status           *string    `json:"status"`
	DueDate          *string    `json:"due_date"`
	OrderWorkSetupID *uuid.UUID `json:"order_work_setup_id,omitempty"`
}

type UpdateWorkOrderCommand struct {
	ActorID     string `json:"-"`
	ID          uuid.UUID
	WorkName    *string    `json:"work_name"`
	Notes       *string    `json:"notes"`
	Status      *string    `json:"status"`
	Priority    *string    `json:"priority"`
	DueDate     *string    `json:"due_date"`
	WorkSetupID *uuid.UUID `json:"work_setup_id"`
}

type UpdateWorkOrderTaskStatusCommand struct {
	ActorID string `json:"-"`
	WorkID  uuid.UUID
	TaskID  uuid.UUID
	Action  string  `json:"action" binding:"required"`
	Notes   *string `json:"notes"`
}

// --- WorkSetup Commands ---

type WorkSetupLineInput struct {
	WorkTypeID     uuid.UUID `json:"work_type_id" binding:"required"`
	PositionID     uuid.UUID `json:"position_id" binding:"required"`
	DesignFilePath *string   `json:"design_file_path"`
	Notes          *string   `json:"notes"`
	Sequence       int       `json:"sequence" binding:"required,min=1"`
}

type CreateWorkSetupCommand struct {
	ActorID         string               `json:"-"`
	Name            string               `json:"name" binding:"required"`
	PartyID         string               `json:"party_id" binding:"required"`
	TangibleGroupID *uuid.UUID           `json:"tangible_group_id"`
	Description     *string              `json:"description"`
	IsActive        *bool                `json:"is_active"`
	Lines           []WorkSetupLineInput `json:"lines"`
}

type UpdateWorkSetupCommand struct {
	ActorID         string `json:"-"`
	ID              uuid.UUID
	Name            *string              `json:"name"`
	PartyID         *string              `json:"party_id"`
	TangibleGroupID *uuid.UUID           `json:"tangible_group_id"`
	Description     *string              `json:"description"`
	IsActive        *bool                `json:"is_active"`
	Lines           []WorkSetupLineInput `json:"lines"`
}

type DeleteWorkSetupCommand struct {
	ActorID string
	ID      uuid.UUID
}
