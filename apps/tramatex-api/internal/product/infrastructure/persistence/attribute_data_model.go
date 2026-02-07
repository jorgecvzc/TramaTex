package persistence

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/joran-cortez/tramatex/internal/product/domain"
)

// AttributeDataModel represents the attribute entity in the database.
type AttributeDataModel struct {
	gorm.Model
	ID           uuid.UUID `gorm:"type:uuid;primary_key;"`
	Name         string    `gorm:"not null"`
	Code         string    `gorm:"uniqueIndex;not null"`
	SortOrder    int       `gorm:"not null;default:0"`
	ScopeBrandID *uuid.UUID
	ScopeGroupID *uuid.UUID
	Values       []AttributeValueDataModel `gorm:"foreignKey:AttributeID"`
	CreatedBy    string
	ModifiedBy   string
}

func (AttributeDataModel) TableName() string {
	return "attributes"
}

// AttributeValueDataModel represents the attribute value entity in the database.
type AttributeValueDataModel struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;primary_key;"`
	AttributeID uuid.UUID `gorm:"not null"`
	Value       string    `gorm:"not null"`
	Code        string    `gorm:"not null"`
	CreatedBy   string
	ModifiedBy  string
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
		ID:           a.ID,
		Name:         a.Name,
		Code:         a.Code,
		SortOrder:    a.SortOrder,
		ScopeBrandID: a.ScopeBrandID,
		ScopeGroupID: a.ScopeGroupID,
		Values:       values,
	}
}

// ToDomain converts the attribute value data model to a domain model.
func (av *AttributeValueDataModel) ToDomain() *domain.AttributeValue {
	return &domain.AttributeValue{
		ID:          av.ID,
		AttributeID: av.AttributeID,
		Value:       av.Value,
		Code:        av.Code,
	}
}

// AttributeFromDomain converts an attribute domain model to a data model.
func AttributeFromDomain(a *domain.Attribute) *AttributeDataModel {
	values := make([]AttributeValueDataModel, len(a.Values))
	for i, v := range a.Values {
		values[i] = *AttributeValueFromDomain(&v)
	}

	return &AttributeDataModel{
		ID:           a.ID,
		Name:         a.Name,
		Code:         a.Code,
		SortOrder:    a.SortOrder,
		ScopeBrandID: a.ScopeBrandID,
		ScopeGroupID: a.ScopeGroupID,
		Values:       values,
	}
}

// AttributeValueFromDomain converts an attribute value domain model to a data model.
func AttributeValueFromDomain(av *domain.AttributeValue) *AttributeValueDataModel {
	return &AttributeValueDataModel{
		ID:          av.ID,
		AttributeID: av.AttributeID,
		Value:       av.Value,
		Code:        av.Code,
	}
}
