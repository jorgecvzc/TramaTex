package application

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/joran-cortez/tramatex/internal/mes/domain"
)

type fakeTaskRepo struct {
	saved   *domain.Task
	byID    *domain.Task
	all     []*domain.Task
	deleted uuid.UUID
}

func (r *fakeTaskRepo) Save(_ context.Context, task *domain.Task) error {
	r.saved = task
	return nil
}

func (r *fakeTaskRepo) FindByID(_ context.Context, _ uuid.UUID) (*domain.Task, error) {
	return r.byID, nil
}

func (r *fakeTaskRepo) FindAll(_ context.Context, _ *domain.TaskFilters) ([]*domain.Task, error) {
	return r.all, nil
}

func (r *fakeTaskRepo) Delete(_ context.Context, id uuid.UUID) error {
	r.deleted = id
	return nil
}

func (r *fakeTaskRepo) Exists(_ context.Context, _ uuid.UUID) (bool, error) {
	return true, nil
}

type fakePositionRepo struct {
	saved   *domain.Position
	byID    *domain.Position
	all     []*domain.Position
	deleted uuid.UUID
}

func (r *fakePositionRepo) Save(_ context.Context, position *domain.Position) error {
	r.saved = position
	return nil
}
func (r *fakePositionRepo) FindByID(_ context.Context, _ uuid.UUID) (*domain.Position, error) {
	return r.byID, nil
}
func (r *fakePositionRepo) FindAll(_ context.Context, _ *domain.PositionFilters) ([]*domain.Position, error) {
	return r.all, nil
}
func (r *fakePositionRepo) Delete(_ context.Context, id uuid.UUID) error {
	r.deleted = id
	return nil
}

type fakeServiceGroupRepo struct {
	saved   *domain.ServiceGroup
	byID    *domain.ServiceGroup
	all     []*domain.ServiceGroup
	deleted uuid.UUID
}

func (r *fakeServiceGroupRepo) Save(_ context.Context, serviceGroup *domain.ServiceGroup) error {
	r.saved = serviceGroup
	return nil
}

func (r *fakeServiceGroupRepo) FindByID(_ context.Context, _ uuid.UUID) (*domain.ServiceGroup, error) {
	return r.byID, nil
}

func (r *fakeServiceGroupRepo) FindAll(_ context.Context, _ *domain.ServiceGroupFilters) ([]*domain.ServiceGroup, error) {
	return r.all, nil
}

func (r *fakeServiceGroupRepo) Delete(_ context.Context, id uuid.UUID) error {
	r.deleted = id
	return nil
}

type fakeMESWorkRepo struct {
	saved     *domain.MESWork
	byID      *domain.MESWork
	all       []*domain.MESWork
	yearCount int64
}

func (r *fakeMESWorkRepo) Save(_ context.Context, work *domain.MESWork) error {
	r.saved = work
	return nil
}

func (r *fakeMESWorkRepo) FindByID(_ context.Context, _ uuid.UUID) (*domain.MESWork, error) {
	return r.byID, nil
}

func (r *fakeMESWorkRepo) FindAll(_ context.Context, _ *domain.MESWorkFilters) ([]*domain.MESWork, error) {
	return r.all, nil
}

func (r *fakeMESWorkRepo) CountByYear(_ context.Context, _ int) (int64, error) {
	return r.yearCount, nil
}

