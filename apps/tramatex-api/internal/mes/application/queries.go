package application

import (
	"github.com/google/uuid"
)

type ListTasksQuery struct {
	IsActive *bool  `form:"is_active"`
	Search   string `form:"search"`
}

type GetTaskByIDQuery struct {
	ID uuid.UUID `uri:"id" binding:"required"`
}

type ListPositionsQuery struct {
	IsActive *bool  `form:"is_active"`
	Search   string `form:"search"`
}

type GetPositionByIDQuery struct {
	ID uuid.UUID `uri:"id" binding:"required"`
}

type ListWorkTypesQuery struct {
	IsActive *bool  `form:"is_active"`
	Search   string `form:"search"`
}

type GetWorkTypeByIDQuery struct {
	ID uuid.UUID `uri:"id" binding:"required"`
}

type ListWorkOrdersQuery struct {
	Status      *string `form:"status"`
	Search      string  `form:"search"`
	PartyID     string  `form:"party_id"`
	WorkSetupID string  `form:"work_setup_id"`
}

type GetWorkOrderByIDQuery struct {
	ID uuid.UUID `uri:"id" binding:"required"`
}

type ListOverdueWorkOrdersQuery struct {
	Limit int `form:"limit"`
}

// --- WorkSetup Queries ---

type ListWorkSetupsQuery struct {
	IsActive *bool  `form:"is_active"`
	Search   string `form:"search"`
	PartyID  string `form:"party_id"`
}

type GetWorkSetupByIDQuery struct {
	ID uuid.UUID `uri:"id" binding:"required"`
}
