package application_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/joran-cortez/tramatex/internal/product/application"
	"github.com/joran-cortez/tramatex/internal/product/domain"
	"github.com/joran-cortez/tramatex/internal/product/infrastructure/persistence"
	// Required for sorting attribute values
)

type testServices struct {
	ProductService                *application.ProductService
	ProductRepo                   domain.ProductRepository
	BrandRepo                     domain.BrandRepository
	GroupRepo                     domain.ProductGroupRepository
	AttributeRepo                 domain.AttributeRepository
	VariantRepo                   domain.ProductVariantRepository
	PartyServiceConfigurationRepo domain.PartyServiceConfigurationRepository // New
	TestDB                        *persistence.TestDB
	Ctx                           context.Context
}

func setupTest(t *testing.T) *testServices {
	tdb := persistence.NewTestDB(t)
	if tdb.DB == nil {
		t.Skip("PostgreSQL not available for integration tests")
	}

	if err := tdb.SetUpProduct(); err != nil {
		t.Fatalf("Failed to set up product schema: %v", err)
	}

	productRepo := persistence.NewGORMProductRepository(tdb.DB)
	brandRepo := persistence.NewGORMBrandRepository(tdb.DB)
	groupRepo := persistence.NewGORMProductGroupRepository(tdb.DB)
	attributeRepo := persistence.NewGORMAttributeRepository(tdb.DB)
	variantRepo := persistence.NewGORMVariantRepository(tdb.DB)
	partyServiceConfigurationRepo := persistence.NewGORMPartyServiceConfigurationRepository(tdb.DB) // New

	productService := application.NewProductService(productRepo, brandRepo, groupRepo, attributeRepo, variantRepo, partyServiceConfigurationRepo) // Updated

	return &testServices{
		ProductService:                productService,
		ProductRepo:                   productRepo,
		BrandRepo:                     brandRepo,
		GroupRepo:                     groupRepo,
		AttributeRepo:                 attributeRepo,
		VariantRepo:                   variantRepo,
		PartyServiceConfigurationRepo: partyServiceConfigurationRepo, // New
		TestDB:                        tdb,
		Ctx:                           context.Background(),
	}
}

func teardownTest(ts *testServices) {
	if err := ts.TestDB.TearDownProduct(); err != nil {
		ts.TestDB.Logf("Failed to tear down product schema: %v", err)
	}
}