func TestCreateTask_AppliesDefaultsAndPersists(t *testing.T) {
	taskRepo := &fakeTaskRepo{}
	service := NewMESService(taskRepo, &fakePositionRepo{}, &fakeServiceGroupRepo{}, &fakeMESWorkRepo{})

	result, err := service.CreateTask(context.Background(), CreateTaskCommand{Name: "Diseñar"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected task dto")
	}
	if result.Name != "Diseñar" {
		t.Fatalf("expected name Diseñar, got %s", result.Name)
	}
	if !result.IsActive {
		t.Fatal("expected default is_active true")
	}
	if result.Description != "" {
		t.Fatalf("expected empty description default, got %s", result.Description)
	}
	if taskRepo.saved == nil {
		t.Fatal("expected repository save call")
	}
}

func TestCreateServiceGroup_MapsTaskAssignments(t *testing.T) {
	repo := &fakeServiceGroupRepo{}
	service := NewMESService(&fakeTaskRepo{}, &fakePositionRepo{}, repo, &fakeMESWorkRepo{})

	taskID := uuid.New()
	result, err := service.CreateServiceGroup(context.Background(), CreateServiceGroupCommand{
		Name: "Serigrafía",
		TaskAssignments: []ServiceGroupTaskInput{
			{TaskID: taskID, Sequence: 1},
		},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected service group dto")
	}
	if len(result.Tasks) != 1 {
		t.Fatalf("expected 1 mapped task, got %d", len(result.Tasks))
	}
	if repo.saved == nil || len(repo.saved.Tasks) != 1 {
		t.Fatal("expected saved entity with mapped assignments")
	}
	if repo.saved.Tasks[0].TaskID != taskID || repo.saved.Tasks[0].Sequence != 1 {
		t.Fatal("expected mapped task assignment to match input")
	}
}

func TestUpdateServiceGroup_DoesNotOverrideTasksWhenNilAssignments(t *testing.T) {
	existingTaskID := uuid.New()
	repo := &fakeServiceGroupRepo{
		byID: &domain.ServiceGroup{
			ID:       uuid.New(),
			Name:     "Bordado",
			IsActive: true,
			Tasks: []domain.ServiceGroupTask{
				{TaskID: existingTaskID, Sequence: 1},
			},
		},
	}
	service := NewMESService(&fakeTaskRepo{}, &fakePositionRepo{}, repo, &fakeMESWorkRepo{})

	newName := "Bordado premium"
	result, err := service.UpdateServiceGroup(context.Background(), UpdateServiceGroupCommand{
		ID:   repo.byID.ID,
		Name: &newName,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Name != newName {
		t.Fatalf("expected updated name, got %s", result.Name)
	}
	if len(repo.saved.Tasks) != 1 || repo.saved.Tasks[0].TaskID != existingTaskID {
		t.Fatal("expected existing tasks to remain unchanged when task_assignments is nil")
	}
}

func TestCreateMESWork_GeneratesTasksFromServiceGroupTemplate(t *testing.T) {
	taskID := uuid.New()
	serviceGroupID := uuid.New()
	positionID := uuid.New()
	workRepo := &fakeMESWorkRepo{yearCount: 2}
	serviceGroupRepo := &fakeServiceGroupRepo{
		byID: &domain.ServiceGroup{
			ID:       serviceGroupID,
			Name:     "Serigrafía",
			IsActive: true,
			Tasks: []domain.ServiceGroupTask{
				{TaskID: taskID, Sequence: 1},
			},
		},
	}
	service := NewMESService(&fakeTaskRepo{}, &fakePositionRepo{}, serviceGroupRepo, workRepo)

	result, err := service.CreateMESWork(context.Background(), CreateMESWorkCommand{
		WorkName:        "Trabajo Cliente A",
		PartyID:         "party-001",
		TangibleGroupID: uuid.New(),
		ServiceGroupAssignments: []CreateMESWorkServiceGroupInput{
			{ServiceGroupID: serviceGroupID, PositionID: positionID, Sequence: 1},
		},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected mes work dto")
	}
	if result.WorkNumber == "" {
		t.Fatal("expected generated work number")
	}
	if workRepo.saved == nil {
		t.Fatal("expected mes work persistence")
	}
	if len(workRepo.saved.ServiceGroups) != 1 {
		t.Fatalf("expected 1 service group, got %d", len(workRepo.saved.ServiceGroups))
	}
	if len(workRepo.saved.ServiceGroups[0].Tasks) != 1 {
		t.Fatalf("expected 1 generated task, got %d", len(workRepo.saved.ServiceGroups[0].Tasks))
	}
	if workRepo.saved.ServiceGroups[0].Tasks[0].TaskID != taskID {
		t.Fatal("expected generated task to match service group template task")
	}
}

func TestTaskFlows_ListUpdateDelete(t *testing.T) {
	taskID := uuid.New()
	taskRepo := &fakeTaskRepo{
		byID: &domain.Task{ID: taskID, Name: "Diseñar", Description: "old", IsActive: true},
		all:  []*domain.Task{{ID: taskID, Name: "Diseñar", IsActive: true}},
	}
	service := NewMESService(taskRepo, &fakePositionRepo{}, &fakeServiceGroupRepo{}, &fakeMESWorkRepo{})

	list, err := service.ListTasks(context.Background(), ListTasksQuery{})
	if err != nil || len(list) != 1 {
		t.Fatalf("expected list tasks success, err=%v len=%d", err, len(list))
	}

	name := "Diseñar Premium"
	isActive := false
	updated, err := service.UpdateTask(context.Background(), UpdateTaskCommand{ID: taskID, Name: &name, IsActive: &isActive})
	if err != nil {
		t.Fatalf("expected update success, got %v", err)
	}
	if updated.Name != name || updated.IsActive != isActive {
		t.Fatal("expected updated task values")
	}

	if err := service.DeleteTask(context.Background(), DeleteTaskCommand{ID: taskID}); err != nil {
		t.Fatalf("expected delete success, got %v", err)
	}
	if taskRepo.deleted != taskID {
		t.Fatal("expected deleted id to match")
	}
}

func TestPositionFlows_CreateGetListUpdateDelete(t *testing.T) {
	positionID := uuid.New()
	positionRepo := &fakePositionRepo{
		byID: &domain.Position{ID: positionID, Name: "Espalda", Code: "BACK", IsActive: true},
		all:  []*domain.Position{{ID: positionID, Name: "Espalda", Code: "BACK", IsActive: true}},
	}
	service := NewMESService(&fakeTaskRepo{}, positionRepo, &fakeServiceGroupRepo{}, &fakeMESWorkRepo{})

	created, err := service.CreatePosition(context.Background(), CreatePositionCommand{Name: "Pecho", Code: "CHEST"})
	if err != nil || created == nil {
		t.Fatalf("expected create position success, err=%v", err)
	}

	got, err := service.GetPositionByID(context.Background(), GetPositionByIDQuery{ID: positionID})
	if err != nil || got == nil {
		t.Fatalf("expected get position success, err=%v", err)
	}

	list, err := service.ListPositions(context.Background(), ListPositionsQuery{})
	if err != nil || len(list) != 1 {
		t.Fatalf("expected list positions success, err=%v len=%d", err, len(list))
	}

	name := "Espalda Completa"
	updated, err := service.UpdatePosition(context.Background(), UpdatePositionCommand{ID: positionID, Name: &name})
	if err != nil || updated.Name != name {
		t.Fatalf("expected update position success, err=%v", err)
	}

	if err := service.DeletePosition(context.Background(), DeletePositionCommand{ID: positionID}); err != nil {
		t.Fatalf("expected delete position success, got %v", err)
	}
	if positionRepo.deleted != positionID {
		t.Fatal("expected deleted position id to match")
	}
}

func TestServiceGroupFlows_ListGetDelete(t *testing.T) {
	id := uuid.New()
	repo := &fakeServiceGroupRepo{
		byID: &domain.ServiceGroup{ID: id, Name: "Bordado", IsActive: true},
		all:  []*domain.ServiceGroup{{ID: id, Name: "Bordado", IsActive: true}},
	}
	service := NewMESService(&fakeTaskRepo{}, &fakePositionRepo{}, repo, &fakeMESWorkRepo{})

	got, err := service.GetServiceGroupByID(context.Background(), GetServiceGroupByIDQuery{ID: id})
	if err != nil || got == nil {
		t.Fatalf("expected get service group success, err=%v", err)
	}

	list, err := service.ListServiceGroups(context.Background(), ListServiceGroupsQuery{})
	if err != nil || len(list) != 1 {
		t.Fatalf("expected list service groups success, err=%v len=%d", err, len(list))
	}

	if err := service.DeleteServiceGroup(context.Background(), DeleteServiceGroupCommand{ID: id}); err != nil {
		t.Fatalf("expected delete service group success, got %v", err)
	}
	if repo.deleted != id {
		t.Fatal("expected deleted service group id to match")
	}
}

func TestMESWorkFlows_ListAndGet(t *testing.T) {
	id := uuid.New()
	workRepo := &fakeMESWorkRepo{
		byID: &domain.MESWork{ID: id, WorkNumber: "MES-2026-001", WorkName: "Trabajo", PartyID: "party-1", TangibleGroupID: uuid.New(), Status: domain.ProductionStatusDraft, Priority: domain.WorkPriorityNormal},
		all:  []*domain.MESWork{{ID: id, WorkNumber: "MES-2026-001", WorkName: "Trabajo", PartyID: "party-1", TangibleGroupID: uuid.New(), Status: domain.ProductionStatusDraft, Priority: domain.WorkPriorityNormal}},
	}
	service := NewMESService(&fakeTaskRepo{}, &fakePositionRepo{}, &fakeServiceGroupRepo{}, workRepo)

	list, err := service.ListMESWorks(context.Background(), ListMESWorksQuery{})
	if err != nil || len(list) != 1 {
		t.Fatalf("expected list mes works success, err=%v len=%d", err, len(list))
	}

	got, err := service.GetMESWorkByID(context.Background(), GetMESWorkByIDQuery{ID: id})
	if err != nil || got == nil {
		t.Fatalf("expected get mes work success, err=%v", err)
	}
}

func TestNotFoundBranches_ReturnError(t *testing.T) {
	service := NewMESService(&fakeTaskRepo{}, &fakePositionRepo{}, &fakeServiceGroupRepo{}, &fakeMESWorkRepo{})

	if _, err := service.GetTaskByID(context.Background(), GetTaskByIDQuery{ID: uuid.New()}); err == nil {
		t.Fatal("expected task not found error")
	}

	if _, err := service.GetPositionByID(context.Background(), GetPositionByIDQuery{ID: uuid.New()}); err == nil {
		t.Fatal("expected position not found error")
	}

	if _, err := service.GetServiceGroupByID(context.Background(), GetServiceGroupByIDQuery{ID: uuid.New()}); err == nil {
		t.Fatal("expected service group not found error")
	}

	if _, err := service.GetMESWorkByID(context.Background(), GetMESWorkByIDQuery{ID: uuid.New()}); err == nil {
		t.Fatal("expected mes work not found error")
	}
}

func TestMESDashboardStatsAndOverdue(t *testing.T) {
	today := time.Now().UTC()
	yesterday := today.Add(-24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)

	overdueWork := &domain.MESWork{
		ID:              uuid.New(),
		WorkNumber:      "MES-2026-010",
		WorkName:        "Trabajo vencido",
		PartyID:         "party-1",
		TangibleGroupID: uuid.New(),
		Status:          domain.ProductionStatusInProgress,
		Priority:        domain.WorkPriorityHigh,
		DueDate:         &yesterday,
	}

	dueTodayWork := &domain.MESWork{
		ID:              uuid.New(),
		WorkNumber:      "MES-2026-011",
		WorkName:        "Trabajo hoy",
		PartyID:         "party-1",
		TangibleGroupID: uuid.New(),
		Status:          domain.ProductionStatusDraft,
		Priority:        domain.WorkPriorityNormal,
		DueDate:         &today,
	}

	notOverdueCompleted := &domain.MESWork{
		ID:              uuid.New(),
		WorkNumber:      "MES-2026-012",
		WorkName:        "Trabajo completado",
		PartyID:         "party-1",
		TangibleGroupID: uuid.New(),
		Status:          domain.ProductionStatusCompleted,
		Priority:        domain.WorkPriorityLow,
		DueDate:         &yesterday,
	}

	futureWork := &domain.MESWork{
		ID:              uuid.New(),
		WorkNumber:      "MES-2026-013",
		WorkName:        "Trabajo futuro",
		PartyID:         "party-1",
		TangibleGroupID: uuid.New(),
		Status:          domain.ProductionStatusPending,
		Priority:        domain.WorkPriorityNormal,
		DueDate:         &tomorrow,
	}

	workRepo := &fakeMESWorkRepo{
		all: []*domain.MESWork{overdueWork, dueTodayWork, notOverdueCompleted, futureWork},
	}
	service := NewMESService(&fakeTaskRepo{}, &fakePositionRepo{}, &fakeServiceGroupRepo{}, workRepo)

	stats, err := service.GetMESWorkDashboardStats(context.Background())
	if err != nil {
		t.Fatalf("expected dashboard stats success, got %v", err)
	}
	if stats.Total != 4 {
		t.Fatalf("expected total 4, got %d", stats.Total)
	}
	if stats.Overdue != 1 {
		t.Fatalf("expected overdue 1, got %d", stats.Overdue)
	}
	if stats.ByStatus[string(domain.ProductionStatusInProgress)] != 1 {
		t.Fatal("expected IN_PROGRESS count 1")
	}

	overdue, err := service.ListOverdueMESWorks(context.Background(), ListOverdueMESWorksQuery{Limit: 5})
	if err != nil {
		t.Fatalf("expected overdue list success, got %v", err)
	}
	if len(overdue) != 1 {
		t.Fatalf("expected 1 overdue work, got %d", len(overdue))
	}
	if overdue[0].ID != overdueWork.ID {
		t.Fatal("expected overdue work id to match")
	}
}

func TestUpdateMESWorkTaskStatus_StartAndCompleteFlow(t *testing.T) {
	taskID := uuid.New()
	workID := uuid.New()
	serviceGroupID := uuid.New()

	workRepo := &fakeMESWorkRepo{
		byID: &domain.MESWork{
			ID:              workID,
			WorkNumber:      "MES-2026-020",
			WorkName:        "Trabajo tablet",
			PartyID:         "party-1",
			TangibleGroupID: uuid.New(),
			Status:          domain.ProductionStatusPending,
			Priority:        domain.WorkPriorityNormal,
			ServiceGroups: []domain.MESWorkServiceGroup{
				{
					ID:             serviceGroupID,
					ServiceGroupID: uuid.New(),
					PositionID:     uuid.New(),
					Sequence:       1,
					Tasks: []domain.MESWorkTask{
						{ID: taskID, TaskID: uuid.New(), Sequence: 1, Status: domain.TaskStatusPending},
					},
				},
			},
		},
	}

	service := NewMESService(&fakeTaskRepo{}, &fakePositionRepo{}, &fakeServiceGroupRepo{}, workRepo)

	startResult, err := service.UpdateMESWorkTaskStatus(context.Background(), UpdateMESWorkTaskStatusCommand{
		WorkID: workID,
		TaskID: taskID,
		Action: "START",
	})
	if err != nil {
		t.Fatalf("expected start action success, got %v", err)
	}
	if startResult.Status != string(domain.ProductionStatusInProgress) {
		t.Fatalf("expected work status IN_PROGRESS after start, got %s", startResult.Status)
	}
	if startResult.ServiceGroups[0].Tasks[0].StartedAt == nil {
		t.Fatal("expected started_at to be set on START")
	}

	completeResult, err := service.UpdateMESWorkTaskStatus(context.Background(), UpdateMESWorkTaskStatusCommand{
		WorkID: workID,
		TaskID: taskID,
		Action: "COMPLETE",
	})
	if err != nil {
		t.Fatalf("expected complete action success, got %v", err)
	}
	if completeResult.Status != string(domain.ProductionStatusCompleted) {
		t.Fatalf("expected work status COMPLETED after complete, got %s", completeResult.Status)
	}
	if completeResult.CompletedDate == nil {
		t.Fatal("expected work completed_date after task completion")
	}
	if workRepo.saved == nil {
		t.Fatal("expected mes work saved after status update")
	}
}
