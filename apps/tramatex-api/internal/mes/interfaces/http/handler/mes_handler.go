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

// Service groups
// MES Works
func (h *MESHandler) CreateServiceGroup(c *gin.Context) {
	actorID, ok := actorIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID is required"})
		return
	}

	var cmd application.CreateServiceGroupCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cmd.ActorID = actorID

	result, err := h.service.CreateServiceGroup(c.Request.Context(), cmd)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *MESHandler) CreateServiceTemplate(c *gin.Context) {
	h.CreateServiceGroup(c)
}

func (h *MESHandler) CreateMESWork(c *gin.Context) {
	actorID, ok := actorIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID is required"})
		return
	}

	var cmd application.CreateMESWorkCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cmd.ActorID = actorID

	result, err := h.service.CreateMESWork(c.Request.Context(), cmd)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *MESHandler) CreateWorkDefinition(c *gin.Context) {
	h.CreateMESWork(c)
}

func (h *MESHandler) UpdateMESWork(c *gin.Context) {
	actorID, ok := actorIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID is required"})
		return
	}

	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	var cmd application.UpdateMESWorkCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd.ID = id
	cmd.ActorID = actorID

	result, err := h.service.UpdateMESWork(c.Request.Context(), cmd)
	if err != nil {
		mapServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *MESHandler) UpdateWorkDefinition(c *gin.Context) {
	h.UpdateMESWork(c)
}

func (h *MESHandler) GetMESWork(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	result, err := h.service.GetMESWorkByID(c.Request.Context(), application.GetMESWorkByIDQuery{ID: id})
	if err != nil {
		mapServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *MESHandler) GetWorkDefinition(c *gin.Context) {
	h.GetMESWork(c)
}

func (h *MESHandler) ListMESWorks(c *gin.Context) {
	var query application.ListMESWorksQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	results, err := h.service.ListMESWorks(c.Request.Context(), query)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, results)
}

func (h *MESHandler) ListWorkDefinitions(c *gin.Context) {
	h.ListMESWorks(c)
}

func (h *MESHandler) GetMESWorkDashboardStats(c *gin.Context) {
	result, err := h.service.GetMESWorkDashboardStats(c.Request.Context())
	if err != nil {
		mapServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *MESHandler) GetWorkDefinitionDashboardStats(c *gin.Context) {
	h.GetMESWorkDashboardStats(c)
}

func (h *MESHandler) ListOverdueMESWorks(c *gin.Context) {
	var query application.ListOverdueMESWorksQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	results, err := h.service.ListOverdueMESWorks(c.Request.Context(), query)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, results)
}

func (h *MESHandler) ListOverdueWorkDefinitions(c *gin.Context) {
	h.ListOverdueMESWorks(c)
}

func (h *MESHandler) UpdateMESWorkTaskStatus(c *gin.Context) {
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

	var cmd application.UpdateMESWorkTaskStatusCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd.ActorID = actorID
	cmd.WorkID = workID
	cmd.TaskID = taskID

	result, err := h.service.UpdateMESWorkTaskStatus(c.Request.Context(), cmd)
	if err != nil {
		mapServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *MESHandler) UpdateWorkDefinitionTaskStatus(c *gin.Context) {
	h.UpdateMESWorkTaskStatus(c)
}

func (h *MESHandler) GetServiceGroup(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	result, err := h.service.GetServiceGroupByID(c.Request.Context(), application.GetServiceGroupByIDQuery{ID: id})
	if err != nil {
		mapServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *MESHandler) GetServiceTemplate(c *gin.Context) {
	h.GetServiceGroup(c)
}

func (h *MESHandler) ListServiceGroups(c *gin.Context) {
	var query application.ListServiceGroupsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	results, err := h.service.ListServiceGroups(c.Request.Context(), query)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, results)
}

func (h *MESHandler) ListServiceTemplates(c *gin.Context) {
	h.ListServiceGroups(c)
}

func (h *MESHandler) UpdateServiceGroup(c *gin.Context) {
	actorID, ok := actorIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID is required"})
		return
	}

	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	var cmd application.UpdateServiceGroupCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cmd.ID = id
	cmd.ActorID = actorID

	result, err := h.service.UpdateServiceGroup(c.Request.Context(), cmd)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *MESHandler) UpdateServiceTemplate(c *gin.Context) {
	h.UpdateServiceGroup(c)
}

func (h *MESHandler) DeleteServiceGroup(c *gin.Context) {
	actorID, ok := actorIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID is required"})
		return
	}

	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	if err := h.service.DeleteServiceGroup(c.Request.Context(), application.DeleteServiceGroupCommand{ActorID: actorID, ID: id}); err != nil {
		mapServiceError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *MESHandler) DeleteServiceTemplate(c *gin.Context) {
	h.DeleteServiceGroup(c)
}
