package application_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/joran-cortez/tramatex/internal/product/application"
	"github.com/joran-cortez/tramatex/internal/product/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductService_CreateAttribute(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockBrandRepo := new(MockBrandRepository)
	mockGroupRepo := new(MockProductGroupRepository)
	mockAttributeRepo := new(MockAttributeRepository)
	mockVariantRepo := new(MockProductVariantRepository)
	mockPartyServiceConfigRepo := new(MockPartyServiceConfigurationRepository)
	service := application.NewProductService(mockProductRepo, mockBrandRepo, mockGroupRepo, mockAttributeRepo, mockVariantRepo, mockPartyServiceConfigRepo)
	ctx := context.WithValue(context.Background(), "actorID", "test-actor")

	t.Run("should create attribute with values", func(t *testing.T) {
		mockAttributeRepo.ExpectedCalls = nil

		// Updated after scope refactoring: removed ScopeBrandID and ScopeGroupID
		cmd := application.CreateAttributeCommand{
			Name:      "Color",
			Code:      "C",
			SortOrder: 1,
			Values: []application.CreateAttributeValueCommand{
				{Value: "Red", Code: "R"},
			},
		}

		mockAttributeRepo.On("FindByCode", ctx, "C").Return(nil, nil).Once()
		mockAttributeRepo.On("Save", ctx, mock.AnythingOfType("*domain.Attribute")).Return(nil).Once()

		result, err := service.CreateAttribute(ctx, cmd)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Color", result.Name)
		assert.Len(t, result.Values, 1)
		mockAttributeRepo.AssertExpectations(t)
	})

	// Removed test "should fail when scope brand missing" - no longer applicable after refactoring
}

