package application_test

import (
	"context"
	"fmt"
	"math"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/joran-cortez/tramatex/internal/pricing/application"
	"github.com/joran-cortez/tramatex/internal/pricing/domain"
)

type fakeBaseRuleRepo struct {
	saved   []*domain.BaseSalesPriceRule
	rules   []*domain.BaseSalesPriceRule
	find    map[uuid.UUID]*domain.BaseSalesPriceRule
	saveErr error
	listErr error
}

func (f *fakeBaseRuleRepo) Save(ctx context.Context, rule *domain.BaseSalesPriceRule) error {
	if f.saveErr != nil {
		return f.saveErr
	}
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
	if f.listErr != nil {
		return nil, f.listErr
	}
	return f.rules, nil
}

type fakeSaleRuleRepo struct {
	rules   []*domain.SaleModificationRule
	rule    *domain.SaleModificationRule
	saveErr error
	listErr error
}

func (f *fakeSaleRuleRepo) Save(ctx context.Context, rule *domain.SaleModificationRule) error {
	if f.saveErr != nil {
		return f.saveErr
	}
	f.rules = append(f.rules, rule)
	return nil
}

func (f *fakeSaleRuleRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.SaleModificationRule, error) {
	return f.rule, nil
}

func (f *fakeSaleRuleRepo) ListActive(ctx context.Context, at time.Time) ([]*domain.SaleModificationRule, error) {
	if f.listErr != nil {
		return nil, f.listErr
	}
	return f.rules, nil
}

func (f *fakeSaleRuleRepo) ListApplicable(ctx context.Context, clientID string, productGroupID *uuid.UUID, orderTotal domain.Money, at time.Time) ([]*domain.SaleModificationRule, error) {
	if f.listErr != nil {
		return nil, f.listErr
	}
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
	info     *application.ProductPricingInfo
	infos    []*application.ProductPricingInfo
	err      error
	infoFunc func(context.Context, uuid.UUID) (*application.ProductPricingInfo, error)
}

func (f *fakeProductProvider) GetVariantPricingInfo(ctx context.Context, variantID uuid.UUID) (*application.ProductPricingInfo, error) {
	if f.err != nil {
		return nil, f.err
	}
	if f.infoFunc != nil {
		return f.infoFunc(ctx, variantID)
	}
	return f.info, nil
}

func (f *fakeProductProvider) GetVariantsPricingInfo(ctx context.Context, variantIDs []uuid.UUID) ([]*application.ProductPricingInfo, error) {
	if f.err != nil {
		return nil, f.err
	}
	if f.infos != nil {
		return f.infos, nil
	}
	if f.info != nil {
		// Fallback for tests that only provide a single info
		return []*application.ProductPricingInfo{f.info}, nil
	}
	return nil, nil
}

func (f *fakeProductProvider) ListVariantsPricingInfo(ctx context.Context, productID uuid.UUID) ([]*application.ProductPricingInfo, error) {
	return f.infos, nil
}

type fakeClientPricingRepo struct {
	override    *domain.ClientPricing
	overrideErr error
	saved       *domain.ClientPricing
}

func (f *fakeClientPricingRepo) Save(_ context.Context, o *domain.ClientPricing) error {
	f.saved = o
	return f.overrideErr
}

func (f *fakeClientPricingRepo) FindApplicable(_ context.Context, _ uuid.UUID, _ uuid.UUID, _ time.Time) (*domain.ClientPricing, error) {
	return f.override, f.overrideErr
}

func (f *fakeClientPricingRepo) FindApplicableBulk(_ context.Context, _ uuid.UUID, _ []uuid.UUID, _ time.Time) (map[uuid.UUID]*domain.ClientPricing, error) {
	if f.overrideErr != nil {
		return nil, f.overrideErr
	}
	res := make(map[uuid.UUID]*domain.ClientPricing)
	if f.override != nil {
		res[f.override.ProductVariantID] = f.override
	}
	return res, nil
}

type fakeCalculationRepo struct {
	mu      sync.Mutex
	calcs   []*domain.PriceCalculation
	saved   []*domain.PriceCalculation
	saveErr error
}

func (f *fakeCalculationRepo) Save(_ context.Context, calc *domain.PriceCalculation) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.saveErr != nil {
		return f.saveErr
	}
	f.saved = append(f.saved, calc)
	return nil
}

func (f *fakeCalculationRepo) GetSaved() []*domain.PriceCalculation {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.saved
}

func (f *fakeCalculationRepo) ListByProductVariantID(_ context.Context, _ uuid.UUID) ([]*domain.PriceCalculation, error) {
	return f.calcs, nil
}

func TestPricingEngineService_CreateBaseSalesPriceRule(t *testing.T) {
	repo := &fakeBaseRuleRepo{}
	service := application.NewPricingEngineService(repo, &fakeSaleRuleRepo{}, &fakeProductProvider{}, nil, nil, nil, nil)

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
	service := application.NewPricingEngineService(repo, &fakeSaleRuleRepo{}, &fakeProductProvider{}, nil, nil, nil, nil)

	_, err := service.UpdateBaseSalesPriceRule(context.Background(), application.UpdateBaseSalesPriceRuleCommand{ID: uuid.New()})
	if err == nil {
		t.Fatalf("expected not found error")
	}
}

