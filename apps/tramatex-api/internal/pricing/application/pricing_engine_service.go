package application

import (
	"context"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/joran-cortez/tramatex/internal/pricing/domain"
)

type PricingEngineService struct {
	baseRuleRepo   domain.BaseSalesPriceRuleRepository
	saleRuleRepo   domain.SaleModificationRuleRepository
	productInfo    ProductInfoProvider
	basePriceCache BasePriceCache
}

type ProductVariantsProvider interface {
	ListVariantsPricingInfo(ctx context.Context, productID uuid.UUID) ([]*ProductPricingInfo, error)
}

func NewPricingEngineService(
	baseRuleRepo domain.BaseSalesPriceRuleRepository,
	saleRuleRepo domain.SaleModificationRuleRepository,
	productInfo ProductInfoProvider,
	basePriceCache BasePriceCache,
) *PricingEngineService {
	return &PricingEngineService{
		baseRuleRepo:   baseRuleRepo,
		saleRuleRepo:   saleRuleRepo,
		productInfo:    productInfo,
		basePriceCache: basePriceCache,
	}
}

func (s *PricingEngineService) CreateBaseSalesPriceRule(ctx context.Context, cmd CreateBaseSalesPriceRuleCommand) (*BaseSalesPriceRuleDTO, error) {
	value, err := toRuleValue(cmd.Value)
	if err != nil {
		return nil, err
	}

	rule, err := domain.NewBaseSalesPriceRule(
		cmd.Name,
		cmd.BrandID,
		cmd.ProductGroupID,
		cmd.ProductID,
		cmd.VariantID,
		value,
	)
	if err != nil {
		return nil, err
	}
	if cmd.IsActive != nil {
		rule.IsActive = *cmd.IsActive
	}

	if err := s.baseRuleRepo.Save(ctx, rule); err != nil {
		return nil, err
	}

	dto := NewBaseSalesPriceRuleDTO(rule)
	return &dto, nil
}

func (s *PricingEngineService) UpdateBaseSalesPriceRule(ctx context.Context, cmd UpdateBaseSalesPriceRuleCommand) (*BaseSalesPriceRuleDTO, error) {
	rule, err := s.baseRuleRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}
	if rule == nil {
		return nil, domain.NewNotFoundError("base sales price rule not found")
	}

	if cmd.Name != nil {
		rule.Name = *cmd.Name
	}
	if cmd.BrandID != nil {
		rule.BrandID = cmd.BrandID
	}
	if cmd.ProductGroupID != nil {
		rule.ProductGroupID = cmd.ProductGroupID
	}
	if cmd.ProductID != nil {
		rule.ProductID = cmd.ProductID
	}
	if cmd.VariantID != nil {
		rule.VariantID = cmd.VariantID
	}
	if cmd.Value != nil {
		value, err := toRuleValue(*cmd.Value)
		if err != nil {
			return nil, err
		}
		rule.Value = value
	}
	if cmd.IsActive != nil {
		rule.IsActive = *cmd.IsActive
	}

	if err := s.baseRuleRepo.Save(ctx, rule); err != nil {
		return nil, err
	}

	dto := NewBaseSalesPriceRuleDTO(rule)
	return &dto, nil
}

func (s *PricingEngineService) CreateSaleModificationRule(ctx context.Context, cmd CreateSaleModificationRuleCommand) (*SaleModificationRuleDTO, error) {
	if cmd.EffectiveFrom.IsZero() {
		cmd.EffectiveFrom = time.Now()
	}

	value, err := toRuleValue(cmd.Value)
	if err != nil {
		return nil, err
	}

	var minOrder *domain.Money
	if cmd.MinOrderTotalAmount != nil {
		money, err := domain.NewMoney(cmd.MinOrderTotalAmount.Amount, cmd.MinOrderTotalAmount.Currency)
		if err != nil {
			return nil, err
		}
		minOrder = &money
	}

	rule, err := domain.NewSaleModificationRule(
		cmd.Name,
		cmd.ClientIDs,
		cmd.ProductGroupID,
		minOrder,
		value,
		cmd.Priority,
		cmd.EffectiveFrom,
		cmd.EffectiveTo,
	)
	if err != nil {
		return nil, err
	}
	if cmd.IsActive != nil {
		rule.IsActive = *cmd.IsActive
	}

	if err := s.saleRuleRepo.Save(ctx, rule); err != nil {
		return nil, err
	}

	dto := NewSaleModificationRuleDTO(rule)
	return &dto, nil
}

