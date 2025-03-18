package handler

import (
	"log"
	"net/http"
	"strconv"

	"shorty/internal/config"
	"shorty/internal/models"
	"shorty/internal/service"
	"shorty/pkg/req"
	"shorty/pkg/res"
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
			log.Printf("[UserHandler] Ошибка получения списка пользователей: %v", err)
			http.Error(w, "Не удалось получить список пользователей", http.StatusInternalServerError)
			return
		}
		res.Json(w, users, http.StatusOK)
	}
}

// FindByID получает пользователя по идентификатору.
func (h *UserHandler) FindByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Printf("[UserHandler] Некорректный ID пользователя: %v", err)
			http.Error(w, "Некорректный ID пользователя", http.StatusBadRequest)
			return
		}
		user, err := h.Service.GetByID(ctx, uint(id))
		if err != nil {
			log.Printf("[UserHandler] Ошибка поиска пользователя (ID: %d): %v", id, err)
			http.Error(w, "Пользователь не найден", http.StatusNotFound)
			return
		}
		res.Json(w, user, http.StatusOK)
	}
}

// Update обновляет пользователя.
func (h *UserHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Printf("[UserHandler] Некорректный ID пользователя: %v", err)
			http.Error(w, "Некорректный ID пользователя", http.StatusBadRequest)
			return
		}
		body, err := req.HandleBody[models.User](&w, r)
		if err != nil {
			log.Printf("[UserHandler] Ошибка обработки тела запроса: %v", err)
			http.Error(w, "Не удалось обработать тело запроса", http.StatusBadRequest)
			return
		}
		body.ID = uint(id)
		updatedUser, err := h.Service.Update(ctx, body)
		if err != nil {
			log.Printf("[UserHandler] Ошибка обновления пользователя (ID: %d): %v", id, err)
			http.Error(w, "Не удалось обновить пользователя", http.StatusInternalServerError)
			return
		}
		res.Json(w, updatedUser, http.StatusOK)
	}
}

// Delete удаляет пользователя по идентификатору.
func (h *UserHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Printf("[UserHandler] Некорректный ID пользователя: %v", err)
			http.Error(w, "Некорректный ID пользователя", http.StatusBadRequest)
			return
		}
		err = h.Service.Delete(ctx, uint(id))
		if err != nil {
			log.Printf("[UserHandler] Ошибка удаления пользователя (ID: %d): %v", id, err)
			http.Error(w, "Не удалось удалить пользователя", http.StatusInternalServerError)
			return
		}
		log.Printf("[UserHandler] Пользователь (ID: %d) успешно удалён", id)
		res.Json(w, map[string]string{"message": "Пользователь удалён"}, http.StatusOK)
	}
}
