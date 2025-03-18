package stat

import (
	"context"
	"errors"
	"log"
	"time"

	"shorty/pkg/event"
)

var ErrInvalidLinkID = errors.New("неверный идентификатор ссылки")

type StatServiceDeps struct {
	Repo     *StatRepository
	EventBus *event.EventBus
}

type StatService struct {
	Repo     *StatRepository
	EventBus *event.EventBus
}

func NewStatService(deps *StatServiceDeps) *StatService {
	return &StatService{Repo: deps.Repo, EventBus: deps.EventBus}
}

func (s *StatService) AddClick(ctx context.Context) {
	go func() {
		for msg := range s.EventBus.Subscribe() {
			if msg.Type == event.EventLinkVisited {
				id, ok := msg.Data.(uint)
				if !ok {
					log.Println("Bad EventLinkVisited Data:", msg.Data)
					continue
				}
				s.Repo.AddClick(ctx, id)
			}
		}
	}()
}

func (s *StatService) GetAllStat(ctx context.Context, by string, from, to time.Time) []GetStatsResponse {
	return s.Repo.GetStats(ctx, by, from, to)
}
