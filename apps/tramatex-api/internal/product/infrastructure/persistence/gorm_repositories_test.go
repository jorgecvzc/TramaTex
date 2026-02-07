package persistence

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/joran-cortez/tramatex/internal/product/domain"
	"github.com/stretchr/testify/assert"
)

func setupProductTestDB(t *testing.T) (*TestDB, func()) {
	tdb := NewTestDB(t)
	if tdb.DB == nil {
		t.Skip("PostgreSQL not available for integration tests")
	}

	db := tdb.DB
	assert.NoError(t, tdb.SetUpProduct())

	cleanup := func() {
		_ = tdb.TearDownProduct()
		sqlDB, _ := db.DB()
		_ = sqlDB.Close()
	}

	return tdb, cleanup
}

func TestDataModelConversions(t *testing.T) {
	brand := &domain.Brand{ID: uuid.New(), Name: "Brand"}
	brandModel := BrandFromDomain(brand)
	assert.Equal(t, brand.ID, brandModel.ID)
	assert.Equal(t, brand.Name, brandModel.Name)
	assert.Equal(t, brand.ID, brandModel.ToDomain().ID)

	group := &domain.ProductGroup{ID: uuid.New(), Name: "Group"}
	groupModel := ProductGroupFromDomain(group)
	assert.Equal(t, group.ID, groupModel.ID)
	assert.Equal(t, group.Name, groupModel.Name)
	assert.Equal(t, group.ID, groupModel.ToDomain().ID)

	attr, _ := domain.NewAttribute("Color", "C", 1, nil, nil)
	val, _ := attr.AddValue("Red", "R")
	attrModel := AttributeFromDomain(attr)
	assert.Equal(t, attr.ID, attrModel.ID)
	assert.Len(t, attrModel.Values, 1)
	assert.Equal(t, val.ID, attrModel.Values[0].ID)
	assert.Equal(t, attr.ID, attrModel.ToDomain().ID)

	product := &domain.Product{ID: uuid.New(), SKU: "P-1", Name: "Prod", BrandID: uuid.New(), IsActive: true}
	productModel := FromDomain(product)
	assert.Equal(t, product.ID, productModel.ID)
	assert.Equal(t, product.ID, productModel.ToDomain().ID)

	variant := &domain.ProductVariant{ID: uuid.New(), ProductID: product.ID, SKU: "P-1-C.R", AttributeValues: []uuid.UUID{val.ID}, Status: domain.StatusConfirmed}
	variantModel := VariantFromDomain(variant)
	assert.Equal(t, variant.ID, variantModel.ID)
	assert.Equal(t, variant.ID, variantModel.ToDomain().ID)

	config, _ := domain.NewPartyServiceConfiguration(uuid.New(), "svc", "Name", json.RawMessage(`{"k":"v"}`))
	configModel, err := PartyServiceConfigurationFromDomain(config)
	assert.NoError(t, err)
	assert.Equal(t, config.ID, configModel.ID)
	_, err = configModel.ToDomain()
	assert.NoError(t, err)
}

