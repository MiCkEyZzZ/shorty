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
	"shorty/pkg/req"
	"shorty/pkg/res"
)

// LinkHandlerDeps - зависимости для создания экземпляра LinkHandler
type LinkHandlerDeps struct {
	Config      *config.Config
	LinkService *service.LinkService
	EventBus    *event.EventBus
}

// LinkHandler - обработчик коротких ссылок.
type LinkHandler struct {
	Config      *config.Config
	LinkService *service.LinkService
	EventBus    *event.EventBus
}

// NewLinkHandler создает новый экземпляр LinkHandler
func NewLinkHandler(router *http.ServeMux, deps LinkHandlerDeps) {
	handler := &LinkHandler{
		Config:      deps.Config,
		LinkService: deps.LinkService,
		EventBus:    deps.EventBus,
	}

	router.HandleFunc("POST /links", handler.Create())
	router.Handle("GET /links", middleware.IsAuth(handler.GetAll(), deps.Config))
	router.HandleFunc("GET /links/{hash}", handler.GoTo())
	router.Handle("PATCH /links/{id}", middleware.IsAuth(handler.Update(), deps.Config))
	router.Handle("DELETE /links/{id}", middleware.IsAuth(handler.Delete(), deps.Config))
}

// Create - создание сокращённого URL.
func (h *LinkHandler) Create() http.HandlerFunc {
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

// GetAll - полученте всех сокращённых ссылок.
func (h *LinkHandler) GetAll() http.HandlerFunc {
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

// Redirect - редирект на оригинальный URL.
func (h *LinkHandler) GoTo() http.HandlerFunc {
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

// Update - обновление сокращённого URL.
func (h *LinkHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, err := parseID(r)
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

// Delete - удаление сокращённого URL.
func (h *LinkHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, err := parseID(r)
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

// parseID парсит идентификатор из строки в uint.
func parseID(r *http.Request) (uint, error) {
	rid := r.PathValue("id")
	id, err := strconv.ParseUint(rid, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
