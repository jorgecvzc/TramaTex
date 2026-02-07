package domain

const (
	ErrCodeValidation = "VALIDATION_ERROR"
	ErrCodeNotFound   = "NOT_FOUND"
	ErrCodeConflict   = "CONFLICT"
)

type DomainError struct {
	Code    string
	Message string
}

func (e DomainError) Error() string {
	return e.Message
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
