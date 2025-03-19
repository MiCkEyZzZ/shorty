package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"

	"shorty/internal/app"
	"shorty/internal/config"
	"shorty/pkg/logger"
)

func main() {
	logger.InitLogger("development")
	defer logger.Sync()
	cfg := config.NewConfig()

	app, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("Ошибка инициализации приложения: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := app.Run(ctx); err != nil {
		fmt.Printf("Ошибка запуска сервера %v", err)
	}
}
