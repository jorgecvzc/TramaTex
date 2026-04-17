package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joran-cortez/tramatex/internal/product/application"
	"github.com/joran-cortez/tramatex/internal/product/domain"
	infra_middleware "github.com/joran-cortez/tramatex/internal/shared/infrastructure/middleware"
	"github.com/stretchr/testify/assert"
)

func uuidPtr(id uuid.UUID) *uuid.UUID { return &id }

type stubProductRepo struct {
	saveFn            func(context.Context, *domain.Product) error
	findByIDFn        func(context.Context, uuid.UUID) (*domain.Product, error)
	findByIDsFn       func(context.Context, []uuid.UUID) ([]*domain.Product, error)
	findBySKFn        func(context.Context, string) (*domain.Product, error)
	findByBarcodeFn   func(context.Context, string) (*domain.Product, error)
	findBySKUPrefixFn func(context.Context, string) ([]*domain.Product, error)
	findAllFn         func(context.Context) ([]*domain.Product, error)
	updateSKUsFn      func(context.Context, uuid.UUID, string) error
}

func (s *stubProductRepo) Save(ctx context.Context, product *domain.Product) error {
	if s.saveFn != nil {
		return s.saveFn(ctx, product)
	}
	return nil
}

func (s *stubProductRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	if s.findByIDFn != nil {
		return s.findByIDFn(ctx, id)
	}
	return nil, nil
}

func (s *stubProductRepo) FindByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Product, error) {
	if s.findByIDsFn != nil {
		return s.findByIDsFn(ctx, ids)
	}
	return nil, nil
}

func (s *stubProductRepo) FindBySKU(ctx context.Context, sku string) (*domain.Product, error) {
	if s.findBySKFn != nil {
		return s.findBySKFn(ctx, sku)
	}
	return nil, nil
}

func (s *stubProductRepo) FindByBarcode(ctx context.Context, barcode string) (*domain.Product, error) {
	if s.findByBarcodeFn != nil {
		return s.findByBarcodeFn(ctx, barcode)
	}
	return nil, nil
}

func (s *stubProductRepo) FindBySKUPrefix(ctx context.Context, prefix string) ([]*domain.Product, error) {
	if s.findBySKUPrefixFn != nil {
		return s.findBySKUPrefixFn(ctx, prefix)
	}
	return nil, nil
}

func (s *stubProductRepo) FindAll(ctx context.Context) ([]*domain.Product, error) {
	if s.findAllFn != nil {
		return s.findAllFn(ctx)
	}
	return nil, nil
}

func (s *stubProductRepo) UpdateSKUs(ctx context.Context, productID uuid.UUID, newSKU string) error {
	if s.updateSKUsFn != nil {
		return s.updateSKUsFn(ctx, productID, newSKU)
	}
	return nil
}

type stubBrandRepo struct {
	saveFn      func(context.Context, *domain.Brand) error
	findByIDFn  func(context.Context, uuid.UUID) (*domain.Brand, error)
	findByIDsFn func(context.Context, []uuid.UUID) ([]*domain.Brand, error)
	findAllFn   func(context.Context) ([]*domain.Brand, error)
	deleteFn    func(context.Context, uuid.UUID) error
}

func (s *stubBrandRepo) Save(ctx context.Context, brand *domain.Brand) error {
	if s.saveFn != nil {
		return s.saveFn(ctx, brand)
	}
	return nil
}

func (s *stubBrandRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Brand, error) {
	if s.findByIDFn != nil {
		return s.findByIDFn(ctx, id)
	}
	return nil, nil
}

func (s *stubBrandRepo) FindByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Brand, error) {
	if s.findByIDsFn != nil {
		return s.findByIDsFn(ctx, ids)
	}
	return nil, nil
}

func (s *stubBrandRepo) FindAll(ctx context.Context) ([]*domain.Brand, error) {
	if s.findAllFn != nil {
		return s.findAllFn(ctx)
	}
	return nil, nil
}

func (s *stubBrandRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if s.deleteFn != nil {
		return s.deleteFn(ctx, id)
	}
	return nil
}

type stubGroupRepo struct {
	saveFn     func(context.Context, *domain.ProductGroup) error
	findByIDFn func(context.Context, uuid.UUID) (*domain.ProductGroup, error)
	findAllFn  func(context.Context) ([]*domain.ProductGroup, error)
	deleteFn   func(context.Context, uuid.UUID) error
}

func (s *stubGroupRepo) Save(ctx context.Context, group *domain.ProductGroup) error {
	if s.saveFn != nil {
		return s.saveFn(ctx, group)
	}
	return nil
}

func (s *stubGroupRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.ProductGroup, error) {
	if s.findByIDFn != nil {
		return s.findByIDFn(ctx, id)
	}
	return nil, nil
}

func (s *stubGroupRepo) FindAll(ctx context.Context) ([]*domain.ProductGroup, error) {
	if s.findAllFn != nil {
		return s.findAllFn(ctx)
	}
	return nil, nil
}

func (s *stubGroupRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if s.deleteFn != nil {
		return s.deleteFn(ctx, id)
	}
	return nil
}

type stubAttributeRepo struct {
	saveFn       func(context.Context, *domain.Attribute) error
	findByIDFn   func(context.Context, uuid.UUID) (*domain.Attribute, error)
	findByCodeFn func(context.Context, string) (*domain.Attribute, error)
	findByIDsFn  func(context.Context, []uuid.UUID) ([]domain.Attribute, error)
	findByScope  func(context.Context, *uuid.UUID, *uuid.UUID) ([]*domain.Attribute, error)
	deleteFn     func(context.Context, uuid.UUID) error
}

func (s *stubAttributeRepo) Save(ctx context.Context, attribute *domain.Attribute) error {
	if s.saveFn != nil {
		return s.saveFn(ctx, attribute)
	}
	return nil
}

func (s *stubAttributeRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Attribute, error) {
	if s.findByIDFn != nil {
		return s.findByIDFn(ctx, id)
	}
	return nil, nil
}

