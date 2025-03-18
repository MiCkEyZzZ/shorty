package service

import (
	"context"
	"errors"
	"log"
	"time"

	"shorty/internal/payload"
	"shorty/internal/repository"
	"shorty/pkg/event"
)

var ErrInvalidLinkID = errors.New("неверный идентификатор ссылки")

type StatServiceDeps struct {
	Repo     *repository.StatRepository
	EventBus *event.EventBus
}

type StatService struct {
	Repo     *repository.StatRepository
	EventBus *event.EventBus
}

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
				log.Println("[StatService] Неверные данные EventLinkVisited:", msg.Data)
				continue
			}
			if err := s.Repo.AddClick(context.Background(), linkID); err != nil {
				log.Printf("[StatService] Ошибка при добавлении клика (LinkID: %d): %v", linkID, err)
			}
		}
	}
}

func (s *StatService) GetStats(ctx context.Context, by string, from, to time.Time) []payload.GetStatsResponse {
	return s.Repo.GetStats(ctx, by, from, to)
}
