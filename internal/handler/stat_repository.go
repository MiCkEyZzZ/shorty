package handler

import (
	"errors"
	"net/http"
	"time"

	"shorty/internal/common"
	"shorty/internal/config"
	"shorty/internal/service"
	"shorty/pkg/res"
)

var (
	ErrInvalidParam = errors.New("неверный параметр")
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
			res.ERROR(w, ErrInvalidParam, http.StatusBadRequest)
			return
		}
		to, err := time.Parse("2006-01-02", r.URL.Query().Get("to"))
		if err != nil {
			res.ERROR(w, ErrInvalidParam, http.StatusBadRequest)
			return
		}
		by := r.URL.Query().Get("by")
		if by != common.GroupByDay && by != common.GroupByMonth {
			res.ERROR(w, ErrInvalidParam, http.StatusBadRequest)
			return
		}
		stats := h.Service.GetStats(ctx, by, from, to)
		res.JSON(w, stats, http.StatusOK)
	}
}
