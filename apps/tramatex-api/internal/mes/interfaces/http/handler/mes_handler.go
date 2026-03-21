package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joran-cortez/tramatex/internal/mes/application"
)

// MESHandler handles HTTP requests for MES master data.
type MESHandler struct {
	service *application.MESService
}

func NewMESHandler(service *application.MESService) *MESHandler {
	return &MESHandler{service: service}
}

func actorIDFromContext(c *gin.Context) (string, bool) {
	value, ok := c.Get("userID")
	if !ok {
		return "", false
	}
	actorID, ok := value.(string)
	if !ok || actorID == "" {
		return "", false
	}
	return actorID, true
}

func mapServiceError(c *gin.Context, err error) {
	if err == nil {
		return
	}
	message := err.Error()
	if strings.Contains(message, "not found") {
		c.JSON(http.StatusNotFound, gin.H{"error": message})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": message})
}

func parseID(c *gin.Context, name string) (uuid.UUID, bool) {
	id, err := uuid.Parse(c.Param(name))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return uuid.Nil, false
	}
	return id, true
}

// Tasks
func (h *MESHandler) CreateTask(c *gin.Context) {
	actorID, ok := actorIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID is required"})
		return
	}

	var cmd application.CreateTaskCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cmd.ActorID = actorID

	result, err := h.service.CreateTask(c.Request.Context(), cmd)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *MESHandler) GetTask(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	result, err := h.service.GetTaskByID(c.Request.Context(), application.GetTaskByIDQuery{ID: id})
	if err != nil {
		mapServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *MESHandler) ListTasks(c *gin.Context) {
	var query application.ListTasksQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	results, err := h.service.ListTasks(c.Request.Context(), query)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, results)
}

func (h *MESHandler) UpdateTask(c *gin.Context) {
	actorID, ok := actorIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID is required"})
		return
	}

	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	var cmd application.UpdateTaskCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cmd.ID = id
	cmd.ActorID = actorID

	result, err := h.service.UpdateTask(c.Request.Context(), cmd)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *MESHandler) DeleteTask(c *gin.Context) {
	actorID, ok := actorIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID is required"})
		return
	}

	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	if err := h.service.DeleteTask(c.Request.Context(), application.DeleteTaskCommand{ActorID: actorID, ID: id}); err != nil {
		mapServiceError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// Positions
func (h *MESHandler) CreatePosition(c *gin.Context) {
	actorID, ok := actorIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID is required"})
		return
	}

	var cmd application.CreatePositionCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cmd.ActorID = actorID

	result, err := h.service.CreatePosition(c.Request.Context(), cmd)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *MESHandler) GetPosition(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	result, err := h.service.GetPositionByID(c.Request.Context(), application.GetPositionByIDQuery{ID: id})
	if err != nil {
		mapServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *MESHandler) ListPositions(c *gin.Context) {
	var query application.ListPositionsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	results, err := h.service.ListPositions(c.Request.Context(), query)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, results)
}

func (h *MESHandler) UpdatePosition(c *gin.Context) {
	actorID, ok := actorIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID is required"})
		return
	}

	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	var cmd application.UpdatePositionCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cmd.ID = id
	cmd.ActorID = actorID

	result, err := h.service.UpdatePosition(c.Request.Context(), cmd)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *MESHandler) DeletePosition(c *gin.Context) {
	actorID, ok := actorIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID is required"})
		return
	}

	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	if err := h.service.DeletePosition(c.Request.Context(), application.DeletePositionCommand{ActorID: actorID, ID: id}); err != nil {
		mapServiceError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// Work Types
// Work Orders
func (h *MESHandler) CreateWorkType(c *gin.Context) {
	actorID, ok := actorIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID is required"})
		return
	}

	var cmd application.CreateWorkTypeCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cmd.ActorID = actorID

	result, err := h.service.CreateWorkType(c.Request.Context(), cmd)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *MESHandler) CreateWorkOrder(c *gin.Context) {
	actorID, ok := actorIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID is required"})
		return
	}

	var cmd application.CreateWorkOrderCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cmd.ActorID = actorID

	result, err := h.service.CreateWorkOrder(c.Request.Context(), cmd)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *MESHandler) UpdateWorkOrder(c *gin.Context) {
	actorID, ok := actorIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID is required"})
		return
	}

	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	var cmd application.UpdateWorkOrderCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd.ID = id
	cmd.ActorID = actorID

	result, err := h.service.UpdateWorkOrder(c.Request.Context(), cmd)
	if err != nil {
		mapServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *MESHandler) GetWorkOrder(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	result, err := h.service.GetWorkOrderByID(c.Request.Context(), application.GetWorkOrderByIDQuery{ID: id})
	if err != nil {
		mapServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *MESHandler) ListWorkOrders(c *gin.Context) {
	var query application.ListWorkOrdersQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	results, err := h.service.ListWorkOrders(c.Request.Context(), query)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, results)
}

