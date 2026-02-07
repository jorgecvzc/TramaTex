package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestBrandProfitMarginValidation(t *testing.T) {
	brandID := uuid.New()
	p, _ := NewPercentage(0.1)

	_, err := NewBrandProfitMargin(uuid.Nil, &p, nil, time.Now(), nil)
	require.Error(t, err)

	_, err = NewBrandProfitMargin(brandID, nil, nil, time.Now(), nil)
	require.Error(t, err)

	end := time.Now()
	start := end.Add(time.Hour)
	_, err = NewBrandProfitMargin(brandID, &p, nil, start, &end)
	require.Error(t, err)
}

func TestBrandProfitMarginAppliesTo(t *testing.T) {
	brandID := uuid.New()
	p, _ := NewPercentage(0.1)
	start := time.Now().Add(-time.Hour)
	end := time.Now().Add(time.Hour)

	margin, err := NewBrandProfitMargin(brandID, &p, nil, start, &end)
	require.NoError(t, err)

	require.False(t, margin.AppliesTo(uuid.New(), time.Now()))
	require.False(t, margin.AppliesTo(brandID, start.Add(-time.Minute)))
	require.False(t, margin.AppliesTo(brandID, end.Add(time.Minute)))
	require.True(t, margin.AppliesTo(brandID, time.Now()))

	margin.IsActive = false
	require.False(t, margin.AppliesTo(brandID, time.Now()))
}
