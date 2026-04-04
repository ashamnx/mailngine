package analytics

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"

	"github.com/hellomail/hellomail/internal/db/sqlcdb"
)

type Service struct {
	db      *pgxpool.Pool
	queries *sqlcdb.Queries
	cache   *redis.Client
	logger  zerolog.Logger
}

func NewService(db *pgxpool.Pool, cache *redis.Client, logger zerolog.Logger) *Service {
	return &Service{
		db:      db,
		queries: sqlcdb.New(db),
		cache:   cache,
		logger:  logger,
	}
}

type OverviewResponse struct {
	TotalSent      int     `json:"total_sent"`
	TotalDelivered int     `json:"total_delivered"`
	TotalBounced   int     `json:"total_bounced"`
	TotalReceived  int     `json:"total_received"`
	DeliveryRate   float64 `json:"delivery_rate"`
	BounceRate     float64 `json:"bounce_rate"`
}

type TimeseriesPoint struct {
	Date            string `json:"date"`
	EmailsSent      int    `json:"emails_sent"`
	EmailsDelivered int    `json:"emails_delivered"`
	EmailsBounced   int    `json:"emails_bounced"`
	EmailsReceived  int    `json:"emails_received"`
}

func (s *Service) GetOverview(ctx context.Context, orgID uuid.UUID, from, to time.Time) (*OverviewResponse, error) {
	summary, err := s.queries.GetUsageSummary(ctx, sqlcdb.GetUsageSummaryParams{
		OrgID:   orgID,
		Date:    pgtype.Date{Time: from, Valid: true},
		Date_2:  pgtype.Date{Time: to, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("get usage summary: %w", err)
	}

	resp := &OverviewResponse{
		TotalSent:      int(summary.TotalSent),
		TotalDelivered: int(summary.TotalDelivered),
		TotalBounced:   int(summary.TotalBounced),
		TotalReceived:  int(summary.TotalReceived),
	}

	if resp.TotalSent > 0 {
		resp.DeliveryRate = float64(resp.TotalDelivered) / float64(resp.TotalSent) * 100
		resp.BounceRate = float64(resp.TotalBounced) / float64(resp.TotalSent) * 100
	}

	return resp, nil
}

func (s *Service) GetTimeseries(ctx context.Context, orgID uuid.UUID, from, to time.Time) ([]TimeseriesPoint, error) {
	rows, err := s.queries.GetUsageDaily(ctx, sqlcdb.GetUsageDailyParams{
		OrgID:  orgID,
		Date:   pgtype.Date{Time: from, Valid: true},
		Date_2: pgtype.Date{Time: to, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("get usage daily: %w", err)
	}

	points := make([]TimeseriesPoint, len(rows))
	for i, r := range rows {
		points[i] = TimeseriesPoint{
			Date:            r.Date.Time.Format("2006-01-02"),
			EmailsSent:      int(r.EmailsSent),
			EmailsDelivered: int(r.EmailsDelivered),
			EmailsBounced:   int(r.EmailsBounced),
			EmailsReceived:  int(r.EmailsReceived),
		}
	}
	return points, nil
}

func (s *Service) GetEventBreakdown(ctx context.Context, orgID uuid.UUID, from, to time.Time) (map[string]int64, error) {
	rows, err := s.queries.CountEmailEventsByType(ctx, sqlcdb.CountEmailEventsByTypeParams{
		OrgID:        orgID,
		OccurredAt:   from,
		OccurredAt_2: to,
	})
	if err != nil {
		return nil, fmt.Errorf("count events by type: %w", err)
	}

	breakdown := make(map[string]int64)
	for _, r := range rows {
		breakdown[r.EventType] = r.Count
	}
	return breakdown, nil
}

func (s *Service) RecordUsage(ctx context.Context, orgID uuid.UUID, field string, count int) error {
	today := time.Now().UTC().Truncate(24 * time.Hour)
	params := sqlcdb.UpsertUsageDailyParams{
		OrgID: orgID,
		Date:  pgtype.Date{Time: today, Valid: true},
	}

	switch field {
	case "emails_sent":
		params.EmailsSent = int32(count)
	case "emails_delivered":
		params.EmailsDelivered = int32(count)
	case "emails_bounced":
		params.EmailsBounced = int32(count)
	case "emails_received":
		params.EmailsReceived = int32(count)
	case "api_calls":
		params.ApiCalls = int32(count)
	}

	return s.queries.UpsertUsageDaily(ctx, params)
}
