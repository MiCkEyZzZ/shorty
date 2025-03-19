package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

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

var (
	ErrRequestBodyParseLink = errors.New("не удалось обработать тело запроса")
	ErrCreateShortURLLink   = errors.New("не удалось создать сокращённый URL")
	ErrInvalidLimit         = errors.New("неверный лимит")
	ErrInvalidOffset        = errors.New("неверное оффсет")
	ErrHashNotProvided      = errors.New("hash не указан")
	ErrURLNotFound          = errors.New("URL-адрес не найден")
	ErrClickWriteFailed     = errors.New("Ошибка записи клика")
	ErrInvalidID            = errors.New("некорректный ID")
	ErrUpdateLinkFailed     = errors.New("не удалось обновить ссылку")
	ErrLinkNotFound         = errors.New("ссылка с таким ID не найдена")
	ErrLinkDeleteFailed     = errors.New("ошибка удаления ссылки")
)

// LinkHandlerDeps - зависимости для создания экземпляра LinkHandler
type LinkHandlerDeps struct {
	Config   *config.Config
	Service  *service.LinkService
	EventBus *event.EventBus
}

type LinkHandler struct {
	Config   *config.Config
	Service  *service.LinkService
	EventBus *event.EventBus
}

// NewLinkHandler создает новый экземпляр LinkHandler
func NewLinkHandler(router *http.ServeMux, deps LinkHandlerDeps) {
	handler := &LinkHandler{
		Config:   deps.Config,
		Service:  deps.Service,
		EventBus: deps.EventBus,
	}
	router.HandleFunc("POST /links", handler.Create())
	router.HandleFunc("GET /links", handler.GetAll())
	router.HandleFunc("GET /links/{hash}", handler.GoTo())
	router.Handle("PATCH /links/{id}", middleware.IsAuth(handler.Update(), deps.Config))
	router.HandleFunc("DELETE /links/{id}", handler.Delete())
}

// Create - создание сокращённого URL.
func (h *LinkHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		body, err := req.HandleBody[payload.CreateLinkRequest](&w, r)
		if err != nil {
			logger.Error("[LinkHandler] Ошибка обработки тела запроса:", zap.Error(err))
			res.ERROR(w, ErrRequestBodyParseLink, http.StatusBadRequest)
			return
		}

		link := models.NewLink(body.URL)
		newLink, err := h.Service.Create(ctx, link)
		if err != nil {
			logger.Error("[LinkHandler] Ошибка создания ссылки:", zap.Error(err))
			res.ERROR(w, ErrCreateShortURLLink, http.StatusBadRequest)
			return
		}

		res.JSON(w, newLink, http.StatusOK)
	}
}

// GetAll - полученте всех сокращённых ссылок.
func (h *LinkHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			logger.Error("[LinkHandler] Ошибка неверный лимит:", zap.Error(err))
			res.ERROR(w, ErrInvalidLimit, http.StatusBadRequest)
			return
		}
		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			logger.Error("[LinkHandler] Ошибка неверный оффсет:", zap.Error(err))
			res.ERROR(w, ErrInvalidOffset, http.StatusBadRequest)
			return
		}
		count, _ := h.Service.Count(ctx)
		links, _ := h.Service.GetAll(ctx, limit, offset)
		res.JSON(w, payload.GetAllLinksResponse{Count: count, Links: links}, http.StatusOK)
	}
}

// Redirect - редирект на оригинальный URL.
func (h *LinkHandler) GoTo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		hash := r.PathValue("hash")
		if hash == "" {
			logger.Error("[LinkHandler] не удалось найти хеш:", zap.String("hash", hash))
			res.ERROR(w, ErrHashNotProvided, http.StatusBadRequest)
			return
		}

		// Получаем ссылку по хешу
		link, err := h.Service.GetByHash(ctx, hash)
		if err != nil {
			logger.Error("[LinkHandler] Ошибка редиректа", zap.String("hash", hash), zap.Error(err))
			res.ERROR(w, ErrURLNotFound, http.StatusInternalServerError)
			return
		}

		// Логируем переход в статистику
		go func(linkID uint) {
			timer := time.NewTimer(time.Second) // Таймаут 1 секунда
			defer timer.Stop()

			done := make(chan struct{}) // Канал завершения

			go func() {
				if err := h.EventBus.Publish(event.Event{Type: event.EventLinkVisited, Data: linkID}); err != nil {
					logger.Error("[StatService] Ошибка записи клика", zap.Uint("linkID", linkID), zap.Error(err))
					res.ERROR(w, ErrClickWriteFailed, http.StatusInternalServerError)
				}
				close(done) // Закрываем канал по завершению
			}()

			select {
			case <-done: // Успешная запись
			case <-timer.C: // Таймаут
				logger.Info("[StatService] Таймаут при записи клика", zap.Uint("linkID", linkID))
			}
		}(link.ID)

		http.Redirect(w, r, link.Url, http.StatusTemporaryRedirect)
	}
}

// Update - обновление сокращённого URL.
func (h *LinkHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, err := parseID(r)
		if err != nil {
			logger.Error("[LinkHandler] Некорректный ID", zap.Error(err))
			res.ERROR(w, ErrInvalidID, http.StatusBadRequest)
			return
		}

		body, err := req.HandleBody[payload.UpdateLinkRequest](&w, r)
		if err != nil {
			logger.Error("[LinkHandler] Ошибка обработки тела запроса", zap.Error(err))
			res.ERROR(w, ErrRequestBodyParseLink, http.StatusInternalServerError)
			return
		}

		link, err := h.Service.Update(ctx, &models.Link{
			Model: gorm.Model{ID: uint(id)},
			Url:   body.URL,
			Hash:  body.Hash,
		})
		if err != nil {
			logger.Error("[LinkHandler] Ошибка обновления ссылки", zap.Uint("id", id), zap.Error(err))
			res.ERROR(w, ErrUpdateLinkFailed, http.StatusInternalServerError)
			return
		}

		res.JSON(w, link, http.StatusOK)
	}
}

// Delete - удаление сокращённого URL.
func (h *LinkHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, err := parseID(r)
		if err != nil {
			logger.Error("[LinkHandler] Некорректный ID", zap.Error(err))
			res.ERROR(w, ErrInvalidID, http.StatusBadRequest)
			return
		}

		_, err = h.Service.FindByID(ctx, uint(id))
		if err != nil {
			logger.Error("[LinkHandler] Попытка удалить несуществующую ссылку", zap.Uint("id", id))
			res.ERROR(w, ErrLinkNotFound, http.StatusNotFound)
			return
		}

		err = h.Service.Delete(ctx, uint(id))
		if err != nil {
			logger.Error("[LinkHandler] Ошибка удаления ссылки (ID: %d): %v", zap.Uint("id", id), zap.Error(err))
			res.ERROR(w, ErrLinkDeleteFailed, http.StatusInternalServerError)
			return
		}
		logger.Info("[LinkHandler] Ссылка (ID: %d) успешно удалена", zap.Uint("id", id))

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
