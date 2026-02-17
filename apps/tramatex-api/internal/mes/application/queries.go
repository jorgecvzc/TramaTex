package application

import (
	"time"

	"github.com/google/uuid"
)

// ============================================================================
// PRODUCTION RECIPE QUERIES
// ============================================================================

// ListProductionRecipesQuery lists recipes with optional filters
type ListProductionRecipesQuery struct {
	ClientID   *uuid.UUID `form:"clientId"`
	ProductID  *uuid.UUID `form:"productId"`
	RecipeType *string    `form:"recipeType"`
	IsMaster   *bool      `form:"isMaster"`
	PageNumber int        `form:"page_number" binding:"min=1"`
	PageSize   int        `form:"page_size" binding:"min=1,max=100"`
}

// GetProductionRecipeByIDQuery gets a single recipe by ID
type GetProductionRecipeByIDQuery struct {
	ID uuid.UUID `uri:"id" binding:"required"`
}

// ============================================================================
// PRODUCTION ORDER QUERIES
// ============================================================================

// ListProductionOrdersQuery lists production orders with filters
type ListProductionOrdersQuery struct {
	SalesOrderID *uuid.UUID `form:"salesOrderId"`
	RecipeID     *uuid.UUID `form:"recipeId"`
	Status       *string    `form:"status"`
	WorkCenterID *uuid.UUID `form:"workCenterId"`
	FromDate     *time.Time `form:"fromDate" time_format:"2006-01-02"`
	ToDate       *time.Time `form:"toDate" time_format:"2006-01-02"`
	PageNumber   int        `form:"page_number" binding:"min=1"`
	PageSize     int        `form:"page_size" binding:"min=1,max=100"`
}

// GetProductionOrderByIDQuery gets a single production order by ID
type GetProductionOrderByIDQuery struct {
	ID uuid.UUID `uri:"id" binding:"required"`
}

// ============================================================================
// WORK CENTER QUERIES
// ============================================================================

// ListWorkCentersQuery lists all work centers
type ListWorkCentersQuery struct {
	IsActive   *bool `form:"isActive"`
	PageNumber int   `form:"page_number" binding:"min=1"`
	PageSize   int   `form:"page_size" binding:"min=1,max=100"`
}

// GetWorkCenterByIDQuery gets a single work center by ID
type GetWorkCenterByIDQuery struct {
	ID uuid.UUID `uri:"id" binding:"required"`
}
