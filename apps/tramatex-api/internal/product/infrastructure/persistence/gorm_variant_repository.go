package persistence

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/joran-cortez/tramatex/internal/product/domain"
)

// GORMVariantRepository is a GORM implementation of the ProductVariantRepository
type GORMVariantRepository struct {
	db *gorm.DB
}

// NewGORMVariantRepository creates a new GORMVariantRepository
func NewGORMVariantRepository(db *gorm.DB) *GORMVariantRepository {
	return &GORMVariantRepository{db: db}
}

// Save saves a product variant to the database
func (r *GORMVariantRepository) Save(ctx context.Context, variant *domain.ProductVariant) error {
	dataModel := VariantFromDomain(variant)
	return r.db.WithContext(ctx).Create(dataModel).Error
}

// FindByID finds a product variant by its ID
func (r *GORMVariantRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.ProductVariant, error) {
	var dataModel VariantDataModel
	err := r.db.WithContext(ctx).First(&dataModel, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Or a custom not found error
		}
		return nil, err
	}
	return dataModel.ToDomain(), nil
}

// FindBySKU finds a product variant by its SKU
func (r *GORMVariantRepository) FindBySKU(ctx context.Context, sku string) (*domain.ProductVariant, error) {
	var dataModel VariantDataModel
	err := r.db.WithContext(ctx).First(&dataModel, "sku = ?", sku).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Or a custom not found error
		}
		return nil, err
	}
	return dataModel.ToDomain(), nil
}

// FindByProductIDAndAttributeValues finds a product variant by its product ID and attribute values
func (r *GORMVariantRepository) FindByProductIDAndAttributeValues(ctx context.Context, productID uuid.UUID, attributeValueIDs []uuid.UUID) (*domain.ProductVariant, error) {
	var dataModel VariantDataModel
	// This is a simplified search. A more robust implementation would handle the array comparison correctly.
	err := r.db.WithContext(ctx).Where("product_id = ? AND attribute_values @> ?", productID, stringArrayFromUUIDArray(attributeValueIDs)).First(&dataModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Or a custom not found error
		}
		return nil, err
	}
	return dataModel.ToDomain(), nil
}
