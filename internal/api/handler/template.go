package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/mailngine/mailngine/internal/api/response"
	"github.com/mailngine/mailngine/internal/auth"
	"github.com/mailngine/mailngine/internal/observability"
	"github.com/mailngine/mailngine/internal/template"
)

// TemplateHandler handles template-related HTTP requests.
type TemplateHandler struct {
	templateSvc *template.Service
}

// NewTemplateHandler creates a new TemplateHandler with the given template service.
func NewTemplateHandler(svc *template.Service) *TemplateHandler {
	return &TemplateHandler{
		templateSvc: svc,
	}
}

type createTemplateRequest struct {
	Name      string   `json:"name" validate:"required"`
	Subject   string   `json:"subject" validate:"required"`
	HTMLBody  string   `json:"html_body" validate:"required"`
	TextBody  string   `json:"text_body"`
	Variables []string `json:"variables"`
}

type updateTemplateRequest struct {
	Name      string   `json:"name" validate:"required"`
	Subject   string   `json:"subject" validate:"required"`
	HTMLBody  string   `json:"html_body" validate:"required"`
	TextBody  string   `json:"text_body"`
	Variables []string `json:"variables"`
}

type previewTemplateRequest struct {
	Data map[string]string `json:"data" validate:"required"`
}

type previewTemplateResponse struct {
	Subject  string `json:"subject"`
	HTMLBody string `json:"html_body"`
	TextBody string `json:"text_body"`
}

// Create handles POST /v1/templates.
func (h *TemplateHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)

	var req createTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if err := validate.Struct(req); err != nil {
		response.BadRequest(w, formatValidationError(err))
		return
	}

	tmpl, err := h.templateSvc.Create(ctx, orgID, req.Name, req.Subject, req.HTMLBody, req.TextBody, req.Variables)
	if err != nil {
		if errors.Is(err, template.ErrTemplateNameConflict) {
			response.Conflict(w, "template name already exists")
			return
		}
		logger.Error().Err(err).Msg("failed to create template")
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusCreated, tmpl)
}

// List handles GET /v1/templates.
func (h *TemplateHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)

	templates, err := h.templateSvc.List(ctx, orgID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to list templates")
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, templates)
}

// Get handles GET /v1/templates/{id}.
func (h *TemplateHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid template ID")
		return
	}

	tmpl, err := h.templateSvc.Get(ctx, orgID, id)
	if err != nil {
		if errors.Is(err, template.ErrTemplateNotFound) {
			response.NotFound(w, "template not found")
			return
		}
		logger.Error().Err(err).Str("template_id", id.String()).Msg("failed to get template")
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, tmpl)
}

// Update handles PATCH /v1/templates/{id}.
func (h *TemplateHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid template ID")
		return
	}

	var req updateTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if err := validate.Struct(req); err != nil {
		response.BadRequest(w, formatValidationError(err))
		return
	}

	tmpl, err := h.templateSvc.Update(ctx, orgID, id, req.Name, req.Subject, req.HTMLBody, req.TextBody, req.Variables)
	if err != nil {
		switch {
		case errors.Is(err, template.ErrTemplateNotFound):
			response.NotFound(w, "template not found")
		case errors.Is(err, template.ErrTemplateNameConflict):
			response.Conflict(w, "template name already exists")
		default:
			logger.Error().Err(err).Str("template_id", id.String()).Msg("failed to update template")
			response.InternalError(w)
		}
		return
	}

	response.JSON(w, http.StatusOK, tmpl)
}

// Delete handles DELETE /v1/templates/{id}.
func (h *TemplateHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid template ID")
		return
	}

	if err := h.templateSvc.Delete(ctx, orgID, id); err != nil {
		logger.Error().Err(err).Str("template_id", id.String()).Msg("failed to delete template")
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "template deleted"})
}

// Preview handles POST /v1/templates/{id}/preview.
func (h *TemplateHandler) Preview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid template ID")
		return
	}

	var req previewTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	subject, html, text, err := h.templateSvc.Render(ctx, orgID, id, req.Data)
	if err != nil {
		if errors.Is(err, template.ErrTemplateNotFound) {
			response.NotFound(w, "template not found")
			return
		}
		logger.Error().Err(err).Str("template_id", id.String()).Msg("failed to preview template")
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, previewTemplateResponse{
		Subject:  subject,
		HTMLBody: html,
		TextBody: text,
	})
}
