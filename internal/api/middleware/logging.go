package middleware

import (
	"net/http"
	"time"

	"github.com/mailngine/mailngine/internal/observability"
	"github.com/rs/zerolog"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func Logging(logger zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			correlationID := observability.CorrelationIDFromContext(r.Context())

			reqLogger := logger.With().
				Str("correlation_id", correlationID).
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Str("remote_addr", r.RemoteAddr).
				Logger()

			ctx := observability.WithLogger(r.Context(), reqLogger)

			wrapped := &wrappedWriter{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(wrapped, r.WithContext(ctx))

			duration := time.Since(start)

			reqLogger.Info().
				Int("status", wrapped.statusCode).
				Dur("duration", duration).
				Msg("request completed")
		})
	}
}
