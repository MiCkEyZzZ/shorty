package handler

import (
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"shorty/internal/common"
	"shorty/internal/config"
	"shorty/internal/models"
	"shorty/internal/payload"
	"shorty/internal/service"
	"shorty/pkg/event"
	"shorty/pkg/logger"
	"shorty/pkg/middleware"
	"shorty/pkg/parse"
	"shorty/pkg/req"
	"shorty/pkg/res"
)

// UserHandlerDeps - зависимости для создания экземпляра UserHandler
type UserHandlerDeps struct {
	Config      *config.Config
	UserService *service.UserService
	LinkService *service.LinkService
	EventBus    *event.EventBus
}

// UserHandler - обработчик для управления пользователями.
type UserHandler struct {
	Config      *config.Config
	UserService *service.UserService
	LinkService *service.LinkService
	EventBus    *event.EventBus
}

// NewUserHandler регистрирует маршруты, связанные с пользователями, и привязывает их к методам UserHandler.
func NewUserHandler(router *http.ServeMux, deps UserHandlerDeps) {
	handler := &UserHandler{
		Config:      deps.Config,
		UserService: deps.UserService,
		LinkService: deps.LinkService,
		EventBus:    deps.EventBus,
	}

	// Управление пользователями.
	router.Handle("PATCH /users/{id}", middleware.IsAuth(handler.Update(), deps.Config))
	router.Handle("DELETE /users/{id}", middleware.IsAuth(handler.Delete(), deps.Config))

	// Управление ссылками.
	router.HandleFunc("POST /users/links", handler.CreateLink())
	router.Handle("GET /users/links", middleware.IsAuth(handler.GetLinks(), deps.Config))
	router.Handle("PATCH /users/links/{id}", middleware.IsAuth(handler.UpdateLink(), deps.Config))
	router.Handle("DELETE /users/links/{id}", middleware.IsAuth(handler.DeleteLink(), deps.Config))
	router.Handle("GET /users/links/{hash}", middleware.IsAuth(handler.Redirect(), deps.Config))
}

// Update метод для обновления пользователя.
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

// Delete метод для удаления пользователя по идентификатору.
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

// CreateLink метод для создания новой ссылку.
func (h *UserHandler) CreateLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		body, err := req.HandleBody[payload.CreateLinkRequest](&w, r)
		if err != nil {
			logger.Error("Ошибка парсинга тела запроса для создания ссылки", zap.Error(err))
			res.ERROR(w, common.ErrRequestBodyParse, http.StatusBadRequest)
			return
		}
		link := models.NewLink(body.URL)
		newLink, err := h.LinkService.Create(ctx, link)
		if err != nil {
			logger.Error("Ошибка создания сокращённого URL", zap.Error(err))
			res.ERROR(w, common.ErrLinkCreateUR, http.StatusBadRequest)
			return
		}
		logger.Info("Сокращённый URL успешно создан", zap.String("short_url", newLink.Url))
		res.JSON(w, newLink, http.StatusOK)
	}
}

// GetLinks метод для получения списка ссылок.
func (h *UserHandler) GetLinks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			logger.Error("Неверный параметр 'limit'", zap.String("limit", r.URL.Query().Get("limit")), zap.Error(err))
			res.ERROR(w, common.ErrInvalidLimit, http.StatusBadRequest)
			return
		}
		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			logger.Error("Неверный параметр 'offset'", zap.String("offset", r.URL.Query().Get("offset")), zap.Error(err))
			res.ERROR(w, common.ErrInvalidOffset, http.StatusBadRequest)
			return
		}
		count, _ := h.LinkService.Count(ctx)
		links, _ := h.LinkService.GetAll(ctx, limit, offset)
		logger.Info("Получение всех сокращённых ссылок", zap.Int("limit", limit), zap.Int("offset", offset), zap.Int64("count", count))
		res.JSON(w, payload.GetAllLinksResponse{Count: count, Links: links}, http.StatusOK)
	}
}

