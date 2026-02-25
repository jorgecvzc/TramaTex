package persistence

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"

	"github.com/joran-cortez/tramatex/internal/product/domain"
)

// VariantDataModel represents the product variant entity in the database.
type VariantDataModel struct {
	gorm.Model
	ID              uuid.UUID      `gorm:"type:uuid;primary_key;"`
	ProductID       uuid.UUID      `gorm:"not null"`
	SKU             string         `gorm:"uniqueIndex;not null"`
	Barcode         *string        `gorm:"uniqueIndex"`
	BaseCost        float64        `gorm:"type:numeric(12,2);not null;default:0"`
	Status          string         `gorm:"type:variant_status;not null"`
	AttributeValues pq.StringArray `gorm:"type:uuid[]"`
	IsActive        bool           `gorm:"not null;default:true"`
	CreatedBy       string
	ModifiedBy      string
}

func (VariantDataModel) TableName() string {
	return "product_variants"
}

// ToDomain converts the variant data model to a domain model.
// Note: BaseCost is ignored as it's now calculated dynamically from Product.BasePrice + Attribute modifiers.
func (v *VariantDataModel) ToDomain() *domain.ProductVariant {
	return &domain.ProductVariant{
		ID:              v.ID,
		ProductID:       v.ProductID,
		SKU:             v.SKU,
		Barcode:         v.Barcode,
		Status:          domain.VariantStatus(v.Status),
		AttributeValues: uuidArrayFromStringArray(v.AttributeValues),
		IsActive:        v.IsActive,
	}
}

// VariantFromDomain converts a variant domain model to a data model.
// Note: BaseCost is set to 0 as it's now calculated dynamically, not stored.
func VariantFromDomain(v *domain.ProductVariant) *VariantDataModel {
	return &VariantDataModel{
		ID:              v.ID,
		ProductID:       v.ProductID,
		SKU:             v.SKU,
		Barcode:         v.Barcode,
		BaseCost:        0, // Legacy field, no longer used
		Status:          string(v.Status),
		AttributeValues: stringArrayFromUUIDArray(v.AttributeValues),
		IsActive:        v.IsActive,
	}
}
