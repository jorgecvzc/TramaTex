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
	ID    uuid.UUID `json:"id"`
	Value string    `json:"value"`
	Code  string    `json:"code"`
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
			ID:    av.ID,
			Value: av.Value,
			Code:  av.Code,
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
	BaseCost            float64              `json:"baseCost"`
	Price               float64              `json:"price"` // Missing in domain, assuming it will be added or fetched from pricing
	Status              domain.VariantStatus `json:"status"`
	OptionConfiguration map[string]string    `json:"optionConfiguration"` // AttributeName -> Value
	IsActive            bool                 `json:"isActive"`
}

// NewProductVariantDTOFromDomain creates a ProductVariantDTO from a domain.ProductVariant entity.
// NOTE: Populating OptionConfiguration requires fetching Attribute and AttributeValue data.
// For now, it will be an empty map or partially populated if attribute data is available.
func NewProductVariantDTOFromDomain(v *domain.ProductVariant, allAttributes []*domain.Attribute) *ProductVariantDTO {
	optionConfig := make(map[string]string)
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

	return &ProductVariantDTO{
		ID:                  v.ID,
		ProductID:           v.ProductID,
		SKU:                 v.SKU,
		Barcode:             v.Barcode,
		BaseCost:            v.BaseCost,
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
