package application

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/joran-cortez/tramatex/internal/pricing/domain"
)

type ProductPricingInfo struct {
	VariantID             uuid.UUID
	ProductID             uuid.UUID
	BaseCost              decimal.Decimal
	Currency              string
	BrandID               uuid.UUID
	BrandMarkupPercentage decimal.Decimal // Brand's default markup percentage (e.g., 30.0 = 30%)
	GroupIDs              []uuid.UUID
	TaxRate               decimal.Decimal // Tax rate as percentage (e.g., 21.0 = 21%)
}

type ProductInfoProvider interface {
	GetVariantPricingInfo(ctx context.Context, variantID uuid.UUID) (*ProductPricingInfo, error)
	GetVariantsPricingInfo(ctx context.Context, variantIDs []uuid.UUID) ([]*ProductPricingInfo, error)
}

type PricingEngineService struct {
	baseRuleRepo      domain.BaseSalesPriceRuleRepository
	saleRuleRepo      domain.SaleModificationRuleRepository
	productInfo       ProductInfoProvider
	basePriceCache    BasePriceCache
	clientInfo        ClientInfoProvider
	clientPricingRepo domain.ClientPricingRepository
	calculationRepo   domain.PriceCalculationRepository
}

type ProductVariantsProvider interface {
	ListVariantsPricingInfo(ctx context.Context, productID uuid.UUID) ([]*ProductPricingInfo, error)
}

// ClientInfoProvider abstracts client/party info for pricing calculations (anti-corruption layer)
type ClientInfoProvider interface {
	GetClientDefaultDiscount(ctx context.Context, clientID string) (float64, error)
}