func (s *PricingEngineService) UpdateSaleModificationRule(ctx context.Context, cmd UpdateSaleModificationRuleCommand) (*SaleModificationRuleDTO, error) {
	rule, err := s.saleRuleRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}
	if rule == nil {
		return nil, domain.NewNotFoundError("sale modification rule not found")
	}

	if cmd.Name != nil {
		rule.Name = *cmd.Name
	}
	if cmd.ClientIDs != nil {
		rule.ClientIDs = cmd.ClientIDs
	}
	if cmd.ProductGroupID != nil {
		rule.ProductGroupID = cmd.ProductGroupID
	}
	if cmd.MinOrderTotalAmount != nil {
		money, err := domain.NewMoney(cmd.MinOrderTotalAmount.Amount, cmd.MinOrderTotalAmount.Currency)
		if err != nil {
			return nil, err
		}
		rule.MinOrderTotal = &money
	}
	if cmd.Value != nil {
		value, err := toRuleValue(*cmd.Value)
		if err != nil {
			return nil, err
		}
		rule.Value = value
	}
	if cmd.Priority != nil {
		rule.Priority = *cmd.Priority
	}
	if cmd.IsActive != nil {
		rule.IsActive = *cmd.IsActive
	}
	if cmd.EffectiveFrom != nil {
		rule.EffectiveFrom = *cmd.EffectiveFrom
	}
	if cmd.EffectiveTo != nil {
		rule.EffectiveTo = cmd.EffectiveTo
	}

	if err := s.saleRuleRepo.Save(ctx, rule); err != nil {
		return nil, err
	}

	dto := NewSaleModificationRuleDTO(rule)
	return &dto, nil
}

func (s *PricingEngineService) CalculateBaseSalesPrice(ctx context.Context, req CalculateBaseSalesPriceRequest) (*CalculatedBaseSalesPriceResponse, error) {
	info, err := s.productInfo.GetVariantPricingInfo(ctx, req.VariantID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, domain.NewNotFoundError("product variant not found")
	}
	baseSalesPrice, err := s.calculateBaseSalesPriceFromInfo(ctx, req.VariantID, info)
	if err != nil {
		return nil, err
	}

	if s.basePriceCache != nil {
		_ = s.basePriceCache.SetBasePrice(ctx, info.ProductID, req.VariantID, baseSalesPrice)
		s.primeCacheForProduct(ctx, info.ProductID)
	}

	return &CalculatedBaseSalesPriceResponse{
		VariantID:      req.VariantID,
		BaseSalesPrice: NewMoneyDTO(baseSalesPrice),
	}, nil
}

func (s *PricingEngineService) CalculateFinalSalePrice(ctx context.Context, req CalculateFinalSalePriceRequest) (*CalculateFinalSalePriceResponse, error) {
	if len(req.SaleItems) == 0 {
		return nil, domain.NewValidationError("saleItems cannot be empty")
	}

	saleDate := req.SaleDate
	if saleDate.IsZero() {
		saleDate = time.Now()
	}

	items := make([]CalculatedSaleItemResponse, 0, len(req.SaleItems))
	orderTotal, err := domain.NewMoney(0, domain.DefaultCurrency)
	if err != nil {
		return nil, err
	}
	basePrices := make([]domain.Money, 0, len(req.SaleItems))
	productInfos := make([]*ProductPricingInfo, 0, len(req.SaleItems))

	for _, item := range req.SaleItems {
		if item.Quantity <= 0 {
			return nil, domain.NewValidationError("quantity must be greater than zero")
		}

		productInfo, err := s.productInfo.GetVariantPricingInfo(ctx, item.ProductVariantID)
		if err != nil {
			return nil, err
		}
		if productInfo == nil {
			return nil, domain.NewNotFoundError("product variant not found")
		}

		var basePrice domain.Money
		if s.basePriceCache != nil {
			cached, err := s.basePriceCache.GetBasePrice(ctx, productInfo.ProductID, item.ProductVariantID)
			if err == nil && cached != nil {
				basePrice = *cached
			}
		}

		if basePrice.Amount() == 0 {
			basePrice, err = s.calculateBaseSalesPriceFromInfo(ctx, item.ProductVariantID, productInfo)
			if err != nil {
				return nil, err
			}
			if s.basePriceCache != nil {
				_ = s.basePriceCache.SetBasePrice(ctx, productInfo.ProductID, item.ProductVariantID, basePrice)
				s.primeCacheForProduct(ctx, productInfo.ProductID)
			}
		}

		basePrices = append(basePrices, basePrice)
		productInfos = append(productInfos, productInfo)

		lineTotal, err := basePrice.Multiply(float64(item.Quantity))
		if err != nil {
			return nil, err
		}
		orderTotal, err = orderTotal.Add(lineTotal)
		if err != nil {
			return nil, err
		}
	}

	saleTotal, err := domain.NewMoney(0, orderTotal.Currency())
	if err != nil {
		return nil, err
	}

	for index, item := range req.SaleItems {
		basePrice := basePrices[index]
		productInfo := productInfos[index]
		productGroupID := firstGroupID(productInfo.GroupIDs)
		rules, err := s.saleRuleRepo.ListApplicable(ctx, req.ClientID, productGroupID, orderTotal, saleDate)
		if err != nil {
			return nil, err
		}

		finalPrice := basePrice
		if len(rules) > 0 {
			sort.SliceStable(rules, func(i, j int) bool {
				return rules[i].Priority > rules[j].Priority
			})
			for _, rule := range rules {
				updated, err := rule.Value.Apply(finalPrice)
				if err != nil {
					return nil, err
				}
				finalPrice = updated
			}
		}

		items = append(items, CalculatedSaleItemResponse{
			ProductVariantID: item.ProductVariantID,
			Quantity:         item.Quantity,
			BaseSalesPrice:   NewMoneyDTO(basePrice),
			FinalPrice:       NewMoneyDTO(finalPrice),
		})

		lineTotal, err := finalPrice.Multiply(float64(item.Quantity))
		if err != nil {
			return nil, err
		}
		saleTotal, err = saleTotal.Add(lineTotal)
		if err != nil {
			return nil, err
		}
	}

	return &CalculateFinalSalePriceResponse{
		CalculatedItems: items,
		SaleTotal:       NewMoneyDTO(saleTotal),
	}, nil
}

