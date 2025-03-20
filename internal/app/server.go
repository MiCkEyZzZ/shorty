package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"shorty/internal/config"
	"shorty/internal/handler"
	"shorty/internal/service"
	"shorty/pkg/event"
	"shorty/pkg/logger"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(
	cfg *config.Config,
	middleware func(http.Handler) http.Handler,
	authService *service.AuthService,
	eventBus *event.EventBus,
	linkService *service.LinkService,
	statService *service.StatService,
	userService *service.UserService,
) *Server {
	router := http.NewServeMux()

	// Обработчики.
	handler.NewAdminHandler(router, handler.AdminHandlerDeps{
		Config:      cfg,
		UserService: userService,
		LinkService: linkService,
		StatService: statService,
	})
	handler.NewAuthHandler(router, handler.AuthHandlerDeps{
		Config:      cfg,
		AuthService: authService,
	})
	handler.NewStatHandler(router, handler.StatHandlerDeps{
		Config:  cfg,
		Service: statService,
	})
	handler.NewUserHandler(router, handler.UserHandlerDeps{
		Config:      cfg,
		UserService: userService,
		LinkService: linkService,
		EventBus:    eventBus,
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: middleware(router),
	}

	return &Server{httpServer: server}
}

func (s *Server) Start(ctx context.Context) error {
	errChan := make(chan error, 1)

	go func() {
		fmt.Printf("Сервер запущен на %s\n", s.httpServer.Addr)
		fmt.Println()
		errChan <- s.httpServer.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		logger.Info("Выключение сервера")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return s.httpServer.Shutdown(shutdownCtx)
	case err := <-errChan:
		return err
	}
}
