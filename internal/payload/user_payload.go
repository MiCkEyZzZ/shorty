package payload

import "shorty/internal/models"

// CreateRequest структура представляет запрос на создание нового пользователя.
type CreateUserRequest struct {
	Name      string      `json:"name" validate:"required"`
	Email     string      `json:"email" validate:"required,email"`
	Password  string      `json:"password" validate:"required"`
	Role      models.Role `json:"role" validate:"required"`
	IsBlocked bool        `json:"is_blocked"`
}

// GetByEmailRequest структура представляет запрос на получения пользователя по адресу электронной почты.
type GetUserByEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// CreateResponse структура представляет ответ на запрос создания нового пользователя.
type CreateUserResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// GetByEmailResponse структура представляет ответ на запрос получения пользователя по его адресу электронной почты.
type GetUserByEmailResponse struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