func TestPricingEngineService_CalculateBaseSalesPrice_CacheSet(t *testing.T) {
	repo := &fakeBaseRuleRepo{}
	cache := &fakeBasePriceCache{}
	provider := &fakeProductProvider{}
	service := application.NewPricingEngineService(repo, &fakeSaleRuleRepo{}, provider, cache, nil, nil, nil)

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
		BaseCost:  decimal.NewFromInt(100),
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
	service := application.NewPricingEngineService(repo, saleRepo, provider, cache, nil, nil, nil)

	variantID := uuid.New()
	productID := uuid.New()
	clientID := uuid.New().String()
	groupID := uuid.New()

	basePrice, _ := domain.NewMoney(100, "EUR")
	cache.SetBasePrice(context.Background(), productID, variantID, basePrice)

	discount, _ := domain.NewMoney(10, "EUR")
	value, _ := domain.NewRuleValue(domain.RuleValueApplyFixedAmountDiscount, nil, &discount)
	rule, _ := domain.NewSaleModificationRule("Promo", []string{clientID}, &groupID, nil, value, 1, time.Now().Add(-time.Hour), nil)
	saleRepo.rules = []*domain.SaleModificationRule{rule}

	provider.info = &application.ProductPricingInfo{
		VariantID: variantID,
		ProductID: productID,
		BaseCost:  decimal.NewFromInt(100),
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
	service := application.NewPricingEngineService(&fakeBaseRuleRepo{}, &fakeSaleRuleRepo{}, &fakeProductProvider{}, nil, nil, nil, nil)

	_, err := service.CalculateFinalSalePrice(context.Background(), application.CalculateFinalSalePriceRequest{})
	if err == nil {
		t.Fatalf("expected validation error")
	}
}

func TestPricingEngineService_UpdateBaseSalesPriceRule_Success(t *testing.T) {
	ruleID := uuid.New()
	brandID := uuid.New()
	percentage, _ := domain.NewPercentage(0.2)
	value, _ := domain.NewRuleValue(domain.RuleValuePercentageMarkup, &percentage, nil)
	existingRule, _ := domain.NewBaseSalesPriceRule("OldRule", &brandID, nil, nil, nil, value)
	existingRule.ID = ruleID

	repo := &fakeBaseRuleRepo{find: map[uuid.UUID]*domain.BaseSalesPriceRule{ruleID: existingRule}}
	service := application.NewPricingEngineService(repo, &fakeSaleRuleRepo{}, &fakeProductProvider{}, nil, nil, nil, nil)

	newPercentage := application.PercentageDTO{Value: 0.3}
	newValue := application.RuleValueDTO{Type: string(domain.RuleValuePercentageMarkup), PercentageValue: &newPercentage}
	newName := "UpdatedRule"
	dto, err := service.UpdateBaseSalesPriceRule(context.Background(), application.UpdateBaseSalesPriceRuleCommand{
		ID:    ruleID,
		Name:  &newName,
		Value: &newValue,
	})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if dto.Name != "UpdatedRule" {
		t.Fatalf("expected updated name")
	}
	if len(repo.saved) != 1 {
		t.Fatalf("expected rule to be saved")
	}
}

func TestPricingEngineService_CreateSaleModificationRule_PercentageDiscount(t *testing.T) {
	repo := &fakeSaleRuleRepo{}
	service := application.NewPricingEngineService(&fakeBaseRuleRepo{}, repo, &fakeProductProvider{}, nil, nil, nil, nil)

	clientID := uuid.New().String()
	groupID := uuid.New()
	percentage := application.PercentageDTO{Value: 0.15}
	value := application.RuleValueDTO{Type: string(domain.RuleValueApplyPercentageDiscount), PercentageValue: &percentage}

	_, err := service.CreateSaleModificationRule(context.Background(), application.CreateSaleModificationRuleCommand{
		Name:           "BulkDiscount",
		ClientIDs:      []string{clientID},
		ProductGroupID: &groupID,
		Value:          value,
		Priority:       10,
	})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if len(repo.rules) != 1 {
		t.Fatalf("expected rule saved")
	}
	if repo.rules[0].Name != "BulkDiscount" {
		t.Fatalf("unexpected rule name")
	}
}

func TestPricingEngineService_CreateSaleModificationRule_FixedDiscount(t *testing.T) {
	repo := &fakeSaleRuleRepo{}
	service := application.NewPricingEngineService(&fakeBaseRuleRepo{}, repo, &fakeProductProvider{}, nil, nil, nil, nil)

	money := application.MoneyDTO{Amount: 25, Currency: "EUR"}
	value := application.RuleValueDTO{Type: string(domain.RuleValueApplyFixedAmountDiscount), MoneyValue: &money}

	minOrder := application.MoneyDTO{Amount: 100, Currency: "EUR"}

	_, err := service.CreateSaleModificationRule(context.Background(), application.CreateSaleModificationRuleCommand{
		Name:                "VolumeDiscount",
		Value:               value,
		MinOrderTotalAmount: &minOrder,
		Priority:            5,
	})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if len(repo.rules) != 1 {
		t.Fatalf("expected rule saved")
	}
}

func TestPricingEngineService_CreateSaleModificationRule_WithEffectiveDates(t *testing.T) {
	repo := &fakeSaleRuleRepo{}
	service := application.NewPricingEngineService(&fakeBaseRuleRepo{}, repo, &fakeProductProvider{}, nil, nil, nil, nil)

	effectiveFrom := time.Now().Add(-24 * time.Hour)
	effectiveTo := time.Now().Add(7 * 24 * time.Hour)

	percentage := application.PercentageDTO{Value: 0.2}
	value := application.RuleValueDTO{Type: string(domain.RuleValueApplyPercentageDiscount), PercentageValue: &percentage}
	isActive := true

	_, err := service.CreateSaleModificationRule(context.Background(), application.CreateSaleModificationRuleCommand{
		Name:          "FlashSale",
		Value:         value,
		Priority:      100,
		EffectiveFrom: effectiveFrom,
		EffectiveTo:   &effectiveTo,
		IsActive:      &isActive,
	})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if len(repo.rules) != 1 {
		t.Fatalf("expected rule saved")
	}
	if repo.rules[0].EffectiveFrom != effectiveFrom {
		t.Fatalf("effectiveFrom not preserved")
	}
}

func TestPricingEngineService_UpdateSaleModificationRule_NotFound(t *testing.T) {
	service := application.NewPricingEngineService(&fakeBaseRuleRepo{}, &fakeSaleRuleRepo{}, &fakeProductProvider{}, nil, nil, nil, nil)

	percentage := application.PercentageDTO{Value: 0.1}
	value := application.RuleValueDTO{Type: string(domain.RuleValueApplyPercentageDiscount), PercentageValue: &percentage}
	name := "NotFound"

	_, err := service.UpdateSaleModificationRule(context.Background(), application.UpdateSaleModificationRuleCommand{
		ID:    uuid.New(),
		Name:  &name,
		Value: &value,
	})
	if err == nil {
		t.Fatalf("expected not found error")
	}
}

func TestPricingEngineService_CalculateBaseSalesPrice_NoRules(t *testing.T) {
	repo := &fakeBaseRuleRepo{rules: []*domain.BaseSalesPriceRule{}}
	provider := &fakeProductProvider{}
	service := application.NewPricingEngineService(repo, &fakeSaleRuleRepo{}, provider, nil, nil, nil, nil)

	variantID := uuid.New()
	provider.info = &application.ProductPricingInfo{
		VariantID: variantID,
		ProductID: uuid.New(),
		BaseCost:  decimal.NewFromInt(50),
		Currency:  "EUR",
		BrandID:   uuid.New(),
	}

	resp, err := service.CalculateBaseSalesPrice(context.Background(), application.CalculateBaseSalesPriceRequest{VariantID: variantID})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if resp.BaseSalesPrice.Amount != 50 {
		t.Fatalf("expected base cost, got %.2f", resp.BaseSalesPrice.Amount)
	}
}

func TestPricingEngineService_CalculateFinalSalePrice_CacheHit(t *testing.T) {
	repo := &fakeBaseRuleRepo{}
	saleRepo := &fakeSaleRuleRepo{rules: []*domain.SaleModificationRule{}}
	cache := &fakeBasePriceCache{}
	provider := &fakeProductProvider{}
	service := application.NewPricingEngineService(repo, saleRepo, provider, cache, nil, nil, nil)

	variantID := uuid.New()
	productID := uuid.New()
	cachedPrice, _ := domain.NewMoney(120, "EUR")
	cache.SetBasePrice(context.Background(), productID, variantID, cachedPrice)

	provider.info = &application.ProductPricingInfo{
		VariantID: variantID,
		ProductID: productID,
		BaseCost:  decimal.NewFromInt(100),
		Currency:  "EUR",
		BrandID:   uuid.New(),
	}

	resp, err := service.CalculateFinalSalePrice(context.Background(), application.CalculateFinalSalePriceRequest{
		ClientID: uuid.New().String(),
		SaleItems: []application.SaleItemRequest{{
			ProductVariantID: variantID,
			Quantity:         1,
		}},
	})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if resp.SaleTotal.Amount != 120 {
		t.Fatalf("expected cached price, got %.2f", resp.SaleTotal.Amount)
	}
}

func TestPricingEngineService_CalculateFinalSalePrice_NoRules(t *testing.T) {
	repo := &fakeBaseRuleRepo{}
	saleRepo := &fakeSaleRuleRepo{rules: []*domain.SaleModificationRule{}}
	cache := &fakeBasePriceCache{}
	provider := &fakeProductProvider{}
	service := application.NewPricingEngineService(repo, saleRepo, provider, cache, nil, nil, nil)

	variantID := uuid.New()
	productID := uuid.New()

	basePrice, _ := domain.NewMoney(100, "EUR")
	cache.SetBasePrice(context.Background(), productID, variantID, basePrice)

	provider.info = &application.ProductPricingInfo{
		VariantID: variantID,
		ProductID: productID,
		BaseCost:  decimal.NewFromInt(100),
		Currency:  "EUR",
		BrandID:   uuid.New(),
	}

	resp, err := service.CalculateFinalSalePrice(context.Background(), application.CalculateFinalSalePriceRequest{
		ClientID: uuid.New().String(),
		SaleItems: []application.SaleItemRequest{{
			ProductVariantID: variantID,
			Quantity:         1,
		}},
	})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if resp.SaleTotal.Amount != 100 {
		t.Fatalf("expected base price, got %.2f", resp.SaleTotal.Amount)
	}
}

func TestPricingEngineService_CreateBaseSalesPriceRule_RepositoryError(t *testing.T) {
	baseRepo := &fakeBaseRuleRepo{saveErr: domain.NewValidationError("db error")}
	service := application.NewPricingEngineService(baseRepo, &fakeSaleRuleRepo{}, nil, nil, nil, nil, nil)

	value := application.RuleValueDTO{Type: string(domain.RuleValuePercentageMarkup), PercentageValue: &application.PercentageDTO{Value: 0.2}}
	_, err := service.CreateBaseSalesPriceRule(context.Background(), application.CreateBaseSalesPriceRuleCommand{
		Name:  "Test Rule",
		Value: value,
	})
	if err == nil {
		t.Fatalf("expected repository error")
	}
}

func TestPricingEngineService_CreateSaleModificationRule_InvalidValue(t *testing.T) {
	service := application.NewPricingEngineService(&fakeBaseRuleRepo{}, &fakeSaleRuleRepo{}, nil, nil, nil, nil, nil)

	discountPercent := -0.5
	clientID := uuid.New().String()
	_, err := service.CreateSaleModificationRule(context.Background(), application.CreateSaleModificationRuleCommand{
		Name:          "Invalid Rule",
		Value:         application.RuleValueDTO{Type: string(domain.RuleValueApplyPercentageDiscount), PercentageValue: &application.PercentageDTO{Value: discountPercent}},
		ClientIDs:     []string{clientID},
		EffectiveFrom: time.Now(),
	})
	if err == nil {
		t.Fatalf("expected validation error for negative discount")
	}
}

func TestPricingEngineService_CalculateFinalSalePrice_MultipleItems(t *testing.T) {
	baseRepo := &fakeBaseRuleRepo{}
	modRepo := &fakeSaleRuleRepo{}
	provider := &fakeProductProvider{}

	service := application.NewPricingEngineService(baseRepo, modRepo, provider, nil, nil, nil, nil)

	variantID1 := uuid.New()
	variantID2 := uuid.New()
	productID1 := uuid.New()
	productID2 := uuid.New()

	provider.infos = []*application.ProductPricingInfo{
		{
			VariantID: variantID1,
			ProductID: productID1,
			BaseCost:  decimal.NewFromInt(50),
			Currency:  "EUR",
			BrandID:   uuid.New(),
		},
		{
			VariantID: variantID2,
			ProductID: productID2,
			BaseCost:  decimal.NewFromInt(100),
			Currency:  "EUR",
			BrandID:   uuid.New(),
		},
	}

	resp, err := service.CalculateFinalSalePrice(context.Background(), application.CalculateFinalSalePriceRequest{
		ClientID: uuid.New().String(),
		SaleItems: []application.SaleItemRequest{
			{ProductVariantID: variantID1, Quantity: 2},
			{ProductVariantID: variantID2, Quantity: 1},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedTotal := (50.0 * 2) + 100.0
	if resp.SaleTotal.Amount != expectedTotal {
		t.Fatalf("expected total %.2f, got %.2f", expectedTotal, resp.SaleTotal.Amount)
	}
	if len(resp.CalculatedItems) != 2 {
		t.Fatalf("expected 2 items, got %d", len(resp.CalculatedItems))
	}
}

func TestPricingEngineService_CalculateFinalSalePrice_ProductInfoError(t *testing.T) {
	provider := &fakeProductProvider{err: domain.NewNotFoundError("product not found")}
	service := application.NewPricingEngineService(&fakeBaseRuleRepo{}, &fakeSaleRuleRepo{}, provider, nil, nil, nil, nil)

	_, err := service.CalculateFinalSalePrice(context.Background(), application.CalculateFinalSalePriceRequest{
		ClientID: uuid.New().String(),
		SaleItems: []application.SaleItemRequest{{
			ProductVariantID: uuid.New(),
			Quantity:         1,
		}},
	})
	if err == nil {
		t.Fatalf("expected product info error")
	}
}

func TestNewRuleValueDTO(t *testing.T) {
	percentage, _ := domain.NewPercentage(0.3)
	ruleValue, _ := domain.NewRuleValue(domain.RuleValuePercentageMarkup, &percentage, nil)

	dto := application.NewRuleValueDTO(ruleValue)

	if dto.Type != string(domain.RuleValuePercentageMarkup) {
		t.Fatalf("expected type %s, got %s", domain.RuleValuePercentageMarkup, dto.Type)
	}
	if dto.PercentageValue == nil {
		t.Fatalf("expected percentage value to be set")
	}
	if math.Abs(dto.PercentageValue.Value-0.3) > 0.0001 {
		t.Fatalf("expected percentage 0.3, got %.2f", dto.PercentageValue.Value)
	}
	if dto.MoneyValue != nil {
		t.Fatalf("expected money value to be nil")
	}
}

func TestNewRuleValueDTO_WithMoney(t *testing.T) {
	money, _ := domain.NewMoney(50, "EUR")
	ruleValue, _ := domain.NewRuleValue(domain.RuleValueFixedAmountIncrease, nil, &money)

	dto := application.NewRuleValueDTO(ruleValue)

	if dto.Type != string(domain.RuleValueFixedAmountIncrease) {
		t.Fatalf("expected type %s, got %s", domain.RuleValueFixedAmountIncrease, dto.Type)
	}
	if dto.MoneyValue == nil {
		t.Fatalf("expected money value to be set")
	}
	if math.Abs(dto.MoneyValue.Amount-50) > 0.0001 {
		t.Fatalf("expected amount 50, got %.2f", dto.MoneyValue.Amount)
	}
	if dto.PercentageValue != nil {
		t.Fatalf("expected percentage value to be nil")
	}
}

func TestNewBaseSalesPriceRuleDTO(t *testing.T) {
	brandID := uuid.New()
	productID := uuid.New()
	percentage, _ := domain.NewPercentage(0.4)
	ruleValue, _ := domain.NewRuleValue(domain.RuleValuePercentageMarkup, &percentage, nil)

	rule, _ := domain.NewBaseSalesPriceRule("Base Rule", &brandID, nil, &productID, nil, ruleValue)

	dto := application.NewBaseSalesPriceRuleDTO(rule)

	if dto.Name != "Base Rule" {
		t.Fatalf("expected name 'Base Rule', got %s", dto.Name)
	}
	if dto.BrandID == nil || *dto.BrandID != brandID {
		t.Fatalf("expected brandID %s", brandID)
	}
	if dto.ProductID == nil || *dto.ProductID != productID {
		t.Fatalf("expected productID %s", productID)
	}
	if !dto.IsActive {
		t.Fatalf("expected rule to be active")
	}
}

func TestNewSaleModificationRuleDTO(t *testing.T) {
	clientID := uuid.New().String()
	percentage, _ := domain.NewPercentage(0.15)
	ruleValue, _ := domain.NewRuleValue(domain.RuleValueApplyPercentageDiscount, &percentage, nil)

	rule, _ := domain.NewSaleModificationRule(
		"Discount Rule",
		[]string{clientID},
		nil,
		nil,
		ruleValue,
		1,
		time.Now(),
		nil,
	)

	dto := application.NewSaleModificationRuleDTO(rule)

	if dto.Name != "Discount Rule" {
		t.Fatalf("expected name 'Discount Rule', got %s", dto.Name)
	}
	if len(dto.ClientIDs) != 1 || dto.ClientIDs[0] != clientID {
		t.Fatalf("expected clientID %s in list", clientID)
	}
	if dto.Priority != 1 {
		t.Fatalf("expected priority 1, got %d", dto.Priority)
	}
	if !dto.IsActive {
		t.Fatalf("expected rule to be active")
	}
}

func TestPricingEngineService_UpdateSaleModificationRule_AllFields(t *testing.T) {
	ruleID := uuid.New()
	clientID1 := uuid.New().String()
	clientID2 := uuid.New().String()
	productGroupID := uuid.New()

	percentage, _ := domain.NewPercentage(0.1)
	ruleValue, _ := domain.NewRuleValue(domain.RuleValueApplyPercentageDiscount, &percentage, nil)

	existingRule, _ := domain.NewSaleModificationRule(
		"Old Name",
		[]string{clientID1},
		nil,
		nil,
		ruleValue,
		1,
		time.Now(),
		nil,
	)

	modRepo := &fakeSaleRuleRepo{rule: existingRule}
	service := application.NewPricingEngineService(&fakeBaseRuleRepo{}, modRepo, nil, nil, nil, nil, nil)

	newName := "Updated Name"
	newPercentage := &application.PercentageDTO{Value: 0.2}
	newValue := &application.RuleValueDTO{Type: string(domain.RuleValueApplyPercentageDiscount), PercentageValue: newPercentage}
	newPriority := 5
	newActive := false
	newFrom := time.Now().Add(24 * time.Hour)
	newTo := time.Now().Add(48 * time.Hour)
	minOrderTotal := &application.MoneyDTO{Amount: 100, Currency: "EUR"}

	resp, err := service.UpdateSaleModificationRule(context.Background(), application.UpdateSaleModificationRuleCommand{
		ID:                  ruleID,
		Name:                &newName,
		ClientIDs:           []string{clientID2},
		ProductGroupID:      &productGroupID,
		MinOrderTotalAmount: minOrderTotal,
		Value:               newValue,
		Priority:            &newPriority,
		IsActive:            &newActive,
		EffectiveFrom:       &newFrom,
		EffectiveTo:         &newTo,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Name != newName {
		t.Fatalf("expected updated name")
	}
	if resp.Priority != newPriority {
		t.Fatalf("expected updated priority")
	}
	if resp.IsActive {
		t.Fatalf("expected rule to be inactive")
	}
}

func TestPricingEngineService_UpdateSaleModificationRule_PartialUpdate(t *testing.T) {
	ruleID := uuid.New()
	clientID := uuid.New().String()

	percentage, _ := domain.NewPercentage(0.15)
	ruleValue, _ := domain.NewRuleValue(domain.RuleValueApplyPercentageDiscount, &percentage, nil)

	existingRule, _ := domain.NewSaleModificationRule(
		"Original",
		[]string{clientID},
		nil,
		nil,
		ruleValue,
		2,
		time.Now(),
		nil,
	)

	modRepo := &fakeSaleRuleRepo{rule: existingRule}
	service := application.NewPricingEngineService(&fakeBaseRuleRepo{}, modRepo, nil, nil, nil, nil, nil)

	newName := "Partially Updated"

	resp, err := service.UpdateSaleModificationRule(context.Background(), application.UpdateSaleModificationRuleCommand{
		ID:   ruleID,
		Name: &newName,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Name != newName {
		t.Fatalf("expected name updated")
	}
	if resp.Priority != 2 {
		t.Fatalf("expected original priority preserved")
	}
}

func TestPricingEngineService_UpdateBaseSalesPriceRule_MultipleFields(t *testing.T) {
	ruleID := uuid.New()
	brandID := uuid.New()
	newBrandID := uuid.New()

	percentage, _ := domain.NewPercentage(0.2)
	ruleValue, _ := domain.NewRuleValue(domain.RuleValuePercentageMarkup, &percentage, nil)

	existingRule, _ := domain.NewBaseSalesPriceRule("Original Rule", &brandID, nil, nil, nil, ruleValue)

	baseRepo := &fakeBaseRuleRepo{find: map[uuid.UUID]*domain.BaseSalesPriceRule{ruleID: existingRule}}
	service := application.NewPricingEngineService(baseRepo, &fakeSaleRuleRepo{}, nil, nil, nil, nil, nil)

	newName := "Updated Base Rule"
	newPercentage := &application.PercentageDTO{Value: 0.3}
	newValue := &application.RuleValueDTO{Type: string(domain.RuleValuePercentageMarkup), PercentageValue: newPercentage}
	isActive := false

	resp, err := service.UpdateBaseSalesPriceRule(context.Background(), application.UpdateBaseSalesPriceRuleCommand{
		ID:       ruleID,
		Name:     &newName,
		BrandID:  &newBrandID,
		Value:    newValue,
		IsActive: &isActive,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Name != newName {
		t.Fatalf("expected updated name")
	}
	if resp.IsActive {
		t.Fatalf("expected rule to be inactive")
	}
}

func TestPricingEngineService_CreateBaseSalesPriceRule_WithProductGroup(t *testing.T) {
	repo := &fakeBaseRuleRepo{}
	service := application.NewPricingEngineService(repo, &fakeSaleRuleRepo{}, nil, nil, nil, nil, nil)

	groupID := uuid.New()
	value := application.RuleValueDTO{Type: string(domain.RuleValuePercentageMarkup), PercentageValue: &application.PercentageDTO{Value: 0.25}}

	_, err := service.CreateBaseSalesPriceRule(context.Background(), application.CreateBaseSalesPriceRuleCommand{
		Name:           "Group Rule",
		ProductGroupID: &groupID,
		Value:          value,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repo.saved) != 1 {
		t.Fatalf("expected rule saved")
	}
}

func TestPricingEngineService_CalculateFinalSalePrice_WithCache(t *testing.T) {
	cache := &fakeBasePriceCache{}
	provider := &fakeProductProvider{}
	service := application.NewPricingEngineService(&fakeBaseRuleRepo{}, &fakeSaleRuleRepo{}, provider, cache, nil, nil, nil)

	variantID := uuid.New()
	productID := uuid.New()
	brandID := uuid.New()
	clientID := uuid.New().String()

	cachedPrice, _ := domain.NewMoney(120, "EUR")
	cache.prices = map[string]domain.Money{
		cache.key(productID, variantID): cachedPrice,
	}

	provider.info = &application.ProductPricingInfo{
		VariantID: variantID,
		ProductID: productID,
		BaseCost:  decimal.NewFromInt(100),
		Currency:  "EUR",
		BrandID:   brandID,
	}

	resp, err := service.CalculateFinalSalePrice(context.Background(), application.CalculateFinalSalePriceRequest{
		ClientID: clientID,
		SaleItems: []application.SaleItemRequest{{
			ProductVariantID: variantID,
			Quantity:         2,
		}},
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.CalculatedItems) != 1 {
		t.Fatalf("expected 1 item")
	}
	if resp.CalculatedItems[0].BaseSalesPrice.Amount != 120 {
		t.Fatalf("expected cached price 120, got %.2f", resp.CalculatedItems[0].BaseSalesPrice.Amount)
	}
}

func TestPricingEngineService_UpdateSaleModificationRule_InvalidPercentage(t *testing.T) {
	ruleID := uuid.New()
	clientID := uuid.New().String()

	percentage, _ := domain.NewPercentage(0.1)
	ruleValue, _ := domain.NewRuleValue(domain.RuleValueApplyPercentageDiscount, &percentage, nil)
	existingRule, _ := domain.NewSaleModificationRule("Rule", []string{clientID}, nil, nil, ruleValue, 1, time.Now(), nil)

	modRepo := &fakeSaleRuleRepo{rule: existingRule}
	service := application.NewPricingEngineService(&fakeBaseRuleRepo{}, modRepo, nil, nil, nil, nil, nil)

	invalidValue := &application.RuleValueDTO{
		Type:            string(domain.RuleValueApplyPercentageDiscount),
		PercentageValue: &application.PercentageDTO{Value: -0.5},
	}

	_, err := service.UpdateSaleModificationRule(context.Background(), application.UpdateSaleModificationRuleCommand{
		ID:    ruleID,
		Value: invalidValue,
	})

	if err == nil {
		t.Fatalf("expected error for invalid percentage")
	}
}

func TestPricingEngineService_UpdateSaleModificationRule_InvalidMoney(t *testing.T) {
	ruleID := uuid.New()
	clientID := uuid.New().String()

	percentage, _ := domain.NewPercentage(0.1)
	ruleValue, _ := domain.NewRuleValue(domain.RuleValueApplyPercentageDiscount, &percentage, nil)
	existingRule, _ := domain.NewSaleModificationRule("Rule", []string{clientID}, nil, nil, ruleValue, 1, time.Now(), nil)

	modRepo := &fakeSaleRuleRepo{rule: existingRule}
	service := application.NewPricingEngineService(&fakeBaseRuleRepo{}, modRepo, nil, nil, nil, nil, nil)

	invalidMinOrder := &application.MoneyDTO{Amount: -100, Currency: "EUR"}

	_, err := service.UpdateSaleModificationRule(context.Background(), application.UpdateSaleModificationRuleCommand{
		ID:                  ruleID,
		MinOrderTotalAmount: invalidMinOrder,
	})

	if err == nil {
		t.Fatalf("expected error for negative money amount")
	}
}

func TestPricingEngineService_UpdateBaseSalesPriceRule_InvalidValue(t *testing.T) {
	ruleID := uuid.New()
	brandID := uuid.New()

	percentage, _ := domain.NewPercentage(0.2)
	ruleValue, _ := domain.NewRuleValue(domain.RuleValuePercentageMarkup, &percentage, nil)
	existingRule, _ := domain.NewBaseSalesPriceRule("Rule", &brandID, nil, nil, nil, ruleValue)

	baseRepo := &fakeBaseRuleRepo{find: map[uuid.UUID]*domain.BaseSalesPriceRule{ruleID: existingRule}}
	service := application.NewPricingEngineService(baseRepo, &fakeSaleRuleRepo{}, nil, nil, nil, nil, nil)

	invalidValue := &application.RuleValueDTO{
		Type:            string(domain.RuleValuePercentageMarkup),
		PercentageValue: &application.PercentageDTO{Value: 2.5},
	}

	_, err := service.UpdateBaseSalesPriceRule(context.Background(), application.UpdateBaseSalesPriceRuleCommand{
		ID:    ruleID,
		Value: invalidValue,
	})

	if err == nil {
		t.Fatalf("expected error for invalid percentage > 1")
	}
}

func TestPricingEngineService_SelectBaseRule_WithNonMatchingProductGroup(t *testing.T) {
	groupID1 := uuid.New()
	groupID2 := uuid.New()
	productID := uuid.New()
	variantID := uuid.New()
	brandID := uuid.New()

	percentage, _ := domain.NewPercentage(0.25)
	ruleValue, _ := domain.NewRuleValue(domain.RuleValuePercentageMarkup, &percentage, nil)

	ruleWithUnmatchedGroup, _ := domain.NewBaseSalesPriceRule("Group Rule", nil, &groupID2, nil, nil, ruleValue)
	ruleBrand, _ := domain.NewBaseSalesPriceRule("Brand Rule", &brandID, nil, nil, nil, ruleValue)

	baseRepo := &fakeBaseRuleRepo{rules: []*domain.BaseSalesPriceRule{ruleWithUnmatchedGroup, ruleBrand}}
	provider := &fakeProductProvider{
		info: &application.ProductPricingInfo{
			VariantID: variantID,
			ProductID: productID,
			BaseCost:  decimal.NewFromInt(100),
			Currency:  "EUR",
			BrandID:   brandID,
			GroupIDs:  []uuid.UUID{groupID1},
		},
	}
	service := application.NewPricingEngineService(baseRepo, &fakeSaleRuleRepo{}, provider, nil, nil, nil, nil)

	resp, err := service.CalculateBaseSalesPrice(context.Background(), application.CalculateBaseSalesPriceRequest{
		VariantID: variantID,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.BaseSalesPrice.Amount != 125 {
		t.Fatalf("expected price 125 from brand rule, got %.2f", resp.BaseSalesPrice.Amount)
	}
}

func TestPricingEngineService_SelectBaseRule_WithMatchingProductGroup(t *testing.T) {
	groupID := uuid.New()
	productID := uuid.New()
	variantID := uuid.New()
	brandID := uuid.New()

	percentageBrand, _ := domain.NewPercentage(0.2)
	ruleValueBrand, _ := domain.NewRuleValue(domain.RuleValuePercentageMarkup, &percentageBrand, nil)
	ruleBrand, _ := domain.NewBaseSalesPriceRule("Brand Rule", &brandID, nil, nil, nil, ruleValueBrand)

	percentageGroup, _ := domain.NewPercentage(0.3)
	ruleValueGroup, _ := domain.NewRuleValue(domain.RuleValuePercentageMarkup, &percentageGroup, nil)
	ruleGroup, _ := domain.NewBaseSalesPriceRule("Group Rule", nil, &groupID, nil, nil, ruleValueGroup)

	baseRepo := &fakeBaseRuleRepo{rules: []*domain.BaseSalesPriceRule{ruleBrand, ruleGroup}}
	provider := &fakeProductProvider{
		info: &application.ProductPricingInfo{
			VariantID: variantID,
			ProductID: productID,
			BaseCost:  decimal.NewFromInt(100),
			Currency:  "EUR",
			BrandID:   brandID,
			GroupIDs:  []uuid.UUID{groupID},
		},
	}
	service := application.NewPricingEngineService(baseRepo, &fakeSaleRuleRepo{}, provider, nil, nil, nil, nil)

	resp, err := service.CalculateBaseSalesPrice(context.Background(), application.CalculateBaseSalesPriceRequest{
		VariantID: variantID,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.BaseSalesPrice.Amount != 130 {
		t.Fatalf("expected price 130 from group rule (higher priority), got %.2f", resp.BaseSalesPrice.Amount)
	}
}

func TestPricingEngineService_CreateBaseSalesPriceRule_SaveError(t *testing.T) {
	repo := &fakeBaseRuleRepo{saveErr: fmt.Errorf("save failed")}
	service := application.NewPricingEngineService(repo, &fakeSaleRuleRepo{}, &fakeProductProvider{}, nil, nil, nil, nil)

	value := application.RuleValueDTO{Type: string(domain.RuleValuePercentageMarkup), PercentageValue: &application.PercentageDTO{Value: 0.2}}
	_, err := service.CreateBaseSalesPriceRule(context.Background(), application.CreateBaseSalesPriceRuleCommand{Name: "Rule", Value: value})

	if err == nil {
		t.Fatalf("expected save error")
	}
}

func TestPricingEngineService_CreateBaseSalesPriceRule_NoIsActive(t *testing.T) {
	repo := &fakeBaseRuleRepo{}
	service := application.NewPricingEngineService(repo, &fakeSaleRuleRepo{}, &fakeProductProvider{}, nil, nil, nil, nil)

	value := application.RuleValueDTO{Type: string(domain.RuleValuePercentageMarkup), PercentageValue: &application.PercentageDTO{Value: 0.15}}
	resp, err := service.CreateBaseSalesPriceRule(context.Background(), application.CreateBaseSalesPriceRuleCommand{
		Name:  "Default Active Rule",
		Value: value,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.IsActive {
		t.Fatalf("expected rule to be active by default")
	}
}

func TestPricingEngineService_CalculateBaseSalesPrice_ProductInfoError(t *testing.T) {
	provider := &fakeProductProvider{err: fmt.Errorf("product lookup failed")}
	service := application.NewPricingEngineService(&fakeBaseRuleRepo{}, &fakeSaleRuleRepo{}, provider, nil, nil, nil, nil)

	_, err := service.CalculateBaseSalesPrice(context.Background(), application.CalculateBaseSalesPriceRequest{
		VariantID: uuid.New(),
	})

	if err == nil {
		t.Fatalf("expected product info error")
	}
}

func TestPricingEngineService_CreateSaleModificationRule_SaveError(t *testing.T) {
	repo := &fakeSaleRuleRepo{saveErr: fmt.Errorf("save failed")}
	service := application.NewPricingEngineService(&fakeBaseRuleRepo{}, repo, nil, nil, nil, nil, nil)

	clientID := uuid.New().String()
	value := application.RuleValueDTO{Type: string(domain.RuleValueApplyPercentageDiscount), PercentageValue: &application.PercentageDTO{Value: 0.1}}

	_, err := service.CreateSaleModificationRule(context.Background(), application.CreateSaleModificationRuleCommand{
		Name:      "Test Rule",
		ClientIDs: []string{clientID},
		Value:     value,
		Priority:  1,
	})

	if err == nil {
		t.Fatalf("expected save error")
	}
}

// --- Tests for Client Pricing Override (G1) and Audit Trail (G2) ---

func TestPricingEngineService_CalculateFinalSalePrice_ClientOverrideFound(t *testing.T) {
	cache := &fakeBasePriceCache{}
	provider := &fakeProductProvider{}

	variantID := uuid.New()
	productID := uuid.New()
	clientID := uuid.New()

	basePrice, _ := domain.NewMoney(100, "EUR")
	cache.SetBasePrice(context.Background(), productID, variantID, basePrice)

	provider.info = &application.ProductPricingInfo{
		VariantID: variantID,
		ProductID: productID,
		BaseCost:  decimal.NewFromInt(80),
		Currency:  "EUR",
		BrandID:   uuid.New(),
	}

	fixedPrice, _ := domain.NewMoney(75, "EUR")
	override, _ := domain.NewClientPricing(clientID, variantID, fixedPrice, time.Now().Add(-time.Hour), nil)

	clientPricingRepo := &fakeClientPricingRepo{override: override}
	calcRepo := &fakeCalculationRepo{}

	service := application.NewPricingEngineService(
		&fakeBaseRuleRepo{}, &fakeSaleRuleRepo{}, provider, cache, nil,
		clientPricingRepo, calcRepo,
	)

	resp, err := service.CalculateFinalSalePrice(context.Background(), application.CalculateFinalSalePriceRequest{
		ClientID: clientID.String(),
		SaleItems: []application.SaleItemRequest{{
			ProductVariantID: variantID,
			Quantity:         2,
		}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.CalculatedItems[0].FinalPrice.Amount != 75 {
		t.Fatalf("expected override price 75, got %.2f", resp.CalculatedItems[0].FinalPrice.Amount)
	}
	if resp.SaleTotal.Amount != 150 {
		t.Fatalf("expected total 150, got %.2f", resp.SaleTotal.Amount)
	}
}

func TestPricingEngineService_CalculateFinalSalePrice_ClientOverrideNotFound(t *testing.T) {
	cache := &fakeBasePriceCache{}
	provider := &fakeProductProvider{}

	variantID := uuid.New()
	productID := uuid.New()
	clientID := uuid.New()
	groupID := uuid.New()

	basePrice, _ := domain.NewMoney(100, "EUR")
	cache.SetBasePrice(context.Background(), productID, variantID, basePrice)

	provider.info = &application.ProductPricingInfo{
		VariantID: variantID,
		ProductID: productID,
		BaseCost:  decimal.NewFromInt(80),
		Currency:  "EUR",
		BrandID:   uuid.New(),
		GroupIDs:  []uuid.UUID{groupID},
	}

	// No override â†’ nil
	clientPricingRepo := &fakeClientPricingRepo{override: nil}

	discount, _ := domain.NewMoney(10, "EUR")
	value, _ := domain.NewRuleValue(domain.RuleValueApplyFixedAmountDiscount, nil, &discount)
	rule, _ := domain.NewSaleModificationRule("Promo", []string{clientID.String()}, &groupID, nil, value, 1, time.Now().Add(-time.Hour), nil)
	saleRepo := &fakeSaleRuleRepo{rules: []*domain.SaleModificationRule{rule}}

	service := application.NewPricingEngineService(
		&fakeBaseRuleRepo{}, saleRepo, provider, cache, nil,
		clientPricingRepo, nil,
	)

	resp, err := service.CalculateFinalSalePrice(context.Background(), application.CalculateFinalSalePriceRequest{
		ClientID: clientID.String(),
		SaleItems: []application.SaleItemRequest{{
			ProductVariantID: variantID,
			Quantity:         1,
		}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should fall through to SaleModificationRule: 100 - 10 = 90
	if resp.CalculatedItems[0].FinalPrice.Amount != 90 {
		t.Fatalf("expected rule-discounted price 90, got %.2f", resp.CalculatedItems[0].FinalPrice.Amount)
	}
}

func TestPricingEngineService_CalculateFinalSalePrice_ClientOverrideRepoError(t *testing.T) {
	cache := &fakeBasePriceCache{}
	provider := &fakeProductProvider{}

	variantID := uuid.New()
	productID := uuid.New()
	clientID := uuid.New()

	basePrice, _ := domain.NewMoney(100, "EUR")
	cache.SetBasePrice(context.Background(), productID, variantID, basePrice)

	provider.info = &application.ProductPricingInfo{
		VariantID: variantID,
		ProductID: productID,
		BaseCost:  decimal.NewFromInt(80),
		Currency:  "EUR",
		BrandID:   uuid.New(),
	}

	// Repository error â†’ should log warning and fall through to normal flow
	clientPricingRepo := &fakeClientPricingRepo{overrideErr: fmt.Errorf("db connection error")}

	service := application.NewPricingEngineService(
		&fakeBaseRuleRepo{}, &fakeSaleRuleRepo{}, provider, cache, nil,
		clientPricingRepo, nil,
	)

	resp, err := service.CalculateFinalSalePrice(context.Background(), application.CalculateFinalSalePriceRequest{
		ClientID: clientID.String(),
		SaleItems: []application.SaleItemRequest{{
			ProductVariantID: variantID,
			Quantity:         1,
		}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Falls through to normal flow (no rules, no discount) â†’ base price
	if resp.CalculatedItems[0].FinalPrice.Amount != 100 {
		t.Fatalf("expected base price 100, got %.2f", resp.CalculatedItems[0].FinalPrice.Amount)
	}
}

func TestPricingEngineService_CalculateFinalSalePrice_AuditTrailSaved(t *testing.T) {
	cache := &fakeBasePriceCache{}
	provider := &fakeProductProvider{}

	variantID := uuid.New()
	productID := uuid.New()
	clientID := uuid.New()

	basePrice, _ := domain.NewMoney(100, "EUR")
	cache.SetBasePrice(context.Background(), productID, variantID, basePrice)

	provider.info = &application.ProductPricingInfo{
		VariantID: variantID,
		ProductID: productID,
		BaseCost:  decimal.NewFromInt(80),
		Currency:  "EUR",
		BrandID:   uuid.New(),
	}

	calcRepo := &fakeCalculationRepo{}

	service := application.NewPricingEngineService(
		&fakeBaseRuleRepo{}, &fakeSaleRuleRepo{}, provider, cache, nil,
		nil, calcRepo,
	)

	_, err := service.CalculateFinalSalePrice(context.Background(), application.CalculateFinalSalePriceRequest{
		ClientID: clientID.String(),
		SaleItems: []application.SaleItemRequest{{
			ProductVariantID: variantID,
			Quantity:         3,
		}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Wait for async audit
	time.Sleep(200 * time.Millisecond)

	savedCalculations := calcRepo.GetSaved()
	if len(savedCalculations) != 1 {
		t.Fatalf("expected 1 audit record, got %d", len(savedCalculations))
	}
	saved := savedCalculations[0]
	if saved.ProductVariantID != variantID {
		t.Fatalf("expected variantID %s, got %s", variantID, saved.ProductVariantID)
	}
	if saved.ClientID != clientID {
		t.Fatalf("expected clientID %s, got %s", clientID, saved.ClientID)
	}
	if saved.Quantity != 3 {
		t.Fatalf("expected quantity 3, got %d", saved.Quantity)
	}
}

func TestPricingEngineService_CalculateFinalSalePrice_AuditTrailSaveErrorDoesNotBreakCalculation(t *testing.T) {
	cache := &fakeBasePriceCache{}
	provider := &fakeProductProvider{}

	variantID := uuid.New()
	productID := uuid.New()
	clientID := uuid.New()

	basePrice, _ := domain.NewMoney(100, "EUR")
	cache.SetBasePrice(context.Background(), productID, variantID, basePrice)

	provider.info = &application.ProductPricingInfo{
		VariantID: variantID,
		ProductID: productID,
		BaseCost:  decimal.NewFromInt(80),
		Currency:  "EUR",
		BrandID:   uuid.New(),
	}

	calcRepo := &fakeCalculationRepo{saveErr: fmt.Errorf("audit save failed")}

	service := application.NewPricingEngineService(
		&fakeBaseRuleRepo{}, &fakeSaleRuleRepo{}, provider, cache, nil,
		nil, calcRepo,
	)

	resp, err := service.CalculateFinalSalePrice(context.Background(), application.CalculateFinalSalePriceRequest{
		ClientID: clientID.String(),
		SaleItems: []application.SaleItemRequest{{
			ProductVariantID: variantID,
			Quantity:         1,
		}},
	})
	if err != nil {
		t.Fatalf("expected no error despite audit save failure, got: %v", err)
	}
	if resp.SaleTotal.Amount != 100 {
		t.Fatalf("expected total 100, got %.2f", resp.SaleTotal.Amount)
	}
}
