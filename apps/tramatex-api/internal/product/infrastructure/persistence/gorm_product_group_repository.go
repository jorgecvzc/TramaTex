package persistence

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/joran-cortez/tramatex/internal/product/domain"
)

// GORMProductGroupRepository is a GORM implementation of the ProductGroupRepository
type GORMProductGroupRepository struct {
	db *gorm.DB
}

// NewGORMProductGroupRepository creates a new GORMProductGroupRepository
func NewGORMProductGroupRepository(db *gorm.DB) *GORMProductGroupRepository {
	return &GORMProductGroupRepository{db: db}
}

// FindByID finds a product group by its ID
func (r *GORMProductGroupRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.ProductGroup, error) {
	var dataModel ProductGroupDataModel
	err := r.db.WithContext(ctx).First(&dataModel, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Or a custom not found error
		}
		return nil, err
	}
	return dataModel.ToDomain(), nil
}

// FindAll finds all product groups
func (r *GORMProductGroupRepository) FindAll(ctx context.Context) ([]*domain.ProductGroup, error) {
	var dataModels []ProductGroupDataModel
	err := r.db.WithContext(ctx).Order("name ASC").Find(&dataModels).Error
	if err != nil {
		return nil, err
	}

	groups := make([]*domain.ProductGroup, len(dataModels))
	for i, dm := range dataModels {
		groups[i] = dm.ToDomain()
	}
	return groups, nil
}

// Save saves or updates a product group
func (r *GORMProductGroupRepository) Save(ctx context.Context, group *domain.ProductGroup) error {
	dataModel := ProductGroupFromDomain(group)
	return r.db.WithContext(ctx).Save(dataModel).Error
}

// Delete deletes a product group by its ID
func (r *GORMProductGroupRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&ProductGroupDataModel{}, "id = ?", id).Error
}
