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

func newTestService() (*application.ProductService, *MockProductRepository, *MockBrandRepository, *MockProductGroupRepository, *MockAttributeRepository, *MockProductVariantRepository, *MockPartyServiceConfigurationRepository) {
	mockProductRepo := new(MockProductRepository)
	mockBrandRepo := new(MockBrandRepository)
	mockGroupRepo := new(MockProductGroupRepository)
	mockAttributeRepo := new(MockAttributeRepository)
	mockVariantRepo := new(MockProductVariantRepository)
	mockPSCRepo := new(MockPartyServiceConfigurationRepository)
	svc := application.NewProductService(mockProductRepo, mockBrandRepo, mockGroupRepo, mockAttributeRepo, mockVariantRepo, mockPSCRepo)
	return svc, mockProductRepo, mockBrandRepo, mockGroupRepo, mockAttributeRepo, mockVariantRepo, mockPSCRepo
}

func actorCtx() context.Context {
	return context.WithValue(context.Background(), "actorID", testActorID)
}

// ============================================================================
// Brand CRUD Tests
// ============================================================================

func TestProductService_ListBrands(t *testing.T) {
	svc, _, mockBrandRepo, _, _, _, _ := newTestService()
	ctx := actorCtx()

	t.Run("should list brands successfully", func(t *testing.T) {
		mockBrandRepo.ExpectedCalls = nil
		brand1 := &domain.Brand{ID: uuid.New(), Name: "Brand A", DefaultMarkupPercentage: 10.0, IsActive: true}
		brand2 := &domain.Brand{ID: uuid.New(), Name: "Brand B", DefaultMarkupPercentage: 20.0, IsActive: false}
		mockBrandRepo.On("FindAll", ctx).Return([]*domain.Brand{brand1, brand2}, nil).Once()

		result, err := svc.ListBrands(ctx)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, brand1.Name, result[0].Name)
		assert.Equal(t, brand2.Name, result[1].Name)
		assert.Equal(t, brand1.DefaultMarkupPercentage, result[0].DefaultMarkupPercentage)
		mockBrandRepo.AssertExpectations(t)
	})

	t.Run("should return empty list when no brands exist", func(t *testing.T) {
		mockBrandRepo.ExpectedCalls = nil
		mockBrandRepo.On("FindAll", ctx).Return([]*domain.Brand{}, nil).Once()

		result, err := svc.ListBrands(ctx)

		assert.NoError(t, err)
		assert.Len(t, result, 0)
	})

	t.Run("should return error on repository failure", func(t *testing.T) {
		mockBrandRepo.ExpectedCalls = nil
		mockBrandRepo.On("FindAll", ctx).Return(nil, errors.New("db error")).Once()

		result, err := svc.ListBrands(ctx)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestProductService_GetBrandByID(t *testing.T) {
	svc, _, mockBrandRepo, _, _, _, _ := newTestService()
	ctx := actorCtx()

	t.Run("should get brand by ID", func(t *testing.T) {
		mockBrandRepo.ExpectedCalls = nil
		brandID := uuid.New()
		brand := &domain.Brand{ID: brandID, Name: "Test Brand", DefaultMarkupPercentage: 15.0, IsActive: true}
		mockBrandRepo.On("FindByID", ctx, brandID).Return(brand, nil).Once()

		result, err := svc.GetBrandByID(ctx, brandID)

		assert.NoError(t, err)
		assert.Equal(t, brandID, result.ID)
		assert.Equal(t, "Test Brand", result.Name)
		assert.Equal(t, 15.0, result.DefaultMarkupPercentage)
		assert.True(t, result.IsActive)
	})

	t.Run("should return error when brand not found", func(t *testing.T) {
		mockBrandRepo.ExpectedCalls = nil
		brandID := uuid.New()
		mockBrandRepo.On("FindByID", ctx, brandID).Return(nil, nil).Once()

		result, err := svc.GetBrandByID(ctx, brandID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "does not exist")
	})

	t.Run("should return error on repository failure", func(t *testing.T) {
		mockBrandRepo.ExpectedCalls = nil
		brandID := uuid.New()
		mockBrandRepo.On("FindByID", ctx, brandID).Return(nil, errors.New("db error")).Once()

		result, err := svc.GetBrandByID(ctx, brandID)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestProductService_CreateBrand(t *testing.T) {
	svc, _, mockBrandRepo, _, _, _, _ := newTestService()
	ctx := actorCtx()

	t.Run("should create brand successfully", func(t *testing.T) {
		mockBrandRepo.ExpectedCalls = nil
		cmd := application.CreateBrandCommand{
			ActorID:                 testActorID,
			Name:                    "New Brand",
			DefaultMarkupPercentage: 25.0,
			IsActive:                true,
		}
		mockBrandRepo.On("Save", ctx, mock.AnythingOfType("*domain.Brand")).Return(nil).Once()

		result, err := svc.CreateBrand(ctx, cmd)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "New Brand", result.Name)
		assert.Equal(t, 25.0, result.DefaultMarkupPercentage)
		assert.True(t, result.IsActive)
		mockBrandRepo.AssertExpectations(t)
	})

	t.Run("should return error for empty name", func(t *testing.T) {
		mockBrandRepo.ExpectedCalls = nil
		cmd := application.CreateBrandCommand{
			ActorID: testActorID,
			Name:    "",
		}

		result, err := svc.CreateBrand(ctx, cmd)

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("should return error for negative markup", func(t *testing.T) {
		mockBrandRepo.ExpectedCalls = nil
		cmd := application.CreateBrandCommand{
			ActorID:                 testActorID,
			Name:                    "Brand",
			DefaultMarkupPercentage: -5.0,
		}

		result, err := svc.CreateBrand(ctx, cmd)

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("should return error on save failure", func(t *testing.T) {
		mockBrandRepo.ExpectedCalls = nil
		cmd := application.CreateBrandCommand{
			ActorID:                 testActorID,
			Name:                    "Brand",
			DefaultMarkupPercentage: 10.0,
			IsActive:                true,
		}
		mockBrandRepo.On("Save", ctx, mock.AnythingOfType("*domain.Brand")).Return(errors.New("save failed")).Once()

		result, err := svc.CreateBrand(ctx, cmd)

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("should return error when no actor ID", func(t *testing.T) {
		cmd := application.CreateBrandCommand{
			ActorID: "",
			Name:    "Brand",
		}

		result, err := svc.CreateBrand(context.Background(), cmd)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestProductService_UpdateBrand(t *testing.T) {
	svc, _, mockBrandRepo, _, _, _, _ := newTestService()
	ctx := actorCtx()

	t.Run("should update brand name", func(t *testing.T) {
		mockBrandRepo.ExpectedCalls = nil
		brandID := uuid.New()
		existing := &domain.Brand{ID: brandID, Name: "Old Name", DefaultMarkupPercentage: 10.0, IsActive: true}
		newName := "New Name"
		cmd := application.UpdateBrandCommand{
			ActorID: testActorID,
			ID:      brandID,
			Name:    &newName,
		}
		mockBrandRepo.On("FindByID", ctx, brandID).Return(existing, nil).Once()
		mockBrandRepo.On("Save", ctx, mock.AnythingOfType("*domain.Brand")).Return(nil).Once()

		result, err := svc.UpdateBrand(ctx, cmd)

		assert.NoError(t, err)
		assert.Equal(t, "New Name", result.Name)
		mockBrandRepo.AssertExpectations(t)
	})

	t.Run("should update brand markup percentage", func(t *testing.T) {
		mockBrandRepo.ExpectedCalls = nil
		brandID := uuid.New()
		existing := &domain.Brand{ID: brandID, Name: "Brand", DefaultMarkupPercentage: 10.0, IsActive: true}
		newMarkup := 30.0
		cmd := application.UpdateBrandCommand{
			ActorID:                 testActorID,
			ID:                      brandID,
			DefaultMarkupPercentage: &newMarkup,
		}
		mockBrandRepo.On("FindByID", ctx, brandID).Return(existing, nil).Once()
		mockBrandRepo.On("Save", ctx, mock.AnythingOfType("*domain.Brand")).Return(nil).Once()

		result, err := svc.UpdateBrand(ctx, cmd)

		assert.NoError(t, err)
		assert.Equal(t, 30.0, result.DefaultMarkupPercentage)
	})

	t.Run("should update brand isActive", func(t *testing.T) {
		mockBrandRepo.ExpectedCalls = nil
		brandID := uuid.New()
		existing := &domain.Brand{ID: brandID, Name: "Brand", DefaultMarkupPercentage: 10.0, IsActive: true}
		isActive := false
		cmd := application.UpdateBrandCommand{
			ActorID:  testActorID,
			ID:       brandID,
			IsActive: &isActive,
		}
		mockBrandRepo.On("FindByID", ctx, brandID).Return(existing, nil).Once()
		mockBrandRepo.On("Save", ctx, mock.AnythingOfType("*domain.Brand")).Return(nil).Once()

		result, err := svc.UpdateBrand(ctx, cmd)

		assert.NoError(t, err)
		assert.False(t, result.IsActive)
	})

	t.Run("should return error for negative markup", func(t *testing.T) {
		mockBrandRepo.ExpectedCalls = nil
		brandID := uuid.New()
		existing := &domain.Brand{ID: brandID, Name: "Brand", DefaultMarkupPercentage: 10.0, IsActive: true}
		negativeMarkup := -5.0
		cmd := application.UpdateBrandCommand{
			ActorID:                 testActorID,
			ID:                      brandID,
			DefaultMarkupPercentage: &negativeMarkup,
		}
		mockBrandRepo.On("FindByID", ctx, brandID).Return(existing, nil).Once()

		result, err := svc.UpdateBrand(ctx, cmd)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "negative")
	})

	t.Run("should return error when brand not found", func(t *testing.T) {
		mockBrandRepo.ExpectedCalls = nil
		brandID := uuid.New()
		newName := "Name"
		cmd := application.UpdateBrandCommand{
			ActorID: testActorID,
			ID:      brandID,
			Name:    &newName,
		}
		mockBrandRepo.On("FindByID", ctx, brandID).Return(nil, nil).Once()

		result, err := svc.UpdateBrand(ctx, cmd)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "does not exist")
	})

	t.Run("should return error for empty name update", func(t *testing.T) {
		mockBrandRepo.ExpectedCalls = nil
		brandID := uuid.New()
		existing := &domain.Brand{ID: brandID, Name: "Brand", DefaultMarkupPercentage: 10.0, IsActive: true}
		emptyName := ""
		cmd := application.UpdateBrandCommand{
			ActorID: testActorID,
			ID:      brandID,
			Name:    &emptyName,
		}
		mockBrandRepo.On("FindByID", ctx, brandID).Return(existing, nil).Once()

		result, err := svc.UpdateBrand(ctx, cmd)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestProductService_DeleteBrand(t *testing.T) {
	svc, mockProductRepo, mockBrandRepo, _, _, _, _ := newTestService()
	ctx := actorCtx()

	t.Run("should delete brand successfully", func(t *testing.T) {
		mockBrandRepo.ExpectedCalls = nil
		mockProductRepo.ExpectedCalls = nil
		brandID := uuid.New()
		existing := &domain.Brand{ID: brandID, Name: "Brand", IsActive: true}
		mockBrandRepo.On("FindByID", ctx, brandID).Return(existing, nil).Once()
		mockProductRepo.On("FindAll", ctx).Return([]*domain.Product{}, nil).Once()
		mockBrandRepo.On("Delete", ctx, brandID).Return(nil).Once()

		err := svc.DeleteBrand(ctx, application.DeleteBrandCommand{ActorID: testActorID, ID: brandID})

		assert.NoError(t, err)
		mockBrandRepo.AssertExpectations(t)
		mockProductRepo.AssertExpectations(t)
	})

	t.Run("should return error when brand not found", func(t *testing.T) {
		mockBrandRepo.ExpectedCalls = nil
		brandID := uuid.New()
		mockBrandRepo.On("FindByID", ctx, brandID).Return(nil, nil).Once()

		err := svc.DeleteBrand(ctx, application.DeleteBrandCommand{ActorID: testActorID, ID: brandID})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "does not exist")
	})

	t.Run("should return error on delete failure", func(t *testing.T) {
		mockBrandRepo.ExpectedCalls = nil
		mockProductRepo.ExpectedCalls = nil
		brandID := uuid.New()
		existing := &domain.Brand{ID: brandID, Name: "Brand", IsActive: true}
		mockBrandRepo.On("FindByID", ctx, brandID).Return(existing, nil).Once()
		mockProductRepo.On("FindAll", ctx).Return([]*domain.Product{}, nil).Once()
		mockBrandRepo.On("Delete", ctx, brandID).Return(errors.New("constraint violation")).Once()

		err := svc.DeleteBrand(ctx, application.DeleteBrandCommand{ActorID: testActorID, ID: brandID})

		assert.Error(t, err)
	})

	t.Run("should return error when brand is used by products", func(t *testing.T) {
		mockBrandRepo.ExpectedCalls = nil
		mockProductRepo.ExpectedCalls = nil
		brandID := uuid.New()
		existing := &domain.Brand{ID: brandID, Name: "Brand", IsActive: true}
		inUseProduct := &domain.Product{ID: uuid.New(), BrandID: brandID}
		mockBrandRepo.On("FindByID", ctx, brandID).Return(existing, nil).Once()
		mockProductRepo.On("FindAll", ctx).Return([]*domain.Product{inUseProduct}, nil).Once()

		err := svc.DeleteBrand(ctx, application.DeleteBrandCommand{ActorID: testActorID, ID: brandID})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot delete brand")
	})
}

// ============================================================================
// ProductGroup CRUD Tests
// ============================================================================

func TestProductService_ListProductGroups(t *testing.T) {
	svc, _, _, mockGroupRepo, _, _, _ := newTestService()
	ctx := actorCtx()

	t.Run("should list product groups successfully", func(t *testing.T) {
		mockGroupRepo.ExpectedCalls = nil
		group1 := &domain.ProductGroup{ID: uuid.New(), Name: "Tangibles", Type: domain.ProductGroupTypeTangible, IsActive: true}
		group2 := &domain.ProductGroup{ID: uuid.New(), Name: "Services", Type: domain.ProductGroupTypeService, IsActive: true}
		mockGroupRepo.On("FindAll", ctx).Return([]*domain.ProductGroup{group1, group2}, nil).Once()

		result, err := svc.ListProductGroups(ctx)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "Tangibles", result[0].Name)
		assert.Equal(t, "TANGIBLE", result[0].Type)
		assert.Equal(t, "Services", result[1].Name)
		assert.Equal(t, "SERVICE", result[1].Type)
	})

	t.Run("should return empty list when no groups exist", func(t *testing.T) {
		mockGroupRepo.ExpectedCalls = nil
		mockGroupRepo.On("FindAll", ctx).Return([]*domain.ProductGroup{}, nil).Once()

		result, err := svc.ListProductGroups(ctx)

		assert.NoError(t, err)
		assert.Len(t, result, 0)
	})

	t.Run("should return error on repository failure", func(t *testing.T) {
		mockGroupRepo.ExpectedCalls = nil
		mockGroupRepo.On("FindAll", ctx).Return(nil, errors.New("db error")).Once()

		result, err := svc.ListProductGroups(ctx)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestProductService_GetProductGroupByID(t *testing.T) {
	svc, _, _, mockGroupRepo, _, _, _ := newTestService()
	ctx := actorCtx()

	t.Run("should get product group by ID", func(t *testing.T) {
		mockGroupRepo.ExpectedCalls = nil
		groupID := uuid.New()
		parentID := uuid.New()
		group := &domain.ProductGroup{ID: groupID, Name: "Test Group", Type: domain.ProductGroupTypeTangible, ParentGroupID: &parentID, IsActive: true}
		mockGroupRepo.On("FindByID", ctx, groupID).Return(group, nil).Once()

		result, err := svc.GetProductGroupByID(ctx, groupID)

		assert.NoError(t, err)
		assert.Equal(t, groupID, result.ID)
		assert.Equal(t, "Test Group", result.Name)
		assert.Equal(t, "TANGIBLE", result.Type)
		assert.Equal(t, &parentID, result.ParentGroupID)
		assert.True(t, result.IsActive)
	})

	t.Run("should return error when group not found", func(t *testing.T) {
		mockGroupRepo.ExpectedCalls = nil
		groupID := uuid.New()
		mockGroupRepo.On("FindByID", ctx, groupID).Return(nil, nil).Once()

		result, err := svc.GetProductGroupByID(ctx, groupID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "does not exist")
	})
}

func TestProductService_CreateProductGroup(t *testing.T) {
	svc, _, _, mockGroupRepo, _, _, _ := newTestService()
	ctx := actorCtx()

	t.Run("should create tangible product group", func(t *testing.T) {
		mockGroupRepo.ExpectedCalls = nil
		cmd := application.CreateProductGroupCommand{
			ActorID:  testActorID,
			Name:     "Fabrics",
			Type:     "TANGIBLE",
			IsActive: true,
		}
		mockGroupRepo.On("Save", ctx, mock.AnythingOfType("*domain.ProductGroup")).Return(nil).Once()

		result, err := svc.CreateProductGroup(ctx, cmd)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Fabrics", result.Name)
		assert.Equal(t, "TANGIBLE", result.Type)
		assert.True(t, result.IsActive)
		mockGroupRepo.AssertExpectations(t)
	})

	t.Run("should create service product group", func(t *testing.T) {
		mockGroupRepo.ExpectedCalls = nil
		cmd := application.CreateProductGroupCommand{
			ActorID:  testActorID,
			Name:     "Treatments",
			Type:     "SERVICE",
			IsActive: true,
		}
		mockGroupRepo.On("Save", ctx, mock.AnythingOfType("*domain.ProductGroup")).Return(nil).Once()

		result, err := svc.CreateProductGroup(ctx, cmd)

		assert.NoError(t, err)
		assert.Equal(t, "SERVICE", result.Type)
	})

	t.Run("should create product group with parent", func(t *testing.T) {
		mockGroupRepo.ExpectedCalls = nil
		parentID := uuid.New()
		parent := &domain.ProductGroup{ID: parentID, Name: "Parent", Type: domain.ProductGroupTypeTangible, IsActive: true}
		cmd := application.CreateProductGroupCommand{
			ActorID:  testActorID,
			Name:     "Child Group",
			Type:     "TANGIBLE",
			ParentID: &parentID,
			IsActive: true,
		}
		mockGroupRepo.On("FindByID", ctx, parentID).Return(parent, nil).Once()
		mockGroupRepo.On("Save", ctx, mock.AnythingOfType("*domain.ProductGroup")).Return(nil).Once()

		result, err := svc.CreateProductGroup(ctx, cmd)

		assert.NoError(t, err)
		assert.Equal(t, "Child Group", result.Name)
	})

	t.Run("should return error for invalid group type", func(t *testing.T) {
		mockGroupRepo.ExpectedCalls = nil
		cmd := application.CreateProductGroupCommand{
			ActorID:  testActorID,
			Name:     "Bad Group",
			Type:     "INVALID",
			IsActive: true,
		}

		result, err := svc.CreateProductGroup(ctx, cmd)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid group type")
	})

	t.Run("should return error for empty name", func(t *testing.T) {
		mockGroupRepo.ExpectedCalls = nil
		cmd := application.CreateProductGroupCommand{
			ActorID:  testActorID,
			Name:     "",
			Type:     "TANGIBLE",
			IsActive: true,
		}

		result, err := svc.CreateProductGroup(ctx, cmd)

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("should return error when parent group not found", func(t *testing.T) {
		mockGroupRepo.ExpectedCalls = nil
		parentID := uuid.New()
		cmd := application.CreateProductGroupCommand{
			ActorID:  testActorID,
			Name:     "Child",
			Type:     "TANGIBLE",
			ParentID: &parentID,
			IsActive: true,
		}
		mockGroupRepo.On("FindByID", ctx, parentID).Return(nil, nil).Once()

		result, err := svc.CreateProductGroup(ctx, cmd)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "does not exist")
	})

	t.Run("should return error on save failure", func(t *testing.T) {
		mockGroupRepo.ExpectedCalls = nil
		cmd := application.CreateProductGroupCommand{
			ActorID:  testActorID,
			Name:     "Good Group",
			Type:     "TANGIBLE",
			IsActive: true,
		}
		mockGroupRepo.On("Save", ctx, mock.AnythingOfType("*domain.ProductGroup")).Return(errors.New("save failed")).Once()

		result, err := svc.CreateProductGroup(ctx, cmd)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestProductService_UpdateProductGroup(t *testing.T) {
	svc, _, _, mockGroupRepo, _, _, _ := newTestService()
	ctx := actorCtx()

	t.Run("should update group name", func(t *testing.T) {
		mockGroupRepo.ExpectedCalls = nil
		groupID := uuid.New()
		existing := &domain.ProductGroup{ID: groupID, Name: "Old Name", Type: domain.ProductGroupTypeTangible, IsActive: true}
		newName := "New Name"
		cmd := application.UpdateProductGroupCommand{
			ActorID: testActorID,
			ID:      groupID,
			Name:    &newName,
		}
		mockGroupRepo.On("FindByID", ctx, groupID).Return(existing, nil).Once()
		mockGroupRepo.On("Save", ctx, mock.AnythingOfType("*domain.ProductGroup")).Return(nil).Once()

		result, err := svc.UpdateProductGroup(ctx, cmd)

		assert.NoError(t, err)
		assert.Equal(t, "New Name", result.Name)
	})

	t.Run("should update group type", func(t *testing.T) {
		mockGroupRepo.ExpectedCalls = nil
		groupID := uuid.New()
		existing := &domain.ProductGroup{ID: groupID, Name: "Group", Type: domain.ProductGroupTypeTangible, IsActive: true}
		newType := "SERVICE"
		cmd := application.UpdateProductGroupCommand{
			ActorID: testActorID,
			ID:      groupID,
			Type:    &newType,
		}
		mockGroupRepo.On("FindByID", ctx, groupID).Return(existing, nil).Once()
		mockGroupRepo.On("Save", ctx, mock.AnythingOfType("*domain.ProductGroup")).Return(nil).Once()

		result, err := svc.UpdateProductGroup(ctx, cmd)

		assert.NoError(t, err)
		assert.Equal(t, "SERVICE", result.Type)
	})

	t.Run("should update group isActive", func(t *testing.T) {
		mockGroupRepo.ExpectedCalls = nil
		groupID := uuid.New()
		existing := &domain.ProductGroup{ID: groupID, Name: "Group", Type: domain.ProductGroupTypeTangible, IsActive: true}
		isActive := false
		cmd := application.UpdateProductGroupCommand{
			ActorID:  testActorID,
			ID:       groupID,
			IsActive: &isActive,
		}
		mockGroupRepo.On("FindByID", ctx, groupID).Return(existing, nil).Once()
		mockGroupRepo.On("Save", ctx, mock.AnythingOfType("*domain.ProductGroup")).Return(nil).Once()

		result, err := svc.UpdateProductGroup(ctx, cmd)

		assert.NoError(t, err)
		assert.False(t, result.IsActive)
	})

	t.Run("should update group parent", func(t *testing.T) {
		mockGroupRepo.ExpectedCalls = nil
		groupID := uuid.New()
		parentID := uuid.New()
		existing := &domain.ProductGroup{ID: groupID, Name: "Group", Type: domain.ProductGroupTypeTangible, IsActive: true}
		parent := &domain.ProductGroup{ID: parentID, Name: "Parent", Type: domain.ProductGroupTypeTangible, IsActive: true}
		cmd := application.UpdateProductGroupCommand{
			ActorID:  testActorID,
			ID:       groupID,
			ParentID: &parentID,
		}
		mockGroupRepo.On("FindByID", ctx, groupID).Return(existing, nil).Once()
		mockGroupRepo.On("FindByID", ctx, parentID).Return(parent, nil).Once()
		mockGroupRepo.On("Save", ctx, mock.AnythingOfType("*domain.ProductGroup")).Return(nil).Once()

		result, err := svc.UpdateProductGroup(ctx, cmd)

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("should return error when group not found", func(t *testing.T) {
		mockGroupRepo.ExpectedCalls = nil
		groupID := uuid.New()
		newName := "Name"
		cmd := application.UpdateProductGroupCommand{
			ActorID: testActorID,
			ID:      groupID,
			Name:    &newName,
		}
		mockGroupRepo.On("FindByID", ctx, groupID).Return(nil, nil).Once()

		result, err := svc.UpdateProductGroup(ctx, cmd)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "does not exist")
	})

	t.Run("should return error for invalid type update", func(t *testing.T) {
		mockGroupRepo.ExpectedCalls = nil
		groupID := uuid.New()
		existing := &domain.ProductGroup{ID: groupID, Name: "Group", Type: domain.ProductGroupTypeTangible, IsActive: true}
		badType := "INVALID"
		cmd := application.UpdateProductGroupCommand{
			ActorID: testActorID,
			ID:      groupID,
			Type:    &badType,
		}
		mockGroupRepo.On("FindByID", ctx, groupID).Return(existing, nil).Once()

		result, err := svc.UpdateProductGroup(ctx, cmd)

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("should return error when parent not found", func(t *testing.T) {
		mockGroupRepo.ExpectedCalls = nil
		groupID := uuid.New()
		parentID := uuid.New()
		existing := &domain.ProductGroup{ID: groupID, Name: "Group", Type: domain.ProductGroupTypeTangible, IsActive: true}
		cmd := application.UpdateProductGroupCommand{
			ActorID:  testActorID,
			ID:       groupID,
			ParentID: &parentID,
		}
		mockGroupRepo.On("FindByID", ctx, groupID).Return(existing, nil).Once()
		mockGroupRepo.On("FindByID", ctx, parentID).Return(nil, nil).Once()

		result, err := svc.UpdateProductGroup(ctx, cmd)

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("should return error for empty name update", func(t *testing.T) {
		mockGroupRepo.ExpectedCalls = nil
		groupID := uuid.New()
		existing := &domain.ProductGroup{ID: groupID, Name: "Group", Type: domain.ProductGroupTypeTangible, IsActive: true}
		emptyName := ""
		cmd := application.UpdateProductGroupCommand{
			ActorID: testActorID,
			ID:      groupID,
			Name:    &emptyName,
		}
		mockGroupRepo.On("FindByID", ctx, groupID).Return(existing, nil).Once()

		result, err := svc.UpdateProductGroup(ctx, cmd)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestProductService_DeleteProductGroup(t *testing.T) {
	svc, mockProductRepo, _, mockGroupRepo, _, _, _ := newTestService()
	ctx := actorCtx()

	t.Run("should delete product group successfully", func(t *testing.T) {
		mockGroupRepo.ExpectedCalls = nil
		mockProductRepo.ExpectedCalls = nil
		groupID := uuid.New()
		existing := &domain.ProductGroup{ID: groupID, Name: "Group", Type: domain.ProductGroupTypeTangible, IsActive: true}
		mockGroupRepo.On("FindByID", ctx, groupID).Return(existing, nil).Once()
		mockGroupRepo.On("FindAll", ctx).Return([]*domain.ProductGroup{existing}, nil).Once()
		mockProductRepo.On("FindAll", ctx).Return([]*domain.Product{}, nil).Once()
		mockGroupRepo.On("Delete", ctx, groupID).Return(nil).Once()

		err := svc.DeleteProductGroup(ctx, application.DeleteProductGroupCommand{ActorID: testActorID, ID: groupID})

		assert.NoError(t, err)
		mockGroupRepo.AssertExpectations(t)
		mockProductRepo.AssertExpectations(t)
	})

	t.Run("should return error when group not found", func(t *testing.T) {
		mockGroupRepo.ExpectedCalls = nil
		groupID := uuid.New()
		mockGroupRepo.On("FindByID", ctx, groupID).Return(nil, nil).Once()

		err := svc.DeleteProductGroup(ctx, application.DeleteProductGroupCommand{ActorID: testActorID, ID: groupID})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "does not exist")
	})

	t.Run("should return error on delete failure", func(t *testing.T) {
		mockGroupRepo.ExpectedCalls = nil
		mockProductRepo.ExpectedCalls = nil
		groupID := uuid.New()
		existing := &domain.ProductGroup{ID: groupID, Name: "Group", Type: domain.ProductGroupTypeTangible, IsActive: true}
		mockGroupRepo.On("FindByID", ctx, groupID).Return(existing, nil).Once()
		mockGroupRepo.On("FindAll", ctx).Return([]*domain.ProductGroup{existing}, nil).Once()
		mockProductRepo.On("FindAll", ctx).Return([]*domain.Product{}, nil).Once()
		mockGroupRepo.On("Delete", ctx, groupID).Return(errors.New("fk constraint")).Once()

		err := svc.DeleteProductGroup(ctx, application.DeleteProductGroupCommand{ActorID: testActorID, ID: groupID})

		assert.Error(t, err)
	})

	t.Run("should return error when group is parent of another group", func(t *testing.T) {
		mockGroupRepo.ExpectedCalls = nil
		mockProductRepo.ExpectedCalls = nil
		groupID := uuid.New()
		childID := uuid.New()
		existing := &domain.ProductGroup{ID: groupID, Name: "Parent", Type: domain.ProductGroupTypeTangible, IsActive: true}
		child := &domain.ProductGroup{ID: childID, Name: "Child", Type: domain.ProductGroupTypeTangible, ParentGroupID: &groupID, IsActive: true}
		mockGroupRepo.On("FindByID", ctx, groupID).Return(existing, nil).Once()
		mockGroupRepo.On("FindAll", ctx).Return([]*domain.ProductGroup{existing, child}, nil).Once()

		err := svc.DeleteProductGroup(ctx, application.DeleteProductGroupCommand{ActorID: testActorID, ID: groupID})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "parent of another group")
	})
}

// ============================================================================
// DeleteAttribute Test
// ============================================================================

func TestProductService_DeleteAttribute(t *testing.T) {
	svc, mockProductRepo, _, _, mockAttributeRepo, mockVariantRepo, _ := newTestService()
	ctx := actorCtx()

	t.Run("should delete attribute successfully", func(t *testing.T) {
		mockAttributeRepo.ExpectedCalls = nil
		mockProductRepo.ExpectedCalls = nil
		mockVariantRepo.ExpectedCalls = nil
		attrID := uuid.New()
		existing := &domain.Attribute{ID: attrID, Name: "Color", Code: "COLOR"}
		mockAttributeRepo.On("FindByID", ctx, attrID).Return(existing, nil).Once()
		mockProductRepo.On("FindAll", ctx).Return([]*domain.Product{}, nil).Once()
		mockAttributeRepo.On("Delete", ctx, attrID).Return(nil).Once()

		err := svc.DeleteAttribute(ctx, application.DeleteAttributeCommand{ActorID: testActorID, ID: attrID})

		assert.NoError(t, err)
		mockAttributeRepo.AssertExpectations(t)
		mockProductRepo.AssertExpectations(t)
	})

	t.Run("should return error when attribute not found", func(t *testing.T) {
		mockAttributeRepo.ExpectedCalls = nil
		attrID := uuid.New()
		mockAttributeRepo.On("FindByID", ctx, attrID).Return(nil, nil).Once()

		err := svc.DeleteAttribute(ctx, application.DeleteAttributeCommand{ActorID: testActorID, ID: attrID})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "does not exist")
	})

	t.Run("should return error on delete failure", func(t *testing.T) {
		mockAttributeRepo.ExpectedCalls = nil
		mockProductRepo.ExpectedCalls = nil
		mockVariantRepo.ExpectedCalls = nil
		attrID := uuid.New()
		existing := &domain.Attribute{ID: attrID, Name: "Color", Code: "COLOR"}
		mockAttributeRepo.On("FindByID", ctx, attrID).Return(existing, nil).Once()
		mockProductRepo.On("FindAll", ctx).Return([]*domain.Product{}, nil).Once()
		mockAttributeRepo.On("Delete", ctx, attrID).Return(errors.New("constraint")).Once()

		err := svc.DeleteAttribute(ctx, application.DeleteAttributeCommand{ActorID: testActorID, ID: attrID})

		assert.Error(t, err)
	})

	t.Run("should return error when attribute is directly assigned to a product", func(t *testing.T) {
		mockAttributeRepo.ExpectedCalls = nil
		mockProductRepo.ExpectedCalls = nil
		mockVariantRepo.ExpectedCalls = nil
		attrID := uuid.New()
		existing := &domain.Attribute{ID: attrID, Name: "Color", Code: "COLOR"}
		product := &domain.Product{ID: uuid.New(), DirectAttributeIDs: []uuid.UUID{attrID}}
		mockAttributeRepo.On("FindByID", ctx, attrID).Return(existing, nil).Once()
		mockProductRepo.On("FindAll", ctx).Return([]*domain.Product{product}, nil).Once()

		err := svc.DeleteAttribute(ctx, application.DeleteAttributeCommand{ActorID: testActorID, ID: attrID})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot delete attribute")
	})

	t.Run("should return error when variants depend on attribute values", func(t *testing.T) {
		mockAttributeRepo.ExpectedCalls = nil
		mockProductRepo.ExpectedCalls = nil
		mockVariantRepo.ExpectedCalls = nil

		attrID := uuid.New()
		valueID := uuid.New()
		productID := uuid.New()
		existing := &domain.Attribute{
			ID:   attrID,
			Name: "Color",
			Code: "COLOR",
			Values: []domain.AttributeValue{
				{ID: valueID, Value: "Red", Code: "RED"},
			},
		}
		product := &domain.Product{ID: productID}
		variant := &domain.ProductVariant{ID: uuid.New(), ProductID: productID, AttributeValues: []uuid.UUID{valueID}}

		mockAttributeRepo.On("FindByID", ctx, attrID).Return(existing, nil).Once()
		mockProductRepo.On("FindAll", ctx).Return([]*domain.Product{product}, nil).Once()
		mockVariantRepo.On("FindByProductID", ctx, productID).Return([]*domain.ProductVariant{variant}, nil).Once()

		err := svc.DeleteAttribute(ctx, application.DeleteAttributeCommand{ActorID: testActorID, ID: attrID})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "product variants depend")
	})
}

// ============================================================================
// PartyServiceConfiguration Tests
// ============================================================================

func TestProductService_ListPartyServiceConfigurationsByPartyID(t *testing.T) {
	svc, _, _, _, _, _, mockPSCRepo := newTestService()
	ctx := actorCtx()

	t.Run("should list configurations by party ID", func(t *testing.T) {
		mockPSCRepo.ExpectedCalls = nil
		partyID := uuid.New()
		config1 := &domain.PartyServiceConfiguration{ID: uuid.New(), PartyID: partyID, ServiceID: "svc1", Name: "Config 1"}
		config2 := &domain.PartyServiceConfiguration{ID: uuid.New(), PartyID: partyID, ServiceID: "svc2", Name: "Config 2"}
		mockPSCRepo.On("FindByPartyID", ctx, partyID).Return([]*domain.PartyServiceConfiguration{config1, config2}, nil).Once()

		result, err := svc.ListPartyServiceConfigurationsByPartyID(ctx, application.ListPartyServiceConfigurationsByPartyIDQuery{PartyID: partyID})

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "Config 1", result[0].Name)
		assert.Equal(t, "Config 2", result[1].Name)
	})

	t.Run("should return empty list when no configs exist", func(t *testing.T) {
		mockPSCRepo.ExpectedCalls = nil
		partyID := uuid.New()
		mockPSCRepo.On("FindByPartyID", ctx, partyID).Return([]*domain.PartyServiceConfiguration{}, nil).Once()

		result, err := svc.ListPartyServiceConfigurationsByPartyID(ctx, application.ListPartyServiceConfigurationsByPartyIDQuery{PartyID: partyID})

		assert.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("should return error on repository failure", func(t *testing.T) {
		mockPSCRepo.ExpectedCalls = nil
		partyID := uuid.New()
		mockPSCRepo.On("FindByPartyID", ctx, partyID).Return(nil, errors.New("db error")).Once()

		result, err := svc.ListPartyServiceConfigurationsByPartyID(ctx, application.ListPartyServiceConfigurationsByPartyIDQuery{PartyID: partyID})

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestProductService_DeletePartyServiceConfiguration(t *testing.T) {
	svc, _, _, _, _, _, mockPSCRepo := newTestService()
	ctx := actorCtx()

	t.Run("should delete configuration successfully", func(t *testing.T) {
		mockPSCRepo.ExpectedCalls = nil
		partyID := uuid.New()
		configID := uuid.New()
		mockPSCRepo.On("Delete", ctx, partyID, configID).Return(nil).Once()

		err := svc.DeletePartyServiceConfiguration(ctx, application.DeletePartyServiceConfigurationCommand{
			ActorID: testActorID,
			ID:      configID,
			PartyID: partyID,
		})

		assert.NoError(t, err)
		mockPSCRepo.AssertExpectations(t)
	})

	t.Run("should return error on delete failure", func(t *testing.T) {
		mockPSCRepo.ExpectedCalls = nil
		partyID := uuid.New()
		configID := uuid.New()
		mockPSCRepo.On("Delete", ctx, partyID, configID).Return(errors.New("not found")).Once()

		err := svc.DeletePartyServiceConfiguration(ctx, application.DeletePartyServiceConfigurationCommand{
			ActorID: testActorID,
			ID:      configID,
			PartyID: partyID,
		})

		assert.Error(t, err)
	})
}
