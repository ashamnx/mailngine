package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/mailngine/mailngine/internal/api/response"
	"github.com/mailngine/mailngine/internal/auth"
	sqlcdb "github.com/mailngine/mailngine/internal/db/sqlcdb"
	"github.com/mailngine/mailngine/internal/observability"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

const (
	sessionPrefix = "session:"
)

// Authenticate returns middleware that validates JWT tokens or API keys on each request.
// For JWT auth, it also verifies an active session exists in Valkey.
// For API key auth, it updates the last_used_at timestamp asynchronously.
func Authenticate(jwtMgr *auth.JWTManager, db *pgxpool.Pool, cache *redis.Client) func(http.Handler) http.Handler {
	queries := sqlcdb.New(db)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := extractToken(r)
			if token == "" {
				response.Unauthorized(w, "missing authentication token")
				return
			}

			ctx := r.Context()
			logger := observability.LoggerFromContext(ctx)

			// API key authentication path
			if strings.HasPrefix(token, "mn_") {
				hash := auth.HashAPIKey(token)
				apiKey, err := queries.GetAPIKeyByHash(ctx, hash)
				if err != nil {
					logger.Warn().Err(err).Msg("invalid api key")
					response.Unauthorized(w, "invalid api key")
					return
				}

				// Determine role from API key permission
				role := permissionToRole(apiKey.Permission)

				keyID := apiKey.ID
				ctx = auth.WithAuthContext(ctx, apiKey.OrgID, apiKey.OrgID, role, auth.AuthTypeAPIKey, &keyID)

				// Update last_used_at asynchronously
				go func() {
					_ = queries.UpdateAPIKeyLastUsed(context.Background(), apiKey.ID)
				}()

				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// JWT authentication path
			claims, err := jwtMgr.Validate(token)
			if err != nil {
				logger.Warn().Err(err).Msg("invalid jwt token")
				response.Unauthorized(w, "invalid or expired token")
				return
			}

			// Verify session exists in Valkey
			sessionKey := sessionPrefix + claims.UserID.String()
			exists, err := cache.Exists(ctx, sessionKey).Result()
			if err != nil {
				logger.Error().Err(err).Msg("failed to check session in valkey")
				response.InternalError(w)
				return
			}
			if exists == 0 {
				response.Unauthorized(w, "session expired")
				return
			}

			ctx = auth.WithAuthContext(ctx, claims.UserID, claims.OrgID, claims.Role, auth.AuthTypeJWT, nil)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// extractToken retrieves the authentication token from the Authorization header
// or falls back to the "token" cookie.
func extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
			return strings.TrimSpace(parts[1])
		}
	}

	cookie, err := r.Cookie("token")
	if err == nil && cookie.Value != "" {
		return cookie.Value
	}

	return ""
}

// permissionToRole maps API key permission levels to organization roles.
func permissionToRole(permission string) string {
	switch permission {
	case "full":
		return "owner"
	case "send_only":
		return "member"
	case "read_only":
		return "viewer"
	default:
		return "viewer"
	}
}
