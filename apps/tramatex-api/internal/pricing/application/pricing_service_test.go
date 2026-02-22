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
	saved   []*domain.PriceCalculation
	list    []*domain.PriceCalculation
	saveErr error
}

func (f *fakeCalculationRepo) Save(ctx context.Context, calc *domain.PriceCalculation) error {
	if f.saveErr != nil {
		return f.saveErr
	}
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

func TestPricingService_ListPricingRules_Success(t *testing.T) {
	percentage, _ := domain.NewPercentage(0.15)
	rule1, _ := domain.NewPricingRule("Rule1", nil, nil, percentage, 1, nil, time.Now(), nil)
	rule2, _ := domain.NewPricingRule("Rule2", nil, nil, percentage, 10, nil, time.Now(), nil)

	repo := &fakePricingRuleRepo{list: []*domain.PricingRule{rule1, rule2}}
	service := newPricingService(repo, &fakeClientPricingRepo{}, &fakeBrandMarginRepo{}, &fakeDiscountRuleRepo{}, &fakeCalculationRepo{}, &fakeProductInfoProvider{})

	results, err := service.ListPricingRules(context.Background(), application.ListPricingRulesQuery{})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(results))
	}
	if results[0].Name != "Rule1" || results[1].Name != "Rule2" {
		t.Fatalf("unexpected rule names")
	}
}

func TestPricingService_ListPricingRules_EmptyList(t *testing.T) {
	repo := &fakePricingRuleRepo{list: []*domain.PricingRule{}}
	service := newPricingService(repo, &fakeClientPricingRepo{}, &fakeBrandMarginRepo{}, &fakeDiscountRuleRepo{}, &fakeCalculationRepo{}, &fakeProductInfoProvider{})

	results, err := service.ListPricingRules(context.Background(), application.ListPricingRulesQuery{})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected empty list, got %d", len(results))
	}
}

func TestPricingService_CreatePricingRule_WithEffectiveDates(t *testing.T) {
	repo := &fakePricingRuleRepo{}
	service := newPricingService(repo, &fakeClientPricingRepo{}, &fakeBrandMarginRepo{}, &fakeDiscountRuleRepo{}, &fakeCalculationRepo{}, &fakeProductInfoProvider{})

	effectiveFrom := time.Now().Add(-24 * time.Hour)
	effectiveTo := time.Now().Add(30 * 24 * time.Hour)
	variantID := uuid.New()
	maxQty := 100

	cmd := application.CreatePricingRuleCommand{
		Name:             "SeasonalRule",
		ProductVariantID: &variantID,
		MarkupPercentage: 0.25,
		MinQuantity:      5,
		MaxQuantity:      &maxQty,
		EffectiveFrom:    effectiveFrom,
		EffectiveTo:      &effectiveTo,
	}
	dto, err := service.CreatePricingRule(context.Background(), cmd)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if len(repo.saved) != 1 {
		t.Fatalf("expected rule to be saved")
	}
	if dto.Name != "SeasonalRule" {
		t.Fatalf("unexpected rule name")
	}
	if repo.saved[0].EffectiveFrom != effectiveFrom {
		t.Fatalf("effectiveFrom not preserved")
	}
}

func TestPricingService_CreateClientPricing_WithEffectiveDates(t *testing.T) {
	clientRepo := &fakeClientPricingRepo{}
	service := newPricingService(&fakePricingRuleRepo{}, clientRepo, &fakeBrandMarginRepo{}, &fakeDiscountRuleRepo{}, &fakeCalculationRepo{}, &fakeProductInfoProvider{})

	effectiveFrom := time.Now().Add(-10 * time.Hour)
	effectiveTo := time.Now().Add(72 * time.Hour)
	clientID := uuid.New()
	variantID := uuid.New()

	cmd := application.CreateClientPricingCommand{
		ClientID:         clientID,
		ProductVariantID: variantID,
		FixedPrice:       150.50,
		Currency:         "EUR",
		EffectiveFrom:    effectiveFrom,
		EffectiveTo:      &effectiveTo,
	}

	dto, err := service.CreateClientPricing(context.Background(), cmd)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if len(clientRepo.saved) != 1 {
		t.Fatalf("expected override to be saved")
	}
	if math.Abs(dto.FixedPrice.Amount-150.50) > 0.001 {
		t.Fatalf("unexpected fixed price: got %.2f", dto.FixedPrice.Amount)
	}
}

