package application_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/joran-cortez/tramatex/internal/product/application"
	"github.com/joran-cortez/tramatex/internal/product/domain"
	"github.com/joran-cortez/tramatex/internal/product/infrastructure/persistence"
	"sort" // Required for sorting attribute values
)

func TestProductService_CreateProduct_Integration(t *testing.T) {
	tdb := persistence.NewTestDB(t)
	if tdb.DB == nil {
		t.Skip("PostgreSQL not available for integration tests")
	}

	if err := tdb.SetUpProduct(); err != nil {
		t.Fatalf("Failed to set up product schema: %v", err)
	}
	defer func() {
		if err := tdb.TearDownProduct(); err != nil {
			t.Logf("Failed to tear down product schema: %v", err)
		}
	}()

	// Repositories
	productRepo := persistence.NewGORMProductRepository(tdb.DB)
	brandRepo := persistence.NewGORMBrandRepository(tdb.DB)
	groupRepo := persistence.NewGORMProductGroupRepository(tdb.DB)
	attributeRepo := persistence.NewGORMAttributeRepository(tdb.DB)
	variantRepo := persistence.NewGORMVariantRepository(tdb.DB)

	// Service
	productService := application.NewProductService(productRepo, brandRepo, groupRepo, attributeRepo, variantRepo)

	ctx := context.Background()

	// Create a brand and a group for the test
	brand := &domain.Brand{ID: uuid.New(), Name: "Test Brand"}
	group := &domain.ProductGroup{ID: uuid.New(), Name: "Test Group"}

	// Save the brand and group to the database
	brandDataModel := persistence.BrandFromDomain(brand)
	if err := tdb.DB.Create(brandDataModel).Error; err != nil {
		t.Fatalf("Failed to create brand: %v", err)
	}
	groupDataModel := persistence.ProductGroupFromDomain(group)
	if err := tdb.DB.Create(groupDataModel).Error; err != nil {
		t.Fatalf("Failed to create group: %v", err)
	}

	barcode := "1234567890123"
	cmd := application.CreateProductCommand{
		SKU:         "INTEGRATION-TEST-SKU",
		Name:        "Integration Test Product",
		LongName:    "Long Name for Integration Test Product",
		Barcode:     &barcode,
		Description: "Description for integration test product",
		ProductType: domain.ProductTypeTangible,
		BrandID:     brand.ID,
		GroupIDs:    []uuid.UUID{group.ID},
	}

	productDTO, err := productService.CreateProduct(ctx, cmd)

	assert.NoError(t, err)
	assert.NotNil(t, productDTO)
	assert.Equal(t, cmd.SKU, productDTO.SKU)
	assert.Equal(t, cmd.Name, productDTO.Name)
	assert.Equal(t, cmd.BrandID, productDTO.BrandID)
	assert.Len(t, productDTO.GroupIDs, 1)

	// Verify that the product was actually saved in the database
	savedProduct, err := productRepo.FindByID(ctx, productDTO.ID)
	assert.NoError(t, err)
	assert.NotNil(t, savedProduct)
	assert.Equal(t, cmd.SKU, savedProduct.SKU)
}

		assert.Nil(t, attributeDTO)
	})
}

		assert.Nil(t, savedAttr.ScopeBrandID)
	})
}

