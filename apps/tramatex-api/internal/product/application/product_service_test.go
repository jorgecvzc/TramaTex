package application_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/joran-cortez/tramatex/internal/product/application"
	"github.com/joran-cortez/tramatex/internal/product/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProductRepository is a mock implementation of domain.ProductRepository
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) Save(ctx context.Context, product *domain.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	args := m.Called(ctx, id)
	var product *domain.Product
	if args.Get(0) != nil {
		product = args.Get(0).(*domain.Product)
	}
	return product, args.Error(1)
}

func (m *MockProductRepository) FindBySKU(ctx context.Context, sku string) (*domain.Product, error) {
	args := m.Called(ctx, sku)
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductRepository) UpdateSKUs(ctx context.Context, productID uuid.UUID, newSKU string) error {
	args := m.Called(ctx, productID, newSKU)
	return args.Error(0)
}

// MockBrandRepository is a mock implementation of domain.BrandRepository
type MockBrandRepository struct {
	mock.Mock
}

func (m *MockBrandRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Brand, error) {
	args := m.Called(ctx, id)
	var brand *domain.Brand
	if args.Get(0) != nil {
		brand = args.Get(0).(*domain.Brand)
	}
	return brand, args.Error(1)
}

// MockProductGroupRepository is a mock implementation of domain.ProductGroupRepository
type MockProductGroupRepository struct {
	mock.Mock
}

func (m *MockProductGroupRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.ProductGroup, error) {
	args := m.Called(ctx, id)
	var group *domain.ProductGroup
	if args.Get(0) != nil {
		group = args.Get(0).(*domain.ProductGroup)
	}
	return group, args.Error(1)
}

// MockAttributeRepository is a mock implementation of domain.AttributeRepository
type MockAttributeRepository struct {
	mock.Mock
}

func (m *MockAttributeRepository) Save(ctx context.Context, attribute *domain.Attribute) error {
	args := m.Called(ctx, attribute)
	return args.Error(0)
}

func (m *MockAttributeRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Attribute, error) {
	args := m.Called(ctx, id)
	var attribute *domain.Attribute
	if args.Get(0) != nil {
		attribute = args.Get(0).(*domain.Attribute)
	}
	return attribute, args.Error(1)
}

func (m *MockAttributeRepository) FindByIDs(ctx context.Context, ids []uuid.UUID) ([]domain.Attribute, error) {
	args := m.Called(ctx, ids)
	var attributes []domain.Attribute
	if args.Get(0) != nil {
		attributes = args.Get(0).([]domain.Attribute)
	}
	return attributes, args.Error(1)
}

func (m *MockAttributeRepository) FindByScope(ctx context.Context, brandID *uuid.UUID, groupID *uuid.UUID) ([]*domain.Attribute, error) {
	args := m.Called(ctx, brandID, groupID)
	var attributes []*domain.Attribute
	if args.Get(0) != nil {
		attributes = args.Get(0).([]*domain.Attribute)
	}
	return attributes, args.Error(1)
}

// MockProductVariantRepository is a mock implementation of domain.ProductVariantRepository
type MockProductVariantRepository struct {
	mock.Mock
}

func (m *MockProductVariantRepository) Save(ctx context.Context, variant *domain.ProductVariant) error {
	args := m.Called(ctx, variant)
	return args.Error(0)
}

func (m *MockProductVariantRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.ProductVariant, error) {
	args := m.Called(ctx, id)
	var variant *domain.ProductVariant
	if args.Get(0) != nil {
		variant = args.Get(0).(*domain.ProductVariant)
	}
	return variant, args.Error(1)
}

func (m *MockProductVariantRepository) FindBySKU(ctx context.Context, sku string) (*domain.ProductVariant, error) {
	args := m.Called(ctx, sku)
	var variant *domain.ProductVariant
	if args.Get(0) != nil {
		variant = args.Get(0).(*domain.ProductVariant)
	}
	return variant, args.Error(1)
}

func (m *MockProductVariantRepository) FindByProductIDAndAttributeValues(ctx context.Context, productID uuid.UUID, attributeValueIDs []uuid.UUID) (*domain.ProductVariant, error) {
	args := m.Called(ctx, productID, attributeValueIDs)
	var variant *domain.ProductVariant
	if args.Get(0) != nil {
		variant = args.Get(0).(*domain.ProductVariant)
	}
	return variant, args.Error(1)
}

