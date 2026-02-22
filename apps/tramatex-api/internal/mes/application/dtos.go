package application

import (
	"time"

	"github.com/google/uuid"
)

type TaskDTO struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	IsActive    bool      `json:"is_active"`
}

type PositionDTO struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description,omitempty"`
	IsActive    bool      `json:"is_active"`
}

type ServiceGroupTaskDTO struct {
	TaskID   uuid.UUID `json:"task_id"`
	Sequence int       `json:"sequence"`
}

type ServiceGroupDTO struct {
	ID             uuid.UUID             `json:"id"`
	Name           string                `json:"name"`
	Description    string                `json:"description,omitempty"`
	ProductGroupID *uuid.UUID            `json:"product_group_id,omitempty"`
	IsActive       bool                  `json:"is_active"`
	Tasks          []ServiceGroupTaskDTO `json:"tasks"`
}

type MESWorkTaskDTO struct {
	ID          uuid.UUID  `json:"id"`
	TaskID      uuid.UUID  `json:"task_id"`
	Sequence    int        `json:"sequence"`
	Status      string     `json:"status"`
	AssignedTo  *uuid.UUID `json:"assigned_to,omitempty"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Notes       string     `json:"notes,omitempty"`
}

type MESWorkServiceGroupDTO struct {
	ID             uuid.UUID        `json:"id"`
	ServiceGroupID uuid.UUID        `json:"service_group_id"`
	PositionID     uuid.UUID        `json:"position_id"`
	DesignFilePath string           `json:"design_file_path,omitempty"`
	Notes          string           `json:"notes,omitempty"`
	Sequence       int              `json:"sequence"`
	Tasks          []MESWorkTaskDTO `json:"tasks"`
}

type MESWorkDTO struct {
	ID              uuid.UUID                `json:"id"`
	WorkNumber      string                   `json:"work_number"`
	WorkName        string                   `json:"work_name"`
	PartyID         string                   `json:"party_id"`
	TangibleGroupID uuid.UUID                `json:"tangible_group_id"`
	GarmentNotes    string                   `json:"garment_notes,omitempty"`
	Status          string                   `json:"status"`
	Priority        string                   `json:"priority"`
	StartDate       *time.Time               `json:"start_date,omitempty"`
	DueDate         *time.Time               `json:"due_date,omitempty"`
	CompletedDate   *time.Time               `json:"completed_date,omitempty"`
	ServiceGroups   []MESWorkServiceGroupDTO `json:"service_groups"`
}

type MESWorkDashboardStatsDTO struct {
	Total    int            `json:"total"`
	ByStatus map[string]int `json:"by_status"`
	Overdue  int            `json:"overdue"`
	DueToday int            `json:"due_today"`
}