func (s *stubAttributeRepo) FindByCode(ctx context.Context, code string) (*domain.Attribute, error) {
	if s.findByCodeFn != nil {
		return s.findByCodeFn(ctx, code)
	}
	return nil, nil
}

func (s *stubAttributeRepo) FindByIDs(ctx context.Context, ids []uuid.UUID) ([]domain.Attribute, error) {
	if s.findByIDsFn != nil {
		return s.findByIDsFn(ctx, ids)
	}
	return nil, nil
}

func (s *stubAttributeRepo) FindByScope(ctx context.Context, brandID *uuid.UUID, groupID *uuid.UUID) ([]*domain.Attribute, error) {
	if s.findByScope != nil {
		return s.findByScope(ctx, brandID, groupID)
	}
	return nil, nil
}

func (s *stubAttributeRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if s.deleteFn != nil {
		return s.deleteFn(ctx, id)
	}
	return nil
}

type stubVariantRepo struct {
	saveFn                     func(context.Context, *domain.ProductVariant) error
	findByIDFn                 func(context.Context, uuid.UUID) (*domain.ProductVariant, error)
	findByIDsFn                func(context.Context, []uuid.UUID) ([]*domain.ProductVariant, error)
	findBySKFn                 func(context.Context, string) (*domain.ProductVariant, error)
	findByBarcodeFn            func(context.Context, string) (*domain.ProductVariant, error)
	findBySKUPrefixFn          func(context.Context, string) ([]*domain.ProductVariant, error)
	findByProductIDFn          func(context.Context, uuid.UUID) ([]*domain.ProductVariant, error)
	findByProductIDAndValuesFn func(context.Context, uuid.UUID, []uuid.UUID) (*domain.ProductVariant, error)
}

func (s *stubVariantRepo) Save(ctx context.Context, variant *domain.ProductVariant) error {
	if s.saveFn != nil {
		return s.saveFn(ctx, variant)
	}
	return nil
}

func (s *stubVariantRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.ProductVariant, error) {
	if s.findByIDFn != nil {
		return s.findByIDFn(ctx, id)
	}
	return nil, nil
}

func (s *stubVariantRepo) FindByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.ProductVariant, error) {
	if s.findByIDsFn != nil {
		return s.findByIDsFn(ctx, ids)
	}
	return nil, nil
}

func (s *stubVariantRepo) FindBySKU(ctx context.Context, sku string) (*domain.ProductVariant, error) {
	if s.findBySKFn != nil {
		return s.findBySKFn(ctx, sku)
	}
	return nil, nil
}

func (s *stubVariantRepo) FindByBarcode(ctx context.Context, barcode string) (*domain.ProductVariant, error) {
	if s.findByBarcodeFn != nil {
		return s.findByBarcodeFn(ctx, barcode)
	}
	return nil, nil
}

func (s *stubVariantRepo) FindBySKUPrefix(ctx context.Context, prefix string) ([]*domain.ProductVariant, error) {
	if s.findBySKUPrefixFn != nil {
		return s.findBySKUPrefixFn(ctx, prefix)
	}
	return nil, nil
}

func (s *stubVariantRepo) FindByProductID(ctx context.Context, productID uuid.UUID) ([]*domain.ProductVariant, error) {
	if s.findByProductIDFn != nil {
		return s.findByProductIDFn(ctx, productID)
	}
	return nil, nil
}

func (s *stubVariantRepo) FindByProductIDAndAttributeValues(ctx context.Context, productID uuid.UUID, attributeValueIDs []uuid.UUID) (*domain.ProductVariant, error) {
	if s.findByProductIDAndValuesFn != nil {
		return s.findByProductIDAndValuesFn(ctx, productID, attributeValueIDs)
	}
	return nil, nil
}

type stubPartyServiceConfigRepo struct {
	saveFn        func(context.Context, *domain.PartyServiceConfiguration) error
	findByIDFn    func(context.Context, uuid.UUID, uuid.UUID) (*domain.PartyServiceConfiguration, error)
	findByPartyFn func(context.Context, uuid.UUID) ([]*domain.PartyServiceConfiguration, error)
	deleteFn      func(context.Context, uuid.UUID, uuid.UUID) error
}

func (s *stubPartyServiceConfigRepo) Save(ctx context.Context, config *domain.PartyServiceConfiguration) error {
	if s.saveFn != nil {
		return s.saveFn(ctx, config)
	}
	return nil
}

func (s *stubPartyServiceConfigRepo) FindByID(ctx context.Context, partyID, id uuid.UUID) (*domain.PartyServiceConfiguration, error) {
	if s.findByIDFn != nil {
		return s.findByIDFn(ctx, partyID, id)
	}
	return nil, nil
}

func (s *stubPartyServiceConfigRepo) FindByPartyID(ctx context.Context, partyID uuid.UUID) ([]*domain.PartyServiceConfiguration, error) {
	if s.findByPartyFn != nil {
		return s.findByPartyFn(ctx, partyID)
	}
	return nil, nil
}

func (s *stubPartyServiceConfigRepo) Delete(ctx context.Context, partyID, id uuid.UUID) error {
	if s.deleteFn != nil {
		return s.deleteFn(ctx, partyID, id)
	}
	return nil
}

