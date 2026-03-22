package domain

import "net/http"

const (
	MESErrValidation = "VALIDATION"
	MESErrNotFound   = "NOT_FOUND"
	MESErrInternal   = "INTERNAL"
)

// MESError is a typed domain error for the MES module.
// Implements shared/domain.HTTPStatuser for global error middleware handling.
type MESError struct {
	Code    string
	Message string
}

func (e MESError) Error() string { return e.Message }

func (e MESError) HTTPStatus() int {
	switch e.Code {
	case MESErrValidation:
		return http.StatusBadRequest
	case MESErrNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

func NewValidationError(message string) error {
	return MESError{Code: MESErrValidation, Message: message}
}

func NewNotFoundError(message string) error {
	return MESError{Code: MESErrNotFound, Message: message}
}
