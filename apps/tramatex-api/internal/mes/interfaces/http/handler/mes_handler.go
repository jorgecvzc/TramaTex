package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joran-cortez/tramatex/internal/mes/application"
)

// MESHandler handles HTTP requests for the MES module
type MESHandler struct {
	service *application.MESService
}

// NewMESHandler creates a new MES handler
func NewMESHandler(service *application.MESService) *MESHandler {
	return &MESHandler{
		service: service,
	}
}

// ============================================================================
// PRODUCTION RECIPE HANDLERS
// ============================================================================

// CreateProductionRecipe creates a new production recipe
// POST /api/mes/recipes
func (h *MESHandler) CreateProductionRecipe(c *gin.Context) {
	var cmd application.CreateProductionRecipeCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recipe, err := h.service.CreateProductionRecipe(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, recipe)
}

// GetProductionRecipe retrieves a single production recipe
// GET /api/mes/recipes/:id
func (h *MESHandler) GetProductionRecipe(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid recipe ID"})
		return
	}

	query := application.GetProductionRecipeByIDQuery{ID: id}
	recipe, err := h.service.GetProductionRecipeByID(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, recipe)
}

// ListProductionRecipes lists all production recipes with optional filters
// GET /api/mes/recipes
func (h *MESHandler) ListProductionRecipes(c *gin.Context) {
	var query application.ListProductionRecipesQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set defaults
	if query.PageNumber == 0 {
		query.PageNumber = 1
	}
	if query.PageSize == 0 {
		query.PageSize = 20
	}

	recipes, err := h.service.ListProductionRecipes(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, recipes)
}

// UpdateProductionRecipe updates an existing production recipe
// PUT /api/mes/recipes/:id
func (h *MESHandler) UpdateProductionRecipe(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid recipe ID"})
		return
	}

	var cmd application.UpdateProductionRecipeCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cmd.ID = id

	recipe, err := h.service.UpdateProductionRecipe(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, recipe)
}

// ============================================================================
// PRODUCTION ORDER HANDLERS
// ============================================================================

// CreateProductionOrder creates a new production order
// POST /api/mes/orders
func (h *MESHandler) CreateProductionOrder(c *gin.Context) {
	var cmd application.CreateProductionOrderCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.service.CreateProductionOrder(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

// GetProductionOrder retrieves a single production order
// GET /api/mes/orders/:id
func (h *MESHandler) GetProductionOrder(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	query := application.GetProductionOrderByIDQuery{ID: id}
	order, err := h.service.GetProductionOrderByID(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// ListProductionOrders lists all production orders with optional filters
// GET /api/mes/orders
func (h *MESHandler) ListProductionOrders(c *gin.Context) {
	var query application.ListProductionOrdersQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set defaults
	if query.PageNumber == 0 {
		query.PageNumber = 1
	}
	if query.PageSize == 0 {
		query.PageSize = 20
	}

	orders, err := h.service.ListProductionOrders(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// UpdateProductionOrderStatus updates the status of a production order
// PATCH /api/mes/orders/:id/status
func (h *MESHandler) UpdateProductionOrderStatus(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	var cmd application.UpdateProductionOrderStatusCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cmd.ID = id

	order, err := h.service.UpdateProductionOrderStatus(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// AssignWorkCenter assigns a work center to a production order
// POST /api/mes/orders/:id/assign-workcenter
func (h *MESHandler) AssignWorkCenter(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	var cmd application.AssignWorkCenterCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cmd.ProductionOrderID = id

	order, err := h.service.AssignWorkCenter(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// ============================================================================
// TASK INSTANCE HANDLERS
// ============================================================================

// UpdateTaskStatus updates a task instance status
// PATCH /api/mes/orders/:id/tasks/:taskId/status
func (h *MESHandler) UpdateTaskStatus(c *gin.Context) {
	orderIdParam := c.Param("id")
	orderId, err := uuid.Parse(orderIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	taskIdParam := c.Param("taskId")
	taskId, err := uuid.Parse(taskIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	var cmd application.UpdateTaskStatusCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cmd.ProductionOrderID = orderId
	cmd.TaskInstanceID = taskId

	order, err := h.service.UpdateTaskStatus(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// AssignOperatorToTask assigns an operator to a task instance
// POST /api/mes/orders/:id/tasks/:taskId/assign-operator
func (h *MESHandler) AssignOperatorToTask(c *gin.Context) {
	orderIdParam := c.Param("id")
	orderId, err := uuid.Parse(orderIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	taskIdParam := c.Param("taskId")
	taskId, err := uuid.Parse(taskIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	var cmd application.AssignOperatorToTaskCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cmd.ProductionOrderID = orderId
	cmd.TaskInstanceID = taskId

	order, err := h.service.AssignOperatorToTask(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// RecordTaskProgress records actual time/completion for a task
// POST /api/mes/orders/:id/tasks/:taskId/progress
func (h *MESHandler) RecordTaskProgress(c *gin.Context) {
	orderIdParam := c.Param("id")
	orderId, err := uuid.Parse(orderIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	taskIdParam := c.Param("taskId")
	taskId, err := uuid.Parse(taskIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	var cmd application.RecordTaskProgressCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cmd.ProductionOrderID = orderId
	cmd.TaskInstanceID = taskId

	order, err := h.service.RecordTaskProgress(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// ============================================================================
// WORK CENTER HANDLERS
// ============================================================================

// CreateWorkCenter creates a new work center
// POST /api/mes/workcenters
func (h *MESHandler) CreateWorkCenter(c *gin.Context) {
	var cmd application.CreateWorkCenterCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	workCenter, err := h.service.CreateWorkCenter(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, workCenter)
}

// GetWorkCenter retrieves a single work center
// GET /api/mes/workcenters/:id
func (h *MESHandler) GetWorkCenter(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid work center ID"})
		return
	}

	query := application.GetWorkCenterByIDQuery{ID: id}
	workCenter, err := h.service.GetWorkCenterByID(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, workCenter)
}

// ListWorkCenters lists all work centers
// GET /api/mes/workcenters
func (h *MESHandler) ListWorkCenters(c *gin.Context) {
	var query application.ListWorkCentersQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set defaults
	if query.PageNumber == 0 {
		query.PageNumber = 1
	}
	if query.PageSize == 0 {
		query.PageSize = 20
	}

	workCenters, err := h.service.ListWorkCenters(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, workCenters)
}

// UpdateWorkCenter updates an existing work center
// PUT /api/mes/workcenters/:id
func (h *MESHandler) UpdateWorkCenter(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid work center ID"})
		return
	}

	var cmd application.UpdateWorkCenterCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cmd.ID = id

	workCenter, err := h.service.UpdateWorkCenter(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, workCenter)
}
