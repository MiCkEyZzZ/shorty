package main

import (
	"fmt"
	"net/http"

	"shorty/internal/auth"
	"shorty/internal/config"
	"shorty/internal/link"
	"shorty/internal/repository"
	"shorty/internal/service"
	"shorty/pkg/db"
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

	// Репохитории
	linkRepository := repository.NewLinkRepository(db)
	// userRepository := repository.NewUserRepository(db)
	// statRepository := repository.NewStatRepository(db)

	// Сервисы
	linkService := service.NewLinkService(linkRepository)
	// userService := service.NewUserService(userRepository)
	// statService := service.NewStatService(statRepository)

	// Обработчики
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{Config: cfg})
	link.NewLinkHandler(router, link.LinkHandlerDeps{Config: cfg, Service: linkService})
	// user.NewUserHandler(router, user.UserHandlerDeps{Config: cfg, Service: userService})
	// stat.NewStatHandler(router, stat.StatHandlerDeps{Config: cfg, Service: statService})

	// Middleware
	middleware := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
		middleware.IsAuth,
	)

	// HTTP-сервер
	server := http.Server{
		Addr:    ":8080",
		Handler: middleware(router),
	}

	fmt.Println("Сервер запущен на http://localhost:8080")
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("Не удалось запустить сервер: %v", err)
	}
}
