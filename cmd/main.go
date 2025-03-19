package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"shorty/internal/app"
	"shorty/internal/config"
	"shorty/pkg/logger"
)

func main() {
	cfg := config.NewConfig()

	app, err := app.NewApp(cfg)
	if err != nil {
		logger.Error("Ошибка инициализации приложения: %v", zap.Error(err))
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := app.Run(ctx); err != nil {
		logger.Error("Ошибка запуска сервера %v", zap.Error(err))
		fmt.Printf("Ошибка запуска сервера %v", err)
	}
}
