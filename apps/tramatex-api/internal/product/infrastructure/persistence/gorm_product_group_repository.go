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
