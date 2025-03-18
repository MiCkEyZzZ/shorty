package main

import (
	"context"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"

	"shorty/internal/config"
	"shorty/internal/handler"
	"shorty/internal/repository"
	"shorty/internal/service"
	"shorty/pkg/db"
	"shorty/pkg/event"
	"shorty/pkg/middleware"
)

func main() {
	cfg := config.NewConfig()
	router := http.NewServeMux()
	db, err := db.NewDatabase(cfg)
	if err != nil {
		fmt.Printf("Не удалось подключиться к базе данных: %v", err)
		return
	}
	eventBus := event.NewEventBus()

	// Репозитории.
	linkRepository := repository.NewLinkRepository(db)
	userRepository := repository.NewUserRepository(db)
	statRepository := repository.NewStatRepository(db)

	// Сервисы.
	linkService := service.NewLinkService(linkRepository)
	userService := service.NewUserService(userRepository)
	authService := service.NewAuthService(userRepository)
	statService := service.NewStatService(&service.StatServiceDeps{
		EventBus: eventBus,
		Repo:     statRepository,
	})

	// Обработчики.
	handler.NewLinkHandler(router, handler.LinkHandlerDeps{Config: cfg, Service: linkService, EventBus: eventBus})
	handler.NewUserHandler(router, handler.UserHandlerDeps{Config: cfg, Service: userService})
	handler.NewAuthHandler(router, handler.AuthHandlerDeps{Config: cfg, Service: authService})
	handler.NewStatHandler(router, handler.StatHandlerDeps{Config: cfg, Service: statService})

	// Промежуточное ПО.
	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)

	ctx := context.Background()
	go statService.AddClick(ctx)

	// HTTP-сервер.
	server := http.Server{
		Addr:    ":8080",
		Handler: stack(router),
	}

	fmt.Println("Сервер запущен на http://localhost:8080")
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("Не удалось запустить сервер: %v", err)
	}
}