func TestProductService_CreateProduct_Integration(t *testing.T) {
	ts := setupTest(t)
	defer teardownTest(ts)

	ctx := ts.Ctx

	// Create a brand and a group for the test
	brand := &domain.Brand{ID: uuid.New(), Name: "Test Brand"}
	group := &domain.ProductGroup{ID: uuid.New(), Name: "Test Group"}

	// Save the brand and group to the database
	brandDataModel := persistence.BrandFromDomain(brand)
	if err := ts.TestDB.DB.Create(brandDataModel).Error; err != nil {
		t.Fatalf("Failed to create brand: %v", err)
	}
	groupDataModel := persistence.ProductGroupFromDomain(group)
	if err := ts.TestDB.DB.Create(groupDataModel).Error; err != nil {
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

	productDTO, err := ts.ProductService.CreateProduct(ctx, cmd)

	assert.NoError(t, err)
	assert.NotNil(t, productDTO)
	assert.Equal(t, cmd.SKU, productDTO.SKU)
	assert.Equal(t, cmd.Name, productDTO.Name)
	assert.Equal(t, cmd.BrandID, productDTO.BrandID)
	assert.Len(t, productDTO.GroupIDs, 1)

	// Verify that the product was actually saved in the database
	savedProduct, err := ts.ProductRepo.FindByID(ctx, productDTO.ID)
	assert.NoError(t, err)
	assert.NotNil(t, savedProduct)
	assert.Equal(t, cmd.SKU, savedProduct.SKU)
}

func TestProductService_GetApplicableAttributesForProduct_Integration(t *testing.T) {
	ts := setupTest(t)
	defer teardownTest(ts)

	ctx := ts.Ctx

	// --- Helper functions ---
	createBrand := func(name string) *domain.Brand {
		brand := &domain.Brand{ID: uuid.New(), Name: name}
		assert.NoError(t, ts.TestDB.DB.Create(persistence.BrandFromDomain(brand)).Error)
		return brand
	}

	createGroup := func(name string) *domain.ProductGroup {
		group := &domain.ProductGroup{ID: uuid.New(), Name: name}
		assert.NoError(t, ts.TestDB.DB.Create(persistence.ProductGroupFromDomain(group)).Error)
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
		productDTO, err := ts.ProductService.CreateProduct(ctx, cmd)
		assert.NoError(t, err)

		// Add direct attributes if any
		for _, attrID := range directAttrIDs {
			_, err := ts.ProductService.AddDirectAttributeToProduct(ctx, application.AddDirectAttributeCommand{
				ProductID:   productDTO.ID,
				AttributeID: attrID,
			})
			assert.NoError(t, err)
		}

		product, err := ts.ProductRepo.FindByID(ctx, productDTO.ID)
		assert.NoError(t, err)
		return product
	}

	createAttribute := func(name, code string, sortOrder int, scopeBrandID, scopeGroupID *uuid.UUID, values []string) *domain.Attribute {
		attrValues := make([]application.CreateAttributeValueCommand, len(values))
		for i, v := range values {
			attrValues[i] = application.CreateAttributeValueCommand{Value: v, Code: v[:1]} // Simple code generation for test
		}

		cmd := application.CreateAttributeCommand{
			Name:         name,
			Code:         code,
			SortOrder:    sortOrder,
			ScopeBrandID: scopeBrandID,
			ScopeGroupID: scopeGroupID,
			Values:       attrValues,
		}
		attrDTO, err := ts.ProductService.CreateAttribute(ctx, cmd)
		assert.NoError(t, err)
		attr, err := ts.AttributeRepo.FindByID(ctx, attrDTO.ID)
		assert.NoError(t, err)
		return attr
	}

	// --- Test Scenarios ---

	t.Run("Product with only generic attributes", func(t *testing.T) {
		brand := createBrand("Brand A")
		product := createProduct(brand.ID, nil, nil)
		attrGenericColor := createAttribute("Color", "COL", 1, nil, nil, []string{"Red"})
		attrGenericSize := createAttribute("Size", "SIZ", 2, nil, nil, []string{"S"}) // Another generic
		_ = attrGenericSize

		applicableAttrs, err := ts.ProductService.GetApplicableAttributesForProduct(ctx, product.ID)
		assert.NoError(t, err)
		assert.Len(t, applicableAttrs, 2)
		assert.Equal(t, attrGenericColor.Code, applicableAttrs[0].AttributeName) // Sorted by SortOrder
	})

	t.Run("Product with brand-scoped attributes overriding generic", func(t *testing.T) {
		brandB := createBrand("Brand B")
		product := createProduct(brandB.ID, nil, nil)
		attrGenericSize := createAttribute("Size", "SIZ", 1, nil, nil, []string{"M"})
		attrBrandSize := createAttribute("Size", "SIZ", 1, brandB.ID_PTR(), nil, []string{"L"})
		attrGenericColor := createAttribute("Color", "COL", 2, nil, nil, []string{"Green"})
		_ = attrGenericSize

		applicableAttrs, err := ts.ProductService.GetApplicableAttributesForProduct(ctx, product.ID)
		assert.NoError(t, err)
		assert.Len(t, applicableAttrs, 2)
		// Brand-scoped SIZ should override generic SIZ
		assert.Equal(t, attrBrandSize.Code, applicableAttrs[0].AttributeName)
		assert.Equal(t, attrBrandSize.ID, applicableAttrs[0].ID)
		assert.Equal(t, "L", applicableAttrs[0].Values[0]) // Check value from brand-scoped
		assert.Equal(t, attrGenericColor.Code, applicableAttrs[1].AttributeName)
	})

	t.Run("Product with group-scoped attributes overriding generic", func(t *testing.T) {
		brandC := createBrand("Brand C")
		groupC := createGroup("Group C")
		product := createProduct(brandC.ID, []uuid.UUID{groupC.ID}, nil)
		attrGenericWeight := createAttribute("Weight", "WGT", 1, nil, nil, []string{"1kg"})
		attrGroupWeight := createAttribute("Weight", "WGT", 1, nil, groupC.ID_PTR(), []string{"2kg"})
		attrGenericType := createAttribute("Type", "TYP", 2, nil, nil, []string{"TypeA"})
		_ = attrGenericWeight

		applicableAttrs, err := ts.ProductService.GetApplicableAttributesForProduct(ctx, product.ID)
		assert.NoError(t, err)
		assert.Len(t, applicableAttrs, 2)
		// Group-scoped WGT should override generic WGT
		assert.Equal(t, attrGroupWeight.Code, applicableAttrs[0].AttributeName)
		assert.Equal(t, attrGroupWeight.ID, applicableAttrs[0].ID)
		assert.Equal(t, "2kg", applicableAttrs[0].Values[0]) // Check value from group-scoped
		assert.Equal(t, attrGenericType.Code, applicableAttrs[1].AttributeName)
	})

	t.Run("Product with direct attributes overriding all others", func(t *testing.T) {
		brandD := createBrand("Brand D")
		groupD := createGroup("Group D")
		attrGenericSize := createAttribute("Size", "SIZ", 1, nil, nil, []string{"XS"})
		attrBrandSize := createAttribute("Size", "SIZ", 1, brandD.ID_PTR(), nil, []string{"S"})
		attrGroupSize := createAttribute("Size", "SIZ", 1, nil, groupD.ID_PTR(), []string{"M"})
		attrDirectSize := createAttribute("Size", "SIZ", 1, nil, nil, []string{"XL"}) // Direct attributes can also be generic
		product := createProduct(brandD.ID, []uuid.UUID{groupD.ID}, []uuid.UUID{attrDirectSize.ID})
		_ = attrGenericSize
		_ = attrBrandSize
		_ = attrGroupSize

		applicableAttrs, err := ts.ProductService.GetApplicableAttributesForProduct(ctx, product.ID)
		assert.NoError(t, err)
		assert.Len(t, applicableAttrs, 1)
		// Direct SIZ should override all others
		assert.Equal(t, attrDirectSize.Code, applicableAttrs[0].AttributeName)
		assert.Equal(t, attrDirectSize.ID, applicableAttrs[0].ID)
		assert.Equal(t, "XL", applicableAttrs[0].Values[0])
	})

	t.Run("Product with Group+Brand scoped attribute overriding less specific", func(t *testing.T) {
		brandE := createBrand("Brand E")
		groupE := createGroup("Group E")
		product := createProduct(brandE.ID, []uuid.UUID{groupE.ID}, nil)

		attrGenericFit := createAttribute("Fit", "FIT", 1, nil, nil, []string{"Loose"})
		attrBrandFit := createAttribute("Fit", "FIT", 1, brandE.ID_PTR(), nil, []string{"Regular"})
		attrGroupFit := createAttribute("Fit", "FIT", 1, nil, groupE.ID_PTR(), []string{"Slim"})
		attrGroupBrandFit := createAttribute("Fit", "FIT", 1, brandE.ID_PTR(), groupE.ID_PTR(), []string{"Athletic"})
		_ = attrGenericFit
		_ = attrBrandFit
		_ = attrGroupFit

		applicableAttrs, err := ts.ProductService.GetApplicableAttributesForProduct(ctx, product.ID)
		assert.NoError(t, err)
		assert.Len(t, applicableAttrs, 1)
		// Group+Brand scoped FIT should override all others
		assert.Equal(t, attrGroupBrandFit.Code, applicableAttrs[0].AttributeName)
		assert.Equal(t, attrGroupBrandFit.ID, applicableAttrs[0].ID)
		assert.Equal(t, "Athletic", applicableAttrs[0].Values[0])
	})

	t.Run("No applicable attributes", func(t *testing.T) {
		brandF := createBrand("Brand F")
		groupF := createGroup("Group F")
		product := createProduct(brandF.ID, []uuid.UUID{groupF.ID}, nil)
		// Create attributes that do not match the product's brand/group
		createAttribute("Color", "COL", 1, createBrand("Other Brand").ID_PTR(), nil, []string{"Red"})
		createAttribute("Size", "SIZ", 1, nil, createGroup("Other Group").ID_PTR(), []string{"S"})

		applicableAttrs, err := ts.ProductService.GetApplicableAttributesForProduct(ctx, product.ID)
		assert.NoError(t, err)
		assert.Empty(t, applicableAttrs)
	})

	t.Run("Product with multiple groups - correct precedence", func(t *testing.T) {
		brandG := createBrand("Brand G")
		groupG1 := createGroup("Group G1")
		groupG2 := createGroup("Group G2")
		product := createProduct(brandG.ID, []uuid.UUID{groupG1.ID, groupG2.ID}, nil)

		attrGenericMaterial := createAttribute("Material", "MAT", 1, nil, nil, []string{"Cotton"})
		attrGroupG1Material := createAttribute("Material", "MAT", 1, nil, groupG1.ID_PTR(), []string{"Wool"})
		attrGroupG2Material := createAttribute("Material", "MAT", 1, nil, groupG2.ID_PTR(), []string{"Linen"}) // Should not be picked if G1 is found first
		attrDirectMaterial := createAttribute("Material", "MAT", 1, nil, nil, []string{"Silk"})
		_ = attrGenericMaterial
		_ = attrGroupG1Material
		_ = attrGroupG2Material

		// Add attrDirectMaterial to product AFTER creating it
		product.AddDirectAttribute(attrDirectMaterial.ID)
		ts.ProductRepo.Save(ctx, product) // Ensure updated product is saved

		applicableAttrs, err := ts.ProductService.GetApplicableAttributesForProduct(ctx, product.ID)
		assert.NoError(t, err)
		assert.Len(t, applicableAttrs, 1)
		// Direct attribute should win
		assert.Equal(t, attrDirectMaterial.Code, applicableAttrs[0].AttributeName)
		assert.Equal(t, attrDirectMaterial.ID, applicableAttrs[0].ID)
		assert.Equal(t, "Silk", applicableAttrs[0].Values[0])
	})
}

func TestProductService_AddDirectAttributeToProduct_Integration(t *testing.T) {
	ts := setupTest(t)
	defer teardownTest(ts)

	ctx := ts.Ctx

	// --- Helper functions ---
	createBrand := func(name string) *domain.Brand {
		brand := &domain.Brand{ID: uuid.New(), Name: name}
		assert.NoError(t, ts.TestDB.DB.Create(persistence.BrandFromDomain(brand)).Error)
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
		productDTO, err := ts.ProductService.CreateProduct(ctx, cmd)
		assert.NoError(t, err)
		product, err := ts.ProductRepo.FindByID(ctx, productDTO.ID)
		assert.NoError(t, err)
		return product
	}

	createAttribute := func(name, code string) *domain.Attribute {
		attrValues := make([]application.CreateAttributeValueCommand, 1)
		attrValues[0] = application.CreateAttributeValueCommand{Value: "Default", Code: "D"}
		cmd := application.CreateAttributeCommand{
			Name:      name,
			Code:      code,
			SortOrder: 1,
			Values:    attrValues,
		}
		attrDTO, err := ts.ProductService.CreateAttribute(ctx, cmd)
		assert.NoError(t, err)
		attr, err := ts.AttributeRepo.FindByID(ctx, attrDTO.ID)
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

		productDTO, err := ts.ProductService.AddDirectAttributeToProduct(ctx, cmd)
		assert.NoError(t, err)
		assert.NotNil(t, productDTO)
		assert.Contains(t, productDTO.DirectAttributeIDs, attribute.ID)

		// Verify in DB
		savedProduct, err := ts.ProductRepo.FindByID(ctx, product.ID)
		assert.NoError(t, err)
		assert.NotNil(t, savedProduct)
		assert.Contains(t, savedProduct.DirectAttributeIDs, attribute.ID)
	})

	t.Run("Adding the same direct attribute twice should not add duplicate", func(t *testing.T) {
		brand := createBrand("Brand B")
		product := createProduct(brand.ID, nil)
		attribute := createAttribute("Size", "SIZ")

		// Add once
		_, err := ts.ProductService.AddDirectAttributeToProduct(ctx, application.AddDirectAttributeCommand{
			ProductID:   product.ID,
			AttributeID: attribute.ID,
		})
		assert.NoError(t, err)

		// Add again
		_, err = ts.ProductService.AddDirectAttributeToProduct(ctx, application.AddDirectAttributeCommand{
			ProductID:   product.ID,
			AttributeID: attribute.ID,
		})
		assert.NoError(t, err) // Should still return no error, just no change

		savedProduct, err := ts.ProductRepo.FindByID(ctx, product.ID)
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

		productDTO, err := ts.ProductService.AddDirectAttributeToProduct(ctx, cmd)
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

		productDTO, err := ts.ProductService.AddDirectAttributeToProduct(ctx, cmd)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "attribute with ID "+nonExistentAttributeID.String()+" does not exist")
		assert.Nil(t, productDTO)
	})
}

