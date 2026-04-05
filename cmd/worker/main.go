package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hibiken/asynq"
	"github.com/mailngine/mailngine/internal/config"
	"github.com/mailngine/mailngine/internal/db"
	sqlcdb "github.com/mailngine/mailngine/internal/db/sqlcdb"
	"github.com/mailngine/mailngine/internal/email"
	"github.com/mailngine/mailngine/internal/inbox"
	"github.com/mailngine/mailngine/internal/observability"
	internalsmtp "github.com/mailngine/mailngine/internal/smtp"
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
	logger.Info().Msg("starting Mailngine worker")

	// Connect to PostgreSQL.
	pool, err := db.NewPool(ctx, cfg.DatabaseURL, logger)
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	defer pool.Close()

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
		ShutdownTimeout: cfg.WorkerShutdownTimeout,
		Logger:          &asynqLogger{logger: logger},
	})

	logger.Info().
		Dur("shutdown_timeout", cfg.WorkerShutdownTimeout).
		Msg("worker configured")

	mux := asynq.NewServeMux()

	// Register task handlers.
	smtpSender := email.NewSMTPSender(cfg.Postfix.SMTPHost, cfg.Postfix.SMTPPort)
	sendHandler := email.NewSendHandler(pool, smtpSender, logger)
	sendHandler.Register(mux)

	logger.Info().Msg("task handlers registered")

	// Start LMTP server for inbound email.
	queries := sqlcdb.New(pool)
	inboxSvc := inbox.NewService(pool, logger)
	lmtpServer := internalsmtp.NewLMTPServer(cfg.LMTPSocketPath, queries, inboxSvc, logger)

	errCh := make(chan error, 2)

	go func() {
		errCh <- lmtpServer.ListenAndServe()
	}()

	go func() {
		errCh <- srv.Run(mux)
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return fmt.Errorf("worker error: %w", err)
	case sig := <-shutdown:
		logger.Info().Str("signal", sig.String()).Msg("shutting down worker")
		lmtpServer.Close()
		srv.Shutdown()
	}

	return nil
}

type asynqLogger struct {
	logger zerolog.Logger
}

func (l *asynqLogger) Debug(args ...any) { l.logger.Debug().Msgf("%v", args) }
func (l *asynqLogger) Info(args ...any)  { l.logger.Info().Msgf("%v", args) }
func (l *asynqLogger) Warn(args ...any)  { l.logger.Warn().Msgf("%v", args) }
func (l *asynqLogger) Error(args ...any) { l.logger.Error().Msgf("%v", args) }
func (l *asynqLogger) Fatal(args ...any) { l.logger.Fatal().Msgf("%v", args) }
