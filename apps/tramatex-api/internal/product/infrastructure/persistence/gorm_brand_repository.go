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
