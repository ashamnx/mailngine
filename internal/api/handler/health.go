package handler

import (
	"net/http"

	"github.com/hellomail/hellomail/internal/api/response"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type HealthHandler struct {
	db    *pgxpool.Pool
	cache *redis.Client
}

func NewHealthHandler(db *pgxpool.Pool, cache *redis.Client) *HealthHandler {
	return &HealthHandler{db: db, cache: cache}
}

func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	status := map[string]string{
		"status":   "healthy",
		"postgres": "up",
		"valkey":   "up",
	}

	if err := h.db.Ping(ctx); err != nil {
		status["status"] = "degraded"
		status["postgres"] = "down"
	}

	if err := h.cache.Ping(ctx).Err(); err != nil {
		status["status"] = "degraded"
		status["valkey"] = "down"
	}

	if status["status"] == "healthy" {
		response.JSON(w, http.StatusOK, status)
	} else {
		response.JSON(w, http.StatusServiceUnavailable, status)
	}
}
