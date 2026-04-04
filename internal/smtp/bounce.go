// Package smtp implements the Postfix integration layer for Hello Mail,
// including bounce processing, FBL handling, and milter business logic.
package smtp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	sqlcdb "github.com/hellomail/hellomail/internal/db/sqlcdb"
	"github.com/hellomail/hellomail/internal/suppression"
)

// BounceProcessor handles DSN (Delivery Status Notification) messages from
// Postfix. It classifies bounces, records events, and manages the suppression
// list for hard bounces.
type BounceProcessor struct {
	db             *pgxpool.Pool
	queries        *sqlcdb.Queries
	suppressionSvc *suppression.Service
	logger         zerolog.Logger
}

// NewBounceProcessor creates a new BounceProcessor with the given dependencies.
func NewBounceProcessor(db *pgxpool.Pool, suppressionSvc *suppression.Service, logger zerolog.Logger) *BounceProcessor {
	return &BounceProcessor{
		db:             db,
		queries:        sqlcdb.New(db),
		suppressionSvc: suppressionSvc,
		logger:         logger.With().Str("component", "bounce_processor").Logger(),
	}
}

// bounceType classifies a bounce into hard or soft based on the SMTP reply code.
type bounceType string

const (
	bounceHard bounceType = "hard"
	bounceSoft bounceType = "soft"
)

// classifyBounce returns whether a bounce is hard (5xx) or soft (4xx).
// Unknown codes default to soft to avoid premature suppression.
func classifyBounce(code string) bounceType {
	if strings.HasPrefix(code, "5") {
		return bounceHard
	}
	return bounceSoft
}

// ProcessBounce handles a single bounce notification. It:
//  1. Classifies the bounce as hard (5xx) or soft (4xx).
//  2. Records an email_event with the appropriate type.
//  3. For hard bounces: updates the email status to "bounced" and adds
//     the recipient to the organization's suppression list.
func (bp *BounceProcessor) ProcessBounce(ctx context.Context, emailID, orgID uuid.UUID, recipient, bounceCode, bounceMessage string) error {
	bt := classifyBounce(bounceCode)

	// Determine event type based on bounce classification.
	eventType := "deferred"
	if bt == bounceHard {
		eventType = "bounced"
	}

	// Build event metadata.
	metadata, err := json.Marshal(map[string]string{
		"bounce_code":    bounceCode,
		"bounce_message": bounceMessage,
		"bounce_type":    string(bt),
	})
	if err != nil {
		return fmt.Errorf("marshal bounce metadata: %w", err)
	}

	// Record the email event.
	if _, err := bp.queries.CreateEmailEvent(ctx, sqlcdb.CreateEmailEventParams{
		EmailID:   emailID,
		OrgID:     orgID,
		EventType: eventType,
		Recipient: recipient,
		Metadata:  metadata,
		IpAddress: "",
		UserAgent: pgtype.Text{},
	}); err != nil {
		return fmt.Errorf("create bounce event: %w", err)
	}

	bp.logger.Info().
		Str("email_id", emailID.String()).
		Str("recipient", recipient).
		Str("bounce_type", string(bt)).
		Str("bounce_code", bounceCode).
		Msg("bounce event recorded")

	// For hard bounces, update the email status and suppress the recipient.
	if bt == bounceHard {
		if err := bp.queries.UpdateEmailStatus(ctx, sqlcdb.UpdateEmailStatusParams{
			ID:     emailID,
			Status: "bounced",
		}); err != nil {
			return fmt.Errorf("update email status to bounced: %w", err)
		}

		if err := bp.suppressionSvc.Add(ctx, orgID, recipient, "hard_bounce", map[string]any{
			"bounce_code":    bounceCode,
			"bounce_message": bounceMessage,
			"email_id":       emailID.String(),
		}); err != nil {
			return fmt.Errorf("add recipient to suppression list: %w", err)
		}

		bp.logger.Info().
			Str("email_id", emailID.String()).
			Str("recipient", recipient).
			Msg("recipient suppressed due to hard bounce")
	}

	return nil
}
