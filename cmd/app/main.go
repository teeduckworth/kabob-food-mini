package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/rashidmailru/kabobfood/internal/app"
	"github.com/rashidmailru/kabobfood/internal/config"
	"github.com/rashidmailru/kabobfood/internal/observability"
	"github.com/rashidmailru/kabobfood/pkg/logger"
)

func main() {
	cfg := config.MustLoad()

	log, err := logger.New(cfg.AppEnv)
	if err != nil {
		panic(err)
	}
	defer func() { _ = log.Sync() }()

	if err := observability.InitSentry(cfg.Sentry.DSN); err != nil {
		log.Warn("sentry init failed", zap.Error(err))
	}
	defer observability.Flush(2 * time.Second)

	application, err := app.New(cfg, log)
	if err != nil {
		observability.CaptureError(err)
		log.Fatal("failed to build app", zap.Error(err))
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		errCh <- application.Run()
	}()

	select {
	case <-ctx.Done():
		log.Info("shutting down...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()
		if err := application.Shutdown(shutdownCtx); err != nil {
			observability.CaptureError(err)
			log.Error("graceful shutdown failed", zap.Error(err))
		}
	case err := <-errCh:
		if err != nil {
			observability.CaptureError(err)
			log.Error("server error", zap.Error(err))
			os.Exit(1)
		}
	}
}
