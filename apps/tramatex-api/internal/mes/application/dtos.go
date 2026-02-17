package application

import (
	"time"

	"github.com/google/uuid"
)

// ============================================================================
// PRODUCTION RECIPE DTOs
// ============================================================================

// ProductionRecipeDTO represents a production recipe for API responses
type ProductionRecipeDTO struct {
	ID              uuid.UUID           `json:"id"`
	Name            string              `json:"name"`
	ClientID        uuid.UUID           `json:"clientId"`
	ProductID       uuid.UUID           `json:"productId"`
	RecipeType      string              `json:"recipeType"`
	Version         int                 `json:"version"`
	IsMaster        bool                `json:"isMaster"`
	TaskDefinitions []TaskDefinitionDTO `json:"taskDefinitions"`
}

// TaskDefinitionDTO represents a task definition
type TaskDefinitionDTO struct {
	ID                  uuid.UUID  `json:"id"`
	Name                string     `json:"name"`
	Description         string     `json:"description"`
	SequenceOrder       int        `json:"sequenceOrder"`
	EstimatedDurationHs float64    `json:"estimatedDurationHs"`
	WorkCenterID        *uuid.UUID `json:"workCenterId,omitempty"`
}

// ============================================================================
// PRODUCTION ORDER DTOs
// ============================================================================

// ProductionOrderDTO represents a production order for API responses
type ProductionOrderDTO struct {
	ID                     uuid.UUID         `json:"id"`
	SalesOrderID           uuid.UUID         `json:"salesOrderId"`
	RecipeID               uuid.UUID         `json:"recipeId"`
	ProductID              uuid.UUID         `json:"productId"`
	Quantity               int               `json:"quantity"`
	Status                 string            `json:"status"`
	StartDate              time.Time         `json:"startDate"`
	EndDate                time.Time         `json:"endDate"`
	AssignedToWorkCenterID *uuid.UUID        `json:"assignedToWorkCenterId,omitempty"`
	TaskInstances          []TaskInstanceDTO `json:"taskInstances"`
}

// TaskInstanceDTO represents a task instance
type TaskInstanceDTO struct {
	ID                  uuid.UUID  `json:"id"`
	TaskDefinitionID    uuid.UUID  `json:"taskDefinitionId"`
	Name                string     `json:"name"`
	Description         string     `json:"description"`
	SequenceOrder       int        `json:"sequenceOrder"`
	Status              string     `json:"status"`
	EstimatedDurationHs float64    `json:"estimatedDurationHs"`
	ActualStartTime     *time.Time `json:"actualStartTime,omitempty"`
	ActualEndTime       *time.Time `json:"actualEndTime,omitempty"`
	AssignedOperatorID  *uuid.UUID `json:"assignedOperatorId,omitempty"`
	WorkCenterID        *uuid.UUID `json:"workCenterId,omitempty"`
	Notes               string     `json:"notes,omitempty"`
}

// ============================================================================
// WORK CENTER DTOs
// ============================================================================

// WorkCenterDTO represents a work center for API responses
type WorkCenterDTO struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsActive    bool      `json:"isActive"`
}

// ============================================================================
// PAGINATED RESPONSES
// ============================================================================

// PaginatedProductionRecipesResponse represents a paginated list of recipes
type PaginatedProductionRecipesResponse struct {
	Data       []ProductionRecipeDTO `json:"data"`
	PageNumber int                   `json:"page_number"`
	PageSize   int                   `json:"page_size"`
	Total      int64                 `json:"total"`
}

// PaginatedProductionOrdersResponse represents a paginated list of orders
type PaginatedProductionOrdersResponse struct {
	Data       []ProductionOrderDTO `json:"data"`
	PageNumber int                  `json:"page_number"`
	PageSize   int                  `json:"page_size"`
	Total      int64                `json:"total"`
}

// PaginatedWorkCentersResponse represents a paginated list of work centers
type PaginatedWorkCentersResponse struct {
	Data       []WorkCenterDTO `json:"data"`
	PageNumber int             `json:"page_number"`
	PageSize   int             `json:"page_size"`
	Total      int64           `json:"total"`
}
