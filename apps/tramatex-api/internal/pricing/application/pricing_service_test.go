package application_test

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/joran-cortez/tramatex/internal/pricing/application"
	"github.com/joran-cortez/tramatex/internal/pricing/domain"
)

type fakePricingRuleRepo struct {
	saved      []*domain.PricingRule
	list       []*domain.PricingRule
	listErr    error
	applicable []*domain.PricingRule
	findErr    error
}

func (f *fakePricingRuleRepo) Save(ctx context.Context, rule *domain.PricingRule) error {
	f.saved = append(f.saved, rule)
	return nil
}

func (f *fakePricingRuleRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.PricingRule, error) {
	return nil, f.findErr
}

func (f *fakePricingRuleRepo) List(ctx context.Context) ([]*domain.PricingRule, error) {
	return f.list, f.listErr
}

func (f *fakePricingRuleRepo) FindApplicable(ctx context.Context, variantID uuid.UUID, quantity int, at time.Time) ([]*domain.PricingRule, error) {
	return f.applicable, nil
}

type fakeClientPricingRepo struct {
	saved       []*domain.ClientPricing
	override    *domain.ClientPricing
	overrideErr error
}

func (f *fakeClientPricingRepo) Save(ctx context.Context, override *domain.ClientPricing) error {
	f.saved = append(f.saved, override)
	return nil
}

func (f *fakeClientPricingRepo) FindApplicable(ctx context.Context, clientID uuid.UUID, variantID uuid.UUID, at time.Time) (*domain.ClientPricing, error) {
	return f.override, f.overrideErr
}

type fakeBrandMarginRepo struct {
	margin *domain.BrandProfitMargin
}

func (f *fakeBrandMarginRepo) Save(ctx context.Context, margin *domain.BrandProfitMargin) error {
	f.margin = margin
	return nil
}

func (f *fakeBrandMarginRepo) FindApplicable(ctx context.Context, brandID uuid.UUID, at time.Time) (*domain.BrandProfitMargin, error) {
	return f.margin, nil
}

type fakeDiscountRuleRepo struct {
	rules []*domain.SalesDiscountRule
}

func (f *fakeDiscountRuleRepo) Save(ctx context.Context, rule *domain.SalesDiscountRule) error {
	f.rules = append(f.rules, rule)
	return nil
}

func (f *fakeDiscountRuleRepo) FindApplicable(ctx context.Context, clientID uuid.UUID, variantID uuid.UUID, quantity int, at time.Time) ([]*domain.SalesDiscountRule, error) {
	return f.rules, nil
}

type fakeCalculationRepo struct {
	saved []*domain.PriceCalculation
	list  []*domain.PriceCalculation
}

func (f *fakeCalculationRepo) Save(ctx context.Context, calc *domain.PriceCalculation) error {
	f.saved = append(f.saved, calc)
	return nil
}

func (f *fakeCalculationRepo) ListByProductVariantID(ctx context.Context, variantID uuid.UUID) ([]*domain.PriceCalculation, error) {
	return f.list, nil
}

type fakeProductInfoProvider struct {
	info *application.ProductPricingInfo
	err  error
}

func (f *fakeProductInfoProvider) GetVariantPricingInfo(ctx context.Context, variantID uuid.UUID) (*application.ProductPricingInfo, error) {
	return f.info, f.err
}

func newPricingService(
	pricingRepo *fakePricingRuleRepo,
	clientRepo *fakeClientPricingRepo,
	brandRepo *fakeBrandMarginRepo,
	discountRepo *fakeDiscountRuleRepo,
	calcRepo *fakeCalculationRepo,
	productInfo *fakeProductInfoProvider,
) *application.PricingService {
	return application.NewPricingService(pricingRepo, clientRepo, brandRepo, discountRepo, calcRepo, productInfo)
}

