package handler

import (
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hellomail/hellomail/internal/api/response"
	"github.com/hellomail/hellomail/internal/auth"
	"github.com/hellomail/hellomail/internal/config"
	sqlcdb "github.com/hellomail/hellomail/internal/db/sqlcdb"
	"github.com/hellomail/hellomail/internal/observability"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

const sessionPrefix = "session:"

// slugRegexp matches characters that are not lowercase alphanumeric or hyphens.
var slugRegexp = regexp.MustCompile(`[^a-z0-9-]`)

// AuthHandler handles authentication-related HTTP requests.
type AuthHandler struct {
	cfg     *config.Config
	db      *pgxpool.Pool
	cache   *redis.Client
	oauth   *auth.GoogleOAuth
	jwtMgr  *auth.JWTManager
	queries *sqlcdb.Queries
}

// NewAuthHandler creates a new AuthHandler with the given dependencies.
func NewAuthHandler(cfg *config.Config, db *pgxpool.Pool, cache *redis.Client, oauth *auth.GoogleOAuth, jwtMgr *auth.JWTManager) *AuthHandler {
	return &AuthHandler{
		cfg:     cfg,
		db:      db,
		cache:   cache,
		oauth:   oauth,
		jwtMgr:  jwtMgr,
		queries: sqlcdb.New(db),
	}
}

// GoogleRedirect initiates the Google OAuth2 flow by redirecting the user to Google's consent page.
// GET /v1/auth/google
func (h *AuthHandler) GoogleRedirect(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)

	url, err := h.oauth.AuthCodeURL(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to generate oauth url")
		response.InternalError(w)
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// GoogleCallback handles the Google OAuth2 callback, creates or updates the user,
// provisions a default organization if needed, and issues a JWT token.
// GET /v1/auth/google/callback
func (h *AuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)

	// Validate state parameter for CSRF protection
	state := r.URL.Query().Get("state")
	if state == "" {
		response.BadRequest(w, "missing state parameter")
		return
	}
	if err := h.oauth.ValidateState(ctx, state); err != nil {
		logger.Warn().Err(err).Msg("invalid oauth state")
		response.BadRequest(w, "invalid or expired state parameter")
		return
	}

	// Exchange authorization code for user info
	code := r.URL.Query().Get("code")
	if code == "" {
		response.BadRequest(w, "missing code parameter")
		return
	}
	info, err := h.oauth.Exchange(ctx, code)
	if err != nil {
		logger.Error().Err(err).Msg("failed to exchange oauth code")
		response.InternalError(w)
		return
	}

	// Upsert user in database
	user, err := h.queries.UpsertUser(ctx, sqlcdb.UpsertUserParams{
		Email:     info.Email,
		Name:      info.Name,
		AvatarUrl: pgtype.Text{String: info.Picture, Valid: info.Picture != ""},
		GoogleID:  info.ID,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to upsert user")
		response.InternalError(w)
		return
	}

	// Check if user belongs to any organization
	orgs, err := h.queries.ListUserOrgs(ctx, user.ID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to list user orgs")
		response.InternalError(w)
		return
	}

	// Create default organization if user has none
	if len(orgs) == 0 {
		slug := deriveSlug(user.Email)
		org, err := h.queries.CreateOrganization(ctx, sqlcdb.CreateOrganizationParams{
			Name:         user.Name + "'s Org",
			Slug:         slug,
			Plan:         "free",
			MonthlyLimit: 100,
		})
		if err != nil {
			logger.Error().Err(err).Msg("failed to create default organization")
			response.InternalError(w)
			return
		}

		_, err = h.queries.AddOrgMember(ctx, sqlcdb.AddOrgMemberParams{
			OrgID:     org.ID,
			UserID:    user.ID,
			Role:      "owner",
			InvitedBy: pgtype.UUID{Valid: false},
		})
		if err != nil {
			logger.Error().Err(err).Msg("failed to add user as org owner")
			response.InternalError(w)
			return
		}

		// Refresh org list
		orgs, err = h.queries.ListUserOrgs(ctx, user.ID)
		if err != nil {
			logger.Error().Err(err).Msg("failed to list user orgs after creation")
			response.InternalError(w)
			return
		}
	}

	// Generate JWT for the user's first organization
	token, err := h.jwtMgr.Generate(user.ID, orgs[0].ID, orgs[0].Role)
	if err != nil {
		logger.Error().Err(err).Msg("failed to generate jwt")
		response.InternalError(w)
		return
	}

	// Store session in Valkey
	sessionKey := sessionPrefix + user.ID.String()
	expirySeconds := int(h.jwtMgr.Expiry().Seconds())
	if err := h.cache.Set(ctx, sessionKey, "1", time.Duration(expirySeconds)*time.Second).Err(); err != nil {
		logger.Error().Err(err).Msg("failed to store session")
		response.InternalError(w)
		return
	}

	// Set HttpOnly cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   !h.cfg.IsDevelopment(),
		SameSite: http.SameSiteLaxMode,
		MaxAge:   expirySeconds,
	})

	// Redirect to frontend with token
	redirectURL := h.cfg.FrontendURL + "/auth/callback?token=" + token
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

// Logout invalidates the user's session and clears the auth cookie.
// POST /v1/auth/logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)

	userID := auth.UserIDFromContext(ctx)
	if userID == uuid.Nil {
		response.Unauthorized(w, "not authenticated")
		return
	}

	// Remove session from Valkey
	sessionKey := sessionPrefix + userID.String()
	if err := h.cache.Del(ctx, sessionKey).Err(); err != nil {
		logger.Error().Err(err).Msg("failed to delete session")
		// Continue with logout even if session deletion fails
	}

	// Clear cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   !h.cfg.IsDevelopment(),
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	response.JSON(w, http.StatusOK, map[string]string{"message": "logged out"})
}