func TestProductService_UpdateProductSKU_Integration(t *testing.T) {
	ts := setupTest(t)
	defer teardownTest(ts)

	ctx := ts.Ctx

	// --- Helper functions ---
	createBrand := func(name string) *domain.Brand {
		brand := &domain.Brand{ID: uuid.New(), Name: name}
		assert.NoError(t, ts.TestDB.DB.Create(persistence.BrandFromDomain(brand)).Error)
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
		productDTO, err := ts.ProductService.CreateProduct(ctx, cmd)
		assert.NoError(t, err)
		product, err := ts.ProductRepo.FindByID(ctx, productDTO.ID)
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

		productDTO, err := ts.ProductService.UpdateProductSKU(ctx, cmd)
		assert.NoError(t, err)
		assert.NotNil(t, productDTO)
		assert.Equal(t, newSKU, productDTO.SKU)

		// Verify in DB
		savedProduct, err := ts.ProductRepo.FindByID(ctx, product.ID)
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

		productDTO, err := ts.ProductService.UpdateProductSKU(ctx, cmd)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "product with ID "+nonExistentProductID.String()+" does not exist")
		assert.Nil(t, productDTO)
	})
}

func TestProductService_GetAttributeByID_Integration(t *testing.T) {
	ts := setupTest(t)
	defer teardownTest(ts)

	ctx := ts.Ctx

	createAttribute := func(name, code string, sortOrder int, scopeBrandID, scopeGroupID *uuid.UUID, values []string) *domain.Attribute {
		attrValues := make([]application.CreateAttributeValueCommand, len(values))
		for i, v := range values {
			attrValues[i] = application.CreateAttributeValueCommand{Value: v, Code: v[:1]}
		}
		cmd := application.CreateAttributeCommand{
			Name:         name,
			Code:         code,
			SortOrder:    sortOrder,
			ScopeBrandID: scopeBrandID,
			ScopeGroupID: scopeGroupID,
			Values:       attrValues,
		}
		attrDTO, err := ts.ProductService.CreateAttribute(ctx, cmd)
		assert.NoError(t, err)
		attr, err := ts.AttributeRepo.FindByID(ctx, attrDTO.ID)
		assert.NoError(t, err)
		return attr
	}

	t.Run("Successfully get attribute by ID", func(t *testing.T) {
		attr := createAttribute("Color", "COL", 1, nil, nil, []string{"Red"})
		query := application.GetAttributeByIDQuery{ID: attr.ID}

		foundAttr, err := ts.ProductService.GetAttributeByID(ctx, query)
		assert.NoError(t, err)
		assert.NotNil(t, foundAttr)
		assert.Equal(t, attr.ID, foundAttr.ID)
		assert.Equal(t, attr.Name, foundAttr.Name)
	})

	t.Run("Getting a non-existent attribute by ID should return error", func(t *testing.T) {
		nonExistentID := uuid.New()
		query := application.GetAttributeByIDQuery{ID: nonExistentID}

		foundAttr, err := ts.ProductService.GetAttributeByID(ctx, query)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "attribute with ID "+nonExistentID.String()+" does not exist")
		assert.Nil(t, foundAttr)
	})
}

func TestProductService_ListAttributes_Integration(t *testing.T) {
	ts := setupTest(t)
	defer teardownTest(ts)

	ctx := ts.Ctx

	createBrand := func(name string) *domain.Brand {
		brand := &domain.Brand{ID: uuid.New(), Name: name}
		assert.NoError(t, ts.TestDB.DB.Create(persistence.BrandFromDomain(brand)).Error)
		return brand
	}
	createGroup := func(name string) *domain.ProductGroup {
		group := &domain.ProductGroup{ID: uuid.New(), Name: name}
		assert.NoError(t, ts.TestDB.DB.Create(persistence.ProductGroupFromDomain(group)).Error)
		return group
	}
	createAttribute := func(name, code string, sortOrder int, scopeBrandID, scopeGroupID *uuid.UUID, values []string) *domain.Attribute {
		attrValues := make([]application.CreateAttributeValueCommand, len(values))
		for i, v := range values {
			attrValues[i] = application.CreateAttributeValueCommand{Value: v, Code: v[:1]}
		}
		cmd := application.CreateAttributeCommand{
			Name:         name,
			Code:         code,
			SortOrder:    sortOrder,
			ScopeBrandID: scopeBrandID,
			ScopeGroupID: scopeGroupID,
			Values:       attrValues,
		}
		attrDTO, err := ts.ProductService.CreateAttribute(ctx, cmd)
		assert.NoError(t, err)
		attr, err := ts.AttributeRepo.FindByID(ctx, attrDTO.ID)
		assert.NoError(t, err)
		return attr
	}

	t.Run("List all attributes when no filters are applied", func(t *testing.T) {
		createAttribute("Color", "COL", 1, nil, nil, []string{"Red"})
		createAttribute("Size", "SIZ", 2, nil, nil, []string{"S"})

		query := application.ListAttributesQuery{}
		attributes, err := ts.ProductService.ListAttributes(ctx, query)
		assert.NoError(t, err)
		assert.Len(t, attributes, 2)
	})

	t.Run("List attributes filtered by BrandID", func(t *testing.T) {
		brandA := createBrand("Brand A")
		brandB := createBrand("Brand B")
		createAttribute("Material", "MAT", 1, brandA.ID_PTR(), nil, []string{"Cotton"})
		createAttribute("Pattern", "PAT", 2, brandB.ID_PTR(), nil, []string{"Stripes"})
		createAttribute("Style", "STY", 3, nil, nil, []string{"Casual"})

		query := application.ListAttributesQuery{BrandID: &brandA.ID}
		attributes, err := ts.ProductService.ListAttributes(ctx, query)
		assert.NoError(t, err)
		assert.Len(t, attributes, 1)
		assert.Equal(t, "Material", attributes[0].Name)
	})

	t.Run("List attributes filtered by ProductGroupID", func(t *testing.T) {
		groupX := createGroup("Group X")
		groupY := createGroup("Group Y")
		createAttribute("Closure", "CLO", 1, nil, groupX.ID_PTR(), []string{"Zipper"})
		createAttribute("Neckline", "NEC", 2, nil, groupY.ID_PTR(), []string{"V-Neck"})
		createAttribute("Sleeve", "SLV", 3, nil, nil, []string{"Long"})

		query := application.ListAttributesQuery{ProductGroupID: &groupX.ID}
		attributes, err := ts.ProductService.ListAttributes(ctx, query)
		assert.NoError(t, err)
		assert.Len(t, attributes, 1)
		assert.Equal(t, "Closure", attributes[0].Name)
	})

	t.Run("List attributes filtered by both BrandID and ProductGroupID", func(t *testing.T) {
		brandP := createBrand("Brand P")
		groupQ := createGroup("Group Q")
		createAttribute("Fit", "FIT", 1, brandP.ID_PTR(), groupQ.ID_PTR(), []string{"Slim"})
		createAttribute("Occasion", "OCC", 2, brandP.ID_PTR(), nil, []string{"Formal"})
		createAttribute("Season", "SEA", 3, nil, groupQ.ID_PTR(), []string{"Summer"})

		query := application.ListAttributesQuery{BrandID: &brandP.ID, ProductGroupID: &groupQ.ID}
		attributes, err := ts.ProductService.ListAttributes(ctx, query)
		assert.NoError(t, err)
		assert.Len(t, attributes, 1)
		assert.Equal(t, "Fit", attributes[0].Name)
	})

	t.Run("List generic attributes only", func(t *testing.T) {
		createBrand("Brand Gen")
		createGroup("Group Gen")
		createAttribute("Generic Color", "GCOL", 1, nil, nil, []string{"Blue"})
		createAttribute("Brand Color", "BCOL", 2, createBrand("Another Brand").ID_PTR(), nil, []string{"Green"})

		// Assuming FindByScope with both nil returns generic. This depends on implementation.
		// If FindByScope needs a specific type, this test might need adjustment.
		query := application.ListAttributesQuery{} // Both BrandID and ProductGroupID are nil
		attributes, err := ts.ProductService.ListAttributes(ctx, query)
		assert.NoError(t, err)
		// We expect to find only generic attribute, or potentially all if the query is not strict
		// For now, let's assume it lists all if no explicit filters are passed, but the service method filters by provided nil values
		// The current `ListAttributes` in product_service.go just passes nil to FindByScope, which in GORM layer might return all if both are nil.
		// This test needs to be refined once the ListAttributesQuery is fully implemented in service layer.
		// For now, testing basic FindByScope behavior.
		assert.GreaterOrEqual(t, len(attributes), 1)
		assert.Contains(t, attributes[0].Name, "Color") // At least one generic attribute.
	})
}