func (s *PricingEngineService) calculateBaseSalesPriceFromInfo(ctx context.Context, variantID uuid.UUID, info *ProductPricingInfo) (domain.Money, error) {
	baseCost, err := domain.NewMoney(info.BaseCost, info.Currency)
	if err != nil {
		return domain.Money{}, err
	}

	rules, err := s.baseRuleRepo.List(ctx)
	if err != nil {
		return domain.Money{}, err
	}

	selected := selectBaseRule(rules, info.BrandID, info.GroupIDs, info.ProductID, variantID)
	baseSalesPrice := baseCost
	if selected != nil {
		// Apply pricing rule if found
		baseSalesPrice, err = selected.Value.Apply(baseCost)
		if err != nil {
			return domain.Money{}, err
		}
	} else if info.BrandMarkupPercentage > 0 {
		// If no rule, apply brand's default markup percentage
		// baseSalesPrice = baseCost * (1 + markup/100)
		multiplier := 1.0 + (info.BrandMarkupPercentage / 100.0)
		baseSalesPrice, err = baseCost.Multiply(multiplier)
		if err != nil {
			return domain.Money{}, err
		}
	}

	return baseSalesPrice, nil
}

func (s *PricingEngineService) primeCacheForProduct(ctx context.Context, productID uuid.UUID) {
	provider, ok := s.productInfo.(ProductVariantsProvider)
	if !ok || s.basePriceCache == nil {
		return
	}
	infos, err := provider.ListVariantsPricingInfo(ctx, productID)
	if err != nil || len(infos) == 0 {
		return
	}
	for _, info := range infos {
		price, err := s.calculateBaseSalesPriceFromInfo(ctx, info.VariantID, info)
		if err != nil {
			continue
		}
		_ = s.basePriceCache.SetBasePrice(ctx, productID, info.VariantID, price)
	}
}

func toRuleValue(dto RuleValueDTO) (domain.RuleValue, error) {
	var percentage *domain.Percentage
	if dto.PercentageValue != nil {
		p, err := domain.NewPercentage(dto.PercentageValue.Value)
		if err != nil {
			return domain.RuleValue{}, err
		}
		percentage = &p
	}

	var money *domain.Money
	if dto.MoneyValue != nil {
		m, err := domain.NewMoney(dto.MoneyValue.Amount, dto.MoneyValue.Currency)
		if err != nil {
			return domain.RuleValue{}, err
		}
		money = &m
	}

	return domain.NewRuleValue(domain.RuleValueType(dto.Type), percentage, money)
}

func selectBaseRule(rules []*domain.BaseSalesPriceRule, brandID uuid.UUID, groupIDs []uuid.UUID, productID uuid.UUID, variantID uuid.UUID) *domain.BaseSalesPriceRule {
	bestScore := -1
	var selected *domain.BaseSalesPriceRule
	for _, rule := range rules {
		if !rule.IsActive {
			continue
		}
		if rule.VariantID != nil && *rule.VariantID != variantID {
			continue
		}
		if rule.ProductID != nil && *rule.ProductID != productID {
			continue
		}
		if rule.ProductGroupID != nil && !containsUUID(groupIDs, *rule.ProductGroupID) {
			continue
		}
		if rule.BrandID != nil && *rule.BrandID != brandID {
			continue
		}

		score := 0
		if rule.BrandID != nil {
			score = 1
		}
		if rule.ProductGroupID != nil {
			score = 2
		}
		if rule.ProductID != nil {
			score = 3
		}
		if rule.VariantID != nil {
			score = 4
		}
		if score > bestScore {
			bestScore = score
			selected = rule
		}
	}
	return selected
}

func containsUUID(values []uuid.UUID, target uuid.UUID) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func firstGroupID(values []uuid.UUID) *uuid.UUID {
	if len(values) == 0 {
		return nil
	}
	value := values[0]
	return &value
}
