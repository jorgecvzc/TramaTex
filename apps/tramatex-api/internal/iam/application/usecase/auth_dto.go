package usecase

// LoginInput contains the input data for the login use case.
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserDTO represents a user in the output DTOs.
type UserDTO struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// LoginOutput contains the output data from the login use case.
// Returns user info + tokens.
type LoginOutput struct {
	User         UserDTO `json:"user"`
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token"`
	ExpiresIn    int     `json:"expires_in"` // Seconds until token expiry
}

// RegisterInput contains input data for user registration.
type RegisterInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RegisterOutput contains output data from user registration.
type RegisterOutput struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// RefreshInput contains input data for token refresh.
type RefreshInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshOutput contains output data for token refresh.
type RefreshOutput struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"` // Seconds until token expiry
}

// LogoutInput contains optional data for logout.
type LogoutInput struct {
	RefreshToken string `json:"refresh_token,omitempty"`
}

// AssignRoleInput contains input data to assign a role to a user.
type AssignRoleInput struct {
	UserID string `json:"user_id" binding:"required"`
	Role   string `json:"role" binding:"required"`
}

// AssignRoleOutput contains output data for role assignment.
type AssignRoleOutput struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

// CheckAuthorizationInput contains data for authorization checks.
type CheckAuthorizationInput struct {
	UserID        string   `json:"user_id,omitempty"`
	RequiredRoles []string `json:"required_roles" binding:"required"`
}

// CheckAuthorizationOutput contains result of authorization check.
type CheckAuthorizationOutput struct {
	Allowed bool   `json:"allowed"`
	Role    string `json:"role"`
}
