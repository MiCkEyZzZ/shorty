package handler

import (
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"

	"shorty/internal/common"
	"shorty/internal/config"
	"shorty/internal/models"
	"shorty/internal/service"
	"shorty/pkg/jwt"
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
	JWTService  *jwt.JWT
}

// AdminHandler - обработчик для управления администратором.
type AdminHandler struct {
	Config      *config.Config
	UserService *service.UserService
	LinkService *service.LinkService
	StatService *service.StatService
	JWTService  *jwt.JWT
}

// NewAdminHandler регистрирует маршруты, связанные с администратором, и привязывает их к методам AdminHandler.
func NewAdminHandler(router *http.ServeMux, deps AdminHandlerDeps) {
	handler := &AdminHandler{
		Config:      deps.Config,
		UserService: deps.UserService,
		LinkService: deps.LinkService,
		StatService: deps.StatService,
		JWTService:  deps.JWTService,
	}

	adminMiddleware := middleware.AdminMiddleware(deps.JWTService, deps.UserService)

	// Управление пользователями.
	router.Handle("GET /admin/users", adminMiddleware(handler.GetUsers()))
	router.Handle("GET /admin/users/{id}", adminMiddleware(handler.GetUser()))
	router.Handle("PATCH /admin/users/{id}", adminMiddleware(handler.UpdateUser()))
	router.Handle("DELETE /admin/users/{id}", adminMiddleware(handler.DeleteUser()))
	router.Handle("PATCH /admin/users/{id}/block", adminMiddleware(handler.BlockUser()))
	router.Handle("PATCH /admin/users/{id}/unblock", adminMiddleware(handler.UnblockUser()))
	router.Handle("GET /admin/users/blocked/count", adminMiddleware(handler.GetBlockedUsersCount()))

	// Управление ссылками.
	router.Handle("PATCH /admin/links/{id}/block", adminMiddleware(handler.BlockLink()))
	router.Handle("PATCH /admin/links/{id}/unblock", adminMiddleware(handler.UnblockLink()))
	router.Handle("DELETE /admin/links/{id}", adminMiddleware(handler.DeleteLink()))
	router.Handle("GET /admin/links/blocked/count", adminMiddleware(handler.GetBlockedLinksCount()))

	// Управление статистикой
	router.HandleFunc("GET /admin/stats", handler.GetStats())
}

// GetStats метод для получения статистики.
func (h *AdminHandler) GetStats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		fromStr := r.URL.Query().Get("from")
		from, err := time.Parse("2006-01-02", fromStr)
		if err != nil {
			logger.Error("Ошибка парсинга параметра 'from'", zap.String("from", fromStr), zap.Error(err))
			res.ERROR(w, common.ErrInvalidParam, http.StatusBadRequest)
			return
		}
		toStr := r.URL.Query().Get("to")
		to, err := time.Parse("2006-01-02", toStr)
		if err != nil {
			logger.Error("Ошибка парсинга параметра 'to'", zap.String("to", toStr), zap.Error(err))
			res.ERROR(w, common.ErrInvalidParam, http.StatusBadRequest)
			return
		}
		by := r.URL.Query().Get("by")
		if by != common.GroupByDay && by != common.GroupByMonth {
			logger.Error("Неверное значение параметра 'by'", zap.String("by", by))
			res.ERROR(w, common.ErrInvalidParam, http.StatusBadRequest)
			return
		}
		logger.Info("Получение статистики", zap.String("by", by), zap.Time("from", from), zap.Time("to", to))
		stats := h.StatService.GetStats(ctx, by, from, to)
		logger.Info("Статистика успешно получена", zap.Int("record_count", len(stats)))
		res.JSON(w, stats, http.StatusOK)
	}
}

// GetUsers метод для получения списка пользователей.
func (a *AdminHandler) GetUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		users, err := a.UserService.GetAll(ctx)
		if err != nil {
			logger.Error("Error getting list of users", zap.Error(err))
			res.ERROR(w, common.ErrorGetUsers, http.StatusInternalServerError)
			return
		}
		logger.Info("User list successfully received.", zap.Int("count", len(users)))
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
			logger.Error("Ошибка при поиска пользователя", zap.Int("userID", userID), zap.Error(err))
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
		ctx := r.Context()
		id, err := parse.ParseID(r)
		if err != nil {
			logger.Error("Ошибка парсинга идентификатора пользователя", zap.Error(err))
			res.ERROR(w, common.ErrInvalidID, http.StatusBadRequest)
			return
		}
		user, err := a.UserService.GetByID(ctx, uint(id))
		if err != nil {
			logger.Error("Ошибка при поиске пользователя", zap.Uint("id", uint(id)), zap.Error(err))
			res.ERROR(w, common.ErrNotFound, http.StatusNotFound)
			return
		}

		updateUser, err := a.UserService.Block(ctx, user.ID)
		if err != nil {
			logger.Error("Ошибка при блокировке пользователя", zap.Uint("id", uint(id)), zap.Error(err))
			res.ERROR(w, common.ErrLinkBlockFailed, http.StatusInternalServerError)
			return
		}

		logger.Info("Пользователь успешно заблокирована", zap.Uint("id", updateUser.ID))
		res.JSON(w, updateUser, http.StatusOK)
	}
}

