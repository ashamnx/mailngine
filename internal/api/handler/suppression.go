package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/hellomail/hellomail/internal/api/response"
	"github.com/hellomail/hellomail/internal/auth"
	sqlcdb "github.com/hellomail/hellomail/internal/db/sqlcdb"
	"github.com/hellomail/hellomail/internal/observability"
	"github.com/hellomail/hellomail/internal/suppression"
)

// SuppressionHandler handles suppression list HTTP requests.
type SuppressionHandler struct {
	suppressionSvc *suppression.Service
}

// NewSuppressionHandler creates a new SuppressionHandler with the given service.
func NewSuppressionHandler(svc *suppression.Service) *SuppressionHandler {
	return &SuppressionHandler{
		suppressionSvc: svc,
	}
}

// createSuppressionRequest represents the request body for manually adding a suppression.
type createSuppressionRequest struct {
	Email  string `json:"email" validate:"required,email"`
	Reason string `json:"reason" validate:"required"`
}

// suppressionResponse represents a suppression in API responses.
type suppressionResponse struct {
	ID        uuid.UUID       `json:"id"`
	OrgID     uuid.UUID       `json:"org_id"`
	Email     string          `json:"email"`
	Reason    string          `json:"reason"`
	Metadata  json.RawMessage `json:"metadata,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
}

// List handles GET /v1/suppressions?page=1&per_page=20.
func (h *SuppressionHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)

	page := parseQueryInt(r, "page", 1)
	perPage := parseQueryInt(r, "per_page", 20)

	suppressions, total, err := h.suppressionSvc.List(ctx, orgID, page, perPage)
	if err != nil {
		logger.Error().Err(err).Msg("failed to list suppressions")
		response.InternalError(w)
		return
	}

	items := make([]suppressionResponse, len(suppressions))
	for i := range suppressions {
		items[i] = toSuppressionResponse(suppressions[i])
	}

	response.JSONWithMeta(w, http.StatusOK, items, &response.Meta{
		Page:    page,
		PerPage: perPage,
		Total:   int(total),
		HasMore: int64(page*perPage) < total,
	})
}

// Create handles POST /v1/suppressions for manually adding an address.
func (h *SuppressionHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)

	var req createSuppressionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if err := validate.Struct(req); err != nil {
		response.BadRequest(w, formatValidationError(err))
		return
	}

	if err := h.suppressionSvc.Add(ctx, orgID, req.Email, req.Reason, map[string]any{
		"source": "manual",
	}); err != nil {
		logger.Error().Err(err).Str("email", req.Email).Msg("failed to create suppression")
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusCreated, map[string]string{
		"message": "email suppressed",
		"email":   req.Email,
	})
}

// Delete handles DELETE /v1/suppressions/{id}.
func (h *SuppressionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)

	idParam := chi.URLParam(r, "id")
	suppressionID, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(w, "invalid suppression id")
		return
	}

	if err := h.suppressionSvc.Remove(ctx, orgID, suppressionID); err != nil {
		if errors.Is(err, suppression.ErrSuppressionNotFound) {
			response.NotFound(w, "suppression not found")
			return
		}
		logger.Error().Err(err).Str("suppression_id", idParam).Msg("failed to delete suppression")
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "suppression removed"})
}

// toSuppressionResponse converts a database Suppression model to an API response.
func toSuppressionResponse(s sqlcdb.Suppression) suppressionResponse {
	resp := suppressionResponse{
		ID:        s.ID,
		OrgID:     s.OrgID,
		Email:     s.Email,
		Reason:    s.Reason,
		CreatedAt: s.CreatedAt,
	}
	if len(s.Metadata) > 0 {
		resp.Metadata = s.Metadata
	}
	return resp
}
