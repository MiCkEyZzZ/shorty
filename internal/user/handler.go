package user

import (
	"net/http"

	"shorty/internal/config"
	"shorty/internal/service"
)

// UserHandlerDeps - зависимости для создания экземпляра UserHandler
type UserHandlerDeps struct {
	*config.Config
	Service *service.UserService
}

type UserHandler struct {
	*config.Config
	Service *service.UserService
}

func NewUserHandler(router *http.ServeMux, deps UserHandlerDeps) {
	handler := UserHandler{
		Config:  deps.Config,
		Service: deps.Service,
	}
	router.HandleFunc("POST /users", handler.SaveUser)
	router.HandleFunc("GET /users", handler.FindUsers)
	router.HandleFunc("GET /users/{id}", handler.FindUserByID)
	router.HandleFunc("GET /users/{email}", handler.GetByEmail)
	router.HandleFunc("PUT /users/{id}", handler.UpdateUser)
	router.HandleFunc("DELETE /users/{id}", handler.DeleteUser)
}

func (u *UserHandler) SaveUser(w http.ResponseWriter, r *http.Request) {}

func (u *UserHandler) FindUsers(w http.ResponseWriter, r *http.Request) {}

func (u *UserHandler) FindUserByID(w http.ResponseWriter, r *http.Request) {}

func (u *UserHandler) GetByEmail(w http.ResponseWriter, r *http.Request) {}

func (u *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {}

func (u *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {}