func TestPricingService_CalculatePrice_NoRulesNoOverride(t *testing.T) {
	variantID := uuid.New()
	productID := uuid.New()
	brandID := uuid.New()
	clientID := uuid.New()

	calcRepo := &fakeCalculationRepo{}
	productInfo := &fakeProductInfoProvider{info: &application.ProductPricingInfo{
		VariantID: variantID,
		ProductID: productID,
		BaseCost:  75,
		Currency:  "EUR",
		BrandID:   brandID,
	}}

	service := newPricingService(&fakePricingRuleRepo{}, &fakeClientPricingRepo{}, &fakeBrandMarginRepo{}, &fakeDiscountRuleRepo{}, calcRepo, productInfo)
	resp, err := service.CalculatePrice(context.Background(), application.CalculatePriceCommand{
		ProductVariantID: variantID,
		ClientID:         clientID,
		Quantity:         1,
	})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if resp.FinalPrice != 75 {
		t.Fatalf("expected base cost as final price, got %.2f", resp.FinalPrice)
	}
	if len(calcRepo.saved) != 1 {
		t.Fatalf("expected calculation saved")
	}
}

func TestPricingService_CalculatePrice_MultipleDiscounts(t *testing.T) {
	variantID := uuid.New()
	clientID := uuid.New()
	brandID := uuid.New()
	productID := uuid.New()

	discount1Value, _ := domain.NewPercentage(0.1)
	discount1, _ := domain.NewSalesDiscountRule("Volume", &clientID, &variantID, nil, domain.DiscountTypePercentage, &discount1Value, nil, 5, time.Now().Add(-time.Hour), nil)

	discount2Fixed, _ := domain.NewMoney(15, "EUR")
	discount2, _ := domain.NewSalesDiscountRule("Loyalty", &clientID, &variantID, nil, domain.DiscountTypeFixed, nil, &discount2Fixed, 1, time.Now().Add(-time.Hour), nil)

	calcRepo := &fakeCalculationRepo{}
	productInfo := &fakeProductInfoProvider{info: &application.ProductPricingInfo{
		VariantID: variantID,
		ProductID: productID,
		BaseCost:  200,
		Currency:  "EUR",
		BrandID:   brandID,
	}}

	service := newPricingService(&fakePricingRuleRepo{}, &fakeClientPricingRepo{}, &fakeBrandMarginRepo{}, &fakeDiscountRuleRepo{rules: []*domain.SalesDiscountRule{discount1, discount2}}, calcRepo, productInfo)
	_, err := service.CalculatePrice(context.Background(), application.CalculatePriceCommand{
		ProductVariantID: variantID,
		ClientID:         clientID,
		Quantity:         10,
	})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if len(calcRepo.saved) != 1 {
		t.Fatalf("expected calculation saved")
	}
	applied := calcRepo.saved[0].AppliedRules
	if len(applied) < 2 {
		t.Fatalf("expected at least 2 discounts applied")
	}
}

func TestPricingService_CreatePricingRule_InvalidPercentage(t *testing.T) {
	service := newPricingService(&fakePricingRuleRepo{}, &fakeClientPricingRepo{}, &fakeBrandMarginRepo{}, &fakeDiscountRuleRepo{}, &fakeCalculationRepo{}, &fakeProductInfoProvider{})

	_, err := service.CreatePricingRule(context.Background(), application.CreatePricingRuleCommand{
		Name:             "Invalid",
		MarkupPercentage: -0.5,
		MinQuantity:      1,
	})
	if err == nil {
		t.Fatalf("expected validation error for negative percentage")
	}
}

func TestPricingService_CreatePricingRule_InvalidQuantity(t *testing.T) {
	service := newPricingService(&fakePricingRuleRepo{}, &fakeClientPricingRepo{}, &fakeBrandMarginRepo{}, &fakeDiscountRuleRepo{}, &fakeCalculationRepo{}, &fakeProductInfoProvider{})

	_, err := service.CreatePricingRule(context.Background(), application.CreatePricingRuleCommand{
		Name:             "Invalid",
		MarkupPercentage: 0.2,
		MinQuantity:      -5,
	})
	if err == nil {
		t.Fatalf("expected validation error for negative min quantity")
	}
}