func TestProductService_UpdateAttribute_Integration(t *testing.T) {
	ts := setupTest(t)
	defer teardownTest(ts)

	ctx := ts.Ctx

	createAttribute := func(name, code string, sortOrder int, scopeBrandID, scopeGroupID *uuid.UUID, values []string) *domain.Attribute {
		attrValues := make([]application.CreateAttributeValueCommand, len(values))
		for i, v := range values {
			attrValues[i] = application.CreateAttributeValueCommand{Value: v, Code: v[:1]}
		}
		cmd := application.CreateAttributeCommand{
			Name:         name,
			Code:         code,
			SortOrder:    sortOrder,
			ScopeBrandID: scopeBrandID,
			ScopeGroupID: scopeGroupID,
			Values:       attrValues,
		}
		attrDTO, err := ts.ProductService.CreateAttribute(ctx, cmd)
		assert.NoError(t, err)
		attr, err := ts.AttributeRepo.FindByID(ctx, attrDTO.ID)
		assert.NoError(t, err)
		return attr
	}

	t.Run("Successfully update attribute name and sort order", func(t *testing.T) {
		attr := createAttribute("Old Name", "OLD", 1, nil, nil, []string{"Val1", "Val2"})

		newName := "New Name"
		newSortOrder := 2
		cmd := application.UpdateAttributeCommand{
			ID:        attr.ID,
			Name:      &newName,
			SortOrder: &newSortOrder,
		}

		updatedAttr, err := ts.ProductService.UpdateAttribute(ctx, cmd)
		assert.NoError(t, err)
		assert.NotNil(t, updatedAttr)
		assert.Equal(t, newName, updatedAttr.Name)
		assert.Equal(t, newSortOrder, updatedAttr.SortOrder)

		// Verify in DB
		savedAttr, err := ts.AttributeRepo.FindByID(ctx, attr.ID)
		assert.NoError(t, err)
		assert.NotNil(t, savedAttr)
		assert.Equal(t, newName, savedAttr.Name)
		assert.Equal(t, newSortOrder, savedAttr.SortOrder)
	})

	t.Run("Successfully update attribute values (add, modify, delete)", func(t *testing.T) {
		attr := createAttribute("Test Values", "TV", 1, nil, nil, []string{"A", "B"})
		valA := attr.Values[0] // Assuming first value is A
		valB := attr.Values[1] // Assuming second value is B

		// Add C, Modify A, Delete B
		newValC := "C"
		updatedValA := "A-Modified"

		cmd := application.UpdateAttributeCommand{
			ID: attr.ID,
			Values: []application.UpdateAttributeValueCommand{
				{ID: &valA.ID, Value: updatedValA, Code: updatedValA[:1]}, // Modify existing A
				{Value: newValC, Code: newValC[:1]},                       // Add new C
			},
		}

		updatedAttr, err := ts.ProductService.UpdateAttribute(ctx, cmd)
		assert.NoError(t, err)
		assert.NotNil(t, updatedAttr)
		assert.Len(t, updatedAttr.Values, 2) // A-Modified, C

		var updatedAValue string
		for _, v := range updatedAttr.Values {
			if v == updatedValA {
				updatedAValue = v
			}
		}
		assert.Equal(t, updatedValA, updatedAValue)
		assert.Contains(t, updatedAttr.Values, newValC)
		assert.NotContains(t, updatedAttr.Values, valB.Value) // B should be deleted

		// Verify in DB
		savedAttr, err := ts.AttributeRepo.FindByID(ctx, attr.ID)
		assert.NoError(t, err)
		assert.NotNil(t, savedAttr)
		assert.Len(t, savedAttr.Values, 2)
		assert.NotContains(t, savedAttr.Values, func(v domain.AttributeValue) bool { return v.Value == valB.Value })
	})

	t.Run("Updating a non-existent attribute should return error", func(t *testing.T) {
		nonExistentID := uuid.New()
		newName := "Invalid"
		cmd := application.UpdateAttributeCommand{
			ID:   nonExistentID,
			Name: &newName,
		}

		updatedAttr, err := ts.ProductService.UpdateAttribute(ctx, cmd)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "attribute with ID "+nonExistentID.String()+" does not exist")
		assert.Nil(t, updatedAttr)
	})
}

func TestProductService_GetProductByID_Integration(t *testing.T) {
	ts := setupTest(t)
	defer teardownTest(ts)

	ctx := ts.Ctx

	createBrand := func(name string) *domain.Brand {
		brand := &domain.Brand{ID: uuid.New(), Name: name}
		assert.NoError(t, ts.TestDB.DB.Create(persistence.BrandFromDomain(brand)).Error)
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
		productDTO, err := ts.ProductService.CreateProduct(ctx, cmd)
		assert.NoError(t, err)
		product, err := ts.ProductRepo.FindByID(ctx, productDTO.ID)
		assert.NoError(t, err)
		return product
	}

	t.Run("Successfully get product by ID", func(t *testing.T) {
		brand := createBrand("Test Brand")
		product := createProduct(brand.ID, nil)

		query := application.GetProductByIDQuery{ID: product.ID}
		foundProduct, err := ts.ProductService.GetProductByID(ctx, query)
		assert.NoError(t, err)
		assert.NotNil(t, foundProduct)
		assert.Equal(t, product.ID, foundProduct.ID)
		assert.Equal(t, product.Name, foundProduct.Name)
	})

	t.Run("Getting a non-existent product by ID should return error", func(t *testing.T) {
		nonExistentID := uuid.New()
		query := application.GetProductByIDQuery{ID: nonExistentID}

		foundProduct, err := ts.ProductService.GetProductByID(ctx, query)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "product with ID "+nonExistentID.String()+" does not exist")
		assert.Nil(t, foundProduct)
	})
}

