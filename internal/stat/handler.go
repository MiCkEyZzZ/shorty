package stat

import (
	"net/http"
	"time"

	"shorty/internal/config"
	"shorty/pkg/res"
)

const (
	GroupByDay   = "day"
	GroupByMonth = "month"
)

// StatHandlerDeps - зависимости для создания экземпляра StatHandler
type StatHandlerDeps struct {
	*config.Config
	Service *StatService
}

type StatHandler struct {
	Service *StatService
}

func NewStatHandler(router *http.ServeMux, deps StatHandlerDeps) {
	handler := StatHandler{
		Service: deps.Service,
	}
	router.HandleFunc("GET /stats", handler.Getstats())
}

func (h *StatHandler) Getstats() http.HandlerFunc {
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
		if by != GroupByDay && by != GroupByMonth {
			http.Error(w, "Invalid by param", http.StatusBadRequest)
			return
		}
		stats := h.Service.GetAllStat(ctx, by, from, to)
		res.Json(w, stats, http.StatusOK)
	}
}
