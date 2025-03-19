package handler

import (
	"errors"
	"net/http"
	"strconv"

	"go.uber.org/zap"

	"shorty/internal/common"
	"shorty/internal/config"
	"shorty/internal/models"
	"shorty/internal/service"
	"shorty/pkg/logger"
	"shorty/pkg/req"
	"shorty/pkg/res"
)

var (
	ErrorGetUsers       = errors.New("не удалось получить список пользователей")
	ErrWrongID          = errors.New("некорректный ID пользователя")
	ErrUserNotFound     = errors.New("пользователь не найден")
	ErrUserUpdateFailed = errors.New("не удалось обновить пользователя")
	ErrUserDeleteFailed = errors.New("не удалось удалить пользователя")
)

// UserHandlerDeps - зависимости для создания экземпляра UserHandler
type UserHandlerDeps struct {
	Config  *config.Config
	Service *service.UserService
}

// UserHandler - обработчик для управления пользователями.
type UserHandler struct {
	Config  *config.Config
	Service *service.UserService
}

// NewUserHandler регистрирует маршруты, связанные с пользователями, и привязывает их к методам UserHandler.
func NewUserHandler(router *http.ServeMux, deps UserHandlerDeps) {
	handler := &UserHandler{
		Config:  deps.Config,
		Service: deps.Service,
	}

	router.HandleFunc("GET /users", handler.FindAll())
	router.HandleFunc("GET /users/{id}", handler.FindByID())
	router.HandleFunc("PATCH /users/{id}", handler.Update())
	router.HandleFunc("DELETE /users/{id}", handler.Delete())
}

// FindAll получает список пользователей.
func (h *UserHandler) FindAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		users, err := h.Service.GetAll(ctx)
		if err != nil {
			logger.Error("Ошибка получения списка пользователей", zap.Error(err))
			res.ERROR(w, ErrorGetUsers, http.StatusInternalServerError)
			return
		}
		logger.Info("Список пользователей успешно получен", zap.Int("count", len(users)))
		res.JSON(w, users, http.StatusOK)
	}
}

// FindByID получает пользователя по идентификатору.
func (h *UserHandler) FindByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			logger.Error("Некорректный ID пользователя", zap.String("id", idStr), zap.Error(err))
			res.ERROR(w, ErrWrongID, http.StatusBadRequest)
			return
		}
		user, err := h.Service.GetByID(ctx, uint(id))
		if err != nil {
			logger.Error("Ошибка поиска пользователя", zap.Int("id", id), zap.Error(err))
			res.ERROR(w, ErrUserNotFound, http.StatusNotFound)
			return
		}
		logger.Info("Пользователь найден", zap.Int("id", id))

		res.JSON(w, user, http.StatusOK)
	}
}

// Update обновляет пользователя.
func (h *UserHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			logger.Error("Некорректный ID пользователя", zap.String("id", idStr), zap.Error(err))
			res.ERROR(w, ErrWrongID, http.StatusBadRequest)
			return
		}
		body, err := req.HandleBody[models.User](&w, r)
		if err != nil {
			logger.Error("Ошибка обработки тела запроса", zap.Error(err))
			res.ERROR(w, common.ErrRequestBodyParse, http.StatusBadRequest)
			return
		}
		body.ID = uint(id)
		updatedUser, err := h.Service.Update(ctx, body)
		if err != nil {
			logger.Error("Ошибка обновления пользователя", zap.Int("id", id), zap.Error(err))
			res.ERROR(w, ErrUserUpdateFailed, http.StatusInternalServerError)
			return
		}
		logger.Info("Пользователь успешно обновлён", zap.Int("id", id))

		res.JSON(w, updatedUser, http.StatusOK)
	}
}

// Delete удаляет пользователя по идентификатору.
func (h *UserHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			logger.Error("Некорректный ID пользователя", zap.String("id", idStr), zap.Error(err))
			res.ERROR(w, ErrWrongID, http.StatusBadRequest)
			return
		}
		err = h.Service.Delete(ctx, uint(id))
		if err != nil {
			logger.Error("Ошибка удаления пользователя", zap.Int("id", id), zap.Error(err))
			res.ERROR(w, ErrUserDeleteFailed, http.StatusInternalServerError)
			return
		}
		logger.Info("Пользователь успешно удалён", zap.Int("id", id))

		res.JSON(w, map[string]string{"message": "Пользователь удалён"}, http.StatusOK)
	}
}
