package persistence

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/joran-cortez/tramatex/internal/product/domain"
)

// AttributeDataModel represents the attribute entity in the database.
// Note: Scope fields removed for MVP simplicity.
type AttributeDataModel struct {
	gorm.Model
	ID        uuid.UUID                 `gorm:"type:uuid;primary_key;"`
	Name      string                    `gorm:"not null"`
	Code      string                    `gorm:"uniqueIndex;not null"`
	SortOrder int                       `gorm:"not null;default:0"`
	Values    []AttributeValueDataModel `gorm:"foreignKey:AttributeID"`
}

func (AttributeDataModel) TableName() string {
	return "attributes"
}

// AttributeValueDataModel represents the attribute value entity in the database.
type AttributeValueDataModel struct {
	gorm.Model
	ID               uuid.UUID `gorm:"type:uuid;primary_key;"`
	AttributeID      uuid.UUID `gorm:"not null"`
	Value            string    `gorm:"not null"`
	Code             string    `gorm:"not null"`
	HasPriceModifier bool      `gorm:"not null;default:false"`
	ModifierType     *string   `gorm:"type:varchar(20);check:modifier_type IN ('FIXED', 'PERCENTAGE')"`
	ModifierAmount   *float64  `gorm:"type:numeric(10,2)"`
}

func (AttributeValueDataModel) TableName() string {
	return "attribute_values"
}

// ToDomain converts the attribute data model to a domain model.
func (a *AttributeDataModel) ToDomain() *domain.Attribute {
	values := make([]domain.AttributeValue, len(a.Values))
	for i, v := range a.Values {
		values[i] = *v.ToDomain()
	}

	return &domain.Attribute{
		ID:        a.ID,
		Name:      a.Name,
		Code:      a.Code,
		SortOrder: a.SortOrder,
		Values:    values,
	}
}

// ToDomain converts the attribute value data model to a domain model.
func (av *AttributeValueDataModel) ToDomain() *domain.AttributeValue {
	var modifierType domain.PriceModifierType
	var modifierAmount float64
	if av.ModifierType != nil {
		modifierType = domain.PriceModifierType(*av.ModifierType)
	}
	if av.ModifierAmount != nil {
		modifierAmount = *av.ModifierAmount
	}

	return &domain.AttributeValue{
		ID:               av.ID,
		AttributeID:      av.AttributeID,
		Value:            av.Value,
		Code:             av.Code,
		HasPriceModifier: av.HasPriceModifier,
		ModifierType:     modifierType,
		ModifierAmount:   modifierAmount,
	}
}

// AttributeFromDomain converts an attribute domain model to a data model.
func AttributeFromDomain(a *domain.Attribute) *AttributeDataModel {
	values := make([]AttributeValueDataModel, len(a.Values))
	for i, v := range a.Values {
		values[i] = *AttributeValueFromDomain(&v)
	}

	return &AttributeDataModel{
		ID:        a.ID,
		Name:      a.Name,
		Code:      a.Code,
		SortOrder: a.SortOrder,
		Values:    values,
	}
}

// AttributeValueFromDomain converts an attribute value domain model to a data model.
func AttributeValueFromDomain(av *domain.AttributeValue) *AttributeValueDataModel {
	var modifierType *string
	var modifierAmount *float64

	// Only set ModifierType and ModifierAmount if HasPriceModifier is true
	if av.HasPriceModifier && av.ModifierType != "" {
		mt := string(av.ModifierType)
		modifierType = &mt
		ma := av.ModifierAmount
		modifierAmount = &ma
	}

	return &AttributeValueDataModel{
		ID:               av.ID,
		AttributeID:      av.AttributeID,
		Value:            av.Value,
		Code:             av.Code,
		HasPriceModifier: av.HasPriceModifier,
		ModifierType:     modifierType,
		ModifierAmount:   modifierAmount,
	}
}