func (h *MESHandler) GetWorkOrderDashboardStats(c *gin.Context) {
	result, err := h.service.GetWorkOrderDashboardStats(c.Request.Context())
	if err != nil {
		mapServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *MESHandler) ListOverdueWorkOrders(c *gin.Context) {
	var query application.ListOverdueWorkOrdersQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	results, err := h.service.ListOverdueWorkOrders(c.Request.Context(), query)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, results)
}

func (h *MESHandler) UpdateWorkOrderTaskStatus(c *gin.Context) {
	actorID, ok := actorIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID is required"})
		return
	}

	workID, ok := parseID(c, "workId")
	if !ok {
		return
	}

	taskID, ok := parseID(c, "taskId")
	if !ok {
		return
	}

	var cmd application.UpdateWorkOrderTaskStatusCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd.ActorID = actorID
	cmd.WorkID = workID
	cmd.TaskID = taskID

	result, err := h.service.UpdateWorkOrderTaskStatus(c.Request.Context(), cmd)
	if err != nil {
		mapServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *MESHandler) GetWorkType(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	result, err := h.service.GetWorkTypeByID(c.Request.Context(), application.GetWorkTypeByIDQuery{ID: id})
	if err != nil {
		mapServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *MESHandler) ListWorkTypes(c *gin.Context) {
	var query application.ListWorkTypesQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	results, err := h.service.ListWorkTypes(c.Request.Context(), query)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, results)
}

func (h *MESHandler) UpdateWorkType(c *gin.Context) {
	actorID, ok := actorIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID is required"})
		return
	}

	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	var cmd application.UpdateWorkTypeCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cmd.ID = id
	cmd.ActorID = actorID

	result, err := h.service.UpdateWorkType(c.Request.Context(), cmd)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *MESHandler) DeleteWorkType(c *gin.Context) {
	actorID, ok := actorIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID is required"})
		return
	}

	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	if err := h.service.DeleteWorkType(c.Request.Context(), application.DeleteWorkTypeCommand{ActorID: actorID, ID: id}); err != nil {
		mapServiceError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// --- WorkSetup ---

func (h *MESHandler) CreateWorkSetup(c *gin.Context) {
	actorID, ok := actorIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID is required"})
		return
	}

	var cmd application.CreateWorkSetupCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cmd.ActorID = actorID

	result, err := h.service.CreateWorkSetup(c.Request.Context(), cmd)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *MESHandler) GetWorkSetup(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	result, err := h.service.GetWorkSetupByID(c.Request.Context(), application.GetWorkSetupByIDQuery{ID: id})
	if err != nil {
		mapServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *MESHandler) ListWorkSetups(c *gin.Context) {
	var query application.ListWorkSetupsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	results, err := h.service.ListWorkSetups(c.Request.Context(), query)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, results)
}

func (h *MESHandler) UpdateWorkSetup(c *gin.Context) {
	actorID, ok := actorIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID is required"})
		return
	}

	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	var cmd application.UpdateWorkSetupCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cmd.ID = id
	cmd.ActorID = actorID

	result, err := h.service.UpdateWorkSetup(c.Request.Context(), cmd)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *MESHandler) DeleteWorkSetup(c *gin.Context) {
	actorID, ok := actorIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID is required"})
		return
	}

	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	if err := h.service.DeleteWorkSetup(c.Request.Context(), application.DeleteWorkSetupCommand{ActorID: actorID, ID: id}); err != nil {
		mapServiceError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// ListPendingWorkSetups returns pending work setups from confirmed Sales orders
// that have no MES WorkOrder yet. Delegates to a Sales adapter.
func (h *MESHandler) ListPendingWorkSetups(c *gin.Context) {
	result, err := h.service.ListPendingWorkSetups(c.Request.Context())
	if err != nil {
		mapServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}
