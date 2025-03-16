package user

import "time"

// CreateRequest структура представляет запрос на создание нового пользователя.
type CreateRequest struct {
	Name      string    `json:"name" validate:"reuired,name"`
	Email     string    `json:"email" validate:"required,email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetByEmailRequest структура представляет запрос на получения пользователя по адресу электронной почты.
type GetByEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// CreateResponse структура представляет ответ на запрос создания нового пользователя.
type CreateResponse struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

// GetByEmailResponse структура представляет ответ на запрос получения пользователя по его адресу электронной почты.
type GetByEmailResponse struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
