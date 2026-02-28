package domain_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/joran-cortez/tramatex/internal/product/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewProduct(t *testing.T) {
	brandID := uuid.New()
	barcode := "1234567890123"

	testCases := []struct {
		name          string
		sku           string
		productName   string
		longName      string
		description   string
		productType   domain.ProductType
		brandID       uuid.UUID
		barcode       *string
		expectError   bool
		expectedError string
	}{
		{
			name:        "Valid Product Creation",
			sku:         "FYR2040",
			productName: "Camiseta",
			longName:    "Camiseta de Algodón",
			description: "Una camiseta cómoda",
			productType: domain.ProductTypeTangible,
			brandID:     brandID,
			barcode:     &barcode,
			expectError: false,
		},
		{
			name:          "Missing SKU",
			sku:           "",
			productName:   "Camiseta",
			productType:   domain.ProductTypeTangible,
			brandID:       brandID,
			expectError:   true,
			expectedError: "product SKU cannot be empty",
		},
		{
			name:          "Missing Name",
			sku:           "FYR2040",
			productName:   "",
			productType:   domain.ProductTypeTangible,
			brandID:       brandID,
			expectError:   true,
			expectedError: "product name cannot be empty",
		},
		{
			name:          "Invalid Product Type",
			sku:           "FYR2040",
			productName:   "Camiseta",
			productType:   "INVALID",
			brandID:       brandID,
			expectError:   true,
			expectedError: "invalid ProductType: INVALID",
		},
		{
			name:          "Missing BrandID",
			sku:           "FYR2040",
			productName:   "Camiseta",
			productType:   domain.ProductTypeTangible,
			brandID:       uuid.Nil,
			expectError:   true,
			expectedError: "product must be associated with a brand",
		},
		{
			name:        "Product without barcode",
			sku:         "FYR2041",
			productName: "Servicio de Diseño",
			productType: domain.ProductTypeService,
			brandID:     brandID,
			barcode:     nil,
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			product, err := domain.NewProduct(
				tc.sku,
				tc.productName,
				tc.longName,
				tc.description,
				tc.productType,
				tc.brandID,
				tc.barcode,
				0,
				21,
			)

			if tc.expectError {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError)
				assert.Nil(t, product)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, product)
				assert.Equal(t, tc.sku, product.SKU)
				assert.Equal(t, tc.productName, product.Name)
				assert.NotEqual(t, uuid.Nil, product.ID)
				assert.True(t, product.IsActive)
			}
		})
	}
}

func TestProduct_AddGroup(t *testing.T) {
	p, _ := domain.NewProduct("SKU1", "Name1", "", "", domain.ProductTypeTangible, uuid.New(), nil, 0, 21)
	groupID1 := uuid.New()
	groupID2 := uuid.New()

	p.AddGroup(groupID1)
	assert.Len(t, p.GroupIDs, 1)
	assert.Contains(t, p.GroupIDs, groupID1)

	p.AddGroup(groupID2)
	assert.Len(t, p.GroupIDs, 2)
	assert.Contains(t, p.GroupIDs, groupID2)

	// Test adding duplicate
	p.AddGroup(groupID1)
	assert.Len(t, p.GroupIDs, 2)
}

func TestProduct_AddDirectAttribute(t *testing.T) {
	p, _ := domain.NewProduct("SKU1", "Name1", "", "", domain.ProductTypeTangible, uuid.New(), nil, 0, 21)
	attrID1 := uuid.New()
	attrID2 := uuid.New()

	p.AddDirectAttribute(attrID1)
	assert.Len(t, p.DirectAttributeIDs, 1)
	assert.Contains(t, p.DirectAttributeIDs, attrID1)

	p.AddDirectAttribute(attrID2)
	assert.Len(t, p.DirectAttributeIDs, 2)
	assert.Contains(t, p.DirectAttributeIDs, attrID2)

	// Test adding duplicate
	p.AddDirectAttribute(attrID1)
	assert.Len(t, p.DirectAttributeIDs, 2)
}

func TestProductHelpers_IDPTR(t *testing.T) {
	brandID := uuid.New()
	groupID := uuid.New()
	brand := domain.Brand{ID: brandID}
	group := domain.ProductGroup{ID: groupID}

	assert.Equal(t, brandID, *brand.ID_PTR())
	assert.Equal(t, groupID, *group.ID_PTR())
}
