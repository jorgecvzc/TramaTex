package persistence

import (
	"context"
	"errors"
	"time"

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

	// Check if record exists to determine if it's an insert or update
	var existing ProductDataModel
	result := r.db.WithContext(ctx).Select("id", "created_at").First(&existing, "id = ?", dataModel.ID)

	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// New record, GORM's Model hook will handle CreatedAt
	} else {
		// Existing record, preserve original CreatedAt
		dataModel.CreatedAt = existing.CreatedAt
	}

	return r.db.WithContext(ctx).Save(dataModel).Error
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

// FindBySKU finds a product by its SKU (case-insensitive)
func (r *GORMProductRepository) FindBySKU(ctx context.Context, sku string) (*domain.Product, error) {
	var dataModel ProductDataModel
	err := r.db.WithContext(ctx).First(&dataModel, "LOWER(sku) = LOWER(?)", sku).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Or a custom not found error
		}
		return nil, err
	}
	return dataModel.ToDomain(), nil
}

// FindByBarcode finds a product by its barcode
func (r *GORMProductRepository) FindByBarcode(ctx context.Context, barcode string) (*domain.Product, error) {
	var dataModel ProductDataModel
	err := r.db.WithContext(ctx).First(&dataModel, "barcode = ?", barcode).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return dataModel.ToDomain(), nil
}

// FindBySKUPrefix finds products whose SKU starts with the given prefix
func (r *GORMProductRepository) FindBySKUPrefix(ctx context.Context, prefix string) ([]*domain.Product, error) {
	var dataModels []ProductDataModel
	err := r.db.WithContext(ctx).Where("sku ILIKE ?", prefix+"%").Order("sku asc").Limit(20).Find(&dataModels).Error
	if err != nil {
		return nil, err
	}
	products := make([]*domain.Product, len(dataModels))
	for i := range dataModels {
		products[i] = dataModels[i].ToDomain()
	}
	return products, nil
}

// FindAll lists all products.
func (r *GORMProductRepository) FindAll(ctx context.Context) ([]*domain.Product, error) {
	var dataModels []ProductDataModel
	err := r.db.WithContext(ctx).Order("created_at desc").Find(&dataModels).Error
	if err != nil {
		return nil, err
	}
	products := make([]*domain.Product, len(dataModels))
	for i := range dataModels {
		products[i] = dataModels[i].ToDomain()
	}
	return products, nil
}

// UpdateSKUs updates the SKU of a product.

func (r *GORMProductRepository) UpdateSKUs(ctx context.Context, productID uuid.UUID, newSKU string) error {
	return r.db.WithContext(ctx).Model(&ProductDataModel{}).Where("id = ?", productID).Updates(map[string]interface{}{
		"sku":        newSKU,
		"updated_at": time.Now(), // Explicitly update updated_at since we're using Updates with a map
	}).Error
}
