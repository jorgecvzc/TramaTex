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

type ServiceGroupDataModel struct {
	ID             uuid.UUID                   `gorm:"column:id;type:uuid;primaryKey"`
	Name           string                      `gorm:"column:name"`
	Description    string                      `gorm:"column:description"`
	ProductGroupID *uuid.UUID                  `gorm:"column:product_group_id;type:uuid"`
	IsActive       bool                        `gorm:"column:is_active"`
	Tasks          []ServiceGroupTaskDataModel `gorm:"foreignKey:ServiceGroupID;references:ID"`
}

func (ServiceGroupDataModel) TableName() string { return "service_groups" }

type ServiceGroupTaskDataModel struct {
	ServiceGroupID uuid.UUID `gorm:"column:service_group_id;type:uuid;primaryKey"`
	TaskID         uuid.UUID `gorm:"column:task_id;type:uuid;primaryKey"`
	Sequence       int       `gorm:"column:sequence"`
}

func (ServiceGroupTaskDataModel) TableName() string { return "service_group_tasks" }

type MESWorkDataModel struct {
	ID              uuid.UUID                      `gorm:"column:id;type:uuid;primaryKey"`
	WorkNumber      string                         `gorm:"column:work_number"`
	WorkName        string                         `gorm:"column:work_name"`
	PartyID         string                         `gorm:"column:party_id"`
	TangibleGroupID uuid.UUID                      `gorm:"column:tangible_group_id;type:uuid"`
	GarmentNotes    string                         `gorm:"column:garment_notes"`
	Status          string                         `gorm:"column:status"`
	Priority        string                         `gorm:"column:priority"`
	StartDate       *time.Time                     `gorm:"column:start_date"`
	DueDate         *time.Time                     `gorm:"column:due_date"`
	CompletedDate   *time.Time                     `gorm:"column:completed_date"`
	ServiceGroups   []MESWorkServiceGroupDataModel `gorm:"foreignKey:MESWorkID;references:ID"`
}

func (MESWorkDataModel) TableName() string { return "mes_works" }

type MESWorkServiceGroupDataModel struct {
	ID             uuid.UUID              `gorm:"column:id;type:uuid;primaryKey"`
	MESWorkID      uuid.UUID              `gorm:"column:mes_work_id;type:uuid"`
	ServiceGroupID uuid.UUID              `gorm:"column:service_group_id;type:uuid"`
	PositionID     uuid.UUID              `gorm:"column:position_id;type:uuid"`
	DesignFilePath string                 `gorm:"column:design_file_path"`
	Notes          string                 `gorm:"column:notes"`
	Sequence       int                    `gorm:"column:sequence"`
	Tasks          []MESWorkTaskDataModel `gorm:"foreignKey:MESWorkServiceGroupID;references:ID"`
}

func (MESWorkServiceGroupDataModel) TableName() string { return "mes_work_service_groups" }

type MESWorkTaskDataModel struct {
	ID                    uuid.UUID  `gorm:"column:id;type:uuid;primaryKey"`
	MESWorkServiceGroupID uuid.UUID  `gorm:"column:mes_work_service_group_id;type:uuid"`
	TaskID                uuid.UUID  `gorm:"column:task_id;type:uuid"`
	Sequence              int        `gorm:"column:sequence"`
	Status                string     `gorm:"column:status"`
	AssignedTo            *uuid.UUID `gorm:"column:assigned_to;type:uuid"`
	StartedAt             *time.Time `gorm:"column:started_at"`
	CompletedAt           *time.Time `gorm:"column:completed_at"`
	Notes                 string     `gorm:"column:notes"`
}

func (MESWorkTaskDataModel) TableName() string { return "mes_work_tasks" }

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

type GORMServiceGroupRepository struct {
	db *gorm.DB
}

func NewGORMServiceGroupRepository(db *gorm.DB) *GORMServiceGroupRepository {
	return &GORMServiceGroupRepository{db: db}
}

