package auth

import (
	"fmt"
	"log"
	"net/http"

	"shorty/internal/config"
	"shorty/pkg/req"
	"shorty/pkg/res"
)

// AuthHandlerDeps - зависимости для обработчика аутентификации.
type AuthHandlerDeps struct {
	*config.Config
}

// AuthHandler - обработчик аутентификации.
type AuthHandler struct {
	*config.Config
}

// NewAuthHandler - создание обработчика аутентификации.
func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		Config: deps.Config,
	}
	router.HandleFunc("POST /auth/signup", handler.SignUp())
	router.HandleFunc("POST /auth/signin", handler.SignIn())
}

// Signup - регистрация нового пользователя.
func (h *AuthHandler) SignUp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		body, err := req.HandleBody[SignupRequest](&w, r)
		if err != nil {
			log.Println("Ошибка при разборе тела запроса:", err)
			http.Error(w, "Некорректный запрос", http.StatusBadRequest)
			return
		}

		fmt.Println(body)
	}
}

// Signin - вход пользователя.
func (h *AuthHandler) SignIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		body, err := req.HandleBody[SigninRequest](&w, r)
		if err != nil {
			log.Println("Ошибка при разборе тела запроса:", err)
			http.Error(w, "Некорректный запрос", http.StatusBadRequest)
			return
		}

		fmt.Println(body)
		data := SinginResponse{
			Token: "123",
		}
		res.Json(w, data, http.StatusOK)
	}
}
