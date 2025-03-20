package handler

import (
	"net/http"
	"strconv"

	"shorty/internal/common"
	"shorty/internal/config"
	"shorty/internal/service"
	"shorty/pkg/logger"
	"shorty/pkg/middleware"
	"shorty/pkg/res"

	"go.uber.org/zap"
)

type AdminHandlerDeps struct {
	Config      *config.Config
	UserService *service.UserService
	LinkService *service.LinkService
	StatService *service.StatService
}

type AdminHandler struct {
	Config      *config.Config
	UserService *service.UserService
	LinkService *service.LinkService
	StatService *service.StatService
}

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
	router.HandleFunc("PATCH /admin/users/{id}", middleware.AdminOnly(handler.UpdateUser()))
	router.HandleFunc("DELETE /admin/users/{id}", middleware.AdminOnly(handler.DeleteUser()))
	router.HandleFunc("PATCH /admin/users/{id}/block", middleware.AdminOnly(handler.BlockUser()))
	router.HandleFunc("PATCH /admin/users/{id}/unblock", middleware.AdminOnly(handler.UnblockUser()))

	// Управление ссылками.
	router.HandleFunc("PATCH /admin/links/{id}/block", handler.BlockLink())
	router.HandleFunc("PATCH /admin/links/{id}/unblock", handler.UnblockLink())
	router.HandleFunc("DELETE /admin/links/{id}", middleware.AdminOnly(handler.DeleteLink()))
}

func (a *AdminHandler) GetStats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Логика для получения статистики
	}
}

func (a *AdminHandler) GetUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		users, err := a.UserService.GetAll(ctx)
		if err != nil {
			logger.Error("Ошибка получения списка пользователей", zap.Error(err))
			res.ERROR(w, ErrorGetUsers, http.StatusInternalServerError)
			return
		}
		logger.Info("Список пользователей успешно получен", zap.Int("count", len(users)))
		res.JSON(w, users, http.StatusOK)

	}
}

func (a *AdminHandler) GetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := r.PathValue("id")
		userID, err := strconv.Atoi(id)
		if err != nil {
			logger.Error("Некорректный идентификатор пользователя", zap.String("id", id), zap.Error(err))
			res.ERROR(w, ErrWrongID, http.StatusBadRequest)
			return
		}
		user, err := a.UserService.GetByID(ctx, uint(userID))
		if err != nil {
			logger.Error("Ошибка поиска пользователя", zap.Int("userID", userID), zap.Error(err))
			res.ERROR(w, ErrUserNotFound, http.StatusNotFound)
			return
		}
		logger.Info("Пользователь найден", zap.Int("id", userID))
		res.JSON(w, user, http.StatusOK)
	}
}

func (a *AdminHandler) UpdateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Логика для получения статистики
	}
}

func (a *AdminHandler) DeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Логика для получения статистики
	}
}

func (a *AdminHandler) BlockUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Логика для получения статистики
	}
}

func (a *AdminHandler) UnblockUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Логика для получения статистики
	}
}

func (a *AdminHandler) BlockLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, err := parseID(r)
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
			res.ERROR(w, common.ErrUpdateFailed, http.StatusInternalServerError)
			return
		}

		logger.Info("Ссылка успешно заблокирована", zap.Uint("id", updatedLink.ID))
		res.JSON(w, updatedLink, http.StatusOK)
	}
}

func (a *AdminHandler) UnblockLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, err := parseID(r)
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
			res.ERROR(w, common.ErrUpdateFailed, http.StatusInternalServerError)
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

// parseID парсит идентификатор из строки в uint.
func parseIDs(r *http.Request) (uint, error) {
	rid := r.PathValue("id")
	id, err := strconv.ParseUint(rid, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
