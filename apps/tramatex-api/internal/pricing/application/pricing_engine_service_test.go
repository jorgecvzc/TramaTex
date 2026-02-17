package application_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/joran-cortez/tramatex/internal/pricing/application"
	"github.com/joran-cortez/tramatex/internal/pricing/domain"
)

type fakeBaseRuleRepo struct {
	saved []*domain.BaseSalesPriceRule
	rules []*domain.BaseSalesPriceRule
	find  map[uuid.UUID]*domain.BaseSalesPriceRule
}

func (f *fakeBaseRuleRepo) Save(ctx context.Context, rule *domain.BaseSalesPriceRule) error {
	f.saved = append(f.saved, rule)
	return nil
}

func (f *fakeBaseRuleRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.BaseSalesPriceRule, error) {
	if f.find == nil {
		return nil, nil
	}
	return f.find[id], nil
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

func (f *fakeSaleRuleRepo) ListApplicable(ctx context.Context, clientID uuid.UUID, productGroupID *uuid.UUID, orderTotal domain.Money, at time.Time) ([]*domain.SaleModificationRule, error) {
	return f.rules, nil
}

type fakeBasePriceCache struct {
	prices map[string]domain.Money
	sets   int
}

func (f *fakeBasePriceCache) key(productID uuid.UUID, variantID uuid.UUID) string {
	return productID.String() + ":" + variantID.String()
}

func (f *fakeBasePriceCache) GetBasePrice(ctx context.Context, productID uuid.UUID, variantID uuid.UUID) (*domain.Money, error) {
	if f.prices == nil {
		return nil, nil
	}
	if price, ok := f.prices[f.key(productID, variantID)]; ok {
		return &price, nil
	}
	return nil, nil
}

func (f *fakeBasePriceCache) SetBasePrice(ctx context.Context, productID uuid.UUID, variantID uuid.UUID, price domain.Money) error {
	if f.prices == nil {
		f.prices = make(map[string]domain.Money)
	}
	f.prices[f.key(productID, variantID)] = price
	f.sets++
	return nil
}

type fakeProductProvider struct {
	info  *application.ProductPricingInfo
	infos []*application.ProductPricingInfo
}

func (f *fakeProductProvider) GetVariantPricingInfo(ctx context.Context, variantID uuid.UUID) (*application.ProductPricingInfo, error) {
	return f.info, nil
}

func (f *fakeProductProvider) ListVariantsPricingInfo(ctx context.Context, productID uuid.UUID) ([]*application.ProductPricingInfo, error) {
	return f.infos, nil
}

func TestPricingEngineService_CreateBaseSalesPriceRule(t *testing.T) {
	repo := &fakeBaseRuleRepo{}
	service := application.NewPricingEngineService(repo, &fakeSaleRuleRepo{}, &fakeProductProvider{}, nil)

	value := application.RuleValueDTO{Type: string(domain.RuleValuePercentageMarkup), PercentageValue: &application.PercentageDTO{Value: 0.1}}
	_, err := service.CreateBaseSalesPriceRule(context.Background(), application.CreateBaseSalesPriceRuleCommand{Name: "Rule", Value: value})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if len(repo.saved) != 1 {
		t.Fatalf("expected rule saved")
	}
}

func TestPricingEngineService_UpdateBaseSalesPriceRule_NotFound(t *testing.T) {
	repo := &fakeBaseRuleRepo{}
	service := application.NewPricingEngineService(repo, &fakeSaleRuleRepo{}, &fakeProductProvider{}, nil)

	_, err := service.UpdateBaseSalesPriceRule(context.Background(), application.UpdateBaseSalesPriceRuleCommand{ID: uuid.New()})
	if err == nil {
		t.Fatalf("expected not found error")
	}
}

func TestPricingEngineService_CalculateBaseSalesPrice_CacheSet(t *testing.T) {
	repo := &fakeBaseRuleRepo{}
	cache := &fakeBasePriceCache{}
	provider := &fakeProductProvider{}
	service := application.NewPricingEngineService(repo, &fakeSaleRuleRepo{}, provider, cache)

	brandID := uuid.New()
	variantID := uuid.New()
	productID := uuid.New()
	percentage, _ := domain.NewPercentage(0.1)
	value, _ := domain.NewRuleValue(domain.RuleValuePercentageMarkup, &percentage, nil)
	rule, _ := domain.NewBaseSalesPriceRule("Rule", &brandID, nil, &productID, &variantID, value)
	repo.rules = []*domain.BaseSalesPriceRule{rule}

	provider.info = &application.ProductPricingInfo{
		VariantID: variantID,
		ProductID: productID,
		BaseCost:  100,
		Currency:  "EUR",
		BrandID:   brandID,
	}
	provider.infos = []*application.ProductPricingInfo{provider.info}

	_, err := service.CalculateBaseSalesPrice(context.Background(), application.CalculateBaseSalesPriceRequest{VariantID: variantID})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if cache.sets == 0 {
		t.Fatalf("expected cache to be set")
	}
}

func TestPricingEngineService_CalculateFinalSalePrice_UsesRule(t *testing.T) {
	repo := &fakeBaseRuleRepo{}
	saleRepo := &fakeSaleRuleRepo{}
	cache := &fakeBasePriceCache{}
	provider := &fakeProductProvider{}
	service := application.NewPricingEngineService(repo, saleRepo, provider, cache)

	variantID := uuid.New()
	productID := uuid.New()
	clientID := uuid.New()
	groupID := uuid.New()

	basePrice, _ := domain.NewMoney(100, "EUR")
	cache.SetBasePrice(context.Background(), productID, variantID, basePrice)

	discount, _ := domain.NewMoney(10, "EUR")
	value, _ := domain.NewRuleValue(domain.RuleValueApplyFixedAmountDiscount, nil, &discount)
	rule, _ := domain.NewSaleModificationRule("Promo", []uuid.UUID{clientID}, &groupID, nil, value, 1, time.Now().Add(-time.Hour), nil)
	saleRepo.rules = []*domain.SaleModificationRule{rule}

	provider.info = &application.ProductPricingInfo{
		VariantID: variantID,
		ProductID: productID,
		BaseCost:  100,
		Currency:  "EUR",
		BrandID:   uuid.New(),
		GroupIDs:  []uuid.UUID{groupID},
	}

	resp, err := service.CalculateFinalSalePrice(context.Background(), application.CalculateFinalSalePriceRequest{
		ClientID: clientID,
		SaleItems: []application.SaleItemRequest{{
			ProductVariantID: variantID,
			Quantity:         2,
		}},
	})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if len(resp.CalculatedItems) != 1 {
		t.Fatalf("expected 1 item")
	}
	if resp.SaleTotal.Amount != 180 {
		t.Fatalf("expected discounted total, got %.2f", resp.SaleTotal.Amount)
	}
}

func TestPricingEngineService_CalculateFinalSalePrice_Validation(t *testing.T) {
	service := application.NewPricingEngineService(&fakeBaseRuleRepo{}, &fakeSaleRuleRepo{}, &fakeProductProvider{}, nil)

	_, err := service.CalculateFinalSalePrice(context.Background(), application.CalculateFinalSalePriceRequest{})
	if err == nil {
		t.Fatalf("expected validation error")
	}
}
