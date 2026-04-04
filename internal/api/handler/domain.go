package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/hellomail/hellomail/internal/api/response"
	"github.com/hellomail/hellomail/internal/auth"
	sqlcdb "github.com/hellomail/hellomail/internal/db/sqlcdb"
	"github.com/hellomail/hellomail/internal/domain"
	"github.com/hellomail/hellomail/internal/observability"
)

// DomainHandler handles domain management HTTP requests.
type DomainHandler struct {
	domainSvc *domain.Service
}

// NewDomainHandler creates a new DomainHandler with the given domain service.
func NewDomainHandler(domainSvc *domain.Service) *DomainHandler {
	return &DomainHandler{
		domainSvc: domainSvc,
	}
}

// createDomainRequest represents the request body for creating a domain.
type createDomainRequest struct {
	Name          string `json:"name"`
	EnableInbound bool   `json:"enable_inbound"`
	SkipDMARC     bool   `json:"skip_dmarc"`
	MergedSPF     string `json:"merged_spf,omitempty"`
}

// updateDomainRequest represents the request body for updating domain settings.
type updateDomainRequest struct {
	OpenTracking  *bool `json:"open_tracking"`
	ClickTracking *bool `json:"click_tracking"`
}

// autoDNSRequest represents the request body for auto-creating DNS records via Cloudflare.
type autoDNSRequest struct {
	ZoneID   string `json:"zone_id"`
	APIToken string `json:"api_token"`
}

// domainResponse represents a domain in API responses.
type domainResponse struct {
	ID            uuid.UUID  `json:"id"`
	OrgID         uuid.UUID  `json:"org_id"`
	Name          string     `json:"name"`
	Status        string     `json:"status"`
	Region        string     `json:"region"`
	DkimSelector  string     `json:"dkim_selector"`
	OpenTracking  bool       `json:"open_tracking"`
	ClickTracking bool       `json:"click_tracking"`
	VerifiedAt    *time.Time `json:"verified_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// dnsRecordResponse represents a DNS record in API responses.
type dnsRecordResponse struct {
	ID         uuid.UUID  `json:"id"`
	DomainID   uuid.UUID  `json:"domain_id"`
	RecordType string     `json:"record_type"`
	Host       string     `json:"host"`
	Value      string     `json:"value"`
	Purpose    string     `json:"purpose"`
	Status     string     `json:"status"`
	VerifiedAt *time.Time `json:"verified_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

// createDomainResponse wraps a domain with its DNS records for the create response.
type createDomainResponse struct {
	Domain     domainResponse      `json:"domain"`
	DNSRecords []dnsRecordResponse `json:"dns_records"`
}

// Analyze checks existing DNS records for a domain before setup.
// POST /v1/domains/analyze
func (h *DomainHandler) Analyze(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)

	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	if req.Name == "" {
		response.BadRequest(w, "name is required")
		return
	}

	analysis, err := domain.AnalyzeDomain(ctx, strings.ToLower(strings.TrimSpace(req.Name)))
	if err != nil {
		logger.Error().Err(err).Str("domain", req.Name).Msg("failed to analyze domain")
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, analysis)
}

// Create registers a new sending domain.
// POST /v1/domains
func (h *DomainHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)

	var req createDomainRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if req.Name == "" {
		response.BadRequest(w, "name is required")
		return
	}

	dom, records, err := h.domainSvc.CreateDomain(ctx, orgID, req.Name, &domain.CreateDomainOptions{
		EnableInbound: req.EnableInbound,
		SkipDMARC:     req.SkipDMARC,
		MergedSPF:     req.MergedSPF,
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidDomain):
			response.BadRequest(w, "invalid domain name format")
		case errors.Is(err, domain.ErrDomainConflict):
			response.Conflict(w, "domain already exists for this organization")
		default:
			logger.Error().Err(err).Str("domain", req.Name).Msg("failed to create domain")
			response.InternalError(w)
		}
		return
	}

	resp := createDomainResponse{
		Domain:     toDomainResponse(dom),
		DNSRecords: toDNSRecordResponses(records),
	}

	response.JSON(w, http.StatusCreated, resp)
}

