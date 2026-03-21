package application

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joran-cortez/tramatex/internal/mes/domain"
)

// SalesOrderLinker updates a Sales order_work_setup record to link it
// with a newly created MES WorkOrder. MES calls this after creating a
// work order so Sales can track production. Implementation lives in
// the Sales infrastructure layer.
type SalesOrderLinker interface {
	LinkWorkOrder(ctx context.Context, orderWorkSetupID uuid.UUID, workOrderID uuid.UUID) error
}

// PendingWorkSetupProvider retrieves confirmed-order WorkSetup configs
// that don't yet have a MES WorkOrder. Implementation lives in the
// Sales infrastructure layer so MES never queries Sales tables directly.
type PendingWorkSetupProvider interface {
	ListPending(ctx context.Context) ([]PendingWorkSetupDTO, error)
}

// WorkOrderSalesInfo holds the Sales order reference for a given WorkOrder.
type WorkOrderSalesInfo struct {
	SalesOrderID     uuid.UUID
	SalesOrderNumber string
}

// WorkOrderSalesInfoProvider returns Sales order info for a batch of
// WorkOrder IDs. Implementation lives in Sales infrastructure.
type WorkOrderSalesInfoProvider interface {
	GetSalesInfoByWorkOrderIDs(ctx context.Context, workOrderIDs []uuid.UUID) (map[uuid.UUID]WorkOrderSalesInfo, error)
}

// MESService provides operations for Manufacturing Execution System
// focused on foundation master data CRUD.
type MESService struct {
	taskRepo             domain.TaskRepository
	positionRepo         domain.PositionRepository
	workTypeRepo         domain.WorkTypeRepository
	workOrderRepo        domain.WorkOrderRepository
	workSetupRepo        domain.WorkSetupRepository
	salesOrderLinker     SalesOrderLinker
	pendingSetupProvider PendingWorkSetupProvider
	salesInfoProvider    WorkOrderSalesInfoProvider
}

// NewMESService creates a new MES service.
func NewMESService(
	taskRepo domain.TaskRepository,
	positionRepo domain.PositionRepository,
	workTypeRepo domain.WorkTypeRepository,
	workOrderRepo domain.WorkOrderRepository,
	workSetupRepo domain.WorkSetupRepository,
) *MESService {
	return &MESService{
		taskRepo:      taskRepo,
		positionRepo:  positionRepo,
		workTypeRepo:  workTypeRepo,
		workOrderRepo: workOrderRepo,
		workSetupRepo: workSetupRepo,
	}
}

// SetSalesOrderLinker configures the cross-module linker (optional, nil-safe).
func (s *MESService) SetSalesOrderLinker(linker SalesOrderLinker) {
	s.salesOrderLinker = linker
}

// SetPendingSetupProvider configures the provider for pending setups from Sales.
func (s *MESService) SetPendingSetupProvider(provider PendingWorkSetupProvider) {
	s.pendingSetupProvider = provider
}

// SetSalesInfoProvider configures the provider for Sales order info on WorkOrders.
func (s *MESService) SetSalesInfoProvider(provider WorkOrderSalesInfoProvider) {
	s.salesInfoProvider = provider
}

func (s *MESService) CreateTask(ctx context.Context, cmd CreateTaskCommand) (*TaskDTO, error) {
	description := ""
	if cmd.Description != nil {
		description = *cmd.Description
	}
	reference := ""
	if cmd.Reference != nil {
		reference = *cmd.Reference
	}
	isActive := true
	if cmd.IsActive != nil {
		isActive = *cmd.IsActive
	}

	task, err := domain.NewTask(cmd.Name, reference, description, isActive)
	if err != nil {
		return nil, err
	}

	if err := s.taskRepo.Save(ctx, task); err != nil {
		return nil, fmt.Errorf("save task: %w", err)
	}

	return toTaskDTO(task), nil
}

func (s *MESService) GetTaskByID(ctx context.Context, query GetTaskByIDQuery) (*TaskDTO, error) {
	task, err := s.taskRepo.FindByID(ctx, query.ID)
	if err != nil {
		return nil, fmt.Errorf("find task by id: %w", err)
	}
	if task == nil {
		return nil, fmt.Errorf("task not found")
	}

	return toTaskDTO(task), nil
}

func (s *MESService) ListTasks(ctx context.Context, query ListTasksQuery) ([]TaskDTO, error) {
	tasks, err := s.taskRepo.FindAll(ctx, &domain.TaskFilters{
		IsActive: query.IsActive,
		Search:   query.Search,
	})
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}

	dtos := make([]TaskDTO, 0, len(tasks))
	for _, task := range tasks {
		dtos = append(dtos, *toTaskDTO(task))
	}
	return dtos, nil
}

