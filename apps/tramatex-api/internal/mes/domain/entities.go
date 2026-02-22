package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID          uuid.UUID
	Name        string
	Description string
	IsActive    bool
}

func NewTask(name, description string, isActive bool) (*Task, error) {
	if name == "" {
		return nil, fmt.Errorf("task name is required")
	}

	return &Task{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		IsActive:    isActive,
	}, nil
}

type Position struct {
	ID          uuid.UUID
	Name        string
	Code        string
	Description string
	IsActive    bool
}

func NewPosition(name, code, description string, isActive bool) (*Position, error) {
	if name == "" {
		return nil, fmt.Errorf("position name is required")
	}
	if code == "" {
		return nil, fmt.Errorf("position code is required")
	}

	return &Position{
		ID:          uuid.New(),
		Name:        name,
		Code:        code,
		Description: description,
		IsActive:    isActive,
	}, nil
}

type ServiceGroupTask struct {
	TaskID   uuid.UUID
	Sequence int
}

type ServiceGroup struct {
	ID             uuid.UUID
	Name           string
	Description    string
	ProductGroupID *uuid.UUID
	IsActive       bool
	Tasks          []ServiceGroupTask
}

func NewServiceGroup(name, description string, productGroupID *uuid.UUID, isActive bool, tasks []ServiceGroupTask) (*ServiceGroup, error) {
	if name == "" {
		return nil, fmt.Errorf("service group name is required")
	}

	for _, task := range tasks {
		if task.TaskID == uuid.Nil {
			return nil, fmt.Errorf("task ID is required")
		}
		if task.Sequence <= 0 {
			return nil, fmt.Errorf("task sequence must be greater than zero")
		}
	}

	return &ServiceGroup{
		ID:             uuid.New(),
		Name:           name,
		Description:    description,
		ProductGroupID: productGroupID,
		IsActive:       isActive,
		Tasks:          tasks,
	}, nil
}

type MESWorkTask struct {
	ID          uuid.UUID
	TaskID      uuid.UUID
	Sequence    int
	Status      TaskStatus
	AssignedTo  *uuid.UUID
	StartedAt   *time.Time
	CompletedAt *time.Time
	Notes       string
}

type MESWorkServiceGroup struct {
	ID             uuid.UUID
	ServiceGroupID uuid.UUID
	PositionID     uuid.UUID
	DesignFilePath string
	Notes          string
	Sequence       int
	Tasks          []MESWorkTask
}

type MESWork struct {
	ID              uuid.UUID
	WorkNumber      string
	WorkName        string
	PartyID         string
	TangibleGroupID uuid.UUID
	GarmentNotes    string
	Status          ProductionStatus
	Priority        WorkPriority
	StartDate       *time.Time
	DueDate         *time.Time
	CompletedDate   *time.Time
	ServiceGroups   []MESWorkServiceGroup
}

func NewMESWork(
	workNumber, workName, partyID string,
	tangibleGroupID uuid.UUID,
	garmentNotes string,
	status ProductionStatus,
	priority WorkPriority,
	startDate, dueDate, completedDate *time.Time,
	serviceGroups []MESWorkServiceGroup,
) (*MESWork, error) {
	if workNumber == "" {
		return nil, fmt.Errorf("work number is required")
	}
	if workName == "" {
		return nil, fmt.Errorf("work name is required")
	}
	if partyID == "" {
		return nil, fmt.Errorf("party id is required")
	}
	if tangibleGroupID == uuid.Nil {
		return nil, fmt.Errorf("tangible group id is required")
	}
	if !status.IsValid() {
		return nil, fmt.Errorf("invalid production status")
	}
	if !priority.IsValid() {
		return nil, fmt.Errorf("invalid work priority")
	}

	for _, group := range serviceGroups {
		if group.ServiceGroupID == uuid.Nil {
			return nil, fmt.Errorf("service group id is required")
		}
		if group.PositionID == uuid.Nil {
			return nil, fmt.Errorf("position id is required")
		}
		if group.Sequence <= 0 {
			return nil, fmt.Errorf("service group sequence must be greater than zero")
		}
		for _, task := range group.Tasks {
			if task.TaskID == uuid.Nil {
				return nil, fmt.Errorf("work task id is required")
			}
			if task.Sequence <= 0 {
				return nil, fmt.Errorf("work task sequence must be greater than zero")
			}
			if !task.Status.IsValid() {
				return nil, fmt.Errorf("invalid work task status")
			}
		}
	}

	return &MESWork{
		ID:              uuid.New(),
		WorkNumber:      workNumber,
		WorkName:        workName,
		PartyID:         partyID,
		TangibleGroupID: tangibleGroupID,
		GarmentNotes:    garmentNotes,
		Status:          status,
		Priority:        priority,
		StartDate:       startDate,
		DueDate:         dueDate,
		CompletedDate:   completedDate,
		ServiceGroups:   serviceGroups,
	}, nil
}