func TestPricingService_CalculatePrice_ProductNotFound(t *testing.T) {
	productInfo := &fakeProductInfoProvider{info: nil}
	service := newPricingService(&fakePricingRuleRepo{}, &fakeClientPricingRepo{}, &fakeBrandMarginRepo{}, &fakeDiscountRuleRepo{}, &fakeCalculationRepo{}, productInfo)

	_, err := service.CalculatePrice(context.Background(), application.CalculatePriceCommand{
		ProductVariantID: uuid.New(),
		ClientID:         uuid.New(),
		Quantity:         1,
	})
	if err == nil {
		t.Fatalf("expected not found error")
	}
}

func TestPricingService_CalculatePrice_WithQuantityThreshold(t *testing.T) {
	variantID := uuid.New()
	clientID := uuid.New()
	brandID := uuid.New()
	productID := uuid.New()

	percentage1, _ := domain.NewPercentage(0.1)
	rule1, _ := domain.NewPricingRule("LowVolume", &variantID, nil, percentage1, 1, nil, time.Now().Add(-time.Hour), nil)

	percentage2, _ := domain.NewPercentage(0.2)
	maxQty := 999
	rule2, _ := domain.NewPricingRule("HighVolume", &variantID, nil, percentage2, 10, &maxQty, time.Now().Add(-time.Hour), nil)

	pricingRepo := &fakePricingRuleRepo{applicable: []*domain.PricingRule{rule1, rule2}}
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
		Quantity:         15,
	})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if math.Abs(resp.FinalPrice-120) > 0.0001 {
		t.Fatalf("expected 20%% markup for high volume, got %.2f", resp.FinalPrice)
	}
}

func TestPricingService_CalculatePrice_BrandMarginFixedAmount(t *testing.T) {
	variantID := uuid.New()
	clientID := uuid.New()
	brandID := uuid.New()
	productID := uuid.New()

	fixed, _ := domain.NewMoney(20, "EUR")
	margin, _ := domain.NewBrandProfitMargin(brandID, nil, &fixed, time.Now().Add(-time.Hour), nil)

	calcRepo := &fakeCalculationRepo{}
	productInfo := &fakeProductInfoProvider{info: &application.ProductPricingInfo{
		VariantID: variantID,
		ProductID: productID,
		BaseCost:  100,
		Currency:  "EUR",
		BrandID:   brandID,
	}}

	service := newPricingService(&fakePricingRuleRepo{}, &fakeClientPricingRepo{}, &fakeBrandMarginRepo{margin: margin}, &fakeDiscountRuleRepo{}, calcRepo, productInfo)
	resp, err := service.CalculatePrice(context.Background(), application.CalculatePriceCommand{
		ProductVariantID: variantID,
		ClientID:         clientID,
		Quantity:         1,
	})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if math.Abs(resp.FinalPrice-120) > 0.0001 {
		t.Fatalf("expected fixed amount applied (100+20), got %.2f", resp.FinalPrice)
	}
}

func TestPricingService_CalculatePrice_BothPercentageAndFixedMargin(t *testing.T) {
	variantID := uuid.New()
	clientID := uuid.New()
	brandID := uuid.New()
	productID := uuid.New()

	percentage, _ := domain.NewPercentage(0.1)
	fixed, _ := domain.NewMoney(10, "EUR")
	margin, _ := domain.NewBrandProfitMargin(brandID, &percentage, &fixed, time.Now().Add(-time.Hour), nil)

	calcRepo := &fakeCalculationRepo{}
	productInfo := &fakeProductInfoProvider{info: &application.ProductPricingInfo{
		VariantID: variantID,
		ProductID: productID,
		BaseCost:  100,
		Currency:  "EUR",
		BrandID:   brandID,
	}}

	service := newPricingService(&fakePricingRuleRepo{}, &fakeClientPricingRepo{}, &fakeBrandMarginRepo{margin: margin}, &fakeDiscountRuleRepo{}, calcRepo, productInfo)
	resp, err := service.CalculatePrice(context.Background(), application.CalculatePriceCommand{
		ProductVariantID: variantID,
		ClientID:         clientID,
		Quantity:         1,
	})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if math.Abs(resp.FinalPrice-120) > 0.0001 {
		t.Fatalf("expected 10%% (10) + 10 EUR = 120, got %.2f", resp.FinalPrice)
	}
}

