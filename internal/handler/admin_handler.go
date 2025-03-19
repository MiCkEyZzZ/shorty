package handler

import (
	"net/http"

	"shorty/internal/config"
	"shorty/internal/service"
	"shorty/pkg/middleware"
)

type AdminHandlerDeps struct {
	Config      *config.Config
	UserService *service.UserService
}

type AdminHandler struct {
	Config      *config.Config
	UserService *service.UserService
}

func NewAdminHandler(router *http.ServeMux, deps AdminHandlerDeps) {
	handler := &AdminHandler{
		Config:      deps.Config,
		UserService: deps.UserService,
	}

	// Статистика.
	router.HandleFunc("GET /admin/stats", middleware.AdminOnly(handler.GetStats()))

	// Управление пользователями.
	router.HandleFunc("GET /admin/users", middleware.AdminOnly(handler.Getusers()))
	router.HandleFunc("GET /admin/users/{id}", middleware.AdminOnly(handler.GetUser()))
	router.HandleFunc("PATCH /admin/users/{id}", middleware.AdminOnly(handler.UpdateUser()))
	router.HandleFunc("DELETE /admin/users/{id}", middleware.AdminOnly(handler.DeleteUser()))
	router.HandleFunc("PATCH /admin/users/{id}/block", middleware.AdminOnly(handler.BlockUser()))
	router.HandleFunc("PATCH /admin/users/{id}/unblock", middleware.AdminOnly(handler.UnblockUser()))

	// Управление ссылками.
	router.HandleFunc("POST /admin/links/{id}/block", middleware.AdminOnly(handler.BlockLink()))
	router.HandleFunc("PATCH /admin/links/{id}/unblock", middleware.AdminOnly(handler.UnblockLink()))
	router.HandleFunc("DELETE /admin/links/{id}", middleware.AdminOnly(handler.DeleteLink()))
}

func (a *AdminHandler) GetStats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Логика для получения статистики
	}
}

func (a *AdminHandler) Getusers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Логика для получения статистики
	}
}

func (a *AdminHandler) GetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Логика для получения статистики
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
		// TODO: Логика для получения статистики
	}
}

func (a *AdminHandler) UnblockLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Логика для получения статистики
	}
}

func (a *AdminHandler) DeleteLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Логика для получения статистики
	}
}
