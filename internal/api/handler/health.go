package handler

import (
	"net/http"

	"github.com/mailngine/mailngine/internal/api/response"
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

	healthy := true

	if err := h.db.Ping(ctx); err != nil {
		healthy = false
	}

	if err := h.cache.Ping(ctx).Err(); err != nil {
		healthy = false
	}

	status := map[string]string{"status": "healthy"}
	if !healthy {
		status["status"] = "degraded"
	}

	if healthy {
		response.JSON(w, http.StatusOK, status)
	} else {
		response.JSON(w, http.StatusServiceUnavailable, status)
	}
}
