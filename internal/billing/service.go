package billing

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	"github.com/mailngine/mailngine/internal/db/sqlcdb"
)

// PlanDefinition describes a billing plan's properties.
type PlanDefinition struct {
	Name         string `json:"name"`
	MonthlyLimit int    `json:"monthly_limit"` // 0 = unlimited
	MaxDomains   int    `json:"max_domains"`   // 0 = unlimited
	Price        int    `json:"price"`         // cents, 0 for free
}

// Plans defines the available billing plans.
var Plans = map[string]PlanDefinition{
	"free":       {Name: "Free", MonthlyLimit: 100, MaxDomains: 1, Price: 0},
	"starter":    {Name: "Starter", MonthlyLimit: 10000, MaxDomains: 5, Price: 2000},
	"pro":        {Name: "Pro", MonthlyLimit: 100000, MaxDomains: 0, Price: 8000},
	"enterprise": {Name: "Enterprise", MonthlyLimit: 0, MaxDomains: 0, Price: 0},
}

// UsageResponse represents the current month's usage summary.
type UsageResponse struct {
	TotalSent       int     `json:"total_sent"`
	TotalReceived   int     `json:"total_received"`
	TotalAPICalls   int     `json:"total_api_calls"`
	PlanLimit       int     `json:"plan_limit"`
	UsagePercentage float64 `json:"usage_percentage"`
}

// PlanResponse represents the organization's plan details and current usage.
type PlanResponse struct {
	Plan           string `json:"plan"`
	Name           string `json:"name"`
	MonthlyLimit   int    `json:"monthly_limit"`
	MaxDomains     int    `json:"max_domains"`
	OverageEnabled bool   `json:"overage_enabled"`
	CurrentUsage   int    `json:"current_usage"`
}

// Service handles billing and usage operations.
type Service struct {
	db      *pgxpool.Pool
	queries *sqlcdb.Queries
	logger  zerolog.Logger
}

// NewService creates a new billing Service.
func NewService(db *pgxpool.Pool, logger zerolog.Logger) *Service {
	return &Service{
		db:      db,
		queries: sqlcdb.New(db),
		logger:  logger,
	}
}

// GetCurrentUsage returns the current month's usage for an organization.
func (s *Service) GetCurrentUsage(ctx context.Context, orgID uuid.UUID) (*UsageResponse, error) {
	org, err := s.queries.GetOrganization(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("get organization: %w", err)
	}

	usage, err := s.queries.GetCurrentMonthUsage(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("get current month usage: %w", err)
	}

	planLimit := int(org.MonthlyLimit)
	var usagePercentage float64
	if planLimit > 0 {
		usagePercentage = float64(usage.TotalSent) / float64(planLimit) * 100
	}

	return &UsageResponse{
		TotalSent:       int(usage.TotalSent),
		TotalReceived:   int(usage.TotalReceived),
		TotalAPICalls:   int(usage.TotalApiCalls),
		PlanLimit:       planLimit,
		UsagePercentage: usagePercentage,
	}, nil
}

// GetUsageHistory returns the monthly usage history for an organization.
func (s *Service) GetUsageHistory(ctx context.Context, orgID uuid.UUID, limit int) ([]sqlcdb.UsageMonthly, error) {
	return s.queries.GetMonthlyUsageHistory(ctx, sqlcdb.GetMonthlyUsageHistoryParams{
		OrgID: orgID,
		Limit: int32(limit),
	})
}

// GetPlan returns the organization's plan details and current usage.
func (s *Service) GetPlan(ctx context.Context, orgID uuid.UUID) (*PlanResponse, error) {
	org, err := s.queries.GetOrganization(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("get organization: %w", err)
	}

	plan, ok := Plans[org.Plan]
	if !ok {
		plan = Plans["free"]
	}

	usage, err := s.queries.GetCurrentMonthUsage(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("get current month usage: %w", err)
	}

	return &PlanResponse{
		Plan:           org.Plan,
		Name:           plan.Name,
		MonthlyLimit:   int(org.MonthlyLimit),
		MaxDomains:     plan.MaxDomains,
		OverageEnabled: org.OverageEnabled,
		CurrentUsage:   int(usage.TotalSent),
	}, nil
}

// CheckLimit verifies whether the organization is within its monthly sending limit.
// Returns true if the organization can still send emails.
func (s *Service) CheckLimit(ctx context.Context, orgID uuid.UUID) (bool, error) {
	org, err := s.queries.GetOrganization(ctx, orgID)
	if err != nil {
		return false, fmt.Errorf("get organization: %w", err)
	}

	// 0 means unlimited (enterprise plan).
	if org.MonthlyLimit == 0 {
		return true, nil
	}

	usage, err := s.queries.GetCurrentMonthUsage(ctx, orgID)
	if err != nil {
		return false, fmt.Errorf("get current month usage: %w", err)
	}

	if int(usage.TotalSent) >= int(org.MonthlyLimit) {
		// Allow overage if enabled.
		return org.OverageEnabled, nil
	}

	return true, nil
}