func TestGORMRepositories(t *testing.T) {
	tdb, cleanup := setupProductTestDB(t)
	defer cleanup()

	ctx := context.Background()
	db := tdb.DB

	brandID := uuid.New()
	groupID := uuid.New()
	productID := uuid.New()
	variantID := uuid.New()

	assert.NoError(t, db.WithContext(ctx).Create(&BrandDataModel{ID: brandID, Name: "Brand"}).Error)
	assert.NoError(t, db.WithContext(ctx).Create(&ProductGroupDataModel{ID: groupID, Name: "Group"}).Error)

	brandRepo := NewGORMBrandRepository(db)
	brand, err := brandRepo.FindByID(ctx, brandID)
	assert.NoError(t, err)
	assert.Equal(t, brandID, brand.ID)

	groupRepo := NewGORMProductGroupRepository(db)
	group, err := groupRepo.FindByID(ctx, groupID)
	assert.NoError(t, err)
	assert.Equal(t, groupID, group.ID)

	attr, _ := domain.NewAttribute("Color", "C", 1, nil, nil)
	_, _ = attr.AddValue("Red", "R")
	attrRepo := NewGORMAttributeRepository(db)
	assert.NoError(t, attrRepo.Save(ctx, attr))
	attrFound, err := attrRepo.FindByID(ctx, attr.ID)
	assert.NoError(t, err)
	assert.Equal(t, attr.ID, attrFound.ID)
	assert.Len(t, attrFound.Values, 1)

	product := &domain.Product{
		ID:          productID,
		SKU:         "P-1",
		Name:        "Product",
		ProductType: domain.ProductTypeTangible,
		BrandID:     brandID,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	productRepo := NewGORMProductRepository(db)
	assert.NoError(t, productRepo.Save(ctx, product))
	foundByID, err := productRepo.FindByID(ctx, productID)
	assert.NoError(t, err)
	assert.Equal(t, productID, foundByID.ID)
	foundBySKU, err := productRepo.FindBySKU(ctx, "P-1")
	assert.NoError(t, err)
	assert.Equal(t, productID, foundBySKU.ID)
	all, err := productRepo.FindAll(ctx)
	assert.NoError(t, err)
	assert.NotEmpty(t, all)
	assert.NoError(t, productRepo.UpdateSKUs(ctx, productID, "P-1B"))
	updated, err := productRepo.FindByID(ctx, productID)
	assert.NoError(t, err)
	assert.Equal(t, "P-1B", updated.SKU)

	variant := &domain.ProductVariant{
		ID:              variantID,
		ProductID:       productID,
		SKU:             "P-1B-C.R",
		AttributeValues: []uuid.UUID{attr.Values[0].ID},
		Status:          domain.StatusConfirmed,
		IsActive:        true,
	}
	variantRepo := NewGORMVariantRepository(db)
	assert.NoError(t, variantRepo.Save(ctx, variant))
	foundVariant, err := variantRepo.FindByID(ctx, variantID)
	assert.NoError(t, err)
	assert.Equal(t, variantID, foundVariant.ID)
	foundVariantBySKU, err := variantRepo.FindBySKU(ctx, variant.SKU)
	assert.NoError(t, err)
	assert.Equal(t, variantID, foundVariantBySKU.ID)
	variantsByProduct, err := variantRepo.FindByProductID(ctx, productID)
	assert.NoError(t, err)
	assert.Len(t, variantsByProduct, 1)
	_, err = variantRepo.FindByProductIDAndAttributeValues(ctx, productID, variant.AttributeValues)
	assert.NoError(t, err)

	config, _ := domain.NewPartyServiceConfiguration(uuid.New(), "svc", "Config", json.RawMessage(`{"k":"v"}`))
	pscRepo := NewGORMPartyServiceConfigurationRepository(db)
	assert.NoError(t, db.Exec("INSERT INTO parties (id) VALUES (?)", config.PartyID).Error)
	assert.NoError(t, pscRepo.Save(ctx, config))
	foundConfig, err := pscRepo.FindByID(ctx, config.PartyID, config.ID)
	assert.NoError(t, err)
	assert.Equal(t, config.ID, foundConfig.ID)
	configs, err := pscRepo.FindByPartyID(ctx, config.PartyID)
	assert.NoError(t, err)
	assert.Len(t, configs, 1)
	assert.NoError(t, pscRepo.Delete(ctx, config.PartyID, config.ID))
}

func TestGORMAttributeRepository_FindByScope(t *testing.T) {
	tdb, cleanup := setupProductTestDB(t)
	defer cleanup()

	ctx := context.Background()
	db := tdb.DB

	brandID := uuid.New()
	groupID := uuid.New()

	assert.NoError(t, db.WithContext(ctx).Create(&BrandDataModel{ID: brandID, Name: "Brand"}).Error)
	assert.NoError(t, db.WithContext(ctx).Create(&ProductGroupDataModel{ID: groupID, Name: "Group"}).Error)

	attrRepo := NewGORMAttributeRepository(db)

	genericAttr, _ := domain.NewAttribute("Generic", "GEN", 0, nil, nil)
	brandAttr, _ := domain.NewAttribute("Brand", "BR", 0, &brandID, nil)
	groupAttr, _ := domain.NewAttribute("Group", "GR", 0, nil, &groupID)
	brandGroupAttr, _ := domain.NewAttribute("BrandGroup", "BG", 0, &brandID, &groupID)

	assert.NoError(t, attrRepo.Save(ctx, genericAttr))
	assert.NoError(t, attrRepo.Save(ctx, brandAttr))
	assert.NoError(t, attrRepo.Save(ctx, groupAttr))
	assert.NoError(t, attrRepo.Save(ctx, brandGroupAttr))

	genericResults, err := attrRepo.FindByScope(ctx, nil, nil)
	assert.NoError(t, err)
	assert.Len(t, genericResults, 1)

	brandResults, err := attrRepo.FindByScope(ctx, &brandID, nil)
	assert.NoError(t, err)
	assert.Len(t, brandResults, 1)

	groupResults, err := attrRepo.FindByScope(ctx, nil, &groupID)
	assert.NoError(t, err)
	assert.Len(t, groupResults, 1)

	brandGroupResults, err := attrRepo.FindByScope(ctx, &brandID, &groupID)
	assert.NoError(t, err)
	assert.Len(t, brandGroupResults, 1)
}

func TestGORMAttributeRepository_FindByIDs(t *testing.T) {
	tdb, cleanup := setupProductTestDB(t)
	defer cleanup()

	ctx := context.Background()
	db := tdb.DB

	attrRepo := NewGORMAttributeRepository(db)
	attrOne, _ := domain.NewAttribute("Color", "C", 0, nil, nil)
	attrTwo, _ := domain.NewAttribute("Size", "S", 1, nil, nil)
	assert.NoError(t, attrRepo.Save(ctx, attrOne))
	assert.NoError(t, attrRepo.Save(ctx, attrTwo))

	results, err := attrRepo.FindByIDs(ctx, []uuid.UUID{attrOne.ID, attrTwo.ID})
	assert.NoError(t, err)
	assert.Len(t, results, 2)
}

func TestGORMProductRepository_NotFound(t *testing.T) {
	tdb, cleanup := setupProductTestDB(t)
	defer cleanup()

	ctx := context.Background()
	db := tdb.DB

	productRepo := NewGORMProductRepository(db)

	product, err := productRepo.FindByID(ctx, uuid.New())
	assert.NoError(t, err)
	assert.Nil(t, product)

	product, err = productRepo.FindBySKU(ctx, "missing")
	assert.NoError(t, err)
	assert.Nil(t, product)

	products, err := productRepo.FindAll(ctx)
	assert.NoError(t, err)
	assert.Len(t, products, 0)
}

func TestGORMVariantRepository_NotFound(t *testing.T) {
	tdb, cleanup := setupProductTestDB(t)
	defer cleanup()

	ctx := context.Background()
	db := tdb.DB

	variantRepo := NewGORMVariantRepository(db)

	variant, err := variantRepo.FindByID(ctx, uuid.New())
	assert.NoError(t, err)
	assert.Nil(t, variant)

	variant, err = variantRepo.FindBySKU(ctx, "missing")
	assert.NoError(t, err)
	assert.Nil(t, variant)

	variants, err := variantRepo.FindByProductID(ctx, uuid.New())
	assert.NoError(t, err)
	assert.Len(t, variants, 0)
}

func TestGORMPartyServiceConfigurationRepository_NotFound(t *testing.T) {
	tdb, cleanup := setupProductTestDB(t)
	defer cleanup()

	ctx := context.Background()
	db := tdb.DB

	repo := NewGORMPartyServiceConfigurationRepository(db)
	partyID := uuid.New()
	configID := uuid.New()

	config, err := repo.FindByID(ctx, partyID, configID)
	assert.NoError(t, err)
	assert.Nil(t, config)

	configs, err := repo.FindByPartyID(ctx, partyID)
	assert.NoError(t, err)
	assert.Len(t, configs, 0)

	assert.NoError(t, repo.Delete(ctx, partyID, configID))
}

func TestGORMPartyServiceConfigurationRepository_Save_InvalidJSON(t *testing.T) {
	tdb, cleanup := setupProductTestDB(t)
	defer cleanup()

	ctx := context.Background()
	db := tdb.DB

	repo := NewGORMPartyServiceConfigurationRepository(db)
	config := &domain.PartyServiceConfiguration{
		ID:                   uuid.New(),
		PartyID:              uuid.New(),
		ServiceID:            "svc",
		Name:                 "Bad",
		ConfigurationDetails: json.RawMessage("{invalid"),
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}
	assert.NoError(t, db.Exec("INSERT INTO parties (id) VALUES (?)", config.PartyID).Error)

	err := repo.Save(ctx, config)
	assert.Error(t, err)
}
