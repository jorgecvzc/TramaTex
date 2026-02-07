package domain

import "sort"

type DiscountResult struct {
	Name   string
	Amount Money
}

type SellingPriceCalculatorService struct{}

func NewSellingPriceCalculatorService() *SellingPriceCalculatorService {
	return &SellingPriceCalculatorService{}
}

func (s *SellingPriceCalculatorService) CalculateSellingPrice(baseCost Money, markup Percentage, fixedMarkup *Money) (Money, error) {
	price, err := baseCost.Multiply(1 + markup.Value())
	if err != nil {
		return Money{}, err
	}
	if fixedMarkup != nil {
		price, err = price.Add(*fixedMarkup)
		if err != nil {
			return Money{}, err
		}
	}
	return price, nil
}

type SalesDiscountCalculatorService struct{}

func NewSalesDiscountCalculatorService() *SalesDiscountCalculatorService {
	return &SalesDiscountCalculatorService{}
}

func (s *SalesDiscountCalculatorService) ApplyDiscounts(basePrice Money, rules []*SalesDiscountRule) (Money, []DiscountResult, error) {
	if len(rules) == 0 {
		return basePrice, nil, nil
	}

	sorted := make([]*SalesDiscountRule, len(rules))
	copy(sorted, rules)
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].Priority > sorted[j].Priority
	})

	current := basePrice
	applied := make([]DiscountResult, 0)
	for _, rule := range sorted {
		var discount Money
		var err error
		switch rule.DiscountType {
		case DiscountTypePercentage:
			discountAmount := rule.Percentage.Apply(current.Amount())
			discount, err = NewMoney(discountAmount, current.Currency())
		case DiscountTypeFixed:
			discount = *rule.FixedAmount
		default:
			return Money{}, nil, NewRuleError("unsupported discount type")
		}
		if err != nil {
			return Money{}, nil, err
		}

		updated, err := current.Subtract(discount)
		if err != nil {
			return Money{}, nil, err
		}
		current = updated
		applied = append(applied, DiscountResult{Name: rule.Name, Amount: discount})
	}

	return current, applied, nil
}
