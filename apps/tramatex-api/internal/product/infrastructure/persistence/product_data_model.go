package persistence

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"

	"github.com/joran-cortez/tramatex/internal/product/domain"
)

// ProductDataModel represents the product entity in the database.
type ProductDataModel struct {
	gorm.Model
	ID                 uuid.UUID `gorm:"type:uuid;primary_key;"`
	SKU                string    `gorm:"uniqueIndex;not null"`
	Name               string    `gorm:"not null"`
	LongName           string
	Barcode            *string `gorm:"uniqueIndex"`
	Description        string
	ProductType        string         `gorm:"type:product_type;not null"`
	BrandID            *uuid.UUID     `gorm:"type:uuid"`
	GroupIDs           pq.StringArray `gorm:"type:uuid[]"`
	DirectAttributeIDs pq.StringArray `gorm:"type:uuid[]"`
	BasePrice          float64        `gorm:"type:numeric(12,2);not null"`
	TaxRate            float64        `gorm:"type:numeric(5,2);not null"`
	IsActive           bool           `gorm:"not null;default:true"`
}

func (ProductDataModel) TableName() string {
	return "products"
}

// ToDomain converts the data model to a domain model.
func (p *ProductDataModel) ToDomain() *domain.Product {
	return &domain.Product{
		ID:                 p.ID,
		SKU:                p.SKU,
		Name:               p.Name,
		LongName:           p.LongName,
		Barcode:            p.Barcode,
		Description:        p.Description,
		ProductType:        domain.ProductType(p.ProductType),
		BrandID:            p.BrandID,
		GroupIDs:           uuidArrayFromStringArray(p.GroupIDs),
		DirectAttributeIDs: uuidArrayFromStringArray(p.DirectAttributeIDs),
		BasePrice:          p.BasePrice,
		TaxRate:            p.TaxRate,
		IsActive:           p.IsActive,
	}
}

// FromDomain converts a domain model to a data model.
func FromDomain(p *domain.Product) *ProductDataModel {
	return &ProductDataModel{
		ID:                 p.ID,
		SKU:                p.SKU,
		Name:               p.Name,
		LongName:           p.LongName,
		Barcode:            p.Barcode,
		Description:        p.Description,
		ProductType:        string(p.ProductType),
		BrandID:            p.BrandID,
		GroupIDs:           stringArrayFromUUIDArray(p.GroupIDs),
		DirectAttributeIDs: stringArrayFromUUIDArray(p.DirectAttributeIDs),
		BasePrice:          p.BasePrice,
		TaxRate:            p.TaxRate,
		IsActive:           p.IsActive,
	}
}

// Helper functions to convert between []uuid.UUID and pq.StringArray
func uuidArrayFromStringArray(src pq.StringArray) []uuid.UUID {
	if src == nil {
		return nil
	}
	dest := make([]uuid.UUID, len(src))
	for i, s := range src {
		dest[i], _ = uuid.Parse(s)
	}
	return dest
}

func stringArrayFromUUIDArray(src []uuid.UUID) pq.StringArray {
	if src == nil {
		return nil
	}
	dest := make(pq.StringArray, len(src))
	for i, u := range src {
		dest[i] = u.String()
	}
	return dest
}
