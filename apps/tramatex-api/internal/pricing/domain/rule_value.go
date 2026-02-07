package domain

type RuleValueType string

const (
	RuleValuePercentageMarkup          RuleValueType = "PERCENTAGE_MARKUP"
	RuleValueFixedAmountIncrease       RuleValueType = "FIXED_AMOUNT_INCREASE"
	RuleValueSetToFixedPrice           RuleValueType = "SET_TO_FIXED_PRICE"
	RuleValueApplyPercentageDiscount   RuleValueType = "APPLY_PERCENTAGE_DISCOUNT"
	RuleValueApplyFixedAmountDiscount  RuleValueType = "APPLY_FIXED_AMOUNT_DISCOUNT"
	RuleValueSetToFixedDiscountedPrice RuleValueType = "SET_TO_FIXED_DISCOUNTED_PRICE"
)

type RuleValue struct {
	Type            RuleValueType
	PercentageValue *Percentage
	MoneyValue      *Money
}

func NewRuleValue(ruleType RuleValueType, percentage *Percentage, money *Money) (RuleValue, error) {
	switch ruleType {
	case RuleValuePercentageMarkup, RuleValueApplyPercentageDiscount:
		if percentage == nil {
			return RuleValue{}, NewValidationError("percentage value is required")
		}
	case RuleValueFixedAmountIncrease, RuleValueApplyFixedAmountDiscount, RuleValueSetToFixedPrice, RuleValueSetToFixedDiscountedPrice:
		if money == nil {
			return RuleValue{}, NewValidationError("money value is required")
		}
	default:
		return RuleValue{}, NewValidationError("unsupported rule value type")
	}

	return RuleValue{Type: ruleType, PercentageValue: percentage, MoneyValue: money}, nil
}

func (v RuleValue) Apply(base Money) (Money, error) {
	switch v.Type {
	case RuleValuePercentageMarkup:
		return base.Multiply(1 + v.PercentageValue.Value())
	case RuleValueFixedAmountIncrease:
		return base.Add(*v.MoneyValue)
	case RuleValueSetToFixedPrice:
		return *v.MoneyValue, nil
	case RuleValueApplyPercentageDiscount:
		discount := v.PercentageValue.Apply(base.Amount())
		discountMoney, err := NewMoney(discount, base.Currency())
		if err != nil {
			return Money{}, err
		}
		return base.Subtract(discountMoney)
	case RuleValueApplyFixedAmountDiscount:
		return base.Subtract(*v.MoneyValue)
	case RuleValueSetToFixedDiscountedPrice:
		return *v.MoneyValue, nil
	default:
		return Money{}, NewRuleError("unsupported rule value type")
	}
}
