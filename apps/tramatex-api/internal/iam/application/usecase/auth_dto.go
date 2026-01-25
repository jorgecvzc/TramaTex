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