// List returns all domains for the authenticated organization.
// GET /v1/domains
func (h *DomainHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)

	domains, err := h.domainSvc.ListDomains(ctx, orgID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to list domains")
		response.InternalError(w)
		return
	}

	items := make([]domainResponse, len(domains))
	for i, d := range domains {
		items[i] = toDomainResponse(&d)
	}

	response.JSON(w, http.StatusOK, items)
}

// Get returns a single domain with its DNS records.
// GET /v1/domains/{id}
func (h *DomainHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)

	domainID, err := parseUUIDParam(r, "id")
	if err != nil {
		response.BadRequest(w, "invalid domain id")
		return
	}

	dom, err := h.domainSvc.GetDomain(ctx, orgID, domainID)
	if err != nil {
		if errors.Is(err, domain.ErrDomainNotFound) {
			response.NotFound(w, "domain not found")
			return
		}
		logger.Error().Err(err).Str("domain_id", domainID.String()).Msg("failed to get domain")
		response.InternalError(w)
		return
	}

	// Fetch DNS records for this domain
	queries := sqlcdb.New(h.domainSvc.DB())
	records, err := queries.ListDNSRecordsByDomain(ctx, dom.ID)
	if err != nil {
		logger.Error().Err(err).Str("domain_id", domainID.String()).Msg("failed to list DNS records")
		response.InternalError(w)
		return
	}

	resp := createDomainResponse{
		Domain:     toDomainResponse(dom),
		DNSRecords: toDNSRecordResponses(records),
	}
	response.JSON(w, http.StatusOK, resp)
}

// Verify triggers DNS verification for a domain.
// POST /v1/domains/{id}/verify
func (h *DomainHandler) Verify(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)

	domainID, err := parseUUIDParam(r, "id")
	if err != nil {
		response.BadRequest(w, "invalid domain id")
		return
	}

	records, err := h.domainSvc.VerifyDomain(ctx, orgID, domainID)
	if err != nil {
		if errors.Is(err, domain.ErrDomainNotFound) {
			response.NotFound(w, "domain not found")
			return
		}
		logger.Error().Err(err).Str("domain_id", domainID.String()).Msg("failed to verify domain")
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, toDNSRecordResponses(records))
}

// Update modifies domain settings (tracking options).
// PATCH /v1/domains/{id}
func (h *DomainHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)

	domainID, err := parseUUIDParam(r, "id")
	if err != nil {
		response.BadRequest(w, "invalid domain id")
		return
	}

	var req updateDomainRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	// Get current domain to use existing values as defaults
	current, err := h.domainSvc.GetDomain(ctx, orgID, domainID)
	if err != nil {
		if errors.Is(err, domain.ErrDomainNotFound) {
			response.NotFound(w, "domain not found")
			return
		}
		logger.Error().Err(err).Str("domain_id", domainID.String()).Msg("failed to get domain for update")
		response.InternalError(w)
		return
	}

	openTracking := current.OpenTracking
	clickTracking := current.ClickTracking
	if req.OpenTracking != nil {
		openTracking = *req.OpenTracking
	}
	if req.ClickTracking != nil {
		clickTracking = *req.ClickTracking
	}

	dom, err := h.domainSvc.UpdateSettings(ctx, orgID, domainID, openTracking, clickTracking)
	if err != nil {
		if errors.Is(err, domain.ErrDomainNotFound) {
			response.NotFound(w, "domain not found")
			return
		}
		logger.Error().Err(err).Str("domain_id", domainID.String()).Msg("failed to update domain")
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, toDomainResponse(dom))
}

// Delete removes a domain and its associated DNS records.
// DELETE /v1/domains/{id}
func (h *DomainHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)

	domainID, err := parseUUIDParam(r, "id")
	if err != nil {
		response.BadRequest(w, "invalid domain id")
		return
	}

	if err := h.domainSvc.DeleteDomain(ctx, orgID, domainID); err != nil {
		if errors.Is(err, domain.ErrDomainNotFound) {
			response.NotFound(w, "domain not found")
			return
		}
		logger.Error().Err(err).Str("domain_id", domainID.String()).Msg("failed to delete domain")
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "domain deleted"})
}