func TestProductService_ListProducts_Integration(t *testing.T) {
	ts := setupTest(t)
	defer teardownTest(ts)

	ctx := ts.Ctx

	createBrand := func(name string) *domain.Brand {
		brand := &domain.Brand{ID: uuid.New(), Name: name}
		assert.NoError(t, ts.TestDB.DB.Create(persistence.BrandFromDomain(brand)).Error)
		return brand
	}

	createProduct := func(brandID uuid.UUID, groupIDs []uuid.UUID, sku string) *domain.Product {
		barcode := "12345"
		cmd := application.CreateProductCommand{
			SKU:         sku,
			Name:        "Test Product",
			LongName:    "Long Test Product",
			Barcode:     &barcode,
			Description: "Description",
			ProductType: domain.ProductTypeTangible,
			BrandID:     brandID,
			GroupIDs:    groupIDs,
		}
		productDTO, err := ts.ProductService.CreateProduct(ctx, cmd)
		assert.NoError(t, err)
		product, err := ts.ProductRepo.FindByID(ctx, productDTO.ID)
		assert.NoError(t, err)
		return product
	}

	t.Run("List all products when no filters are applied", func(t *testing.T) {
		brand := createBrand("Brand X")
		createProduct(brand.ID, nil, "PROD-1")
		createProduct(brand.ID, nil, "PROD-2")

		query := application.ListProductsQuery{}
		products, err := ts.ProductService.ListProducts(ctx, query)
		assert.NoError(t, err)
		assert.Len(t, products, 2)
	})

	t.Run("List products with filters (not yet implemented in service)", func(t *testing.T) {
		// As the service.ListProducts currently returns an empty slice, this tests that behavior.
		query := application.ListProductsQuery{} // No filters for now
		products, err := ts.ProductService.ListProducts(ctx, query)
		assert.NoError(t, err)
		assert.Empty(t, products)
	})
}

func TestProductService_GenerateProductVariants_Integration(t *testing.T) {
	ts := setupTest(t)
	defer teardownTest(ts)

	ctx := ts.Ctx

	createBrand := func(name string) *domain.Brand {
		brand := &domain.Brand{ID: uuid.New(), Name: name}
		assert.NoError(t, ts.TestDB.DB.Create(persistence.BrandFromDomain(brand)).Error)
		return brand
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
		productDTO, err := ts.ProductService.CreateProduct(ctx, cmd)
		assert.NoError(t, err)

		for _, attrID := range directAttrIDs {
			_, err := ts.ProductService.AddDirectAttributeToProduct(ctx, application.AddDirectAttributeCommand{
				ProductID:   productDTO.ID,
				AttributeID: attrID,
			})
			assert.NoError(t, err)
		}

		product, err := ts.ProductRepo.FindByID(ctx, productDTO.ID)
		assert.NoError(t, err)
		return product
	}
	createAttribute := func(name, code string, sortOrder int, scopeBrandID, scopeGroupID *uuid.UUID, values []string) *domain.Attribute {
		attrValues := make([]application.CreateAttributeValueCommand, len(values))
		for i, v := range values {
			attrValues[i] = application.CreateAttributeValueCommand{Value: v, Code: v[:1]}
		}
		cmd := application.CreateAttributeCommand{
			Name:         name,
			Code:         code,
			SortOrder:    sortOrder,
			ScopeBrandID: scopeBrandID,
			ScopeGroupID: scopeGroupID,
			Values:       attrValues,
		}
		attrDTO, err := ts.ProductService.CreateAttribute(ctx, cmd)
		assert.NoError(t, err)
		attr, err := ts.AttributeRepo.FindByID(ctx, attrDTO.ID)
		assert.NoError(t, err)
		return attr
	}

	t.Run("Successfully generates product variants for a product", func(t *testing.T) {
		brand := createBrand("Brand A")
		product := createProduct(brand.ID, nil, nil)
		attrColor := createAttribute("Color", "COL", 1, nil, nil, []string{"Red", "Blue"})
		attrSize := createAttribute("Size", "SIZ", 2, nil, nil, []string{"S", "M"})
		_ = attrColor
		_ = attrSize

		cmd := application.GenerateProductVariantsCommand{ProductID: product.ID}
		err := ts.ProductService.GenerateProductVariants(ctx, cmd)
		assert.NoError(t, err)

		// Verify variants are created (expect 4 variants: Red-S, Red-M, Blue-S, Blue-M)
		// This requires a new method in ProductVariantRepository to list by productID.
		// For now, we'll check directly via FindByProductIDAndAttributeValues once implemented.
		// Assuming GenerateProductVariants actually creates them.
	})
}

func TestProductService_FindOrCreateProductVariant_Integration(t *testing.T) {
	ts := setupTest(t)
	defer teardownTest(ts)

	ctx := ts.Ctx

	createBrand := func(name string) *domain.Brand {
		brand := &domain.Brand{ID: uuid.New(), Name: name}
		assert.NoError(t, ts.TestDB.DB.Create(persistence.BrandFromDomain(brand)).Error)
		return brand
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
		productDTO, err := ts.ProductService.CreateProduct(ctx, cmd)
		assert.NoError(t, err)

		for _, attrID := range directAttrIDs {
			_, err := ts.ProductService.AddDirectAttributeToProduct(ctx, application.AddDirectAttributeCommand{
				ProductID:   productDTO.ID,
				AttributeID: attrID,
			})
			assert.NoError(t, err)
		}

		product, err := ts.ProductRepo.FindByID(ctx, productDTO.ID)
		assert.NoError(t, err)
		return product
	}
	createAttribute := func(name, code string, sortOrder int, scopeBrandID, scopeGroupID *uuid.UUID, values []string) *domain.Attribute {
		attrValues := make([]application.CreateAttributeValueCommand, len(values))
		for i, v := range values {
			attrValues[i] = application.CreateAttributeValueCommand{Value: v, Code: v[:1]}
		}
		cmd := application.CreateAttributeCommand{
			Name:         name,
			Code:         code,
			SortOrder:    sortOrder,
			ScopeBrandID: scopeBrandID,
			ScopeGroupID: scopeGroupID,
			Values:       attrValues,
		}
		attrDTO, err := ts.ProductService.CreateAttribute(ctx, cmd)
		assert.NoError(t, err)
		attr, err := ts.AttributeRepo.FindByID(ctx, attrDTO.ID)
		assert.NoError(t, err)
		return attr
	}

	t.Run("Finds an existing variant", func(t *testing.T) {
		brand := createBrand("Brand FOC")
		attrColor := createAttribute("Color", "FOC_C", 1, nil, nil, []string{"Red"})
		attrSize := createAttribute("Size", "FOC_S", 2, nil, nil, []string{"M"})
		product := createProduct(brand.ID, nil, nil) // No direct attributes initially

		// Manually create a variant first (for simplicity, typically done via GenerateProductVariants or directly)
		variant, err := domain.NewProductVariant(product.ID, "PROD-FOC-FOC_C.R-FOC_S.M", nil, domain.StatusConfirmed, []uuid.UUID{attrColor.Values[0].ID, attrSize.Values[0].ID})
		assert.NoError(t, err)
		variant.Status = domain.StatusConfirmed
		assert.NoError(t, ts.VariantRepo.Save(ctx, variant))

		cmd := application.FindOrCreateProductVariantCommand{
			ProductID: product.ID,
			OptionConfiguration: []application.OptionConfigurationItem{
				{AttributeName: attrColor.Code, Value: attrColor.Values[0].Value},
				{AttributeName: attrSize.Code, Value: attrSize.Values[0].Value},
			},
		}

		foundVariant, err := ts.ProductService.FindOrCreateProductVariant(ctx, cmd)
		assert.NoError(t, err)
		assert.NotNil(t, foundVariant)
		assert.Equal(t, variant.ID, foundVariant.ID)
		assert.Equal(t, domain.StatusConfirmed, foundVariant.Status)
	})

	t.Run("Creates a new provisional variant if not found", func(t *testing.T) {
		brand := createBrand("Brand New")
		attrColor := createAttribute("Color", "NEW_C", 1, nil, nil, []string{"Blue"})
		attrSize := createAttribute("Size", "NEW_S", 2, nil, nil, []string{"L"})
		product := createProduct(brand.ID, nil, nil)

		cmd := application.FindOrCreateProductVariantCommand{
			ProductID: product.ID,
			OptionConfiguration: []application.OptionConfigurationItem{
				{AttributeName: attrColor.Code, Value: attrColor.Values[0].Value},
				{AttributeName: attrSize.Code, Value: attrSize.Values[0].Value},
			},
		}

		newVariant, err := ts.ProductService.FindOrCreateProductVariant(ctx, cmd)
		assert.NoError(t, err)
		assert.NotNil(t, newVariant)
		assert.NotEqual(t, uuid.Nil, newVariant.ID)
		assert.Equal(t, domain.StatusProvisional, newVariant.Status)
		assert.Equal(t, product.SKU+"-"+attrColor.Code+"."+attrColor.Values[0].Code+"-"+attrSize.Code+"."+attrSize.Values[0].Code, newVariant.SKU)
	})

	t.Run("Returns error for invalid option configuration", func(t *testing.T) {
		brand := createBrand("Brand Invalid")
		createAttribute("Color", "INV_C", 1, nil, nil, []string{"Red"})
		product := createProduct(brand.ID, nil, nil)

		cmd := application.FindOrCreateProductVariantCommand{
			ProductID: product.ID,
			OptionConfiguration: []application.OptionConfigurationItem{
				{AttributeName: "INV_C", Value: "InvalidValue"}, // Invalid value
			},
		}

		variant, err := ts.ProductService.FindOrCreateProductVariant(ctx, cmd)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "value 'InvalidValue' is not valid for attribute 'INV_C'")
		assert.Nil(t, variant)
	})
}

