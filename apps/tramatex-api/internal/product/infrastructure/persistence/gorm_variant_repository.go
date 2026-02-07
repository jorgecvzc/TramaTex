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
	actorID := actorIDFromContext(ctx)

	// Check if record exists to determine if it's an insert or update
	var existing VariantDataModel
	result := r.db.WithContext(ctx).Select("id", "created_at", "created_by").First(&existing, "id = ?", dataModel.ID)

	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// New record, GORM's Model hook will handle CreatedAt
		dataModel.CreatedBy = actorID
	} else {
		// Existing record, preserve original CreatedAt and CreatedBy
		dataModel.CreatedAt = existing.CreatedAt
		dataModel.CreatedBy = existing.CreatedBy
		// Update ModifiedBy for existing record
		dataModel.ModifiedBy = actorID
	}

	return r.db.WithContext(ctx).Save(dataModel).Error
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

// FindByProductID lists product variants by product ID.
func (r *GORMVariantRepository) FindByProductID(ctx context.Context, productID uuid.UUID) ([]*domain.ProductVariant, error) {
	var dataModels []VariantDataModel
	err := r.db.WithContext(ctx).Where("product_id = ?", productID).Order("created_at desc").Find(&dataModels).Error
	if err != nil {
		return nil, err
	}
	variants := make([]*domain.ProductVariant, len(dataModels))
	for i := range dataModels {
		variants[i] = dataModels[i].ToDomain()
	}
	return variants, nil
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
