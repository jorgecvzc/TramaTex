package persistence

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/joran-cortez/tramatex/internal/product/domain"
)

// GORMProductRepository is a GORM implementation of the ProductRepository
type GORMProductRepository struct {
	db *gorm.DB
}

// NewGORMProductRepository creates a new GORMProductRepository
func NewGORMProductRepository(db *gorm.DB) *GORMProductRepository {
	return &GORMProductRepository{db: db}
}

// Save saves a product to the database
func (r *GORMProductRepository) Save(ctx context.Context, product *domain.Product) error {
	dataModel := FromDomain(product)
	return r.db.WithContext(ctx).Create(dataModel).Error
}

// FindByID finds a product by its ID
func (r *GORMProductRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	var dataModel ProductDataModel
	err := r.db.WithContext(ctx).First(&dataModel, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Or a custom not found error
		}
		return nil, err
	}
	return dataModel.ToDomain(), nil
}

// FindBySKU finds a product by its SKU
func (r *GORMProductRepository) FindBySKU(ctx context.Context, sku string) (*domain.Product, error) {
	var dataModel ProductDataModel
	err := r.db.WithContext(ctx).First(&dataModel, "sku = ?", sku).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Or a custom not found error
		}
		return nil, err
	}
	return dataModel.ToDomain(), nil
}

// UpdateSKUs updates the SKU of a product.
func (r *GORMProductRepository) UpdateSKUs(ctx context.Context, productID uuid.UUID, newSKU string) error {
	return r.db.WithContext(ctx).Model(&ProductDataModel{}).Where("id = ?", productID).Update("sku", newSKU).Error
}
