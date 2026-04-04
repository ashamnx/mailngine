package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/hellomail/hellomail/internal/api/response"
	"github.com/hellomail/hellomail/internal/auth"
	"github.com/hellomail/hellomail/internal/observability"
	"github.com/redis/go-redis/v9"
)

// RateLimitConfig defines the rate limiting parameters.
type RateLimitConfig struct {
	// APIKeyRate is the maximum number of requests per window for API key auth.
	APIKeyRate int
	// SessionRate is the maximum number of requests per window for JWT session auth.
	SessionRate int
	// Window is the duration of the rate limit window.
	Window time.Duration
}

// RateLimit returns middleware that enforces per-identity rate limits using Valkey.
// It uses a fixed-window counter per second with INCR and EXPIRE.
func RateLimit(cache *redis.Client, cfg RateLimitConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			logger := observability.LoggerFromContext(ctx)

			var identity string
			var limit int

			authType := auth.AuthTypeFromContext(ctx)
			switch authType {
			case auth.AuthTypeAPIKey:
				keyID := auth.APIKeyIDFromContext(ctx)
				if keyID != nil {
					identity = "rl:apikey:" + keyID.String()
				}
				limit = cfg.APIKeyRate
			case auth.AuthTypeJWT:
				userID := auth.UserIDFromContext(ctx)
				identity = "rl:session:" + userID.String()
				limit = cfg.SessionRate
			default:
				// No auth context; skip rate limiting
				next.ServeHTTP(w, r)
				return
			}

			if identity == "" {
				next.ServeHTTP(w, r)
				return
			}

			windowSeconds := int(cfg.Window.Seconds())
			if windowSeconds < 1 {
				windowSeconds = 1
			}

			now := time.Now().Unix()
			windowKey := fmt.Sprintf("%s:%d", identity, now/int64(windowSeconds))

			count, err := cache.Incr(ctx, windowKey).Result()
			if err != nil {
				logger.Error().Err(err).Msg("rate limit incr failed")
				// Fail open: allow the request if Valkey is unavailable
				next.ServeHTTP(w, r)
				return
			}

			// Set expiry on first request in the window
			if count == 1 {
				cache.Expire(ctx, windowKey, cfg.Window+time.Second)
			}

			remaining := limit - int(count)
			if remaining < 0 {
				remaining = 0
			}

			resetTime := ((now / int64(windowSeconds)) + 1) * int64(windowSeconds)

			// Set rate limit headers
			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limit))
			w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
			w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(resetTime, 10))

			if int(count) > limit {
				response.TooManyRequests(w, "rate limit exceeded")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
