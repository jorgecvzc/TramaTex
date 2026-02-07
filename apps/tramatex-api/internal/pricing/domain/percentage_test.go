package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPercentageValidation(t *testing.T) {
	_, err := NewPercentage(-0.1)
	require.Error(t, err)

	_, err = NewPercentage(1.1)
	require.Error(t, err)
}

func TestPercentageApply(t *testing.T) {
	p, err := NewPercentage(0.2)
	require.NoError(t, err)
	require.Equal(t, 0.2, p.Value())
	require.Equal(t, 20.0, p.Apply(100))
}