// UpdateLink метод для обновления ссылки.
func (h *UserHandler) UpdateLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, err := parse.ParseID(r)
		if err != nil {
			logger.Error("Неверный ID для обновления ссылки", zap.Error(err))
			res.ERROR(w, common.ErrInvalidID, http.StatusBadRequest)
			return
		}
		body, err := req.HandleBody[payload.UpdateLinkRequest](&w, r)
		if err != nil {
			logger.Error("Ошибка парсинга тела запроса для обновления ссылки", zap.Error(err))
			res.ERROR(w, common.ErrRequestBodyParse, http.StatusInternalServerError)
			return
		}
		link, err := h.LinkService.Update(ctx, &models.Link{
			Model: gorm.Model{ID: uint(id)},
			Url:   body.URL,
			Hash:  body.Hash,
		})
		if err != nil {
			logger.Error("Ошибка обновления ссылки", zap.Uint("id", uint(id)), zap.Error(err))
			res.ERROR(w, common.ErrLinkUpdateLinkFailed, http.StatusInternalServerError)
			return
		}
		logger.Info("Ссылка успешно обновлена", zap.Uint("id", uint(id)), zap.String("url", body.URL))
		res.JSON(w, link, http.StatusOK)
	}
}

// DeleteLink метод для удаления ссылки по идентификатору.
func (h *UserHandler) DeleteLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, err := parse.ParseID(r)
		if err != nil {
			logger.Error("Неверный ID для удаления ссылки", zap.Error(err))
			res.ERROR(w, common.ErrInvalidID, http.StatusBadRequest)
			return
		}
		_, err = h.LinkService.FindByID(ctx, uint(id))
		if err != nil {
			logger.Error("Ссылка не найдена для удаления", zap.Uint("id", uint(id)), zap.Error(err))
			res.ERROR(w, common.ErrLinkNotFound, http.StatusNotFound)
			return
		}
		err = h.LinkService.Delete(ctx, uint(id))
		if err != nil {
			logger.Error("Ошибка удаления ссылки", zap.Uint("id", uint(id)), zap.Error(err))
			res.ERROR(w, common.ErrLinkDeleteFailed, http.StatusInternalServerError)
			return
		}
		logger.Info("Ссылка успешно удалена", zap.Uint("id", uint(id)))
		res.JSON(w, map[string]string{"message": "ссылка удалена"}, http.StatusOK)
	}
}

// Redirect - редирект на оригинальный URL.
func (h *UserHandler) Redirect() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		hash := r.PathValue("hash")
		if hash == "" {
			logger.Error("Не указан hash для редиректа")
			res.ERROR(w, common.ErrLinkHashNotProvided, http.StatusBadRequest)
			return
		}

		// Получаем ссылку по хешу
		link, err := h.LinkService.GetByHash(ctx, hash)
		if err != nil {
			logger.Error("Ошибка получения ссылки по хешу", zap.String("hash", hash), zap.Error(err))
			res.ERROR(w, common.ErrURLNotFound, http.StatusInternalServerError)
			return
		}

		// Логируем переход в статистику
		go func(linkID uint) {
			timer := time.NewTimer(time.Second)
			defer timer.Stop()

			done := make(chan struct{})

			go func() {
				if err := h.EventBus.Publish(event.Event{Type: event.EventLinkVisited, Data: linkID}); err != nil {
					logger.Error("Ошибка записи события о переходе по ссылке", zap.Uint("linkID", linkID), zap.Error(err))
					res.ERROR(w, common.ErrClickWriteFailed, http.StatusInternalServerError)
				}
				close(done)
			}()

			select {
			case <-done:
			case <-timer.C:
				logger.Info("Таймаут при записи клика", zap.Uint("linkID", linkID))
			}
		}(link.ID)

		logger.Info("Переход по ссылке", zap.String("url", link.Url), zap.String("hash", hash))
		http.Redirect(w, r, link.Url, http.StatusTemporaryRedirect)
	}
}
