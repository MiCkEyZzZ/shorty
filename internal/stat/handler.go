package stat

import (
	"net/http"
	"shorty/internal/config"
	"shorty/internal/service"
)

// StatHandlerDeps - зависимости для создания экземпляра StatHandler
type StatHandlerDeps struct {
	*config.Config
	Service *service.StatService
}

type StatHandler struct {
	*config.Config
	Service *service.StatService
}

func NewStatHandler(router *http.ServeMux, deps StatHandlerDeps) {
	handler := StatHandler{
		Config:  deps.Config,
		Service: deps.Service,
	}
	router.HandleFunc("GET /stat", handler.Getstats)
}

func (s *StatHandler) Getstats(w http.ResponseWriter, r *http.Request) {}