func (s *MESService) UpdateTask(ctx context.Context, cmd UpdateTaskCommand) (*TaskDTO, error) {
	task, err := s.taskRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("find task for update: %w", err)
	}
	if task == nil {
		return nil, fmt.Errorf("task not found")
	}

	if cmd.Name != nil {
		task.Name = *cmd.Name
	}
	if cmd.Reference != nil {
		task.Reference = *cmd.Reference
	}
	if cmd.Description != nil {
		task.Description = *cmd.Description
	}
	if cmd.IsActive != nil {
		task.IsActive = *cmd.IsActive
	}

	if err := s.taskRepo.Save(ctx, task); err != nil {
		return nil, fmt.Errorf("save task update: %w", err)
	}

	return toTaskDTO(task), nil
}

func (s *MESService) DeleteTask(ctx context.Context, cmd DeleteTaskCommand) error {
	if err := s.taskRepo.Delete(ctx, cmd.ID); err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	return nil
}

func (s *MESService) CreatePosition(ctx context.Context, cmd CreatePositionCommand) (*PositionDTO, error) {
	description := ""
	if cmd.Description != nil {
		description = *cmd.Description
	}
	isActive := true
	if cmd.IsActive != nil {
		isActive = *cmd.IsActive
	}

	position, err := domain.NewPosition(cmd.Name, cmd.Code, description, isActive)
	if err != nil {
		return nil, err
	}

	if err := s.positionRepo.Save(ctx, position); err != nil {
		return nil, fmt.Errorf("save position: %w", err)
	}

	return toPositionDTO(position), nil
}

func (s *MESService) GetPositionByID(ctx context.Context, query GetPositionByIDQuery) (*PositionDTO, error) {
	position, err := s.positionRepo.FindByID(ctx, query.ID)
	if err != nil {
		return nil, fmt.Errorf("find position by id: %w", err)
	}
	if position == nil {
		return nil, fmt.Errorf("position not found")
	}

	return toPositionDTO(position), nil
}

func (s *MESService) ListPositions(ctx context.Context, query ListPositionsQuery) ([]PositionDTO, error) {
	positions, err := s.positionRepo.FindAll(ctx, &domain.PositionFilters{
		IsActive: query.IsActive,
		Search:   query.Search,
	})
	if err != nil {
		return nil, fmt.Errorf("list positions: %w", err)
	}

	dtos := make([]PositionDTO, 0, len(positions))
	for _, position := range positions {
		dtos = append(dtos, *toPositionDTO(position))
	}
	return dtos, nil
}

func (s *MESService) UpdatePosition(ctx context.Context, cmd UpdatePositionCommand) (*PositionDTO, error) {
	position, err := s.positionRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("find position for update: %w", err)
	}
	if position == nil {
		return nil, fmt.Errorf("position not found")
	}

	if cmd.Name != nil {
		position.Name = *cmd.Name
	}
	if cmd.Code != nil {
		position.Code = *cmd.Code
	}
	if cmd.Description != nil {
		position.Description = *cmd.Description
	}
	if cmd.IsActive != nil {
		position.IsActive = *cmd.IsActive
	}

	if err := s.positionRepo.Save(ctx, position); err != nil {
		return nil, fmt.Errorf("save position update: %w", err)
	}

	return toPositionDTO(position), nil
}

func (s *MESService) DeletePosition(ctx context.Context, cmd DeletePositionCommand) error {
	if err := s.positionRepo.Delete(ctx, cmd.ID); err != nil {
		return fmt.Errorf("delete position: %w", err)
	}
	return nil
}

func (s *MESService) CreateWorkType(ctx context.Context, cmd CreateWorkTypeCommand) (*WorkTypeDTO, error) {
	description := ""
	if cmd.Description != nil {
		description = *cmd.Description
	}
	reference := ""
	if cmd.Reference != nil {
		reference = *cmd.Reference
	}
	isActive := true
	if cmd.IsActive != nil {
		isActive = *cmd.IsActive
	}
	assignments := mapTaskAssignments(cmd.TaskAssignments)

	workType, err := domain.NewWorkType(cmd.Name, reference, description, isActive, assignments)
	if err != nil {
		return nil, err
	}

	if err := s.workTypeRepo.Save(ctx, workType); err != nil {
		return nil, fmt.Errorf("save work type: %w", err)
	}

	return toWorkTypeDTO(workType), nil
}

func (s *MESService) GetWorkTypeByID(ctx context.Context, query GetWorkTypeByIDQuery) (*WorkTypeDTO, error) {
	workType, err := s.workTypeRepo.FindByID(ctx, query.ID)
	if err != nil {
		return nil, fmt.Errorf("find work type by id: %w", err)
	}
	if workType == nil {
		return nil, fmt.Errorf("work type not found")
	}

	return toWorkTypeDTO(workType), nil
}