// MockPartyServiceConfigurationRepository is a mock implementation of domain.PartyServiceConfigurationRepository
type MockPartyServiceConfigurationRepository struct {
	mock.Mock
}

func (m *MockPartyServiceConfigurationRepository) Save(ctx context.Context, config *domain.PartyServiceConfiguration) error {
	args := m.Called(ctx, config)
	return args.Error(0)
}

func (m *MockPartyServiceConfigurationRepository) FindByID(ctx context.Context, partyID, id uuid.UUID) (*domain.PartyServiceConfiguration, error) {
	args := m.Called(ctx, partyID, id)
	var config *domain.PartyServiceConfiguration
	if args.Get(0) != nil {
		config = args.Get(0).(*domain.PartyServiceConfiguration)
	}
	return config, args.Error(1)
}

func (m *MockPartyServiceConfigurationRepository) FindByPartyID(ctx context.Context, partyID uuid.UUID) ([]*domain.PartyServiceConfiguration, error) {
	args := m.Called(ctx, partyID)
	var configs []*domain.PartyServiceConfiguration
	if args.Get(0) != nil {
		configs = args.Get(0).([]*domain.PartyServiceConfiguration)
	}
	return configs, args.Error(1)
}

func (m *MockPartyServiceConfigurationRepository) Delete(ctx context.Context, partyID, id uuid.UUID) error {
	args := m.Called(ctx, partyID, id)
	return args.Error(0)
}