func TestPricingService_ListPricingRules_RepositoryError(t *testing.T) {
	repo := &fakePricingRuleRepo{listErr: domain.NewValidationError("db error")}
	service := newPricingService(repo, &fakeClientPricingRepo{}, &fakeBrandMarginRepo{}, &fakeDiscountRuleRepo{}, &fakeCalculationRepo{}, &fakeProductInfoProvider{})

	_, err := service.ListPricingRules(context.Background(), application.ListPricingRulesQuery{})
	if err == nil {
		t.Fatalf("expected repository error")
	}
}

func TestPricingService_CalculatePrice_ClientPricingRepoError(t *testing.T) {
	variantID := uuid.New()
	clientID := uuid.New()
	brandID := uuid.New()
	productID := uuid.New()

	clientRepo := &fakeClientPricingRepo{overrideErr: domain.NewValidationError("db error")}
	productInfo := &fakeProductInfoProvider{info: &application.ProductPricingInfo{
		VariantID: variantID,
		ProductID: productID,
		BaseCost:  100,
		Currency:  "EUR",
		BrandID:   brandID,
	}}

	service := newPricingService(&fakePricingRuleRepo{}, clientRepo, &fakeBrandMarginRepo{}, &fakeDiscountRuleRepo{}, &fakeCalculationRepo{}, productInfo)
	_, err := service.CalculatePrice(context.Background(), application.CalculatePriceCommand{
		ProductVariantID: variantID,
		ClientID:         clientID,
		Quantity:         1,
	})
	if err == nil {
		t.Fatalf("expected repository error")
	}
}

func TestPricingService_CalculatePrice_SaveCalculationError(t *testing.T) {
	variantID := uuid.New()
	clientID := uuid.New()
	brandID := uuid.New()
	productID := uuid.New()

	overridePrice, _ := domain.NewMoney(90, "EUR")
	override, _ := domain.NewClientPricing(clientID, variantID, overridePrice, time.Now(), nil)

	calcRepo := &fakeCalculationRepo{saveErr: domain.NewValidationError("save failed")}
	productInfo := &fakeProductInfoProvider{info: &application.ProductPricingInfo{
		VariantID: variantID,
		ProductID: productID,
		BaseCost:  50,
		Currency:  "EUR",
		BrandID:   brandID,
	}}

	service := newPricingService(&fakePricingRuleRepo{}, &fakeClientPricingRepo{override: override}, &fakeBrandMarginRepo{}, &fakeDiscountRuleRepo{}, calcRepo, productInfo)
	_, err := service.CalculatePrice(context.Background(), application.CalculatePriceCommand{
		ProductVariantID: variantID,
		ClientID:         clientID,
		Quantity:         1,
	})
	if err == nil {
		t.Fatalf("expected calculation save error")
	}
}

func TestPricingService_CreatePricingRule_InvalidName(t *testing.T) {
	service := newPricingService(&fakePricingRuleRepo{}, &fakeClientPricingRepo{}, &fakeBrandMarginRepo{}, &fakeDiscountRuleRepo{}, &fakeCalculationRepo{}, &fakeProductInfoProvider{})

	_, err := service.CreatePricingRule(context.Background(), application.CreatePricingRuleCommand{
		Name:             "",
		MarkupPercentage: 0.2,
		MinQuantity:      1,
	})
	if err == nil {
		t.Fatalf("expected validation error for empty name")
	}
}

func TestPricingService_CreateClientPricingInvalidPrice(t *testing.T) {
	service := newPricingService(&fakePricingRuleRepo{}, &fakeClientPricingRepo{}, &fakeBrandMarginRepo{}, &fakeDiscountRuleRepo{}, &fakeCalculationRepo{}, &fakeProductInfoProvider{})

	_, err := service.CreateClientPricing(context.Background(), application.CreateClientPricingCommand{
		ClientID:         uuid.New(),
		ProductVariantID: uuid.New(),
		FixedPrice:       -100,
		Currency:         "EUR",
	})
	if err == nil {
		t.Fatalf("expected validation error for negative price")
	}
}

