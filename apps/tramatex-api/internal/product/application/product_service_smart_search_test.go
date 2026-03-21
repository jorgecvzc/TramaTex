package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/joran-cortez/tramatex/internal/product/application"
	"github.com/joran-cortez/tramatex/internal/product/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============================================================================
// SmartSearch Tests
// ============================================================================

func newSmartSearchService() (*application.ProductService, *MockProductRepository, *MockBrandRepository, *MockProductGroupRepository, *MockAttributeRepository, *MockProductVariantRepository, *MockPartyServiceConfigurationRepository) {
	return newTestService()
}

func TestProductService_SmartSearch(t *testing.T) {
	ctx := actorCtx()

	productID := uuid.New()
	brandID := uuid.New()
	product := &domain.Product{
		ID:          productID,
		SKU:         "FY5678",
		Name:        "Test Fabric",
		LongName:    "Test Fabric Long Name",
		ProductType: domain.ProductTypeTangible,
		BrandID:     brandID,
		GroupIDs:    []uuid.UUID{uuid.New()},
		BasePrice:   10.0,
		IsActive:    true,
	}

	variantID := uuid.New()
	variant := &domain.ProductVariant{
		ID:              variantID,
		ProductID:       productID,
		SKU:             "FY5678-SIZE.M-COLOR.RED",
		Status:          domain.StatusConfirmed,
		AttributeValues: []uuid.UUID{uuid.New(), uuid.New()},
		IsActive:        true,
	}

	t.Run("should return no_match for empty query", func(t *testing.T) {
		svc, _, _, _, _, _, _ := newSmartSearchService()
		result, err := svc.SmartSearch(ctx, application.SmartSearchQuery{Query: ""})

		assert.NoError(t, err)
		assert.Equal(t, "no_match", result.Type)
	})

	t.Run("should return no_match for whitespace-only query", func(t *testing.T) {
		svc, _, _, _, _, _, _ := newSmartSearchService()
		result, err := svc.SmartSearch(ctx, application.SmartSearchQuery{Query: "   "})

		assert.NoError(t, err)
		assert.Equal(t, "no_match", result.Type)
	})

	t.Run("should find exact variant by SKU", func(t *testing.T) {
		svc, mockProductRepo, _, _, mockAttributeRepo, mockVariantRepo, _ := newSmartSearchService()

		mockVariantRepo.On("FindBySKU", ctx, "FY5678-SIZE.M-COLOR.RED").Return(variant, nil).Once()
		mockProductRepo.On("FindByID", ctx, productID).Return(product, nil).Once()
		mockAttributeRepo.On("FindByScope", ctx, (*uuid.UUID)(nil), (*uuid.UUID)(nil)).Return([]*domain.Attribute{}, nil).Once()

		result, err := svc.SmartSearch(ctx, application.SmartSearchQuery{Query: "FY5678-SIZE.M-COLOR.RED"})

		assert.NoError(t, err)
		assert.Equal(t, "exact_variant", result.Type)
		assert.NotNil(t, result.Variant)
		assert.NotNil(t, result.Product)
		assert.Equal(t, product.Name, result.Product.Name)
	})

	t.Run("should find exact variant by barcode", func(t *testing.T) {
		svc, mockProductRepo, _, _, mockAttributeRepo, mockVariantRepo, _ := newSmartSearchService()

		mockVariantRepo.On("FindBySKU", ctx, "1234567890").Return(nil, nil).Once()
		mockVariantRepo.On("FindByBarcode", ctx, "1234567890").Return(variant, nil).Once()
		mockProductRepo.On("FindByID", ctx, productID).Return(product, nil).Once()
		mockAttributeRepo.On("FindByScope", ctx, (*uuid.UUID)(nil), (*uuid.UUID)(nil)).Return([]*domain.Attribute{}, nil).Once()

		result, err := svc.SmartSearch(ctx, application.SmartSearchQuery{Query: "1234567890"})

		assert.NoError(t, err)
		assert.Equal(t, "exact_variant", result.Type)
	})

	t.Run("should find exact product by SKU", func(t *testing.T) {
		svc, mockProductRepo, _, _, mockAttributeRepo, mockVariantRepo, _ := newSmartSearchService()

		mockVariantRepo.On("FindBySKU", ctx, "FY5678").Return(nil, nil).Once()
		mockVariantRepo.On("FindByBarcode", ctx, "FY5678").Return(nil, nil).Once()
		mockProductRepo.On("FindBySKU", ctx, "FY5678").Return(product, nil).Once()
		// buildExactProductResult calls GetApplicableAttributesForProduct
		mockProductRepo.On("FindByID", ctx, productID).Return(product, nil).Once()
		mockAttributeRepo.On("FindByIDs", ctx, mock.Anything).Return([]domain.Attribute{}, nil).Once()
		mockAttributeRepo.On("FindByScope", ctx, mock.Anything, mock.Anything).Return([]*domain.Attribute{}, nil).Once()

		result, err := svc.SmartSearch(ctx, application.SmartSearchQuery{Query: "FY5678"})

		assert.NoError(t, err)
		assert.Equal(t, "exact_product", result.Type)
		assert.NotNil(t, result.Product)
		assert.Equal(t, "FY5678", result.Product.SKU)
	})

	t.Run("should find exact product by barcode", func(t *testing.T) {
		svc, mockProductRepo, _, _, mockAttributeRepo, mockVariantRepo, _ := newSmartSearchService()

		mockVariantRepo.On("FindBySKU", ctx, "BARCODE123").Return(nil, nil).Once()
		mockVariantRepo.On("FindByBarcode", ctx, "BARCODE123").Return(nil, nil).Once()
		mockProductRepo.On("FindBySKU", ctx, "BARCODE123").Return(nil, nil).Once()
		mockProductRepo.On("FindByBarcode", ctx, "BARCODE123").Return(product, nil).Once()
		// buildExactProductResult calls
		mockProductRepo.On("FindByID", ctx, productID).Return(product, nil).Once()
		mockAttributeRepo.On("FindByIDs", ctx, mock.Anything).Return([]domain.Attribute{}, nil).Once()
		mockAttributeRepo.On("FindByScope", ctx, mock.Anything, mock.Anything).Return([]*domain.Attribute{}, nil).Once()

		result, err := svc.SmartSearch(ctx, application.SmartSearchQuery{Query: "BARCODE123"})

		assert.NoError(t, err)
		assert.Equal(t, "exact_product", result.Type)
	})

	t.Run("should find by variant SKU prefix with dash", func(t *testing.T) {
		svc, mockProductRepo, _, _, mockAttributeRepo, mockVariantRepo, _ := newSmartSearchService()

		mockVariantRepo.On("FindBySKU", ctx, "FY5678-SIZE.M").Return(nil, nil).Once()
		mockVariantRepo.On("FindByBarcode", ctx, "FY5678-SIZE.M").Return(nil, nil).Once()
		mockProductRepo.On("FindBySKU", ctx, "FY5678-SIZE.M").Return(nil, nil).Once()
		mockProductRepo.On("FindByBarcode", ctx, "FY5678-SIZE.M").Return(nil, nil).Once()
		// Prefix match branch
		mockVariantRepo.On("FindBySKUPrefix", ctx, "FY5678-SIZE.M").Return([]*domain.ProductVariant{variant}, nil).Once()
		// buildPartialMatchResult calls: first FindByID for the product from variant, then GetApplicableAttributesForProduct which calls FindByID again
		mockProductRepo.On("FindByID", ctx, productID).Return(product, nil).Times(2)
		mockAttributeRepo.On("FindByIDs", ctx, mock.Anything).Return([]domain.Attribute{}, nil).Once()
		mockAttributeRepo.On("FindByScope", ctx, mock.Anything, mock.Anything).Return([]*domain.Attribute{}, nil).Once()

		result, err := svc.SmartSearch(ctx, application.SmartSearchQuery{Query: "FY5678-SIZE.M"})

		assert.NoError(t, err)
		assert.Equal(t, "partial_match", result.Type)
		assert.NotNil(t, result.Product)
		assert.NotEmpty(t, result.SelectedAttributes)
	})

	t.Run("should find by product SKU prefix (single match)", func(t *testing.T) {
		svc, mockProductRepo, _, _, mockAttributeRepo, mockVariantRepo, _ := newSmartSearchService()

		mockVariantRepo.On("FindBySKU", ctx, "FY56").Return(nil, nil).Once()
		mockVariantRepo.On("FindByBarcode", ctx, "FY56").Return(nil, nil).Once()
		mockProductRepo.On("FindBySKU", ctx, "FY56").Return(nil, nil).Once()
		mockProductRepo.On("FindByBarcode", ctx, "FY56").Return(nil, nil).Once()
		// Product prefix match
		mockProductRepo.On("FindBySKUPrefix", ctx, "FY56").Return([]*domain.Product{product}, nil).Once()
		// buildExactProductResult
		mockProductRepo.On("FindByID", ctx, productID).Return(product, nil).Once()
		mockAttributeRepo.On("FindByIDs", ctx, mock.Anything).Return([]domain.Attribute{}, nil).Once()
		mockAttributeRepo.On("FindByScope", ctx, mock.Anything, mock.Anything).Return([]*domain.Attribute{}, nil).Once()

		result, err := svc.SmartSearch(ctx, application.SmartSearchQuery{Query: "FY56"})

		assert.NoError(t, err)
		assert.Equal(t, "exact_product", result.Type)
	})

	t.Run("should return product_list for multiple prefix matches", func(t *testing.T) {
		svc, mockProductRepo, _, _, _, mockVariantRepo, _ := newSmartSearchService()

		product2 := &domain.Product{ID: uuid.New(), SKU: "FY9999", Name: "Another Fabric", BrandID: brandID, IsActive: true}

		mockVariantRepo.On("FindBySKU", ctx, "FY").Return(nil, nil).Once()
		mockVariantRepo.On("FindByBarcode", ctx, "FY").Return(nil, nil).Once()
		mockProductRepo.On("FindBySKU", ctx, "FY").Return(nil, nil).Once()
		mockProductRepo.On("FindByBarcode", ctx, "FY").Return(nil, nil).Once()
		mockProductRepo.On("FindBySKUPrefix", ctx, "FY").Return([]*domain.Product{product, product2}, nil).Once()

		result, err := svc.SmartSearch(ctx, application.SmartSearchQuery{Query: "FY"})

		assert.NoError(t, err)
		assert.Equal(t, "product_list", result.Type)
		assert.Len(t, result.Products, 2)
	})

	t.Run("should fall back to text search by name", func(t *testing.T) {
		svc, mockProductRepo, _, _, _, mockVariantRepo, _ := newSmartSearchService()

		mockVariantRepo.On("FindBySKU", ctx, "fabric").Return(nil, nil).Once()
		mockVariantRepo.On("FindByBarcode", ctx, "fabric").Return(nil, nil).Once()
		mockProductRepo.On("FindBySKU", ctx, "fabric").Return(nil, nil).Once()
		mockProductRepo.On("FindByBarcode", ctx, "fabric").Return(nil, nil).Once()
		mockProductRepo.On("FindBySKUPrefix", ctx, "fabric").Return([]*domain.Product{}, nil).Once()
		mockProductRepo.On("FindAll", ctx).Return([]*domain.Product{product}, nil).Once()

		result, err := svc.SmartSearch(ctx, application.SmartSearchQuery{Query: "fabric"})

		assert.NoError(t, err)
		assert.Equal(t, "product_list", result.Type)
		assert.Len(t, result.Products, 1)
		assert.Equal(t, "Test Fabric", result.Products[0].Name)
	})

	t.Run("should return no_match when nothing found", func(t *testing.T) {
		svc, mockProductRepo, _, _, _, mockVariantRepo, _ := newSmartSearchService()

		mockVariantRepo.On("FindBySKU", ctx, "ZZZZZZZ").Return(nil, nil).Once()
		mockVariantRepo.On("FindByBarcode", ctx, "ZZZZZZZ").Return(nil, nil).Once()
		mockProductRepo.On("FindBySKU", ctx, "ZZZZZZZ").Return(nil, nil).Once()
		mockProductRepo.On("FindByBarcode", ctx, "ZZZZZZZ").Return(nil, nil).Once()
		mockProductRepo.On("FindBySKUPrefix", ctx, "ZZZZZZZ").Return([]*domain.Product{}, nil).Once()
		mockProductRepo.On("FindAll", ctx).Return([]*domain.Product{}, nil).Once()

		result, err := svc.SmartSearch(ctx, application.SmartSearchQuery{Query: "ZZZZZZZ"})

		assert.NoError(t, err)
		assert.Equal(t, "no_match", result.Type)
	})

	t.Run("should return error on variant SKU lookup failure", func(t *testing.T) {
		svc, _, _, _, _, mockVariantRepo, _ := newSmartSearchService()
		mockVariantRepo.On("FindBySKU", ctx, "ERR").Return(nil, errors.New("db error")).Once()

		result, err := svc.SmartSearch(ctx, application.SmartSearchQuery{Query: "ERR"})

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("should extract product SKU from partial reference with attr.value", func(t *testing.T) {
		svc, mockProductRepo, _, _, mockAttributeRepo, mockVariantRepo, _ := newSmartSearchService()

		// Query contains dash with attribute-like segments
		mockVariantRepo.On("FindBySKU", ctx, "ABC-SIZE.L").Return(nil, nil).Once()
		mockVariantRepo.On("FindByBarcode", ctx, "ABC-SIZE.L").Return(nil, nil).Once()
		mockProductRepo.On("FindBySKU", ctx, "ABC-SIZE.L").Return(nil, nil).Once()
		mockProductRepo.On("FindByBarcode", ctx, "ABC-SIZE.L").Return(nil, nil).Once()
		// variant prefix lookup
		mockVariantRepo.On("FindBySKUPrefix", ctx, "ABC-SIZE.L").Return([]*domain.ProductVariant{}, nil).Once()
		// extractProductSKUFromPartialRef("ABC-SIZE.L") → "ABC"
		abcProduct := &domain.Product{ID: uuid.New(), SKU: "ABC", Name: "ABC Product", BrandID: brandID, IsActive: true}
		mockProductRepo.On("FindBySKU", ctx, "ABC").Return(abcProduct, nil).Once()
		// buildPartialMatchResult → GetApplicableAttributesForProduct
		mockProductRepo.On("FindByID", ctx, abcProduct.ID).Return(abcProduct, nil).Once()
		mockAttributeRepo.On("FindByIDs", ctx, mock.Anything).Return([]domain.Attribute{}, nil).Once()
		mockAttributeRepo.On("FindByScope", ctx, mock.Anything, mock.Anything).Return([]*domain.Attribute{}, nil).Once()

		result, err := svc.SmartSearch(ctx, application.SmartSearchQuery{Query: "ABC-SIZE.L"})

		assert.NoError(t, err)
		assert.Equal(t, "partial_match", result.Type)
		assert.Equal(t, "ABC", result.Product.SKU)
	})
}

// ============================================================================
// GenerateProductVariants Tests
// ============================================================================

func TestProductService_GenerateProductVariants(t *testing.T) {
	svc, mockProductRepo, _, _, mockAttributeRepo, mockVariantRepo, _ := newTestService()
	ctx := actorCtx()

	productID := uuid.New()
	product := &domain.Product{
		ID:                 productID,
		SKU:                "FY5678",
		Name:               "Test Fabric",
		ProductType:        domain.ProductTypeTangible,
		BrandID:            uuid.New(),
		GroupIDs:           []uuid.UUID{uuid.New()},
		DirectAttributeIDs: []uuid.UUID{},
		BasePrice:          10.0,
		IsActive:           true,
	}

	t.Run("should return error when product not found", func(t *testing.T) {
		mockProductRepo.ExpectedCalls = nil
		nonExistentID := uuid.New()
		mockProductRepo.On("FindByID", ctx, nonExistentID).Return(nil, nil).Once()

		err := svc.GenerateProductVariants(ctx, application.GenerateProductVariantsCommand{
			ActorID:   testActorID,
			ProductID: nonExistentID,
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "does not exist")
	})

	t.Run("should return nil when no applicable attributes", func(t *testing.T) {
		mockProductRepo.ExpectedCalls = nil
		mockAttributeRepo.ExpectedCalls = nil
		mockVariantRepo.ExpectedCalls = nil

		mockProductRepo.On("FindByID", ctx, productID).Return(product, nil).Once()
		// GetApplicableAttributesForProduct → FindByID (for product) + FindByIDs (for direct attrs) + FindByScope
		mockProductRepo.On("FindByID", ctx, productID).Return(product, nil).Once()
		mockAttributeRepo.On("FindByIDs", ctx, mock.Anything).Return([]domain.Attribute{}, nil).Once()
		mockAttributeRepo.On("FindByScope", ctx, mock.Anything, mock.Anything).Return([]*domain.Attribute{}, nil).Once()
		// ensureDefaultVariant calls
		mockVariantRepo.On("FindByProductIDAndAttributeValues", ctx, productID, []uuid.UUID{}).Return(nil, nil).Once()
		mockVariantRepo.On("Save", ctx, mock.AnythingOfType("*domain.ProductVariant")).Return(nil).Once()

		err := svc.GenerateProductVariants(ctx, application.GenerateProductVariantsCommand{
			ActorID:   testActorID,
			ProductID: productID,
		})

		assert.NoError(t, err)
	})

	t.Run("should generate variants for single attribute with values", func(t *testing.T) {
		mockProductRepo.ExpectedCalls = nil
		mockAttributeRepo.ExpectedCalls = nil
		mockVariantRepo.ExpectedCalls = nil
		mockVariantRepo.Calls = nil

		sizeAttrID := uuid.New()
		valS := domain.AttributeValue{ID: uuid.New(), Value: "S", Code: "S"}
		valM := domain.AttributeValue{ID: uuid.New(), Value: "M", Code: "M"}
		sizeAttr := &domain.Attribute{ID: sizeAttrID, Name: "Size", Code: "SIZE", SortOrder: 1, Values: []domain.AttributeValue{valS, valM}}

		// Product has direct attribute
		productWithAttr := &domain.Product{
			ID:                 productID,
			SKU:                "FY5678",
			Name:               "Test Fabric",
			ProductType:        domain.ProductTypeTangible,
			BrandID:            uuid.New(),
			GroupIDs:           []uuid.UUID{uuid.New()},
			DirectAttributeIDs: []uuid.UUID{sizeAttrID},
			BasePrice:          10.0,
			IsActive:           true,
		}

		mockProductRepo.On("FindByID", ctx, productID).Return(productWithAttr, nil)
		// GetApplicableAttributesForProduct
		mockAttributeRepo.On("FindByIDs", ctx, mock.Anything).Return([]domain.Attribute{*sizeAttr}, nil).Once()
		mockAttributeRepo.On("FindByScope", ctx, mock.Anything, mock.Anything).Return([]*domain.Attribute{}, nil).Once()
		// Full domain attribute retrieval
		mockAttributeRepo.On("FindByID", ctx, sizeAttrID).Return(sizeAttr, nil).Once()

		// For each combination: check existing, create new
		mockVariantRepo.On("FindByProductIDAndAttributeValues", ctx, productID, mock.Anything).Return(nil, nil)
		mockVariantRepo.On("Save", ctx, mock.AnythingOfType("*domain.ProductVariant")).Return(nil)

		err := svc.GenerateProductVariants(ctx, application.GenerateProductVariantsCommand{
			ActorID:   testActorID,
			ProductID: productID,
		})

		assert.NoError(t, err)
		// Should have saved 2 variants (S and M)
		mockVariantRepo.AssertNumberOfCalls(t, "Save", 2)
	})

	t.Run("should skip existing variants and not re-save unchanged", func(t *testing.T) {
		svc2, mockProductRepo2, _, _, mockAttributeRepo2, mockVariantRepo2, _ := newTestService()

		sizeAttrID := uuid.New()
		valS := domain.AttributeValue{ID: uuid.New(), Value: "S", Code: "S"}
		sizeAttr := &domain.Attribute{ID: sizeAttrID, Name: "Size", Code: "SIZE", SortOrder: 1, Values: []domain.AttributeValue{valS}}

		productWithAttr := &domain.Product{
			ID:                 productID,
			SKU:                "FY5678",
			Name:               "Test Fabric",
			ProductType:        domain.ProductTypeTangible,
			BrandID:            uuid.New(),
			GroupIDs:           []uuid.UUID{uuid.New()},
			DirectAttributeIDs: []uuid.UUID{sizeAttrID},
			BasePrice:          10.0,
			IsActive:           true,
		}

		existingVariant := &domain.ProductVariant{
			ID:              uuid.New(),
			ProductID:       productID,
			SKU:             "FY5678-SIZE.S",
			Status:          domain.StatusConfirmed,
			AttributeValues: []uuid.UUID{valS.ID},
			IsActive:        true,
		}

		mockProductRepo2.On("FindByID", ctx, productID).Return(productWithAttr, nil)
		mockAttributeRepo2.On("FindByIDs", ctx, mock.Anything).Return([]domain.Attribute{*sizeAttr}, nil).Once()
		mockAttributeRepo2.On("FindByScope", ctx, mock.Anything, mock.Anything).Return([]*domain.Attribute{}, nil).Once()
		mockAttributeRepo2.On("FindByID", ctx, sizeAttrID).Return(sizeAttr, nil).Once()
		mockVariantRepo2.On("FindByProductIDAndAttributeValues", ctx, productID, []uuid.UUID{valS.ID}).Return(existingVariant, nil).Once()

		err := svc2.GenerateProductVariants(ctx, application.GenerateProductVariantsCommand{
			ActorID:   testActorID,
			ProductID: productID,
		})

		assert.NoError(t, err)
		// Should NOT save since variant already exists with correct SKU, status and isActive
		mockVariantRepo2.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	t.Run("should return nil when attribute has no values", func(t *testing.T) {
		svc3, mockProductRepo3, _, _, mockAttributeRepo3, mockVariantRepo3, _ := newTestService()

		emptyAttrID := uuid.New()
		emptyAttr := &domain.Attribute{ID: emptyAttrID, Name: "Color", Code: "COLOR", SortOrder: 1, Values: []domain.AttributeValue{}}

		productWithAttr := &domain.Product{
			ID:                 productID,
			SKU:                "FY5678",
			ProductType:        domain.ProductTypeTangible,
			BrandID:            uuid.New(),
			GroupIDs:           []uuid.UUID{uuid.New()},
			DirectAttributeIDs: []uuid.UUID{emptyAttrID},
			BasePrice:          10.0,
			IsActive:           true,
		}

		mockProductRepo3.On("FindByID", ctx, productID).Return(productWithAttr, nil)
		mockAttributeRepo3.On("FindByIDs", ctx, mock.Anything).Return([]domain.Attribute{*emptyAttr}, nil).Once()
		mockAttributeRepo3.On("FindByScope", ctx, mock.Anything, mock.Anything).Return([]*domain.Attribute{}, nil).Once()
		mockAttributeRepo3.On("FindByID", ctx, emptyAttrID).Return(emptyAttr, nil).Once()

		err := svc3.GenerateProductVariants(ctx, application.GenerateProductVariantsCommand{
			ActorID:   testActorID,
			ProductID: productID,
		})

		assert.NoError(t, err)
		mockVariantRepo3.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	t.Run("should return error on product repo failure", func(t *testing.T) {
		mockProductRepo.ExpectedCalls = nil
		mockProductRepo.On("FindByID", ctx, productID).Return(nil, errors.New("db error")).Once()

		err := svc.GenerateProductVariants(ctx, application.GenerateProductVariantsCommand{
			ActorID:   testActorID,
			ProductID: productID,
		})

		assert.Error(t, err)
	})
}

// ============================================================================
// GetApplicableAttributesForProductLegacy test
// ============================================================================

func TestProductService_GetApplicableAttributesForProductLegacy(t *testing.T) {
	svc, _, _, _, _, _, _ := newTestService()
	ctx := context.WithValue(context.Background(), "actorID", testActorID)

	t.Run("should return error since legacy method is not implemented", func(t *testing.T) {
		productID := uuid.New()

		result, err := svc.GetApplicableAttributesForProductLegacy(ctx, productID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "legacy method not implemented")
	})
}