func NewPricingEngineService(
	baseRuleRepo domain.BaseSalesPriceRuleRepository,
	saleRuleRepo domain.SaleModificationRuleRepository,
	productInfo ProductInfoProvider,
	basePriceCache BasePriceCache,
	clientInfo ClientInfoProvider,
	clientPricingRepo domain.ClientPricingRepository,
	calculationRepo domain.PriceCalculationRepository,
) *PricingEngineService {
	return &PricingEngineService{
		baseRuleRepo:      baseRuleRepo,
		saleRuleRepo:      saleRuleRepo,
		productInfo:       productInfo,
		basePriceCache:    basePriceCache,
		clientInfo:        clientInfo,
		clientPricingRepo: clientPricingRepo,
		calculationRepo:   calculationRepo,
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

	// Fetch all rules once for this calculation
	rules, err := s.baseRuleRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	baseSalesPrice, err := s.calculateBaseSalesPriceFromInfoWithRules(ctx, req.VariantID, info, rules)
	if err != nil {
		return nil, err
	}

	baseCost, err := domain.NewMoneyFromDecimal(info.BaseCost, info.Currency)
	if err != nil {
		return nil, err
	}

	if s.basePriceCache != nil {
		_ = s.basePriceCache.SetBasePrice(ctx, info.ProductID, req.VariantID, baseSalesPrice)
		s.primeCacheForProduct(ctx, info.ProductID)
	}

	tr, _ := info.TaxRate.Float64()
	return &CalculatedBaseSalesPriceResponse{
		VariantID:      req.VariantID,
		BaseCost:       NewMoneyDTO(baseCost),
		BaseSalesPrice: NewMoneyDTO(baseSalesPrice),
		TaxRate:        tr,
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

	// 1. Bulk fetch all product variant info
	variantIDs := make([]uuid.UUID, len(req.SaleItems))
	for i, item := range req.SaleItems {
		variantIDs[i] = item.ProductVariantID
	}
	infoList, err := s.productInfo.GetVariantsPricingInfo(ctx, variantIDs)
	if err != nil {
		return nil, err
	}
	infoMap := make(map[uuid.UUID]*ProductPricingInfo)
	for _, info := range infoList {
		infoMap[info.VariantID] = info
	}

	// 2. Pre-fetch all base pricing rules once
	baseRules, err := s.baseRuleRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	// 3. Parse ClientID if present
	var clientUUID uuid.UUID
	if req.ClientID != "" {
		clientUUID, _ = uuid.Parse(req.ClientID)
	}

	// 4. Pre-fetch all client pricing overrides (Bulk)
	clientOverrides := make(map[uuid.UUID]*domain.ClientPricing)
	if clientUUID != uuid.Nil && s.clientPricingRepo != nil {
		clientOverrides, err = s.clientPricingRepo.FindApplicableBulk(ctx, clientUUID, variantIDs, saleDate)
		if err != nil {
			slog.Warn("failed to bulk query client pricing overrides", "clientID", req.ClientID, "error", err)
		}
	}

	// 5. Pre-fetch ALL active sale modification rules once
	allSaleRules, err := s.saleRuleRepo.ListActive(ctx, saleDate)
	if err != nil {
		return nil, err
	}

	items := make([]CalculatedSaleItemResponse, 0, len(req.SaleItems))
	initialOrderTotal, err := domain.NewMoney(0, domain.DefaultCurrency)
	if err != nil {
		return nil, err
	}
	
	type itemState struct {
		basePrice   domain.Money
		productInfo *ProductPricingInfo
	}
	states := make([]itemState, len(req.SaleItems))

	// First pass: Calculate Base Prices and Order Total
	for i, item := range req.SaleItems {
		if item.Quantity <= 0 {
			return nil, domain.NewValidationError("quantity must be greater than zero")
		}

		productInfo, found := infoMap[item.ProductVariantID]
		if !found {
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
			basePrice, err = s.calculateBaseSalesPriceFromInfoWithRules(ctx, item.ProductVariantID, productInfo, baseRules)
			if err != nil {
				return nil, err
			}
			if s.basePriceCache != nil {
				_ = s.basePriceCache.SetBasePrice(ctx, productInfo.ProductID, item.ProductVariantID, basePrice)
				// Prime cache as a background task to avoid blocking
				go s.primeCacheForProduct(context.Background(), productInfo.ProductID)
			}
		}

		states[i] = itemState{basePrice: basePrice, productInfo: productInfo}

		lineTotal, err := basePrice.Multiply(float64(item.Quantity))
		if err != nil {
			return nil, err
		}
		initialOrderTotal, err = initialOrderTotal.Add(lineTotal)
		if err != nil {
			return nil, err
		}
	}

	// Fetch client default discount as fallback
	var clientDefaultDiscount float64
	if req.ClientID != "" && s.clientInfo != nil {
		clientDefaultDiscount, _ = s.clientInfo.GetClientDefaultDiscount(ctx, req.ClientID)
	}

	saleTotal, err := domain.NewMoney(0, initialOrderTotal.Currency())
	if err != nil {
		return nil, err
	}
	
	auditCalculations := make([]*domain.PriceCalculation, 0, len(req.SaleItems))

	// Second pass: Apply overrides and sale rules
	for i, item := range req.SaleItems {
		state := states[i]
		basePrice := state.basePrice
		productInfo := state.productInfo
		
		baseCost, err := domain.NewMoneyFromDecimal(productInfo.BaseCost, productInfo.Currency)
		if err != nil {
			return nil, err
		}

		finalPrice := basePrice
		var appliedRules []string

		// --- LOGIC FOR MANUAL OVERRIDES ---
		var manualOverride bool
		if item.ManualUnitPrice != nil {
			finalPrice, _ = domain.NewMoney(item.ManualUnitPrice.Amount, item.ManualUnitPrice.Currency)
			appliedRules = append(appliedRules, "ManualPriceOverride")
			manualOverride = true
		} else if item.ManualDiscountPercent != nil {
			discountMultiplier := 1.0 - (*item.ManualDiscountPercent / 100.0)
			finalPrice, err = basePrice.Multiply(discountMultiplier)
			if err != nil {
				return nil, err
			}
			appliedRules = append(appliedRules, "ManualDiscountOverride")
			manualOverride = true
		}

		// --- AUTOMATIC RULES (Only applied if no manual override exists) ---
		if !manualOverride {
			// 1. Client Pricing Override (Highest Priority)
			if override, exists := clientOverrides[item.ProductVariantID]; exists {
				finalPrice = override.FixedPrice
				appliedRules = append(appliedRules, "ClientPricingOverride")
			} else {
				// 2. Sale Modification Rules (In-memory filtering)
				productGroupID := firstGroupID(productInfo.GroupIDs)
				var applicableRules []*domain.SaleModificationRule
				for _, rule := range allSaleRules {
					if rule.AppliesTo(req.ClientID, productGroupID, initialOrderTotal, saleDate) {
						applicableRules = append(applicableRules, rule)
					}
				}

				if len(applicableRules) > 0 {
					// Rules are already sorted by priority desc from ListActive
					for _, rule := range applicableRules {
						updated, err := rule.Value.Apply(finalPrice)
						if err != nil {
							return nil, err
						}
						finalPrice = updated
						appliedRules = append(appliedRules, rule.Name)
					}
				} else if clientDefaultDiscount > 0 {
					// 3. Client Default Discount (Fallback)
					discountMultiplier := 1.0 - (clientDefaultDiscount / 100.0)
					discounted, err := finalPrice.Multiply(discountMultiplier)
					if err != nil {
						return nil, err
					}
					finalPrice = discounted
					appliedRules = append(appliedRules, "ClientDefaultDiscount")
				}
			}
		}

		discountAmount, err := basePrice.Subtract(finalPrice)
		if err != nil {
			discountAmount, _ = domain.NewMoney(0, basePrice.Currency())
		}
		
		var discountPercent float64
		if basePrice.Amount() > 0 && discountAmount.Amount() > 0 {
			// Precision calculation for percentage
			dp := discountAmount.Decimal().Div(basePrice.Decimal()).Mul(decimal.NewFromInt(100))
			discountPercent, _ = dp.Float64()
		}

		taxRateFactor, _ := productInfo.TaxRate.Div(decimal.NewFromInt(100)).Float64()
		taxAmountPerUnit, err := finalPrice.Multiply(taxRateFactor)
		if err != nil {
			return nil, err
		}

		finalPriceWithTax, err := finalPrice.Add(taxAmountPerUnit)
		if err != nil {
			return nil, err
		}

		lineSubtotal, err := finalPrice.Multiply(float64(item.Quantity))
		if err != nil {
			return nil, err
		}
		lineTaxAmount, err := taxAmountPerUnit.Multiply(float64(item.Quantity))
		if err != nil {
			return nil, err
		}
		lineTotal, err := lineSubtotal.Add(lineTaxAmount)
		if err != nil {
			return nil, err
		}

		tr, _ := productInfo.TaxRate.Float64()
		items = append(items, CalculatedSaleItemResponse{
			ProductVariantID:  item.ProductVariantID,
			Quantity:          item.Quantity,
			BaseCost:          NewMoneyDTO(baseCost),
			BaseSalesPrice:    NewMoneyDTO(basePrice),
			FinalPrice:        NewMoneyDTO(finalPrice),
			DiscountPercent:   discountPercent,
			DiscountAmount:    NewMoneyDTO(discountAmount),
			TaxRate:           tr,
			TaxAmountPerUnit:  NewMoneyDTO(taxAmountPerUnit),
			FinalPriceWithTax: NewMoneyDTO(finalPriceWithTax),
			LineSubtotal:      NewMoneyDTO(lineSubtotal),
			LineTaxAmount:     NewMoneyDTO(lineTaxAmount),
			LineTotal:         NewMoneyDTO(lineTotal),
		})

		// Prepare for async audit trail
		if s.calculationRepo != nil && clientUUID != uuid.Nil {
			calc, calcErr := domain.NewPriceCalculation(item.ProductVariantID, clientUUID, item.Quantity, baseCost, finalPrice, appliedRules)
			if calcErr == nil {
				auditCalculations = append(auditCalculations, calc)
			}
		}

		saleTotal, err = saleTotal.Add(lineSubtotal)
		if err != nil {
			return nil, err
		}
	}

	// --- ASYNC AUDIT TRAIL ---
	if len(auditCalculations) > 0 {
		go func(calcs []*domain.PriceCalculation, repo domain.PriceCalculationRepository) {
			for _, calc := range calcs {
				_ = repo.Save(context.Background(), calc)
			}
		}(auditCalculations, s.calculationRepo)
	}

	saleTotalWithTax, err := domain.NewMoney(0, saleTotal.Currency())
	if err != nil {
		return nil, err
	}
	for _, calcItem := range items {
		lineTotalWithTax, err := domain.NewMoney(calcItem.LineTotal.Amount, calcItem.LineTotal.Currency)
		if err != nil {
			return nil, err
		}
		saleTotalWithTax, err = saleTotalWithTax.Add(lineTotalWithTax)
		if err != nil {
			return nil, err
		}
	}

	return &CalculateFinalSalePriceResponse{
		CalculatedItems:  items,
		SaleTotal:        NewMoneyDTO(saleTotal),
		SaleTotalWithTax: NewMoneyDTO(saleTotalWithTax),
	}, nil
}

func (s *PricingEngineService) CreateClientPricingOverride(ctx context.Context, cmd CreateClientPricingCommand) (*ClientPricingDTO, error) {
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

func (s *PricingEngineService) GetPricingHistory(ctx context.Context, query GetPricingHistoryQuery) ([]*PriceCalculationDTO, error) {
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

func (s *PricingEngineService) calculateBaseSalesPriceFromInfo(ctx context.Context, variantID uuid.UUID, info *ProductPricingInfo) (domain.Money, error) {
	rules, err := s.baseRuleRepo.List(ctx)
	if err != nil {
		return domain.Money{}, err
	}
	return s.calculateBaseSalesPriceFromInfoWithRules(ctx, variantID, info, rules)
}

func (s *PricingEngineService) calculateBaseSalesPriceFromInfoWithRules(ctx context.Context, variantID uuid.UUID, info *ProductPricingInfo, rules []*domain.BaseSalesPriceRule) (domain.Money, error) {
	// 1. Empezamos con el coste base de la variante (ya incluye modificadores de atributos desde el provider)
	baseCost, err := domain.NewMoneyFromDecimal(info.BaseCost, info.Currency)
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
	} else if info.BrandMarkupPercentage.IsPositive() {
		// If no rule, apply brand's default markup percentage
		// baseSalesPrice = baseCost * (1 + markup/100)
		markupFactor, _ := info.BrandMarkupPercentage.Div(decimal.NewFromInt(100)).Float64()
		multiplier := 1.0 + markupFactor
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
