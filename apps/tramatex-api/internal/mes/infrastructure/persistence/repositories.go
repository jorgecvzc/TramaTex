package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/joran-cortez/tramatex/internal/mes/domain"
	"gorm.io/gorm"
)

type TaskDataModel struct {
	ID          uuid.UUID `gorm:"column:id;type:uuid;primaryKey"`
	Name        string    `gorm:"column:name"`
	Reference   string    `gorm:"column:reference"`
	Description string    `gorm:"column:description"`
	IsActive    bool      `gorm:"column:is_active"`
}

func (TaskDataModel) TableName() string { return "tasks" }

type PositionDataModel struct {
	ID          uuid.UUID `gorm:"column:id;type:uuid;primaryKey"`
	Name        string    `gorm:"column:name"`
	Code        string    `gorm:"column:code"`
	Description string    `gorm:"column:description"`
	IsActive    bool      `gorm:"column:is_active"`
}

func (PositionDataModel) TableName() string { return "positions" }

type WorkTypeDataModel struct {
	ID          uuid.UUID               `gorm:"column:id;type:uuid;primaryKey"`
	Name        string                  `gorm:"column:name"`
	Reference   string                  `gorm:"column:reference"`
	Description string                  `gorm:"column:description"`
	IsActive    bool                    `gorm:"column:is_active"`
	Tasks       []WorkTypeTaskDataModel `gorm:"foreignKey:WorkTypeID;references:ID"`
}

func (WorkTypeDataModel) TableName() string { return "work_types" }

type WorkTypeTaskDataModel struct {
	WorkTypeID uuid.UUID `gorm:"column:work_type_id;type:uuid;primaryKey"`
	TaskID     uuid.UUID `gorm:"column:task_id;type:uuid;primaryKey"`
	Sequence   int       `gorm:"column:sequence"`
}

func (WorkTypeTaskDataModel) TableName() string { return "work_type_tasks" }

type WorkOrderDataModel struct {
	ID            uuid.UUID                `gorm:"column:id;type:uuid;primaryKey"`
	WorkNumber    string                   `gorm:"column:work_number"`
	WorkName      string                   `gorm:"column:work_name"`
	PartyID       string                   `gorm:"column:party_id"`
	WorkSetupID   *uuid.UUID               `gorm:"column:work_setup_id;type:uuid"`
	Notes         string                   `gorm:"column:notes"`
	Status        string                   `gorm:"column:status"`
	Priority      string                   `gorm:"column:priority"`
	StartDate     *time.Time               `gorm:"column:start_date"`
	DueDate       *time.Time               `gorm:"column:due_date"`
	CompletedDate *time.Time               `gorm:"column:completed_date"`
	Lines         []WorkOrderLineDataModel `gorm:"foreignKey:WorkOrderID;references:ID"`
}

func (WorkOrderDataModel) TableName() string { return "work_orders" }

type WorkOrderLineDataModel struct {
	ID             uuid.UUID                `gorm:"column:id;type:uuid;primaryKey"`
	WorkOrderID    uuid.UUID                `gorm:"column:work_order_id;type:uuid"`
	WorkTypeID     uuid.UUID                `gorm:"column:work_type_id;type:uuid"`
	PositionID     uuid.UUID                `gorm:"column:position_id;type:uuid"`
	DesignFilePath string                   `gorm:"column:design_file_path"`
	Notes          string                   `gorm:"column:notes"`
	Sequence       int                      `gorm:"column:sequence"`
	Tasks          []WorkOrderTaskDataModel `gorm:"foreignKey:WorkOrderLineID;references:ID"`
}

func (WorkOrderLineDataModel) TableName() string { return "work_order_lines" }

type WorkOrderTaskDataModel struct {
	ID              uuid.UUID  `gorm:"column:id;type:uuid;primaryKey"`
	WorkOrderLineID uuid.UUID  `gorm:"column:work_order_line_id;type:uuid"`
	TaskID          uuid.UUID  `gorm:"column:task_id;type:uuid"`
	Sequence        int        `gorm:"column:sequence"`
	Status          string     `gorm:"column:status"`
	AssignedTo      *uuid.UUID `gorm:"column:assigned_to;type:uuid"`
	StartedAt       *time.Time `gorm:"column:started_at"`
	CompletedAt     *time.Time `gorm:"column:completed_at"`
	Notes           string     `gorm:"column:notes"`
}

func (WorkOrderTaskDataModel) TableName() string { return "work_order_tasks" }

type GORMTaskRepository struct {
	db *gorm.DB
}

func NewGORMTaskRepository(db *gorm.DB) *GORMTaskRepository {
	return &GORMTaskRepository{db: db}
}