// GetConnectURL checks if the domain's DNS provider supports Domain Connect
// and returns a redirect URL for automatic DNS configuration.
// GET /v1/domains/{id}/connect-url
func (h *DomainHandler) GetConnectURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := observability.LoggerFromContext(ctx)
	orgID := auth.OrgIDFromContext(ctx)

	domainID, err := parseUUIDParam(r, "id")
	if err != nil {
		response.BadRequest(w, "invalid domain id")
		return
	}

	dom, err := h.domainSvc.GetDomain(ctx, orgID, domainID)
	if err != nil {
		if errors.Is(err, domain.ErrDomainNotFound) {
			response.NotFound(w, "domain not found")
			return
		}
		logger.Error().Err(err).Msg("failed to get domain")
		response.InternalError(w)
		return
	}

	// Check if Domain Connect is supported
	provider, err := domain.DiscoverProvider(ctx, dom.Name)
	if err != nil {
		response.JSON(w, http.StatusOK, domain.DomainConnectResult{
			Supported: false,
			Message:   "Your DNS provider does not support automatic configuration. Please add DNS records manually.",
		})
		return
	}

	// Get the DKIM public key from the domain's private key
	publicKey := ""
	if dom.DkimPrivateKey.Valid {
		pk, err := domain.PublicKeyFromPEM(dom.DkimPrivateKey.String)
		if err == nil {
			publicKey = pk
		}
	}

	// Generate CSRF state
	state, _ := domain.GenerateState()

	// Build the callback URL
	callbackURL := r.URL.Query().Get("callback_url")
	if callbackURL == "" {
		callbackURL = fmt.Sprintf("https://app.mailngine.com/domains/%s?connected=true", domainID.String())
	}

	// Build the redirect URL
	signingKey := "" // TODO: Load from config when registered as DC provider
	redirectURL, err := domain.GenerateRedirectURL(provider, dom.Name, dom.DkimSelector, publicKey, callbackURL, state, signingKey)
	if err != nil {
		logger.Error().Err(err).Msg("failed to generate Domain Connect URL")
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, domain.DomainConnectResult{
		Supported:   true,
		Provider:    provider.ProviderName,
		RedirectURL: redirectURL,
	})
}

// parseUUIDParam extracts and parses a UUID from a chi URL parameter.
func parseUUIDParam(r *http.Request, param string) (uuid.UUID, error) {
	return uuid.Parse(chi.URLParam(r, param))
}

// toDomainResponse converts a database domain model to an API response.
func toDomainResponse(d *sqlcdb.Domain) domainResponse {
	resp := domainResponse{
		ID:            d.ID,
		OrgID:         d.OrgID,
		Name:          d.Name,
		Status:        d.Status,
		Region:        d.Region,
		DkimSelector:  d.DkimSelector,
		OpenTracking:  d.OpenTracking,
		ClickTracking: d.ClickTracking,
		CreatedAt:     d.CreatedAt,
		UpdatedAt:     d.UpdatedAt,
	}
	if d.VerifiedAt.Valid {
		resp.VerifiedAt = &d.VerifiedAt.Time
	}
	return resp
}

// toDNSRecordResponses converts a slice of database DNS records to API responses.
func toDNSRecordResponses(records []sqlcdb.DnsRecord) []dnsRecordResponse {
	items := make([]dnsRecordResponse, len(records))
	for i, r := range records {
		items[i] = toDNSRecordResponse(r)
	}
	return items
}

// toDNSRecordResponse converts a database DNS record model to an API response.
func toDNSRecordResponse(r sqlcdb.DnsRecord) dnsRecordResponse {
	resp := dnsRecordResponse{
		ID:         r.ID,
		DomainID:   r.DomainID,
		RecordType: r.RecordType,
		Host:       r.Host,
		Value:      r.Value,
		Purpose:    r.Purpose,
		Status:     r.Status,
		CreatedAt:  r.CreatedAt,
	}
	if r.VerifiedAt.Valid {
		resp.VerifiedAt = &r.VerifiedAt.Time
	}
	return resp
}
