package persistence

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/joran-cortez/tramatex/internal/product/domain"
	"github.com/stretchr/testify/assert"
)

func TestGORMProductRepository_Save_Update(t *testing.T) {
	tdb, cleanup := setupProductTestDB(t)
	defer cleanup()

	ctxCreate := context.WithValue(context.Background(), "actorID", "actor-1")
	ctxUpdate := context.WithValue(context.Background(), "actorID", "actor-2")
	db := tdb.DB

	brandID := uuid.New()
	assert.NoError(t, db.WithContext(ctxCreate).Create(&BrandDataModel{ID: brandID, Name: "Brand"}).Error)

	product := &domain.Product{
		ID:          uuid.New(),
		SKU:         "P-100",
		Name:        "Product",
		ProductType: domain.ProductTypeTangible,
		BrandID:     brandID,
		IsActive:    true,
	}

	repo := NewGORMProductRepository(db)
	assert.NoError(t, repo.Save(ctxCreate, product))

	product.Name = "Product Updated"
	assert.NoError(t, repo.Save(ctxUpdate, product))

	var dataModel ProductDataModel
	assert.NoError(t, db.WithContext(ctxUpdate).First(&dataModel, "id = ?", product.ID).Error)
	// TODO: CreatedBy and ModifiedBy fields not yet implemented in ProductDataModel
	// assert.Equal(t, "actor-1", dataModel.CreatedBy)
	// assert.Equal(t, "actor-2", dataModel.ModifiedBy)
	assert.Equal(t, "Product Updated", dataModel.Name)
}

func TestGORMVariantRepository_Save_Update(t *testing.T) {
	tdb, cleanup := setupProductTestDB(t)
	defer cleanup()

	ctxCreate := context.WithValue(context.Background(), "actorID", "actor-1")
	ctxUpdate := context.WithValue(context.Background(), "actorID", "actor-2")
	db := tdb.DB

	brandID := uuid.New()
	productID := uuid.New()
	assert.NoError(t, db.WithContext(ctxCreate).Create(&BrandDataModel{ID: brandID, Name: "Brand"}).Error)
	assert.NoError(t, db.WithContext(ctxCreate).Create(&ProductDataModel{ID: productID, SKU: "P-200", Name: "Product", ProductType: "TANGIBLE", BrandID: brandID, IsActive: true}).Error)

	variant := &domain.ProductVariant{
		ID:              uuid.New(),
		ProductID:       productID,
		SKU:             "P-200-C.R",
		AttributeValues: []uuid.UUID{uuid.New()},
		Status:          domain.StatusConfirmed,
		IsActive:        true,
	}

	repo := NewGORMVariantRepository(db)
	assert.NoError(t, repo.Save(ctxCreate, variant))

	variant.SKU = "P-200-C.R-UPDATED"
	assert.NoError(t, repo.Save(ctxUpdate, variant))

	var dataModel VariantDataModel
	assert.NoError(t, db.WithContext(ctxUpdate).First(&dataModel, "id = ?", variant.ID).Error)
	assert.Equal(t, "actor-1", dataModel.CreatedBy)
	assert.Equal(t, "actor-2", dataModel.ModifiedBy)
	assert.Equal(t, "P-200-C.R-UPDATED", dataModel.SKU)
}

func TestGORMAttributeRepository_Save_UpdateValues(t *testing.T) {
	tdb, cleanup := setupProductTestDB(t)
	defer cleanup()

	ctxCreate := context.WithValue(context.Background(), "actorID", "actor-1")
	ctxUpdate := context.WithValue(context.Background(), "actorID", "actor-2")
	repo := NewGORMAttributeRepository(tdb.DB)

	attr, _ := domain.NewAttribute("Color", "C", 1)
	valRed, _ := attr.AddValue("Red", "R")
	assert.NoError(t, repo.Save(ctxCreate, attr))

	_ = attr.UpdateValue(valRed.ID, "Red-Updated", "R2")
	valBlue, _ := attr.AddValue("Blue", "B")
	_ = attr.RemoveValue(valRed.ID)
	_ = valBlue

	assert.NoError(t, repo.Save(ctxUpdate, attr))

	found, err := repo.FindByID(ctxUpdate, attr.ID)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Len(t, found.Values, 1)
	assert.Equal(t, "Blue", found.Values[0].Value)
}

func TestGORMBrandAndGroupRepositories_NotFound(t *testing.T) {
	tdb, cleanup := setupProductTestDB(t)
	defer cleanup()

	ctx := context.Background()
	brandRepo := NewGORMBrandRepository(tdb.DB)
	groupRepo := NewGORMProductGroupRepository(tdb.DB)

	brand, err := brandRepo.FindByID(ctx, uuid.New())
	assert.NoError(t, err)
	assert.Nil(t, brand)

	group, err := groupRepo.FindByID(ctx, uuid.New())
	assert.NoError(t, err)
	assert.Nil(t, group)
}
