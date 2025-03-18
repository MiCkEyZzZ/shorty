package main

import (
	"context"
	"fmt"
	"net/http"

	"shorty/internal/auth"
	"shorty/internal/config"
	"shorty/internal/link"
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

	// Репозитории.
	linkRepository := link.NewLinkRepository(db)
	userRepository := user.NewUserRepository(db)
	statRepository := stat.NewStatRepository(db)

	// Сервисы.
	linkService := link.NewLinkService(linkRepository)
	userService := user.NewUserService(userRepository)
	authService := auth.NewAuthService(userRepository)
	statService := stat.NewStatService(&stat.StatServiceDeps{
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
