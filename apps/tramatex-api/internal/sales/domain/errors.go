package domain

const (
	ErrCodeValidation    = "VALIDATION_ERROR"
	ErrCodeNotFound      = "NOT_FOUND"
	ErrCodeConflict      = "CONFLICT"
	ErrCodeConfiguration = "CONFIGURATION_ERROR"
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
		return 400
	case ErrCodeNotFound:
		return 404
	case ErrCodeConflict:
		return 409
	case ErrCodeConfiguration:
		return 500
	default:
		return 400
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

func NewConfigurationError(message string) error {
	return DomainError{Code: ErrCodeConfiguration, Message: message}
}