func registerProductRoutes(router *gin.Engine, handler *ProductHandler) {
	router.POST("/products", handler.CreateProduct)
	router.POST("/products/:id/groups", handler.AddGroupToProduct)
	router.POST("/products/:id/direct-attributes", handler.AddDirectAttributeToProduct)
	router.GET("/products", handler.ListProducts)
	router.GET("/products/smart-search", handler.SmartSearch)
	router.GET("/products/:id", handler.GetProductByID)
	router.PUT("/products/:id/sku", handler.UpdateProductSKU)
	router.GET("/products/:id/options", handler.GetCalculatedOptionSetsForProduct)
	router.POST("/products/:id/variants", handler.GenerateProductVariants)
	router.POST("/products/:id/variants/find-or-create", handler.FindOrCreateProductVariant)
	router.GET("/products/:id/variants", handler.ListProductVariantsByProductID)
	router.GET("/attributes", handler.ListAttributes)
	router.POST("/attributes", handler.CreateAttribute)
	router.GET("/attributes/:id", handler.GetAttributeByID)
	router.PUT("/attributes/:id", handler.UpdateAttribute)
	router.GET("/variants/:id", handler.GetProductVariantByID)
	router.GET("/variants", handler.GetProductVariantBySKU)
	router.PUT("/variants/:id", handler.UpdateProductVariant)
	router.POST("/parties/:id/configurations", handler.CreatePartyServiceConfiguration)
	router.GET("/parties/:id/configurations", handler.ListPartyServiceConfigurationsByPartyID)
	router.GET("/parties/:id/configurations/:configId", handler.GetPartyServiceConfigurationByID)
	router.PUT("/parties/:id/configurations/:configId", handler.UpdatePartyServiceConfiguration)
	router.DELETE("/parties/:id/configurations/:configId", handler.DeletePartyServiceConfiguration)
}

func newTestRouter(handler *ProductHandler) *gin.Engine {
	router := gin.New()
	router.Use(infra_middleware.ErrorHandlerMiddleware("development"))
	router.Use(func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "actorID", "test-actor")
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	})
	registerProductRoutes(router, handler)
	return router
}

func newTestRouterWithoutActor(handler *ProductHandler) *gin.Engine {
	router := gin.New()
	registerProductRoutes(router, handler)
	return router
}

