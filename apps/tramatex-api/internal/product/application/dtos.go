package application

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/joran-cortez/tramatex/internal/product/domain"
)

// ProductDTO represents the data transfer object for a Product.
type ProductDTO struct {
	ID                 uuid.UUID          `json:"id"`
	SKU                string             `json:"sku"`
	Name               string             `json:"name"`
	LongName           string             `json:"longName"`
	Barcode            *string            `json:"barcode,omitempty"`
	Description        string             `json:"description"`
	ProductType        domain.ProductType `json:"productType"`
	BrandID            uuid.UUID          `json:"brandId"`
	GroupIDs           []uuid.UUID        `json:"groupIds"` // Note: api-contracts has single 'groupId', but domain has multiple
	DirectAttributeIDs []uuid.UUID        `json:"directAttributeIds"`
	BasePrice          float64            `json:"basePrice"`
	TaxRate            float64            `json:"taxRate"`
	IsActive           bool               `json:"isActive"`
}

// NewProductDTOFromDomain creates a ProductDTO from a domain.Product entity.
func NewProductDTOFromDomain(p *domain.Product) *ProductDTO {
	var barcode *string
	if p.Barcode != nil {
		barcode = p.Barcode
	}
	return &ProductDTO{
		ID:                 p.ID,
		SKU:                p.SKU,
		Name:               p.Name,
		LongName:           p.LongName,
		Barcode:            barcode,
		Description:        p.Description,
		ProductType:        p.ProductType,
		BrandID:            p.BrandID,
		GroupIDs:           p.GroupIDs,
		DirectAttributeIDs: p.DirectAttributeIDs,
		BasePrice:          p.BasePrice,
		TaxRate:            p.TaxRate,
		IsActive:           p.IsActive,
	}
}

// AttributeValueDTO represents the data transfer object for an AttributeValue.
type AttributeValueDTO struct {
	ID               uuid.UUID `json:"id"`
	Value            string    `json:"value"`
	Code             string    `json:"code"`
	HasPriceModifier bool      `json:"hasPriceModifier"`
	ModifierType     string    `json:"modifierType,omitempty"` // "FIXED" or "PERCENTAGE"
	ModifierAmount   float64   `json:"modifierAmount,omitempty"`
}

// AttributeDTO represents the data transfer object for an Attribute (ProductOptionSet in API).
// Note: Scope fields removed for MVP simplicity.
type AttributeDTO struct {
	ID        uuid.UUID           `json:"id"`
	Name      string              `json:"name"`
	Code      string              `json:"code"` // Changed from AttributeName for consistency
	SortOrder int                 `json:"sortOrder"`
	Values    []AttributeValueDTO `json:"values"` // Full value objects with ID, value, and code
}

// NewAttributeDTOFromDomain creates an AttributeDTO from a domain.Attribute entity.
func NewAttributeDTOFromDomain(a *domain.Attribute) *AttributeDTO {
	values := make([]AttributeValueDTO, len(a.Values))
	for i, av := range a.Values {
		values[i] = AttributeValueDTO{
			ID:               av.ID,
			Value:            av.Value,
			Code:             av.Code,
			HasPriceModifier: av.HasPriceModifier,
			ModifierType:     string(av.ModifierType),
			ModifierAmount:   av.ModifierAmount,
		}
	}

	return &AttributeDTO{
		ID:        a.ID,
		Name:      a.Name,
		Code:      a.Code,
		SortOrder: a.SortOrder,
		Values:    values,
	}
}

// ProductVariantDTO represents the data transfer object for a ProductVariant.
type ProductVariantDTO struct {
	ID                  uuid.UUID            `json:"id"`
	ProductID           uuid.UUID            `json:"productId"`
	SKU                 string               `json:"sku"`
	Barcode             *string              `json:"barcode,omitempty"`
	BaseCost            float64              `json:"baseCost"` // Calculated: Product.BasePrice + AttributeValue modifiers (NOT stored)
	Price               float64              `json:"price"`    // Final sales price from pricing module
	Status              domain.VariantStatus `json:"status"`
	OptionConfiguration map[string]string    `json:"optionConfiguration"` // AttributeName -> Value
	IsActive            bool                 `json:"isActive"`
	ProductName         string               `json:"productName,omitempty"` // Enriched field for smart search results
}

// NewProductVariantDTOFromDomain creates a ProductVariantDTO from a domain.ProductVariant entity.
// NOTE: Populating OptionConfiguration requires fetching Attribute and AttributeValue data.
// The BaseCost is calculated dynamically from Product.BasePrice + AttributeValue price modifiers.
func NewProductVariantDTOFromDomain(v *domain.ProductVariant, product *domain.Product, allAttributes []*domain.Attribute) *ProductVariantDTO {
	optionConfig := make(map[string]string)
	var variantAttributeValues []domain.AttributeValue // To calculate baseCost

	if allAttributes != nil {
		attributeMap := make(map[uuid.UUID]*domain.Attribute)
		attributeValueMap := make(map[uuid.UUID]*domain.AttributeValue)

		for _, attr := range allAttributes {
			attributeMap[attr.ID] = attr
			for _, val := range attr.Values {
				attributeValueMap[val.ID] = &val
			}
		}

		for _, avID := range v.AttributeValues {
			if av, found := attributeValueMap[avID]; found {
				// Collect attribute values for price calculation
				variantAttributeValues = append(variantAttributeValues, *av)

				// Find the parent attribute for this attribute value
				for _, attr := range allAttributes {
					for _, value := range attr.Values {
						if value.ID == avID {
							optionConfig[attr.Name] = av.Value // Using Attribute Name as key
							break
						}
					}
					if _, exists := optionConfig[attr.Name]; exists {
						break
					}
				}
			}
		}
	}

	// Calculate BaseCost: Product.BasePrice + AttributeValue modifiers
	baseCost := 0.0
	if product != nil {
		baseCost = domain.CalculateBaseCost(product.BasePrice, variantAttributeValues)
	}

	return &ProductVariantDTO{
		ID:                  v.ID,
		ProductID:           v.ProductID,
		SKU:                 v.SKU,
		Barcode:             v.Barcode,
		BaseCost:            baseCost,
		Price:               0.0, // Placeholder, needs to come from pricing service/domain
		Status:              v.Status,
		OptionConfiguration: optionConfig,
		IsActive:            v.IsActive,
	}
}

// PartyServiceConfigurationDTO represents the data transfer object for a PartyServiceConfiguration.
type PartyServiceConfigurationDTO struct {
	ID                   uuid.UUID       `json:"id"`
	PartyID              uuid.UUID       `json:"partyId"`
	ServiceID            string          `json:"serviceId"`
	Name                 string          `json:"name"`
	ConfigurationDetails json.RawMessage `json:"configurationDetails"` // Flexible JSON object
}

// SmartSearchResultDTO represents the result of a smart search for products/variants.
// Type values: "exact_variant", "exact_product", "partial_match", "product_list", "no_match"
type SmartSearchResultDTO struct {
	Type               string               `json:"type"`
	Product            *ProductDTO          `json:"product,omitempty"`
	Variant            *ProductVariantDTO   `json:"variant,omitempty"`
	Products           []*ProductDTO        `json:"products,omitempty"`
	OptionSets         []*AttributeDTO      `json:"optionSets,omitempty"`
	SelectedAttributes map[string]string    `json:"selectedAttributes,omitempty"`
	MatchingVariants   []*ProductVariantDTO `json:"matchingVariants,omitempty"`
}
