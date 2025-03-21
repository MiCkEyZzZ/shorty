package service

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"

	"shorty/internal/payload"
	"shorty/internal/repository"
	"shorty/pkg/event"
	"shorty/pkg/logger"
)

var ErrInvalidLinkID = errors.New("неверный идентификатор ссылки")

type StatServiceDeps struct {
	Repo     *repository.StatRepository
	EventBus *event.EventBus
}

// StatService предоставляет методы для работы с статистикой.
type StatService struct {
	Repo     *repository.StatRepository
	EventBus *event.EventBus
}

// NewUserService создаёт новый экземпляр StatService.
func NewStatService(deps *StatServiceDeps) *StatService {
	ctx := context.Background()
	service := &StatService{Repo: deps.Repo, EventBus: deps.EventBus}
	go service.AddClick(ctx)
	return service
}

func (s *StatService) AddClick(ctx context.Context) {
	for msg := range s.EventBus.Subscribe() {
		if msg.Type == event.EventLinkVisited {
			linkID, ok := msg.Data.(uint)
			if !ok {
				logger.Error("Неверные данные при получении события EventLinkVisited", zap.Any("data", msg.Data))
				continue
			}
			if err := s.Repo.AddClick(context.Background(), linkID); err != nil {
				logger.Error("Ошибка при добавлении клика", zap.Uint("linkID", linkID), zap.Error(err))
			} else {
				logger.Info("Добавлен клик", zap.Uint("linkID", linkID))
			}
		}
	}
}

// GetStats метод для получения статистики.
func (s *StatService) GetStats(ctx context.Context, by string, from, to time.Time) []payload.GetStatsResponse {
	logger.Info("Запрос статистики", zap.String("by", by), zap.Time("from", from), zap.Time("to", to))
	stats := s.Repo.GetStats(ctx, by, from, to)
	logger.Info("Статистика получена", zap.Int("count", len(stats)))
	return stats
}

func (s *StatService) GetAllLinksStats(ctx context.Context, from, to time.Time) []payload.LinkStatsResponse {
	logger.Info("Запрос статистики по всем ссылкам", zap.Time("from", from), zap.Time("to", to))
	stats := s.Repo.GetAllLinksStats(ctx, from, to)
	return stats
}
