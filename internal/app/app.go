package app

import (
	"context"
	"fmt"

	"shorty/internal/config"
	"shorty/internal/repository"
	"shorty/internal/service"
	"shorty/pkg/db"
	"shorty/pkg/event"
	"shorty/pkg/middleware"
)

type App struct {
	Server *Server
}

func NewApp(cfg *config.Config) (*App, error) {
	// Инициализация БД
	db, err := db.NewDatabase(cfg)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к базе данных: %w", err)
	}

	// Создаём EventBus.
	eventBus := event.NewEventBus()

	// Репозитории.
	linkRepository := repository.NewLinkRepository(db)
	userRepository := repository.NewUserRepository(db)
	statRepository := repository.NewStatRepository(db)

	// Сервисы.
	linkService := service.NewLinkService(linkRepository)
	userService := service.NewUserService(userRepository)
	authService := service.NewAuthService(userRepository)
	statService := service.NewStatService(&service.StatServiceDeps{EventBus: eventBus, Repo: statRepository})

	// Промежуточное ПО.
	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)

	// Создаём сервер с обработчиками.
	server := NewServer(cfg, stack, authService, eventBus, linkService, statService, userService)

	return &App{Server: server}, nil
}

func (a *App) Run(ctx context.Context) error {
	return a.Server.Start(ctx)
}