func (s *MESService) ListWorkTypes(ctx context.Context, query ListWorkTypesQuery) ([]WorkTypeDTO, error) {
	workTypes, err := s.workTypeRepo.FindAll(ctx, &domain.WorkTypeFilters{
		IsActive: query.IsActive,
		Search:   query.Search,
	})
	if err != nil {
		return nil, fmt.Errorf("list work types: %w", err)
	}

	dtos := make([]WorkTypeDTO, 0, len(workTypes))
	for _, wt := range workTypes {
		dtos = append(dtos, *toWorkTypeDTO(wt))
	}
	return dtos, nil
}

func (s *MESService) UpdateWorkType(ctx context.Context, cmd UpdateWorkTypeCommand) (*WorkTypeDTO, error) {
	workType, err := s.workTypeRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("find work type for update: %w", err)
	}
	if workType == nil {
		return nil, fmt.Errorf("work type not found")
	}

	if cmd.Name != nil {
		workType.Name = *cmd.Name
	}
	if cmd.Reference != nil {
		workType.Reference = *cmd.Reference
	}
	if cmd.Description != nil {
		workType.Description = *cmd.Description
	}
	if cmd.IsActive != nil {
		workType.IsActive = *cmd.IsActive
	}
	if cmd.TaskAssignments != nil {
		workType.Tasks = mapTaskAssignments(cmd.TaskAssignments)
	}

	if err := s.workTypeRepo.Save(ctx, workType); err != nil {
		return nil, fmt.Errorf("save work type update: %w", err)
	}

	return toWorkTypeDTO(workType), nil
}

func (s *MESService) DeleteWorkType(ctx context.Context, cmd DeleteWorkTypeCommand) error {
	if err := s.workTypeRepo.Delete(ctx, cmd.ID); err != nil {
		return fmt.Errorf("delete work type: %w", err)
	}
	return nil
}

func mapTaskAssignments(inputs []WorkTypeTaskInput) []domain.WorkTypeTask {
	if len(inputs) == 0 {
		return []domain.WorkTypeTask{}
	}

	assignments := make([]domain.WorkTypeTask, 0, len(inputs))
	for _, input := range inputs {
		assignments = append(assignments, domain.WorkTypeTask{
			TaskID:   input.TaskID,
			Sequence: input.Sequence,
		})
	}
	return assignments
}

func toTaskDTO(task *domain.Task) *TaskDTO {
	if task == nil {
		return nil
	}

	return &TaskDTO{
		ID:          task.ID,
		Name:        task.Name,
		Reference:   task.Reference,
		Description: task.Description,
		IsActive:    task.IsActive,
	}
}

func toPositionDTO(position *domain.Position) *PositionDTO {
	if position == nil {
		return nil
	}

	return &PositionDTO{
		ID:          position.ID,
		Name:        position.Name,
		Code:        position.Code,
		Description: position.Description,
		IsActive:    position.IsActive,
	}
}

func toWorkTypeDTO(workType *domain.WorkType) *WorkTypeDTO {
	if workType == nil {
		return nil
	}

	tasks := make([]WorkTypeTaskDTO, 0, len(workType.Tasks))
	for _, task := range workType.Tasks {
		tasks = append(tasks, WorkTypeTaskDTO{
			TaskID:   task.TaskID,
			Sequence: task.Sequence,
		})
	}

	return &WorkTypeDTO{
		ID:          workType.ID,
		Name:        workType.Name,
		Reference:   workType.Reference,
		Description: workType.Description,
		IsActive:    workType.IsActive,
		Tasks:       tasks,
	}
}