func TestProductService_CreateProduct(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockBrandRepo := new(MockBrandRepository)
	mockGroupRepo := new(MockProductGroupRepository)
	mockAttributeRepo := new(MockAttributeRepository)
	mockVariantRepo := new(MockProductVariantRepository)
	mockPartyServiceConfigRepo := new(MockPartyServiceConfigurationRepository)

	productService := application.NewProductService(mockProductRepo, mockBrandRepo, mockGroupRepo, mockAttributeRepo, mockVariantRepo, mockPartyServiceConfigRepo)
	ctx := context.Background()

	brandID := uuid.New()
	groupID1 := uuid.New()
	groupID2 := uuid.New()
	barcode := "1234567890123"

	t.Run("should create product successfully", func(t *testing.T) {
		mockProductRepo.ExpectedCalls = nil
		mockBrandRepo.ExpectedCalls = nil
		mockGroupRepo.ExpectedCalls = nil
		cmd := application.CreateProductCommand{
			SKU:         "TEST-SKU",
			Name:        "Test Product",
			LongName:    "Long Test Product Name",
			Barcode:     &barcode,
			Description: "Description for test product",
			ProductType: domain.ProductTypeTangible,
			BrandID:     brandID,
			GroupIDs:    []uuid.UUID{groupID1, groupID2},
		}

		mockBrandRepo.On("FindByID", ctx, brandID).Return(&domain.Brand{ID: brandID, Name: "TestBrand"}, nil).Once()
		mockGroupRepo.On("FindByID", ctx, groupID1).Return(&domain.ProductGroup{ID: groupID1, Name: "Group1"}, nil).Once()
		mockGroupRepo.On("FindByID", ctx, groupID2).Return(&domain.ProductGroup{ID: groupID2, Name: "Group2"}, nil).Once()

		mockProductRepo.On("Save", ctx, mock.AnythingOfType("*domain.Product")).Return(nil).Once()

		productDTO, err := productService.CreateProduct(ctx, cmd)

		assert.NoError(t, err)
		assert.NotNil(t, productDTO)
		assert.Equal(t, cmd.SKU, productDTO.SKU)
		assert.Equal(t, cmd.Name, productDTO.Name)
		assert.Equal(t, cmd.BrandID, productDTO.BrandID)
		assert.Len(t, productDTO.GroupIDs, 2)
		mockProductRepo.AssertExpectations(t)
		mockBrandRepo.AssertExpectations(t)
		mockGroupRepo.AssertExpectations(t)
	})

	t.Run("should return error if brand not found", func(t *testing.T) {
		mockProductRepo.ExpectedCalls = nil
		mockBrandRepo.ExpectedCalls = nil
		mockGroupRepo.ExpectedCalls = nil
		invalidBrandID := uuid.New()
		cmd := application.CreateProductCommand{
			SKU:         "TEST-SKU",
			Name:        "Test Product",
			ProductType: domain.ProductTypeTangible,
			BrandID:     invalidBrandID,
		}

		mockBrandRepo.On("FindByID", ctx, invalidBrandID).Return(nil, nil).Once() // Simulate not found

		productDTO, err := productService.CreateProduct(ctx, cmd)

		assert.Error(t, err)
		assert.EqualError(t, err, fmt.Sprintf("brand with ID %s does not exist", invalidBrandID))
		assert.Nil(t, productDTO)
		mockBrandRepo.AssertExpectations(t)
	})

	t.Run("should return error if group not found", func(t *testing.T) {
		mockProductRepo.ExpectedCalls = nil
		mockBrandRepo.ExpectedCalls = nil
		mockGroupRepo.ExpectedCalls = nil
		invalidGroupID := uuid.New()
		cmd := application.CreateProductCommand{
			SKU:         "TEST-SKU",
			Name:        "Test Product",
			ProductType: domain.ProductTypeTangible,
			BrandID:     brandID,
			GroupIDs:    []uuid.UUID{invalidGroupID},
		}

		mockBrandRepo.On("FindByID", ctx, brandID).Return(&domain.Brand{ID: brandID, Name: "TestBrand"}, nil).Once()
		mockGroupRepo.On("FindByID", ctx, invalidGroupID).Return(nil, nil).Once() // Simulate not found

		productDTO, err := productService.CreateProduct(ctx, cmd)

		assert.Error(t, err)
		assert.EqualError(t, err, fmt.Sprintf("product group with ID %s does not exist", invalidGroupID))
		assert.Nil(t, productDTO)
		mockBrandRepo.AssertExpectations(t)
		mockGroupRepo.AssertExpectations(t)
	})

	t.Run("should return error if product domain entity creation fails", func(t *testing.T) {
		mockProductRepo.ExpectedCalls = nil
		mockBrandRepo.ExpectedCalls = nil
		mockGroupRepo.ExpectedCalls = nil
		cmd := application.CreateProductCommand{
			SKU:         "", // Invalid SKU
			Name:        "Test Product",
			ProductType: domain.ProductTypeTangible,
			BrandID:     brandID,
		}

		mockBrandRepo.On("FindByID", ctx, brandID).Return(&domain.Brand{ID: brandID, Name: "TestBrand"}, nil).Once()

		productDTO, err := productService.CreateProduct(ctx, cmd)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "product SKU cannot be empty")
		assert.Nil(t, productDTO)
		mockBrandRepo.AssertExpectations(t)
	})

	t.Run("should return error if product repository save fails", func(t *testing.T) {
		mockProductRepo.ExpectedCalls = nil
		mockBrandRepo.ExpectedCalls = nil
		mockGroupRepo.ExpectedCalls = nil
		cmd := application.CreateProductCommand{
			SKU:         "TEST-SKU",
			Name:        "Test Product",
			ProductType: domain.ProductTypeTangible,
			BrandID:     brandID,
		}

		mockBrandRepo.On("FindByID", ctx, brandID).Return(&domain.Brand{ID: brandID, Name: "TestBrand"}, nil).Once()
		mockProductRepo.On("Save", ctx, mock.AnythingOfType("*domain.Product")).Return(assert.AnError).Once()

		productDTO, err := productService.CreateProduct(ctx, cmd)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to save product")
		assert.Nil(t, productDTO)
		mockProductRepo.AssertExpectations(t)
		mockBrandRepo.AssertExpectations(t)
	})
}

