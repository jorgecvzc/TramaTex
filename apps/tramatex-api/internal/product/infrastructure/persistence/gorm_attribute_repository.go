package persistence

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/joran-cortez/tramatex/internal/product/domain"
)

// GORMAttributeRepository is a GORM implementation of the AttributeRepository
type GORMAttributeRepository struct {
	db *gorm.DB
}

// NewGORMAttributeRepository creates a new GORMAttributeRepository
func NewGORMAttributeRepository(db *gorm.DB) *GORMAttributeRepository {
	return &GORMAttributeRepository{db: db}
}

// Save saves an attribute to the database
func (r *GORMAttributeRepository) Save(ctx context.Context, attribute *domain.Attribute) error {
	dataModel := AttributeFromDomain(attribute)
	actorID := actorIDFromContext(ctx)
	dataModel.CreatedBy = actorID
	dataModel.ModifiedBy = actorID

	// Check if record exists to determine if it's an insert or update
	var existing AttributeDataModel
	result := r.db.WithContext(ctx).Select("id", "created_at", "created_by").First(&existing, "id = ?", dataModel.ID)

	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// New record, set CreatedAt (GORM's Model hook will handle this if not explicitly set)
		// and use the provided createdBy
		dataModel.CreatedBy = actorID
	} else {
		// Existing record, preserve original CreatedAt and CreatedBy
		dataModel.CreatedAt = existing.CreatedAt
		dataModel.CreatedBy = existing.CreatedBy
		dataModel.ModifiedBy = actorID
	}

	values := dataModel.Values
	dataModel.Values = nil // Clear values to avoid GORM attempting to save them directly with the attribute
	if err := r.db.WithContext(ctx).Save(dataModel).Error; err != nil {
		return err
	}
	if len(values) == 0 {
		return r.db.WithContext(ctx).Where("attribute_id = ?", dataModel.ID).Delete(&AttributeValueDataModel{}).Error
	}
	for i := range values {
		values[i].AttributeID = dataModel.ID
		values[i].CreatedBy = actorID // Default for new values
		values[i].ModifiedBy = actorID

		// Check if AttributeValue record exists to determine if it's an insert or update
		var existingValue AttributeValueDataModel
		valueResult := r.db.WithContext(ctx).Select("id", "created_at", "created_by").First(&existingValue, "id = ?", values[i].ID)

		if valueResult.Error != nil && !errors.Is(valueResult.Error, gorm.ErrRecordNotFound) {
			return valueResult.Error
		}

		if errors.Is(valueResult.Error, gorm.ErrRecordNotFound) {
			values[i].CreatedBy = actorID
		} else {
			values[i].CreatedAt = existingValue.CreatedAt
			values[i].CreatedBy = existingValue.CreatedBy
			values[i].ModifiedBy = actorID
		}
	}
	ids := make([]uuid.UUID, 0, len(values))
	for i := range values {
		ids = append(ids, values[i].ID)
	}
	if err := r.db.WithContext(ctx).Where("attribute_id = ? AND id NOT IN ?", dataModel.ID, ids).Delete(&AttributeValueDataModel{}).Error; err != nil {
		return err
	}
	return r.db.WithContext(ctx).Save(&values).Error
}

// FindByID finds an attribute by its ID
func (r *GORMAttributeRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Attribute, error) {
	var dataModel AttributeDataModel
	err := r.db.WithContext(ctx).Preload("Values").First(&dataModel, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Or a custom not found error
		}
		return nil, err
	}
	return dataModel.ToDomain(), nil
}

// FindByIDs finds attributes by their IDs
func (r *GORMAttributeRepository) FindByIDs(ctx context.Context, ids []uuid.UUID) ([]domain.Attribute, error) {
	var dataModels []AttributeDataModel
	err := r.db.WithContext(ctx).Preload("Values").Find(&dataModels, "id IN ?", ids).Error
	if err != nil {
		return nil, err
	}

	attributes := make([]domain.Attribute, len(dataModels))
	for i, dm := range dataModels {
		attributes[i] = *dm.ToDomain()
	}

	return attributes, nil
}

// FindByScope finds attributes based on optional brandID and groupID.
// If both are nil, it returns generic attributes.
func (r *GORMAttributeRepository) FindByScope(ctx context.Context, brandID *uuid.UUID, groupID *uuid.UUID) ([]*domain.Attribute, error) {
	var dataModels []AttributeDataModel
	query := r.db.WithContext(ctx).Preload("Values")

	if brandID != nil && groupID != nil {
		// Scoped to both brand and group
		query = query.Where("scope_brand_id = ? AND scope_group_id = ?", *brandID, *groupID)
	} else if brandID != nil {
		// Scoped to only brand
		query = query.Where("scope_brand_id = ? AND scope_group_id IS NULL", *brandID)
	} else if groupID != nil {
		// Scoped to only group
		query = query.Where("scope_group_id = ? AND scope_brand_id IS NULL", *groupID)
	} else {
		// Generic attributes (not scoped to any brand or group)
		query = query.Where("scope_brand_id IS NULL AND scope_group_id IS NULL")
	}

	err := query.Find(&dataModels).Error
	if err != nil {
		return nil, err
	}

	attributes := make([]*domain.Attribute, len(dataModels))
	for i, dm := range dataModels {
		attributes[i] = dm.ToDomain()
	}

	return attributes, nil
}
