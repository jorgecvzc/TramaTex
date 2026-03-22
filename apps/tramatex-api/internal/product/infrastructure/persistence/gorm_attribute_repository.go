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

	// Check if record exists to determine if it's an insert or update
	var existing AttributeDataModel
	result := r.db.WithContext(ctx).Select("id", "created_at").First(&existing, "id = ?", dataModel.ID)

	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}

	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// Existing record, preserve original CreatedAt
		dataModel.CreatedAt = existing.CreatedAt
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

		// Check if AttributeValue record exists to determine if it's an insert or update
		var existingValue AttributeValueDataModel
		valueResult := r.db.WithContext(ctx).Select("id", "created_at").First(&existingValue, "id = ?", values[i].ID)

		if valueResult.Error != nil && !errors.Is(valueResult.Error, gorm.ErrRecordNotFound) {
			return valueResult.Error
		}

		if !errors.Is(valueResult.Error, gorm.ErrRecordNotFound) {
			values[i].CreatedAt = existingValue.CreatedAt
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

// FindByCode finds an attribute by its code
func (r *GORMAttributeRepository) FindByCode(ctx context.Context, code string) (*domain.Attribute, error) {
	var dataModel AttributeDataModel
	err := r.db.WithContext(ctx).Preload("Values").First(&dataModel, "code = ?", code).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return dataModel.ToDomain(), nil
}

// FindByIDs finds attributes by their IDs, preserving the order of the input IDs slice.
// This is critical for SKU generation where attribute order comes from product.DirectAttributeIDs.
func (r *GORMAttributeRepository) FindByIDs(ctx context.Context, ids []uuid.UUID) ([]domain.Attribute, error) {
	var dataModels []AttributeDataModel
	err := r.db.WithContext(ctx).Preload("Values").Find(&dataModels, "id IN ?", ids).Error
	if err != nil {
		return nil, err
	}

	// Build lookup map for reordering by input ID order
	dmByID := make(map[uuid.UUID]*AttributeDataModel, len(dataModels))
	for i := range dataModels {
		dmByID[dataModels[i].ID] = &dataModels[i]
	}

	attributes := make([]domain.Attribute, 0, len(ids))
	for _, id := range ids {
		if dm, ok := dmByID[id]; ok {
			attributes = append(attributes, *dm.ToDomain())
		}
	}

	return attributes, nil
}

// FindByScope returns all attributes.
// Note: Scope filtering removed for MVP simplicity. This method now returns all attributes
// regardless of brandID/groupID parameters (kept for backwards compatibility).
func (r *GORMAttributeRepository) FindByScope(ctx context.Context, brandID *uuid.UUID, groupID *uuid.UUID) ([]*domain.Attribute, error) {
	var dataModels []AttributeDataModel
	err := r.db.WithContext(ctx).Preload("Values").Find(&dataModels).Error
	if err != nil {
		return nil, err
	}

	attributes := make([]*domain.Attribute, len(dataModels))
	for i, dm := range dataModels {
		attributes[i] = dm.ToDomain()
	}

	return attributes, nil
}

// Delete deletes an attribute and its values by its ID
func (r *GORMAttributeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// First delete all associated attribute values
	if err := r.db.WithContext(ctx).Where("attribute_id = ?", id).Delete(&AttributeValueDataModel{}).Error; err != nil {
		return err
	}
	// Then delete the attribute itself
	return r.db.WithContext(ctx).Delete(&AttributeDataModel{}, "id = ?", id).Error
}
