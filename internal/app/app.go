package app

import (
	"context"
	"fmt"

	"shorty/internal/config"
	"shorty/internal/repository"
	"shorty/internal/service"
	"shorty/pkg/db"
	"shorty/pkg/event"
	"shorty/pkg/jwt"
	"shorty/pkg/logger"
	"shorty/pkg/middleware"
)

type App struct {
	Server *Server
}

func NewApp(cfg *config.Config) (*App, error) {
	// Инициализируем логгер в соответствии с указанной средой
	logger.InitLogger(logger.Env(cfg.Env))
	defer logger.Sync()

	// Инициализация БД
	db, err := db.NewDatabase(cfg)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to the database: %w", err)
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
	jwtService := jwt.NewJWT(cfg.Auth.Secret)

	// Промежуточное ПО.
	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)

	// Создаём сервер с обработчиками.
	server := NewServer(cfg, stack, authService, eventBus, linkService, statService, userService, jwtService)

	return &App{Server: server}, nil
}

func (a *App) Run(ctx context.Context) error {
	return a.Server.Start(ctx)
}
