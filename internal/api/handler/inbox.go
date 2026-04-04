package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/hellomail/hellomail/internal/api/response"
	"github.com/hellomail/hellomail/internal/auth"
	"github.com/hellomail/hellomail/internal/inbox"
	"github.com/hellomail/hellomail/internal/observability"
)

// InboxHandler handles inbox-related HTTP requests for threads, messages, and labels.
type InboxHandler struct {
	inboxSvc *inbox.Service
}

// NewInboxHandler creates a new InboxHandler with the given inbox service.
func NewInboxHandler(inboxSvc *inbox.Service) *InboxHandler {
	return &InboxHandler{
		inboxSvc: inboxSvc,
	}
}

// ListThreads handles GET /v1/inbox/threads.
func (h *InboxHandler) ListThreads(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)

	page := parseQueryInt(r, "page", 1)
	perPage := parseQueryInt(r, "per_page", 20)

	threads, total, err := h.inboxSvc.ListThreads(ctx, orgID, page, perPage)
	if err != nil {
		logger.Error().Err(err).Msg("failed to list threads")
		response.InternalError(w)
		return
	}

	response.JSONWithMeta(w, http.StatusOK, threads, &response.Meta{
		Page:    page,
		PerPage: perPage,
		Total:   int(total),
		HasMore: int64(page*perPage) < total,
	})
}

// GetThread handles GET /v1/inbox/threads/{id}.
func (h *InboxHandler) GetThread(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)

	threadID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid thread id")
		return
	}

	thread, messages, err := h.inboxSvc.GetThread(ctx, orgID, threadID)
	if err != nil {
		if errors.Is(err, inbox.ErrThreadNotFound) {
			response.NotFound(w, "thread not found")
			return
		}
		logger.Error().Err(err).Str("thread_id", threadID.String()).Msg("failed to get thread")
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"thread":   thread,
		"messages": messages,
	})
}

// DeleteThread handles DELETE /v1/inbox/threads/{id}.
func (h *InboxHandler) DeleteThread(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)

	threadID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid thread id")
		return
	}

	if err := h.inboxSvc.DeleteThread(ctx, orgID, threadID); err != nil {
		if errors.Is(err, inbox.ErrThreadNotFound) {
			response.NotFound(w, "thread not found")
			return
		}
		logger.Error().Err(err).Str("thread_id", threadID.String()).Msg("failed to delete thread")
		response.InternalError(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetMessage handles GET /v1/inbox/messages/{id}.
func (h *InboxHandler) GetMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)

	messageID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid message id")
		return
	}

	msg, err := h.inboxSvc.GetMessage(ctx, orgID, messageID)
	if err != nil {
		if errors.Is(err, inbox.ErrMessageNotFound) {
			response.NotFound(w, "message not found")
			return
		}
		logger.Error().Err(err).Str("message_id", messageID.String()).Msg("failed to get message")
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, msg)
}

// UpdateMessageRequest represents the request body for PATCH /v1/inbox/messages/{id}.
type UpdateMessageRequest struct {
	IsRead     *bool `json:"is_read,omitempty"`
	IsStarred  *bool `json:"is_starred,omitempty"`
	IsArchived *bool `json:"is_archived,omitempty"`
	IsTrashed  *bool `json:"is_trashed,omitempty"`
}

// UpdateMessage handles PATCH /v1/inbox/messages/{id}.
func (h *InboxHandler) UpdateMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)

	messageID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid message id")
		return
	}

	var req UpdateMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	// Require at least one field to update.
	if req.IsRead == nil && req.IsStarred == nil && req.IsArchived == nil && req.IsTrashed == nil {
		response.BadRequest(w, "at least one flag must be provided")
		return
	}

	if err := h.inboxSvc.UpdateMessageFlags(ctx, orgID, messageID, req.IsRead, req.IsStarred, req.IsArchived, req.IsTrashed); err != nil {
		if errors.Is(err, inbox.ErrMessageNotFound) {
			response.NotFound(w, "message not found")
			return
		}
		logger.Error().Err(err).Str("message_id", messageID.String()).Msg("failed to update message flags")
		response.InternalError(w)
		return
	}

	// Return the updated message.
	msg, err := h.inboxSvc.GetMessage(ctx, orgID, messageID)
	if err != nil {
		logger.Error().Err(err).Str("message_id", messageID.String()).Msg("failed to get updated message")
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, msg)
}

