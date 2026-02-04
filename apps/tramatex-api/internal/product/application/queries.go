package application

import (
	"github.com/google/uuid"
)

// GetAttributeByIDQuery defines the query for fetching a single attribute by ID.
type GetAttributeByIDQuery struct {
	ID uuid.UUID
}

// ListAttributesQuery defines the query for listing attributes with optional filtering.
type ListAttributesQuery struct {
	ScopeType      *string    // GENERIC | BRAND | BRAND_GROUP
	BrandID        *uuid.UUID
	ProductGroupID *uuid.UUID
}

// GetProductByIDQuery defines the query for fetching a single product by ID.
type GetProductByIDQuery struct {
	ID uuid.UUID
}

// ListProductsQuery defines the query for listing products with optional filtering.
type ListProductsQuery struct {
	// Add filtering parameters here (e.g., BrandID, GroupID, IsActive)
	// For example:
	// BrandID *uuid.UUID
	// GroupID *uuid.UUID
	// IsActive *bool
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

// GetPartyServiceConfigurationByIDQuery defines the query for fetching a single party service configuration by ID.
type GetPartyServiceConfigurationByIDQuery struct {
	PartyID uuid.UUID
	ID      uuid.UUID
}

// ListPartyServiceConfigurationsByPartyIDQuery defines the query for listing party service configurations by party ID.
type ListPartyServiceConfigurationsByPartyIDQuery struct {
	PartyID uuid.UUID
}