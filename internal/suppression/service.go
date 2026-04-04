// Package suppression manages the per-organization email suppression list.
// Suppressed addresses are those that should not receive further emails,
// typically due to hard bounces, spam complaints, or manual operator action.
package suppression

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"

	sqlcdb "github.com/hellomail/hellomail/internal/db/sqlcdb"
)

// cacheTTL is the duration a suppression check result is cached in Valkey.
const cacheTTL = 1 * time.Hour

// ErrSuppressionNotFound is returned when a suppression record does not exist.
var ErrSuppressionNotFound = errors.New("suppression not found")

// ErrAlreadySuppressed is returned when an address is already on the suppression list.
var ErrAlreadySuppressed = errors.New("email is already suppressed")

// Service provides operations on the suppression list backed by PostgreSQL
// with a Valkey (Redis-compatible) read-through cache layer.
type Service struct {
	db      *pgxpool.Pool
	queries *sqlcdb.Queries
	cache   *redis.Client
	logger  zerolog.Logger
}

// NewService creates a new suppression Service with the given dependencies.
func NewService(db *pgxpool.Pool, cache *redis.Client, logger zerolog.Logger) *Service {
	return &Service{
		db:      db,
		queries: sqlcdb.New(db),
		cache:   cache,
		logger:  logger.With().Str("component", "suppression_service").Logger(),
	}
}

// cacheKey returns the Valkey key for a suppression lookup.
func cacheKey(orgID uuid.UUID, email string) string {
	return fmt.Sprintf("suppression:%s:%s", orgID.String(), strings.ToLower(email))
}

// Check determines whether the given email address is suppressed for the
// specified organization. It consults the Valkey cache first and falls back
// to the database. Positive results (suppressed) are cached for one hour.
func (s *Service) Check(ctx context.Context, orgID uuid.UUID, email string) (bool, error) {
	key := cacheKey(orgID, email)

	// Check cache first.
	val, err := s.cache.Get(ctx, key).Result()
	if err == nil {
		return val == "1", nil
	}
	if !errors.Is(err, redis.Nil) {
		// Log but do not fail; fall through to the database.
		s.logger.Warn().Err(err).Str("key", key).Msg("valkey get failed, falling back to database")
	}

	// Query the database.
	suppressed, err := s.queries.CheckSuppressed(ctx, sqlcdb.CheckSuppressedParams{
		OrgID: orgID,
		Email: strings.ToLower(email),
	})
	if err != nil {
		return false, fmt.Errorf("check suppression: %w", err)
	}

	// Cache positive results so subsequent checks are fast.
	if suppressed {
		if cacheErr := s.cache.Set(ctx, key, "1", cacheTTL).Err(); cacheErr != nil {
			s.logger.Warn().Err(cacheErr).Str("key", key).Msg("failed to cache suppression result")
		}
	}

	return suppressed, nil
}

// Add inserts a new suppression record and updates the cache. If the email
// is already suppressed for this organization, the operation is idempotent
// and returns nil.
func (s *Service) Add(ctx context.Context, orgID uuid.UUID, email, reason string, metadata map[string]any) error {
	var metadataBytes []byte
	if metadata != nil {
		var err error
		metadataBytes, err = json.Marshal(metadata)
		if err != nil {
			return fmt.Errorf("marshal suppression metadata: %w", err)
		}
	}

	_, err := s.queries.CreateSuppression(ctx, sqlcdb.CreateSuppressionParams{
		OrgID:    orgID,
		Email:    strings.ToLower(email),
		Reason:   reason,
		Metadata: metadataBytes,
	})
	if err != nil {
		// ON CONFLICT DO NOTHING means pgx.ErrNoRows when the row already exists.
		if errors.Is(err, pgx.ErrNoRows) {
			// Already suppressed; this is not an error.
			return nil
		}
		return fmt.Errorf("create suppression: %w", err)
	}

	// Update cache.
	key := cacheKey(orgID, email)
	if cacheErr := s.cache.Set(ctx, key, "1", cacheTTL).Err(); cacheErr != nil {
		s.logger.Warn().Err(cacheErr).Str("key", key).Msg("failed to cache new suppression")
	}

	return nil
}

// Remove deletes a suppression record by ID and invalidates the cache.
func (s *Service) Remove(ctx context.Context, orgID uuid.UUID, suppressionID uuid.UUID) error {
	// Fetch the email address before deletion so we can invalidate the cache.
	row := s.db.QueryRow(ctx,
		"SELECT email FROM suppressions WHERE id = $1 AND org_id = $2",
		suppressionID, orgID,
	)
	var email string
	if err := row.Scan(&email); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrSuppressionNotFound
		}
		return fmt.Errorf("lookup suppression for removal: %w", err)
	}

	// Delete the record.
	if err := s.queries.DeleteSuppression(ctx, sqlcdb.DeleteSuppressionParams{
		ID:    suppressionID,
		OrgID: orgID,
	}); err != nil {
		return fmt.Errorf("delete suppression: %w", err)
	}

	// Invalidate cache.
	key := cacheKey(orgID, email)
	if cacheErr := s.cache.Del(ctx, key).Err(); cacheErr != nil {
		s.logger.Warn().Err(cacheErr).Str("key", key).Msg("failed to invalidate suppression cache")
	}

	return nil
}

// List returns a paginated list of suppressions for an organization along
// with the total count.
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, perPage int) ([]sqlcdb.Suppression, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}

	offset := (page - 1) * perPage

	suppressions, err := s.queries.ListSuppressionsByOrg(ctx, sqlcdb.ListSuppressionsByOrgParams{
		OrgID:  orgID,
		Limit:  int32(perPage),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("list suppressions: %w", err)
	}

	total, err := s.queries.CountSuppressionsByOrg(ctx, orgID)
	if err != nil {
		return nil, 0, fmt.Errorf("count suppressions: %w", err)
	}

	return suppressions, total, nil
}