func (r *GORMTaskRepository) Save(ctx context.Context, task *domain.Task) error {
	data := TaskDataModel{
		ID:          task.ID,
		Name:        task.Name,
		Reference:   task.Reference,
		Description: task.Description,
		IsActive:    task.IsActive,
	}
	return r.db.WithContext(ctx).Save(&data).Error
}

func (r *GORMTaskRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	var data TaskDataModel
	err := r.db.WithContext(ctx).First(&data, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &domain.Task{
		ID:          data.ID,
		Name:        data.Name,
		Reference:   data.Reference,
		Description: data.Description,
		IsActive:    data.IsActive,
	}, nil
}

func (r *GORMTaskRepository) FindAll(ctx context.Context, filters *domain.TaskFilters) ([]*domain.Task, error) {
	query := r.db.WithContext(ctx).Model(&TaskDataModel{}).Order("name ASC")
	if filters != nil {
		if filters.IsActive != nil {
			query = query.Where("is_active = ?", *filters.IsActive)
		}
		if filters.Search != "" {
			query = query.Where("name ILIKE ?", "%"+filters.Search+"%")
		}
	}

	var data []TaskDataModel
	if err := query.Find(&data).Error; err != nil {
		return nil, err
	}

	results := make([]*domain.Task, 0, len(data))
	for _, row := range data {
		results = append(results, &domain.Task{
			ID:          row.ID,
			Name:        row.Name,
			Reference:   row.Reference,
			Description: row.Description,
			IsActive:    row.IsActive,
		})
	}
	return results, nil
}

func (r *GORMTaskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&TaskDataModel{}, "id = ?", id).Error
}