func TestPricingService_CalculatePrice_WithClientOverride(t *testing.T) {
	variantID := uuid.New()
	clientID := uuid.New()
	brandID := uuid.New()
	productID := uuid.New()

	overridePrice, _ := domain.NewMoney(200, "EUR")
	override, _ := domain.NewClientPricing(clientID, variantID, overridePrice, time.Now(), nil)

	productInfo := &fakeProductInfoProvider{info: &application.ProductPricingInfo{
		VariantID: variantID,
		ProductID: productID,
		BaseCost:  100,
		Currency:  "EUR",
		BrandID:   brandID,
	}}

	service := newPricingService(&fakePricingRuleRepo{}, &fakeClientPricingRepo{override: override}, &fakeBrandMarginRepo{}, &fakeDiscountRuleRepo{}, &fakeCalculationRepo{}, productInfo)

	resp, err := service.CalculatePrice(context.Background(), application.CalculatePriceCommand{
		ProductVariantID: variantID,
		ClientID:         clientID,
		Quantity:         1,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if math.Abs(resp.FinalPrice-200) > 0.0001 {
		t.Fatalf("expected override price 200, got %.2f", resp.FinalPrice)
	}
}

func TestNewMoneyDTO(t *testing.T) {
	money, _ := domain.NewMoney(150.50, "EUR")
	dto := application.NewMoneyDTO(money)

	if dto.Amount != 150.50 {
		t.Fatalf("expected amount 150.50, got %.2f", dto.Amount)
	}
	if dto.Currency != "EUR" {
		t.Fatalf("expected currency EUR, got %s", dto.Currency)
	}
}

func TestNewPricingRuleDTO(t *testing.T) {
	percentage, _ := domain.NewPercentage(0.25)
	rule, _ := domain.NewPricingRule("Test Rule", nil, nil, percentage, 10, nil, time.Now(), nil)

	dto := application.NewPricingRuleDTO(rule)

	if dto.Name != "Test Rule" {
		t.Fatalf("expected name 'Test Rule', got %s", dto.Name)
	}
	if math.Abs(dto.Markup.Value-0.25) > 0.0001 {
		t.Fatalf("expected markup 0.25, got %.2f", dto.Markup.Value)
	}
	if dto.MinQuantity != 10 {
		t.Fatalf("expected minQuantity 10, got %d", dto.MinQuantity)
	}
}

func TestNewClientPricingDTO(t *testing.T) {
	clientID := uuid.New()
	variantID := uuid.New()
	price, _ := domain.NewMoney(99.99, "EUR")
	override, _ := domain.NewClientPricing(clientID, variantID, price, time.Now(), nil)

	dto := application.NewClientPricingDTO(override)

	if dto.ClientID != clientID {
		t.Fatalf("expected clientID %s, got %s", clientID, dto.ClientID)
	}
	if dto.ProductVariantID != variantID {
		t.Fatalf("expected variantID %s, got %s", variantID, dto.ProductVariantID)
	}
	if math.Abs(dto.FixedPrice.Amount-99.99) > 0.0001 {
		t.Fatalf("expected price 99.99, got %.2f", dto.FixedPrice.Amount)
	}
}

func TestNewPriceCalculationDTO(t *testing.T) {
	variantID := uuid.New()
	clientID := uuid.New()
	baseCost, _ := domain.NewMoney(100, "EUR")
	finalPrice, _ := domain.NewMoney(150, "EUR")

	calc, _ := domain.NewPriceCalculation(variantID, clientID, 5, baseCost, finalPrice, []string{"rule1"})

	dto := application.NewPriceCalculationDTO(calc)

	if dto.ProductVariantID != variantID {
		t.Fatalf("expected variantID %s, got %s", variantID, dto.ProductVariantID)
	}
	if dto.Quantity != 5 {
		t.Fatalf("expected quantity 5, got %d", dto.Quantity)
	}
	if math.Abs(dto.BaseCost.Amount-100) > 0.0001 {
		t.Fatalf("expected baseCost 100, got %.2f", dto.BaseCost.Amount)
	}
	if math.Abs(dto.FinalPrice.Amount-150) > 0.0001 {
		t.Fatalf("expected finalPrice 150, got %.2f", dto.FinalPrice.Amount)
	}
}
