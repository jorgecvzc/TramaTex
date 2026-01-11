package errors

import "fmt"

// Domain errors
var (
	ErrInvalidEmail      = fmt.Errorf("invalid email format")
	ErrInvalidPassword   = fmt.Errorf("invalid password")
	ErrUserNotFound      = fmt.Errorf("user not found")
	ErrInvalidCredentials = fmt.Errorf("invalid credentials")
	ErrInvalidInput      = fmt.Errorf("invalid input")
	ErrUnauthorized      = fmt.Errorf("unauthorized")
	ErrForbidden         = fmt.Errorf("forbidden")
	ErrConflict          = fmt.Errorf("resource already exists")
	ErrInternal          = fmt.Errorf("internal server error")
)
