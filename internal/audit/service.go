package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	"github.com/mailngine/mailngine/internal/auth"
	"github.com/mailngine/mailngine/internal/db/sqlcdb"
)

type Service struct {
	db      *pgxpool.Pool
	queries *sqlcdb.Queries
	logger  zerolog.Logger
}

func NewService(db *pgxpool.Pool, logger zerolog.Logger) *Service {
	return &Service{
		db:      db,
		queries: sqlcdb.New(db),
		logger:  logger,
	}
}

func (s *Service) Log(ctx context.Context, action, resourceType string, resourceID *uuid.UUID, metadata map[string]any) {
	orgID := auth.OrgIDFromContext(ctx)
	userID := auth.UserIDFromContext(ctx)
	apiKeyID := auth.APIKeyIDFromContext(ctx)

	var metadataJSON []byte
	if metadata != nil {
		metadataJSON, _ = json.Marshal(metadata)
	}

	var userIDField pgtype.UUID
	if userID != uuid.Nil {
		userIDField = pgtype.UUID{Bytes: userID, Valid: true}
	}

	var apiKeyIDField pgtype.UUID
	if apiKeyID != nil {
		apiKeyIDField = pgtype.UUID{Bytes: *apiKeyID, Valid: true}
	}

	var resourceIDField pgtype.UUID
	if resourceID != nil {
		resourceIDField = pgtype.UUID{Bytes: *resourceID, Valid: true}
	}

	var ipAddress string
	if r, ok := ctx.Value(httpRequestKey).(*http.Request); ok {
		ipAddress = r.RemoteAddr
	}

	var userAgent pgtype.Text
	if r, ok := ctx.Value(httpRequestKey).(*http.Request); ok {
		userAgent = pgtype.Text{String: r.UserAgent(), Valid: true}
	}

	// Fire-and-forget: don't block the request
	go func() {
		if _, err := s.queries.CreateAuditLog(context.Background(), sqlcdb.CreateAuditLogParams{
			OrgID:        orgID,
			UserID:       userIDField,
			ApiKeyID:     apiKeyIDField,
			Action:       action,
			ResourceType: resourceType,
			ResourceID:   resourceIDField,
			Metadata:     metadataJSON,
			IpAddress:    ipAddress,
			UserAgent:    userAgent,
		}); err != nil {
			s.logger.Error().Err(err).Str("action", action).Msg("failed to create audit log")
		}
	}()
}

func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, perPage int) ([]sqlcdb.AuditLog, int64, error) {
	offset := (page - 1) * perPage
	logs, err := s.queries.ListAuditLogs(ctx, sqlcdb.ListAuditLogsParams{
		OrgID:  orgID,
		Limit:  int32(perPage),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("list audit logs: %w", err)
	}

	count, err := s.queries.CountAuditLogs(ctx, orgID)
	if err != nil {
		return nil, 0, fmt.Errorf("count audit logs: %w", err)
	}

	return logs, count, nil
}

func (s *Service) Get(ctx context.Context, orgID, logID uuid.UUID) (*sqlcdb.AuditLog, error) {
	log, err := s.queries.GetAuditLog(ctx, sqlcdb.GetAuditLogParams{ID: logID, OrgID: orgID})
	if err != nil {
		return nil, fmt.Errorf("get audit log: %w", err)
	}
	return &log, nil
}

type ctxKey string

const httpRequestKey ctxKey = "http_request"

func WithHTTPRequest(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, httpRequestKey, r)
}
