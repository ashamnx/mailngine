package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hibiken/asynq"
	"github.com/hellomail/hellomail/internal/api"
	"github.com/hellomail/hellomail/internal/cache"
	"github.com/hellomail/hellomail/internal/config"
	"github.com/hellomail/hellomail/internal/db"
	"github.com/hellomail/hellomail/internal/observability"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	// Initialize logger
	logger := observability.NewLogger(cfg.Env)
	logger.Info().
		Str("env", cfg.Env).
		Str("port", cfg.Port).
		Msg("starting Hello Mail server")

	// Connect to PostgreSQL
	pool, err := db.NewPool(ctx, cfg.DatabaseURL, logger)
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	defer pool.Close()

	// Connect to Valkey
	valkey, err := cache.NewClient(ctx, cfg.ValkeyURL, logger)
	if err != nil {
		return fmt.Errorf("connect to valkey: %w", err)
	}
	defer valkey.Close()

	// Create asynq client for task queue
	asynqOpt, err := asynq.ParseRedisURI(cfg.ValkeyURL)
	if err != nil {
		return fmt.Errorf("parse valkey url for asynq: %w", err)
	}
	asynqClient := asynq.NewClient(asynqOpt)
	defer asynqClient.Close()
	logger.Info().Msg("connected to asynq task queue")

	// Build HTTP router
	router := api.NewRouter(cfg, pool, valkey, asynqClient, logger)

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	errCh := make(chan error, 1)
	go func() {
		logger.Info().Str("addr", server.Addr).Msg("HTTP server listening")
		errCh <- server.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		if err != http.ErrServerClosed {
			return fmt.Errorf("server error: %w", err)
		}
	case sig := <-shutdown:
		logger.Info().Str("signal", sig.String()).Msg("shutting down")

		shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 30*time.Second)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("graceful shutdown failed: %w", err)
		}
		logger.Info().Msg("server stopped gracefully")
	}

	return nil
}
