package application

import (
	"github.com/google/uuid"
	"github.com/joran-cortez/tramatex/internal/product/domain"
)

// CreateProductCommand is the input DTO for the CreateProduct use case.
type CreateProductCommand struct {
	ActorID            string             `json:"-"`
	SKU                string             `json:"sku"`
	Name               string             `json:"name"`
	LongName           string             `json:"long_name"`
	Barcode            *string            `json:"barcode,omitempty"`
	Description        string             `json:"description"`
	ProductType        domain.ProductType `json:"product_type"`
	BrandID            uuid.UUID          `json:"brand_id"`
	GroupIDs           []uuid.UUID        `json:"group_ids"`
	DirectAttributeIDs []uuid.UUID        `json:"direct_attribute_ids"`
	BasePrice          float64            `json:"base_price"`
	TaxRate            float64            `json:"tax_rate"` // Tax rate as percentage (e.g., 21.0 = 21%)
}

// UpdateProductCommand is the input DTO for updating a product's general information.
type UpdateProductCommand struct {
	ActorID            string
	ProductID          uuid.UUID
	Name               *string             `json:"name,omitempty"`
	LongName           *string             `json:"long_name,omitempty"`
	SKU                *string             `json:"sku,omitempty"`
	Barcode            *string             `json:"barcode,omitempty"`
	BasePrice          *float64            `json:"base_price,omitempty"`
	TaxRate            *float64            `json:"tax_rate,omitempty"` // Tax rate as percentage
	IsActive           *bool               `json:"is_active,omitempty"`
	ProductType        *domain.ProductType `json:"product_type,omitempty"`
	Description        *string             `json:"description,omitempty"`
	BrandID            *uuid.UUID          `json:"brand_id,omitempty"`
	GroupIDs           []uuid.UUID         `json:"group_ids,omitempty"`
	DirectAttributeIDs []uuid.UUID         `json:"direct_attribute_ids,omitempty"`
}

// AddGroupCommand is the input DTO for adding a group to a product.
type AddGroupCommand struct {
	ActorID   string
	ProductID uuid.UUID
	GroupID   uuid.UUID
}

// AddDirectAttributeCommand is the input DTO for adding a direct attribute to a product.
type AddDirectAttributeCommand struct {
	ActorID     string
	ProductID   uuid.UUID
	AttributeID uuid.UUID
}

// UpdateProductSKUCommand is the input DTO for updating a product's SKU.
type UpdateProductSKUCommand struct {
	ActorID   string
	ProductID uuid.UUID
	NewSKU    string
}

// CreateAttributeValueCommand is the input DTO for creating an attribute value.
type CreateAttributeValueCommand struct {
	Value            string
	Code             string
	HasPriceModifier bool
	ModifierType     string // "FIXED" or "PERCENTAGE"
	ModifierAmount   float64
}

// CreateAttributeCommand is the input DTO for the CreateAttribute use case.
// Note: Scope fields removed for MVP simplicity.
type CreateAttributeCommand struct {
	ActorID string
	Name    string
	Code    string
	Values  []CreateAttributeValueCommand
}

// UpdateAttributeValueCommand is the input DTO for updating an attribute value.
type UpdateAttributeValueCommand struct {
	ID               *uuid.UUID // ID is nil for new values, present for existing values
	Value            string
	Code             string
	HasPriceModifier bool
	ModifierType     string // "FIXED" or "PERCENTAGE"
	ModifierAmount   float64
}

// UpdateAttributeCommand is the input DTO for the UpdateAttribute use case.
// Note: Scope fields removed for MVP simplicity.
type UpdateAttributeCommand struct {
	ActorID string
	ID      uuid.UUID
	Name    *string
	Code    *string
	Values  []UpdateAttributeValueCommand // Replaces existing values
}

// PreGenerateProductVariantsCommand is the input DTO for the PreGenerateProductVariants use case.
type PreGenerateProductVariantsCommand struct {
	ProductID uuid.UUID
}

// UpdateProductVariantCommand is the input DTO for the UpdateProductVariant use case.
// Note: BaseCost is not updateable - it's calculated dynamically from Product.BasePrice + Attribute modifiers.
type UpdateProductVariantCommand struct {
	ActorID  string
	ID       uuid.UUID
	Barcode  *string               `json:"barcode"`
	IsActive *bool                 `json:"is_active"`
	Status   *domain.VariantStatus `json:"status"` // Explicitly set status, otherwise implies CONFIRMED if other fields updated from PROVISIONAL
}

// GenerateProductVariantsCommand is the input DTO for the GenerateProductVariants use case.
type GenerateProductVariantsCommand struct {
	ActorID   string
	ProductID uuid.UUID
}

// FindOrCreateProductVariantCommand is the input DTO for the FindOrCreateProductVariant use case.
type FindOrCreateProductVariantCommand struct {
	ActorID             string
	ProductID           uuid.UUID
	OptionConfiguration map[string]string // AttributeName -> Value
}

// CreatePartyServiceConfigurationCommand is the input DTO for creating a new party service configuration.
type CreatePartyServiceConfigurationCommand struct {
	ActorID              string
	PartyID              uuid.UUID
	ServiceID            string
	Name                 string
	ConfigurationDetails map[string]interface{}
}

// UpdatePartyServiceConfigurationCommand is the input DTO for updating an existing party service configuration.
type UpdatePartyServiceConfigurationCommand struct {
	ActorID              string
	ID                   uuid.UUID
	PartyID              uuid.UUID
	ServiceID            *string
	Name                 *string
	ConfigurationDetails map[string]interface{}
}

// DeletePartyServiceConfigurationCommand is the input DTO for deleting a party service configuration.
type DeletePartyServiceConfigurationCommand struct {
	ActorID string
	ID      uuid.UUID
	PartyID uuid.UUID
}

// ============================================================================
// Brand Commands
// ============================================================================

// CreateBrandCommand is the input DTO for creating a brand.
type CreateBrandCommand struct {
	ActorID                 string
	Name                    string  `json:"name"`
	DefaultMarkupPercentage float64 `json:"defaultMarkupPercentage"`
	IsActive                bool    `json:"isActive"`
}

// UpdateBrandCommand is the input DTO for updating a brand.
type UpdateBrandCommand struct {
	ActorID                 string
	ID                      uuid.UUID
	Name                    *string  `json:"name,omitempty"`
	DefaultMarkupPercentage *float64 `json:"defaultMarkupPercentage,omitempty"`
	IsActive                *bool    `json:"isActive,omitempty"`
}

// DeleteBrandCommand is the input DTO for deleting a brand.
type DeleteBrandCommand struct {
	ActorID string
	ID      uuid.UUID
}

// ============================================================================
// Product Group Commands
// ============================================================================

// CreateProductGroupCommand is the input DTO for creating a product group.
type CreateProductGroupCommand struct {
	ActorID  string
	Name     string
	Type     string // TANGIBLE or SERVICE
	ParentID *uuid.UUID
	IsActive bool
}

// UpdateProductGroupCommand is the input DTO for updating a product group.
type UpdateProductGroupCommand struct {
	ActorID  string
	ID       uuid.UUID
	Name     *string
	Type     *string // TANGIBLE or SERVICE
	ParentID *uuid.UUID
	IsActive *bool
}

// DeleteProductGroupCommand is the input DTO for deleting a product group.
type DeleteProductGroupCommand struct {
	ActorID string
	ID      uuid.UUID
}

// ============================================================================
// Attribute Commands (additional)
// ============================================================================

// DeleteAttributeCommand is the input DTO for deleting an attribute.
type DeleteAttributeCommand struct {
	ActorID string
	ID      uuid.UUID
}
