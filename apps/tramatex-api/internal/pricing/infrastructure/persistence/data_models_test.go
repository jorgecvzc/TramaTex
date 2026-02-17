package persistence

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/joran-cortez/tramatex/internal/pricing/domain"
)

func TestPricingDataModels_ToDomainAndFromDomain(t *testing.T) {
	percentage, _ := domain.NewPercentage(0.1)
	money, _ := domain.NewMoney(5, "EUR")

	pricingRule, _ := domain.NewPricingRule("Rule", nil, nil, percentage, 1, nil, time.Now(), nil)
	pricingModel := PricingRuleFromDomain(pricingRule)
	mappedRule, err := pricingModel.ToDomain()
	if err != nil || mappedRule.Name != pricingRule.Name {
		t.Fatalf("expected pricing rule mapping")
	}

	clientPricing, _ := domain.NewClientPricing(uuid.New(), uuid.New(), money, time.Now(), nil)
	clientModel := ClientPricingFromDomain(clientPricing)
	mappedClient, err := clientModel.ToDomain()
	if err != nil || mappedClient.ClientID != clientPricing.ClientID {
		t.Fatalf("expected client pricing mapping")
	}

	margin, _ := domain.NewBrandProfitMargin(uuid.New(), &percentage, &money, time.Now(), nil)
	marginModel := BrandProfitMarginFromDomain(margin)
	mappedMargin, err := marginModel.ToDomain()
	if err != nil || mappedMargin.BrandID != margin.BrandID {
		t.Fatalf("expected margin mapping")
	}

	discount, _ := domain.NewSalesDiscountRule("Discount", nil, nil, nil, domain.DiscountTypeFixed, nil, &money, 1, time.Now(), nil)
	discountModel := SalesDiscountRuleFromDomain(discount)
	mappedDiscount, err := discountModel.ToDomain()
	if err != nil || mappedDiscount.Name != discount.Name {
		t.Fatalf("expected discount mapping")
	}

	value, _ := domain.NewRuleValue(domain.RuleValueFixedAmountIncrease, nil, &money)
	baseRule, _ := domain.NewBaseSalesPriceRule("Base", nil, nil, nil, nil, value)
	baseModel := BaseSalesPriceRuleFromDomain(baseRule)
	mappedBase, err := baseModel.ToDomain()
	if err != nil || mappedBase.Name != baseRule.Name {
		t.Fatalf("expected base rule mapping")
	}

	saleRule, _ := domain.NewSaleModificationRule("Sale", nil, nil, &money, value, 1, time.Now(), nil)
	saleModel := SaleModificationRuleFromDomain(saleRule)
	mappedSale, err := saleModel.ToDomain()
	if err != nil || mappedSale.Name != saleRule.Name {
		t.Fatalf("expected sale rule mapping")
	}

	calc, _ := domain.NewPriceCalculation(uuid.New(), uuid.New(), 1, money, money, []string{"Rule"})
	calcModel, err := PriceCalculationFromDomain(calc)
	if err != nil {
		t.Fatalf("expected price calculation mapping")
	}
	mappedCalc, err := calcModel.ToDomain()
	if err != nil || mappedCalc.Quantity != calc.Quantity {
		t.Fatalf("expected price calculation mapping")
	}

	calcModel.AppliedRules = "{invalid"
	if _, err := calcModel.ToDomain(); err == nil {
		t.Fatalf("expected applied rules json error")
	}

	calcModel.AppliedRules = string(mustJSON([]string{"Rule"}))
	if _, err := calcModel.ToDomain(); err != nil {
		t.Fatalf("expected valid applied rules")
	}
}

func mustJSON(value []string) []byte {
	data, _ := json.Marshal(value)
	return data
}
