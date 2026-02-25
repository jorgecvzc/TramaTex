package domain

import "github.com/google/uuid"

// PriceModifierType represents the type of price modification to apply.
type PriceModifierType string

const (
	ModifierTypeFixed      PriceModifierType = "FIXED"
	ModifierTypePercentage PriceModifierType = "PERCENTAGE"
)

// Attribute represents a configurable characteristic of a product (e.g., "Size", "Color").
// Note: Scope (brand/group) restrictions removed for MVP simplicity.
// Users are responsible for assigning appropriate attributes to products.
type Attribute struct {
	ID        uuid.UUID
	Name      string
	Code      string
	SortOrder int
	Values    []AttributeValue
}

// AttributeValue represents a specific value for an Attribute (e.g., "Large", "Red").
type AttributeValue struct {
	ID               uuid.UUID
	AttributeID      uuid.UUID
	Value            string
	Code             string
	HasPriceModifier bool
	ModifierType     PriceModifierType
	ModifierAmount   float64
}

// NewAttribute creates a new Attribute with validation.
func NewAttribute(name, code string, sortOrder int) (*Attribute, error) {
	if name == "" {
		return nil, NewValidationError("attribute name cannot be empty")
	}
	if code == "" {
		return nil, NewValidationError("attribute code cannot be empty")
	}

	return &Attribute{
		ID:        uuid.New(),
		Name:      name,
		Code:      code,
		SortOrder: sortOrder,
		Values:    make([]AttributeValue, 0),
	}, nil
}

// AddValue creates a new AttributeValue and adds it to the Attribute.
func (a *Attribute) AddValue(value, code string) (*AttributeValue, error) {
	return a.AddValueWithModifier(value, code, false, "", 0)
}

// AddValueWithModifier creates a new AttributeValue with optional price modifier and adds it to the Attribute.
func (a *Attribute) AddValueWithModifier(value, code string, hasPriceModifier bool, modifierType PriceModifierType, modifierAmount float64) (*AttributeValue, error) {
	if value == "" {
		return nil, NewValidationError("attribute value cannot be empty")
	}
	if code == "" {
		return nil, NewValidationError("attribute value code cannot be empty")
	}

	if hasPriceModifier {
		if modifierType != ModifierTypeFixed && modifierType != ModifierTypePercentage {
			return nil, NewValidationError("invalid modifier type: must be FIXED or PERCENTAGE")
		}
		if modifierType == ModifierTypePercentage && (modifierAmount < -100 || modifierAmount > 1000) {
			return nil, NewValidationError("percentage modifier must be between -100 and 1000")
		}
	}

	newValue := AttributeValue{
		ID:               uuid.New(),
		AttributeID:      a.ID,
		Value:            value,
		Code:             code,
		HasPriceModifier: hasPriceModifier,
		ModifierType:     modifierType,
		ModifierAmount:   modifierAmount,
	}

	a.Values = append(a.Values, newValue)
	return &newValue, nil
}

// UpdateValue updates an existing AttributeValue.
func (a *Attribute) UpdateValue(id uuid.UUID, newValue, newCode string) error {
	if newValue == "" {
		return NewValidationError("attribute value cannot be empty")
	}
	if newCode == "" {
		return NewValidationError("attribute value code cannot be empty")
	}

	for i, val := range a.Values {
		if val.ID == id {
			a.Values[i].Value = newValue
			a.Values[i].Code = newCode
			return nil
		}
	}
	return NewNotFoundErrorf("attribute value with ID %s not found", id)
}

// RemoveValue removes an AttributeValue by its ID.
func (a *Attribute) RemoveValue(id uuid.UUID) error {
	for i, val := range a.Values {
		if val.ID == id {
			a.Values = append(a.Values[:i], a.Values[i+1:]...)
			return nil
		}
	}
	return NewNotFoundErrorf("attribute value with ID %s not found", id)
}

// CalculateBaseCost calculates the base cost by applying attribute value price modifiers to a base price.
// This is used to dynamically calculate variant costs without storing them.
func CalculateBaseCost(basePrice float64, attributeValues []AttributeValue) float64 {
	baseCost := basePrice

	for _, attrValue := range attributeValues {
		if !attrValue.HasPriceModifier {
			continue
		}

		switch attrValue.ModifierType {
		case ModifierTypeFixed:
			// Fixed amount: simply add or subtract the modifier
			baseCost += attrValue.ModifierAmount
		case ModifierTypePercentage:
			// Percentage: apply as percentage of current base cost
			baseCost += baseCost * (attrValue.ModifierAmount / 100.0)
		}
	}

	// Ensure the result is not negative
	if baseCost < 0 {
		baseCost = 0
	}

	return baseCost
}
