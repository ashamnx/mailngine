package domain

import (
	"context"
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	sqlcdb "github.com/mailngine/mailngine/internal/db/sqlcdb"
)

const (
	dkimSelector = "mn1"

	statusPending  = "pending"
	statusVerified = "verified"
	statusFailed   = "failed"

	purposeSPF        = "spf"
	purposeDKIM       = "dkim"
	purposeDMARC      = "dmarc"
	purposeMX         = "mx"
	purposeReturnPath = "return_path"
)

// domainRegexp validates a hostname: labels separated by dots, each label 1-63 chars,
// total length up to 253 chars. No protocol, no trailing dot.
var domainRegexp = regexp.MustCompile(`^([a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`)

// ErrDomainConflict indicates the domain already exists for this organization.
var ErrDomainConflict = errors.New("domain already exists for this organization")

// ErrDomainNotFound indicates the requested domain was not found.
var ErrDomainNotFound = errors.New("domain not found")

// ErrInvalidDomain indicates the domain name format is invalid.
var ErrInvalidDomain = errors.New("invalid domain name")

// Service provides domain management operations.
type Service struct {
	db      *pgxpool.Pool
	queries *sqlcdb.Queries
	logger  zerolog.Logger
}

// NewService creates a new domain Service.
func NewService(db *pgxpool.Pool, logger zerolog.Logger) *Service {
	return &Service{
		db:      db,
		queries: sqlcdb.New(db),
		logger:  logger,
	}
}

// DB returns the underlying database pool for use in handlers that need
// to perform additional queries (e.g., listing DNS records for auto-dns).
func (s *Service) DB() *pgxpool.Pool {
	return s.db
}

// CreateDomainOptions configures domain creation behavior.
type CreateDomainOptions struct {
	// EnableInbound generates MX records for receiving email on this domain.
	// Should be false if the domain is used with Google Workspace, Office 365, etc.
	EnableInbound bool
	// SkipDMARC skips DMARC record generation (when domain already has DMARC).
	SkipDMARC bool
	// MergedSPF is the pre-computed merged SPF record (from analysis).
	// If empty, a standalone SPF include is generated.
	MergedSPF string
}

