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
	"shorty/pkg/middleware"
	"shorty/pkg/parse"
	"shorty/pkg/req"
	"shorty/pkg/res"
)

// AdminHandlerDeps - зависимости для создания экземпляра AdminHandler
type AdminHandlerDeps struct {
	Config      *config.Config
	UserService *service.UserService
	LinkService *service.LinkService
	StatService *service.StatService
}

// AdminHandler - обработчик для управления администратором.
type AdminHandler struct {
	Config      *config.Config
	UserService *service.UserService
	LinkService *service.LinkService
	StatService *service.StatService
}

// NewAdminHandler регистрирует маршруты, связанные с администратором, и привязывает их к методам AdminHandler.
func NewAdminHandler(router *http.ServeMux, deps AdminHandlerDeps) {
	handler := &AdminHandler{
		Config:      deps.Config,
		UserService: deps.UserService,
		LinkService: deps.LinkService,
		StatService: deps.StatService,
	}

	// Статистика.
	router.HandleFunc("GET /admin/stats", middleware.AdminOnly(handler.GetStats()))

	// Управление пользователями.
	router.HandleFunc("/admin/users", handler.GetUsers())
	router.HandleFunc("GET /admin/users/{id}", handler.GetUser())
	router.HandleFunc("PATCH /admin/users/{id}", handler.UpdateUser())
	router.HandleFunc("DELETE /admin/users/{id}", handler.DeleteUser())
	router.HandleFunc("PATCH /admin/users/{id}/block", middleware.AdminOnly(handler.BlockUser()))
	router.HandleFunc("PATCH /admin/users/{id}/unblock", middleware.AdminOnly(handler.UnblockUser()))

	// Управление ссылками.
	router.HandleFunc("PATCH /admin/links/{id}/block", handler.BlockLink())
	router.HandleFunc("PATCH /admin/links/{id}/unblock", handler.UnblockLink())
	router.HandleFunc("DELETE /admin/links/{id}", middleware.AdminOnly(handler.DeleteLink()))
}

// GetStats метод для получения статистики.
func (a *AdminHandler) GetStats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Логика для получения статистики
	}
}

// GetUsers метод для получения списка пользователей.
func (a *AdminHandler) GetUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		users, err := a.UserService.GetAll(ctx)
		if err != nil {
			logger.Error("Ошибка получения списка пользователей", zap.Error(err))
			res.ERROR(w, common.ErrorGetUsers, http.StatusInternalServerError)
			return
		}
		logger.Info("Список пользователей успешно получен", zap.Int("count", len(users)))
		res.JSON(w, users, http.StatusOK)

	}
}

// GetUser метод для получения пользователя по идентификатору.
func (a *AdminHandler) GetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := r.PathValue("id")
		userID, err := strconv.Atoi(id)
		if err != nil {
			logger.Error("Некорректный идентификатор пользователя", zap.String("id", id), zap.Error(err))
			res.ERROR(w, common.ErrInvalidID, http.StatusBadRequest)
			return
		}
		user, err := a.UserService.GetByID(ctx, uint(userID))
		if err != nil {
			logger.Error("Ошибка поиска пользователя", zap.Int("userID", userID), zap.Error(err))
			res.ERROR(w, common.ErrUserNotFound, http.StatusNotFound)
			return
		}
		logger.Info("Пользователь найден", zap.Int("id", userID))
		res.JSON(w, user, http.StatusOK)
	}
}

// UpdateUser метод для обновления пользователя по идентификатору.
func (a *AdminHandler) UpdateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := r.PathValue("id")
		userID, err := strconv.Atoi(id)
		if err != nil {
			logger.Error("Некорректный ID пользователя", zap.String("id", id), zap.Error(err))
			res.ERROR(w, common.ErrInvalidID, http.StatusBadRequest)
			return
		}
		body, err := req.HandleBody[models.User](&w, r)
		if err != nil {
			logger.Error("Ошибка обработки тела запроса", zap.Error(err))
			res.ERROR(w, common.ErrRequestBodyParse, http.StatusBadRequest)
			return
		}
		body.ID = uint(userID)
		updatedUser, err := a.UserService.Update(ctx, body)
		if err != nil {
			logger.Error("Ошибка обновления пользователя", zap.Int("userID", userID), zap.Error(err))
			res.ERROR(w, common.ErrUserUpdateFailed, http.StatusInternalServerError)
			return
		}
		logger.Info("Пользователь успешно обновлён", zap.Int("userID", userID))
		res.JSON(w, updatedUser, http.StatusOK)
	}
}

