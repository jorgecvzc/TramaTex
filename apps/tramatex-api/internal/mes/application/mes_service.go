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

// MESService provides operations for Manufacturing Execution System
// focused on foundation master data CRUD.
type MESService struct {
	taskRepo         domain.TaskRepository
	positionRepo     domain.PositionRepository
	serviceGroupRepo domain.ServiceGroupRepository
	mesWorkRepo      domain.MESWorkRepository
}

// NewMESService creates a new MES service.
func NewMESService(
	taskRepo domain.TaskRepository,
	positionRepo domain.PositionRepository,
	serviceGroupRepo domain.ServiceGroupRepository,
	mesWorkRepo domain.MESWorkRepository,
) *MESService {
	return &MESService{
		taskRepo:         taskRepo,
		positionRepo:     positionRepo,
		serviceGroupRepo: serviceGroupRepo,
		mesWorkRepo:      mesWorkRepo,
	}
}

func (s *MESService) CreateTask(ctx context.Context, cmd CreateTaskCommand) (*TaskDTO, error) {
	description := ""
	if cmd.Description != nil {
		description = *cmd.Description
	}
	isActive := true
	if cmd.IsActive != nil {
		isActive = *cmd.IsActive
	}

	task, err := domain.NewTask(cmd.Name, description, isActive)
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

func (s *MESService) CreateServiceGroup(ctx context.Context, cmd CreateServiceGroupCommand) (*ServiceGroupDTO, error) {
	description := ""
	if cmd.Description != nil {
		description = *cmd.Description
	}
	isActive := true
	if cmd.IsActive != nil {
		isActive = *cmd.IsActive
	}
	assignments := mapTaskAssignments(cmd.TaskAssignments)

	serviceGroup, err := domain.NewServiceGroup(cmd.Name, description, cmd.ProductGroupID, isActive, assignments)
	if err != nil {
		return nil, err
	}

	if err := s.serviceGroupRepo.Save(ctx, serviceGroup); err != nil {
		return nil, fmt.Errorf("save service group: %w", err)
	}

	return toServiceGroupDTO(serviceGroup), nil
}

func (s *MESService) CreateServiceTemplate(ctx context.Context, cmd CreateServiceTemplateCommand) (*ServiceTemplateDTO, error) {
	result, err := s.CreateServiceGroup(ctx, CreateServiceGroupCommand(cmd))
	if err != nil {
		return nil, err
	}
	return (*ServiceTemplateDTO)(result), nil
}

func (s *MESService) GetServiceGroupByID(ctx context.Context, query GetServiceGroupByIDQuery) (*ServiceGroupDTO, error) {
	serviceGroup, err := s.serviceGroupRepo.FindByID(ctx, query.ID)
	if err != nil {
		return nil, fmt.Errorf("find service group by id: %w", err)
	}
	if serviceGroup == nil {
		return nil, fmt.Errorf("service group not found")
	}

	return toServiceGroupDTO(serviceGroup), nil
}

func (s *MESService) GetServiceTemplateByID(ctx context.Context, query GetServiceTemplateByIDQuery) (*ServiceTemplateDTO, error) {
	result, err := s.GetServiceGroupByID(ctx, GetServiceGroupByIDQuery(query))
	if err != nil {
		return nil, err
	}
	return (*ServiceTemplateDTO)(result), nil
}

func (s *MESService) ListServiceGroups(ctx context.Context, query ListServiceGroupsQuery) ([]ServiceGroupDTO, error) {
	serviceGroups, err := s.serviceGroupRepo.FindAll(ctx, &domain.ServiceGroupFilters{
		IsActive: query.IsActive,
		Search:   query.Search,
	})
	if err != nil {
		return nil, fmt.Errorf("list service groups: %w", err)
	}

	dtos := make([]ServiceGroupDTO, 0, len(serviceGroups))
	for _, serviceGroup := range serviceGroups {
		dtos = append(dtos, *toServiceGroupDTO(serviceGroup))
	}
	return dtos, nil
}

func (s *MESService) ListServiceTemplates(ctx context.Context, query ListServiceTemplatesQuery) ([]ServiceTemplateDTO, error) {
	results, err := s.ListServiceGroups(ctx, ListServiceGroupsQuery(query))
	if err != nil {
		return nil, err
	}

	aliases := make([]ServiceTemplateDTO, 0, len(results))
	for _, result := range results {
		aliases = append(aliases, ServiceTemplateDTO(result))
	}

	return aliases, nil
}

func (s *MESService) UpdateServiceGroup(ctx context.Context, cmd UpdateServiceGroupCommand) (*ServiceGroupDTO, error) {
	serviceGroup, err := s.serviceGroupRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("find service group for update: %w", err)
	}
	if serviceGroup == nil {
		return nil, fmt.Errorf("service group not found")
	}

	if cmd.Name != nil {
		serviceGroup.Name = *cmd.Name
	}
	if cmd.Description != nil {
		serviceGroup.Description = *cmd.Description
	}
	// ProductGroupID was removed in the domain refactor (WorkType no longer has it).
	// The command field is kept for API compatibility but is no longer stored.
	if cmd.IsActive != nil {
		serviceGroup.IsActive = *cmd.IsActive
	}
	if cmd.TaskAssignments != nil {
		serviceGroup.Tasks = mapTaskAssignments(cmd.TaskAssignments)
	}

	if err := s.serviceGroupRepo.Save(ctx, serviceGroup); err != nil {
		return nil, fmt.Errorf("save service group update: %w", err)
	}

	return toServiceGroupDTO(serviceGroup), nil
}

