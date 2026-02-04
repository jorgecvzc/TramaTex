package domain

import "fmt"

// ProductType represents the type of a product, either tangible or a service.
type ProductType string

const (
	ProductTypeTangible ProductType = "TANGIBLE"
	ProductTypeService  ProductType = "SERVICE"
)

// IsValid checks if the ProductType is one of the predefined constants.
func (pt ProductType) IsValid() error {
	switch pt {
	case ProductTypeTangible, ProductTypeService:
		return nil
	}
	return fmt.Errorf("invalid ProductType: %s", pt)
}

// VariantStatus represents the state of a product variant.
type VariantStatus string

const (
	StatusProvisional VariantStatus = "PROVISIONAL"
	StatusConfirmed   VariantStatus = "CONFIRMED"
)

// IsValid checks if the VariantStatus is one of the predefined constants.
func (vs VariantStatus) IsValid() error {
	switch vs {
	case StatusProvisional, StatusConfirmed:
		return nil
	}
	return fmt.Errorf("invalid VariantStatus: %s", vs)
}
