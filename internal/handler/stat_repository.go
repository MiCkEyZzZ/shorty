package handler

import (
	"net/http"
	"time"

	"shorty/internal/common"
	"shorty/internal/config"
	"shorty/internal/service"
	"shorty/pkg/res"
)

// StatHandlerDeps - зависимости для создания экземпляра StatHandler
type StatHandlerDeps struct {
	Config  *config.Config
	Service *service.StatService
}

type StatHandler struct {
	Service *service.StatService
}

func NewStatHandler(router *http.ServeMux, deps StatHandlerDeps) {
	handler := StatHandler{
		Service: deps.Service,
	}
	router.HandleFunc("GET /stats", handler.GetStats())
}

func (h *StatHandler) GetStats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		from, err := time.Parse("2006-01-02", r.URL.Query().Get("from"))
		if err != nil {
			http.Error(w, "Invalid from param", http.StatusBadRequest)
			return
		}
		to, err := time.Parse("2006-01-02", r.URL.Query().Get("to"))
		if err != nil {
			http.Error(w, "Invalid to param", http.StatusBadRequest)
			return
		}
		by := r.URL.Query().Get("by")
		if by != common.GroupByDay && by != common.GroupByMonth {
			http.Error(w, "Invalid by param", http.StatusBadRequest)
			return
		}
		stats := h.Service.GetStats(ctx, by, from, to)
		res.Json(w, stats, http.StatusOK)
	}
}
