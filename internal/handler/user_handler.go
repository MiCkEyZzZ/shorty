package handler

import (
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

// UserHandlerDeps - зависимости для создания экземпляра UserHandler
type UserHandlerDeps struct {
	Config      *config.Config
	UserService *service.UserService
}

// UserHandler - обработчик для управления пользователями.
type UserHandler struct {
	Config      *config.Config
	UserService *service.UserService
}

// NewUserHandler регистрирует маршруты, связанные с пользователями, и привязывает их к методам UserHandler.
func NewUserHandler(router *http.ServeMux, deps UserHandlerDeps) {
	handler := &UserHandler{
		Config:      deps.Config,
		UserService: deps.UserService,
	}

	router.HandleFunc("PATCH /users/{id}", handler.Update())
	router.HandleFunc("DELETE /users/{id}", handler.Delete())
}

// Update обновляет пользователя.
func (h *UserHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			logger.Error("Некорректный ID пользователя", zap.String("id", idStr), zap.Error(err))
			res.ERROR(w, common.ErrInvalidID, http.StatusBadRequest)
			return
		}
		body, err := req.HandleBody[models.User](&w, r)
		if err != nil {
			logger.Error("Ошибка обработки тела запроса", zap.Error(err))
			res.ERROR(w, common.ErrRequestBodyParse, http.StatusBadRequest)
			return
		}
		body.ID = uint(id)
		updatedUser, err := h.UserService.Update(ctx, body)
		if err != nil {
			logger.Error("Ошибка обновления пользователя", zap.Int("id", id), zap.Error(err))
			res.ERROR(w, common.ErrUserUpdateFailed, http.StatusInternalServerError)
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
			res.ERROR(w, common.ErrInvalidID, http.StatusBadRequest)
			return
		}
		err = h.UserService.Delete(ctx, uint(id))
		if err != nil {
			logger.Error("Ошибка удаления пользователя", zap.Int("id", id), zap.Error(err))
			res.ERROR(w, common.ErrUserDeleteFailed, http.StatusInternalServerError)
			return
		}
		logger.Info("Пользователь успешно удалён", zap.Int("id", id))

		res.JSON(w, map[string]string{"message": "Пользователь удалён"}, http.StatusOK)
	}
}