// UnblockUser метод для разблокировки пользователя по идентификатору.
func (a *AdminHandler) UnblockUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, err := parse.ParseID(r)
		if err != nil {
			logger.Error("Ошибка парсинга идентификатора пользователя", zap.Error(err))
			res.ERROR(w, common.ErrInvalidID, http.StatusBadRequest)
			return
		}

		user, err := a.UserService.GetByID(ctx, uint(id))
		if err != nil {
			logger.Error("Ошибка при поиске пользователя", zap.Uint("id", uint(id)), zap.Error(err))
			res.ERROR(w, common.ErrNotFound, http.StatusNotFound)
			return
		}
		updatedUser, err := a.UserService.UnBlock(ctx, user.ID)
		if err != nil {
			logger.Error("Ошибка при разблокировке пользователя", zap.Uint("id", uint(id)), zap.Error(err))
			res.ERROR(w, common.ErrUnBlockFailed, http.StatusInternalServerError)
			return
		}

		logger.Info("Пользователь успешно разблокирован", zap.Uint("id", updatedUser.ID))
		res.JSON(w, updatedUser, http.StatusOK)
	}
}

// BlockLink метод для блокировки ссылки.
func (a *AdminHandler) BlockLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, err := parse.ParseID(r)
		if err != nil {
			logger.Error("Ошибка парсинга идентификатора ссылки", zap.Error(err))
			res.ERROR(w, common.ErrInvalidID, http.StatusBadRequest)
			return
		}

		// Получаем ссылку из базы
		link, err := a.LinkService.FindByID(ctx, uint(id))
		if err != nil {
			logger.Error("Ошибка при поиске ссылки", zap.Uint("id", uint(id)), zap.Error(err))
			res.ERROR(w, common.ErrNotFound, http.StatusNotFound)
			return
		}

		// Блокируем ссылку
		updatedLink, err := a.LinkService.Block(ctx, link.ID)
		if err != nil {
			logger.Error("Ошибка при блокировке ссылки", zap.Uint("id", uint(id)), zap.Error(err))
			res.ERROR(w, common.ErrLinkBlockFailed, http.StatusInternalServerError)
			return
		}

		logger.Info("Ссылка успешно заблокирована", zap.Uint("id", updatedLink.ID))
		res.JSON(w, updatedLink, http.StatusOK)
	}
}

// UnblockLink метод для разблокировки ссылки.
func (a *AdminHandler) UnblockLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, err := parse.ParseID(r)
		if err != nil {
			logger.Error("Ошибка парсинга идентификатора ссылки", zap.Error(err))
			res.ERROR(w, common.ErrInvalidID, http.StatusBadRequest)
			return
		}

		// Получаем ссылку из базы
		link, err := a.LinkService.FindByID(ctx, uint(id))
		if err != nil {
			logger.Error("Ошибка при поиске ссылки", zap.Uint("id", uint(id)), zap.Error(err))
			res.ERROR(w, common.ErrNotFound, http.StatusNotFound)
			return
		}

		// Блокируем ссылку
		updatedLink, err := a.LinkService.UnBlock(ctx, link.ID)
		if err != nil {
			logger.Error("Ошибка при разблокировке ссылки", zap.Uint("id", uint(id)), zap.Error(err))
			res.ERROR(w, common.ErrUnBlockFailed, http.StatusInternalServerError)
			return
		}

		logger.Info("Ссылка успешно разблокирована", zap.Uint("id", updatedLink.ID))
		res.JSON(w, updatedLink, http.StatusOK)
	}
}

func (a *AdminHandler) DeleteLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, err := parse.ParseID(r)
		if err != nil {
			logger.Error("Неверный идентификатор для удаления ссылки", zap.Error(err))
			res.ERROR(w, common.ErrInvalidID, http.StatusBadRequest)
			return
		}

		_, err = a.LinkService.FindByID(ctx, uint(id))
		if err != nil {
			logger.Error("Ссылка не найдена для удаления", zap.Uint("id", uint(id)), zap.Error(err))
			res.ERROR(w, common.ErrLinkNotFound, http.StatusNotFound)
			return
		}

		err = a.LinkService.Delete(ctx, uint(id))
		if err != nil {
			logger.Error("Ошибка при удалении ссылки", zap.Uint("id", uint(id)), zap.Error(err))
			res.ERROR(w, common.ErrLinkDeleteFailed, http.StatusInternalServerError)
			return
		}

		logger.Info("Ссылка успешно удалена", zap.Uint("id", uint(id)))
		res.JSON(w, map[string]string{"message": "ссылка удалена"}, http.StatusOK)
	}
}

// GetBlockedUsersCount метод для получения количества заблокированных пользователей.
func (a *AdminHandler) GetBlockedUsersCount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		count, err := a.UserService.GetBlockedUsersCount(ctx)
		if err != nil {
			logger.Error("Ошибка при получении количества заблокированных пользователей", zap.Error(err))
			res.ERROR(w, common.ErrInternal, http.StatusInternalServerError)
			return
		}
		res.JSON(w, map[string]int64{"blocked_users": count}, http.StatusOK)
	}
}

// GetBlockedLinksCount метод для получения количества заблокированных ссылок.
func (a *AdminHandler) GetBlockedLinksCount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		count, err := a.LinkService.GetBlockedLinksCount(ctx)
		if err != nil {
			if err != nil {
				logger.Error("Ошибка при получении количества заблокированных ссылок", zap.Error(err))
				res.ERROR(w, common.ErrInternal, http.StatusInternalServerError)
				return
			}
		}
		res.JSON(w, map[string]int64{"blocked_links": count}, http.StatusOK)
	}
}