func TestProductService_GetApplicableAttributesForProduct_Integration(t *testing.T) {
	tdb := persistence.NewTestDB(t)
	if tdb.DB == nil {
		t.Skip("PostgreSQL not available for integration tests")
	}

	if err := tdb.SetUpProduct(); err != nil {
		t.Fatalf("Failed to set up product schema: %v", err)
	}
	defer func() {
		if err := tdb.TearDownProduct(); err != nil {
			t.Logf("Failed to tear down product schema: %v", err)
		}
	}()

	// Repositories
	productRepo := persistence.NewGORMProductRepository(tdb.DB)
	brandRepo := persistence.NewGORMBrandRepository(tdb.DB)
	groupRepo := persistence.NewGORMProductGroupRepository(tdb.DB)
	attributeRepo := persistence.NewGORMAttributeRepository(tdb.DB)
	variantRepo := persistence.NewGORMVariantRepository(tdb.DB)

	// Service
	productService := application.NewProductService(productRepo, brandRepo, groupRepo, attributeRepo, variantRepo)

	ctx := context.Background()

	// --- Helper functions ---
	createBrand := func(name string) *domain.Brand {
		brand := &domain.Brand{ID: uuid.New(), Name: name}
		assert.NoError(t, tdb.DB.Create(persistence.BrandFromDomain(brand)).Error)
		return brand
	}

	createGroup := func(name string) *domain.ProductGroup {
		group := &domain.ProductGroup{ID: uuid.New(), Name: name}
		assert.NoError(t, tdb.DB.Create(persistence.ProductGroupFromDomain(group)).Error)
		return group
	}

	createProduct := func(brandID uuid.UUID, groupIDs []uuid.UUID, directAttrIDs []uuid.UUID) *domain.Product {
		barcode := "12345"
		cmd := application.CreateProductCommand{
			SKU:         "PROD-" + uuid.New().String()[:4],
			Name:        "Test Product",
			LongName:    "Long Test Product",
			Barcode:     &barcode,
			Description: "Description",
			ProductType: domain.ProductTypeTangible,
			BrandID:     brandID,
			GroupIDs:    groupIDs,
		}
		productDTO, err := productService.CreateProduct(ctx, cmd)
		assert.NoError(t, err)

		// Add direct attributes if any
		for _, attrID := range directAttrIDs {
			_, err := productService.AddDirectAttributeToProduct(ctx, application.AddDirectAttributeCommand{
				ProductID:   productDTO.ID,
				AttributeID: attrID,
			})
			assert.NoError(t, err)
		}

		product, err := productRepo.FindByID(ctx, productDTO.ID)
		assert.NoError(t, err)
		return product
	}

	createAttribute := func(name, code string, sortOrder int, scopeBrandID, scopeGroupID *uuid.UUID, values []application.CreateAttributeValueCommand) *domain.Attribute {
		cmd := application.CreateAttributeCommand{
			Name:         name,
			Code:         code,
			SortOrder:    sortOrder,
			ScopeBrandID: scopeBrandID,
			ScopeGroupID: scopeGroupID,
			Values:       values,
		}
		attrDTO, err := productService.CreateAttribute(ctx, cmd)
		assert.NoError(t, err)
		attr, err := attributeRepo.FindByID(ctx, attrDTO.ID)
		assert.NoError(t, err)
		return attr
	}

	// --- Test Scenarios ---

	t.Run("Product with only generic attributes", func(t *testing.T) {
		brand := createBrand("Brand A")
		product := createProduct(brand.ID, nil, nil)
		attrGenericColor := createAttribute("Color", "COL", 1, nil, nil, []application.CreateAttributeValueCommand{{Value: "Red", Code: "R"}})
		createAttribute("Size", "SIZ", 2, nil, nil, []application.CreateAttributeValueCommand{{Value: "S", Code: "S"}}) // Another generic

		applicableAttrs, err := productService.GetApplicableAttributesForProduct(ctx, product.ID)
		assert.NoError(t, err)
		assert.Len(t, applicableAttrs, 2)
		assert.Equal(t, attrGenericColor.Code, applicableAttrs[0].Code) // Sorted by SortOrder
	})

	t.Run("Product with brand-scoped attributes overriding generic", func(t *testing.T) {
		brandB := createBrand("Brand B")
		product := createProduct(brandB.ID, nil, nil)
		attrGenericSize := createAttribute("Size", "SIZ", 1, nil, nil, []application.CreateAttributeValueCommand{{Value: "M", Code: "M"}})
		attrBrandSize := createAttribute("Size", "SIZ", 1, brandB.ID_PTR(), nil, []application.CreateAttributeValueCommand{{Value: "L", Code: "L"}})
		attrGenericColor := createAttribute("Color", "COL", 2, nil, nil, []application.CreateAttributeValueCommand{{Value: "Green", Code: "G"}})

		applicableAttrs, err := productService.GetApplicableAttributesForProduct(ctx, product.ID)
		assert.NoError(t, err)
		assert.Len(t, applicableAttrs, 2)
		// Brand-scoped SIZ should override generic SIZ
		assert.Equal(t, attrBrandSize.Code, applicableAttrs[0].Code)
		assert.Equal(t, attrBrandSize.ID, applicableAttrs[0].ID)
		assert.Equal(t, "L", applicableAttrs[0].Values[0].Value) // Check value from brand-scoped
		assert.Equal(t, attrGenericColor.Code, applicableAttrs[1].Code)
	})

	t.Run("Product with group-scoped attributes overriding generic", func(t *testing.T) {
		brandC := createBrand("Brand C")
		groupC := createGroup("Group C")
		product := createProduct(brandC.ID, []uuid.UUID{groupC.ID}, nil)
		attrGenericWeight := createAttribute("Weight", "WGT", 1, nil, nil, []application.CreateAttributeValueCommand{{Value: "1kg", Code: "1K"}})
		attrGroupWeight := createAttribute("Weight", "WGT", 1, nil, groupC.ID_PTR(), []application.CreateAttributeValueCommand{{Value: "2kg", Code: "2K"}})
		attrGenericType := createAttribute("Type", "TYP", 2, nil, nil, []application.CreateAttributeValueCommand{{Value: "TypeA", Code: "TA"}})

		applicableAttrs, err := productService.GetApplicableAttributesForProduct(ctx, product.ID)
		assert.NoError(t, err)
		assert.Len(t, applicableAttrs, 2)
		// Group-scoped WGT should override generic WGT
		assert.Equal(t, attrGroupWeight.Code, applicableAttrs[0].Code)
		assert.Equal(t, attrGroupWeight.ID, applicableAttrs[0].ID)
		assert.Equal(t, "2kg", applicableAttrs[0].Values[0].Value) // Check value from group-scoped
		assert.Equal(t, attrGenericType.Code, applicableAttrs[1].Code)
	})

	t.Run("Product with direct attributes overriding all others", func(t *testing.T) {
		brandD := createBrand("Brand D")
		groupD := createGroup("Group D")
		attrGenericSize := createAttribute("Size", "SIZ", 1, nil, nil, []application.CreateAttributeValueCommand{{Value: "XS", Code: "XS"}})
		attrBrandSize := createAttribute("Size", "SIZ", 1, brandD.ID_PTR(), nil, []application.CreateAttributeValueCommand{{Value: "S", Code: "S"}})
		attrGroupSize := createAttribute("Size", "SIZ", 1, nil, groupD.ID_PTR(), []application.CreateAttributeValueCommand{{Value: "M", Code: "M"}})
		attrDirectSize := createAttribute("Size", "SIZ", 1, nil, nil, []application.CreateAttributeValueCommand{{Value: "XL", Code: "XL"}}) // Direct attributes can also be generic
		product := createProduct(brandD.ID, []uuid.UUID{groupD.ID}, []uuid.UUID{attrDirectSize.ID})

		applicableAttrs, err := productService.GetApplicableAttributesForProduct(ctx, product.ID)
		assert.NoError(t, err)
		assert.Len(t, applicableAttrs, 1)
		// Direct SIZ should override all others
		assert.Equal(t, attrDirectSize.Code, applicableAttrs[0].Code)
		assert.Equal(t, attrDirectSize.ID, applicableAttrs[0].ID)
		assert.Equal(t, "XL", applicableAttrs[0].Values[0].Value)
	})

	t.Run("Product with Group+Brand scoped attribute overriding less specific", func(t *testing.T) {
		brandE := createBrand("Brand E")
		groupE := createGroup("Group E")
		product := createProduct(brandE.ID, []uuid.UUID{groupE.ID}, nil)

		attrGenericFit := createAttribute("Fit", "FIT", 1, nil, nil, []application.CreateAttributeValueCommand{{Value: "Loose", Code: "L"}})
		attrBrandFit := createAttribute("Fit", "FIT", 1, brandE.ID_PTR(), nil, []application.CreateAttributeValueCommand{{Value: "Regular", Code: "R"}})
		attrGroupFit := createAttribute("Fit", "FIT", 1, nil, groupE.ID_PTR(), []application.CreateAttributeValueCommand{{Value: "Slim", Code: "S"}})
		attrGroupBrandFit := createAttribute("Fit", "FIT", 1, brandE.ID_PTR(), groupE.ID_PTR(), []application.CreateAttributeValueCommand{{Value: "Athletic", Code: "A"}})

		applicableAttrs, err := productService.GetApplicableAttributesForProduct(ctx, product.ID)
		assert.NoError(t, err)
		assert.Len(t, applicableAttrs, 1)
		// Group+Brand scoped FIT should override all others
		assert.Equal(t, attrGroupBrandFit.Code, applicableAttrs[0].Code)
		assert.Equal(t, attrGroupBrandFit.ID, applicableAttrs[0].ID)
		assert.Equal(t, "Athletic", applicableAttrs[0].Values[0].Value)
	})

	t.Run("No applicable attributes", func(t *testing.T) {
		brandF := createBrand("Brand F")
		groupF := createGroup("Group F")
		product := createProduct(brandF.ID, []uuid.UUID{groupF.ID}, nil)
		// Create attributes that do not match the product's brand/group
		createAttribute("Color", "COL", 1, createBrand("Other Brand").ID_PTR(), nil, []application.CreateAttributeValueCommand{{Value: "Red", Code: "R"}})
		createAttribute("Size", "SIZ", 1, nil, createGroup("Other Group").ID_PTR(), []application.CreateAttributeValueCommand{{Value: "S", Code: "S"}})

		applicableAttrs, err := productService.GetApplicableAttributesForProduct(ctx, product.ID)
		assert.NoError(t, err)
		assert.Empty(t, applicableAttrs)
	})

	t.Run("Product with multiple groups - correct precedence", func(t *testing.T) {
		brandG := createBrand("Brand G")
		groupG1 := createGroup("Group G1")
		groupG2 := createGroup("Group G2")
		product := createProduct(brandG.ID, []uuid.UUID{groupG1.ID, groupG2.ID}, nil)

		attrGenericMaterial := createAttribute("Material", "MAT", 1, nil, nil, []application.CreateAttributeValueCommand{{Value: "Cotton", Code: "C"}})
		attrGroupG1Material := createAttribute("Material", "MAT", 1, nil, groupG1.ID_PTR(), []application.CreateAttributeValueCommand{{Value: "Wool", Code: "W"}})
		attrGroupG2Material := createAttribute("Material", "MAT", 1, nil, groupG2.ID_PTR(), []application.CreateAttributeValueCommand{{Value: "Linen", Code: "L"}}) // Should not be picked if G1 is found first
		attrDirectMaterial := createAttribute("Material", "MAT", 1, nil, nil, []application.CreateAttributeValueCommand{{Value: "Silk", Code: "S"}})

		// Add attrDirectMaterial to product AFTER creating it
		product.AddDirectAttribute(attrDirectMaterial.ID)
		productRepo.Save(ctx, product) // Ensure updated product is saved

		applicableAttrs, err := productService.GetApplicableAttributesForProduct(ctx, product.ID)
		assert.NoError(t, err)
		assert.Len(t, applicableAttrs, 1)
		// Direct attribute should win
		assert.Equal(t, attrDirectMaterial.Code, applicableAttrs[0].Code)
		assert.Equal(t, attrDirectMaterial.ID, applicableAttrs[0].ID)
		assert.Equal(t, "Silk", applicableAttrs[0].Values[0].Value)
	})
}

