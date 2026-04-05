package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/mailngine/mailngine/internal/api/response"
	"github.com/mailngine/mailngine/internal/auth"
	sqlcdb "github.com/mailngine/mailngine/internal/db/sqlcdb"
	"github.com/mailngine/mailngine/internal/observability"
	"github.com/mailngine/mailngine/internal/team"
)

// TeamHandler handles team and organization management HTTP requests.
type TeamHandler struct {
	teamSvc *team.Service
}

// NewTeamHandler creates a new TeamHandler with the given team service.
func NewTeamHandler(svc *team.Service) *TeamHandler {
	return &TeamHandler{
		teamSvc: svc,
	}
}

// updateOrgRequest represents the request body for updating an organization.
type updateOrgRequest struct {
	Name string `json:"name"`
}

// inviteMemberRequest represents the request body for inviting a member.
type inviteMemberRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

// updateRoleRequest represents the request body for updating a member's role.
type updateRoleRequest struct {
	Role string `json:"role"`
}

// orgResponse represents an organization in API responses.
type orgResponse struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Slug           string    `json:"slug"`
	Plan           string    `json:"plan"`
	MonthlyLimit   int32     `json:"monthly_limit"`
	OverageEnabled bool      `json:"overage_enabled"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// memberResponse represents an organization member in API responses.
type memberResponse struct {
	ID        uuid.UUID  `json:"id"`
	OrgID     uuid.UUID  `json:"org_id"`
	UserID    uuid.UUID  `json:"user_id"`
	Role      string     `json:"role"`
	InvitedBy *uuid.UUID `json:"invited_by,omitempty"`
	JoinedAt  time.Time  `json:"joined_at"`
	Email     string     `json:"email"`
	Name      string     `json:"name"`
	AvatarURL *string    `json:"avatar_url,omitempty"`
}

// GetOrg returns the current organization.
// GET /v1/org
func (h *TeamHandler) GetOrg(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)

	org, err := h.teamSvc.GetOrg(ctx, orgID)
	if err != nil {
		if errors.Is(err, team.ErrOrgNotFound) {
			response.NotFound(w, "organization not found")
			return
		}
		logger.Error().Err(err).Msg("failed to get organization")
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, toOrgResponse(org))
}

// UpdateOrg updates the current organization's name.
// PATCH /v1/org
func (h *TeamHandler) UpdateOrg(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)
	role := auth.RoleFromContext(ctx)

	if role != "admin" && role != "owner" {
		response.Forbidden(w, "only admins and owners can update the organization")
		return
	}

	var req updateOrgRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	org, err := h.teamSvc.UpdateOrg(ctx, orgID, req.Name)
	if err != nil {
		switch {
		case errors.Is(err, team.ErrNameRequired):
			response.BadRequest(w, err.Error())
		case errors.Is(err, team.ErrOrgNotFound):
			response.NotFound(w, "organization not found")
		default:
			logger.Error().Err(err).Msg("failed to update organization")
			response.InternalError(w)
		}
		return
	}

	response.JSON(w, http.StatusOK, toOrgResponse(org))
}

// ListMembers returns all members of the current organization.
// GET /v1/org/members
func (h *TeamHandler) ListMembers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)

	members, err := h.teamSvc.ListMembers(ctx, orgID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to list members")
		response.InternalError(w)
		return
	}

	items := make([]memberResponse, len(members))
	for i, m := range members {
		items[i] = toMemberResponse(m)
	}

	response.JSON(w, http.StatusOK, items)
}

// InviteMember adds a registered user to the current organization.
// POST /v1/org/members/invite
func (h *TeamHandler) InviteMember(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)
	userID := auth.UserIDFromContext(ctx)
	role := auth.RoleFromContext(ctx)

	if role != "admin" && role != "owner" {
		response.Forbidden(w, "only admins and owners can invite members")
		return
	}

	var req inviteMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if req.Email == "" {
		response.BadRequest(w, "email is required")
		return
	}
	if req.Role == "" {
		response.BadRequest(w, "role is required")
		return
	}

	err := h.teamSvc.InviteMember(ctx, orgID, userID, req.Email, req.Role)
	if err != nil {
		switch {
		case errors.Is(err, team.ErrUserNotRegistered):
			response.BadRequest(w, err.Error())
		case errors.Is(err, team.ErrAlreadyMember):
			response.Conflict(w, err.Error())
		case errors.Is(err, team.ErrInvalidRole):
			response.BadRequest(w, err.Error())
		default:
			logger.Error().Err(err).Str("email", req.Email).Msg("failed to invite member")
			response.InternalError(w)
		}
		return
	}

	response.JSON(w, http.StatusCreated, map[string]string{"message": "member invited"})
}

// UpdateRole changes the role of an organization member.
// PATCH /v1/org/members/{id}
func (h *TeamHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)
	role := auth.RoleFromContext(ctx)

	if role != "admin" && role != "owner" {
		response.Forbidden(w, "only admins and owners can update member roles")
		return
	}

	memberUserID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid member id")
		return
	}

	var req updateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if req.Role == "" {
		response.BadRequest(w, "role is required")
		return
	}

	err = h.teamSvc.UpdateRole(ctx, orgID, memberUserID, req.Role)
	if err != nil {
		switch {
		case errors.Is(err, team.ErrMemberNotFound):
			response.NotFound(w, "member not found")
		case errors.Is(err, team.ErrInvalidRole):
			response.BadRequest(w, err.Error())
		case errors.Is(err, team.ErrLastOwner):
			response.BadRequest(w, err.Error())
		default:
			logger.Error().Err(err).Str("member_user_id", memberUserID.String()).Msg("failed to update member role")
			response.InternalError(w)
		}
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "role updated"})
}

// RemoveMember removes a member from the current organization.
// DELETE /v1/org/members/{id}
func (h *TeamHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)
	callerUserID := auth.UserIDFromContext(ctx)
	role := auth.RoleFromContext(ctx)

	if role != "admin" && role != "owner" {
		response.Forbidden(w, "only admins and owners can remove members")
		return
	}

	memberUserID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid member id")
		return
	}

	// Prevent self-removal via this endpoint.
	if memberUserID == callerUserID {
		response.BadRequest(w, team.ErrCannotRemoveSelf.Error())
		return
	}

	err = h.teamSvc.RemoveMember(ctx, orgID, memberUserID)
	if err != nil {
		switch {
		case errors.Is(err, team.ErrMemberNotFound):
			response.NotFound(w, "member not found")
		case errors.Is(err, team.ErrLastOwner):
			response.BadRequest(w, err.Error())
		default:
			logger.Error().Err(err).Str("member_user_id", memberUserID.String()).Msg("failed to remove member")
			response.InternalError(w)
		}
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "member removed"})
}

// toOrgResponse converts a database organization to an API response.
func toOrgResponse(org *sqlcdb.Organization) orgResponse {
	return orgResponse{
		ID:             org.ID,
		Name:           org.Name,
		Slug:           org.Slug,
		Plan:           org.Plan,
		MonthlyLimit:   org.MonthlyLimit,
		OverageEnabled: org.OverageEnabled,
		CreatedAt:      org.CreatedAt,
		UpdatedAt:      org.UpdatedAt,
	}
}

// toMemberResponse converts a database org member row to an API response.
func toMemberResponse(m sqlcdb.ListOrgMembersRow) memberResponse {
	resp := memberResponse{
		ID:       m.ID,
		OrgID:    m.OrgID,
		UserID:   m.UserID,
		Role:     m.Role,
		JoinedAt: m.JoinedAt,
		Email:    m.Email,
		Name:     m.Name,
	}
	if m.InvitedBy.Valid {
		id := uuid.UUID(m.InvitedBy.Bytes)
		resp.InvitedBy = &id
	}
	if m.AvatarUrl.Valid {
		resp.AvatarURL = &m.AvatarUrl.String
	}
	return resp
}
