package domain_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/joran-cortez/tramatex/internal/product/domain"
)

func TestNewAttribute(t *testing.T) {
	t.Run("should create a new attribute successfully", func(t *testing.T) {
		attr, err := domain.NewAttribute("Color", "C", 1, nil, nil)
		assert.NoError(t, err)
		assert.NotNil(t, attr)
		assert.Equal(t, "Color", attr.Name)
		assert.Equal(t, "C", attr.Code)
		assert.Equal(t, 1, attr.SortOrder)
		assert.NotEqual(t, uuid.Nil, attr.ID)
		assert.Empty(t, attr.Values)
	})

	t.Run("should return error if name is empty", func(t *testing.T) {
		_, err := domain.NewAttribute("", "C", 1, nil, nil)
		assert.Error(t, err)
		assert.EqualError(t, err, "attribute name cannot be empty")
	})

	t.Run("should return error if code is empty", func(t *testing.T) {
		_, err := domain.NewAttribute("Color", "", 1, nil, nil)
		assert.Error(t, err)
		assert.EqualError(t, err, "attribute code cannot be empty")
	})
}

func TestAttribute_AddValue(t *testing.T) {
	attr, _ := domain.NewAttribute("Size", "S", 0, nil, nil)

	t.Run("should add a value successfully", func(t *testing.T) {
		val, err := attr.AddValue("Large", "L")
		assert.NoError(t, err)
		assert.NotNil(t, val)
		assert.Equal(t, "Large", val.Value)
		assert.Equal(t, "L", val.Code)
		assert.Equal(t, attr.ID, val.AttributeID)
		assert.Len(t, attr.Values, 1)
		assert.Equal(t, *val, attr.Values[0])
	})

	t.Run("should return error if value is empty", func(t *testing.T) {
		_, err := attr.AddValue("", "XL")
		assert.Error(t, err)
		assert.EqualError(t, err, "attribute value cannot be empty")
	})

	t.Run("should return error if value code is empty", func(t *testing.T) {
		_, err := attr.AddValue("Extra Large", "")
		assert.Error(t, err)
		assert.EqualError(t, err, "attribute value code cannot be empty")
	})

	t.Run("should add multiple values", func(t *testing.T) {
		_, _ = attr.AddValue("Medium", "M")
		_, _ = attr.AddValue("Small", "S")
		assert.Len(t, attr.Values, 3) // Large, Medium, Small
	})
}
