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
	return r.db.WithContext(ctx).Create(dataModel).Error
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
