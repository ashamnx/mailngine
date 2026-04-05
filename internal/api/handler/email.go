package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/mailngine/mailngine/internal/api/response"
	"github.com/mailngine/mailngine/internal/auth"
	"github.com/mailngine/mailngine/internal/email"
	"github.com/mailngine/mailngine/internal/observability"
)

// validate is a package-level validator instance, safe for concurrent use.
var validate = validator.New()

// EmailHandler handles email-related HTTP requests.
type EmailHandler struct {
	emailSvc *email.Service
}

// NewEmailHandler creates a new EmailHandler with the given email service.
func NewEmailHandler(emailSvc *email.Service) *EmailHandler {
	return &EmailHandler{
		emailSvc: emailSvc,
	}
}

// Send handles POST /v1/emails.
// It decodes the request, validates it, and delegates to the email service.
func (h *EmailHandler) Send(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)
	apiKeyID := auth.APIKeyIDFromContext(ctx)

	var req email.SendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if err := validate.Struct(req); err != nil {
		response.BadRequest(w, formatValidationError(err))
		return
	}

	// At least one of HTML or Text must be provided.
	if req.HTML == "" && req.Text == "" {
		response.BadRequest(w, "at least one of 'html' or 'text' is required")
		return
	}

	result, err := h.emailSvc.SendEmail(ctx, orgID, apiKeyID, &req)
	if err != nil {
		switch {
		case errors.Is(err, email.ErrDomainNotVerified):
			response.BadRequest(w, err.Error())
		case errors.Is(err, email.ErrInvalidFromAddress):
			response.BadRequest(w, err.Error())
		default:
			logger.Error().Err(err).Msg("failed to send email")
			response.InternalError(w)
		}
		return
	}

	resp := email.ToResponse(result)
	response.JSON(w, http.StatusCreated, resp)
}

// Get handles GET /v1/emails/{id}.
func (h *EmailHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)

	idParam := chi.URLParam(r, "id")
	emailID, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(w, "invalid email id")
		return
	}

	result, err := h.emailSvc.GetEmail(ctx, orgID, emailID)
	if err != nil {
		if errors.Is(err, email.ErrEmailNotFound) {
			response.NotFound(w, "email not found")
			return
		}
		logger.Error().Err(err).Str("email_id", idParam).Msg("failed to get email")
		response.InternalError(w)
		return
	}

	resp := email.ToResponse(result)
	response.JSON(w, http.StatusOK, resp)
}

// List handles GET /v1/emails?page=1&per_page=20.
func (h *EmailHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)

	page := parseQueryInt(r, "page", 1)
	perPage := parseQueryInt(r, "per_page", 20)

	emails, total, err := h.emailSvc.ListEmails(ctx, orgID, page, perPage)
	if err != nil {
		logger.Error().Err(err).Msg("failed to list emails")
		response.InternalError(w)
		return
	}

	items := make([]email.EmailResponse, len(emails))
	for i := range emails {
		items[i] = email.ToResponse(&emails[i])
	}

	response.JSONWithMeta(w, http.StatusOK, items, &response.Meta{
		Page:    page,
		PerPage: perPage,
		Total:   int(total),
		HasMore: int64(page*perPage) < total,
	})
}

// parseQueryInt extracts an integer query parameter with a default fallback.
func parseQueryInt(r *http.Request, key string, fallback int) int {
	v := r.URL.Query().Get(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil || n < 1 {
		return fallback
	}
	return n
}

// formatValidationError converts validator errors into a human-readable message.
func formatValidationError(err error) string {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) && len(validationErrors) > 0 {
		fe := validationErrors[0]
		switch fe.Tag() {
		case "required":
			return fe.Field() + " is required"
		case "email":
			return fe.Field() + " must be a valid email address"
		case "min":
			return fe.Field() + " must have at least " + fe.Param() + " item(s)"
		default:
			return fe.Field() + " is invalid"
		}
	}
	return "validation failed"
}
