package domain

import "net/http"

const (
	ErrCodeValidation = "VALIDATION_ERROR"
	ErrCodeNotFound   = "NOT_FOUND"
	ErrCodeConflict   = "CONFLICT"
	ErrCodeRule       = "RULE_ERROR"
)

type DomainError struct {
	Code    string
	Message string
}

func (e DomainError) Error() string {
	return e.Message
}

// HTTPStatus returns the HTTP status code corresponding to this error.
// Implements shared/domain.HTTPStatuser for global middleware handling.
func (e DomainError) HTTPStatus() int {
	switch e.Code {
	case ErrCodeValidation:
		return http.StatusBadRequest
	case ErrCodeNotFound:
		return http.StatusNotFound
	case ErrCodeConflict:
		return http.StatusConflict
	case ErrCodeRule:
		return http.StatusUnprocessableEntity
	default:
		return http.StatusBadRequest
	}
}

func NewValidationError(message string) error {
	return DomainError{Code: ErrCodeValidation, Message: message}
}

func NewNotFoundError(message string) error {
	return DomainError{Code: ErrCodeNotFound, Message: message}
}

func NewConflictError(message string) error {
	return DomainError{Code: ErrCodeConflict, Message: message}
}

func NewRuleError(message string) error {
	return DomainError{Code: ErrCodeRule, Message: message}
}