func TestProductService_AddDirectAttributeToProduct_Integration(t *testing.T) {
	tdb := persistence.NewTestDB(t)
	if tdb.DB == nil {
		t.Skip("PostgreSQL not available for integration tests")
	}

	if err := tdb.SetUpProduct(); err != nil {
		t.Fatalf("Failed to set up product schema: %v", err)
	}
	defer func() {
		if err := tdb.TearDownProduct(); err != nil {
			t.Logf("Failed to tear down product schema: %v", err)
		}
	}()

	// Repositories
	productRepo := persistence.NewGORMProductRepository(tdb.DB)
	brandRepo := persistence.NewGORMBrandRepository(tdb.DB)
	groupRepo := persistence.NewGORMProductGroupRepository(tdb.DB)
	attributeRepo := persistence.NewGORMAttributeRepository(tdb.DB)
	variantRepo := persistence.NewGORMVariantRepository(tdb.DB)

	// Service
	productService := application.NewProductService(productRepo, brandRepo, groupRepo, attributeRepo, variantRepo)

	ctx := context.Background()

	// --- Helper functions ---
	createBrand := func(name string) *domain.Brand {
		brand := &domain.Brand{ID: uuid.New(), Name: name}
		assert.NoError(t, tdb.DB.Create(persistence.BrandFromDomain(brand)).Error)
		return brand
	}

	createProduct := func(brandID uuid.UUID, groupIDs []uuid.UUID) *domain.Product {
		barcode := "12345"
		cmd := application.CreateProductCommand{
			SKU:         "PROD-" + uuid.New().String()[:4],
			Name:        "Test Product",
			LongName:    "Long Test Product",
			Barcode:     &barcode,
			Description: "Description",
			ProductType: domain.ProductTypeTangible,
			BrandID:     brandID,
			GroupIDs:    groupIDs,
		}
		productDTO, err := productService.CreateProduct(ctx, cmd)
		assert.NoError(t, err)
		product, err := productRepo.FindByID(ctx, productDTO.ID)
		assert.NoError(t, err)
		return product
	}

	createAttribute := func(name, code string) *domain.Attribute {
		cmd := application.CreateAttributeCommand{
			Name:      name,
			Code:      code,
			SortOrder: 1,
			Values:    []application.CreateAttributeValueCommand{},
		}
		attrDTO, err := productService.CreateAttribute(ctx, cmd)
		assert.NoError(t, err)
		attr, err := attributeRepo.FindByID(ctx, attrDTO.ID)
		assert.NoError(t, err)
		return attr
	}

	t.Run("Successfully add a direct attribute to a product", func(t *testing.T) {
		brand := createBrand("Brand A")
		product := createProduct(brand.ID, nil)
		attribute := createAttribute("Color", "COL")

		cmd := application.AddDirectAttributeCommand{
			ProductID:   product.ID,
			AttributeID: attribute.ID,
		}

		productDTO, err := productService.AddDirectAttributeToProduct(ctx, cmd)
		assert.NoError(t, err)
		assert.NotNil(t, productDTO)
		assert.Contains(t, productDTO.DirectAttributeIDs, attribute.ID)

		// Verify in DB
		savedProduct, err := productRepo.FindByID(ctx, product.ID)
		assert.NoError(t, err)
		assert.NotNil(t, savedProduct)
		assert.Contains(t, savedProduct.DirectAttributeIDs, attribute.ID)
	})

	t.Run("Adding the same direct attribute twice should not add duplicate", func(t *testing.T) {
		brand := createBrand("Brand B")
		product := createProduct(brand.ID, nil)
		attribute := createAttribute("Size", "SIZ")

		// Add once
		_, err := productService.AddDirectAttributeToProduct(ctx, application.AddDirectAttributeCommand{
			ProductID:   product.ID,
			AttributeID: attribute.ID,
		})
		assert.NoError(t, err)

		// Add again
		_, err = productService.AddDirectAttributeToProduct(ctx, application.AddDirectAttributeCommand{
			ProductID:   product.ID,
			AttributeID: attribute.ID,
		})
		assert.NoError(t, err) // Should still return no error, just no change

		savedProduct, err := productRepo.FindByID(ctx, product.ID)
		assert.NoError(t, err)
		assert.NotNil(t, savedProduct)
		assert.Len(t, savedProduct.DirectAttributeIDs, 1) // Should only have one entry
		assert.Contains(t, savedProduct.DirectAttributeIDs, attribute.ID)
	})

	t.Run("Adding direct attribute to a non-existent product should return error", func(t *testing.T) {
		nonExistentProductID := uuid.New()
		attribute := createAttribute("Material", "MAT")

		cmd := application.AddDirectAttributeCommand{
			ProductID:   nonExistentProductID,
			AttributeID: attribute.ID,
		}

		productDTO, err := productService.AddDirectAttributeToProduct(ctx, cmd)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "product with ID "+nonExistentProductID.String()+" does not exist")
		assert.Nil(t, productDTO)
	})

	t.Run("Adding non-existent direct attribute should return error", func(t *testing.T) {
		brand := createBrand("Brand C")
		product := createProduct(brand.ID, nil)
		nonExistentAttributeID := uuid.New()

		cmd := application.AddDirectAttributeCommand{
			ProductID:   product.ID,
			AttributeID: nonExistentAttributeID,
		}

		productDTO, err := productService.AddDirectAttributeToProduct(ctx, cmd)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "attribute with ID "+nonExistentAttributeID.String()+" does not exist")
		assert.Nil(t, productDTO)
	})
}

