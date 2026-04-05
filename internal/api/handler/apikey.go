package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/mailngine/mailngine/internal/api/response"
	"github.com/mailngine/mailngine/internal/auth"
	sqlcdb "github.com/mailngine/mailngine/internal/db/sqlcdb"
	"github.com/mailngine/mailngine/internal/observability"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// APIKeyHandler handles API key management HTTP requests.
type APIKeyHandler struct {
	db      *pgxpool.Pool
	queries *sqlcdb.Queries
}

// NewAPIKeyHandler creates a new APIKeyHandler with the given dependencies.
func NewAPIKeyHandler(db *pgxpool.Pool) *APIKeyHandler {
	return &APIKeyHandler{
		db:      db,
		queries: sqlcdb.New(db),
	}
}

// CreateAPIKeyRequest represents the request body for creating an API key.
type CreateAPIKeyRequest struct {
	Name       string     `json:"name"`
	Permission string     `json:"permission"`
	DomainID   *uuid.UUID `json:"domain_id,omitempty"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
}

// CreateAPIKeyResponse represents the response after creating an API key.
// The full key is only returned once at creation time.
type CreateAPIKeyResponse struct {
	ID         uuid.UUID  `json:"id"`
	Name       string     `json:"name"`
	Prefix     string     `json:"prefix"`
	Key        string     `json:"key"`
	Permission string     `json:"permission"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

// APIKeyListItem represents an API key in list responses (without the full key).
type APIKeyListItem struct {
	ID         uuid.UUID  `json:"id"`
	Name       string     `json:"name"`
	Prefix     string     `json:"prefix"`
	Permission string     `json:"permission"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

// Create generates a new API key for the authenticated user's organization.
// POST /v1/api-keys
func (h *APIKeyHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)

	var req CreateAPIKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if req.Name == "" {
		response.BadRequest(w, "name is required")
		return
	}

	if req.Permission == "" {
		req.Permission = "full"
	}

	validPermissions := map[string]bool{"full": true, "send_only": true, "read_only": true}
	if !validPermissions[req.Permission] {
		response.BadRequest(w, "permission must be one of: full, send_only, read_only")
		return
	}

	// Generate API key
	fullKey, prefix, hash, err := auth.GenerateAPIKey()
	if err != nil {
		logger.Error().Err(err).Msg("failed to generate api key")
		response.InternalError(w)
		return
	}

	// Build domain ID parameter
	var domainID pgtype.UUID
	if req.DomainID != nil {
		domainID = pgtype.UUID{Bytes: *req.DomainID, Valid: true}
	}

	// Build expires_at parameter
	var expiresAt pgtype.Timestamptz
	if req.ExpiresAt != nil {
		expiresAt = pgtype.Timestamptz{Time: *req.ExpiresAt, Valid: true}
	}

	apiKey, err := h.queries.CreateAPIKey(ctx, sqlcdb.CreateAPIKeyParams{
		OrgID:      orgID,
		DomainID:   domainID,
		Name:       req.Name,
		Prefix:     prefix,
		KeyHash:    hash,
		Permission: req.Permission,
		ExpiresAt:  expiresAt,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to create api key")
		response.InternalError(w)
		return
	}

	var respExpiresAt *time.Time
	if apiKey.ExpiresAt.Valid {
		respExpiresAt = &apiKey.ExpiresAt.Time
	}

	resp := CreateAPIKeyResponse{
		ID:         apiKey.ID,
		Name:       apiKey.Name,
		Prefix:     apiKey.Prefix,
		Key:        fullKey,
		Permission: apiKey.Permission,
		ExpiresAt:  respExpiresAt,
		CreatedAt:  apiKey.CreatedAt,
	}

	response.JSON(w, http.StatusCreated, resp)
}

// List returns all active API keys for the authenticated user's organization.
// GET /v1/api-keys
func (h *APIKeyHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)

	keys, err := h.queries.ListAPIKeysByOrg(ctx, orgID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to list api keys")
		response.InternalError(w)
		return
	}

	items := make([]APIKeyListItem, len(keys))
	for i, k := range keys {
		var lastUsedAt *time.Time
		if k.LastUsedAt.Valid {
			lastUsedAt = &k.LastUsedAt.Time
		}
		var expiresAt *time.Time
		if k.ExpiresAt.Valid {
			expiresAt = &k.ExpiresAt.Time
		}

		items[i] = APIKeyListItem{
			ID:         k.ID,
			Name:       k.Name,
			Prefix:     k.Prefix,
			Permission: k.Permission,
			LastUsedAt: lastUsedAt,
			ExpiresAt:  expiresAt,
			CreatedAt:  k.CreatedAt,
		}
	}

	response.JSON(w, http.StatusOK, items)
}

// Revoke deactivates an API key by setting its revoked_at timestamp.
// DELETE /v1/api-keys/{id}
func (h *APIKeyHandler) Revoke(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)

	idParam := chi.URLParam(r, "id")
	keyID, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(w, "invalid api key id")
		return
	}

	// Verify the key belongs to this org
	_, err = h.queries.GetAPIKey(ctx, sqlcdb.GetAPIKeyParams{
		ID:    keyID,
		OrgID: orgID,
	})
	if err != nil {
		logger.Warn().Err(err).Str("key_id", idParam).Msg("api key not found")
		response.NotFound(w, "api key not found")
		return
	}

	if err := h.queries.RevokeAPIKey(ctx, sqlcdb.RevokeAPIKeyParams{
		ID:    keyID,
		OrgID: orgID,
	}); err != nil {
		logger.Error().Err(err).Msg("failed to revoke api key")
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "api key revoked"})
}