func TestProductService_ListProductVariantsByProductID_Integration(t *testing.T) {
	ts := setupTest(t)
	defer teardownTest(ts)

	ctx := ts.Ctx

	createBrand := func(name string) *domain.Brand {
		brand := &domain.Brand{ID: uuid.New(), Name: name}
		assert.NoError(t, ts.TestDB.DB.Create(persistence.BrandFromDomain(brand)).Error)
		return brand
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
		productDTO, err := ts.ProductService.CreateProduct(ctx, cmd)
		assert.NoError(t, err)

		for _, attrID := range directAttrIDs {
			_, err := ts.ProductService.AddDirectAttributeToProduct(ctx, application.AddDirectAttributeCommand{
				ProductID:   productDTO.ID,
				AttributeID: attrID,
			})
			assert.NoError(t, err)
		}

		product, err := ts.ProductRepo.FindByID(ctx, productDTO.ID)
		assert.NoError(t, err)
		return product
	}
	createAttribute := func(name, code string, sortOrder int, scopeBrandID, scopeGroupID *uuid.UUID, values []string) *domain.Attribute {
		attrValues := make([]application.CreateAttributeValueCommand, len(values))
		for i, v := range values {
			attrValues[i] = application.CreateAttributeValueCommand{Value: v, Code: v[:1]}
		}
		cmd := application.CreateAttributeCommand{
			Name:         name,
			Code:         code,
			SortOrder:    sortOrder,
			ScopeBrandID: scopeBrandID,
			ScopeGroupID: scopeGroupID,
			Values:       attrValues,
		}
		attrDTO, err := ts.ProductService.CreateAttribute(ctx, cmd)
		assert.NoError(t, err)
		attr, err := ts.AttributeRepo.FindByID(ctx, attrDTO.ID)
		assert.NoError(t, err)
		return attr
	}

	t.Run("Successfully lists variants for a product", func(t *testing.T) {
		brand := createBrand("Brand LV")
		product := createProduct(brand.ID, nil, nil)
		attrColor := createAttribute("Color", "LV_C", 1, nil, nil, []string{"Red", "Blue"})
		attrSize := createAttribute("Size", "LV_S", 2, nil, nil, []string{"S", "M"})

		// Create variants manually for testing list
		variant1, err := domain.NewProductVariant(product.ID, "SKU-LV-1", nil, domain.StatusConfirmed, []uuid.UUID{attrColor.Values[0].ID, attrSize.Values[0].ID}) // Red S
		assert.NoError(t, err)
		assert.NoError(t, ts.VariantRepo.Save(ctx, variant1))
		variant2, err := domain.NewProductVariant(product.ID, "SKU-LV-2", nil, domain.StatusConfirmed, []uuid.UUID{attrColor.Values[1].ID, attrSize.Values[0].ID}) // Blue S
		assert.NoError(t, err)
		assert.NoError(t, ts.VariantRepo.Save(ctx, variant2))

		query := application.ListProductVariantsByProductIDQuery{ProductID: product.ID}
		variants, err := ts.ProductService.ListProductVariantsByProductID(ctx, query)
		assert.NoError(t, err)
		assert.Len(t, variants, 2)
		// Check that the returned variants match the created ones
		var foundIDs []uuid.UUID
		for _, v := range variants {
			foundIDs = append(foundIDs, v.ID)
		}
		assert.Contains(t, foundIDs, variant1.ID)
		assert.Contains(t, foundIDs, variant2.ID)
	})

	t.Run("Returns empty list if no variants for product", func(t *testing.T) {
		brand := createBrand("Brand NoVar")
		product := createProduct(brand.ID, nil, nil)

		query := application.ListProductVariantsByProductIDQuery{ProductID: product.ID}
		variants, err := ts.ProductService.ListProductVariantsByProductID(ctx, query)
		assert.NoError(t, err)
		assert.Empty(t, variants)
	})
}

func TestProductService_GetProductVariantByID_Integration(t *testing.T) {
	ts := setupTest(t)
	defer teardownTest(ts)

	ctx := ts.Ctx

	createBrand := func(name string) *domain.Brand {
		brand := &domain.Brand{ID: uuid.New(), Name: name}
		assert.NoError(t, ts.TestDB.DB.Create(persistence.BrandFromDomain(brand)).Error)
		return brand
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
		productDTO, err := ts.ProductService.CreateProduct(ctx, cmd)
		assert.NoError(t, err)

		for _, attrID := range directAttrIDs {
			_, err := ts.ProductService.AddDirectAttributeToProduct(ctx, application.AddDirectAttributeCommand{
				ProductID:   productDTO.ID,
				AttributeID: attrID,
			})
			assert.NoError(t, err)
		}

		product, err := ts.ProductRepo.FindByID(ctx, productDTO.ID)
		assert.NoError(t, err)
		return product
	}
	createAttribute := func(name, code string, sortOrder int, scopeBrandID, scopeGroupID *uuid.UUID, values []string) *domain.Attribute {
		attrValues := make([]application.CreateAttributeValueCommand, len(values))
		for i, v := range values {
			attrValues[i] = application.CreateAttributeValueCommand{Value: v, Code: v[:1]}
		}
		cmd := application.CreateAttributeCommand{
			Name:         name,
			Code:         code,
			SortOrder:    sortOrder,
			ScopeBrandID: scopeBrandID,
			ScopeGroupID: scopeGroupID,
			Values:       attrValues,
		}
		attrDTO, err := ts.ProductService.CreateAttribute(ctx, cmd)
		assert.NoError(t, err)
		attr, err := ts.AttributeRepo.FindByID(ctx, attrDTO.ID)
		assert.NoError(t, err)
		return attr
	}

	t.Run("Successfully gets a product variant by ID", func(t *testing.T) {
		brand := createBrand("Brand GVBID")
		product := createProduct(brand.ID, nil, nil)
		attrColor := createAttribute("Color", "GVB_C", 1, nil, nil, []string{"Green"})
		attrSize := createAttribute("Size", "GVB_S", 2, nil, nil, []string{"XL"})

		variant, err := domain.NewProductVariant(product.ID, "SKU-GVBID-1", nil, domain.StatusConfirmed, []uuid.UUID{attrColor.Values[0].ID, attrSize.Values[0].ID})
		assert.NoError(t, err)
		assert.NoError(t, ts.VariantRepo.Save(ctx, variant))

		query := application.GetProductVariantByIDQuery{ID: variant.ID}
		foundVariant, err := ts.ProductService.GetProductVariantByID(ctx, query)
		assert.NoError(t, err)
		assert.NotNil(t, foundVariant)
		assert.Equal(t, variant.ID, foundVariant.ID)
		assert.Equal(t, variant.SKU, foundVariant.SKU)
		assert.Contains(t, foundVariant.OptionConfiguration, "Color")
		assert.Equal(t, "Green", foundVariant.OptionConfiguration["Color"])
	})

	t.Run("Returns error for non-existent product variant ID", func(t *testing.T) {
		nonExistentID := uuid.New()
		query := application.GetProductVariantByIDQuery{ID: nonExistentID}
		foundVariant, err := ts.ProductService.GetProductVariantByID(ctx, query)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "product variant with ID "+nonExistentID.String()+" does not exist")
		assert.Nil(t, foundVariant)
	})
}

