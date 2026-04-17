package domain

import (
	"context"

	"github.com/google/uuid"
)

// BrandRepository defines the interface for interacting with Brand data.
type BrandRepository interface {
	Save(ctx context.Context, brand *Brand) error
	FindByID(ctx context.Context, id uuid.UUID) (*Brand, error)
	FindByIDs(ctx context.Context, ids []uuid.UUID) ([]*Brand, error)
	FindAll(ctx context.Context) ([]*Brand, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// ProductGroupRepository defines the interface for interacting with ProductGroup data.
type ProductGroupRepository interface {
	Save(ctx context.Context, group *ProductGroup) error
	FindByID(ctx context.Context, id uuid.UUID) (*ProductGroup, error)
	FindAll(ctx context.Context) ([]*ProductGroup, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// AttributeRepository defines the interface for interacting with Attribute data.
type AttributeRepository interface {
	Save(ctx context.Context, attribute *Attribute) error
	FindByID(ctx context.Context, id uuid.UUID) (*Attribute, error)
	FindByCode(ctx context.Context, code string) (*Attribute, error)
	FindByIDs(ctx context.Context, ids []uuid.UUID) ([]Attribute, error)
	FindByScope(ctx context.Context, brandID *uuid.UUID, groupID *uuid.UUID) ([]*Attribute, error)
	Delete(ctx context.Context, id uuid.UUID) error
	// FindApplicableAttributes(ctx context.Context, productID uuid.UUID, brandID uuid.UUID, groupIDs []uuid.UUID) ([]Attribute, error)
	// Add other necessary methods for Attribute (e.g., GetAttributeValues)
}

// ProductRepository defines the interface for interacting with Product data.
type ProductRepository interface {
	Save(ctx context.Context, product *Product) error
	FindByID(ctx context.Context, id uuid.UUID) (*Product, error)
	FindByIDs(ctx context.Context, ids []uuid.UUID) ([]*Product, error)
	FindBySKU(ctx context.Context, sku string) (*Product, error)
	FindByBarcode(ctx context.Context, barcode string) (*Product, error)
	FindBySKUPrefix(ctx context.Context, prefix string) ([]*Product, error)
	FindAll(ctx context.Context) ([]*Product, error)
	UpdateSKUs(ctx context.Context, productID uuid.UUID, newSKU string) error // For the SKU cascade
}

// ProductVariantRepository defines the interface for interacting with ProductVariant data.
type ProductVariantRepository interface {
	Save(ctx context.Context, variant *ProductVariant) error
	FindByID(ctx context.Context, id uuid.UUID) (*ProductVariant, error)
	FindByIDs(ctx context.Context, ids []uuid.UUID) ([]*ProductVariant, error)
	FindBySKU(ctx context.Context, sku string) (*ProductVariant, error)
	FindByBarcode(ctx context.Context, barcode string) (*ProductVariant, error)
	FindBySKUPrefix(ctx context.Context, prefix string) ([]*ProductVariant, error)
	FindByProductID(ctx context.Context, productID uuid.UUID) ([]*ProductVariant, error)
	FindByProductIDAndAttributeValues(ctx context.Context, productID uuid.UUID, attributeValueIDs []uuid.UUID) (*ProductVariant, error)
}

// PartyServiceConfigurationRepository defines the interface for interacting with PartyServiceConfiguration data.
type PartyServiceConfigurationRepository interface {
	Save(ctx context.Context, config *PartyServiceConfiguration) error
	FindByID(ctx context.Context, partyID, id uuid.UUID) (*PartyServiceConfiguration, error)
	FindByPartyID(ctx context.Context, partyID uuid.UUID) ([]*PartyServiceConfiguration, error)
	Delete(ctx context.Context, partyID, id uuid.UUID) error
}
