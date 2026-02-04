package application

import (
	"encoding/json"
	"time"

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
	IsActive           bool               `json:"isActive"`
	CreatedAt          time.Time          `json:"createdAt"`
	UpdatedAt          time.Time          `json:"updatedAt"`
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
		IsActive:           p.IsActive,
		CreatedAt:          p.CreatedAt,
		UpdatedAt:          p.UpdatedAt,
	}
}

// AttributeDTO represents the data transfer object for an Attribute (ProductOptionSet in API).
type AttributeDTO struct {
	ID            uuid.UUID  `json:"id"`
	Name          string     `json:"name"`
	AttributeName string     `json:"attributeName"` // Corresponds to domain.Attribute.Code
	SortOrder     int        `json:"sortOrder"`
	ScopeBrandID  *uuid.UUID `json:"scopeBrandId,omitempty"`
	ScopeGroupID  *uuid.UUID `json:"scopeGroupId,omitempty"`
	Values        []string   `json:"values"` // Just the value strings
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}

// NewAttributeDTOFromDomain creates an AttributeDTO from a domain.Attribute entity.
func NewAttributeDTOFromDomain(a *domain.Attribute) *AttributeDTO {
	values := make([]string, len(a.Values))
	for i, av := range a.Values {
		values[i] = av.Value
	}

	return &AttributeDTO{
		ID:            a.ID,
		Name:          a.Name,
		AttributeName: a.Code,
		SortOrder:     a.SortOrder,
		ScopeBrandID:  a.ScopeBrandID,
		ScopeGroupID:  a.ScopeGroupID,
		Values:        values,
		CreatedAt:     a.CreatedAt,
		UpdatedAt:     a.UpdatedAt,
	}
}

// ProductVariantDTO represents the data transfer object for a ProductVariant.
type ProductVariantDTO struct {
	ID                  uuid.UUID            `json:"id"`
	ProductID           uuid.UUID            `json:"productId"`
	SKU                 string               `json:"sku"`
	Barcode             *string              `json:"barcode,omitempty"`
	Price               float64              `json:"price"` // Missing in domain, assuming it will be added or fetched from pricing
	Status              domain.VariantStatus `json:"status"`
	OptionConfiguration map[string]string    `json:"optionConfiguration"` // AttributeName -> Value
	IsActive            bool                 `json:"isActive"`
	CreatedAt           time.Time            `json:"createdAt"`
	UpdatedAt           time.Time            `json:"updatedAt"`
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
		Price:               0.0, // Placeholder, needs to come from pricing service/domain
		Status:              v.Status,
		OptionConfiguration: optionConfig,
		IsActive:            v.IsActive,
		CreatedAt:           v.CreatedAt,
		UpdatedAt:           v.UpdatedAt,
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
