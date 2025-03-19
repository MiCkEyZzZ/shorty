package handler

import (
	"errors"
	"log"
	"net/http"

	"shorty/internal/config"
	"shorty/internal/payload"
	"shorty/internal/service"
	"shorty/pkg/jwt"
	"shorty/pkg/req"
	"shorty/pkg/res"
)

var (
	ErrMethodNotAllowed       = errors.New("метод не поддерживается")
	ErrBadRequest             = errors.New("некорректный запрос")
	ErrUserRegistrationFailed = errors.New("ошибка регистрации пользователя")
	ErrAuthFailed             = errors.New("ошибка при авторизации")
)

// AuthHandlerDeps - зависимости для обработчика аутентификации.
type AuthHandlerDeps struct {
	Config  *config.Config
	Service *service.AuthService
}

// AuthHandler - обработчик аутентификации.
type AuthHandler struct {
	Config  *config.Config
	Service *service.AuthService
}

// NewAuthHandler - создание обработчика аутентификации.
func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		Config:  deps.Config,
		Service: deps.Service,
	}
	router.HandleFunc("POST /auth/signup", handler.SignUp())
	router.HandleFunc("POST /auth/signin", handler.SignIn())
}

// Signup - регистрация нового пользователя.
func (h *AuthHandler) SignUp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if r.Method != http.MethodPost {
			res.ERROR(w, ErrMethodNotAllowed, http.StatusMethodNotAllowed)
			return
		}

		body, err := req.HandleBody[payload.SignupRequest](&w, r)
		if err != nil {
			log.Printf("[AuthHandler] Ошибка обработки тела запроса: %v", err)
			res.ERROR(w, ErrBadRequest, http.StatusBadRequest)
			return
		}

		email, err := h.Service.Registration(ctx, body.Name, body.Email, body.Password)
		if err != nil {
			log.Printf("[AuthHandler] Ошибка регистрации: %v", err)
			res.ERROR(w, ErrUserRegistrationFailed, http.StatusInternalServerError)
			return
		}

		token, err := jwt.NewJWT(h.Config.Auth.Secret).Create(jwt.JWTData{Email: email})
		if err != nil {
			log.Println("[AuthHandler] Ошибка при авторизации:", err)
			res.ERROR(w, ErrAuthFailed, http.StatusInternalServerError)
			return
		}
		data := payload.SignupResponse{
			Token: token,
		}
		res.JSON(w, data, http.StatusOK)
	}
}

// Signin - вход пользователя.
func (h *AuthHandler) SignIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if r.Method != http.MethodPost {
			res.ERROR(w, ErrMethodNotAllowed, http.StatusMethodNotAllowed)
			return
		}

		body, err := req.HandleBody[payload.SigninRequest](&w, r)
		if err != nil {
			log.Println("[AuthHandler] Ошибка при разборе тела запроса:", err)
			res.ERROR(w, ErrBadRequest, http.StatusBadRequest)
			return
		}

		email, err := h.Service.Login(ctx, body.Email, body.Password)
		if err != nil {
			log.Println("[AuthHandler] Ошибка при авторизации:", err)
			res.ERROR(w, ErrAuthFailed, http.StatusInternalServerError)
			return
		}
		token, err := jwt.NewJWT(h.Config.Auth.Secret).Create(jwt.JWTData{Email: email})
		if err != nil {
			log.Println("[AuthHandler] Ошибка при авторизации:", err)
			res.ERROR(w, ErrAuthFailed, http.StatusInternalServerError)
			return
		}
		data := payload.SinginResponse{
			Token: token,
		}
		res.JSON(w, data, http.StatusOK)
	}
}
