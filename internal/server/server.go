package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/rashidmailru/kabobfood/internal/config"
)

// Server wraps the HTTP server lifecycle.
type Server struct {
	httpServer      *http.Server
	log             *zap.Logger
	shutdownTimeout time.Duration
}

// New constructs the HTTP server configured from Config.
func New(cfg *config.Config, handler http.Handler, log *zap.Logger) *Server {
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port),
		Handler:      handler,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
	}

	return &Server{httpServer: srv, log: log, shutdownTimeout: cfg.ShutdownTimeout}
}

// Run starts the HTTP server.
func (s *Server) Run() error {
	s.log.Info("HTTP server listening", zap.String("addr", s.httpServer.Addr))
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Shutdown gracefully stops the HTTP server.
func (s *Server) Shutdown(ctx context.Context) error {
	timeout := s.shutdownTimeout
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}
