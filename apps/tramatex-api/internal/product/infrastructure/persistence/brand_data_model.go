package persistence

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/joran-cortez/tramatex/internal/product/domain"
)

// BrandDataModel represents the brand entity in the database.
type BrandDataModel struct {
	gorm.Model
	ID                      uuid.UUID `gorm:"type:uuid;primary_key;"`
	Name                    string    `gorm:"uniqueIndex;not null"`
	DefaultMarkupPercentage float64   `gorm:"type:numeric(5,2);not null;default:0"`
	IsActive                bool      `gorm:"not null;default:true"`
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

func (BrandDataModel) TableName() string {
	return "brands"
}

// ToDomain converts the brand data model to a domain model.
func (b *BrandDataModel) ToDomain() *domain.Brand {
	return &domain.Brand{
		ID:                      b.ID,
		Name:                    b.Name,
		DefaultMarkupPercentage: b.DefaultMarkupPercentage,
		IsActive:                b.IsActive,
	}
}

// BrandFromDomain converts a brand domain model to a data model.
func BrandFromDomain(b *domain.Brand) *BrandDataModel {
	return &BrandDataModel{
		ID:                      b.ID,
		Name:                    b.Name,
		DefaultMarkupPercentage: b.DefaultMarkupPercentage,
		IsActive:                b.IsActive,
	}
}
