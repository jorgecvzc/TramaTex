package application_test

import (
	"context"
	"testing"
	"time"

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
	ctx := context.Background()

	brandID := uuid.New()
	groupID := uuid.New()

	t.Run("should create attribute with values", func(t *testing.T) {
		mockBrandRepo.ExpectedCalls = nil
		mockGroupRepo.ExpectedCalls = nil
		mockAttributeRepo.ExpectedCalls = nil

		cmd := application.CreateAttributeCommand{
			Name:         "Color",
			Code:         "C",
			SortOrder:    1,
			ScopeBrandID: &brandID,
			ScopeGroupID: &groupID,
			Values: []application.CreateAttributeValueCommand{
				{Value: "Red", Code: "R"},
			},
		}

		mockBrandRepo.On("FindByID", ctx, brandID).Return(&domain.Brand{ID: brandID}, nil).Once()
		mockGroupRepo.On("FindByID", ctx, groupID).Return(&domain.ProductGroup{ID: groupID}, nil).Once()
		mockAttributeRepo.On("Save", ctx, mock.AnythingOfType("*domain.Attribute")).Return(nil).Once()

		result, err := service.CreateAttribute(ctx, cmd)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Color", result.Name)
		assert.Len(t, result.Values, 1)
		mockBrandRepo.AssertExpectations(t)
		mockGroupRepo.AssertExpectations(t)
		mockAttributeRepo.AssertExpectations(t)
	})

	t.Run("should fail when scope brand missing", func(t *testing.T) {
		mockBrandRepo.ExpectedCalls = nil

		cmd := application.CreateAttributeCommand{
			Name:         "Color",
			Code:         "C",
			SortOrder:    1,
			ScopeBrandID: &brandID,
		}

		mockBrandRepo.On("FindByID", ctx, brandID).Return(nil, nil).Once()

		result, err := service.CreateAttribute(ctx, cmd)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockBrandRepo.AssertExpectations(t)
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
	ctx := context.Background()

	attr, _ := domain.NewAttribute("Size", "S", 1, nil, nil)
	value1, _ := attr.AddValue("Large", "L")
	value2, _ := attr.AddValue("Medium", "M")
	newName := "Size Updated"
	newCode := "SZ"
	newSort := 2
	brandID := uuid.New()

	t.Run("should update attribute values and remove extras", func(t *testing.T) {
		mockAttributeRepo.ExpectedCalls = nil
		mockBrandRepo.ExpectedCalls = nil

		cmd := application.UpdateAttributeCommand{
			ID:           attr.ID,
			Name:         &newName,
			Code:         &newCode,
			SortOrder:    &newSort,
			ScopeBrandID: &brandID,
			Values: []application.UpdateAttributeValueCommand{
				{ID: &value1.ID, Value: "XL", Code: "XL"},
				{Value: "Small", Code: "S"},
			},
		}

		mockAttributeRepo.On("FindByID", ctx, attr.ID).Return(attr, nil).Once()
		mockBrandRepo.On("FindByID", ctx, brandID).Return(&domain.Brand{ID: brandID}, nil).Once()
		mockAttributeRepo.On("Save", ctx, mock.AnythingOfType("*domain.Attribute")).Return(nil).Once()

		result, err := service.UpdateAttribute(ctx, cmd)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, newName, result.Name)
		assert.Equal(t, newCode, result.AttributeName)
		assert.Len(t, result.Values, 2)
		assert.NotContains(t, result.Values, value2.Value)
		mockAttributeRepo.AssertExpectations(t)
		mockBrandRepo.AssertExpectations(t)
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

func TestProductService_GetApplicableAttributesForProduct(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockBrandRepo := new(MockBrandRepository)
	mockGroupRepo := new(MockProductGroupRepository)
	mockAttributeRepo := new(MockAttributeRepository)
	mockVariantRepo := new(MockProductVariantRepository)
	mockPartyServiceConfigRepo := new(MockPartyServiceConfigurationRepository)
	service := application.NewProductService(mockProductRepo, mockBrandRepo, mockGroupRepo, mockAttributeRepo, mockVariantRepo, mockPartyServiceConfigRepo)
	ctx := context.Background()

	brandID := uuid.New()
	groupID := uuid.New()
	directAttrID := uuid.New()
	product := &domain.Product{
		ID:                 uuid.New(),
		SKU:                "P-1",
		Name:               "Product",
		BrandID:            brandID,
		GroupIDs:           []uuid.UUID{groupID},
		DirectAttributeIDs: []uuid.UUID{directAttrID},
		IsActive:           true,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	genericAttr, _ := domain.NewAttribute("Generic", "G", 3, nil, nil)
	brandAttr, _ := domain.NewAttribute("Brand", "B", 2, &brandID, nil)
	groupAttr, _ := domain.NewAttribute("Group", "GR", 4, nil, &groupID)
	groupBrandAttr, _ := domain.NewAttribute("GroupBrand", "GB", 1, &brandID, &groupID)
	directAttr, _ := domain.NewAttribute("Direct", "D", 0, nil, nil)
	directAttr.ID = directAttrID

	mockProductRepo.On("FindByID", ctx, product.ID).Return(product, nil).Once()
	mockAttributeRepo.On("FindByScope", ctx, (*uuid.UUID)(nil), (*uuid.UUID)(nil)).Return([]*domain.Attribute{genericAttr}, nil).Once()
	mockAttributeRepo.On("FindByScope", ctx, &brandID, (*uuid.UUID)(nil)).Return([]*domain.Attribute{brandAttr}, nil).Once()
	mockAttributeRepo.On("FindByScope", ctx, (*uuid.UUID)(nil), &groupID).Return([]*domain.Attribute{groupAttr}, nil).Once()
	mockAttributeRepo.On("FindByScope", ctx, &brandID, &groupID).Return([]*domain.Attribute{groupBrandAttr}, nil).Once()
	mockAttributeRepo.On("FindByID", ctx, directAttrID).Return(directAttr, nil).Once()

	result, err := service.GetApplicableAttributesForProduct(ctx, product.ID)

	assert.NoError(t, err)
	assert.Len(t, result, 5)
	mockProductRepo.AssertExpectations(t)
	mockAttributeRepo.AssertExpectations(t)
}

func TestProductService_GetAndListAttributes(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockBrandRepo := new(MockBrandRepository)
	mockGroupRepo := new(MockProductGroupRepository)
	mockAttributeRepo := new(MockAttributeRepository)
	mockVariantRepo := new(MockProductVariantRepository)
	mockPartyServiceConfigRepo := new(MockPartyServiceConfigurationRepository)
	service := application.NewProductService(mockProductRepo, mockBrandRepo, mockGroupRepo, mockAttributeRepo, mockVariantRepo, mockPartyServiceConfigRepo)
	ctx := context.Background()

	attr, _ := domain.NewAttribute("Generic", "G", 0, nil, nil)

	t.Run("get attribute by id", func(t *testing.T) {
		mockAttributeRepo.ExpectedCalls = nil
		query := application.GetAttributeByIDQuery{ID: attr.ID}

		mockAttributeRepo.On("FindByID", ctx, attr.ID).Return(attr, nil).Once()

		result, err := service.GetAttributeByID(ctx, query)
		assert.NoError(t, err)
		assert.Equal(t, attr.ID, result.ID)
		mockAttributeRepo.AssertExpectations(t)
	})

	t.Run("list attributes with scope filter", func(t *testing.T) {
		mockAttributeRepo.ExpectedCalls = nil
		scopeType := "GENERIC"
		query := application.ListAttributesQuery{ScopeType: &scopeType}

		mockAttributeRepo.On("FindByScope", ctx, (*uuid.UUID)(nil), (*uuid.UUID)(nil)).Return([]*domain.Attribute{attr}, nil).Once()

		result, err := service.ListAttributes(ctx, query)
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		mockAttributeRepo.AssertExpectations(t)
	})

	t.Run("list attributes invalid scope", func(t *testing.T) {
		scopeType := "INVALID"
		query := application.ListAttributesQuery{ScopeType: &scopeType}

		result, err := service.ListAttributes(ctx, query)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestProductService_GetProductAndVariants(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockBrandRepo := new(MockBrandRepository)
	mockGroupRepo := new(MockProductGroupRepository)
	mockAttributeRepo := new(MockAttributeRepository)
	mockVariantRepo := new(MockProductVariantRepository)
	mockPartyServiceConfigRepo := new(MockPartyServiceConfigurationRepository)
	service := application.NewProductService(mockProductRepo, mockBrandRepo, mockGroupRepo, mockAttributeRepo, mockVariantRepo, mockPartyServiceConfigRepo)
	ctx := context.Background()

	product := &domain.Product{ID: uuid.New(), SKU: "P-1", Name: "Product", BrandID: uuid.New(), IsActive: true}
	attr, _ := domain.NewAttribute("Color", "C", 0, nil, nil)
	val, _ := attr.AddValue("Red", "R")
	variant := &domain.ProductVariant{ID: uuid.New(), ProductID: product.ID, SKU: "P-1-C.R", AttributeValues: []uuid.UUID{val.ID}, Status: domain.StatusConfirmed, IsActive: true}

	mockProductRepo.On("FindByID", ctx, product.ID).Return(product, nil).Once()
	result, err := service.GetProductByID(ctx, application.GetProductByIDQuery{ID: product.ID})
	assert.NoError(t, err)
	assert.Equal(t, product.ID, result.ID)

	mockVariantRepo.On("FindByID", ctx, variant.ID).Return(variant, nil).Once()
	mockAttributeRepo.On("FindByScope", ctx, (*uuid.UUID)(nil), (*uuid.UUID)(nil)).Return([]*domain.Attribute{attr}, nil).Once()
	variantResult, err := service.GetProductVariantByID(ctx, application.GetProductVariantByIDQuery{ID: variant.ID})
	assert.NoError(t, err)
	assert.Equal(t, variant.ID, variantResult.ID)

	mockVariantRepo.On("FindBySKU", ctx, variant.SKU).Return(variant, nil).Once()
	mockAttributeRepo.On("FindByScope", ctx, (*uuid.UUID)(nil), (*uuid.UUID)(nil)).Return([]*domain.Attribute{attr}, nil).Once()
	variantSKUResult, err := service.GetProductVariantBySKU(ctx, application.GetProductVariantBySKUQuery{SKU: variant.SKU})
	assert.NoError(t, err)
	assert.Equal(t, variant.SKU, variantSKUResult.SKU)

	mockProductRepo.AssertExpectations(t)
	mockVariantRepo.AssertExpectations(t)
	mockAttributeRepo.AssertExpectations(t)
}

func TestProductService_FindOrCreateProductVariant(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockBrandRepo := new(MockBrandRepository)
	mockGroupRepo := new(MockProductGroupRepository)
	mockAttributeRepo := new(MockAttributeRepository)
	mockVariantRepo := new(MockProductVariantRepository)
	mockPartyServiceConfigRepo := new(MockPartyServiceConfigurationRepository)
	service := application.NewProductService(mockProductRepo, mockBrandRepo, mockGroupRepo, mockAttributeRepo, mockVariantRepo, mockPartyServiceConfigRepo)
	ctx := context.Background()

	brandID := uuid.New()
	product := &domain.Product{ID: uuid.New(), SKU: "P-1", Name: "Product", BrandID: brandID}
	attr, _ := domain.NewAttribute("Color", "C", 1, nil, nil)
	val, _ := attr.AddValue("Red", "R")
	option := map[string]string{attr.Code: val.Value}

	mockProductRepo.On("FindByID", ctx, product.ID).Return(product, nil).Times(4)
	mockAttributeRepo.On("FindByScope", ctx, (*uuid.UUID)(nil), (*uuid.UUID)(nil)).Return([]*domain.Attribute{attr}, nil).Twice()
	mockAttributeRepo.On("FindByScope", ctx, &brandID, (*uuid.UUID)(nil)).Return([]*domain.Attribute{}, nil).Twice()
	mockAttributeRepo.On("FindByScope", ctx, (*uuid.UUID)(nil), mock.AnythingOfType("*uuid.UUID")).Return([]*domain.Attribute{}, nil).Maybe()
	mockAttributeRepo.On("FindByScope", ctx, &brandID, mock.AnythingOfType("*uuid.UUID")).Return([]*domain.Attribute{}, nil).Maybe()
	mockAttributeRepo.On("FindByID", ctx, attr.ID).Return(attr, nil).Twice()

	existingVariant := &domain.ProductVariant{ID: uuid.New(), ProductID: product.ID, SKU: "P-1-C.R", AttributeValues: []uuid.UUID{val.ID}, Status: domain.StatusConfirmed}
	mockVariantRepo.On("FindByProductIDAndAttributeValues", ctx, product.ID, mock.Anything).Return(existingVariant, nil).Once()

	result, err := service.FindOrCreateProductVariant(ctx, application.FindOrCreateProductVariantCommand{ProductID: product.ID, OptionConfiguration: option})
	assert.NoError(t, err)
	assert.Equal(t, existingVariant.ID, result.ID)

	mockVariantRepo.On("FindByProductIDAndAttributeValues", ctx, product.ID, mock.Anything).Return(nil, nil).Once()
	mockVariantRepo.On("Save", ctx, mock.AnythingOfType("*domain.ProductVariant")).Return(nil).Once()

	result, err = service.FindOrCreateProductVariant(ctx, application.FindOrCreateProductVariantCommand{ProductID: product.ID, OptionConfiguration: option})
	assert.NoError(t, err)
	assert.Equal(t, product.ID, result.ProductID)

	mockVariantRepo.AssertExpectations(t)
	mockAttributeRepo.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
}

func TestProductService_UpdateProductVariant(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockBrandRepo := new(MockBrandRepository)
	mockGroupRepo := new(MockProductGroupRepository)
	mockAttributeRepo := new(MockAttributeRepository)
	mockVariantRepo := new(MockProductVariantRepository)
	mockPartyServiceConfigRepo := new(MockPartyServiceConfigurationRepository)
	service := application.NewProductService(mockProductRepo, mockBrandRepo, mockGroupRepo, mockAttributeRepo, mockVariantRepo, mockPartyServiceConfigRepo)
	ctx := context.Background()

	attr, _ := domain.NewAttribute("Color", "C", 0, nil, nil)
	val, _ := attr.AddValue("Red", "R")
	variant := &domain.ProductVariant{ID: uuid.New(), SKU: "P-1-C.R", AttributeValues: []uuid.UUID{val.ID}, Status: domain.StatusProvisional}
	newBarcode := "123"
	active := true

	mockVariantRepo.On("FindByID", ctx, variant.ID).Return(variant, nil).Once()
	mockVariantRepo.On("Save", ctx, mock.AnythingOfType("*domain.ProductVariant")).Return(nil).Once()
	mockAttributeRepo.On("FindByScope", ctx, (*uuid.UUID)(nil), (*uuid.UUID)(nil)).Return([]*domain.Attribute{attr}, nil).Once()

	result, err := service.UpdateProductVariant(ctx, application.UpdateProductVariantCommand{ID: variant.ID, Barcode: &newBarcode, IsActive: &active})

	assert.NoError(t, err)
	assert.Equal(t, domain.StatusConfirmed, variant.Status)
	assert.Equal(t, newBarcode, *result.Barcode)

	mockVariantRepo.AssertExpectations(t)
	mockAttributeRepo.AssertExpectations(t)
}

func TestProductService_PartyServiceConfigurations(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockBrandRepo := new(MockBrandRepository)
	mockGroupRepo := new(MockProductGroupRepository)
	mockAttributeRepo := new(MockAttributeRepository)
	mockVariantRepo := new(MockProductVariantRepository)
	mockPartyServiceConfigRepo := new(MockPartyServiceConfigurationRepository)
	service := application.NewProductService(mockProductRepo, mockBrandRepo, mockGroupRepo, mockAttributeRepo, mockVariantRepo, mockPartyServiceConfigRepo)
	ctx := context.Background()

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

func TestProductService_GenerateProductVariants(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockBrandRepo := new(MockBrandRepository)
	mockGroupRepo := new(MockProductGroupRepository)
	mockAttributeRepo := new(MockAttributeRepository)
	mockVariantRepo := new(MockProductVariantRepository)
	mockPartyServiceConfigRepo := new(MockPartyServiceConfigurationRepository)
	service := application.NewProductService(mockProductRepo, mockBrandRepo, mockGroupRepo, mockAttributeRepo, mockVariantRepo, mockPartyServiceConfigRepo)
	ctx := context.Background()

	brandID := uuid.New()
	product := &domain.Product{ID: uuid.New(), SKU: "P-1", Name: "Product", BrandID: brandID, GroupIDs: []uuid.UUID{uuid.New()}}
	attr, _ := domain.NewAttribute("Color", "C", 1, nil, nil)
	_, _ = attr.AddValue("Red", "R")

	mockProductRepo.On("FindByID", ctx, product.ID).Return(product, nil).Times(2)
	mockAttributeRepo.On("FindByScope", ctx, (*uuid.UUID)(nil), (*uuid.UUID)(nil)).Return([]*domain.Attribute{attr}, nil).Once()
	mockAttributeRepo.On("FindByScope", ctx, &brandID, (*uuid.UUID)(nil)).Return([]*domain.Attribute{}, nil).Once()
	mockAttributeRepo.On("FindByScope", ctx, (*uuid.UUID)(nil), mock.AnythingOfType("*uuid.UUID")).Return([]*domain.Attribute{}, nil).Maybe()
	mockAttributeRepo.On("FindByScope", ctx, &brandID, mock.AnythingOfType("*uuid.UUID")).Return([]*domain.Attribute{}, nil).Maybe()
	mockAttributeRepo.On("FindByID", ctx, attr.ID).Return(attr, nil).Once()

	err := service.GenerateProductVariants(ctx, application.GenerateProductVariantsCommand{ProductID: product.ID})
	assert.NoError(t, err)

	mockProductRepo.AssertExpectations(t)
	mockAttributeRepo.AssertExpectations(t)
}
