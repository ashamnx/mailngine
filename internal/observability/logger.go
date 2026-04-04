package observability

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type ctxKey string

const (
	correlationIDKey ctxKey = "correlation_id"
	loggerKey        ctxKey = "logger"
)

func NewLogger(env string) zerolog.Logger {
	var output io.Writer = os.Stdout

	if env == "development" {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
	}

	return zerolog.New(output).
		With().
		Timestamp().
		Caller().
		Logger().
		Level(zerolog.InfoLevel)
}

func WithCorrelationID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, correlationIDKey, id)
}

func CorrelationIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(correlationIDKey).(string); ok {
		return id
	}
	return ""
}

func WithLogger(ctx context.Context, logger zerolog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func LoggerFromContext(ctx context.Context) zerolog.Logger {
	if logger, ok := ctx.Value(loggerKey).(zerolog.Logger); ok {
		return logger
	}
	return zerolog.Nop()
}
