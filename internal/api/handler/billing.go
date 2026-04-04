package handler

import (
	"net/http"

	"github.com/hellomail/hellomail/internal/api/response"
	"github.com/hellomail/hellomail/internal/auth"
	"github.com/hellomail/hellomail/internal/billing"
	"github.com/hellomail/hellomail/internal/observability"
)

// BillingHandler handles billing-related HTTP requests.
type BillingHandler struct {
	billingSvc *billing.Service
}

// NewBillingHandler creates a new BillingHandler with the given billing service.
func NewBillingHandler(svc *billing.Service) *BillingHandler {
	return &BillingHandler{
		billingSvc: svc,
	}
}

// Usage handles GET /v1/billing/usage.
func (h *BillingHandler) Usage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)

	usage, err := h.billingSvc.GetCurrentUsage(ctx, orgID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get billing usage")
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, usage)
}

// History handles GET /v1/billing/usage/history?limit=12.
func (h *BillingHandler) History(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)

	limit := parseQueryInt(r, "limit", 12)

	history, err := h.billingSvc.GetUsageHistory(ctx, orgID, limit)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get billing usage history")
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, history)
}

// Plan handles GET /v1/billing/plan.
func (h *BillingHandler) Plan(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)

	plan, err := h.billingSvc.GetPlan(ctx, orgID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get billing plan")
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, plan)
}
