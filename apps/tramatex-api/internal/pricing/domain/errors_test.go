package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDomainError(t *testing.T) {
	err := DomainError{Code: ErrCodeValidation, Message: "bad input"}
	require.Equal(t, "bad input", err.Error())

	require.Equal(t, ErrCodeValidation, NewValidationError("x").(DomainError).Code)
	require.Equal(t, ErrCodeNotFound, NewNotFoundError("x").(DomainError).Code)
	require.Equal(t, ErrCodeConflict, NewConflictError("x").(DomainError).Code)
	require.Equal(t, ErrCodeRule, NewRuleError("x").(DomainError).Code)
}
