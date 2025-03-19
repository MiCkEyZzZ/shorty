package handler

import (
	"net/http"

	"shorty/internal/config"
	"shorty/internal/service"
)

type AdminHandlerDeps struct {
	Config  *config.Config
	Service *service.UserService
}

type AdminHandler struct {
	Config  *config.Config
	Service *service.UserService
}

func NewAdminHandler(router *http.ServeMux, deps AdminHandlerDeps) {
	handler := &AdminHandler{
		Config:  deps.Config,
		Service: deps.Service,
	}

	// Статистика.
	router.HandleFunc("GET /admin/stats", handler.GetStats())

	// Управление пользователями.
	router.HandleFunc("GET /admin/users", handler.Getusers())
	router.HandleFunc("GET /admin/users/{id}", handler.GetUser())
	router.HandleFunc("PATCH /admin/users/{id}", handler.UpdateUser())
	router.HandleFunc("DELETE /admin/users/{id}", handler.DeleteUser())
	router.HandleFunc("POST /admin/users", handler.BlockUser())

	// Управление ссылками.
	router.HandleFunc("POST /admin/{id}", handler.BlockLink())
	router.HandleFunc("DELETE /admin/{id}", handler.DeleteLink())
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

func (a *AdminHandler) BlockLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Логика для получения статистики
	}
}

func (a *AdminHandler) DeleteLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Логика для получения статистики
	}
}
