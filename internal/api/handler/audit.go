package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/hellomail/hellomail/internal/api/response"
	"github.com/hellomail/hellomail/internal/audit"
	"github.com/hellomail/hellomail/internal/auth"
	"github.com/hellomail/hellomail/internal/observability"
)

type AuditHandler struct {
	auditSvc *audit.Service
}

func NewAuditHandler(svc *audit.Service) *AuditHandler {
	return &AuditHandler{auditSvc: svc}
}

func (h *AuditHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)

	page := parseQueryInt(r, "page", 1)
	perPage := parseQueryInt(r, "per_page", 20)

	logs, total, err := h.auditSvc.List(ctx, orgID, page, perPage)
	if err != nil {
		logger.Error().Err(err).Msg("failed to list audit logs")
		response.InternalError(w)
		return
	}

	response.JSONWithMeta(w, http.StatusOK, logs, &response.Meta{
		Page:    page,
		PerPage: perPage,
		Total:   int(total),
	})
}

func (h *AuditHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	orgID := auth.OrgIDFromContext(ctx)
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid audit log ID")
		return
	}

	log, err := h.auditSvc.Get(ctx, orgID, id)
	if err != nil {
		response.NotFound(w, "audit log not found")
		return
	}
	response.JSON(w, http.StatusOK, log)
}