// CreateDomain registers a new sending domain for the given organization.
// It generates a DKIM keypair and creates the required DNS records.
// By default, only sending-related records are created (DKIM, SPF include, Return-Path).
// MX records are only created if EnableInbound is true.
func (s *Service) CreateDomain(ctx context.Context, orgID uuid.UUID, name string, opts *CreateDomainOptions) (*sqlcdb.Domain, []sqlcdb.DnsRecord, error) {
	name = strings.ToLower(strings.TrimSpace(name))

	if opts == nil {
		opts = &CreateDomainOptions{}
	}

	if !isValidDomain(name) {
		return nil, nil, ErrInvalidDomain
	}

	// Check for existing domain within this org
	_, err := s.queries.GetDomainByName(ctx, sqlcdb.GetDomainByNameParams{
		Name:  name,
		OrgID: orgID,
	})
	if err == nil {
		return nil, nil, ErrDomainConflict
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, nil, fmt.Errorf("checking existing domain: %w", err)
	}

	// Generate DKIM keypair
	privatePEM, publicKeyB64, err := GenerateDKIMKeyPair()
	if err != nil {
		return nil, nil, fmt.Errorf("generating DKIM keypair: %w", err)
	}

	// Use a transaction to ensure atomicity
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := s.queries.WithTx(tx)

	// Insert domain
	domain, err := qtx.CreateDomain(ctx, sqlcdb.CreateDomainParams{
		OrgID:          orgID,
		Name:           name,
		DkimPrivateKey: pgtype.Text{String: privatePEM, Valid: true},
		DkimSelector:   dkimSelector,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("inserting domain: %w", err)
	}

	// Build DNS records based on options.
	// SPF: We generate an "include" directive — users with existing SPF records
	// (e.g., Google Workspace, Office 365) must merge this into their existing record.
	// DKIM: Uses a unique selector (mn1) so it safely coexists with other providers.
	// DMARC: Only suggested if user doesn't already have one.
	// MX: Only generated when inbound email is enabled (to avoid breaking existing mailboxes).
	// Return-Path: CNAME for bounce handling, always safe to add.
	// SPF: Use pre-computed merged SPF if provided (from domain analysis),
	// otherwise generate a standalone SPF include.
	spfValue := "v=spf1 include:spf.mailngine.com ~all"
	if opts.MergedSPF != "" {
		spfValue = opts.MergedSPF
	}

	recordDefs := []sqlcdb.CreateDNSRecordParams{
		{
			DomainID:   domain.ID,
			RecordType: "TXT",
			Host:       name,
			Value:      spfValue,
			Purpose:    purposeSPF,
		},
		{
			DomainID:   domain.ID,
			RecordType: "TXT",
			Host:       dkimSelector + "._domainkey." + name,
			Value:      "v=DKIM1; k=rsa; p=" + publicKeyB64,
			Purpose:    purposeDKIM,
		},
		{
			DomainID:   domain.ID,
			RecordType: "CNAME",
			Host:       "bounce." + name,
			Value:      "bounces.mailngine.com",
			Purpose:    purposeReturnPath,
		},
	}

	// Only add DMARC if the domain doesn't already have one
	if !opts.SkipDMARC {
		recordDefs = append(recordDefs, sqlcdb.CreateDNSRecordParams{
			DomainID:   domain.ID,
			RecordType: "TXT",
			Host:       "_dmarc." + name,
			Value:      "v=DMARC1; p=none; rua=mailto:dmarc@mailngine.com",
			Purpose:    purposeDMARC,
		})
	}

	// Only add MX record if inbound email is enabled
	if opts.EnableInbound {
		recordDefs = append(recordDefs, sqlcdb.CreateDNSRecordParams{
			DomainID:   domain.ID,
			RecordType: "MX",
			Host:       name,
			Value:      "10 mx.mailngine.com",
			Purpose:    purposeMX,
		})
	}

	records := make([]sqlcdb.DnsRecord, 0, len(recordDefs))
	for _, def := range recordDefs {
		rec, err := qtx.CreateDNSRecord(ctx, def)
		if err != nil {
			return nil, nil, fmt.Errorf("inserting DNS record (%s): %w", def.Purpose, err)
		}
		records = append(records, rec)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, nil, fmt.Errorf("committing transaction: %w", err)
	}

	return &domain, records, nil
}

// GetDomain retrieves a domain by ID scoped to the organization.
func (s *Service) GetDomain(ctx context.Context, orgID, domainID uuid.UUID) (*sqlcdb.Domain, error) {
	domain, err := s.queries.GetDomain(ctx, sqlcdb.GetDomainParams{
		ID:    domainID,
		OrgID: orgID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrDomainNotFound
		}
		return nil, fmt.Errorf("getting domain: %w", err)
	}
	return &domain, nil
}

// ListDomains retrieves all domains for the given organization.
func (s *Service) ListDomains(ctx context.Context, orgID uuid.UUID) ([]sqlcdb.Domain, error) {
	domains, err := s.queries.ListDomainsByOrg(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("listing domains: %w", err)
	}
	return domains, nil
}

// UpdateSettings updates the tracking settings for a domain.
func (s *Service) UpdateSettings(ctx context.Context, orgID, domainID uuid.UUID, openTracking, clickTracking bool) (*sqlcdb.Domain, error) {
	// Verify domain exists and belongs to org
	_, err := s.GetDomain(ctx, orgID, domainID)
	if err != nil {
		return nil, err
	}

	domain, err := s.queries.UpdateDomainSettings(ctx, sqlcdb.UpdateDomainSettingsParams{
		ID:            domainID,
		OrgID:         orgID,
		OpenTracking:  openTracking,
		ClickTracking: clickTracking,
	})
	if err != nil {
		return nil, fmt.Errorf("updating domain settings: %w", err)
	}
	return &domain, nil
}

// DeleteDomain removes a domain and its associated DNS records.
func (s *Service) DeleteDomain(ctx context.Context, orgID, domainID uuid.UUID) error {
	// Verify domain exists and belongs to org
	_, err := s.GetDomain(ctx, orgID, domainID)
	if err != nil {
		return err
	}

	// DNS records are deleted via ON DELETE CASCADE, but we explicitly delete
	// them to be clear about intent.
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := s.queries.WithTx(tx)

	if err := qtx.DeleteDNSRecordsByDomain(ctx, domainID); err != nil {
		return fmt.Errorf("deleting DNS records: %w", err)
	}

	if err := qtx.DeleteDomain(ctx, sqlcdb.DeleteDomainParams{
		ID:    domainID,
		OrgID: orgID,
	}); err != nil {
		return fmt.Errorf("deleting domain: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}

// VerifyDomain performs DNS lookups for each record associated with the domain
// and updates their verification status. If all records pass, the domain status
// is set to "verified".
func (s *Service) VerifyDomain(ctx context.Context, orgID, domainID uuid.UUID) ([]sqlcdb.DnsRecord, error) {
	// Verify domain exists and belongs to org
	_, err := s.GetDomain(ctx, orgID, domainID)
	if err != nil {
		return nil, err
	}

	records, err := s.queries.ListDNSRecordsByDomain(ctx, domainID)
	if err != nil {
		return nil, fmt.Errorf("listing DNS records: %w", err)
	}

	allVerified := true
	updatedRecords := make([]sqlcdb.DnsRecord, 0, len(records))

	for _, rec := range records {
		verified := verifyDNSRecord(ctx, rec)

		newStatus := statusFailed
		if verified {
			newStatus = statusVerified
		} else {
			allVerified = false
		}

		if err := s.queries.UpdateDNSRecordStatus(ctx, sqlcdb.UpdateDNSRecordStatusParams{
			ID:       rec.ID,
			Status:   newStatus,
			DomainID: rec.DomainID,
		}); err != nil {
			s.logger.Error().Err(err).
				Str("record_id", rec.ID.String()).
				Msg("failed to update DNS record status")
			allVerified = false
		}

		rec.Status = newStatus
		updatedRecords = append(updatedRecords, rec)
	}

	// Update domain status based on verification results
	domainStatus := statusFailed
	if allVerified {
		domainStatus = statusVerified
	}

	if err := s.queries.UpdateDomainStatus(ctx, sqlcdb.UpdateDomainStatusParams{
		ID:     domainID,
		OrgID:  orgID,
		Status: domainStatus,
	}); err != nil {
		return nil, fmt.Errorf("updating domain status: %w", err)
	}

	return updatedRecords, nil
}

// verifyDNSRecord performs a DNS lookup for the given record and checks
// whether the expected value is present.
func verifyDNSRecord(ctx context.Context, rec sqlcdb.DnsRecord) bool {
	resolver := net.DefaultResolver

	switch rec.RecordType {
	case "TXT":
		txts, err := resolver.LookupTXT(ctx, rec.Host)
		if err != nil {
			return false
		}
		for _, txt := range txts {
			if strings.Contains(txt, rec.Value) {
				return true
			}
			// For long TXT records, DNS may split them; check if
			// the concatenated value contains our expected value.
		}
		// For DKIM records, the value may be split across multiple strings.
		// Concatenate all TXT record strings and check.
		concatenated := strings.Join(txts, "")
		return strings.Contains(concatenated, rec.Value)

	case "MX":
		mxs, err := resolver.LookupMX(ctx, rec.Host)
		if err != nil {
			return false
		}
		// Expected value format: "10 mx.mailngine.com"
		parts := strings.Fields(rec.Value)
		if len(parts) < 2 {
			return false
		}
		expectedHost := strings.TrimSuffix(parts[1], ".") + "."
		for _, mx := range mxs {
			if strings.EqualFold(mx.Host, expectedHost) || strings.EqualFold(strings.TrimSuffix(mx.Host, "."), strings.TrimSuffix(parts[1], ".")) {
				return true
			}
		}
		return false

	case "CNAME":
		cname, err := resolver.LookupCNAME(ctx, rec.Host)
		if err != nil {
			return false
		}
		expected := strings.TrimSuffix(rec.Value, ".") + "."
		return strings.EqualFold(cname, expected) || strings.EqualFold(strings.TrimSuffix(cname, "."), strings.TrimSuffix(rec.Value, "."))

	default:
		return false
	}
}

// isValidDomain checks whether the given string is a valid domain name.
func isValidDomain(name string) bool {
	if len(name) == 0 || len(name) > 253 {
		return false
	}
	// Reject protocols
	if strings.Contains(name, "://") {
		return false
	}
	// Reject trailing dots
	if strings.HasSuffix(name, ".") {
		return false
	}
	return domainRegexp.MatchString(name)
}