func (r *GORMTaskRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&TaskDataModel{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

type GORMPositionRepository struct {
	db *gorm.DB
}

func NewGORMPositionRepository(db *gorm.DB) *GORMPositionRepository {
	return &GORMPositionRepository{db: db}
}

func (r *GORMPositionRepository) Save(ctx context.Context, position *domain.Position) error {
	data := PositionDataModel{
		ID:          position.ID,
		Name:        position.Name,
		Code:        position.Code,
		Description: position.Description,
		IsActive:    position.IsActive,
	}
	return r.db.WithContext(ctx).Save(&data).Error
}

func (r *GORMPositionRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Position, error) {
	var data PositionDataModel
	err := r.db.WithContext(ctx).First(&data, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &domain.Position{
		ID:          data.ID,
		Name:        data.Name,
		Code:        data.Code,
		Description: data.Description,
		IsActive:    data.IsActive,
	}, nil
}

func (r *GORMPositionRepository) FindAll(ctx context.Context, filters *domain.PositionFilters) ([]*domain.Position, error) {
	query := r.db.WithContext(ctx).Model(&PositionDataModel{}).Order("name ASC")
	if filters != nil {
		if filters.IsActive != nil {
			query = query.Where("is_active = ?", *filters.IsActive)
		}
		if filters.Search != "" {
			query = query.Where("name ILIKE ? OR code ILIKE ?", "%"+filters.Search+"%", "%"+filters.Search+"%")
		}
	}

	var data []PositionDataModel
	if err := query.Find(&data).Error; err != nil {
		return nil, err
	}

	results := make([]*domain.Position, 0, len(data))
	for _, row := range data {
		results = append(results, &domain.Position{
			ID:          row.ID,
			Name:        row.Name,
			Code:        row.Code,
			Description: row.Description,
			IsActive:    row.IsActive,
		})
	}
	return results, nil
}

func (r *GORMPositionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&PositionDataModel{}, "id = ?", id).Error
}

type GORMWorkTypeRepository struct {
	db *gorm.DB
}

func NewGORMWorkTypeRepository(db *gorm.DB) *GORMWorkTypeRepository {
	return &GORMWorkTypeRepository{db: db}
}

func (r *GORMWorkTypeRepository) Save(ctx context.Context, workType *domain.WorkType) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		data := WorkTypeDataModel{
			ID:          workType.ID,
			Name:        workType.Name,
			Reference:   workType.Reference,
			Description: workType.Description,
			IsActive:    workType.IsActive,
		}
		if err := tx.Save(&data).Error; err != nil {
			return err
		}

		if err := tx.Where("work_type_id = ?", workType.ID).Delete(&WorkTypeTaskDataModel{}).Error; err != nil {
			return err
		}

		if len(workType.Tasks) > 0 {
			assignments := make([]WorkTypeTaskDataModel, 0, len(workType.Tasks))
			for _, assignment := range workType.Tasks {
				assignments = append(assignments, WorkTypeTaskDataModel{
					WorkTypeID: workType.ID,
					TaskID:     assignment.TaskID,
					Sequence:   assignment.Sequence,
				})
			}
			if err := tx.Create(&assignments).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *GORMWorkTypeRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.WorkType, error) {
	var data WorkTypeDataModel
	err := r.db.WithContext(ctx).
		Preload("Tasks", func(db *gorm.DB) *gorm.DB { return db.Order("sequence ASC") }).
		First(&data, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return mapWorkTypeToDomain(data), nil
}

func (r *GORMWorkTypeRepository) FindAll(ctx context.Context, filters *domain.WorkTypeFilters) ([]*domain.WorkType, error) {
	query := r.db.WithContext(ctx).
		Model(&WorkTypeDataModel{}).
		Preload("Tasks", func(db *gorm.DB) *gorm.DB { return db.Order("sequence ASC") }).
		Order("name ASC")

	if filters != nil {
		if filters.IsActive != nil {
			query = query.Where("is_active = ?", *filters.IsActive)
		}
		if filters.Search != "" {
			query = query.Where("name ILIKE ?", "%"+filters.Search+"%")
		}
	}

	var data []WorkTypeDataModel
	if err := query.Find(&data).Error; err != nil {
		return nil, err
	}

	results := make([]*domain.WorkType, 0, len(data))
	for _, row := range data {
		results = append(results, mapWorkTypeToDomain(row))
	}
	return results, nil
}

func (r *GORMWorkTypeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&WorkTypeDataModel{}, "id = ?", id).Error
}

func mapWorkTypeToDomain(data WorkTypeDataModel) *domain.WorkType {
	tasks := make([]domain.WorkTypeTask, 0, len(data.Tasks))
	for _, assignment := range data.Tasks {
		tasks = append(tasks, domain.WorkTypeTask{
			TaskID:   assignment.TaskID,
			Sequence: assignment.Sequence,
		})
	}

	return &domain.WorkType{
		ID:          data.ID,
		Name:        data.Name,
		Reference:   data.Reference,
		Description: data.Description,
		IsActive:    data.IsActive,
		Tasks:       tasks,
	}
}

type GORMWorkOrderRepository struct {
	db *gorm.DB
}

func NewGORMWorkOrderRepository(db *gorm.DB) *GORMWorkOrderRepository {
	return &GORMWorkOrderRepository{db: db}
}

func (r *GORMWorkOrderRepository) Save(ctx context.Context, work *domain.WorkOrder) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		data := WorkOrderDataModel{
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
		}
		if err := tx.Save(&data).Error; err != nil {
			return err
		}

		if err := tx.Exec(`
			DELETE FROM work_order_tasks
			WHERE work_order_line_id IN (
				SELECT id FROM work_order_lines WHERE work_order_id = ?
			)
		`, work.ID).Error; err != nil {
			return err
		}

		if err := tx.Where("work_order_id = ?", work.ID).Delete(&WorkOrderLineDataModel{}).Error; err != nil {
			return err
		}

		for _, line := range work.Lines {
			lineData := WorkOrderLineDataModel{
				ID:             line.ID,
				WorkOrderID:    work.ID,
				WorkTypeID:     line.WorkTypeID,
				PositionID:     line.PositionID,
				DesignFilePath: line.DesignFilePath,
				Notes:          line.Notes,
				Sequence:       line.Sequence,
			}
			if err := tx.Create(&lineData).Error; err != nil {
				return err
			}

			if len(line.Tasks) > 0 {
				taskRows := make([]WorkOrderTaskDataModel, 0, len(line.Tasks))
				for _, task := range line.Tasks {
					taskRows = append(taskRows, WorkOrderTaskDataModel{
						ID:              task.ID,
						WorkOrderLineID: line.ID,
						TaskID:          task.TaskID,
						Sequence:        task.Sequence,
						Status:          string(task.Status),
						AssignedTo:      task.AssignedTo,
						StartedAt:       task.StartedAt,
						CompletedAt:     task.CompletedAt,
						Notes:           task.Notes,
					})
				}
				if err := tx.Create(&taskRows).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (r *GORMWorkOrderRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.WorkOrder, error) {
	var data WorkOrderDataModel
	err := r.db.WithContext(ctx).
		Preload("Lines", func(db *gorm.DB) *gorm.DB { return db.Order("sequence ASC") }).
		Preload("Lines.Tasks", func(db *gorm.DB) *gorm.DB { return db.Order("sequence ASC") }).
		First(&data, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return mapWorkOrderToDomain(data)
}

func (r *GORMWorkOrderRepository) FindAll(ctx context.Context, filters *domain.WorkOrderFilters) ([]*domain.WorkOrder, error) {
	query := r.db.WithContext(ctx).Model(&WorkOrderDataModel{}).Order("work_number DESC")
	if filters != nil {
		if filters.Status != nil {
			query = query.Where("status = ?", string(*filters.Status))
		}
		if filters.Search != "" {
			searchTerm := "%" + filters.Search + "%"
			query = query.Where(`
				work_name ILIKE ?
				OR work_number ILIKE ?
				OR EXISTS (
					SELECT 1
					FROM organization_profiles op
					WHERE op.party_id = work_orders.party_id
					  AND op.name ILIKE ?
				)
				OR EXISTS (
					SELECT 1
					FROM person_profiles pp
					WHERE pp.party_id = work_orders.party_id
					  AND (
						pp.first_name ILIKE ?
						OR pp.last_name ILIKE ?
						OR (pp.first_name || ' ' || pp.last_name) ILIKE ?
					  )
				)
				OR EXISTS (
					SELECT 1
					FROM party_roles pr
					WHERE pr.party_id = work_orders.party_id
					  AND pr.creation_identifier ILIKE ?
				)
			`,
				searchTerm,
				searchTerm,
				searchTerm,
				searchTerm,
				searchTerm,
				searchTerm,
				searchTerm,
			)
		}
		if filters.PartyID != "" {
			query = query.Where("party_id = ?", filters.PartyID)
		}
		if filters.WorkSetupID != nil {
			query = query.Where("work_setup_id = ?", *filters.WorkSetupID)
		}
	}

	var rows []WorkOrderDataModel
	if err := query.
		Preload("Lines", func(db *gorm.DB) *gorm.DB { return db.Order("sequence ASC") }).
		Preload("Lines.Tasks", func(db *gorm.DB) *gorm.DB { return db.Order("sequence ASC") }).
		Find(&rows).Error; err != nil {
		return nil, err
	}

	result := make([]*domain.WorkOrder, 0, len(rows))
	for _, row := range rows {
		wo, err := mapWorkOrderToDomain(row)
		if err != nil {
			return nil, err
		}
		result = append(result, wo)
	}
	return result, nil
}

func (r *GORMWorkOrderRepository) CountByYear(ctx context.Context, year int) (int64, error) {
	var count int64
	prefix := fmt.Sprintf("MES-%d-", year)
	err := r.db.WithContext(ctx).Model(&WorkOrderDataModel{}).Where("work_number LIKE ?", prefix+"%").Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func mapWorkOrderToDomain(data WorkOrderDataModel) (*domain.WorkOrder, error) {
	status := domain.ProductionStatus(data.Status)
	priority := domain.WorkPriority(data.Priority)
	lines := make([]domain.WorkOrderLine, 0, len(data.Lines))

	for _, line := range data.Lines {
		tasks := make([]domain.WorkOrderTask, 0, len(line.Tasks))
		for _, task := range line.Tasks {
			tasks = append(tasks, domain.WorkOrderTask{
				ID:          task.ID,
				TaskID:      task.TaskID,
				Sequence:    task.Sequence,
				Status:      domain.TaskStatus(task.Status),
				AssignedTo:  task.AssignedTo,
				StartedAt:   task.StartedAt,
				CompletedAt: task.CompletedAt,
				Notes:       task.Notes,
			})
		}

		lines = append(lines, domain.WorkOrderLine{
			ID:             line.ID,
			WorkTypeID:     line.WorkTypeID,
			PositionID:     line.PositionID,
			DesignFilePath: line.DesignFilePath,
			Notes:          line.Notes,
			Sequence:       line.Sequence,
			Tasks:          tasks,
		})
	}

	return &domain.WorkOrder{
		ID:            data.ID,
		OrderNumber:   data.WorkNumber,
		OrderName:     data.WorkName,
		PartyID:       data.PartyID,
		WorkSetupID:   data.WorkSetupID,
		Notes:         data.Notes,
		Status:        status,
		Priority:      priority,
		StartDate:     data.StartDate,
		DueDate:       data.DueDate,
		CompletedDate: data.CompletedDate,
		Lines:         lines,
	}, nil
}

// --- WorkSetup Persistence ---

type WorkSetupDataModel struct {
	ID              uuid.UUID                `gorm:"column:id;type:uuid;primaryKey"`
	Name            string                   `gorm:"column:name"`
	PartyID         string                   `gorm:"column:party_id"`
	TangibleGroupID *uuid.UUID               `gorm:"column:tangible_group_id;type:uuid"`
	Description     string                   `gorm:"column:description"`
	IsActive        bool                     `gorm:"column:is_active"`
	Lines           []WorkSetupLineDataModel `gorm:"foreignKey:WorkSetupID;references:ID"`
}

func (WorkSetupDataModel) TableName() string { return "work_setups" }

type WorkSetupLineDataModel struct {
	ID             uuid.UUID `gorm:"column:id;type:uuid;primaryKey"`
	WorkSetupID    uuid.UUID `gorm:"column:work_setup_id;type:uuid"`
	WorkTypeID     uuid.UUID `gorm:"column:work_type_id;type:uuid"`
	PositionID     uuid.UUID `gorm:"column:position_id;type:uuid"`
	DesignFilePath string    `gorm:"column:design_file_path"`
	Notes          string    `gorm:"column:notes"`
	Sequence       int       `gorm:"column:sequence"`
}

func (WorkSetupLineDataModel) TableName() string { return "work_setup_lines" }

type GORMWorkSetupRepository struct {
	db *gorm.DB
}

func NewGORMWorkSetupRepository(db *gorm.DB) *GORMWorkSetupRepository {
	return &GORMWorkSetupRepository{db: db}
}

func (r *GORMWorkSetupRepository) Save(ctx context.Context, ws *domain.WorkSetup) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		data := WorkSetupDataModel{
			ID:              ws.ID,
			Name:            ws.Name,
			PartyID:         ws.PartyID,
			TangibleGroupID: ws.TangibleGroupID,
			Description:     ws.Description,
			IsActive:        ws.IsActive,
		}
		if err := tx.Save(&data).Error; err != nil {
			return err
		}

		if err := tx.Where("work_setup_id = ?", ws.ID).Delete(&WorkSetupLineDataModel{}).Error; err != nil {
			return err
		}

		if len(ws.Lines) > 0 {
			lines := make([]WorkSetupLineDataModel, 0, len(ws.Lines))
			for _, line := range ws.Lines {
				lineID := line.ID
				if lineID == uuid.Nil {
					lineID = uuid.New()
				}
				lines = append(lines, WorkSetupLineDataModel{
					ID:             lineID,
					WorkSetupID:    ws.ID,
					WorkTypeID:     line.WorkTypeID,
					PositionID:     line.PositionID,
					DesignFilePath: line.DesignFilePath,
					Notes:          line.Notes,
					Sequence:       line.Sequence,
				})
			}
			if err := tx.Create(&lines).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *GORMWorkSetupRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.WorkSetup, error) {
	var data WorkSetupDataModel
	err := r.db.WithContext(ctx).
		Preload("Lines", func(db *gorm.DB) *gorm.DB { return db.Order("sequence ASC") }).
		First(&data, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return mapWorkSetupToDomain(data), nil
}

func (r *GORMWorkSetupRepository) FindAll(ctx context.Context, filters *domain.WorkSetupFilters) ([]*domain.WorkSetup, error) {
	query := r.db.WithContext(ctx).
		Model(&WorkSetupDataModel{}).
		Preload("Lines", func(db *gorm.DB) *gorm.DB { return db.Order("sequence ASC") }).
		Order("name ASC")

	if filters != nil {
		if filters.IsActive != nil {
			query = query.Where("is_active = ?", *filters.IsActive)
		}
		if filters.Search != "" {
			query = query.Where("name ILIKE ?", "%"+filters.Search+"%")
		}
		if filters.PartyID != "" {
			query = query.Where("party_id = ?", filters.PartyID)
		}
	}

	var data []WorkSetupDataModel
	if err := query.Find(&data).Error; err != nil {
		return nil, err
	}

	results := make([]*domain.WorkSetup, 0, len(data))
	for _, row := range data {
		results = append(results, mapWorkSetupToDomain(row))
	}
	return results, nil
}

func (r *GORMWorkSetupRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&WorkSetupDataModel{}, "id = ?", id).Error
}

func mapWorkSetupToDomain(data WorkSetupDataModel) *domain.WorkSetup {
	lines := make([]domain.WorkSetupLine, 0, len(data.Lines))
	for _, l := range data.Lines {
		lines = append(lines, domain.WorkSetupLine{
			ID:             l.ID,
			WorkTypeID:     l.WorkTypeID,
			PositionID:     l.PositionID,
			DesignFilePath: l.DesignFilePath,
			Notes:          l.Notes,
			Sequence:       l.Sequence,
		})
	}

	return &domain.WorkSetup{
		ID:              data.ID,
		Name:            data.Name,
		PartyID:         data.PartyID,
		TangibleGroupID: data.TangibleGroupID,
		Description:     data.Description,
		IsActive:        data.IsActive,
		Lines:           lines,
	}
}
