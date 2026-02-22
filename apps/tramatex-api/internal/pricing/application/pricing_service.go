package application

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/joran-cortez/tramatex/internal/pricing/domain"
)

type ProductPricingInfo struct {
	VariantID             uuid.UUID
	ProductID             uuid.UUID
	BaseCost              float64
	Currency              string
	BrandID               uuid.UUID
	BrandMarkupPercentage float64 // Brand's default markup percentage (e.g., 30.0 = 30%)
	GroupIDs              []uuid.UUID
	TaxRate               float64 // Tax rate as percentage (e.g., 21.0 = 21%)
}

type ProductInfoProvider interface {
	GetVariantPricingInfo(ctx context.Context, variantID uuid.UUID) (*ProductPricingInfo, error)
}

type PricingService struct {
	pricingRuleRepo    domain.PricingRuleRepository
	clientPricingRepo  domain.ClientPricingRepository
	brandMarginRepo    domain.BrandProfitMarginRepository
	discountRuleRepo   domain.SalesDiscountRuleRepository
	calculationRepo    domain.PriceCalculationRepository
	productInfo        ProductInfoProvider
	sellingCalculator  *domain.SellingPriceCalculatorService
	discountCalculator *domain.SalesDiscountCalculatorService
}

func NewPricingService(
	pricingRuleRepo domain.PricingRuleRepository,
	clientPricingRepo domain.ClientPricingRepository,
	brandMarginRepo domain.BrandProfitMarginRepository,
	discountRuleRepo domain.SalesDiscountRuleRepository,
	calculationRepo domain.PriceCalculationRepository,
	productInfo ProductInfoProvider,
) *PricingService {
	return &PricingService{
		pricingRuleRepo:    pricingRuleRepo,
		clientPricingRepo:  clientPricingRepo,
		brandMarginRepo:    brandMarginRepo,
		discountRuleRepo:   discountRuleRepo,
		calculationRepo:    calculationRepo,
		productInfo:        productInfo,
		sellingCalculator:  domain.NewSellingPriceCalculatorService(),
		discountCalculator: domain.NewSalesDiscountCalculatorService(),
	}
}

func (s *PricingService) CreatePricingRule(ctx context.Context, cmd CreatePricingRuleCommand) (*PricingRuleDTO, error) {
	percentage, err := domain.NewPercentage(cmd.MarkupPercentage)
	if err != nil {
		return nil, err
	}
	if cmd.EffectiveFrom.IsZero() {
		cmd.EffectiveFrom = time.Now()
	}

	rule, err := domain.NewPricingRule(
		cmd.Name,
		cmd.ProductVariantID,
		cmd.PartyCategory,
		percentage,
		cmd.MinQuantity,
		cmd.MaxQuantity,
		cmd.EffectiveFrom,
		cmd.EffectiveTo,
	)
	if err != nil {
		return nil, err
	}

	if err := s.pricingRuleRepo.Save(ctx, rule); err != nil {
		return nil, err
	}

	dto := NewPricingRuleDTO(rule)
	return &dto, nil
}

func (s *PricingService) ListPricingRules(ctx context.Context, query ListPricingRulesQuery) ([]*PricingRuleDTO, error) {
	rules, err := s.pricingRuleRepo.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]*PricingRuleDTO, 0, len(rules))
	for _, rule := range rules {
		dto := NewPricingRuleDTO(rule)
		result = append(result, &dto)
	}
	return result, nil
}

func (s *PricingService) CreateClientPricing(ctx context.Context, cmd CreateClientPricingCommand) (*ClientPricingDTO, error) {
	if cmd.EffectiveFrom.IsZero() {
		cmd.EffectiveFrom = time.Now()
	}

	price, err := domain.NewMoney(cmd.FixedPrice, cmd.Currency)
	if err != nil {
		return nil, err
	}

	override, err := domain.NewClientPricing(cmd.ClientID, cmd.ProductVariantID, price, cmd.EffectiveFrom, cmd.EffectiveTo)
	if err != nil {
		return nil, err
	}

	if err := s.clientPricingRepo.Save(ctx, override); err != nil {
		return nil, err
	}

	dto := NewClientPricingDTO(override)
	return &dto, nil
}

