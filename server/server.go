package server

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type config struct {
	hostname     string
	port         string
	forwardedURL string
}

type Server struct {
	logger *slog.Logger
	config *config
	cache  Cache
}

func NewServer(hostname, port, forwardedURL string) *Server {
	return &Server{
		logger: slog.Default(),
		cache:  Cache{},
		config: &config{hostname: hostname, port: port, forwardedURL: forwardedURL},
	}
}

func (srv *Server) Run() error {
	httpServer := &http.Server{
		Addr:         net.JoinHostPort(srv.config.hostname, srv.config.port),
		Handler:      srv.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		srv.logger.Info("shutting down server", slog.Group("server", "signal", s.String(), "addr", httpServer.Addr))

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		srv.cache.data.Clear()

		shutdownError <- httpServer.Shutdown(ctx)
	}()

	srv.logger.Info("starting server", slog.Group("server", "addr", httpServer.Addr))

	if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	if err := <-shutdownError; err != nil {
		return err
	}
	// srv.logger.Info("stopped server", slog.Group("server", "addr", httpServer.Addr))

	return nil
}
