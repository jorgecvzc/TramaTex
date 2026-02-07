package application_test

import (
	"testing"

	"github.com/google/uuid"
	app "github.com/joran-cortez/tramatex/internal/product/application"
	"github.com/joran-cortez/tramatex/internal/product/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewAttributeDTOFromDomain(t *testing.T) {
	brandID := uuid.New()
	groupID := uuid.New()
	attr, err := domain.NewAttribute("Color", "C", 1, &brandID, &groupID)
	assert.NoError(t, err)

	_, err = attr.AddValue("Red", "R")
	assert.NoError(t, err)

	dto := app.NewAttributeDTOFromDomain(attr)
	assert.NotNil(t, dto)
	assert.Equal(t, attr.ID, dto.ID)
	assert.Equal(t, "Color", dto.Name)
	assert.Equal(t, "C", dto.AttributeName)
	assert.Equal(t, 1, dto.SortOrder)
	assert.Equal(t, &brandID, dto.ScopeBrandID)
	assert.Equal(t, &groupID, dto.ScopeGroupID)
	assert.Len(t, dto.Values, 1)
	assert.Equal(t, "Red", dto.Values[0])
}
