package handler

import (
	"net/http"
	"time"

	"github.com/hellomail/hellomail/internal/analytics"
	"github.com/hellomail/hellomail/internal/api/response"
	"github.com/hellomail/hellomail/internal/auth"
	"github.com/hellomail/hellomail/internal/observability"
)

type AnalyticsHandler struct {
	analyticsSvc *analytics.Service
}

func NewAnalyticsHandler(svc *analytics.Service) *AnalyticsHandler {
	return &AnalyticsHandler{analyticsSvc: svc}
}

func (h *AnalyticsHandler) Overview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)
	from, to := parseDateRange(r)

	overview, err := h.analyticsSvc.GetOverview(ctx, orgID, from, to)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get analytics overview")
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusOK, overview)
}

func (h *AnalyticsHandler) Timeseries(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)
	from, to := parseDateRange(r)

	points, err := h.analyticsSvc.GetTimeseries(ctx, orgID, from, to)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get timeseries")
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusOK, points)
}

func (h *AnalyticsHandler) Events(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)
	from, to := parseDateRange(r)

	breakdown, err := h.analyticsSvc.GetEventBreakdown(ctx, orgID, from, to)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get event breakdown")
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusOK, breakdown)
}

func parseDateRange(r *http.Request) (time.Time, time.Time) {
	now := time.Now().UTC()
	from := now.AddDate(0, 0, -30)
	to := now

	if f := r.URL.Query().Get("from"); f != "" {
		if t, err := time.Parse("2006-01-02", f); err == nil {
			from = t
		}
	}
	if t := r.URL.Query().Get("to"); t != "" {
		if parsed, err := time.Parse("2006-01-02", t); err == nil {
			to = parsed
		}
	}
	return from, to
}