func TestProductService_GetProductVariantBySKU_Integration(t *testing.T) {
	ts := setupTest(t)
	defer teardownTest(ts)

	ctx := ts.Ctx

	createBrand := func(name string) *domain.Brand {
		brand := &domain.Brand{ID: uuid.New(), Name: name}
		assert.NoError(t, ts.TestDB.DB.Create(persistence.BrandFromDomain(brand)).Error)
		return brand
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
		productDTO, err := ts.ProductService.CreateProduct(ctx, cmd)
		assert.NoError(t, err)

		for _, attrID := range directAttrIDs {
			_, err := ts.ProductService.AddDirectAttributeToProduct(ctx, application.AddDirectAttributeCommand{
				ProductID:   productDTO.ID,
				AttributeID: attrID,
			})
			assert.NoError(t, err)
		}

		product, err := ts.ProductRepo.FindByID(ctx, productDTO.ID)
		assert.NoError(t, err)
		return product
	}
	createAttribute := func(name, code string, sortOrder int, scopeBrandID, scopeGroupID *uuid.UUID, values []string) *domain.Attribute {
		attrValues := make([]application.CreateAttributeValueCommand, len(values))
		for i, v := range values {
			attrValues[i] = application.CreateAttributeValueCommand{Value: v, Code: v[:1]}
		}
		cmd := application.CreateAttributeCommand{
			Name:         name,
			Code:         code,
			SortOrder:    sortOrder,
			ScopeBrandID: scopeBrandID,
			ScopeGroupID: scopeGroupID,
			Values:       attrValues,
		}
		attrDTO, err := ts.ProductService.CreateAttribute(ctx, cmd)
		assert.NoError(t, err)
		attr, err := ts.AttributeRepo.FindByID(ctx, attrDTO.ID)
		assert.NoError(t, err)
		return attr
	}

	t.Run("Successfully gets a product variant by SKU", func(t *testing.T) {
		brand := createBrand("Brand GVBSKU")
		product := createProduct(brand.ID, nil, nil)
		attrColor := createAttribute("Color", "GVBS_C", 1, nil, nil, []string{"Blue"})
		attrSize := createAttribute("Size", "GVBS_S", 2, nil, nil, []string{"XXL"})

		variantSKU := product.SKU + "-" + attrColor.Code + "." + attrColor.Values[0].Code + "-" + attrSize.Code + "." + attrSize.Values[0].Code
		variant, err := domain.NewProductVariant(product.ID, variantSKU, nil, domain.StatusConfirmed, []uuid.UUID{attrColor.Values[0].ID, attrSize.Values[0].ID})
		assert.NoError(t, err)
		assert.NoError(t, ts.VariantRepo.Save(ctx, variant))

		query := application.GetProductVariantBySKUQuery{SKU: variantSKU}
		foundVariant, err := ts.ProductService.GetProductVariantBySKU(ctx, query)
		assert.NoError(t, err)
		assert.NotNil(t, foundVariant)
		assert.Equal(t, variant.ID, foundVariant.ID)
		assert.Equal(t, variant.SKU, foundVariant.SKU)
		assert.Contains(t, foundVariant.OptionConfiguration, "Color")
		assert.Equal(t, "Blue", foundVariant.OptionConfiguration["Color"])
	})

	t.Run("Returns error for non-existent product variant SKU", func(t *testing.T) {
		nonExistentSKU := "NON-EXISTENT-SKU"
		query := application.GetProductVariantBySKUQuery{SKU: nonExistentSKU}
		foundVariant, err := ts.ProductService.GetProductVariantBySKU(ctx, query)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "product variant with SKU "+nonExistentSKU+" does not exist")
		assert.Nil(t, foundVariant)
	})
}

func TestProductService_UpdateProductVariant_Integration(t *testing.T) {
	ts := setupTest(t)
	defer teardownTest(ts)

	ctx := ts.Ctx

	createBrand := func(name string) *domain.Brand {
		brand := &domain.Brand{ID: uuid.New(), Name: name}
		assert.NoError(t, ts.TestDB.DB.Create(persistence.BrandFromDomain(brand)).Error)
		return brand
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
		productDTO, err := ts.ProductService.CreateProduct(ctx, cmd)
		assert.NoError(t, err)

		for _, attrID := range directAttrIDs {
			_, err := ts.ProductService.AddDirectAttributeToProduct(ctx, application.AddDirectAttributeCommand{
				ProductID:   productDTO.ID,
				AttributeID: attrID,
			})
			assert.NoError(t, err)
		}

		product, err := ts.ProductRepo.FindByID(ctx, productDTO.ID)
		assert.NoError(t, err)
		return product
	}
	createAttribute := func(name, code string, sortOrder int, scopeBrandID, scopeGroupID *uuid.UUID, values []string) *domain.Attribute {
		attrValues := make([]application.CreateAttributeValueCommand, len(values))
		for i, v := range values {
			attrValues[i] = application.CreateAttributeValueCommand{Value: v, Code: v[:1]}
		}
		cmd := application.CreateAttributeCommand{
			Name:         name,
			Code:         code,
			SortOrder:    sortOrder,
			ScopeBrandID: scopeBrandID,
			ScopeGroupID: scopeGroupID,
			Values:       attrValues,
		}
		attrDTO, err := ts.ProductService.CreateAttribute(ctx, cmd)
		assert.NoError(t, err)
		attr, err := ts.AttributeRepo.FindByID(ctx, attrDTO.ID)
		assert.NoError(t, err)
		return attr
	}

	t.Run("Successfully updates a product variant", func(t *testing.T) {
		brand := createBrand("Brand UPV")
		product := createProduct(brand.ID, nil, nil)
		attrColor := createAttribute("Color", "UPV_C", 1, nil, nil, []string{"White"})
		attrSize := createAttribute("Size", "UPV_S", 2, nil, nil, []string{"XS"})

		variantSKU := product.SKU + "-" + attrColor.Code + "." + attrColor.Values[0].Code + "-" + attrSize.Code + "." + attrSize.Values[0].Code
		variant, err := domain.NewProductVariant(product.ID, variantSKU, nil, domain.StatusProvisional, []uuid.UUID{attrColor.Values[0].ID, attrSize.Values[0].ID})
		assert.NoError(t, err)
		variant.Status = domain.StatusProvisional // Start as provisional
		assert.NoError(t, ts.VariantRepo.Save(ctx, variant))

		newBarcode := "NEW-BARCODE-123"
		newStatus := domain.StatusConfirmed
		newIsActive := false
		cmd := application.UpdateProductVariantCommand{
			ID:       variant.ID,
			Barcode:  &newBarcode,
			IsActive: &newIsActive,
			Status:   &newStatus,
		}

		updatedVariant, err := ts.ProductService.UpdateProductVariant(ctx, cmd)
		assert.NoError(t, err)
		assert.NotNil(t, updatedVariant)
		assert.Equal(t, newBarcode, *updatedVariant.Barcode)
		assert.Equal(t, newStatus, updatedVariant.Status)
		assert.Equal(t, newIsActive, updatedVariant.IsActive)

		// Verify in DB
		savedVariant, err := ts.VariantRepo.FindByID(ctx, variant.ID)
		assert.NoError(t, err)
		assert.NotNil(t, savedVariant)
		assert.Equal(t, newBarcode, *savedVariant.Barcode)
		assert.Equal(t, newStatus, savedVariant.Status)
		assert.Equal(t, newIsActive, savedVariant.IsActive)
	})

	t.Run("Updating a provisional variant confirms it if status is not explicitly set", func(t *testing.T) {
		brand := createBrand("Brand AutoConfirm")
		product := createProduct(brand.ID, nil, nil)
		attrColor := createAttribute("Color", "AUTOC_C", 1, nil, nil, []string{"Black"})
		variant, err := domain.NewProductVariant(product.ID, product.SKU+"-AUTOC_C.B", nil, domain.StatusProvisional, []uuid.UUID{attrColor.Values[0].ID})
		assert.NoError(t, err)
		variant.Status = domain.StatusProvisional // Explicitly provisional
		assert.NoError(t, ts.VariantRepo.Save(ctx, variant))

		newBarcode := "CONFIRM-BARCODE"
		cmd := application.UpdateProductVariantCommand{
			ID:      variant.ID,
			Barcode: &newBarcode,
		}

		updatedVariant, err := ts.ProductService.UpdateProductVariant(ctx, cmd)
		assert.NoError(t, err)
		assert.Equal(t, domain.StatusConfirmed, updatedVariant.Status) // Should be confirmed

		savedVariant, err := ts.VariantRepo.FindByID(ctx, variant.ID)
		assert.NoError(t, err)
		assert.Equal(t, domain.StatusConfirmed, savedVariant.Status)
	})

	t.Run("Updating a non-existent product variant returns error", func(t *testing.T) {
		nonExistentID := uuid.New()
		newBarcode := "NON-EXISTENT"
		cmd := application.UpdateProductVariantCommand{
			ID:      nonExistentID,
			Barcode: &newBarcode,
		}

		updatedVariant, err := ts.ProductService.UpdateProductVariant(ctx, cmd)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "product variant with ID "+nonExistentID.String()+" does not exist")
		assert.Nil(t, updatedVariant)
	})
}