func (s *PricingService) CalculatePrice(ctx context.Context, cmd CalculatePriceCommand) (*CalculatePriceResponse, error) {
	if cmd.Quantity <= 0 {
		return nil, domain.NewValidationError("quantity must be greater than zero")
	}

	info, err := s.productInfo.GetVariantPricingInfo(ctx, cmd.ProductVariantID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, domain.NewNotFoundError("product variant not found")
	}

	baseCost, err := domain.NewMoney(info.BaseCost, info.Currency)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	appliedRules := make([]string, 0)

	override, err := s.clientPricingRepo.FindApplicable(ctx, cmd.ClientID, cmd.ProductVariantID, now)
	if err != nil {
		return nil, err
	}
	if override != nil {
		calc, err := domain.NewPriceCalculation(cmd.ProductVariantID, cmd.ClientID, cmd.Quantity, baseCost, override.FixedPrice, []string{"ClientOverride"})
		if err != nil {
			return nil, err
		}
		if err := s.calculationRepo.Save(ctx, calc); err != nil {
			return nil, err
		}

		return &CalculatePriceResponse{
			FinalPrice: override.FixedPrice.Amount(),
			Currency:   override.FixedPrice.Currency(),
			Breakdown: PriceBreakdownDTO{
				BaseCost:      baseCost.Amount(),
				MarginApplied: "0%",
				Discounts:     []DiscountBreakdownDTO{},
			},
		}, nil
	}

	markup := domain.Percentage{}
	markupSet := false
	var fixedMarkup *domain.Money

	rules, err := s.pricingRuleRepo.FindApplicable(ctx, cmd.ProductVariantID, cmd.Quantity, now)
	if err != nil {
		return nil, err
	}
	if len(rules) > 0 {
		sort.SliceStable(rules, func(i, j int) bool {
			return rules[i].MinQuantity > rules[j].MinQuantity
		})
		markup = rules[0].Markup
		markupSet = true
		appliedRules = append(appliedRules, rules[0].Name)
	}

	if !markupSet {
		margin, err := s.brandMarginRepo.FindApplicable(ctx, info.BrandID, now)
		if err != nil {
			return nil, err
		}
		if margin != nil {
			if margin.Percentage != nil {
				markup = *margin.Percentage
				markupSet = true
			}
			if margin.FixedAmount != nil {
				fixedMarkup = margin.FixedAmount
			}
			appliedRules = append(appliedRules, "BrandMargin")
		}
	}

	if !markupSet {
		markup, err = domain.NewPercentage(0)
		if err != nil {
			return nil, err
		}
	}

	basePrice, err := s.sellingCalculator.CalculateSellingPrice(baseCost, markup, fixedMarkup)
	if err != nil {
		return nil, err
	}

	discountRules, err := s.discountRuleRepo.FindApplicable(ctx, cmd.ClientID, cmd.ProductVariantID, cmd.Quantity, now)
	if err != nil {
		return nil, err
	}

	finalPrice, discounts, err := s.discountCalculator.ApplyDiscounts(basePrice, discountRules)
	if err != nil {
		return nil, err
	}

	discountBreakdown := make([]DiscountBreakdownDTO, 0, len(discounts))
	for _, discount := range discounts {
		discountBreakdown = append(discountBreakdown, DiscountBreakdownDTO{
			Name:   discount.Name,
			Amount: discount.Amount.Amount(),
		})
		appliedRules = append(appliedRules, discount.Name)
	}

	calc, err := domain.NewPriceCalculation(cmd.ProductVariantID, cmd.ClientID, cmd.Quantity, baseCost, finalPrice, appliedRules)
	if err != nil {
		return nil, err
	}
	if err := s.calculationRepo.Save(ctx, calc); err != nil {
		return nil, err
	}

	marginLabel := fmt.Sprintf("%.2f%%", markup.Value()*100)
	return &CalculatePriceResponse{
		FinalPrice: finalPrice.Amount(),
		Currency:   finalPrice.Currency(),
		Breakdown: PriceBreakdownDTO{
			BaseCost:      baseCost.Amount(),
			MarginApplied: marginLabel,
			Discounts:     discountBreakdown,
		},
	}, nil
}

func (s *PricingService) GetPricingHistory(ctx context.Context, query GetPricingHistoryQuery) ([]*PriceCalculationDTO, error) {
	calcs, err := s.calculationRepo.ListByProductVariantID(ctx, query.ProductVariantID)
	if err != nil {
		return nil, err
	}
	result := make([]*PriceCalculationDTO, 0, len(calcs))
	for _, calc := range calcs {
		dto := NewPriceCalculationDTO(calc)
		result = append(result, &dto)
	}
	return result, nil
}