func (s *MESService) UpdateServiceTemplate(ctx context.Context, cmd UpdateServiceTemplateCommand) (*ServiceTemplateDTO, error) {
	result, err := s.UpdateServiceGroup(ctx, UpdateServiceGroupCommand(cmd))
	if err != nil {
		return nil, err
	}
	return (*ServiceTemplateDTO)(result), nil
}

func (s *MESService) DeleteServiceGroup(ctx context.Context, cmd DeleteServiceGroupCommand) error {
	if err := s.serviceGroupRepo.Delete(ctx, cmd.ID); err != nil {
		return fmt.Errorf("delete service group: %w", err)
	}
	return nil
}

func (s *MESService) DeleteServiceTemplate(ctx context.Context, cmd DeleteServiceTemplateCommand) error {
	return s.DeleteServiceGroup(ctx, DeleteServiceGroupCommand(cmd))
}

func mapTaskAssignments(inputs []ServiceGroupTaskInput) []domain.ServiceGroupTask {
	if len(inputs) == 0 {
		return []domain.ServiceGroupTask{}
	}

	assignments := make([]domain.ServiceGroupTask, 0, len(inputs))
	for _, input := range inputs {
		assignments = append(assignments, domain.ServiceGroupTask{
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

func toServiceGroupDTO(serviceGroup *domain.ServiceGroup) *ServiceGroupDTO {
	if serviceGroup == nil {
		return nil
	}

	tasks := make([]ServiceGroupTaskDTO, 0, len(serviceGroup.Tasks))
	for _, task := range serviceGroup.Tasks {
		tasks = append(tasks, ServiceGroupTaskDTO{
			TaskID:   task.TaskID,
			Sequence: task.Sequence,
		})
	}

	return &ServiceGroupDTO{
		ID:             serviceGroup.ID,
		Name:           serviceGroup.Name,
		Description:    serviceGroup.Description,
		ProductGroupID: nil, // Removed in domain refactor; kept in DTO for API compatibility.
		IsActive:       serviceGroup.IsActive,
		Tasks:          tasks,
	}
}

func (s *MESService) CreateMESWork(ctx context.Context, cmd CreateMESWorkCommand) (*MESWorkDTO, error) {
	if len(cmd.ServiceGroupAssignments) == 0 {
		return nil, fmt.Errorf("at least one service group assignment is required")
	}

	status := domain.ProductionStatusDraft
	if cmd.Status != nil && *cmd.Status != "" {
		status = domain.ProductionStatus(strings.ToUpper(*cmd.Status))
	}

	priority := domain.WorkPriorityNormal
	if cmd.Priority != nil && *cmd.Priority != "" {
		priority = domain.WorkPriority(strings.ToUpper(*cmd.Priority))
	}

	year := time.Now().UTC().Year()
	count, err := s.mesWorkRepo.CountByYear(ctx, year)
	if err != nil {
		return nil, fmt.Errorf("count works by year: %w", err)
	}
	workNumber := fmt.Sprintf("MES-%d-%03d", year, count+1)

	garmentNotes := ""
	if cmd.GarmentNotes != nil {
		garmentNotes = *cmd.GarmentNotes
	}

	groups := make([]domain.MESWorkServiceGroup, 0, len(cmd.ServiceGroupAssignments))
	for _, assignment := range cmd.ServiceGroupAssignments {
		serviceGroup, err := s.serviceGroupRepo.FindByID(ctx, assignment.ServiceGroupID)
		if err != nil {
			return nil, fmt.Errorf("find service group for mes work: %w", err)
		}
		if serviceGroup == nil {
			return nil, fmt.Errorf("service group not found")
		}

		designFilePath := ""
		if assignment.DesignFilePath != nil {
			designFilePath = *assignment.DesignFilePath
		}
		notes := ""
		if assignment.Notes != nil {
			notes = *assignment.Notes
		}

		generatedTasks := make([]domain.MESWorkTask, 0, len(serviceGroup.Tasks))
		for _, taskAssignment := range serviceGroup.Tasks {
			generatedTasks = append(generatedTasks, domain.MESWorkTask{
				ID:       uuid.New(),
				TaskID:   taskAssignment.TaskID,
				Sequence: taskAssignment.Sequence,
				Status:   domain.TaskStatusPending,
			})
		}

		groups = append(groups, domain.MESWorkServiceGroup{
			ID:             uuid.New(),
			WorkTypeID:     assignment.ServiceGroupID,
			PositionID:     assignment.PositionID,
			DesignFilePath: designFilePath,
			Notes:          notes,
			Sequence:       assignment.Sequence,
			Tasks:          generatedTasks,
		})
	}

	work, err := domain.NewMESWork(
		workNumber,
		cmd.WorkName,
		cmd.PartyID,
		cmd.TangibleGroupID,
		garmentNotes,
		status,
		priority,
		nil,
		nil,
		nil,
		groups,
	)
	if err != nil {
		return nil, err
	}

	if err := s.mesWorkRepo.Save(ctx, work); err != nil {
		return nil, fmt.Errorf("save mes work: %w", err)
	}

	return toMESWorkDTO(work), nil
}

func (s *MESService) CreateWorkDefinition(ctx context.Context, cmd CreateWorkDefinitionCommand) (*MESWorkDefinitionDTO, error) {
	result, err := s.CreateMESWork(ctx, CreateMESWorkCommand(cmd))
	if err != nil {
		return nil, err
	}
	return (*MESWorkDefinitionDTO)(result), nil
}

func (s *MESService) UpdateMESWork(ctx context.Context, cmd UpdateMESWorkCommand) (*MESWorkDTO, error) {
	work, err := s.mesWorkRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("find mes work for update: %w", err)
	}
	if work == nil {
		return nil, fmt.Errorf("mes work not found")
	}

	if cmd.WorkName != nil {
		trimmed := strings.TrimSpace(*cmd.WorkName)
		if trimmed == "" {
			return nil, fmt.Errorf("work name is required")
		}
		work.OrderName = trimmed
	}

	if cmd.PartyID != nil {
		trimmed := strings.TrimSpace(*cmd.PartyID)
		if trimmed == "" {
			return nil, fmt.Errorf("party id is required")
		}
		work.PartyID = trimmed
	}

	if cmd.TangibleGroupID != nil {
		if *cmd.TangibleGroupID == uuid.Nil {
			return nil, fmt.Errorf("tangible group id is required")
		}
		work.TangibleGroupID = *cmd.TangibleGroupID
	}

	if cmd.GarmentNotes != nil {
		work.GarmentNotes = strings.TrimSpace(*cmd.GarmentNotes)
	}

	if cmd.Status != nil {
		status := domain.ProductionStatus(strings.ToUpper(strings.TrimSpace(*cmd.Status)))
		if !status.IsValid() {
			return nil, fmt.Errorf("invalid production status")
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

	if err := s.mesWorkRepo.Save(ctx, work); err != nil {
		return nil, fmt.Errorf("save mes work update: %w", err)
	}

	return toMESWorkDTO(work), nil
}

func (s *MESService) UpdateWorkDefinition(ctx context.Context, cmd UpdateWorkDefinitionCommand) (*MESWorkDefinitionDTO, error) {
	result, err := s.UpdateMESWork(ctx, UpdateMESWorkCommand(cmd))
	if err != nil {
		return nil, err
	}
	return (*MESWorkDefinitionDTO)(result), nil
}

func (s *MESService) UpdateMESWorkTaskStatus(ctx context.Context, cmd UpdateMESWorkTaskStatusCommand) (*MESWorkDTO, error) {
	work, err := s.mesWorkRepo.FindByID(ctx, cmd.WorkID)
	if err != nil {
		return nil, fmt.Errorf("find mes work for task update: %w", err)
	}
	if work == nil {
		return nil, fmt.Errorf("mes work not found")
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
		return nil, fmt.Errorf("mes work task not found")
	}

	recalculateWorkStatus(work, now)

	if err := s.mesWorkRepo.Save(ctx, work); err != nil {
		return nil, fmt.Errorf("save mes work task update: %w", err)
	}

	return toMESWorkDTO(work), nil
}

func (s *MESService) UpdateWorkDefinitionTaskStatus(ctx context.Context, cmd UpdateWorkDefinitionTaskStatusCommand) (*MESWorkDefinitionDTO, error) {
	result, err := s.UpdateMESWorkTaskStatus(ctx, UpdateMESWorkTaskStatusCommand(cmd))
	if err != nil {
		return nil, err
	}
	return (*MESWorkDefinitionDTO)(result), nil
}

func (s *MESService) ListMESWorks(ctx context.Context, query ListMESWorksQuery) ([]MESWorkDTO, error) {
	var status *domain.ProductionStatus
	if query.Status != nil && *query.Status != "" {
		s := domain.ProductionStatus(strings.ToUpper(*query.Status))
		status = &s
	}

	works, err := s.mesWorkRepo.FindAll(ctx, &domain.MESWorkFilters{
		Status:  status,
		Search:  query.Search,
		PartyID: query.PartyID,
	})
	if err != nil {
		return nil, fmt.Errorf("list mes works: %w", err)
	}

	result := make([]MESWorkDTO, 0, len(works))
	for _, work := range works {
		result = append(result, *toMESWorkDTO(work))
	}
	return result, nil
}

func (s *MESService) ListWorkDefinitions(ctx context.Context, query ListWorkDefinitionsQuery) ([]MESWorkDefinitionDTO, error) {
	results, err := s.ListMESWorks(ctx, ListMESWorksQuery(query))
	if err != nil {
		return nil, err
	}

	aliases := make([]MESWorkDefinitionDTO, 0, len(results))
	for _, result := range results {
		aliases = append(aliases, MESWorkDefinitionDTO(result))
	}

	return aliases, nil
}

func (s *MESService) GetMESWorkByID(ctx context.Context, query GetMESWorkByIDQuery) (*MESWorkDTO, error) {
	work, err := s.mesWorkRepo.FindByID(ctx, query.ID)
	if err != nil {
		return nil, fmt.Errorf("find mes work by id: %w", err)
	}
	if work == nil {
		return nil, fmt.Errorf("mes work not found")
	}
	return toMESWorkDTO(work), nil
}

func (s *MESService) GetWorkDefinitionByID(ctx context.Context, query GetWorkDefinitionByIDQuery) (*MESWorkDefinitionDTO, error) {
	result, err := s.GetMESWorkByID(ctx, GetMESWorkByIDQuery(query))
	if err != nil {
		return nil, err
	}
	return (*MESWorkDefinitionDTO)(result), nil
}

func (s *MESService) GetMESWorkDashboardStats(ctx context.Context) (*MESWorkDashboardStatsDTO, error) {
	works, err := s.mesWorkRepo.FindAll(ctx, &domain.MESWorkFilters{})
	if err != nil {
		return nil, fmt.Errorf("get mes work dashboard stats: %w", err)
	}

	now := time.Now().UTC()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	tomorrowStart := todayStart.Add(24 * time.Hour)

	stats := &MESWorkDashboardStatsDTO{
		Total:    len(works),
		ByStatus: map[string]int{},
	}

	for _, work := range works {
		stats.ByStatus[string(work.Status)]++

		if work.DueDate != nil {
			if work.DueDate.Before(todayStart) && isOpenMESWorkStatus(work.Status) {
				stats.Overdue++
			}
			if (work.DueDate.Equal(todayStart) || work.DueDate.After(todayStart)) && work.DueDate.Before(tomorrowStart) {
				stats.DueToday++
			}
		}
	}

	return stats, nil
}

func (s *MESService) GetWorkDefinitionDashboardStats(ctx context.Context) (*MESWorkDashboardStatsDTO, error) {
	return s.GetMESWorkDashboardStats(ctx)
}

func (s *MESService) ListOverdueMESWorks(ctx context.Context, query ListOverdueMESWorksQuery) ([]MESWorkDTO, error) {
	works, err := s.mesWorkRepo.FindAll(ctx, &domain.MESWorkFilters{})
	if err != nil {
		return nil, fmt.Errorf("list overdue mes works: %w", err)
	}

	todayStart := time.Date(time.Now().UTC().Year(), time.Now().UTC().Month(), time.Now().UTC().Day(), 0, 0, 0, 0, time.UTC)
	overdue := make([]*domain.MESWork, 0)
	for _, work := range works {
		if work.DueDate != nil && work.DueDate.Before(todayStart) && isOpenMESWorkStatus(work.Status) {
			overdue = append(overdue, work)
		}
	}

	sort.SliceStable(overdue, func(i, j int) bool {
		return overdue[i].DueDate.Before(*overdue[j].DueDate)
	})

	if query.Limit > 0 && len(overdue) > query.Limit {
		overdue = overdue[:query.Limit]
	}

	result := make([]MESWorkDTO, 0, len(overdue))
	for _, work := range overdue {
		result = append(result, *toMESWorkDTO(work))
	}

	return result, nil
}

func (s *MESService) ListOverdueWorkDefinitions(ctx context.Context, query ListOverdueWorkDefinitionsQuery) ([]MESWorkDefinitionDTO, error) {
	results, err := s.ListOverdueMESWorks(ctx, ListOverdueMESWorksQuery(query))
	if err != nil {
		return nil, err
	}

	aliases := make([]MESWorkDefinitionDTO, 0, len(results))
	for _, result := range results {
		aliases = append(aliases, MESWorkDefinitionDTO(result))
	}

	return aliases, nil
}

func isOpenMESWorkStatus(status domain.ProductionStatus) bool {
	return status != domain.ProductionStatusCompleted && status != domain.ProductionStatusCancelled
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

func recalculateWorkStatus(work *domain.MESWork, now time.Time) {
	if work == nil {
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

func toMESWorkDTO(work *domain.MESWork) *MESWorkDTO {
	if work == nil {
		return nil
	}

	groups := make([]MESWorkServiceGroupDTO, 0, len(work.Lines))
	for _, group := range work.Lines {
		tasks := make([]MESWorkTaskDTO, 0, len(group.Tasks))
		for _, task := range group.Tasks {
			tasks = append(tasks, MESWorkTaskDTO{
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

		groups = append(groups, MESWorkServiceGroupDTO{
			ID:             group.ID,
			ServiceGroupID: group.WorkTypeID,
			PositionID:     group.PositionID,
			DesignFilePath: group.DesignFilePath,
			Notes:          group.Notes,
			Sequence:       group.Sequence,
			Tasks:          tasks,
		})
	}

	return &MESWorkDTO{
		ID:              work.ID,
		WorkNumber:      work.OrderNumber,
		WorkName:        work.OrderName,
		PartyID:         work.PartyID,
		TangibleGroupID: work.TangibleGroupID,
		GarmentNotes:    work.GarmentNotes,
		Status:          string(work.Status),
		Priority:        string(work.Priority),
		StartDate:       work.StartDate,
		DueDate:         work.DueDate,
		CompletedDate:   work.CompletedDate,
		ServiceGroups:   groups,
	}
}
