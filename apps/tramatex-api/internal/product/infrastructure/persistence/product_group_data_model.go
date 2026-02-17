package persistence

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/joran-cortez/tramatex/internal/product/domain"
)

// ProductGroupDataModel represents the product group entity in the database.
type ProductGroupDataModel struct {
	gorm.Model
	ID            uuid.UUID `gorm:"type:uuid;primary_key;"`
	Name          string    `gorm:"uniqueIndex;not null"`
	ParentGroupID *uuid.UUID
	IsActive      bool `gorm:"not null;default:true"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (ProductGroupDataModel) TableName() string {
	return "product_groups"
}

// ToDomain converts the product group data model to a domain model.
func (pg *ProductGroupDataModel) ToDomain() *domain.ProductGroup {
	return &domain.ProductGroup{
		ID:            pg.ID,
		Name:          pg.Name,
		ParentGroupID: pg.ParentGroupID,
		IsActive:      pg.IsActive,
	}
}

// ProductGroupFromDomain converts a product group domain model to a data model.
func ProductGroupFromDomain(pg *domain.ProductGroup) *ProductGroupDataModel {
	return &ProductGroupDataModel{
		ID:            pg.ID,
		Name:          pg.Name,
		ParentGroupID: pg.ParentGroupID,
		IsActive:      pg.IsActive,
	}
}