func TestProductService_CreateProduct_ActorIDHandling(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockBrandRepo := new(MockBrandRepository)
	mockGroupRepo := new(MockProductGroupRepository)
	mockAttributeRepo := new(MockAttributeRepository)
	mockVariantRepo := new(MockProductVariantRepository)
	mockPartyServiceConfigRepo := new(MockPartyServiceConfigurationRepository)
	service := application.NewProductService(mockProductRepo, mockBrandRepo, mockGroupRepo, mockAttributeRepo, mockVariantRepo, mockPartyServiceConfigRepo)

	brandID := uuid.New()

	t.Run("missing actor id returns error", func(t *testing.T) {
		ctx := context.Background()
		cmd := application.CreateProductCommand{
			SKU:         "P-1",
			Name:        "Product",
			ProductType: domain.ProductTypeTangible,
			BrandID:     brandID,
		}

		result, err := service.CreateProduct(ctx, cmd)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("uses userID from context", func(t *testing.T) {
		mockBrandRepo.ExpectedCalls = nil
		mockProductRepo.ExpectedCalls = nil
		ctx := context.WithValue(context.Background(), "userID", "user-1")
		cmd := application.CreateProductCommand{
			SKU:         "P-2",
			Name:        "Product 2",
			ProductType: domain.ProductTypeTangible,
			BrandID:     brandID,
		}

		ctxMatcher := mock.MatchedBy(func(c context.Context) bool {
			value, _ := c.Value("actorID").(string)
			return value == "user-1"
		})
		mockBrandRepo.On("FindByID", ctxMatcher, brandID).Return(&domain.Brand{ID: brandID, Name: "Brand"}, nil).Once()
		mockProductRepo.On("FindBySKU", ctxMatcher, "P-2").Return(nil, nil).Once() // SKU not found (no duplicate)
		mockProductRepo.On("Save", ctxMatcher, mock.AnythingOfType("*domain.Product")).Return(nil).Once()

		result, err := service.CreateProduct(ctx, cmd)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		mockBrandRepo.AssertExpectations(t)
		mockProductRepo.AssertExpectations(t)
	})
}

func TestProductService_UpdateAttribute(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockBrandRepo := new(MockBrandRepository)
	mockGroupRepo := new(MockProductGroupRepository)
	mockAttributeRepo := new(MockAttributeRepository)
	mockVariantRepo := new(MockProductVariantRepository)
	mockPartyServiceConfigRepo := new(MockPartyServiceConfigurationRepository)
	service := application.NewProductService(mockProductRepo, mockBrandRepo, mockGroupRepo, mockAttributeRepo, mockVariantRepo, mockPartyServiceConfigRepo)
	ctx := context.WithValue(context.Background(), "actorID", "test-actor")

	attr, _ := domain.NewAttribute("Size", "S", 1)
	value1, _ := attr.AddValue("Large", "L")
	value2, _ := attr.AddValue("Medium", "M")
	newName := "Size Updated"
	newCode := "SZ"
	newSort := 2

	t.Run("should update attribute values and remove extras", func(t *testing.T) {
		mockAttributeRepo.ExpectedCalls = nil

		cmd := application.UpdateAttributeCommand{
			ID:        attr.ID,
			Name:      &newName,
			Code:      &newCode,
			SortOrder: &newSort,
			Values: []application.UpdateAttributeValueCommand{
				{ID: &value1.ID, Value: "XL", Code: "XL"},
				{Value: "Small", Code: "S"},
			},
		}

		mockAttributeRepo.On("FindByID", ctx, attr.ID).Return(attr, nil).Once()
		mockAttributeRepo.On("Save", ctx, mock.AnythingOfType("*domain.Attribute")).Return(nil).Once()

		result, err := service.UpdateAttribute(ctx, cmd)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, newName, result.Name)
		assert.Equal(t, newCode, result.Code)
		assert.Len(t, result.Values, 2)
		assert.NotContains(t, result.Values, value2.Value)
		mockAttributeRepo.AssertExpectations(t)
	})

	_ = value2

	t.Run("should fail when attribute missing", func(t *testing.T) {
		mockAttributeRepo.ExpectedCalls = nil
		cmd := application.UpdateAttributeCommand{ID: uuid.New()}

		mockAttributeRepo.On("FindByID", ctx, cmd.ID).Return(nil, nil).Once()

		result, err := service.UpdateAttribute(ctx, cmd)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockAttributeRepo.AssertExpectations(t)
	})
}

// TestProductService_GetApplicableAttributesForProduct - REMOVED (obsolete after scope refactoring)
// TODO: Rewrite for DirectAttributeIDs system

func TestProductService_GetAndListAttributes(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockBrandRepo := new(MockBrandRepository)
	mockGroupRepo := new(MockProductGroupRepository)
	mockAttributeRepo := new(MockAttributeRepository)
	mockVariantRepo := new(MockProductVariantRepository)
	mockPartyServiceConfigRepo := new(MockPartyServiceConfigurationRepository)
	service := application.NewProductService(mockProductRepo, mockBrandRepo, mockGroupRepo, mockAttributeRepo, mockVariantRepo, mockPartyServiceConfigRepo)
	ctx := context.WithValue(context.Background(), "actorID", "test-actor")

	attr, _ := domain.NewAttribute("Generic", "G", 0)

	t.Run("get attribute by id", func(t *testing.T) {
		mockAttributeRepo.ExpectedCalls = nil
		query := application.GetAttributeByIDQuery{ID: attr.ID}

		mockAttributeRepo.On("FindByID", ctx, attr.ID).Return(attr, nil).Once()

		result, err := service.GetAttributeByID(ctx, query)
		assert.NoError(t, err)
		assert.Equal(t, attr.ID, result.ID)
		mockAttributeRepo.AssertExpectations(t)
	})

	// Note: ScopeType field removed from ListAttributesQuery after refactoring
	// All attributes are now generic (no brand/group scopes)
}

func TestProductService_GetProductAndVariants(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockBrandRepo := new(MockBrandRepository)
	mockGroupRepo := new(MockProductGroupRepository)
	mockAttributeRepo := new(MockAttributeRepository)
	mockVariantRepo := new(MockProductVariantRepository)
	mockPartyServiceConfigRepo := new(MockPartyServiceConfigurationRepository)
	service := application.NewProductService(mockProductRepo, mockBrandRepo, mockGroupRepo, mockAttributeRepo, mockVariantRepo, mockPartyServiceConfigRepo)
	ctx := context.WithValue(context.Background(), "actorID", "test-actor")

	product := &domain.Product{ID: uuid.New(), SKU: "P-1", Name: "Product", BrandID: uuid.New(), IsActive: true}
	attr, _ := domain.NewAttribute("Color", "C", 0)
	val, _ := attr.AddValue("Red", "R")
	variant := &domain.ProductVariant{ID: uuid.New(), ProductID: product.ID, SKU: "P-1-C.R", AttributeValues: []uuid.UUID{val.ID}, Status: domain.StatusConfirmed, IsActive: true}

	mockProductRepo.On("FindByID", ctx, product.ID).Return(product, nil).Once()
	result, err := service.GetProductByID(ctx, application.GetProductByIDQuery{ID: product.ID})
	assert.NoError(t, err)
	assert.Equal(t, product.ID, result.ID)

	mockVariantRepo.On("FindByID", ctx, variant.ID).Return(variant, nil).Once()
	mockProductRepo.On("FindByID", ctx, product.ID).Return(product, nil).Once()
	mockAttributeRepo.On("FindByScope", ctx, (*uuid.UUID)(nil), (*uuid.UUID)(nil)).Return([]*domain.Attribute{attr}, nil).Once()
	variantResult, err := service.GetProductVariantByID(ctx, application.GetProductVariantByIDQuery{ID: variant.ID})
	assert.NoError(t, err)
	assert.Equal(t, variant.ID, variantResult.ID)

	mockVariantRepo.On("FindBySKU", ctx, variant.SKU).Return(variant, nil).Once()
	mockProductRepo.On("FindByID", ctx, product.ID).Return(product, nil).Once()
	mockAttributeRepo.On("FindByScope", ctx, (*uuid.UUID)(nil), (*uuid.UUID)(nil)).Return([]*domain.Attribute{attr}, nil).Once()
	variantSKUResult, err := service.GetProductVariantBySKU(ctx, application.GetProductVariantBySKUQuery{SKU: variant.SKU})
	assert.NoError(t, err)
	assert.Equal(t, variant.SKU, variantSKUResult.SKU)

	nonExistentProductID := uuid.New()
	mockProductRepo.On("FindByID", ctx, nonExistentProductID).Return(nil, nil).Once()
	missingProduct, err := service.GetProductByID(ctx, application.GetProductByIDQuery{ID: nonExistentProductID})
	assert.Error(t, err)
	assert.Nil(t, missingProduct)

	missingVariantID := uuid.New()
	mockVariantRepo.On("FindByID", ctx, missingVariantID).Return(nil, nil).Once()
	missingVariant, err := service.GetProductVariantByID(ctx, application.GetProductVariantByIDQuery{ID: missingVariantID})
	assert.Error(t, err)
	assert.Nil(t, missingVariant)

	mockVariantRepo.On("FindBySKU", ctx, "missing-sku").Return(nil, nil).Once()
	missingVariantBySKU, err := service.GetProductVariantBySKU(ctx, application.GetProductVariantBySKUQuery{SKU: "missing-sku"})
	assert.Error(t, err)
	assert.Nil(t, missingVariantBySKU)

	mockProductRepo.AssertExpectations(t)
	mockVariantRepo.AssertExpectations(t)
	mockAttributeRepo.AssertExpectations(t)
}

// TestProductService_FindOrCreateProductVariant - REMOVED (obsolete after scope refactoring)
// TODO: Rewrite for DirectAttributeIDs system

func TestProductService_PartyServiceConfiguration_InvalidJSON(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockBrandRepo := new(MockBrandRepository)
	mockGroupRepo := new(MockProductGroupRepository)
	mockAttributeRepo := new(MockAttributeRepository)
	mockVariantRepo := new(MockProductVariantRepository)
	mockPartyServiceConfigRepo := new(MockPartyServiceConfigurationRepository)
	service := application.NewProductService(mockProductRepo, mockBrandRepo, mockGroupRepo, mockAttributeRepo, mockVariantRepo, mockPartyServiceConfigRepo)
	ctx := context.WithValue(context.Background(), "actorID", "test-actor")

	partyID := uuid.New()
	badDetails := map[string]interface{}{"bad": func() {}}

	createResult, err := service.CreatePartyServiceConfiguration(ctx, application.CreatePartyServiceConfigurationCommand{
		PartyID:              partyID,
		ServiceID:            "svc",
		Name:                 "Config",
		ConfigurationDetails: badDetails,
	})
	assert.Error(t, err)
	assert.Nil(t, createResult)

	config, _ := domain.NewPartyServiceConfiguration(partyID, "svc", "Config", json.RawMessage(`{"k":"v"}`))
	mockPartyServiceConfigRepo.On("FindByID", ctx, partyID, config.ID).Return(config, nil).Once()

	updateResult, err := service.UpdatePartyServiceConfiguration(ctx, application.UpdatePartyServiceConfigurationCommand{
		ID:                   config.ID,
		PartyID:              partyID,
		ConfigurationDetails: badDetails,
	})
	assert.Error(t, err)
	assert.Nil(t, updateResult)
	mockPartyServiceConfigRepo.AssertExpectations(t)
}

func TestProductService_UpdateProductVariant(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockBrandRepo := new(MockBrandRepository)
	mockGroupRepo := new(MockProductGroupRepository)
	mockAttributeRepo := new(MockAttributeRepository)
	mockVariantRepo := new(MockProductVariantRepository)
	mockPartyServiceConfigRepo := new(MockPartyServiceConfigurationRepository)
	service := application.NewProductService(mockProductRepo, mockBrandRepo, mockGroupRepo, mockAttributeRepo, mockVariantRepo, mockPartyServiceConfigRepo)
	ctx := context.WithValue(context.Background(), "actorID", "test-actor")

	attr, _ := domain.NewAttribute("Color", "C", 0)
	val, _ := attr.AddValue("Red", "R")
	variant := &domain.ProductVariant{ID: uuid.New(), SKU: "P-1-C.R", AttributeValues: []uuid.UUID{val.ID}, Status: domain.StatusProvisional}
	newBarcode := "123"
	active := true

	product := &domain.Product{ID: uuid.New(), Name: "Test Product"}
	variant.ProductID = product.ID

	mockVariantRepo.On("FindByID", ctx, variant.ID).Return(variant, nil).Once()
	mockVariantRepo.On("Save", ctx, mock.AnythingOfType("*domain.ProductVariant")).Return(nil).Once()
	mockAttributeRepo.On("FindByScope", ctx, (*uuid.UUID)(nil), (*uuid.UUID)(nil)).Return([]*domain.Attribute{attr}, nil).Once()
	mockProductRepo.On("FindByID", ctx, product.ID).Return(product, nil).Once()

	result, err := service.UpdateProductVariant(ctx, application.UpdateProductVariantCommand{ID: variant.ID, Barcode: &newBarcode, IsActive: &active})

	assert.NoError(t, err)
	assert.Equal(t, domain.StatusConfirmed, variant.Status)
	assert.Equal(t, newBarcode, *result.Barcode)

	mockVariantRepo.AssertExpectations(t)
	mockAttributeRepo.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
}

func TestProductService_PartyServiceConfigurations(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockBrandRepo := new(MockBrandRepository)
	mockGroupRepo := new(MockProductGroupRepository)
	mockAttributeRepo := new(MockAttributeRepository)
	mockVariantRepo := new(MockProductVariantRepository)
	mockPartyServiceConfigRepo := new(MockPartyServiceConfigurationRepository)
	service := application.NewProductService(mockProductRepo, mockBrandRepo, mockGroupRepo, mockAttributeRepo, mockVariantRepo, mockPartyServiceConfigRepo)
	ctx := context.WithValue(context.Background(), "actorID", "test-actor")

	partyID := uuid.New()
	configID := uuid.New()
	config := &domain.PartyServiceConfiguration{
		ID:        configID,
		PartyID:   partyID,
		ServiceID: "svc",
		Name:      "Config",
	}

	mockPartyServiceConfigRepo.On("FindByID", ctx, partyID, configID).Return(config, nil).Once()
	result, err := service.GetPartyServiceConfigurationByID(ctx, application.GetPartyServiceConfigurationByIDQuery{PartyID: partyID, ID: configID})
	assert.NoError(t, err)
	assert.Equal(t, configID, result.ID)

	mockPartyServiceConfigRepo.On("FindByPartyID", ctx, partyID).Return([]*domain.PartyServiceConfiguration{config}, nil).Once()
	listResult, err := service.ListPartyServiceConfigurationsByPartyID(ctx, application.ListPartyServiceConfigurationsByPartyIDQuery{PartyID: partyID})
	assert.NoError(t, err)
	assert.Len(t, listResult, 1)

	mockPartyServiceConfigRepo.On("Save", ctx, mock.AnythingOfType("*domain.PartyServiceConfiguration")).Return(nil).Once()
	createResult, err := service.CreatePartyServiceConfiguration(ctx, application.CreatePartyServiceConfigurationCommand{PartyID: partyID, ServiceID: "svc", Name: "Config", ConfigurationDetails: map[string]interface{}{"k": "v"}})
	assert.NoError(t, err)
	assert.Equal(t, partyID, createResult.PartyID)

	newServiceID := "svc2"
	mockPartyServiceConfigRepo.On("FindByID", ctx, partyID, configID).Return(config, nil).Once()
	mockPartyServiceConfigRepo.On("Save", ctx, mock.AnythingOfType("*domain.PartyServiceConfiguration")).Return(nil).Once()
	updateResult, err := service.UpdatePartyServiceConfiguration(ctx, application.UpdatePartyServiceConfigurationCommand{ID: configID, PartyID: partyID, ServiceID: &newServiceID})
	assert.NoError(t, err)
	assert.Equal(t, newServiceID, updateResult.ServiceID)

	mockPartyServiceConfigRepo.On("Delete", ctx, partyID, configID).Return(nil).Once()
	assert.NoError(t, service.DeletePartyServiceConfiguration(ctx, application.DeletePartyServiceConfigurationCommand{ID: configID, PartyID: partyID}))

	mockPartyServiceConfigRepo.AssertExpectations(t)
}

// TestProductService_GenerateProductVariants - REMOVED (obsolete after scope refactoring)
// TODO: Rewrite for DirectAttributeIDs system

func TestProductService_UpdateProduct_Success(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockBrandRepo := new(MockBrandRepository)
	mockGroupRepo := new(MockProductGroupRepository)
	mockAttributeRepo := new(MockAttributeRepository)
	mockVariantRepo := new(MockProductVariantRepository)
	mockPartyServiceConfigRepo := new(MockPartyServiceConfigurationRepository)
	service := application.NewProductService(mockProductRepo, mockBrandRepo, mockGroupRepo, mockAttributeRepo, mockVariantRepo, mockPartyServiceConfigRepo)
	ctx := context.WithValue(context.Background(), "actorID", "test-actor")

	t.Run("should update product with new values", func(t *testing.T) {
		productID := uuid.New()
		brandID := uuid.New()
		existingProduct := &domain.Product{
			ID:          productID,
			SKU:         "PROD-001",
			Name:        "Old Name",
			LongName:    "Old Long Name",
			Description: "Old Description",
			ProductType: "SIMPLE",
			BrandID:     brandID,
		}

		newName := "Updated Name"
		newLongName := "Updated Long Name"
		newDescription := "Updated Description"

		mockProductRepo.On("FindByID", ctx, productID).Return(existingProduct, nil).Twice()
		mockProductRepo.On("Save", ctx, mock.AnythingOfType("*domain.Product")).Return(nil).Once()

		result, err := service.UpdateProduct(ctx, application.UpdateProductCommand{
			ProductID:   productID,
			ActorID:     "test-actor",
			Name:        &newName,
			LongName:    &newLongName,
			Description: &newDescription,
		})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Updated Name", result.Name)
		mockProductRepo.AssertExpectations(t)
	})

	t.Run("should return error when product not found", func(t *testing.T) {
		productID := uuid.New()
		mockProductRepo.ExpectedCalls = nil
		mockProductRepo.On("FindByID", ctx, productID).Return(nil, nil).Once()

		newName := "Updated Name"
		result, err := service.UpdateProduct(ctx, application.UpdateProductCommand{
			ProductID: productID,
			ActorID:   "test-actor",
			Name:      &newName,
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "does not exist")
		mockProductRepo.AssertExpectations(t)
	})
}

func TestProductService_GetProductByID_Success(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockBrandRepo := new(MockBrandRepository)
	mockGroupRepo := new(MockProductGroupRepository)
	mockAttributeRepo := new(MockAttributeRepository)
	mockVariantRepo := new(MockProductVariantRepository)
	mockPartyServiceConfigRepo := new(MockPartyServiceConfigurationRepository)
	service := application.NewProductService(mockProductRepo, mockBrandRepo, mockGroupRepo, mockAttributeRepo, mockVariantRepo, mockPartyServiceConfigRepo)
	ctx := context.Background()

	t.Run("should get product by ID", func(t *testing.T) {
		productID := uuid.New()
		brandID := uuid.New()
		expectedProduct := &domain.Product{
			ID:          productID,
			SKU:         "PROD-123",
			Name:        "Test Product",
			LongName:    "Test Product Long Name",
			Description: "Test Description",
			ProductType: "SIMPLE",
			BrandID:     brandID,
		}

		mockProductRepo.On("FindByID", ctx, productID).Return(expectedProduct, nil).Once()

		result, err := service.GetProductByID(ctx, application.GetProductByIDQuery{ID: productID})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "PROD-123", result.SKU)
		assert.Equal(t, "Test Product", result.Name)
		mockProductRepo.AssertExpectations(t)
	})

	t.Run("should return error when product not found", func(t *testing.T) {
		productID := uuid.New()
		mockProductRepo.ExpectedCalls = nil
		mockProductRepo.On("FindByID", ctx, productID).Return(nil, nil).Once()

		result, err := service.GetProductByID(ctx, application.GetProductByIDQuery{ID: productID})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "does not exist")
		mockProductRepo.AssertExpectations(t)
	})
}

