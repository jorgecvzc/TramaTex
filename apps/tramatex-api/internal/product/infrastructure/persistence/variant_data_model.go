package persistence

import (
	"time"

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
	Status          string         `gorm:"type:variant_status;not null"`
	AttributeValues pq.StringArray `gorm:"type:uuid[]"`
	IsActive        bool           `gorm:"not null;default:true"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (VariantDataModel) TableName() string {
	return "product_variants"
}

// ToDomain converts the variant data model to a domain model.
func (v *VariantDataModel) ToDomain() *domain.ProductVariant {
	return &domain.ProductVariant{
		ID:              v.ID,
		ProductID:       v.ProductID,
		SKU:             v.SKU,
		Barcode:         v.Barcode,
		Status:          domain.VariantStatus(v.Status),
		AttributeValues: uuidArrayFromStringArray(v.AttributeValues),
		IsActive:        v.IsActive,
		CreatedAt:       v.CreatedAt,
		UpdatedAt:       v.UpdatedAt,
	}
}

// VariantFromDomain converts a variant domain model to a data model.
func VariantFromDomain(v *domain.ProductVariant) *VariantDataModel {
	return &VariantDataModel{
		ID:              v.ID,
		ProductID:       v.ProductID,
		SKU:             v.SKU,
		Barcode:         v.Barcode,
		Status:          string(v.Status),
		AttributeValues: stringArrayFromUUIDArray(v.AttributeValues),
		IsActive:        v.IsActive,
		CreatedAt:       v.CreatedAt,
		UpdatedAt:       v.UpdatedAt,
	}
}
