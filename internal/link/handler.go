package link

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"gorm.io/gorm"

	"shorty/internal/config"
	"shorty/pkg/event"
	"shorty/pkg/middleware"
	"shorty/pkg/req"
	"shorty/pkg/res"
)

// LinkHandlerDeps - зависимости для создания экземпляра LinkHandler
type LinkHandlerDeps struct {
	*config.Config
	Service  *LinkService
	EventBus *event.EventBus
}

type LinkHandler struct {
	*config.Config
	Service  *LinkService
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

// parseID парсит идентификатор из строки в uint.
func parseID(r *http.Request) (uint, error) {
	rid := r.PathValue("id")
	id, err := strconv.ParseUint(rid, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

// Create - создание сокращённого URL.
func (h *LinkHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		body, err := req.HandleBody[CreateLinkRequest](&w, r)
		if err != nil {
			log.Printf("[LinkHandler] Ошибка обработки тела запроса: %v", err)
			http.Error(w, "не удалось обработать тело запроса", http.StatusBadRequest)
			return
		}

		link := NewLink(body.URL)
		newLink, err := h.Service.Create(ctx, link)
		if err != nil {
			log.Printf("[LinkHandler] Ошибка создания ссылки: %v", err)
			http.Error(w, "не удалось создать сокращённый URL", http.StatusBadRequest)
			return
		}

		res.Json(w, newLink, http.StatusOK)
	}
}

func (h *LinkHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			http.Error(w, "Invalid limi", http.StatusBadRequest)
			return
		}
		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			http.Error(w, "Invalid offset", http.StatusBadRequest)
			return
		}
		count := h.Service.Count(ctx)
		links := h.Service.GetLinks(ctx, limit, offset)
		res.Json(w, GetAllLinksResponse{Count: count, Links: links}, http.StatusOK)
	}
}

// Redirect - редирект на оригинальный URL.
func (h *LinkHandler) GoTo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		hash := r.URL.Path[len("/links/"):]
		if hash == "" {
			log.Printf("[LinkHandler] не удалось найти хеш %s", hash)
			http.Error(w, "hash не указан", http.StatusBadRequest)
			return
		}

		// Получаем ссылку по хешу
		link, err := h.Service.GetByHash(ctx, hash)
		if err != nil {
			log.Printf("[LinkHandler] Ошибка редиректа %s: %v", hash, err)
			http.Error(w, "URL-адрес не найден", http.StatusNotFound)
			return
		}

		// Логируем переход в статистику
		go func(linkID uint) {
			timer := time.NewTimer(time.Second) // Таймаут 1 секунда
			defer timer.Stop()

			done := make(chan struct{}) // Канал завершения

			go func() {
				if err := h.EventBus.Publish(event.Event{Type: event.EventLinkVisited, Data: linkID}); err != nil {
					log.Printf("[StatService] Ошибка записи клика (linkID: %d): %v", linkID, err)
				}
				close(done) // Закрываем канал по завершению
			}()

			select {
			case <-done: // Успешная запись
			case <-timer.C: // Таймаут
				log.Printf("[StatService] Таймаут при записи клика (linkID: %d)", linkID)
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
			log.Printf("[LinkHandler] Некорректный ID: %v", err)
			http.Error(w, "некорректный ID", http.StatusBadRequest)
			return
		}

		body, err := req.HandleBody[UpdateLinkRequest](&w, r)
		if err != nil {
			log.Printf("[LinkHandler] Ошибка обработки тела запроса: %v", err)
			http.Error(w, "не удалось обработать тело запроса", http.StatusBadRequest)
			return
		}

		link, err := h.Service.Update(ctx, &Link{
			Model: gorm.Model{ID: uint(id)},
			Url:   body.URL,
			Hash:  body.Hash,
		})
		if err != nil {
			log.Printf("[LinkHandler] Ошибка обновления ссылки (ID: %d): %v", id, err)
			http.Error(w, "не удалось обновить ссылку", http.StatusBadRequest)
			return
		}

		res.Json(w, link, http.StatusOK)
	}
}

// Delete - удаление сокращённого URL.
func (h *LinkHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, err := parseID(r)
		if err != nil {
			log.Printf("[LinkHandler] Некорректный ID: %v", err)
			http.Error(w, "некорректный ID", http.StatusBadRequest)
			return
		}

		_, err = h.Service.FindByID(ctx, uint(id))
		if err != nil {
			log.Printf("[LinkHandler] Попытка удалить несуществующую ссылку (ID: %d)", id)
			http.Error(w, "ссылка с таким ID не найдена", http.StatusNotFound)
			return
		}

		err = h.Service.Delete(ctx, uint(id))
		if err != nil {
			log.Printf("[LinkHandler] Ошибка удаления ссылки (ID: %d): %v", id, err)
			http.Error(w, "ошибка сервера", http.StatusInternalServerError)
			return
		}
		log.Printf("[LinkHandler] Ссылка (ID: %d) успешно удалена", id)
		res.Json(w, map[string]string{"message": "ссылка удалена"}, http.StatusOK)
	}
}