func TestProductHandler_CreateProduct_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewProductHandler(nil)
	router := newTestRouter(handler)

	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBufferString("{"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProductHandler_CreateProduct_MissingActorID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewProductHandler(nil)
	router := newTestRouterWithoutActor(handler)

	body := map[string]interface{}{
		"sku":         "P-1",
		"name":        "Product",
		"productType": string(domain.ProductTypeTangible),
		"brandId":     uuid.New().String(),
	}
	data, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestProductHandler_GetProductByID_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewProductHandler(nil)
	router := newTestRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/products/invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProductHandler_UpdateProductSKU_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewProductHandler(nil)
	router := newTestRouter(handler)

	req := httptest.NewRequest(http.MethodPut, "/products/invalid/sku", bytes.NewBufferString("{\"NewSKU\":\"P-2\"}"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProductHandler_UpdateProductSKU_MissingActorID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewProductHandler(nil)
	router := newTestRouterWithoutActor(handler)

	req := httptest.NewRequest(http.MethodPut, "/products/"+uuid.New().String()+"/sku", bytes.NewBufferString("{\"NewSKU\":\"P-2\"}"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestProductHandler_ListProducts_InvalidBrandID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewProductHandler(nil)
	router := newTestRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/products?brandId=invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProductHandler_ListProducts_InvalidIsActive(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewProductHandler(nil)
	router := newTestRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/products?isActive=maybe", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProductHandler_ListProducts_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	product := &domain.Product{
		ID:          uuid.New(),
		SKU:         "P-1",
		Name:        "Product",
		BrandID:     uuidPtr(uuid.New()),
		IsActive:    true,
		ProductType: domain.ProductTypeTangible,
	}
	service := application.NewProductService(
		&stubProductRepo{findAllFn: func(ctx context.Context) ([]*domain.Product, error) {
			return []*domain.Product{product}, nil
		}},
		&stubBrandRepo{},
		&stubGroupRepo{},
		&stubAttributeRepo{},
		&stubVariantRepo{},
		&stubPartyServiceConfigRepo{},
	)
	handler := NewProductHandler(service)
	router := newTestRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}
	var response []application.ProductDTO
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	if assert.Len(t, response, 1) {
		assert.Equal(t, "P-1", response[0].SKU)
	}
}

func TestProductHandler_CreateProduct_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	brandID := uuid.New()
	groupID := uuid.New()
	service := application.NewProductService(
		&stubProductRepo{},
		&stubBrandRepo{findByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Brand, error) {
			return &domain.Brand{ID: id, Name: "Brand"}, nil
		}},
		&stubGroupRepo{findByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.ProductGroup, error) {
			return &domain.ProductGroup{ID: id, Name: "Group"}, nil
		}},
		&stubAttributeRepo{},
		&stubVariantRepo{},
		&stubPartyServiceConfigRepo{},
	)
	handler := NewProductHandler(service)
	router := newTestRouter(handler)

	body := map[string]interface{}{
		"sku":          "P-1",
		"name":         "Product",
		"product_type": string(domain.ProductTypeTangible),
		"brand_id":     brandID.String(),
		"group_ids":    []string{groupID.String()},
	}
	data, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestProductHandler_AddGroupToProduct_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	productID := uuid.New()
	groupID := uuid.New()
	product := &domain.Product{ID: productID, SKU: "P-1", Name: "Product", BrandID: uuidPtr(uuid.New()), IsActive: true}
	service := application.NewProductService(
		&stubProductRepo{findByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
			return product, nil
		}},
		&stubBrandRepo{},
		&stubGroupRepo{findByIDFn: func(context.Context, uuid.UUID) (*domain.ProductGroup, error) {
			return &domain.ProductGroup{ID: groupID, Name: "Group"}, nil
		}},
		&stubAttributeRepo{},
		&stubVariantRepo{},
		&stubPartyServiceConfigRepo{},
	)
	handler := NewProductHandler(service)
	router := newTestRouter(handler)

	body := map[string]string{
		"GroupID": groupID.String(),
	}
	data, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/products/"+productID.String()+"/groups", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestProductHandler_AddGroupToProduct_MissingActorID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewProductHandler(nil)
	router := newTestRouterWithoutActor(handler)

	body := map[string]string{"GroupID": uuid.New().String()}
	data, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/products/"+uuid.New().String()+"/groups", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestProductHandler_AddDirectAttributeToProduct_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	productID := uuid.New()
	attributeID := uuid.New()
	product := &domain.Product{ID: productID, SKU: "P-1", Name: "Product", BrandID: uuidPtr(uuid.New()), IsActive: true}
	attribute, _ := domain.NewAttribute("Color", "C")
	attribute.ID = attributeID
	service := application.NewProductService(
		&stubProductRepo{findByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
			return product, nil
		}},
		&stubBrandRepo{},
		&stubGroupRepo{},
		&stubAttributeRepo{findByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Attribute, error) {
			return attribute, nil
		}},
		&stubVariantRepo{},
		&stubPartyServiceConfigRepo{},
	)
	handler := NewProductHandler(service)
	router := newTestRouter(handler)

	body := map[string]string{
		"AttributeID": attributeID.String(),
	}
	data, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/products/"+productID.String()+"/direct-attributes", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestProductHandler_AddDirectAttributeToProduct_MissingActorID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewProductHandler(nil)
	router := newTestRouterWithoutActor(handler)

	body := map[string]string{"AttributeID": uuid.New().String()}
	data, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/products/"+uuid.New().String()+"/direct-attributes", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestProductHandler_UpdateProductSKU_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	productID := uuid.New()
	product := &domain.Product{ID: productID, SKU: "P-1", Name: "Product", BrandID: uuidPtr(uuid.New()), IsActive: true}
	updatedProduct := &domain.Product{ID: productID, SKU: "P-2", Name: "Product", BrandID: product.BrandID, IsActive: true}
	callCount := 0
	service := application.NewProductService(
		&stubProductRepo{
			findByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
				callCount++
				if callCount == 1 {
					return product, nil
				}
				return updatedProduct, nil
			},
		},
		&stubBrandRepo{},
		&stubGroupRepo{},
		&stubAttributeRepo{},
		&stubVariantRepo{},
		&stubPartyServiceConfigRepo{},
	)
	handler := NewProductHandler(service)
	router := newTestRouter(handler)

	req := httptest.NewRequest(http.MethodPut, "/products/"+productID.String()+"/sku", bytes.NewBufferString("{\"NewSKU\":\"P-2\"}"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductHandler_CreateAttribute_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := application.NewProductService(
		&stubProductRepo{},
		&stubBrandRepo{},
		&stubGroupRepo{},
		&stubAttributeRepo{},
		&stubVariantRepo{},
		&stubPartyServiceConfigRepo{},
	)
	handler := NewProductHandler(service)
	router := newTestRouter(handler)

	body := map[string]interface{}{
		"Name": "Color",
		"Code": "C",
		"Values": []map[string]string{
			{"Value": "Red", "Code": "R"},
		},
	}
	data, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/attributes", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestProductHandler_CreateAttribute_MissingActorID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewProductHandler(nil)
	router := newTestRouterWithoutActor(handler)

	data, _ := json.Marshal(map[string]interface{}{"Name": "Color", "Code": "C"})
	req := httptest.NewRequest(http.MethodPost, "/attributes", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestProductHandler_UpdateAttribute_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	attribute, _ := domain.NewAttribute("Color", "C")
	service := application.NewProductService(
		&stubProductRepo{},
		&stubBrandRepo{},
		&stubGroupRepo{},
		&stubAttributeRepo{findByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Attribute, error) {
			return attribute, nil
		}},
		&stubVariantRepo{},
		&stubPartyServiceConfigRepo{},
	)
	handler := NewProductHandler(service)
	router := newTestRouter(handler)

	req := httptest.NewRequest(http.MethodPut, "/attributes/"+attribute.ID.String(), bytes.NewBufferString("{\"Name\":\"Shade\"}"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductHandler_UpdateAttribute_MissingActorID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewProductHandler(nil)
	router := newTestRouterWithoutActor(handler)

	req := httptest.NewRequest(http.MethodPut, "/attributes/"+uuid.New().String(), bytes.NewBufferString("{}"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestProductHandler_GetAttributeByID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	attribute, _ := domain.NewAttribute("Color", "C")
	service := application.NewProductService(
		&stubProductRepo{},
		&stubBrandRepo{},
		&stubGroupRepo{},
		&stubAttributeRepo{findByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Attribute, error) {
			return attribute, nil
		}},
		&stubVariantRepo{},
		&stubPartyServiceConfigRepo{},
	)
	handler := NewProductHandler(service)
	router := newTestRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/attributes/"+attribute.ID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductHandler_ListAttributes_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	attribute, _ := domain.NewAttribute("Color", "C")
	service := application.NewProductService(
		&stubProductRepo{},
		&stubBrandRepo{},
		&stubGroupRepo{},
		&stubAttributeRepo{findByScope: func(ctx context.Context, brandID *uuid.UUID, groupID *uuid.UUID) ([]*domain.Attribute, error) {
			return []*domain.Attribute{attribute}, nil
		}},
		&stubVariantRepo{},
		&stubPartyServiceConfigRepo{},
	)
	handler := NewProductHandler(service)
	router := newTestRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/attributes", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductHandler_GetProductByID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	product := &domain.Product{ID: uuid.New(), SKU: "P-1", Name: "Product", BrandID: uuidPtr(uuid.New()), IsActive: true}
	service := application.NewProductService(
		&stubProductRepo{findByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
			return product, nil
		}},
		&stubBrandRepo{},
		&stubGroupRepo{},
		&stubAttributeRepo{},
		&stubVariantRepo{},
		&stubPartyServiceConfigRepo{},
	)
	handler := NewProductHandler(service)
	router := newTestRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/products/"+product.ID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductHandler_VariantEndpoints_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	productID := uuid.New()
	attribute, _ := domain.NewAttribute("Color", "C")
	value, _ := attribute.AddValue("Red", "R")
	variant := &domain.ProductVariant{ID: uuid.New(), ProductID: productID, SKU: "P-1-C.R", AttributeValues: []uuid.UUID{value.ID}, Status: domain.StatusConfirmed, IsActive: true}
	service := application.NewProductService(
		&stubProductRepo{findByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
			return &domain.Product{ID: productID, SKU: "P-1", Name: "Product", BrandID: uuidPtr(uuid.New()), IsActive: true, BasePrice: 10.0}, nil
		}},
		&stubBrandRepo{},
		&stubGroupRepo{},
		&stubAttributeRepo{findByScope: func(ctx context.Context, brandID *uuid.UUID, groupID *uuid.UUID) ([]*domain.Attribute, error) {
			return []*domain.Attribute{attribute}, nil
		}},
		&stubVariantRepo{
			findByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.ProductVariant, error) {
				return variant, nil
			},
			findBySKFn: func(ctx context.Context, sku string) (*domain.ProductVariant, error) {
				return variant, nil
			},
			findByProductIDFn: func(ctx context.Context, id uuid.UUID) ([]*domain.ProductVariant, error) {
				return []*domain.ProductVariant{variant}, nil
			},
			saveFn: func(ctx context.Context, variant *domain.ProductVariant) error {
				return nil
			},
		},
		&stubPartyServiceConfigRepo{},
	)
	handler := NewProductHandler(service)
	router := newTestRouter(handler)

	listReq := httptest.NewRequest(http.MethodGet, "/products/"+productID.String()+"/variants", nil)
	listRec := httptest.NewRecorder()
	router.ServeHTTP(listRec, listReq)
	assert.Equal(t, http.StatusOK, listRec.Code)

	getReq := httptest.NewRequest(http.MethodGet, "/variants/"+variant.ID.String(), nil)
	getRec := httptest.NewRecorder()
	router.ServeHTTP(getRec, getReq)
	assert.Equal(t, http.StatusOK, getRec.Code)

	getSKUReq := httptest.NewRequest(http.MethodGet, "/variants?sku="+variant.SKU, nil)
	getSKURec := httptest.NewRecorder()
	router.ServeHTTP(getSKURec, getSKUReq)
	assert.Equal(t, http.StatusOK, getSKURec.Code)

	updateReq := httptest.NewRequest(http.MethodPut, "/variants/"+variant.ID.String(), bytes.NewBufferString("{\"Barcode\":\"123\"}"))
	updateReq.Header.Set("Content-Type", "application/json")
	updateRec := httptest.NewRecorder()
	router.ServeHTTP(updateRec, updateReq)
	assert.Equal(t, http.StatusOK, updateRec.Code)
}

func TestProductHandler_PartyServiceConfigurationEndpoints_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	partyID := uuid.New()
	config, _ := domain.NewPartyServiceConfiguration(partyID, "svc", "Config", json.RawMessage(`{"k":"v"}`))
	service := application.NewProductService(
		&stubProductRepo{},
		&stubBrandRepo{},
		&stubGroupRepo{},
		&stubAttributeRepo{},
		&stubVariantRepo{},
		&stubPartyServiceConfigRepo{
			findByIDFn: func(ctx context.Context, pid, id uuid.UUID) (*domain.PartyServiceConfiguration, error) {
				return config, nil
			},
			findByPartyFn: func(ctx context.Context, id uuid.UUID) ([]*domain.PartyServiceConfiguration, error) {
				return []*domain.PartyServiceConfiguration{config}, nil
			},
		},
	)
	handler := NewProductHandler(service)
	router := newTestRouter(handler)

	createReq := httptest.NewRequest(http.MethodPost, "/parties/"+partyID.String()+"/configurations", bytes.NewBufferString("{\"ServiceID\":\"svc\",\"Name\":\"Config\",\"ConfigurationDetails\":{\"k\":\"v\"}}"))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()
	router.ServeHTTP(createRec, createReq)
	assert.Equal(t, http.StatusCreated, createRec.Code)

	listReq := httptest.NewRequest(http.MethodGet, "/parties/"+partyID.String()+"/configurations", nil)
	listRec := httptest.NewRecorder()
	router.ServeHTTP(listRec, listReq)
	assert.Equal(t, http.StatusOK, listRec.Code)

	getReq := httptest.NewRequest(http.MethodGet, "/parties/"+partyID.String()+"/configurations/"+config.ID.String(), nil)
	getRec := httptest.NewRecorder()
	router.ServeHTTP(getRec, getReq)
	assert.Equal(t, http.StatusOK, getRec.Code)

	updateReq := httptest.NewRequest(http.MethodPut, "/parties/"+partyID.String()+"/configurations/"+config.ID.String(), bytes.NewBufferString("{\"ServiceID\":\"svc2\"}"))
	updateReq.Header.Set("Content-Type", "application/json")
	updateRec := httptest.NewRecorder()
	router.ServeHTTP(updateRec, updateReq)
	assert.Equal(t, http.StatusOK, updateRec.Code)

	deleteReq := httptest.NewRequest(http.MethodDelete, "/parties/"+partyID.String()+"/configurations/"+config.ID.String(), nil)
	deleteRec := httptest.NewRecorder()
	router.ServeHTTP(deleteRec, deleteReq)
	assert.Equal(t, http.StatusNoContent, deleteRec.Code)
}

