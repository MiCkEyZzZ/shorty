package handler

import (
	"net/http"
	"time"

	"go.uber.org/zap"

	"shorty/internal/common"
	"shorty/internal/config"
	"shorty/internal/service"
	"shorty/pkg/logger"
	"shorty/pkg/middleware"
	"shorty/pkg/res"
)

// StatHandlerDeps - структура для хранения зависимостей обработчика статистики
type StatHandlerDeps struct {
	Config  *config.Config
	Service *service.StatService
}

// StatHandler - обработчик для работы с запросами статистики.
type StatHandler struct {
	Service *service.StatService
}

// NewStatHandler создает новый обработчик статистики и регистрирует его в маршрутизаторе.
func NewStatHandler(router *http.ServeMux, deps StatHandlerDeps) {
	handler := StatHandler{
		Service: deps.Service,
	}

	router.Handle("GET /stats", middleware.IsAuth(handler.GetStats(), deps.Config))
}

func (h *StatHandler) GetStats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		fromStr := r.URL.Query().Get("from")
		from, err := time.Parse("2006-01-02", fromStr)
		if err != nil {
			logger.Error("Ошибка парсинга параметра 'from'", zap.String("from", fromStr), zap.Error(err))
			res.ERROR(w, common.ErrInvalidParam, http.StatusBadRequest)
			return
		}
		toStr := r.URL.Query().Get("to")
		to, err := time.Parse("2006-01-02", toStr)
		if err != nil {
			logger.Error("Ошибка парсинга параметра 'to'", zap.String("to", toStr), zap.Error(err))
			res.ERROR(w, common.ErrInvalidParam, http.StatusBadRequest)
			return
		}
		by := r.URL.Query().Get("by")
		if by != common.GroupByDay && by != common.GroupByMonth {
			logger.Error("Неверное значение параметра 'by'", zap.String("by", by))
			res.ERROR(w, common.ErrInvalidParam, http.StatusBadRequest)
			return
		}
		logger.Info("Получение статистики", zap.String("by", by), zap.Time("from", from), zap.Time("to", to))
		stats := h.Service.GetStats(ctx, by, from, to)
		logger.Info("Статистика успешно получена", zap.Int("record_count", len(stats)))
		res.JSON(w, stats, http.StatusOK)
	}
}