func (s *MESService) CreateWorkOrder(ctx context.Context, cmd CreateWorkOrderCommand) (*WorkOrderDTO, error) {
	priority := domain.WorkPriorityNormal
	if cmd.Priority != nil && *cmd.Priority != "" {
		priority = domain.WorkPriority(strings.ToUpper(*cmd.Priority))
	}

	var dueDate *time.Time
	if cmd.DueDate != nil {
		parsed, parseErr := parseOptionalDate(*cmd.DueDate)
		if parseErr != nil {
			return nil, parseErr
		}
		dueDate = parsed
	}

	year := time.Now().UTC().Year()
	count, err := s.workOrderRepo.CountByYear(ctx, year)
	if err != nil {
		return nil, fmt.Errorf("count works by year: %w", err)
	}
	workNumber := fmt.Sprintf("MES-%d-%03d", year, count+1)

	notes := ""
	if cmd.Notes != nil {
		notes = *cmd.Notes
	}

	// Build lines from WorkSetup if provided, otherwise create without lines.
	var lines []domain.WorkOrderLine
	if cmd.WorkSetupID != nil && *cmd.WorkSetupID != uuid.Nil {
		workSetup, wsErr := s.workSetupRepo.FindByID(ctx, *cmd.WorkSetupID)
		if wsErr != nil {
			return nil, fmt.Errorf("find work setup: %w", wsErr)
		}
		if workSetup == nil {
			return nil, fmt.Errorf("work setup not found")
		}
		if len(workSetup.Lines) == 0 {
			return nil, fmt.Errorf("work setup has no lines configured")
		}

		lines = make([]domain.WorkOrderLine, 0, len(workSetup.Lines))
		for _, setupLine := range workSetup.Lines {
			workType, wtErr := s.workTypeRepo.FindByID(ctx, setupLine.WorkTypeID)
			if wtErr != nil {
				return nil, fmt.Errorf("find work type for line: %w", wtErr)
			}
			if workType == nil {
				return nil, fmt.Errorf("work type %s not found", setupLine.WorkTypeID)
			}

			tasks := make([]domain.WorkOrderTask, 0, len(workType.Tasks))
			for _, wtt := range workType.Tasks {
				tasks = append(tasks, domain.WorkOrderTask{
					ID:       uuid.New(),
					TaskID:   wtt.TaskID,
					Sequence: wtt.Sequence,
					Status:   domain.TaskStatusPending,
				})
			}

			lines = append(lines, domain.WorkOrderLine{
				ID:             uuid.New(),
				WorkTypeID:     setupLine.WorkTypeID,
				PositionID:     setupLine.PositionID,
				DesignFilePath: setupLine.DesignFilePath,
				Notes:          setupLine.Notes,
				Sequence:       setupLine.Sequence,
				Tasks:          tasks,
			})
		}
	}

	work, err := domain.NewWorkOrder(
		workNumber,
		cmd.WorkName,
		cmd.PartyID,
		cmd.WorkSetupID,
		notes,
		priority,
		dueDate,
		lines,
	)
	if err != nil {
		return nil, err
	}

	if err := s.workOrderRepo.Save(ctx, work); err != nil {
		return nil, fmt.Errorf("save work order: %w", err)
	}

	// Link back to Sales order_work_setup if provided
	if cmd.OrderWorkSetupID != nil && *cmd.OrderWorkSetupID != uuid.Nil && s.salesOrderLinker != nil {
		if linkErr := s.salesOrderLinker.LinkWorkOrder(ctx, *cmd.OrderWorkSetupID, work.ID); linkErr != nil {
			return nil, fmt.Errorf("link work order to sales: %w", linkErr)
		}
	}

	return toWorkOrderDTO(work), nil
}

func (s *MESService) UpdateWorkOrder(ctx context.Context, cmd UpdateWorkOrderCommand) (*WorkOrderDTO, error) {
	work, err := s.workOrderRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("find work order for update: %w", err)
	}
	if work == nil {
		return nil, fmt.Errorf("work order not found")
	}

	if cmd.WorkName != nil {
		trimmed := strings.TrimSpace(*cmd.WorkName)
		if trimmed == "" {
			return nil, fmt.Errorf("work name is required")
		}
		work.OrderName = trimmed
	}

	if cmd.Notes != nil {
		work.Notes = strings.TrimSpace(*cmd.Notes)
	}

	if cmd.Status != nil {
		status := domain.ProductionStatus(strings.ToUpper(strings.TrimSpace(*cmd.Status)))
		if !status.IsValid() {
			return nil, fmt.Errorf("invalid production status")
		}
		if status == domain.ProductionStatusCancelled && work.Status != domain.ProductionStatusSuspended {
			return nil, fmt.Errorf("only suspended work orders can be cancelled")
		}
		if status == domain.ProductionStatusPending && work.Status != domain.ProductionStatusSuspended && work.Status != domain.ProductionStatusCancelled {
			return nil, fmt.Errorf("only suspended or cancelled work orders can be reactivated to pending")
		}
		work.Status = status
	}

	if cmd.Priority != nil {
		priority := domain.WorkPriority(strings.ToUpper(strings.TrimSpace(*cmd.Priority)))
		if !priority.IsValid() {
			return nil, fmt.Errorf("invalid work priority")
		}
		work.Priority = priority
	}

	if cmd.DueDate != nil {
		parsed, err := parseOptionalDate(*cmd.DueDate)
		if err != nil {
			return nil, err
		}
		work.DueDate = parsed
	}

	// Assign WorkSetup and generate lines from it
	if cmd.WorkSetupID != nil && *cmd.WorkSetupID != uuid.Nil {
		workSetup, wsErr := s.workSetupRepo.FindByID(ctx, *cmd.WorkSetupID)
		if wsErr != nil {
			return nil, fmt.Errorf("find work setup: %w", wsErr)
		}
		if workSetup == nil {
			return nil, fmt.Errorf("work setup not found")
		}
		if len(workSetup.Lines) == 0 {
			return nil, fmt.Errorf("work setup has no lines configured")
		}

		lines := make([]domain.WorkOrderLine, 0, len(workSetup.Lines))
		for _, setupLine := range workSetup.Lines {
			workType, wtErr := s.workTypeRepo.FindByID(ctx, setupLine.WorkTypeID)
			if wtErr != nil {
				return nil, fmt.Errorf("find work type for line: %w", wtErr)
			}
			if workType == nil {
				return nil, fmt.Errorf("work type %s not found", setupLine.WorkTypeID)
			}

			tasks := make([]domain.WorkOrderTask, 0, len(workType.Tasks))
			for _, wtt := range workType.Tasks {
				tasks = append(tasks, domain.WorkOrderTask{
					ID:       uuid.New(),
					TaskID:   wtt.TaskID,
					Sequence: wtt.Sequence,
					Status:   domain.TaskStatusPending,
				})
			}

			lines = append(lines, domain.WorkOrderLine{
				ID:             uuid.New(),
				WorkTypeID:     setupLine.WorkTypeID,
				PositionID:     setupLine.PositionID,
				DesignFilePath: setupLine.DesignFilePath,
				Notes:          setupLine.Notes,
				Sequence:       setupLine.Sequence,
				Tasks:          tasks,
			})
		}

		work.WorkSetupID = cmd.WorkSetupID
		work.Lines = lines
	}

	if err := s.workOrderRepo.Save(ctx, work); err != nil {
		return nil, fmt.Errorf("save work order update: %w", err)
	}

	return toWorkOrderDTO(work), nil
}