func TestProductService_UpdateProductSKU_Integration(t *testing.T) {
	tdb := persistence.NewTestDB(t)
	if tdb.DB == nil {
		t.Skip("PostgreSQL not available for integration tests")
	}

	if err := tdb.SetUpProduct(); err != nil {
		t.Fatalf("Failed to set up product schema: %v", err)
	}
	defer func() {
		if err := tdb.TearDownProduct(); err != nil {
			t.Logf("Failed to tear down product schema: %v", err)
		}
	}()

	// Repositories
	productRepo := persistence.NewGORMProductRepository(tdb.DB)
	brandRepo := persistence.NewGORMBrandRepository(tdb.DB)
	groupRepo := persistence.NewGORMProductGroupRepository(tdb.DB)
	attributeRepo := persistence.NewGORMAttributeRepository(tdb.DB)
	variantRepo := persistence.NewGORMVariantRepository(tdb.DB)

	// Service
	productService := application.NewProductService(productRepo, brandRepo, groupRepo, attributeRepo, variantRepo)

	ctx := context.Background()

	// --- Helper functions ---
	createBrand := func(name string) *domain.Brand {
		brand := &domain.Brand{ID: uuid.New(), Name: name}
		assert.NoError(t, tdb.DB.Create(persistence.BrandFromDomain(brand)).Error)
		return brand
	}

	createProduct := func(brandID uuid.UUID, groupIDs []uuid.UUID) *domain.Product {
		barcode := "12345"
		cmd := application.CreateProductCommand{
			SKU:         "PROD-" + uuid.New().String()[:4],
			Name:        "Test Product",
			LongName:    "Long Test Product",
			Barcode:     &barcode,
			Description: "Description",
			ProductType: domain.ProductTypeTangible,
			BrandID:     brandID,
			GroupIDs:    groupIDs,
		}
		productDTO, err := productService.CreateProduct(ctx, cmd)
		assert.NoError(t, err)
		product, err := productRepo.FindByID(ctx, productDTO.ID)
		assert.NoError(t, err)
		return product
	}

	t.Run("Successfully update product SKU", func(t *testing.T) {
		brand := createBrand("Brand A")
		product := createProduct(brand.ID, nil)

		newSKU := "UPDATED-SKU"
		cmd := application.UpdateProductSKUCommand{
			ProductID: product.ID,
			NewSKU:    newSKU,
		}

		productDTO, err := productService.UpdateProductSKU(ctx, cmd)
		assert.NoError(t, err)
		assert.NotNil(t, productDTO)
		assert.Equal(t, newSKU, productDTO.SKU)

		// Verify in DB
		savedProduct, err := productRepo.FindByID(ctx, product.ID)
		assert.NoError(t, err)
		assert.NotNil(t, savedProduct)
		assert.Equal(t, newSKU, savedProduct.SKU)
		// TODO: Add ProductVariant cascade test once variant creation is implemented.
	})

	t.Run("Updating SKU of a non-existent product should return error", func(t *testing.T) {
		nonExistentProductID := uuid.New()
		newSKU := "NON-EXISTENT-SKU"

		cmd := application.UpdateProductSKUCommand{
			ProductID: nonExistentProductID,
			NewSKU:    newSKU,
		}

		productDTO, err := productService.UpdateProductSKU(ctx, cmd)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "product with ID "+nonExistentProductID.String()+" does not exist")
		assert.Nil(t, productDTO)
	})
}