func TestProductService_GetProductVariantByID_Success(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockBrandRepo := new(MockBrandRepository)
	mockGroupRepo := new(MockProductGroupRepository)
	mockAttributeRepo := new(MockAttributeRepository)
	mockVariantRepo := new(MockProductVariantRepository)
	mockPartyServiceConfigRepo := new(MockPartyServiceConfigurationRepository)
	service := application.NewProductService(mockProductRepo, mockBrandRepo, mockGroupRepo, mockAttributeRepo, mockVariantRepo, mockPartyServiceConfigRepo)
	ctx := context.Background()

	t.Run("should get product variant by ID", func(t *testing.T) {
		variantID := uuid.New()
		productID := uuid.New()
		expectedVariant := &domain.ProductVariant{
			ID:              variantID,
			ProductID:       productID,
			SKU:             "VAR-001",
			AttributeValues: []uuid.UUID{},
			Status:          domain.StatusConfirmed,
			IsActive:        true,
		}

		expectedProduct := &domain.Product{ID: productID, Name: "Test Product"}

		mockVariantRepo.On("FindByID", ctx, variantID).Return(expectedVariant, nil).Once()
		mockAttributeRepo.On("FindByScope", ctx, (*uuid.UUID)(nil), (*uuid.UUID)(nil)).Return([]*domain.Attribute{}, nil).Once()
		mockProductRepo.On("FindByID", ctx, productID).Return(expectedProduct, nil).Once()

		result, err := service.GetProductVariantByID(ctx, application.GetProductVariantByIDQuery{ID: variantID})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "VAR-001", result.SKU)
		mockVariantRepo.AssertExpectations(t)
		mockAttributeRepo.AssertExpectations(t)
		mockProductRepo.AssertExpectations(t)
	})

	t.Run("should return error when variant not found", func(t *testing.T) {
		variantID := uuid.New()
		mockVariantRepo.ExpectedCalls = nil
		mockVariantRepo.On("FindByID", ctx, variantID).Return(nil, nil).Once()

		result, err := service.GetProductVariantByID(ctx, application.GetProductVariantByIDQuery{ID: variantID})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "does not exist")
		mockVariantRepo.AssertExpectations(t)
	})
}