func (r *GORMServiceGroupRepository) Save(ctx context.Context, serviceGroup *domain.ServiceGroup) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		data := ServiceGroupDataModel{
			ID:             serviceGroup.ID,
			Name:           serviceGroup.Name,
			Description:    serviceGroup.Description,
			ProductGroupID: serviceGroup.ProductGroupID,
			IsActive:       serviceGroup.IsActive,
		}
		if err := tx.Save(&data).Error; err != nil {
			return err
		}

		if err := tx.Where("service_group_id = ?", serviceGroup.ID).Delete(&ServiceGroupTaskDataModel{}).Error; err != nil {
			return err
		}

		if len(serviceGroup.Tasks) > 0 {
			assignments := make([]ServiceGroupTaskDataModel, 0, len(serviceGroup.Tasks))
			for _, assignment := range serviceGroup.Tasks {
				assignments = append(assignments, ServiceGroupTaskDataModel{
					ServiceGroupID: serviceGroup.ID,
					TaskID:         assignment.TaskID,
					Sequence:       assignment.Sequence,
				})
			}
			if err := tx.Create(&assignments).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *GORMServiceGroupRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.ServiceGroup, error) {
	var data ServiceGroupDataModel
	err := r.db.WithContext(ctx).
		Preload("Tasks", func(db *gorm.DB) *gorm.DB { return db.Order("sequence ASC") }).
		First(&data, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return mapServiceGroupToDomain(data), nil
}

func (r *GORMServiceGroupRepository) FindAll(ctx context.Context, filters *domain.ServiceGroupFilters) ([]*domain.ServiceGroup, error) {
	query := r.db.WithContext(ctx).
		Model(&ServiceGroupDataModel{}).
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

	var data []ServiceGroupDataModel
	if err := query.Find(&data).Error; err != nil {
		return nil, err
	}

	results := make([]*domain.ServiceGroup, 0, len(data))
	for _, row := range data {
		results = append(results, mapServiceGroupToDomain(row))
	}
	return results, nil
}

func (r *GORMServiceGroupRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&ServiceGroupDataModel{}, "id = ?", id).Error
}

func mapServiceGroupToDomain(data ServiceGroupDataModel) *domain.ServiceGroup {
	tasks := make([]domain.ServiceGroupTask, 0, len(data.Tasks))
	for _, assignment := range data.Tasks {
		tasks = append(tasks, domain.ServiceGroupTask{
			TaskID:   assignment.TaskID,
			Sequence: assignment.Sequence,
		})
	}

	return &domain.ServiceGroup{
		ID:             data.ID,
		Name:           data.Name,
		Description:    data.Description,
		ProductGroupID: data.ProductGroupID,
		IsActive:       data.IsActive,
		Tasks:          tasks,
	}
}

type GORMMESWorkRepository struct {
	db *gorm.DB
}

func NewGORMMESWorkRepository(db *gorm.DB) *GORMMESWorkRepository {
	return &GORMMESWorkRepository{db: db}
}

