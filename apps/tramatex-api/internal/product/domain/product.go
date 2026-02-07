package domain

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

// Brand represents an external aggregate root, defined here for context.
type Brand struct {
	ID   uuid.UUID
	Name string
}

// ID_PTR returns a pointer to the Brand's ID.
func (b *Brand) ID_PTR() *uuid.UUID {
	return &b.ID
}

// ProductGroup represents an external aggregate root, defined here for context.
type ProductGroup struct {
	ID   uuid.UUID
	Name string
}

// ID_PTR returns a pointer to the ProductGroup's ID.
func (pg *ProductGroup) ID_PTR() *uuid.UUID {
	return &pg.ID
}

// Product is the aggregate root for the product concept.
type Product struct {
	ID                 uuid.UUID
	SKU                string
	Name               string
	LongName           string
	Barcode            *string // Optional
	Description        string
	ProductType        ProductType
	BrandID            uuid.UUID
	GroupIDs           []uuid.UUID
	DirectAttributeIDs []uuid.UUID
	IsActive           bool
}

// NewProduct creates a new Product with validation.
func NewProduct(
	sku, name, longName, description string,
	productType ProductType,
	brandID uuid.UUID,
	barcode *string,
) (*Product, error) {
	if sku == "" {
		return nil, fmt.Errorf("product SKU cannot be empty")
	}
	if name == "" {
		return nil, fmt.Errorf("product name cannot be empty")
	}
	if err := productType.IsValid(); err != nil {
		return nil, err
	}
	if brandID == uuid.Nil {
		return nil, fmt.Errorf("product must be associated with a brand")
	}

	return &Product{
		ID:                 uuid.New(),
		SKU:                sku,
		Name:               name,
		LongName:           longName,
		Barcode:            barcode,
		Description:        description,
		ProductType:        productType,
		BrandID:            brandID,
		GroupIDs:           make([]uuid.UUID, 0),
		DirectAttributeIDs: make([]uuid.UUID, 0),
		IsActive:           true,
	}, nil
}

// AddGroup adds a product group to the product.
func (p *Product) AddGroup(groupID uuid.UUID) {
	// Avoid duplicates
	for _, id := range p.GroupIDs {
		if id == groupID {
			return
		}
	}
	p.GroupIDs = append(p.GroupIDs, groupID)
}

// AddDirectAttribute adds a direct attribute to the product for overriding inheritance.
func (p *Product) AddDirectAttribute(attributeID uuid.UUID) {
	// Avoid duplicates
	for _, id := range p.DirectAttributeIDs {
		if id == attributeID {
			return
		}
	}
	p.DirectAttributeIDs = append(p.DirectAttributeIDs, attributeID)
}
