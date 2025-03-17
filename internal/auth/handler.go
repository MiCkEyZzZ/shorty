package auth

import (
	"log"
	"net/http"

	"shorty/internal/config"
	"shorty/internal/service"
	"shorty/pkg/jwt"
	"shorty/pkg/req"
	"shorty/pkg/res"
)

// AuthHandlerDeps - зависимости для обработчика аутентификации.
type AuthHandlerDeps struct {
	*config.Config
	Service *service.AuthService
}

// AuthHandler - обработчик аутентификации.
type AuthHandler struct {
	*config.Config
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
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		body, err := req.HandleBody[SignupRequest](&w, r)
		if err != nil {
			log.Printf("[AuthHandler] Ошибка обработки тела запроса: %v", err)
			http.Error(w, "Некорректный запрос", http.StatusBadRequest)
			return
		}

		email, err := h.Service.Registration(ctx, body.Name, body.Email, body.Password)
		if err != nil {
			log.Printf("[AuthHandler] Ошибка регистрации: %v", err)
			http.Error(w, "Ошибка регистрации пользователя", http.StatusInternalServerError)
			return
		}

		token, err := jwt.NewJWT(h.Config.Auth.Secret).Create(email)
		if err != nil {
			log.Println("[AuthHandler] Ошибка при авторизации:", err)
			http.Error(w, "ошибка при авторизации", http.StatusInternalServerError)
			return
		}
		data := SignupResponse{
			Token: token,
		}
		res.Json(w, data, http.StatusOK)
	}
}

// Signin - вход пользователя.
func (h *AuthHandler) SignIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		body, err := req.HandleBody[SigninRequest](&w, r)
		if err != nil {
			log.Println("[AuthHandler] Ошибка при разборе тела запроса:", err)
			http.Error(w, "Некорректный запрос", http.StatusBadRequest)
			return
		}

		email, err := h.Service.Login(ctx, body.Email, body.Password)
		if err != nil {
			log.Println("[AuthHandler] Ошибка при авторизации:", err)
			http.Error(w, "ошибка при авторизации", http.StatusUnauthorized)
			return
		}
		token, err := jwt.NewJWT(h.Config.Auth.Secret).Create(email)
		if err != nil {
			log.Println("[AuthHandler] Ошибка при авторизации:", err)
			http.Error(w, "ошибка при авторизации", http.StatusInternalServerError)
			return
		}
		data := SinginResponse{
			Token: token,
		}
		res.Json(w, data, http.StatusOK)
	}
}
