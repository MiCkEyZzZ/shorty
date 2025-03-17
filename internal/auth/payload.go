package auth

import "shorty/internal/user"

// SigninRequest - представляет запрос на вход пользователя.
type SigninRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// SignupRequest - представляет запрос на регистрацию нового пользователя.
type SignupRequest = user.CreateUserRequest

// SinginResponse - представляет ответ на запрос на вход пользователя.
type SinginResponse struct {
	Token string `json:"token"`
}

// SignupResponse - представляет ответ на запрос на регистрацию нового пользователя.
type SignupResponse struct {
	Token string `json:"token"`
}
