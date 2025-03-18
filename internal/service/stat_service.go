package service

import (
	"context"
	"errors"
	"log"

	"shorty/internal/repository"
	"shorty/pkg/event"
)

var ErrInvalidLinkID = errors.New("неверный идентификатор ссылки")

type StatServiceDeps struct {
	EventBus *event.EventBus
	Repo     *repository.StatRepository
}

type StatService struct {
	EventBus *event.EventBus
	Repo     *repository.StatRepository
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
