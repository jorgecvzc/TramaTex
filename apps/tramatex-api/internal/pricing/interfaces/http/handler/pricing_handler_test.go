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

type fakePricingRuleRepo struct {
	saved      []*domain.PricingRule
	list       []*domain.PricingRule
	applicable []*domain.PricingRule
}

func (f *fakePricingRuleRepo) Save(ctx context.Context, rule *domain.PricingRule) error {
	f.saved = append(f.saved, rule)
	return nil
}

func (f *fakePricingRuleRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.PricingRule, error) {
	return nil, nil
}

func (f *fakePricingRuleRepo) List(ctx context.Context) ([]*domain.PricingRule, error) {
	return f.list, nil
}

func (f *fakePricingRuleRepo) FindApplicable(ctx context.Context, variantID uuid.UUID, quantity int, at time.Time) ([]*domain.PricingRule, error) {
	return f.applicable, nil
}

type fakeClientPricingRepo struct {
	override *domain.ClientPricing
}

func (f *fakeClientPricingRepo) Save(ctx context.Context, override *domain.ClientPricing) error {
	f.override = override
	return nil
}

func (f *fakeClientPricingRepo) FindApplicable(ctx context.Context, clientID uuid.UUID, variantID uuid.UUID, at time.Time) (*domain.ClientPricing, error) {
	return f.override, nil
}

type fakeBrandMarginRepo struct{}

func (f *fakeBrandMarginRepo) Save(ctx context.Context, margin *domain.BrandProfitMargin) error {
	return nil
}

func (f *fakeBrandMarginRepo) FindApplicable(ctx context.Context, brandID uuid.UUID, at time.Time) (*domain.BrandProfitMargin, error) {
	return nil, nil
}

type fakeDiscountRuleRepo struct{}

func (f *fakeDiscountRuleRepo) Save(ctx context.Context, rule *domain.SalesDiscountRule) error {
	return nil
}

func (f *fakeDiscountRuleRepo) FindApplicable(ctx context.Context, clientID uuid.UUID, variantID uuid.UUID, quantity int, at time.Time) ([]*domain.SalesDiscountRule, error) {
	return nil, nil
}

type fakeCalcRepo struct{}

func (f *fakeCalcRepo) Save(ctx context.Context, calc *domain.PriceCalculation) error {
	return nil
}

func (f *fakeCalcRepo) ListByProductVariantID(ctx context.Context, variantID uuid.UUID) ([]*domain.PriceCalculation, error) {
	return nil, nil
}

type fakeProductInfoProvider struct {
	info *application.ProductPricingInfo
}

func (f *fakeProductInfoProvider) GetVariantPricingInfo(ctx context.Context, variantID uuid.UUID) (*application.ProductPricingInfo, error) {
	return f.info, nil
}

func newPricingService(info *application.ProductPricingInfo, override *domain.ClientPricing) *application.PricingService {
	return application.NewPricingService(
		&fakePricingRuleRepo{},
		&fakeClientPricingRepo{override: override},
		&fakeBrandMarginRepo{},
		&fakeDiscountRuleRepo{},
		&fakeCalcRepo{},
		&fakeProductInfoProvider{info: info},
	)
}

func performRequest(t *testing.T, handlerFunc func(*gin.Context), method, path string, body interface{}) *httptest.ResponseRecorder {
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

func TestPricingHandler_CalculatePrice(t *testing.T) {
	service := newPricingService(nil, nil)
	h := NewPricingHandler(service)

	resp := performRequest(t, h.CalculatePrice, http.MethodPost, "/pricing/calculate", "{")
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.Code)
	}

	resp = performRequest(t, h.CalculatePrice, http.MethodPost, "/pricing/calculate", application.CalculatePriceCommand{})
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.Code)
	}

	variantID := uuid.New()
	clientID := uuid.New()
	info := &application.ProductPricingInfo{VariantID: variantID, ProductID: uuid.New(), BaseCost: 10, Currency: "EUR", BrandID: uuid.New()}
	price, _ := domain.NewMoney(12, "EUR")
	override, _ := domain.NewClientPricing(clientID, variantID, price, time.Now().Add(-time.Hour), nil)
	service = newPricingService(info, override)
	h = NewPricingHandler(service)

	resp = performRequest(t, h.CalculatePrice, http.MethodPost, "/pricing/calculate", application.CalculatePriceCommand{
		ProductVariantID: variantID,
		ClientID:         clientID,
		Quantity:         1,
	})
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}
}

func TestPricingHandler_CreatePricingRule(t *testing.T) {
	service := newPricingService(nil, nil)
	h := NewPricingHandler(service)

	resp := performRequest(t, h.CreatePricingRule, http.MethodPost, "/pricing/rules", "{")
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.Code)
	}

	resp = performRequest(t, h.CreatePricingRule, http.MethodPost, "/pricing/rules", application.CreatePricingRuleCommand{
		Name:             "Rule",
		MarkupPercentage: 0.1,
		MinQuantity:      1,
	})
	if resp.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.Code)
	}
}

func TestPricingHandler_ListPricingRules(t *testing.T) {
	repo := &fakePricingRuleRepo{}
	service := application.NewPricingService(repo, &fakeClientPricingRepo{}, &fakeBrandMarginRepo{}, &fakeDiscountRuleRepo{}, &fakeCalcRepo{}, &fakeProductInfoProvider{})
	h := NewPricingHandler(service)

	resp := performRequest(t, h.ListPricingRules, http.MethodGet, "/pricing/rules", nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}
}

func TestPricingHandler_CreateClientPricingOverride(t *testing.T) {
	service := newPricingService(nil, nil)
	h := NewPricingHandler(service)

	resp := performRequest(t, h.CreateClientPricingOverride, http.MethodPost, "/pricing/overrides", "{")
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.Code)
	}

	resp = performRequest(t, h.CreateClientPricingOverride, http.MethodPost, "/pricing/overrides", application.CreateClientPricingCommand{
		ClientID:         uuid.New(),
		ProductVariantID: uuid.New(),
		FixedPrice:       20,
		Currency:         "EUR",
	})
	if resp.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.Code)
	}
}

func TestPricingHandler_GetPricingHistory(t *testing.T) {
	service := newPricingService(nil, nil)
	h := NewPricingHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/pricing/history/invalid", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = []gin.Param{{Key: "variantId", Value: "invalid"}}

	h.GetPricingHistory(c)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
