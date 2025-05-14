package payload

import "shorty/internal/models"

// CreateUserRequest represents the request payload for creating a new user.
type CreateUserRequest struct {
	Name      string      `json:"name" validate:"required"`
	Email     string      `json:"email" validate:"required,email"`
	Password  string      `json:"password" validate:"required"`
	Role      models.Role `json:"role" validate:"required"`
	IsBlocked bool        `json:"is_blocked"`
}

// GetUserByEmailRequest represents the request payload for retrieving a user by email address.
type GetUserByEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// CreateUserResponse represents the response payload for a successful user creation request.
type CreateUserResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// GetUserByEmailResponse represents the response payload for retrieving a user by email address.
type GetUserByEmailResponse struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