func (s *MESService) UpdateWorkOrderTaskStatus(ctx context.Context, cmd UpdateWorkOrderTaskStatusCommand) (*WorkOrderDTO, error) {
	work, err := s.workOrderRepo.FindByID(ctx, cmd.WorkID)
	if err != nil {
		return nil, fmt.Errorf("find work order for task update: %w", err)
	}
	if work == nil {
		return nil, fmt.Errorf("work order not found")
	}

	action := strings.ToUpper(strings.TrimSpace(cmd.Action))
	now := time.Now().UTC()
	assignedTo := parseActorUUID(cmd.ActorID)

	found := false
	for groupIndex := range work.Lines {
		for taskIndex := range work.Lines[groupIndex].Tasks {
			task := &work.Lines[groupIndex].Tasks[taskIndex]
			if task.ID != cmd.TaskID {
				continue
			}

			found = true
			switch action {
			case "START":
				task.Status = domain.TaskStatusInProgress
				if task.StartedAt == nil {
					task.StartedAt = &now
				}
				task.CompletedAt = nil
				if assignedTo != nil {
					task.AssignedTo = assignedTo
				}
				if work.StartDate == nil {
					work.StartDate = &now
				}
			case "PAUSE":
				task.Status = domain.TaskStatusPending
			case "COMPLETE":
				task.Status = domain.TaskStatusCompleted
				if task.StartedAt == nil {
					task.StartedAt = &now
				}
				task.CompletedAt = &now
				if assignedTo != nil {
					task.AssignedTo = assignedTo
				}
			case "BLOCK":
				task.Status = domain.TaskStatusBlocked
			default:
				return nil, fmt.Errorf("invalid action: %s", cmd.Action)
			}

			if cmd.Notes != nil {
				task.Notes = strings.TrimSpace(*cmd.Notes)
			}
		}
	}

	if !found {
		return nil, fmt.Errorf("work order task not found")
	}

	recalculateWorkStatus(work, now)

	if err := s.workOrderRepo.Save(ctx, work); err != nil {
		return nil, fmt.Errorf("save work order task update: %w", err)
	}

	return toWorkOrderDTO(work), nil
}

func (s *MESService) ListWorkOrders(ctx context.Context, query ListWorkOrdersQuery) ([]WorkOrderDTO, error) {
	var status *domain.ProductionStatus
	if query.Status != nil && *query.Status != "" {
		s := domain.ProductionStatus(strings.ToUpper(*query.Status))
		status = &s
	}

	filters := &domain.WorkOrderFilters{
		Status:  status,
		Search:  query.Search,
		PartyID: query.PartyID,
	}
	if query.WorkSetupID != "" {
		parsed, err := uuid.Parse(query.WorkSetupID)
		if err == nil {
			filters.WorkSetupID = &parsed
		}
	}

	works, err := s.workOrderRepo.FindAll(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("list work orders: %w", err)
	}

	result := make([]WorkOrderDTO, 0, len(works))
	for _, work := range works {
		result = append(result, *toWorkOrderDTO(work))
	}
	s.enrichWorkOrderDTOs(ctx, result)
	return result, nil
}

func (s *MESService) GetWorkOrderByID(ctx context.Context, query GetWorkOrderByIDQuery) (*WorkOrderDTO, error) {
	work, err := s.workOrderRepo.FindByID(ctx, query.ID)
	if err != nil {
		return nil, fmt.Errorf("find work order by id: %w", err)
	}
	if work == nil {
		return nil, fmt.Errorf("work order not found")
	}
	dto := toWorkOrderDTO(work)
	dtos := []WorkOrderDTO{*dto}
	s.enrichWorkOrderDTOs(ctx, dtos)
	return &dtos[0], nil
}

