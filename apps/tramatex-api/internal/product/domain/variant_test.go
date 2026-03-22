package domain_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/joran-cortez/tramatex/internal/product/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewProductVariant(t *testing.T) {
	productID := uuid.New()
	barcode := "PVM123456789"
	attributeValueIDs := []uuid.UUID{uuid.New(), uuid.New()}

	testCases := []struct {
		name              string
		productID         uuid.UUID
		sku               string
		barcode           *string
		status            domain.VariantStatus
		attributeValueIDs []uuid.UUID
		expectError       bool
		expectedErrorMsg  string
	}{
		{
			name:              "Valid ProductVariant Creation",
			productID:         productID,
			sku:               "PROD-RED-L",
			barcode:           &barcode,
			status:            domain.StatusProvisional,
			attributeValueIDs: attributeValueIDs,
			expectError:       false,
		},
		{
			name:              "Missing ProductID",
			productID:         uuid.Nil,
			sku:               "PROD-RED-L",
			barcode:           &barcode,
			status:            domain.StatusProvisional,
			attributeValueIDs: attributeValueIDs,
			expectError:       true,
			expectedErrorMsg:  "product variant must be associated with a product",
		},
		{
			name:              "Empty SKU",
			productID:         productID,
			sku:               "",
			barcode:           &barcode,
			status:            domain.StatusProvisional,
			attributeValueIDs: attributeValueIDs,
			expectError:       true,
			expectedErrorMsg:  "product variant SKU cannot be empty",
		},
		{
			name:              "Invalid Status",
			productID:         productID,
			sku:               "PROD-RED-L",
			barcode:           &barcode,
			status:            "INVALID_STATUS",
			attributeValueIDs: attributeValueIDs,
			expectError:       true,
			expectedErrorMsg:  "invalid VariantStatus: INVALID_STATUS",
		},
		{
			name:              "No Attribute Value IDs (default variant)",
			productID:         productID,
			sku:               "PROD-RED-L",
			barcode:           &barcode,
			status:            domain.StatusProvisional,
			attributeValueIDs: []uuid.UUID{},
			expectError:       false,
		},
		{
			name:              "ProductVariant without barcode",
			productID:         productID,
			sku:               "PROD-BLUE-M",
			barcode:           nil,
			status:            domain.StatusConfirmed,
			attributeValueIDs: attributeValueIDs,
			expectError:       false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			variant, err := domain.NewProductVariant(
				tc.productID,
				tc.sku,
				tc.barcode,
				tc.status,
				tc.attributeValueIDs,
			)

			if tc.expectError {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedErrorMsg)
				assert.Nil(t, variant)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, variant)
				assert.Equal(t, tc.sku, variant.SKU)
				assert.Equal(t, tc.status, variant.Status)
				assert.Equal(t, tc.productID, variant.ProductID)
				assert.Equal(t, len(tc.attributeValueIDs), len(variant.AttributeValues))
				assert.NotEqual(t, uuid.Nil, variant.ID)
				assert.True(t, variant.IsActive)
			}
		})
	}
}

func TestNewDefaultProductVariant(t *testing.T) {
	productID := uuid.New()
	productSKU := "SERV-001"

	variant, err := domain.NewDefaultProductVariant(productID, productSKU)
	assert.NoError(t, err)
	assert.NotNil(t, variant)
	assert.Equal(t, productID, variant.ProductID)
	assert.Equal(t, productSKU, variant.SKU)
	assert.Equal(t, domain.StatusConfirmed, variant.Status)
	assert.Empty(t, variant.AttributeValues)
	assert.True(t, variant.IsActive)
	assert.True(t, variant.IsDefault())
}

func TestIsDefault(t *testing.T) {
	// Variant with attributes → not default
	variant, err := domain.NewProductVariant(uuid.New(), "PROD-A.1", nil, domain.StatusConfirmed, []uuid.UUID{uuid.New()})
	assert.NoError(t, err)
	assert.False(t, variant.IsDefault())

	// Default variant (no attributes)
	defaultVariant, err := domain.NewDefaultProductVariant(uuid.New(), "PROD")
	assert.NoError(t, err)
	assert.True(t, defaultVariant.IsDefault())
}

func TestGenerateVariantSKU(t *testing.T) {
	testCases := []struct {
		name                    string
		productSKU              string
		attributeCodeValuePairs []struct{ AttributeCode, ValueCode string }
		expectedSKU             string
		expectError             bool
		expectedErrorMsg        string
	}{
		{
			name:       "Basic SKU Generation",
			productSKU: "FYR2040",
			attributeCodeValuePairs: []struct{ AttributeCode, ValueCode string }{
				{AttributeCode: "T", ValueCode: "L"},
				{AttributeCode: "C", ValueCode: "R"},
			},
			expectedSKU: "FYR2040-T.L-C.R",
			expectError: false,
		},
		{
			name:       "Single Attribute",
			productSKU: "PROD-A",
			attributeCodeValuePairs: []struct{ AttributeCode, ValueCode string }{
				{AttributeCode: "S", ValueCode: "M"},
			},
			expectedSKU: "PROD-A-S.M",
			expectError: false,
		},
		{
			name:                    "No Attributes (Should be handled by NewProductVariant if disallowed)",
			productSKU:              "PROD-B",
			attributeCodeValuePairs: []struct{ AttributeCode, ValueCode string }{},
			expectedSKU:             "PROD-B",
			expectError:             false,
		},
		{
			name:       "Empty Product SKU",
			productSKU: "",
			attributeCodeValuePairs: []struct{ AttributeCode, ValueCode string }{
				{AttributeCode: "T", ValueCode: "L"},
			},
			expectError:      true,
			expectedErrorMsg: "product SKU cannot be empty for variant SKU generation",
		},
		{
			name:       "Empty Attribute Code",
			productSKU: "FYR2040",
			attributeCodeValuePairs: []struct{ AttributeCode, ValueCode string }{
				{AttributeCode: "", ValueCode: "L"},
			},
			expectError:      true,
			expectedErrorMsg: "attribute code and value code cannot be empty for variant SKU generation",
		},
		{
			name:       "Empty Value Code",
			productSKU: "FYR2040",
			attributeCodeValuePairs: []struct{ AttributeCode, ValueCode string }{
				{AttributeCode: "T", ValueCode: ""},
			},
			expectError:      true,
			expectedErrorMsg: "attribute code and value code cannot be empty for variant SKU generation",
		},
		{
			name:       "Order of attributes (pre-sorted input assumed)",
			productSKU: "BASE",
			attributeCodeValuePairs: []struct{ AttributeCode, ValueCode string }{
				{AttributeCode: "C", ValueCode: "B"},
				{AttributeCode: "A", ValueCode: "Z"}, // Order comes from product.DirectAttributeIDs
			},
			expectedSKU: "BASE-C.B-A.Z", // Assumes input is already sorted as per the note in GenerateVariantSKU
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sku, err := domain.GenerateVariantSKU(tc.productSKU, tc.attributeCodeValuePairs)

			if tc.expectError {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedErrorMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedSKU, sku)
			}
		})
	}
}