func TestPricingService_CreatePricingRule_DefaultEffectiveFrom(t *testing.T) {
	repo := &fakePricingRuleRepo{}
	service := newPricingService(repo, &fakeClientPricingRepo{}, &fakeBrandMarginRepo{}, &fakeDiscountRuleRepo{}, &fakeCalculationRepo{}, &fakeProductInfoProvider{})

	cmd := application.CreatePricingRuleCommand{
		Name:             "Rule",
		MarkupPercentage: 0.1,
		MinQuantity:      1,
	}
	_, err := service.CreatePricingRule(context.Background(), cmd)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if len(repo.saved) != 1 {
		t.Fatalf("expected rule to be saved")
	}
	if repo.saved[0].EffectiveFrom.IsZero() {
		t.Fatalf("expected effectiveFrom to be set")
	}
}

func TestPricingService_CreateClientPricing_DefaultEffectiveFrom(t *testing.T) {
	clientRepo := &fakeClientPricingRepo{}
	service := newPricingService(&fakePricingRuleRepo{}, clientRepo, &fakeBrandMarginRepo{}, &fakeDiscountRuleRepo{}, &fakeCalculationRepo{}, &fakeProductInfoProvider{})

	clientID := uuid.New()
	variantID := uuid.New()
	cmd := application.CreateClientPricingCommand{
		ClientID:         clientID,
		ProductVariantID: variantID,
		FixedPrice:       100,
		Currency:         "EUR",
	}

	_, err := service.CreateClientPricing(context.Background(), cmd)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if len(clientRepo.saved) != 1 {
		t.Fatalf("expected override to be saved")
	}
	if clientRepo.saved[0].EffectiveFrom.IsZero() {
		t.Fatalf("expected effectiveFrom to be set")
	}
}

func TestPricingService_CalculatePrice_ClientOverride(t *testing.T) {
	clientID := uuid.New()
	variantID := uuid.New()
	productID := uuid.New()
	brandID := uuid.New()

	overridePrice, _ := domain.NewMoney(90, "EUR")
	override, _ := domain.NewClientPricing(clientID, variantID, overridePrice, time.Now(), nil)

	clientRepo := &fakeClientPricingRepo{override: override}
	calcRepo := &fakeCalculationRepo{}
	productInfo := &fakeProductInfoProvider{info: &application.ProductPricingInfo{
		VariantID: variantID,
		ProductID: productID,
		BaseCost:  50,
		Currency:  "EUR",
		BrandID:   brandID,
	}}

	service := newPricingService(&fakePricingRuleRepo{}, clientRepo, &fakeBrandMarginRepo{}, &fakeDiscountRuleRepo{}, calcRepo, productInfo)
	resp, err := service.CalculatePrice(context.Background(), application.CalculatePriceCommand{
		ProductVariantID: variantID,
		ClientID:         clientID,
		Quantity:         1,
	})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if resp.FinalPrice != overridePrice.Amount() {
		t.Fatalf("expected override price applied")
	}
	if len(calcRepo.saved) != 1 {
		t.Fatalf("expected calculation saved")
	}
	if len(calcRepo.saved[0].AppliedRules) == 0 || calcRepo.saved[0].AppliedRules[0] != "ClientOverride" {
		t.Fatalf("expected ClientOverride rule applied")
	}
}

func TestPricingService_CalculatePrice_PricingRule(t *testing.T) {
	variantID := uuid.New()
	clientID := uuid.New()
	brandID := uuid.New()
	productID := uuid.New()

	percentage, _ := domain.NewPercentage(0.1)
	rule, _ := domain.NewPricingRule("Rule", &variantID, nil, percentage, 1, nil, time.Now().Add(-time.Hour), nil)

	pricingRepo := &fakePricingRuleRepo{applicable: []*domain.PricingRule{rule}}
	calcRepo := &fakeCalculationRepo{}
	productInfo := &fakeProductInfoProvider{info: &application.ProductPricingInfo{
		VariantID: variantID,
		ProductID: productID,
		BaseCost:  100,
		Currency:  "EUR",
		BrandID:   brandID,
	}}

	service := newPricingService(pricingRepo, &fakeClientPricingRepo{}, &fakeBrandMarginRepo{}, &fakeDiscountRuleRepo{}, calcRepo, productInfo)
	resp, err := service.CalculatePrice(context.Background(), application.CalculatePriceCommand{
		ProductVariantID: variantID,
		ClientID:         clientID,
		Quantity:         1,
	})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if math.Abs(resp.FinalPrice-110) > 0.0001 {
		t.Fatalf("expected markup applied, got %.2f", resp.FinalPrice)
	}
}

