package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/joran-cortez/tramatex/internal/product/domain" // Adjust this import path
)

func TestProductType_IsValid(t *testing.T) {
	testCases := []struct {
		name        string
		productType domain.ProductType
		expectError bool
	}{
		{
			name:        "Valid Type Tangible",
			productType: domain.ProductTypeTangible,
			expectError: false,
		},
		{
			name:        "Valid Type Service",
			productType: domain.ProductTypeService,
			expectError: false,
		},
		{
			name:        "Invalid Type",
			productType: "INVALID_TYPE",
			expectError: true,
		},
		{
			name:        "Empty Type",
			productType: "",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.productType.IsValid()
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestVariantStatus_IsValid(t *testing.T) {
	testCases := []struct {
		name          string
		variantStatus domain.VariantStatus
		expectError   bool
	}{
		{
			name:          "Valid Status Provisional",
			variantStatus: domain.StatusProvisional,
			expectError:   false,
		},
		{
			name:          "Valid Status Confirmed",
			variantStatus: domain.StatusConfirmed,
			expectError:   false,
		},
		{
			name:          "Invalid Status",
			variantStatus: "INVALID_STATUS",
			expectError:   true,
		},
		{
			name:          "Empty Status",
			variantStatus: "",
			expectError:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.variantStatus.IsValid()
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
