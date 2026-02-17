package persistence

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/joran-cortez/tramatex/internal/product/domain"
)

// GORMBrandRepository is a GORM implementation of the BrandRepository
type GORMBrandRepository struct {
	db *gorm.DB
}

// NewGORMBrandRepository creates a new GORMBrandRepository
func NewGORMBrandRepository(db *gorm.DB) *GORMBrandRepository {
	return &GORMBrandRepository{db: db}
}

// FindByID finds a brand by its ID
func (r *GORMBrandRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Brand, error) {
	var dataModel BrandDataModel
	err := r.db.WithContext(ctx).First(&dataModel, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Or a custom not found error
		}
		return nil, err
	}
	return dataModel.ToDomain(), nil
}

// FindAll finds all brands
func (r *GORMBrandRepository) FindAll(ctx context.Context) ([]*domain.Brand, error) {
	var dataModels []BrandDataModel
	err := r.db.WithContext(ctx).Order("name ASC").Find(&dataModels).Error
	if err != nil {
		return nil, err
	}

	brands := make([]*domain.Brand, len(dataModels))
	for i, dm := range dataModels {
		brands[i] = dm.ToDomain()
	}
	return brands, nil
}

// Save saves or updates a brand
func (r *GORMBrandRepository) Save(ctx context.Context, brand *domain.Brand) error {
	dataModel := BrandFromDomain(brand)
	return r.db.WithContext(ctx).Save(dataModel).Error
}

// Delete deletes a brand by its ID
func (r *GORMBrandRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&BrandDataModel{}, "id = ?", id).Error
}