func (s *MESService) GetWorkOrderDashboardStats(ctx context.Context) (*WorkOrderDashboardStatsDTO, error) {
	works, err := s.workOrderRepo.FindAll(ctx, &domain.WorkOrderFilters{})
	if err != nil {
		return nil, fmt.Errorf("get work order dashboard stats: %w", err)
	}

	now := time.Now().UTC()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	tomorrowStart := todayStart.Add(24 * time.Hour)

	stats := &WorkOrderDashboardStatsDTO{
		Total:    len(works),
		ByStatus: map[string]int{},
	}

	for _, work := range works {
		stats.ByStatus[string(work.Status)]++

		if work.DueDate != nil {
			if work.DueDate.Before(todayStart) && isOpenWorkOrderStatus(work.Status) {
				stats.Overdue++
			}
			if (work.DueDate.Equal(todayStart) || work.DueDate.After(todayStart)) && work.DueDate.Before(tomorrowStart) {
				stats.DueToday++
			}
		}
	}

	return stats, nil
}

func (s *MESService) ListOverdueWorkOrders(ctx context.Context, query ListOverdueWorkOrdersQuery) ([]WorkOrderDTO, error) {
	works, err := s.workOrderRepo.FindAll(ctx, &domain.WorkOrderFilters{})
	if err != nil {
		return nil, fmt.Errorf("list overdue work orders: %w", err)
	}

	todayStart := time.Date(time.Now().UTC().Year(), time.Now().UTC().Month(), time.Now().UTC().Day(), 0, 0, 0, 0, time.UTC)
	overdue := make([]*domain.WorkOrder, 0)
	for _, work := range works {
		if work.DueDate != nil && work.DueDate.Before(todayStart) && isOpenWorkOrderStatus(work.Status) {
			overdue = append(overdue, work)
		}
	}

	sort.SliceStable(overdue, func(i, j int) bool {
		return overdue[i].DueDate.Before(*overdue[j].DueDate)
	})

	if query.Limit > 0 && len(overdue) > query.Limit {
		overdue = overdue[:query.Limit]
	}

	result := make([]WorkOrderDTO, 0, len(overdue))
	for _, work := range overdue {
		result = append(result, *toWorkOrderDTO(work))
	}
	s.enrichWorkOrderDTOs(ctx, result)

	return result, nil
}

func isOpenWorkOrderStatus(status domain.ProductionStatus) bool {
	return status != domain.ProductionStatusCompleted && status != domain.ProductionStatusCancelled
}

// ListPendingWorkSetups returns confirmed-order work setups without a WorkOrder,
// delegating to the PendingSetupProvider adapter (Sales infrastructure).
func (s *MESService) ListPendingWorkSetups(ctx context.Context) ([]PendingWorkSetupDTO, error) {
	if s.pendingSetupProvider == nil {
		return []PendingWorkSetupDTO{}, nil
	}
	return s.pendingSetupProvider.ListPending(ctx)
}

// SuspendWorkOrders puts the given WorkOrders into SUSPENDED state.
// Only affects orders in PENDING, IN_PROGRESS or ON_HOLD — skips
// COMPLETED and CANCELLED silently.
func (s *MESService) SuspendWorkOrders(ctx context.Context, ids []uuid.UUID) error {
	for _, id := range ids {
		work, err := s.workOrderRepo.FindByID(ctx, id)
		if err != nil {
			return fmt.Errorf("find work order %s for suspend: %w", id, err)
		}
		if work == nil {
			continue
		}
		switch work.Status {
		case domain.ProductionStatusPending, domain.ProductionStatusInProgress, domain.ProductionStatusOnHold:
			work.Status = domain.ProductionStatusSuspended
			if err := s.workOrderRepo.Save(ctx, work); err != nil {
				return fmt.Errorf("save suspended work order %s: %w", id, err)
			}
		default:
			// COMPLETED, CANCELLED, already SUSPENDED → skip
		}
	}
	return nil
}

// ReactivateWorkOrders moves WorkOrders back to PENDING so production
// can resume.  Skips COMPLETED (finished) and IN_PROGRESS (MES already
// decided to continue).
func (s *MESService) ReactivateWorkOrders(ctx context.Context, ids []uuid.UUID) error {
	for _, id := range ids {
		work, err := s.workOrderRepo.FindByID(ctx, id)
		if err != nil {
			return fmt.Errorf("find work order %s for reactivate: %w", id, err)
		}
		if work == nil {
			continue
		}
		switch work.Status {
		case domain.ProductionStatusCompleted, domain.ProductionStatusInProgress:
			// don't touch — finished or MES chose to keep working
		case domain.ProductionStatusPending:
			// already pending — no-op
		default:
			// SUSPENDED, ON_HOLD, CANCELLED → PENDING
			work.Status = domain.ProductionStatusPending
			if err := s.workOrderRepo.Save(ctx, work); err != nil {
				return fmt.Errorf("save reactivated work order %s: %w", id, err)
			}
		}
	}
	return nil
}