// DeleteUser метод для удаления пользователя по идентификатору.
func (a *AdminHandler) DeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := r.PathValue("id")
		userID, err := strconv.Atoi(id)
		if err != nil {
			logger.Error("Некорректный ID пользователя", zap.String("id", id), zap.Error(err))
			res.ERROR(w, common.ErrInvalidID, http.StatusBadRequest)
			return
		}
		err = a.UserService.Delete(ctx, uint(userID))
		if err != nil {
			logger.Error("Ошибка удаления пользователя", zap.Int("userID", userID), zap.Error(err))
			res.ERROR(w, common.ErrUserDeleteFailed, http.StatusInternalServerError)
			return
		}
		logger.Info("Пользователь успешно удалён", zap.Int("userID", userID))
		res.JSON(w, map[string]string{"message": "Пользователь удалён"}, http.StatusOK)
	}
}

// BlockUser метод для блокировки пользователя по идентификатору.
func (a *AdminHandler) BlockUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Логика для получения статистики
	}
}

// UnblockUser метод для разблокировки пользователя по идентификатору.
func (a *AdminHandler) UnblockUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Логика для получения статистики
	}
}

func (a *AdminHandler) BlockLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, err := parse.ParseID(r)
		if err != nil {
			logger.Error("Ошибка парсинга ID ссылки", zap.Error(err))
			res.ERROR(w, common.ErrInvalidID, http.StatusBadRequest)
			return
		}

		// Получаем ссылку из базы
		link, err := a.LinkService.FindByID(ctx, uint(id))
		if err != nil {
			logger.Error("Ошибка поиска ссылки", zap.Uint("id", uint(id)), zap.Error(err))
			res.ERROR(w, common.ErrNotFound, http.StatusNotFound)
			return
		}

		// Блокируем ссылку
		updatedLink, err := a.LinkService.Block(ctx, link.ID)
		if err != nil {
			logger.Error("Ошибка блокировки ссылки", zap.Uint("id", uint(id)), zap.Error(err))
			res.ERROR(w, common.ErrLinkBlockFailed, http.StatusInternalServerError)
			return
		}

		logger.Info("Ссылка успешно заблокирована", zap.Uint("id", updatedLink.ID))
		res.JSON(w, updatedLink, http.StatusOK)
	}
}

func (a *AdminHandler) UnblockLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, err := parse.ParseID(r)
		if err != nil {
			logger.Error("Ошибка парсинга ID ссылки", zap.Error(err))
			res.ERROR(w, common.ErrInvalidID, http.StatusBadRequest)
			return
		}

		// Получаем ссылку из базы
		link, err := a.LinkService.FindByID(ctx, uint(id))
		if err != nil {
			logger.Error("Ошибка поиска ссылки", zap.Uint("id", uint(id)), zap.Error(err))
			res.ERROR(w, common.ErrNotFound, http.StatusNotFound)
			return
		}

		// Блокируем ссылку
		updatedLink, err := a.LinkService.UnBlock(ctx, link.ID)
		if err != nil {
			logger.Error("Ошибка при разблокировки ссылки", zap.Uint("id", uint(id)), zap.Error(err))
			res.ERROR(w, common.ErrUnBlockFailed, http.StatusInternalServerError)
			return
		}

		logger.Info("Ссылка успешно разблокирована", zap.Uint("id", updatedLink.ID))
		res.JSON(w, updatedLink, http.StatusOK)
	}
}

func (a *AdminHandler) DeleteLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Логика для получения статистики
	}
}