func TestProductService_AddGroupToProduct(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockBrandRepo := new(MockBrandRepository)
	mockGroupRepo := new(MockProductGroupRepository)
	mockAttributeRepo := new(MockAttributeRepository)
	mockVariantRepo := new(MockProductVariantRepository)
	mockPartyServiceConfigRepo := new(MockPartyServiceConfigurationRepository)
	productService := application.NewProductService(mockProductRepo, mockBrandRepo, mockGroupRepo, mockAttributeRepo, mockVariantRepo, mockPartyServiceConfigRepo)
	ctx := context.Background()

	productID := uuid.New()
	existingBrandID := uuid.New()
	newGroupID := uuid.New()
	existingProduct := &domain.Product{
		ID:        productID,
		SKU:       "PROD1",
		Name:      "Product 1",
		BrandID:   existingBrandID,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("should add group to product successfully", func(t *testing.T) {
		mockProductRepo.ExpectedCalls = nil
		mockBrandRepo.ExpectedCalls = nil
		mockGroupRepo.ExpectedCalls = nil
		cmd := application.AddGroupCommand{
			ProductID: productID,
			GroupID:   newGroupID,
		}

		mockProductRepo.On("FindByID", ctx, productID).Return(existingProduct, nil).Once()
		mockGroupRepo.On("FindByID", ctx, newGroupID).Return(&domain.ProductGroup{ID: newGroupID, Name: "NewGroup"}, nil).Once()
		mockProductRepo.On("Save", ctx, mock.AnythingOfType("*domain.Product")).Return(nil).Once()

		productDTO, err := productService.AddGroupToProduct(ctx, cmd)

		assert.NoError(t, err)
		assert.NotNil(t, productDTO)
		assert.Contains(t, productDTO.GroupIDs, newGroupID)
		mockProductRepo.AssertExpectations(t)
		mockGroupRepo.AssertExpectations(t)
	})

	t.Run("should return error if product not found", func(t *testing.T) {
		mockProductRepo.ExpectedCalls = nil
		mockBrandRepo.ExpectedCalls = nil
		mockGroupRepo.ExpectedCalls = nil
		cmd := application.AddGroupCommand{
			ProductID: productID,
			GroupID:   newGroupID,
		}

		mockProductRepo.On("FindByID", ctx, productID).Return(nil, nil).Once() // Product not found

		productDTO, err := productService.AddGroupToProduct(ctx, cmd)

		assert.Error(t, err)
		assert.EqualError(t, err, fmt.Sprintf("product with ID %s does not exist", cmd.ProductID))
		assert.Nil(t, productDTO)
		mockProductRepo.AssertExpectations(t)
	})

	t.Run("should return error if group not found", func(t *testing.T) {
		mockProductRepo.ExpectedCalls = nil
		mockBrandRepo.ExpectedCalls = nil
		mockGroupRepo.ExpectedCalls = nil
		cmd := application.AddGroupCommand{
			ProductID: productID,
			GroupID:   newGroupID,
		}

		mockProductRepo.On("FindByID", ctx, productID).Return(existingProduct, nil).Once()
		mockGroupRepo.On("FindByID", ctx, newGroupID).Return(nil, nil).Once() // Group not found

		productDTO, err := productService.AddGroupToProduct(ctx, cmd)

		assert.Error(t, err)
		assert.EqualError(t, err, fmt.Sprintf("product group with ID %s does not exist", cmd.GroupID))
		assert.Nil(t, productDTO)
		mockProductRepo.AssertExpectations(t)
		mockGroupRepo.AssertExpectations(t)
	})
}

func TestProductService_AddDirectAttributeToProduct(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockBrandRepo := new(MockBrandRepository)
	mockGroupRepo := new(MockProductGroupRepository)
	mockAttributeRepo := new(MockAttributeRepository)
	mockVariantRepo := new(MockProductVariantRepository)
	mockPartyServiceConfigRepo := new(MockPartyServiceConfigurationRepository)
	productService := application.NewProductService(mockProductRepo, mockBrandRepo, mockGroupRepo, mockAttributeRepo, mockVariantRepo, mockPartyServiceConfigRepo)
	ctx := context.Background()

	productID := uuid.New()
	existingBrandID := uuid.New()
	newAttributeID := uuid.New()
	existingProduct := &domain.Product{
		ID:        productID,
		SKU:       "PROD1",
		Name:      "Product 1",
		BrandID:   existingBrandID,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("should add direct attribute to product successfully", func(t *testing.T) {
		mockProductRepo.ExpectedCalls = nil
		mockBrandRepo.ExpectedCalls = nil
		mockGroupRepo.ExpectedCalls = nil
		mockAttributeRepo.ExpectedCalls = nil
		cmd := application.AddDirectAttributeCommand{
			ProductID:   productID,
			AttributeID: newAttributeID,
		}

		mockProductRepo.On("FindByID", ctx, productID).Return(existingProduct, nil).Once()
		mockAttributeRepo.On("FindByID", ctx, newAttributeID).Return(&domain.Attribute{ID: newAttributeID}, nil).Once()
		mockProductRepo.On("Save", ctx, mock.AnythingOfType("*domain.Product")).Return(nil).Once()

		productDTO, err := productService.AddDirectAttributeToProduct(ctx, cmd)

		assert.NoError(t, err)
		assert.NotNil(t, productDTO)
		assert.Contains(t, productDTO.DirectAttributeIDs, newAttributeID)
		mockProductRepo.AssertExpectations(t)
	})

	t.Run("should return error if product not found", func(t *testing.T) {
		mockProductRepo.ExpectedCalls = nil
		mockBrandRepo.ExpectedCalls = nil
		mockGroupRepo.ExpectedCalls = nil
		cmd := application.AddDirectAttributeCommand{
			ProductID:   productID,
			AttributeID: newAttributeID,
		}

		mockProductRepo.On("FindByID", ctx, productID).Return(nil, nil).Once() // Product not found

		productDTO, err := productService.AddDirectAttributeToProduct(ctx, cmd)

		assert.Error(t, err)
		assert.EqualError(t, err, fmt.Sprintf("product with ID %s does not exist", cmd.ProductID))
		assert.Nil(t, productDTO)
		mockProductRepo.AssertExpectations(t)
	})
}

func TestProductService_UpdateProductSKU(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockBrandRepo := new(MockBrandRepository)
	mockGroupRepo := new(MockProductGroupRepository)
	mockAttributeRepo := new(MockAttributeRepository)
	mockVariantRepo := new(MockProductVariantRepository)
	mockPartyServiceConfigRepo := new(MockPartyServiceConfigurationRepository)
	productService := application.NewProductService(mockProductRepo, mockBrandRepo, mockGroupRepo, mockAttributeRepo, mockVariantRepo, mockPartyServiceConfigRepo)
	ctx := context.Background()

	productID := uuid.New()
	oldSKU := "OLD-SKU"
	newSKU := "NEW-SKU"
	existingProduct := &domain.Product{
		ID:        productID,
		SKU:       oldSKU,
		Name:      "Product to Update",
		BrandID:   uuid.New(),
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	updatedProduct := &domain.Product{
		ID:        productID,
		SKU:       newSKU,
		Name:      "Product to Update",
		BrandID:   uuid.New(),
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("should update product SKU and cascade successfully", func(t *testing.T) {
		mockProductRepo.ExpectedCalls = nil
		mockBrandRepo.ExpectedCalls = nil
		mockGroupRepo.ExpectedCalls = nil
		cmd := application.UpdateProductSKUCommand{
			ProductID: productID,
			NewSKU:    newSKU,
		}

		mockProductRepo.On("FindByID", ctx, productID).Return(existingProduct, nil).Once()
		mockProductRepo.On("UpdateSKUs", ctx, productID, newSKU).Return(nil).Once()
		mockProductRepo.On("FindByID", ctx, productID).Return(updatedProduct, nil).Once() // Fetch updated product

		productDTO, err := productService.UpdateProductSKU(ctx, cmd)

		assert.NoError(t, err)
		assert.NotNil(t, productDTO)
		assert.Equal(t, newSKU, productDTO.SKU)
		mockProductRepo.AssertExpectations(t)
	})

	t.Run("should return error if product not found", func(t *testing.T) {
		mockProductRepo.ExpectedCalls = nil
		mockBrandRepo.ExpectedCalls = nil
		mockGroupRepo.ExpectedCalls = nil
		cmd := application.UpdateProductSKUCommand{
			ProductID: productID,
			NewSKU:    newSKU,
		}

		mockProductRepo.On("FindByID", ctx, productID).Return(nil, nil).Once()

		productDTO, err := productService.UpdateProductSKU(ctx, cmd)

		assert.Error(t, err)
		assert.EqualError(t, err, fmt.Sprintf("product with ID %s does not exist", cmd.ProductID))
		assert.Nil(t, productDTO)
		mockProductRepo.AssertExpectations(t)
	})

	t.Run("should return error if cascade update fails", func(t *testing.T) {
		cmd := application.UpdateProductSKUCommand{
			ProductID: productID,
			NewSKU:    newSKU,
		}

		mockProductRepo.ExpectedCalls = nil
		mockProductRepo.On("FindByID", ctx, productID).Return(existingProduct, nil).Once()
		mockProductRepo.On("UpdateSKUs", ctx, productID, newSKU).Return(assert.AnError).Once()

		productDTO, err := productService.UpdateProductSKU(ctx, cmd)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to update product SKU and cascade to variants")
		assert.Nil(t, productDTO)
		mockProductRepo.AssertExpectations(t)
	})
}