func parseActorUUID(actorID string) *uuid.UUID {
	if actorID == "" {
		return nil
	}

	parsed, err := uuid.Parse(actorID)
	if err != nil {
		return nil
	}

	return &parsed
}

func recalculateWorkStatus(work *domain.WorkOrder, now time.Time) {
	if work == nil {
		return
	}
	// SUSPENDED is managed externally (e.g. order cancellation); do not overwrite.
	if work.Status == domain.ProductionStatusSuspended {
		return
	}

	total := 0
	completedOrSkipped := 0
	hasInProgress := false
	hasBlocked := false
	hasPending := false

	for _, group := range work.Lines {
		for _, task := range group.Tasks {
			total++
			switch task.Status {
			case domain.TaskStatusCompleted, domain.TaskStatusSkipped:
				completedOrSkipped++
			case domain.TaskStatusInProgress:
				hasInProgress = true
			case domain.TaskStatusBlocked:
				hasBlocked = true
			case domain.TaskStatusPending:
				hasPending = true
			}
		}
	}

	if total > 0 && completedOrSkipped == total {
		work.Status = domain.ProductionStatusCompleted
		work.CompletedDate = &now
		if work.StartDate == nil {
			work.StartDate = &now
		}
		return
	}

	work.CompletedDate = nil

	if hasBlocked {
		work.Status = domain.ProductionStatusOnHold
		return
	}

	if hasInProgress {
		work.Status = domain.ProductionStatusInProgress
		if work.StartDate == nil {
			work.StartDate = &now
		}
		return
	}

	if hasPending {
		work.Status = domain.ProductionStatusPending
		return
	}
}

func parseOptionalDate(raw string) (*time.Time, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil, nil
	}

	if parsed, err := time.Parse(time.RFC3339, trimmed); err == nil {
		value := parsed.UTC()
		return &value, nil
	}

	if parsed, err := time.Parse("2006-01-02", trimmed); err == nil {
		value := time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, time.UTC)
		return &value, nil
	}

	return nil, fmt.Errorf("invalid due date format, expected RFC3339 or YYYY-MM-DD")
}

func toWorkOrderDTO(work *domain.WorkOrder) *WorkOrderDTO {
	if work == nil {
		return nil
	}

	lines := make([]WorkOrderLineDTO, 0, len(work.Lines))
	for _, line := range work.Lines {
		tasks := make([]WorkOrderTaskDTO, 0, len(line.Tasks))
		for _, task := range line.Tasks {
			tasks = append(tasks, WorkOrderTaskDTO{
				ID:          task.ID,
				TaskID:      task.TaskID,
				Sequence:    task.Sequence,
				Status:      string(task.Status),
				AssignedTo:  task.AssignedTo,
				StartedAt:   task.StartedAt,
				CompletedAt: task.CompletedAt,
				Notes:       task.Notes,
			})
		}

		lines = append(lines, WorkOrderLineDTO{
			ID:             line.ID,
			WorkTypeID:     line.WorkTypeID,
			PositionID:     line.PositionID,
			DesignFilePath: line.DesignFilePath,
			Notes:          line.Notes,
			Sequence:       line.Sequence,
			Tasks:          tasks,
		})
	}

	return &WorkOrderDTO{
		ID:            work.ID,
		WorkNumber:    work.OrderNumber,
		WorkName:      work.OrderName,
		PartyID:       work.PartyID,
		WorkSetupID:   work.WorkSetupID,
		Notes:         work.Notes,
		Status:        string(work.Status),
		Priority:      string(work.Priority),
		StartDate:     work.StartDate,
		DueDate:       work.DueDate,
		CompletedDate: work.CompletedDate,
		Lines:         lines,
	}
}

// enrichWorkOrderDTOs populates SalesOrderID/SalesOrderNumber on a slice of DTOs
// using the optional salesInfoProvider. Nil-safe: no-op when provider is nil.
func (s *MESService) enrichWorkOrderDTOs(ctx context.Context, dtos []WorkOrderDTO) {
	if s.salesInfoProvider == nil || len(dtos) == 0 {
		return
	}
	ids := make([]uuid.UUID, len(dtos))
	for i := range dtos {
		ids[i] = dtos[i].ID
	}
	infoMap, err := s.salesInfoProvider.GetSalesInfoByWorkOrderIDs(ctx, ids)
	if err != nil {
		return // best-effort: don't fail the whole list
	}
	for i := range dtos {
		if info, ok := infoMap[dtos[i].ID]; ok {
			dtos[i].SalesOrderID = &info.SalesOrderID
			dtos[i].SalesOrderNumber = info.SalesOrderNumber
		}
	}
}

// --- WorkSetup ---

