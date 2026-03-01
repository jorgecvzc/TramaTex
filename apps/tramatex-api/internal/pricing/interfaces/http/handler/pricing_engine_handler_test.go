package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/joran-cortez/tramatex/internal/pricing/application"
	"github.com/joran-cortez/tramatex/internal/pricing/domain"
)

type fakeBaseRuleRepo struct {
	rules []*domain.BaseSalesPriceRule
}

func (f *fakeBaseRuleRepo) Save(ctx context.Context, rule *domain.BaseSalesPriceRule) error {
	f.rules = append(f.rules, rule)
	return nil
}

func (f *fakeBaseRuleRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.BaseSalesPriceRule, error) {
	return nil, nil
}

func (f *fakeBaseRuleRepo) List(ctx context.Context) ([]*domain.BaseSalesPriceRule, error) {
	return f.rules, nil
}

type fakeSaleRuleRepo struct {
	rules []*domain.SaleModificationRule
}

func (f *fakeSaleRuleRepo) Save(ctx context.Context, rule *domain.SaleModificationRule) error {
	f.rules = append(f.rules, rule)
	return nil
}

func (f *fakeSaleRuleRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.SaleModificationRule, error) {
	return nil, nil
}

func (f *fakeSaleRuleRepo) ListApplicable(ctx context.Context, clientID string, productGroupID *uuid.UUID, orderTotal domain.Money, at time.Time) ([]*domain.SaleModificationRule, error) {
	return f.rules, nil
}

type fakeProductProvider struct {
	info *application.ProductPricingInfo
}

func (f *fakeProductProvider) GetVariantPricingInfo(ctx context.Context, variantID uuid.UUID) (*application.ProductPricingInfo, error) {
	return f.info, nil
}

func performEngineRequest(t *testing.T, handlerFunc func(*gin.Context), method, path string, body interface{}) *httptest.ResponseRecorder {
	t.Helper()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	var payload []byte
	switch v := body.(type) {
	case string:
		payload = []byte(v)
	case nil:
		payload = nil
	default:
		data, err := json.Marshal(v)
		if err != nil {
			t.Fatalf("failed to marshal body: %v", err)
		}
		payload = data
	}

	req := httptest.NewRequest(method, path, bytes.NewReader(payload))
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	c.Request = req

	handlerFunc(c)
	c.Writer.WriteHeaderNow()
	return w
}

func TestPricingEngineHandler_CalculateBaseSalesPrice(t *testing.T) {
	baseRepo := &fakeBaseRuleRepo{}
	service := application.NewPricingEngineService(baseRepo, &fakeSaleRuleRepo{}, &fakeProductProvider{}, nil)
	h := NewPricingEngineHandler(service)

	resp := performEngineRequest(t, h.CalculateBaseSalesPrice, http.MethodPost, "/pricing/base", "{")
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.Code)
	}

	resp = performEngineRequest(t, h.CalculateBaseSalesPrice, http.MethodPost, "/pricing/base", application.CalculateBaseSalesPriceRequest{})
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.Code)
	}
}

func TestPricingEngineHandler_CreateBaseSalesPriceRule(t *testing.T) {
	baseRepo := &fakeBaseRuleRepo{}
	service := application.NewPricingEngineService(baseRepo, &fakeSaleRuleRepo{}, &fakeProductProvider{}, nil)
	h := NewPricingEngineHandler(service)

	resp := performEngineRequest(t, h.CreateBaseSalesPriceRule, http.MethodPost, "/pricing/base-rules", "{")
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.Code)
	}

	req := application.CreateBaseSalesPriceRuleCommand{
		Name: "Rule",
		Value: application.RuleValueDTO{
			Type:            string(domain.RuleValuePercentageMarkup),
			PercentageValue: &application.PercentageDTO{Value: 0.1},
		},
	}
	resp = performEngineRequest(t, h.CreateBaseSalesPriceRule, http.MethodPost, "/pricing/base-rules", req)
	if resp.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.Code)
	}
}

func TestPricingEngineHandler_UpdateBaseSalesPriceRule_InvalidID(t *testing.T) {
	service := application.NewPricingEngineService(&fakeBaseRuleRepo{}, &fakeSaleRuleRepo{}, &fakeProductProvider{}, nil)
	h := NewPricingEngineHandler(service)

	req := httptest.NewRequest(http.MethodPut, "/pricing/base-rules/invalid", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = []gin.Param{{Key: "id", Value: "invalid"}}

	h.UpdateBaseSalesPriceRule(c)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestPricingEngineHandler_CalculateFinalSalePrice(t *testing.T) {
	variantID := uuid.New()
	productID := uuid.New()
	clientID := uuid.New().String()

	provider := &fakeProductProvider{info: &application.ProductPricingInfo{
		VariantID: variantID,
		ProductID: productID,
		BaseCost:  10,
		Currency:  "EUR",
		BrandID:   uuid.New(),
	}}
	service := application.NewPricingEngineService(&fakeBaseRuleRepo{}, &fakeSaleRuleRepo{}, provider, nil)
	h := NewPricingEngineHandler(service)

	resp := performEngineRequest(t, h.CalculateFinalSalePrice, http.MethodPost, "/pricing/final", "{")
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.Code)
	}

	resp = performEngineRequest(t, h.CalculateFinalSalePrice, http.MethodPost, "/pricing/final", application.CalculateFinalSalePriceRequest{})
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.Code)
	}

	resp = performEngineRequest(t, h.CalculateFinalSalePrice, http.MethodPost, "/pricing/final", application.CalculateFinalSalePriceRequest{
		ClientID:  clientID,
		SaleItems: []application.SaleItemRequest{{ProductVariantID: variantID, Quantity: 1}},
	})
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}
}
