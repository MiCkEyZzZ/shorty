package main

import (
	"context"
	"fmt"
	"net/http"

	"shorty/internal/auth"
	"shorty/internal/config"
	"shorty/internal/link"
	"shorty/internal/repository"
	"shorty/internal/service"
	"shorty/internal/stat"
	"shorty/internal/user"
	"shorty/pkg/db"
	"shorty/pkg/event"
	"shorty/pkg/middleware"

	_ "github.com/lib/pq"
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

	// Репохитории.
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
	link.NewLinkHandler(router, link.LinkHandlerDeps{Config: cfg, Service: linkService, EventBus: eventBus})
	user.NewUserHandler(router, user.UserHandlerDeps{Config: cfg, Service: userService})
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{Config: cfg, Service: authService})
	stat.NewStatHandler(router, stat.StatHandlerDeps{Config: cfg, Service: statService})

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
