package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"shorty/internal/config"
	"shorty/internal/handler"
	"shorty/internal/service"
	"shorty/pkg/event"
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
	handler.NewLinkHandler(router, handler.LinkHandlerDeps{Config: cfg, Service: linkService, EventBus: eventBus})
	handler.NewUserHandler(router, handler.UserHandlerDeps{Config: cfg, Service: userService})
	handler.NewAuthHandler(router, handler.AuthHandlerDeps{Config: cfg, Service: authService})
	handler.NewStatHandler(router, handler.StatHandlerDeps{Config: cfg, Service: statService})

	server := &http.Server{
		Addr:    ":8080",
		Handler: middleware(router),
	}

	return &Server{httpServer: server}
}

func (s *Server) Start(ctx context.Context) error {
	errChan := make(chan error, 1)

	go func() {
		fmt.Printf("Сервер запущен на %s", s.httpServer.Addr)
		errChan <- s.httpServer.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		log.Println("Выключение сервера...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return s.httpServer.Shutdown(shutdownCtx)
	case err := <-errChan:
		return err
	}
}
