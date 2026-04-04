package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hibiken/asynq"
	"github.com/hellomail/hellomail/internal/config"
	"github.com/hellomail/hellomail/internal/observability"
	"github.com/rs/zerolog"
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

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	logger := observability.NewLogger(cfg.Env)
	logger.Info().Msg("starting Hello Mail worker")

	redisOpt, err := asynq.ParseRedisURI(cfg.ValkeyURL)
	if err != nil {
		return fmt.Errorf("parse valkey url for asynq: %w", err)
	}

	srv := asynq.NewServer(redisOpt, asynq.Config{
		Concurrency: 10,
		Queues: map[string]int{
			"critical": 6,
			"default":  3,
			"low":      1,
		},
		Logger: &asynqLogger{logger: logger},
	})

	mux := asynq.NewServeMux()
	// Task handlers will be registered here in Phase 2+

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Run(mux)
	}()

	_ = ctx // used in future phases for DB/cache connections

	select {
	case err := <-errCh:
		return fmt.Errorf("worker error: %w", err)
	case sig := <-shutdown:
		logger.Info().Str("signal", sig.String()).Msg("shutting down worker")
		srv.Shutdown()
	}

	return nil
}

type asynqLogger struct {
	logger zerolog.Logger
}

func (l *asynqLogger) Debug(args ...any)                    { l.logger.Debug().Msgf("%v", args) }
func (l *asynqLogger) Info(args ...any)                     { l.logger.Info().Msgf("%v", args) }
func (l *asynqLogger) Warn(args ...any)                     { l.logger.Warn().Msgf("%v", args) }
func (l *asynqLogger) Error(args ...any)                    { l.logger.Error().Msgf("%v", args) }
func (l *asynqLogger) Fatal(args ...any)                    { l.logger.Fatal().Msgf("%v", args) }
