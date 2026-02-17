package domain

import "fmt"

type ErrorCode string

const (
	ErrCodeValidation  ErrorCode = "VALIDATION_ERROR"
	ErrCodeNotFound    ErrorCode = "NOT_FOUND"
	ErrCodeConflict    ErrorCode = "CONFLICT"
	ErrCodePersistence ErrorCode = "PERSISTENCE_ERROR"
)

type PartyError struct {
	Code    ErrorCode
	Message string
	Cause   error
}

func (e PartyError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s", e.Message, e.Cause.Error())
	}
	return e.Message
}

func (e PartyError) Unwrap() error {
	return e.Cause
}

func NewValidationError(message string) error {
	return PartyError{Code: ErrCodeValidation, Message: message}
}

func NewValidationErrorf(format string, args ...interface{}) error {
	return PartyError{Code: ErrCodeValidation, Message: fmt.Sprintf(format, args...)}
}

func WrapValidation(message string, cause error) error {
	return PartyError{Code: ErrCodeValidation, Message: message, Cause: cause}
}

func WrapValidationf(cause error, format string, args ...interface{}) error {
	return PartyError{Code: ErrCodeValidation, Message: fmt.Sprintf(format, args...), Cause: cause}
}

func NewNotFoundError(message string) error {
	return PartyError{Code: ErrCodeNotFound, Message: message}
}

func NewNotFoundErrorf(format string, args ...interface{}) error {
	return PartyError{Code: ErrCodeNotFound, Message: fmt.Sprintf(format, args...)}
}

func WrapNotFound(message string, cause error) error {
	return PartyError{Code: ErrCodeNotFound, Message: message, Cause: cause}
}

func WrapNotFoundf(cause error, format string, args ...interface{}) error {
	return PartyError{Code: ErrCodeNotFound, Message: fmt.Sprintf(format, args...), Cause: cause}
}

func NewConflictError(message string) error {
	return PartyError{Code: ErrCodeConflict, Message: message}
}

func NewConflictErrorf(format string, args ...interface{}) error {
	return PartyError{Code: ErrCodeConflict, Message: fmt.Sprintf(format, args...)}
}

func WrapConflict(message string, cause error) error {
	return PartyError{Code: ErrCodeConflict, Message: message, Cause: cause}
}

func WrapConflictf(cause error, format string, args ...interface{}) error {
	return PartyError{Code: ErrCodeConflict, Message: fmt.Sprintf(format, args...), Cause: cause}
}

func NewPersistenceError(message string) error {
	return PartyError{Code: ErrCodePersistence, Message: message}
}

func NewPersistenceErrorf(format string, args ...interface{}) error {
	return PartyError{Code: ErrCodePersistence, Message: fmt.Sprintf(format, args...)}
}

func WrapPersistence(message string, cause error) error {
	return PartyError{Code: ErrCodePersistence, Message: message, Cause: cause}
}

func WrapPersistencef(cause error, format string, args ...interface{}) error {
	return PartyError{Code: ErrCodePersistence, Message: fmt.Sprintf(format, args...), Cause: cause}
}