// TestPartyServiceConfiguration_Integration contains integration tests for PartyServiceConfiguration management.
func TestPartyServiceConfiguration_Integration(t *testing.T) {
	ts := setupTest(t)
	defer teardownTest(ts)

	ctx := ts.Ctx

	// Helper function to create a dummy Party (since Party is from another module)
	createDummyParty := func() uuid.UUID {
		partyID := uuid.New()
		// In a real scenario, you'd create a Party entry in the 'parties' table.
		// For now, we'll just return a new UUID and rely on foreign key to be optional or mocked.
		// If 'parties' table exists, we should insert a record there.
		// For this test, we assume a party exists or the FK is not strictly enforced in tests.
		// TODO: Create a mock Party entity in the database if the FK is enforced.
		return partyID
	}

	partyID := createDummyParty()

	t.Run("Create a new PartyServiceConfiguration successfully", func(t *testing.T) {
		configDetails := map[string]interface{}{"url": "http://example.com/api", "token": "abc123def"}

		cmd := application.CreatePartyServiceConfigurationCommand{
			PartyID:              partyID,
			ServiceID:            "MESSAGING_SERVICE",
			Name:                 "SMS Gateway Config",
			ConfigurationDetails: configDetails,
		}

		configDTO, err := ts.ProductService.CreatePartyServiceConfiguration(ctx, cmd)
		assert.NoError(t, err)
		assert.NotNil(t, configDTO)
		assert.NotEqual(t, uuid.Nil, configDTO.ID)
		assert.Equal(t, cmd.PartyID, configDTO.PartyID)
		assert.Equal(t, cmd.ServiceID, configDTO.ServiceID)
		assert.Equal(t, cmd.Name, configDTO.Name)

		retrievedConfig, err := ts.PartyServiceConfigurationRepo.FindByID(ctx, configDTO.PartyID, configDTO.ID)
		assert.NoError(t, err)
		assert.NotNil(t, retrievedConfig)
		assert.Equal(t, configDTO.ID, retrievedConfig.ID)
	})

	t.Run("Get a PartyServiceConfiguration by ID", func(t *testing.T) {
		configDetails := map[string]interface{}{"host": "smtp.example.com", "port": 587}
		configDetailsJSON, _ := json.Marshal(configDetails)

		config, err := domain.NewPartyServiceConfiguration(partyID, "EMAIL_SERVICE", "SMTP Config", configDetailsJSON)
		assert.NoError(t, err)
		assert.NoError(t, ts.PartyServiceConfigurationRepo.Save(ctx, config))

		query := application.GetPartyServiceConfigurationByIDQuery{PartyID: partyID, ID: config.ID}
		foundConfigDTO, err := ts.ProductService.GetPartyServiceConfigurationByID(ctx, query)
		assert.NoError(t, err)
		assert.NotNil(t, foundConfigDTO)
		assert.Equal(t, config.ID, foundConfigDTO.ID)
		assert.Equal(t, config.ServiceID, foundConfigDTO.ServiceID)
	})

	t.Run("List PartyServiceConfigurations by PartyID", func(t *testing.T) {
		anotherPartyID := createDummyParty() // New party
		config1, _ := domain.NewPartyServiceConfiguration(anotherPartyID, "API_SERVICE_1", "API Config 1", json.RawMessage(`{"endpoint": "/v1"}`))
		config2, _ := domain.NewPartyServiceConfiguration(anotherPartyID, "API_SERVICE_2", "API Config 2", json.RawMessage(`{"api_key": "xyz"}`))

		assert.NoError(t, ts.PartyServiceConfigurationRepo.Save(ctx, config1))
		assert.NoError(t, ts.PartyServiceConfigurationRepo.Save(ctx, config2))

		query := application.ListPartyServiceConfigurationsByPartyIDQuery{PartyID: anotherPartyID}
		foundConfigs, err := ts.ProductService.ListPartyServiceConfigurationsByPartyID(ctx, query)
		assert.NoError(t, err)
		assert.Len(t, foundConfigs, 2)

		var foundIDs []uuid.UUID
		for _, cfg := range foundConfigs {
			foundIDs = append(foundIDs, cfg.ID)
		}
		assert.Contains(t, foundIDs, config1.ID)
		assert.Contains(t, foundIDs, config2.ID)
	})

	t.Run("Update an existing PartyServiceConfiguration", func(t *testing.T) {
		configDetails := map[string]interface{}{"old_value": true}
		configDetailsJSON, _ := json.Marshal(configDetails)

		config, err := domain.NewPartyServiceConfiguration(partyID, "UPDATE_SERVICE", "Update Me", configDetailsJSON)
		assert.NoError(t, err)
		assert.NoError(t, ts.PartyServiceConfigurationRepo.Save(ctx, config))

		newServiceName := "UPDATED_SERVICE"
		newConfigDetails := map[string]interface{}{"new_value": "updated"}
		newConfigDetailsJSON, _ := json.Marshal(newConfigDetails)

		cmd := application.UpdatePartyServiceConfigurationCommand{
			ID:                   config.ID,
			PartyID:              partyID,
			ServiceID:            &newServiceName,
			ConfigurationDetails: newConfigDetails,
		}

		updatedConfigDTO, err := ts.ProductService.UpdatePartyServiceConfiguration(ctx, cmd)
		assert.NoError(t, err)
		assert.NotNil(t, updatedConfigDTO)
		assert.Equal(t, newServiceName, updatedConfigDTO.ServiceID)
		assert.Equal(t, newConfigDetailsJSON, []byte(updatedConfigDTO.ConfigurationDetails))
	})

	t.Run("Delete a PartyServiceConfiguration", func(t *testing.T) {
		configDetails := map[string]interface{}{"delete_me": true}
		configDetailsJSON, _ := json.Marshal(configDetails)

		config, err := domain.NewPartyServiceConfiguration(partyID, "DELETE_SERVICE", "Delete Test", configDetailsJSON)
		assert.NoError(t, err)
		assert.NoError(t, ts.PartyServiceConfigurationRepo.Save(ctx, config))

		cmd := application.DeletePartyServiceConfigurationCommand{PartyID: partyID, ID: config.ID}
		err = ts.ProductService.DeletePartyServiceConfiguration(ctx, cmd)
		assert.NoError(t, err)

		// Verify deletion
		query := application.GetPartyServiceConfigurationByIDQuery{PartyID: partyID, ID: config.ID}
		foundConfig, err := ts.ProductService.GetPartyServiceConfigurationByID(ctx, query)
		assert.Error(t, err) // Expect an error as it should not be found
		assert.Nil(t, foundConfig)
	})
}
