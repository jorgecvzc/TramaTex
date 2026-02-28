package domain

import (
	"fmt"
	"sort"
	"strings"

	"github.com/google/uuid"
)

// ProductVariant represents a specific, sellable instance of a Product.
// Note: BaseCost is NOT stored - it is calculated dynamically from Product.BasePrice + AttributeValue modifiers.
type ProductVariant struct {
	ID              uuid.UUID
	ProductID       uuid.UUID
	SKU             string
	Barcode         *string
	Status          VariantStatus
	AttributeValues []uuid.UUID // IDs of associated AttributeValues
	IsActive        bool
}

// NewProductVariant creates a new ProductVariant with validation.
// It requires the composed SKU to be passed in, ensuring the domain service or use case
// is responsible for its generation based on Product and Attribute data.
// Note: BaseCost is NOT a parameter - it will be calculated dynamically when needed.
func NewProductVariant(
	productID uuid.UUID,
	sku string,
	barcode *string,
	status VariantStatus,
	attributeValueIDs []uuid.UUID,
) (*ProductVariant, error) {
	if productID == uuid.Nil {
		return nil, NewValidationError("product variant must be associated with a product")
	}
	if sku == "" {
		return nil, NewValidationError("product variant SKU cannot be empty")
	}
	if err := status.IsValid(); err != nil {
		return nil, err
	}
	// A variant must be defined by at least one attribute value to distinguish it
	if len(attributeValueIDs) == 0 {
		return nil, NewValidationError("product variant must have at least one attribute value")
	}

	// Sort attributeValueIDs for deterministic comparison and consistency
	// This helps in uniquely identifying a variant by its attributes, useful for lookups.
	sort.Slice(attributeValueIDs, func(i, j int) bool {
		return attributeValueIDs[i].String() < attributeValueIDs[j].String()
	})

	return &ProductVariant{
		ID:              uuid.New(),
		ProductID:       productID,
		SKU:             sku,
		Barcode:         barcode,
		Status:          status,
		AttributeValues: attributeValueIDs,
		IsActive:        true,
	}, nil
}

// GenerateVariantSKU composes a deterministic SKU for a ProductVariant.
// It requires the base Product SKU and a list of AttributeCode-AttributeValueCode pairs.
// The order of attributeCodeValuePairs is CRUCIAL and must be pre-sorted
// by Attribute.SortOrder to ensure a deterministic SKU.
func GenerateVariantSKU(productSKU string, attributeCodeValuePairs []struct{ AttributeCode, ValueCode string }) (string, error) {
	if productSKU == "" {
		return "", NewValidationError("product SKU cannot be empty for variant SKU generation")
	}
	// A variant SKU will always start with the product SKU
	if len(attributeCodeValuePairs) == 0 {
		return productSKU, nil // This case should ideally be prevented by NewProductVariant if variants must have attributes.
	}

	var parts []string
	for _, pair := range attributeCodeValuePairs {
		if pair.AttributeCode == "" || pair.ValueCode == "" {
			return "", NewValidationError("attribute code and value code cannot be empty for variant SKU generation")
		}
		parts = append(parts, fmt.Sprintf("%s.%s", pair.AttributeCode, pair.ValueCode))
	}

	// The input `attributeCodeValuePairs` slice *must already be sorted* by Attribute.SortOrder
	// to ensure a deterministic SKU. This function assumes that pre-sorting has occurred
	// in the calling layer (e.g., an application service).
	return fmt.Sprintf("%s-%s", productSKU, strings.Join(parts, "-")), nil
}
