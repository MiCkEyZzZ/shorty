package handler

import (
	"net/http"

	"go.uber.org/zap"

	"shorty/internal/common"
	"shorty/internal/config"
	"shorty/internal/payload"
	"shorty/internal/service"
	"shorty/pkg/jwt"
	"shorty/pkg/logger"
	"shorty/pkg/req"
	"shorty/pkg/res"
)

// AuthHandlerDeps - зависимости для обработчика аутентификации.
type AuthHandlerDeps struct {
	Config      *config.Config
	AuthService service.AuthServ
}

// AuthHandler - обработчик аутентификации.
type AuthHandler struct {
	Config      *config.Config
	AuthService service.AuthServ
}

// NewAuthHandler - создание обработчика аутентификации.
func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		Config:      deps.Config,
		AuthService: deps.AuthService,
	}

	// Управление авторизацией.
	router.HandleFunc("POST /auth/signup", handler.SignUp())
	router.HandleFunc("POST /auth/signin", handler.SignIn())
}

// Signup - регистрация нового пользователя.
func (h *AuthHandler) SignUp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if r.Method != http.MethodPost {
			logger.Error("Неверный метод запроса", zap.String("method", r.Method))
			res.ERROR(w, common.ErrMethodNotAllowed, http.StatusMethodNotAllowed)
			return
		}

		body, err := req.HandleBody[payload.SignupRequest](&w, r)
		if err != nil {
			logger.Error("Ошибка при парсинге тела запроса для регистрации", zap.Error(err))
			res.ERROR(w, common.ErrBadRequest, http.StatusBadRequest)
			return
		}

		user, err := h.AuthService.Registration(ctx, body.Name, body.Email, body.Password, body.Role, body.IsBlocked)
		if err != nil {
			logger.Error("Ошибка регистрации пользователя", zap.String("email", body.Email), zap.Error(err))
			res.ERROR(w, common.ErrUserRegistrationFailed, http.StatusInternalServerError)
			return
		}

		token, err := jwt.NewJWT(h.Config.Auth.Secret).CreateToken(user)
		if err != nil {
			logger.Error("Ошибка при создании токена", zap.String("email", user.Email), zap.Error(err))
			res.ERROR(w, common.ErrAuthFailed, http.StatusInternalServerError)
			return
		}

		data := payload.SignupResponse{
			Token: token,
		}
		logger.Info("Пользователь успешно зарегистрирован", zap.String("email", user.Email))
		res.JSON(w, data, http.StatusOK)
	}
}

// Signin - вход пользователя.
func (h *AuthHandler) SignIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if r.Method != http.MethodPost {
			logger.Error("Неверный метод запроса", zap.String("method", r.Method))
			res.ERROR(w, common.ErrMethodNotAllowed, http.StatusMethodNotAllowed)
			return
		}

		body, err := req.HandleBody[payload.SigninRequest](&w, r)
		if err != nil {
			logger.Error("Ошибка при парсинге тела запроса для авторизации", zap.Error(err))
			res.ERROR(w, common.ErrBadRequest, http.StatusBadRequest)
			return
		}

		user, err := h.AuthService.Login(ctx, body.Email, body.Password, body.Role)
		if err != nil {
			logger.Error("Ошибка авторизации пользователя", zap.String("email", body.Email), zap.Error(err))
			res.ERROR(w, common.ErrAuthFailed, http.StatusInternalServerError)
			return
		}

		token, err := jwt.NewJWT(h.Config.Auth.Secret).CreateToken(user)
		if err != nil {
			logger.Error("Ошибка при создании токена для авторизованного пользователя", zap.String("email", user.Email), zap.Error(err))
			res.ERROR(w, common.ErrAuthFailed, http.StatusInternalServerError)
			return
		}

		data := payload.SinginResponse{
			Token: token,
		}
		logger.Info("Пользователь успешно авторизован", zap.String("email", body.Email))
		res.JSON(w, data, http.StatusOK)
	}
}
