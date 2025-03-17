package link

import (
	"log"
	"net/http"
	"strconv"

	"gorm.io/gorm"

	"shorty/internal/config"
	"shorty/internal/models"
	"shorty/internal/service"
	"shorty/pkg/middleware"
	"shorty/pkg/req"
	"shorty/pkg/res"
)

// LinkHandlerDeps - зависимости для создания экземпляра LinkHandler
type LinkHandlerDeps struct {
	*config.Config
	Service *service.LinkService
}

type LinkHandler struct {
	*config.Config
	Service *service.LinkService
}

// NewLinkHandler создает новый экземпляр LinkHandler
func NewLinkHandler(router *http.ServeMux, deps LinkHandlerDeps) {
	handler := &LinkHandler{
		Config:  deps.Config,
		Service: deps.Service,
	}
	router.HandleFunc("POST /links", handler.Create())
	router.HandleFunc("GET /links/{hash}", handler.GoTo())
	router.Handle("PATCH /links/{id}", middleware.IsAuth(handler.Update()))
	router.HandleFunc("DELETE /links/{id}", handler.Delete())
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

		link := models.NewLink(body.URL)
		newLink, err := h.Service.Create(ctx, link)
		if err != nil {
			log.Printf("[LinkHandler] Ошибка создания ссылки: %v", err)
			http.Error(w, "не удалось создать сокращённый URL", http.StatusBadRequest)
			return
		}

		res.Json(w, newLink, http.StatusOK)
	}
}

// Redirect - редирект на оригинальный URL.
func (h *LinkHandler) GoTo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		hash := r.URL.Path[len("/links/"):]
		if hash == "" {
			log.Printf("[LinkHandler] не удалось найти хеш %s", hash)
			http.Error(w, "Hash не указан", http.StatusBadRequest)
			return
		}

		originalURL, err := h.Service.GetByHash(ctx, hash)
		if err != nil {
			log.Printf("[LinkHandler] Ошибка редиректа %s: %v", hash, err)
			http.Error(w, "URL-адрес не найден", http.StatusNotFound)
			return
		}
		http.Redirect(w, r, originalURL.Url, http.StatusTemporaryRedirect)
	}
}

// Update - обновление сокращённого URL.
func (h *LinkHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		body, err := req.HandleBody[UpdateLinkRequest](&w, r)
		if err != nil {
			log.Printf("[LinkHandler] Ошибка обработки тела запроса: %v", err)
			http.Error(w, "не удалось обработать тело запроса", http.StatusBadRequest)
			return
		}
		rid := r.PathValue("id")
		id, err := strconv.ParseUint(rid, 10, 32)
		if err != nil {
			log.Printf("[LinkHandler] Некорректный ID: %v", err)
			http.Error(w, "Некорректный ID", http.StatusBadRequest)
			return
		}
		link, err := h.Service.Update(ctx, &models.Link{
			Model: gorm.Model{ID: uint(id)},
			Url:   body.URL,
			Hash:  body.Hash,
		})
		if err != nil {
			log.Printf("[LinkHandler] Ошибка обновления ссылки (ID: %d): %v", id, err)
			http.Error(w, "Не удалось обновить ссылку", http.StatusBadRequest)
			return
		}

		res.Json(w, link, http.StatusOK)
	}
}

// Delete - удаление сокращённого URL.
func (h *LinkHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		rid := r.PathValue("id")
		id, err := strconv.ParseUint(rid, 10, 32)
		if err != nil {
			log.Printf("[LinkHandler] Некорректный ID: %v", err)
			http.Error(w, "Некорректный ID", http.StatusBadRequest)
			return
		}
		_, err = h.Service.FindByID(ctx, uint(id))
		if err != nil {
			log.Printf("[LinkHandler] Попытка удалить несуществующую ссылку (ID: %d)", id)
			http.Error(w, "Ссылка с таким ID не найдена", http.StatusNotFound)
			return
		}
		err = h.Service.Delete(ctx, uint(id))
		if err != nil {
			log.Printf("[LinkHandler] Ошибка удаления ссылки (ID: %d): %v", id, err)
			http.Error(w, "Не удалось удалить ссылку", http.StatusInternalServerError)
			return
		}
		log.Printf("[LinkHandler] Ссылка (ID: %d) успешно удалена", id)
		res.Json(w, map[string]string{"message": "Ссылка удалена"}, http.StatusOK)
	}
}