// DeleteMessage handles DELETE /v1/inbox/messages/{id}.
func (h *InboxHandler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)

	messageID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid message id")
		return
	}

	if err := h.inboxSvc.DeleteMessage(ctx, orgID, messageID); err != nil {
		if errors.Is(err, inbox.ErrMessageNotFound) {
			response.NotFound(w, "message not found")
			return
		}
		logger.Error().Err(err).Str("message_id", messageID.String()).Msg("failed to delete message")
		response.InternalError(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// SearchMessages handles GET /v1/inbox/search?q=...
func (h *InboxHandler) SearchMessages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)

	query := r.URL.Query().Get("q")
	if query == "" {
		response.BadRequest(w, "search query 'q' is required")
		return
	}

	page := parseQueryInt(r, "page", 1)
	perPage := parseQueryInt(r, "per_page", 20)

	messages, err := h.inboxSvc.SearchMessages(ctx, orgID, query, page, perPage)
	if err != nil {
		logger.Error().Err(err).Str("query", query).Msg("failed to search messages")
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, messages)
}

// ListLabels handles GET /v1/inbox/labels.
func (h *InboxHandler) ListLabels(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)

	labels, err := h.inboxSvc.ListLabels(ctx, orgID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to list labels")
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, labels)
}

// CreateLabelRequest represents the request body for POST /v1/inbox/labels.
type CreateLabelRequest struct {
	Name  string `json:"name" validate:"required"`
	Color string `json:"color,omitempty"`
}

// CreateLabel handles POST /v1/inbox/labels.
func (h *InboxHandler) CreateLabel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)

	var req CreateLabelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if err := validate.Struct(req); err != nil {
		response.BadRequest(w, formatValidationError(err))
		return
	}

	label, err := h.inboxSvc.CreateLabel(ctx, orgID, req.Name, req.Color)
	if err != nil {
		if errors.Is(err, inbox.ErrLabelExists) {
			response.Conflict(w, "label with this name already exists")
			return
		}
		logger.Error().Err(err).Str("name", req.Name).Msg("failed to create label")
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusCreated, label)
}

// DeleteLabel handles DELETE /v1/inbox/labels/{id}.
func (h *InboxHandler) DeleteLabel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)

	labelID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid label id")
		return
	}

	if err := h.inboxSvc.DeleteLabel(ctx, orgID, labelID); err != nil {
		logger.Error().Err(err).Str("label_id", labelID.String()).Msg("failed to delete label")
		response.InternalError(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AddLabelRequest represents the request body for POST /v1/inbox/messages/{id}/labels.
type AddLabelRequest struct {
	LabelID uuid.UUID `json:"label_id" validate:"required"`
}

// AddLabel handles POST /v1/inbox/messages/{id}/labels.
func (h *InboxHandler) AddLabel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)

	messageID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid message id")
		return
	}

	var req AddLabelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if req.LabelID == uuid.Nil {
		response.BadRequest(w, "label_id is required")
		return
	}

	// Verify the message exists and belongs to this org.
	if _, err := h.inboxSvc.GetMessage(ctx, orgID, messageID); err != nil {
		if errors.Is(err, inbox.ErrMessageNotFound) {
			response.NotFound(w, "message not found")
			return
		}
		logger.Error().Err(err).Str("message_id", messageID.String()).Msg("failed to verify message")
		response.InternalError(w)
		return
	}

	if err := h.inboxSvc.AddMessageLabel(ctx, messageID, req.LabelID); err != nil {
		logger.Error().Err(err).
			Str("message_id", messageID.String()).
			Str("label_id", req.LabelID.String()).
			Msg("failed to add label to message")
		response.InternalError(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RemoveLabel handles DELETE /v1/inbox/messages/{id}/labels/{labelId}.
func (h *InboxHandler) RemoveLabel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)

	messageID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid message id")
		return
	}

	labelID, err := uuid.Parse(chi.URLParam(r, "labelId"))
	if err != nil {
		response.BadRequest(w, "invalid label id")
		return
	}

	if err := h.inboxSvc.RemoveMessageLabel(ctx, messageID, labelID); err != nil {
		logger.Error().Err(err).
			Str("message_id", messageID.String()).
			Str("label_id", labelID.String()).
			Msg("failed to remove label from message")
		response.InternalError(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
