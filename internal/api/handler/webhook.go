package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/hellomail/hellomail/internal/api/response"
	"github.com/hellomail/hellomail/internal/auth"
	"github.com/hellomail/hellomail/internal/observability"
	"github.com/hellomail/hellomail/internal/webhook"
)

type WebhookHandler struct {
	webhookSvc *webhook.Service
}

func NewWebhookHandler(svc *webhook.Service) *WebhookHandler {
	return &WebhookHandler{webhookSvc: svc}
}

type createWebhookRequest struct {
	URL    string   `json:"url" validate:"required,url"`
	Events []string `json:"events" validate:"required,min=1"`
}

func (h *WebhookHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)

	var req createWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	wh, err := h.webhookSvc.Create(ctx, orgID, req.URL, req.Events)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create webhook")
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusCreated, wh)
}

func (h *WebhookHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)

	webhooks, err := h.webhookSvc.List(ctx, orgID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to list webhooks")
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusOK, webhooks)
}

func (h *WebhookHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	orgID := auth.OrgIDFromContext(ctx)
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid webhook ID")
		return
	}

	wh, err := h.webhookSvc.Get(ctx, orgID, id)
	if err != nil {
		response.NotFound(w, "webhook not found")
		return
	}
	response.JSON(w, http.StatusOK, wh)
}

type updateWebhookRequest struct {
	URL      string   `json:"url" validate:"required,url"`
	Events   []string `json:"events" validate:"required,min=1"`
	IsActive bool     `json:"is_active"`
}

func (h *WebhookHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid webhook ID")
		return
	}

	var req updateWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	wh, err := h.webhookSvc.Update(ctx, orgID, id, req.URL, req.Events, req.IsActive)
	if err != nil {
		logger.Error().Err(err).Msg("failed to update webhook")
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusOK, wh)
}

func (h *WebhookHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid webhook ID")
		return
	}

	if err := h.webhookSvc.Delete(ctx, orgID, id); err != nil {
		logger.Error().Err(err).Msg("failed to delete webhook")
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "webhook deleted"})
}

func (h *WebhookHandler) ListDeliveries(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid webhook ID")
		return
	}

	page := parseQueryInt(r, "page", 1)
	perPage := parseQueryInt(r, "per_page", 20)

	deliveries, err := h.webhookSvc.ListDeliveries(ctx, id, page, perPage)
	if err != nil {
		logger.Error().Err(err).Msg("failed to list deliveries")
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusOK, deliveries)
}