func TestProductHandler_GetAttributeByID_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewProductHandler(nil)
	router := newTestRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/attributes/invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProductHandler_GetCalculatedOptionSets_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewProductHandler(nil)
	router := newTestRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/products/invalid/options", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProductHandler_GenerateProductVariants_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewProductHandler(nil)
	router := newTestRouter(handler)

	req := httptest.NewRequest(http.MethodPost, "/products/invalid/variants", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProductHandler_ListProducts_InvalidGroupID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewProductHandler(nil)
	router := newTestRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/products?groupId=invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProductHandler_ListProductVariantsByProductID_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewProductHandler(nil)
	router := newTestRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/products/invalid/variants", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProductHandler_GetProductVariantByID_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewProductHandler(nil)
	router := newTestRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/variants/invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProductHandler_GetProductVariantBySKU_MissingSKU(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewProductHandler(nil)
	router := newTestRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/variants", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProductHandler_UpdateProductVariant_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewProductHandler(nil)
	router := newTestRouter(handler)

	req := httptest.NewRequest(http.MethodPut, "/variants/invalid", bytes.NewBufferString("{}"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProductHandler_UpdateProductVariant_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewProductHandler(nil)
	router := newTestRouter(handler)

	req := httptest.NewRequest(http.MethodPut, "/variants/"+uuid.New().String(), bytes.NewBufferString("{"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProductHandler_UpdateProductVariant_MissingActorID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewProductHandler(nil)
	router := newTestRouterWithoutActor(handler)

	req := httptest.NewRequest(http.MethodPut, "/variants/"+uuid.New().String(), bytes.NewBufferString("{}"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestProductHandler_FindOrCreateProductVariant_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewProductHandler(nil)
	router := newTestRouter(handler)

	req := httptest.NewRequest(http.MethodPost, "/products/"+uuid.New().String()+"/variants/find-or-create", bytes.NewBufferString("{"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProductHandler_GetCalculatedOptionSetsForProduct_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	productID := uuid.New()
	brandID := uuid.New()
	groupID := uuid.New()
	attribute, _ := domain.NewAttribute("Color", "C")
	product := &domain.Product{
		ID:                 productID,
		SKU:                "P-1",
		Name:               "Product",
		BrandID:            uuidPtr(brandID),
		GroupIDs:           []uuid.UUID{groupID},
		DirectAttributeIDs: []uuid.UUID{attribute.ID},
		IsActive:           true,
		ProductType:        domain.ProductTypeTangible,
	}

	service := application.NewProductService(
		&stubProductRepo{findByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
			return product, nil
		}},
		&stubBrandRepo{},
		&stubGroupRepo{},
		&stubAttributeRepo{
			findByScope: func(ctx context.Context, brandID *uuid.UUID, groupID *uuid.UUID) ([]*domain.Attribute, error) {
				return []*domain.Attribute{attribute}, nil
			},
			findByIDsFn: func(ctx context.Context, ids []uuid.UUID) ([]domain.Attribute, error) {
				return []domain.Attribute{*attribute}, nil
			},
			findByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Attribute, error) {
				if id == attribute.ID {
					return attribute, nil
				}
				return nil, nil
			},
		},
		&stubVariantRepo{},
		&stubPartyServiceConfigRepo{},
	)

	handler := NewProductHandler(service)
	router := newTestRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/products/"+productID.String()+"/options", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductHandler_GenerateProductVariants_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	productID := uuid.New()
	brandID := uuid.New()
	groupID := uuid.New()
	attribute, _ := domain.NewAttribute("Color", "C")
	product := &domain.Product{
		ID:                 productID,
		SKU:                "P-1",
		Name:               "Product",
		BrandID:            uuidPtr(brandID),
		GroupIDs:           []uuid.UUID{groupID},
		DirectAttributeIDs: []uuid.UUID{attribute.ID},
		IsActive:           true,
		ProductType:        domain.ProductTypeTangible,
	}

	service := application.NewProductService(
		&stubProductRepo{findByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
			return product, nil
		}},
		&stubBrandRepo{},
		&stubGroupRepo{},
		&stubAttributeRepo{
			findByScope: func(ctx context.Context, brandID *uuid.UUID, groupID *uuid.UUID) ([]*domain.Attribute, error) {
				return []*domain.Attribute{attribute}, nil
			},
			findByIDsFn: func(ctx context.Context, ids []uuid.UUID) ([]domain.Attribute, error) {
				return []domain.Attribute{*attribute}, nil
			},
			findByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Attribute, error) {
				if id == attribute.ID {
					return attribute, nil
				}
				return nil, nil
			},
		},
		&stubVariantRepo{},
		&stubPartyServiceConfigRepo{},
	)

	handler := NewProductHandler(service)
	router := newTestRouter(handler)

	req := httptest.NewRequest(http.MethodPost, "/products/"+productID.String()+"/variants", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)
}

func TestProductHandler_FindOrCreateProductVariant_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	productID := uuid.New()
	brandID := uuid.New()
	groupID := uuid.New()
	attribute, _ := domain.NewAttribute("Color", "C")
	value, _ := attribute.AddValue("Red", "R")
	product := &domain.Product{
		ID:                 productID,
		SKU:                "P-1",
		Name:               "Product",
		BrandID:            uuidPtr(brandID),
		GroupIDs:           []uuid.UUID{groupID},
		DirectAttributeIDs: []uuid.UUID{attribute.ID},
		IsActive:           true,
		ProductType:        domain.ProductTypeTangible,
	}

	service := application.NewProductService(
		&stubProductRepo{findByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
			return product, nil
		}},
		&stubBrandRepo{},
		&stubGroupRepo{},
		&stubAttributeRepo{
			findByScope: func(ctx context.Context, brandID *uuid.UUID, groupID *uuid.UUID) ([]*domain.Attribute, error) {
				return []*domain.Attribute{attribute}, nil
			},
			findByIDsFn: func(ctx context.Context, ids []uuid.UUID) ([]domain.Attribute, error) {
				return []domain.Attribute{*attribute}, nil
			},
			findByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Attribute, error) {
				if id == attribute.ID {
					return attribute, nil
				}
				return nil, nil
			},
		},
		&stubVariantRepo{
			findByProductIDAndValuesFn: func(ctx context.Context, id uuid.UUID, valueIDs []uuid.UUID) (*domain.ProductVariant, error) {
				return nil, nil
			},
			saveFn: func(ctx context.Context, variant *domain.ProductVariant) error {
				return nil
			},
		},
		&stubPartyServiceConfigRepo{},
	)

	handler := NewProductHandler(service)
	router := newTestRouter(handler)

	body := map[string]interface{}{
		"optionConfiguration": map[string]string{
			attribute.Code: value.Value,
		},
	}
	data, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/products/"+productID.String()+"/variants/find-or-create", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestProductHandler_FindOrCreateProductVariant_MissingActorID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewProductHandler(nil)
	router := newTestRouterWithoutActor(handler)

	data, _ := json.Marshal(map[string]interface{}{"optionConfiguration": map[string]string{"C": "R"}})
	req := httptest.NewRequest(http.MethodPost, "/products/"+uuid.New().String()+"/variants/find-or-create", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestProductHandler_CreatePartyServiceConfiguration_InvalidPartyID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewProductHandler(nil)
	router := newTestRouter(handler)

	req := httptest.NewRequest(http.MethodPost, "/parties/invalid/configurations", bytes.NewBufferString("{}"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProductHandler_CreatePartyServiceConfiguration_MissingActorID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewProductHandler(nil)
	router := newTestRouterWithoutActor(handler)

	body := map[string]interface{}{"serviceId": "svc", "name": "Config"}
	data, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/parties/"+uuid.New().String()+"/configurations", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestProductHandler_GetPartyServiceConfigurationByID_InvalidPartyID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewProductHandler(nil)
	router := newTestRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/parties/invalid/configurations/"+uuid.New().String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProductHandler_GetPartyServiceConfigurationByID_InvalidConfigID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewProductHandler(nil)
	router := newTestRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/parties/"+uuid.New().String()+"/configurations/invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProductHandler_UpdatePartyServiceConfiguration_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewProductHandler(nil)
	router := newTestRouter(handler)

	req := httptest.NewRequest(http.MethodPut, "/parties/"+uuid.New().String()+"/configurations/"+uuid.New().String(), bytes.NewBufferString("{"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProductHandler_DeletePartyServiceConfiguration_InvalidConfigID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewProductHandler(nil)
	router := newTestRouter(handler)

	req := httptest.NewRequest(http.MethodDelete, "/parties/"+uuid.New().String()+"/configurations/invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProductHandler_CreateProduct_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := application.NewProductService(
		&stubProductRepo{},
		&stubBrandRepo{findByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Brand, error) {
			return nil, assert.AnError
		}},
		&stubGroupRepo{},
		&stubAttributeRepo{},
		&stubVariantRepo{},
		&stubPartyServiceConfigRepo{},
	)
	handler := NewProductHandler(service)
	router := newTestRouter(handler)

	body := map[string]interface{}{
		"sku":          "P-1",
		"name":         "Product",
		"product_type": string(domain.ProductTypeTangible),
		"brand_id":     uuid.New().String(),
	}
	data, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestProductHandler_AddGroupToProduct_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	productID := uuid.New()
	service := application.NewProductService(
		&stubProductRepo{findByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
			return nil, nil
		}},
		&stubBrandRepo{},
		&stubGroupRepo{},
		&stubAttributeRepo{},
		&stubVariantRepo{},
		&stubPartyServiceConfigRepo{},
	)
	handler := NewProductHandler(service)
	router := newTestRouter(handler)

	body := map[string]string{
		"GroupID": uuid.New().String(),
	}
	data, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/products/"+productID.String()+"/groups", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestProductHandler_AddDirectAttributeToProduct_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	productID := uuid.New()
	service := application.NewProductService(
		&stubProductRepo{findByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
			return &domain.Product{ID: id, SKU: "P-1", Name: "Product", BrandID: uuidPtr(uuid.New())}, nil
		}},
		&stubBrandRepo{},
		&stubGroupRepo{},
		&stubAttributeRepo{findByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Attribute, error) {
			return nil, nil
		}},
		&stubVariantRepo{},
		&stubPartyServiceConfigRepo{},
	)
	handler := NewProductHandler(service)
	router := newTestRouter(handler)

	body := map[string]string{
		"AttributeID": uuid.New().String(),
	}
	data, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/products/"+productID.String()+"/direct-attributes", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestProductHandler_CreateAttribute_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := application.NewProductService(
		&stubProductRepo{},
		&stubBrandRepo{},
		&stubGroupRepo{},
		&stubAttributeRepo{findByCodeFn: func(ctx context.Context, code string) (*domain.Attribute, error) {
			return nil, assert.AnError
		}},
		&stubVariantRepo{},
		&stubPartyServiceConfigRepo{},
	)
	handler := NewProductHandler(service)
	router := newTestRouter(handler)

	body := map[string]interface{}{
		"Name": "Color",
		"Code": "C",
	}
	data, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/attributes", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestProductHandler_UpdateAttribute_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	attrID := uuid.New()
	service := application.NewProductService(
		&stubProductRepo{},
		&stubBrandRepo{},
		&stubGroupRepo{},
		&stubAttributeRepo{findByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Attribute, error) {
			return nil, nil
		}},
		&stubVariantRepo{},
		&stubPartyServiceConfigRepo{},
	)
	handler := NewProductHandler(service)
	router := newTestRouter(handler)

	req := httptest.NewRequest(http.MethodPut, "/attributes/"+attrID.String(), bytes.NewBufferString("{\"Name\":\"Shade\"}"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestProductHandler_GetProductVariantBySKU_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := application.NewProductService(
		&stubProductRepo{},
		&stubBrandRepo{},
		&stubGroupRepo{},
		&stubAttributeRepo{},
		&stubVariantRepo{findBySKFn: func(ctx context.Context, sku string) (*domain.ProductVariant, error) {
			return nil, nil
		}},
		&stubPartyServiceConfigRepo{},
	)
	handler := NewProductHandler(service)
	router := newTestRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/variants?sku=missing", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestProductHandler_ListPartyServiceConfigurations_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	partyID := uuid.New()
	service := application.NewProductService(
		&stubProductRepo{},
		&stubBrandRepo{},
		&stubGroupRepo{},
		&stubAttributeRepo{},
		&stubVariantRepo{},
		&stubPartyServiceConfigRepo{findByPartyFn: func(ctx context.Context, id uuid.UUID) ([]*domain.PartyServiceConfiguration, error) {
			return nil, assert.AnError
		}},
	)
	handler := NewProductHandler(service)
	router := newTestRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/parties/"+partyID.String()+"/configurations", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestProductHandler_ListProducts_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := application.NewProductService(
		&stubProductRepo{findAllFn: func(ctx context.Context) ([]*domain.Product, error) {
			return nil, assert.AnError
		}},
		&stubBrandRepo{},
		&stubGroupRepo{},
		&stubAttributeRepo{},
		&stubVariantRepo{},
		&stubPartyServiceConfigRepo{},
	)
	handler := NewProductHandler(service)
	router := newTestRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestProductHandler_GetProductByID_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := application.NewProductService(
		&stubProductRepo{findByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
			return nil, nil
		}},
		&stubBrandRepo{},
		&stubGroupRepo{},
		&stubAttributeRepo{},
		&stubVariantRepo{},
		&stubPartyServiceConfigRepo{},
	)
	handler := NewProductHandler(service)
	router := newTestRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/products/"+uuid.New().String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestProductHandler_GetProductByID_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := application.NewProductService(
		&stubProductRepo{findByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
			return nil, assert.AnError
		}},
		&stubBrandRepo{},
		&stubGroupRepo{},
		&stubAttributeRepo{},
		&stubVariantRepo{},
		&stubPartyServiceConfigRepo{},
	)
	handler := NewProductHandler(service)
	router := newTestRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/products/"+uuid.New().String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestProductHandler_ListAttributes_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := application.NewProductService(
		&stubProductRepo{},
		&stubBrandRepo{},
		&stubGroupRepo{},
		&stubAttributeRepo{findByScope: func(ctx context.Context, brandID *uuid.UUID, groupID *uuid.UUID) ([]*domain.Attribute, error) {
			return nil, assert.AnError
		}},
		&stubVariantRepo{},
		&stubPartyServiceConfigRepo{},
	)
	handler := NewProductHandler(service)
	router := newTestRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/attributes", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestProductHandler_UpdatePartyServiceConfiguration_MissingActorID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewProductHandler(nil)
	router := newTestRouterWithoutActor(handler)

	req := httptest.NewRequest(http.MethodPut, "/parties/"+uuid.New().String()+"/configurations/"+uuid.New().String(), bytes.NewBufferString("{}"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestProductHandler_UpdatePartyServiceConfiguration_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	partyID := uuid.New()
	config, _ := domain.NewPartyServiceConfiguration(partyID, "svc", "Config", json.RawMessage(`{"k":"v"}`))
	service := application.NewProductService(
		&stubProductRepo{},
		&stubBrandRepo{},
		&stubGroupRepo{},
		&stubAttributeRepo{},
		&stubVariantRepo{},
		&stubPartyServiceConfigRepo{
			findByIDFn: func(ctx context.Context, pid, id uuid.UUID) (*domain.PartyServiceConfiguration, error) {
				return config, nil
			},
			saveFn: func(ctx context.Context, cfg *domain.PartyServiceConfiguration) error {
				return assert.AnError
			},
		},
	)
	handler := NewProductHandler(service)
	router := newTestRouter(handler)

	req := httptest.NewRequest(http.MethodPut, "/parties/"+partyID.String()+"/configurations/"+config.ID.String(), bytes.NewBufferString("{\"name\":\"New\"}"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestProductHandler_DeletePartyServiceConfiguration_MissingActorID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewProductHandler(nil)
	router := newTestRouterWithoutActor(handler)

	req := httptest.NewRequest(http.MethodDelete, "/parties/"+uuid.New().String()+"/configurations/"+uuid.New().String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestProductHandler_DeletePartyServiceConfiguration_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	partyID := uuid.New()
	configID := uuid.New()
	service := application.NewProductService(
		&stubProductRepo{},
		&stubBrandRepo{},
		&stubGroupRepo{},
		&stubAttributeRepo{},
		&stubVariantRepo{},
		&stubPartyServiceConfigRepo{deleteFn: func(ctx context.Context, pid, id uuid.UUID) error {
			return assert.AnError
		}},
	)
	handler := NewProductHandler(service)
	router := newTestRouter(handler)

	req := httptest.NewRequest(http.MethodDelete, "/parties/"+partyID.String()+"/configurations/"+configID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestProductHandler_GetPartyServiceConfigurationByID_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	partyID := uuid.New()
	configID := uuid.New()
	service := application.NewProductService(
		&stubProductRepo{},
		&stubBrandRepo{},
		&stubGroupRepo{},
		&stubAttributeRepo{},
		&stubVariantRepo{},
		&stubPartyServiceConfigRepo{findByIDFn: func(ctx context.Context, pid, id uuid.UUID) (*domain.PartyServiceConfiguration, error) {
			return nil, nil
		}},
	)
	handler := NewProductHandler(service)
	router := newTestRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/parties/"+partyID.String()+"/configurations/"+configID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestProductHandler_CreatePartyServiceConfiguration_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	partyID := uuid.New()
	service := application.NewProductService(
		&stubProductRepo{},
		&stubBrandRepo{},
		&stubGroupRepo{},
		&stubAttributeRepo{},
		&stubVariantRepo{},
		&stubPartyServiceConfigRepo{saveFn: func(ctx context.Context, cfg *domain.PartyServiceConfiguration) error {
			return assert.AnError
		}},
	)
	handler := NewProductHandler(service)
	router := newTestRouter(handler)

	req := httptest.NewRequest(http.MethodPost, "/parties/"+partyID.String()+"/configurations", bytes.NewBufferString("{\"serviceId\":\"svc\",\"name\":\"Config\"}"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