func (s *MESService) CreateWorkSetup(ctx context.Context, cmd CreateWorkSetupCommand) (*WorkSetupDTO, error) {
	description := ""
	if cmd.Description != nil {
		description = *cmd.Description
	}
	isActive := true
	if cmd.IsActive != nil {
		isActive = *cmd.IsActive
	}

	lines := make([]domain.WorkSetupLine, 0, len(cmd.Lines))
	for _, l := range cmd.Lines {
		designFilePath := ""
		if l.DesignFilePath != nil {
			designFilePath = *l.DesignFilePath
		}
		notes := ""
		if l.Notes != nil {
			notes = *l.Notes
		}
		lines = append(lines, domain.WorkSetupLine{
			WorkTypeID:     l.WorkTypeID,
			PositionID:     l.PositionID,
			DesignFilePath: designFilePath,
			Notes:          notes,
			Sequence:       l.Sequence,
		})
	}

	ws, err := domain.NewWorkSetup(cmd.Name, cmd.PartyID, cmd.TangibleGroupID, description, isActive, lines)
	if err != nil {
		return nil, err
	}

	if err := s.workSetupRepo.Save(ctx, ws); err != nil {
		return nil, fmt.Errorf("save work setup: %w", err)
	}

	return toWorkSetupDTO(ws), nil
}

func (s *MESService) GetWorkSetupByID(ctx context.Context, query GetWorkSetupByIDQuery) (*WorkSetupDTO, error) {
	ws, err := s.workSetupRepo.FindByID(ctx, query.ID)
	if err != nil {
		return nil, fmt.Errorf("find work setup by id: %w", err)
	}
	if ws == nil {
		return nil, fmt.Errorf("work setup not found")
	}

	return toWorkSetupDTO(ws), nil
}

func (s *MESService) ListWorkSetups(ctx context.Context, query ListWorkSetupsQuery) ([]WorkSetupDTO, error) {
	setups, err := s.workSetupRepo.FindAll(ctx, &domain.WorkSetupFilters{
		IsActive: query.IsActive,
		Search:   query.Search,
		PartyID:  query.PartyID,
	})
	if err != nil {
		return nil, fmt.Errorf("list work setups: %w", err)
	}

	dtos := make([]WorkSetupDTO, 0, len(setups))
	for _, ws := range setups {
		dtos = append(dtos, *toWorkSetupDTO(ws))
	}
	return dtos, nil
}

func (s *MESService) UpdateWorkSetup(ctx context.Context, cmd UpdateWorkSetupCommand) (*WorkSetupDTO, error) {
	ws, err := s.workSetupRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("find work setup for update: %w", err)
	}
	if ws == nil {
		return nil, fmt.Errorf("work setup not found")
	}

	if cmd.Name != nil {
		ws.Name = *cmd.Name
	}
	if cmd.PartyID != nil {
		ws.PartyID = *cmd.PartyID
	}
	if cmd.TangibleGroupID != nil {
		ws.TangibleGroupID = cmd.TangibleGroupID
	}
	if cmd.Description != nil {
		ws.Description = *cmd.Description
	}
	if cmd.IsActive != nil {
		ws.IsActive = *cmd.IsActive
	}
	if cmd.Lines != nil {
		lines := make([]domain.WorkSetupLine, 0, len(cmd.Lines))
		for _, l := range cmd.Lines {
			designFilePath := ""
			if l.DesignFilePath != nil {
				designFilePath = *l.DesignFilePath
			}
			notes := ""
			if l.Notes != nil {
				notes = *l.Notes
			}
			lines = append(lines, domain.WorkSetupLine{
				WorkTypeID:     l.WorkTypeID,
				PositionID:     l.PositionID,
				DesignFilePath: designFilePath,
				Notes:          notes,
				Sequence:       l.Sequence,
			})
		}
		ws.Lines = lines
	}

	if err := s.workSetupRepo.Save(ctx, ws); err != nil {
		return nil, fmt.Errorf("save work setup update: %w", err)
	}

	return toWorkSetupDTO(ws), nil
}

func (s *MESService) DeleteWorkSetup(ctx context.Context, cmd DeleteWorkSetupCommand) error {
	if err := s.workSetupRepo.Delete(ctx, cmd.ID); err != nil {
		return fmt.Errorf("delete work setup: %w", err)
	}
	return nil
}

func toWorkSetupDTO(ws *domain.WorkSetup) *WorkSetupDTO {
	if ws == nil {
		return nil
	}

	lines := make([]WorkSetupLineDTO, 0, len(ws.Lines))
	for _, l := range ws.Lines {
		lines = append(lines, WorkSetupLineDTO{
			ID:             l.ID,
			WorkTypeID:     l.WorkTypeID,
			PositionID:     l.PositionID,
			DesignFilePath: l.DesignFilePath,
			Notes:          l.Notes,
			Sequence:       l.Sequence,
		})
	}

	return &WorkSetupDTO{
		ID:              ws.ID,
		Name:            ws.Name,
		PartyID:         ws.PartyID,
		TangibleGroupID: ws.TangibleGroupID,
		Description:     ws.Description,
		IsActive:        ws.IsActive,
		Lines:           lines,
	}
}