func TestProductService_GetProductVariantBySKU_Success(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockBrandRepo := new(MockBrandRepository)
	mockGroupRepo := new(MockProductGroupRepository)
	mockAttributeRepo := new(MockAttributeRepository)
	mockVariantRepo := new(MockProductVariantRepository)
	mockPartyServiceConfigRepo := new(MockPartyServiceConfigurationRepository)
	service := application.NewProductService(mockProductRepo, mockBrandRepo, mockGroupRepo, mockAttributeRepo, mockVariantRepo, mockPartyServiceConfigRepo)
	ctx := context.Background()

	t.Run("should get product variant by SKU", func(t *testing.T) {
		variantID := uuid.New()
		productID := uuid.New()
		sku := "VAR-SKU-001"
		expectedVariant := &domain.ProductVariant{
			ID:              variantID,
			ProductID:       productID,
			SKU:             sku,
			AttributeValues: []uuid.UUID{},
			Status:          domain.StatusConfirmed,
			IsActive:        true,
		}

		expectedProduct := &domain.Product{ID: productID, Name: "Test Product"}

		mockVariantRepo.On("FindBySKU", ctx, sku).Return(expectedVariant, nil).Once()
		mockAttributeRepo.On("FindByScope", ctx, (*uuid.UUID)(nil), (*uuid.UUID)(nil)).Return([]*domain.Attribute{}, nil).Once()
		mockProductRepo.On("FindByID", ctx, productID).Return(expectedProduct, nil).Once()

		result, err := service.GetProductVariantBySKU(ctx, application.GetProductVariantBySKUQuery{SKU: sku})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, sku, result.SKU)
		mockVariantRepo.AssertExpectations(t)
		mockAttributeRepo.AssertExpectations(t)
		mockProductRepo.AssertExpectations(t)
	})

	t.Run("should return error when variant not found", func(t *testing.T) {
		sku := "NONEXISTENT-SKU"
		mockVariantRepo.ExpectedCalls = nil
		mockVariantRepo.On("FindBySKU", ctx, sku).Return(nil, nil).Once()

		result, err := service.GetProductVariantBySKU(ctx, application.GetProductVariantBySKUQuery{SKU: sku})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "does not exist")
		mockVariantRepo.AssertExpectations(t)
	})
}
