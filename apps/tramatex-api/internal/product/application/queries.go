package application

import (
	"github.com/google/uuid"
)

// GetAttributeByIDQuery defines the query for fetching a single attribute by ID.
type GetAttributeByIDQuery struct {
	ID uuid.UUID
}

// ListAttributesQuery defines the query for listing attributes with optional filtering.
// Note: Scope-based filtering removed for MVP simplicity.
type ListAttributesQuery struct {
	// No filters needed for MVP - all attributes are generic
}

// GetProductByIDQuery defines the query for fetching a single product by ID.
type GetProductByIDQuery struct {
	ID uuid.UUID
}

// ListProductsQuery defines the query for listing products with optional filtering.
type ListProductsQuery struct {
	Search      *string
	BrandID     *uuid.UUID
	GroupID     *uuid.UUID
	IsActive    *bool
	ProductType *string
}

// ListProductVariantsByProductIDQuery defines the query for listing product variants by product ID.
type ListProductVariantsByProductIDQuery struct {
	ProductID uuid.UUID
}

// GetProductVariantByIDQuery defines the query for fetching a single product variant by ID.
type GetProductVariantByIDQuery struct {
	ID uuid.UUID
}

// GetProductVariantBySKUQuery defines the query for fetching a single product variant by SKU.
type GetProductVariantBySKUQuery struct {
	SKU string
}

// SmartSearchQuery defines the query for intelligent product/variant search by SKU, barcode, or partial reference.
type SmartSearchQuery struct {
	Query string
}

// GetPartyServiceConfigurationByIDQuery defines the query for fetching a single party service configuration by ID.
type GetPartyServiceConfigurationByIDQuery struct {
	PartyID uuid.UUID
	ID      uuid.UUID
}

// ListPartyServiceConfigurationsByPartyIDQuery defines the query for listing party service configurations by party ID.
type ListPartyServiceConfigurationsByPartyIDQuery struct {
	PartyID uuid.UUID
}
