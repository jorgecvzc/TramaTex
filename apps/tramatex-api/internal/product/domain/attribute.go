package domain

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

// Attribute represents a configurable characteristic of a product (e.g., "Size", "Color").
type Attribute struct {
	ID             uuid.UUID
	Name           string
	Code           string
	SortOrder      int
	ScopeBrandID   *uuid.UUID
	ScopeGroupID   *uuid.UUID
	Values         []AttributeValue
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// AttributeValue represents a specific value for an Attribute (e.g., "Large", "Red").
type AttributeValue struct {
	ID          uuid.UUID
	AttributeID uuid.UUID
	Value       string
	Code        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewAttribute creates a new Attribute with validation.
func NewAttribute(name, code string, sortOrder int, scopeBrandID, scopeGroupID *uuid.UUID) (*Attribute, error) {
	if name == "" {
		return nil, fmt.Errorf("attribute name cannot be empty")
	}
	if code == "" {
		return nil, fmt.Errorf("attribute code cannot be empty")
	}

	return &Attribute{
		ID:           uuid.New(),
		Name:         name,
		Code:         code,
		SortOrder:    sortOrder,
		ScopeBrandID: scopeBrandID,
		ScopeGroupID: scopeGroupID,
		Values:       make([]AttributeValue, 0),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

// AddValue creates a new AttributeValue and adds it to the Attribute.
func (a *Attribute) AddValue(value, code string) (*AttributeValue, error) {
	if value == "" {
		return nil, fmt.Errorf("attribute value cannot be empty")
	}
	if code == "" {
		return nil, fmt.Errorf("attribute value code cannot be empty")
	}

	newValue := AttributeValue{
		ID:          uuid.New(),
		AttributeID: a.ID,
		Value:       value,
		Code:        code,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	a.Values = append(a.Values, newValue)
	return &newValue, nil
}

// UpdateValue updates an existing AttributeValue.
func (a *Attribute) UpdateValue(id uuid.UUID, newValue, newCode string) error {
	if newValue == "" {
		return fmt.Errorf("attribute value cannot be empty")
	}
	if newCode == "" {
		return fmt.Errorf("attribute value code cannot be empty")
	}

	for i, val := range a.Values {
		if val.ID == id {
			a.Values[i].Value = newValue
			a.Values[i].Code = newCode
			a.Values[i].UpdatedAt = time.Now()
			return nil
		}
	}
	return fmt.Errorf("attribute value with ID %s not found", id)
}

// RemoveValue removes an AttributeValue by its ID.
func (a *Attribute) RemoveValue(id uuid.UUID) error {
	for i, val := range a.Values {
		if val.ID == id {
			a.Values = append(a.Values[:i], a.Values[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("attribute value with ID %s not found", id)
}