func TestPricingService_CalculatePrice_BrandMarginAndDiscount(t *testing.T) {
	variantID := uuid.New()
	clientID := uuid.New()
	brandID := uuid.New()
	productID := uuid.New()

	percentage, _ := domain.NewPercentage(0.2)
	fixed, _ := domain.NewMoney(5, "EUR")
	margin, _ := domain.NewBrandProfitMargin(brandID, &percentage, &fixed, time.Now().Add(-time.Hour), nil)

	discountValue, _ := domain.NewPercentage(0.1)
	discount, _ := domain.NewSalesDiscountRule("Seasonal", &clientID, &variantID, nil, domain.DiscountTypePercentage, &discountValue, nil, 1, time.Now().Add(-time.Hour), nil)

	calcRepo := &fakeCalculationRepo{}
	productInfo := &fakeProductInfoProvider{info: &application.ProductPricingInfo{
		VariantID: variantID,
		ProductID: productID,
		BaseCost:  100,
		Currency:  "EUR",
		BrandID:   brandID,
	}}

	service := newPricingService(&fakePricingRuleRepo{}, &fakeClientPricingRepo{}, &fakeBrandMarginRepo{margin: margin}, &fakeDiscountRuleRepo{rules: []*domain.SalesDiscountRule{discount}}, calcRepo, productInfo)
	_, err := service.CalculatePrice(context.Background(), application.CalculatePriceCommand{
		ProductVariantID: variantID,
		ClientID:         clientID,
		Quantity:         2,
	})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if len(calcRepo.saved) != 1 {
		t.Fatalf("expected calculation saved")
	}
	applied := calcRepo.saved[0].AppliedRules
	foundMargin := false
	foundDiscount := false
	for _, name := range applied {
		if name == "BrandMargin" {
			foundMargin = true
		}
		if name == "Seasonal" {
			foundDiscount = true
		}
	}
	if !foundMargin || !foundDiscount {
		t.Fatalf("expected BrandMargin and discount in applied rules")
	}
}

func TestPricingService_CalculatePrice_InvalidQuantity(t *testing.T) {
	service := newPricingService(&fakePricingRuleRepo{}, &fakeClientPricingRepo{}, &fakeBrandMarginRepo{}, &fakeDiscountRuleRepo{}, &fakeCalculationRepo{}, &fakeProductInfoProvider{})

	_, err := service.CalculatePrice(context.Background(), application.CalculatePriceCommand{Quantity: 0})
	if err == nil {
		t.Fatalf("expected validation error")
	}
}

func TestPricingService_GetPricingHistory(t *testing.T) {
	variantID := uuid.New()
	money, _ := domain.NewMoney(10, "EUR")
	calc, _ := domain.NewPriceCalculation(variantID, uuid.New(), 1, money, money, []string{"Rule"})

	calcRepo := &fakeCalculationRepo{list: []*domain.PriceCalculation{calc}}
	service := newPricingService(&fakePricingRuleRepo{}, &fakeClientPricingRepo{}, &fakeBrandMarginRepo{}, &fakeDiscountRuleRepo{}, calcRepo, &fakeProductInfoProvider{})

	results, err := service.GetPricingHistory(context.Background(), application.GetPricingHistoryQuery{ProductVariantID: variantID})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}
