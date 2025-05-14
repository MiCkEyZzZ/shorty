package payload

import "shorty/internal/models"

// SigninRequest represents a user sign-in request.
type SigninRequest struct {
	Email    string      `json:"email" validate:"required,email"`
	Password string      `json:"password" validate:"required"`
	Role     models.Role `json:"role" validate:"required"`
}

// SignupRequest is an alias for the request to create a new user.
type SignupRequest = CreateUserRequest

// SigninResponse represents the response to a user sign-in request.
type SinginResponse struct {
	Token string `json:"token"`
}

// SignupResponse represents the response to a user registration request.
type SignupResponse struct {
	Token string `json:"token"`
}
