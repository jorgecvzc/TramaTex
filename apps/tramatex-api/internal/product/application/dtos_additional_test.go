package application_test

import (
	"testing"

	app "github.com/joran-cortez/tramatex/internal/product/application"
	"github.com/joran-cortez/tramatex/internal/product/domain"
	"github.com/stretchr/testify/assert"
)

// TestNewAttributeDTOFromDomain updated after scope system refactoring
// Attributes no longer have BrandID/GroupID scopes - using direct assignment instead
func TestNewAttributeDTOFromDomain(t *testing.T) {
	attr, err := domain.NewAttribute("Color", "C", 1)
	assert.NoError(t, err)

	_, err = attr.AddValue("Red", "R")
	assert.NoError(t, err)

	dto := app.NewAttributeDTOFromDomain(attr)
	assert.NotNil(t, dto)
	assert.Equal(t, attr.ID, dto.ID)
	assert.Equal(t, "Color", dto.Name)
	assert.Equal(t, "C", dto.Code) // Changed from AttributeName to Code
	assert.Equal(t, 1, dto.SortOrder)
	// Scope fields removed after refactoring
	assert.Len(t, dto.Values, 1)
	assert.Equal(t, "Red", dto.Values[0].Value) // Values is now array of AttributeValueDTO
	assert.Equal(t, "R", dto.Values[0].Code)
}
