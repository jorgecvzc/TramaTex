package application

import (
	"github.com/google/uuid"
	"github.com/joran-cortez/tramatex/internal/product/domain"
)

// CreateProductCommand is the input DTO for the CreateProduct use case.
type CreateProductCommand struct {
	SKU         string
	Name        string
	LongName    string
	Barcode     *string
	Description string
	ProductType domain.ProductType
	BrandID     uuid.UUID
	GroupIDs    []uuid.UUID
}

// AddGroupCommand is the input DTO for adding a group to a product.
type AddGroupCommand struct {
	ProductID uuid.UUID
	GroupID   uuid.UUID
}

// AddDirectAttributeCommand is the input DTO for adding a direct attribute to a product.
type AddDirectAttributeCommand struct {
	ProductID   uuid.UUID
	AttributeID uuid.UUID
}

// UpdateProductSKUCommand is the input DTO for updating a product's SKU.
type UpdateProductSKUCommand struct {
	ProductID uuid.UUID
	NewSKU    string
}

// CreateAttributeValueCommand is the input DTO for creating an attribute value.
type CreateAttributeValueCommand struct {
	Value string
	Code  string
}

// CreateAttributeCommand is the input DTO for the CreateAttribute use case.
type CreateAttributeCommand struct {
	Name         string
	Code         string
	SortOrder    int
	ScopeBrandID *uuid.UUID
	ScopeGroupID *uuid.UUID
	Values       []CreateAttributeValueCommand
}

// UpdateAttributeValueCommand is the input DTO for updating an attribute value.
type UpdateAttributeValueCommand struct {
	ID    *uuid.UUID // ID is nil for new values, present for existing values
	Value string
	Code  string
}

// UpdateAttributeCommand is the input DTO for the UpdateAttribute use case.
type UpdateAttributeCommand struct {
	ID           uuid.UUID
	Name         *string
	Code         *string
	SortOrder    *int
	ScopeBrandID *uuid.UUID
	ScopeGroupID *uuid.UUID
	Values       []UpdateAttributeValueCommand // Replaces existing values
}

// PreGenerateProductVariantsCommand is the input DTO for the PreGenerateProductVariants use case.
type PreGenerateProductVariantsCommand struct {
	ProductID uuid.UUID
}

// UpdateProductVariantCommand is the input DTO for the UpdateProductVariant use case.
type UpdateProductVariantCommand struct {
	ID       uuid.UUID
	Barcode  *string
	IsActive *bool
	Status   *domain.VariantStatus // Explicitly set status, otherwise implies CONFIRMED if other fields updated from PROVISIONAL
}

// GenerateProductVariantsCommand is the input DTO for the GenerateProductVariants use case.
type GenerateProductVariantsCommand struct {
	ProductID uuid.UUID
}

// OptionConfigurationItem represents a single attribute-value pair for variant creation/finding.
type OptionConfigurationItem struct {
	AttributeName string
	Value         string
}

// FindOrCreateProductVariantCommand is the input DTO for the FindOrCreateProductVariant use case.
type FindOrCreateProductVariantCommand struct {
	ProductID           uuid.UUID
	OptionConfiguration []OptionConfigurationItem // AttributeName -> Value
}

// CreatePartyServiceConfigurationCommand is the input DTO for creating a new party service configuration.
type CreatePartyServiceConfigurationCommand struct {
	PartyID              uuid.UUID
	ServiceID            string
	Name                 string
	ConfigurationDetails map[string]interface{}
}

// UpdatePartyServiceConfigurationCommand is the input DTO for updating an existing party service configuration.
type UpdatePartyServiceConfigurationCommand struct {
	ID                   uuid.UUID
	PartyID              uuid.UUID
	ServiceID            *string
	Name                 *string
	ConfigurationDetails map[string]interface{}
}

// DeletePartyServiceConfigurationCommand is the input DTO for deleting a party service configuration.
type DeletePartyServiceConfigurationCommand struct {
	ID      uuid.UUID
	PartyID uuid.UUID
}