func (r *GORMMESWorkRepository) Save(ctx context.Context, work *domain.MESWork) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		data := MESWorkDataModel{
			ID:              work.ID,
			WorkNumber:      work.WorkNumber,
			WorkName:        work.WorkName,
			PartyID:         work.PartyID,
			TangibleGroupID: work.TangibleGroupID,
			GarmentNotes:    work.GarmentNotes,
			Status:          string(work.Status),
			Priority:        string(work.Priority),
			StartDate:       work.StartDate,
			DueDate:         work.DueDate,
			CompletedDate:   work.CompletedDate,
		}
		if err := tx.Save(&data).Error; err != nil {
			return err
		}

		if err := tx.Exec(`
			DELETE FROM mes_work_tasks
			WHERE mes_work_service_group_id IN (
				SELECT id FROM mes_work_service_groups WHERE mes_work_id = ?
			)
		`, work.ID).Error; err != nil {
			return err
		}

		if err := tx.Where("mes_work_id = ?", work.ID).Delete(&MESWorkServiceGroupDataModel{}).Error; err != nil {
			return err
		}

		for _, group := range work.ServiceGroups {
			groupData := MESWorkServiceGroupDataModel{
				ID:             group.ID,
				MESWorkID:      work.ID,
				ServiceGroupID: group.ServiceGroupID,
				PositionID:     group.PositionID,
				DesignFilePath: group.DesignFilePath,
				Notes:          group.Notes,
				Sequence:       group.Sequence,
			}
			if err := tx.Create(&groupData).Error; err != nil {
				return err
			}

			if len(group.Tasks) > 0 {
				taskRows := make([]MESWorkTaskDataModel, 0, len(group.Tasks))
				for _, task := range group.Tasks {
					taskRows = append(taskRows, MESWorkTaskDataModel{
						ID:                    task.ID,
						MESWorkServiceGroupID: group.ID,
						TaskID:                task.TaskID,
						Sequence:              task.Sequence,
						Status:                string(task.Status),
						AssignedTo:            task.AssignedTo,
						StartedAt:             task.StartedAt,
						CompletedAt:           task.CompletedAt,
						Notes:                 task.Notes,
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

func (r *GORMMESWorkRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.MESWork, error) {
	var data MESWorkDataModel
	err := r.db.WithContext(ctx).
		Preload("ServiceGroups", func(db *gorm.DB) *gorm.DB { return db.Order("sequence ASC") }).
		Preload("ServiceGroups.Tasks", func(db *gorm.DB) *gorm.DB { return db.Order("sequence ASC") }).
		First(&data, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return mapMESWorkToDomain(data)
}

func (r *GORMMESWorkRepository) FindAll(ctx context.Context, filters *domain.MESWorkFilters) ([]*domain.MESWork, error) {
	query := r.db.WithContext(ctx).Model(&MESWorkDataModel{}).Order("created_at DESC")
	if filters != nil {
		if filters.Status != nil {
			query = query.Where("status = ?", string(*filters.Status))
		}
		if filters.Search != "" {
			query = query.Where("work_name ILIKE ? OR work_number ILIKE ?", "%"+filters.Search+"%", "%"+filters.Search+"%")
		}
		if filters.PartyID != "" {
			query = query.Where("party_id = ?", filters.PartyID)
		}
	}

	var rows []MESWorkDataModel
	if err := query.Find(&rows).Error; err != nil {
		return nil, err
	}

	result := make([]*domain.MESWork, 0, len(rows))
	for _, row := range rows {
		status := domain.ProductionStatus(row.Status)
		priority := domain.WorkPriority(row.Priority)
		result = append(result, &domain.MESWork{
			ID:              row.ID,
			WorkNumber:      row.WorkNumber,
			WorkName:        row.WorkName,
			PartyID:         row.PartyID,
			TangibleGroupID: row.TangibleGroupID,
			GarmentNotes:    row.GarmentNotes,
			Status:          status,
			Priority:        priority,
			StartDate:       row.StartDate,
			DueDate:         row.DueDate,
			CompletedDate:   row.CompletedDate,
			ServiceGroups:   []domain.MESWorkServiceGroup{},
		})
	}
	return result, nil
}

func (r *GORMMESWorkRepository) CountByYear(ctx context.Context, year int) (int64, error) {
	var count int64
	prefix := fmt.Sprintf("MES-%d-", year)
	err := r.db.WithContext(ctx).Model(&MESWorkDataModel{}).Where("work_number LIKE ?", prefix+"%").Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func mapMESWorkToDomain(data MESWorkDataModel) (*domain.MESWork, error) {
	status := domain.ProductionStatus(data.Status)
	priority := domain.WorkPriority(data.Priority)
	groups := make([]domain.MESWorkServiceGroup, 0, len(data.ServiceGroups))

	for _, group := range data.ServiceGroups {
		tasks := make([]domain.MESWorkTask, 0, len(group.Tasks))
		for _, task := range group.Tasks {
			tasks = append(tasks, domain.MESWorkTask{
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

		groups = append(groups, domain.MESWorkServiceGroup{
			ID:             group.ID,
			ServiceGroupID: group.ServiceGroupID,
			PositionID:     group.PositionID,
			DesignFilePath: group.DesignFilePath,
			Notes:          group.Notes,
			Sequence:       group.Sequence,
			Tasks:          tasks,
		})
	}

	return &domain.MESWork{
		ID:              data.ID,
		WorkNumber:      data.WorkNumber,
		WorkName:        data.WorkName,
		PartyID:         data.PartyID,
		TangibleGroupID: data.TangibleGroupID,
		GarmentNotes:    data.GarmentNotes,
		Status:          status,
		Priority:        priority,
		StartDate:       data.StartDate,
		DueDate:         data.DueDate,
		CompletedDate:   data.CompletedDate,
		ServiceGroups:   groups,
	}, nil
}