// MeResponse represents the response for the /auth/me endpoint.
type MeResponse struct {
	User          UserResponse  `json:"user"`
	Organization  OrgResponse   `json:"organization"`
	Role          string        `json:"role"`
	Organizations []OrgListItem `json:"organizations"`
}

// UserResponse represents a user in API responses.
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	AvatarURL string    `json:"avatar_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// OrgResponse represents an organization in API responses.
type OrgResponse struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Slug           string    `json:"slug"`
	Plan           string    `json:"plan"`
	MonthlyLimit   int32     `json:"monthly_limit"`
	OverageEnabled bool      `json:"overage_enabled"`
}

// OrgListItem represents an organization list item with the user's role.
type OrgListItem struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Slug string    `json:"slug"`
	Role string    `json:"role"`
}

// Me returns the authenticated user's profile, current organization, and all organizations.
// GET /v1/auth/me
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)

	userID := auth.UserIDFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)
	role := auth.RoleFromContext(ctx)

	user, err := h.queries.GetUser(ctx, userID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get user")
		response.InternalError(w)
		return
	}

	org, err := h.queries.GetOrganization(ctx, orgID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get organization")
		response.InternalError(w)
		return
	}

	orgs, err := h.queries.ListUserOrgs(ctx, userID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to list user orgs")
		response.InternalError(w)
		return
	}

	orgList := make([]OrgListItem, len(orgs))
	for i, o := range orgs {
		orgList[i] = OrgListItem{
			ID:   o.ID,
			Name: o.Name,
			Slug: o.Slug,
			Role: o.Role,
		}
	}

	var avatarURL string
	if user.AvatarUrl.Valid {
		avatarURL = user.AvatarUrl.String
	}

	resp := MeResponse{
		User: UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			AvatarURL: avatarURL,
			CreatedAt: user.CreatedAt,
		},
		Organization: OrgResponse{
			ID:             org.ID,
			Name:           org.Name,
			Slug:           org.Slug,
			Plan:           org.Plan,
			MonthlyLimit:   org.MonthlyLimit,
			OverageEnabled: org.OverageEnabled,
		},
		Role:          role,
		Organizations: orgList,
	}

	response.JSON(w, http.StatusOK, resp)
}

// deriveSlug creates a URL-safe slug from an email address.
func deriveSlug(email string) string {
	parts := strings.SplitN(email, "@", 2)
	slug := strings.ToLower(parts[0])
	slug = slugRegexp.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	if slug == "" {
		slug = "org"
	}
	return slug
}
