package smtp

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/hellomail/hellomail/internal/suppression"
)

// FBLProcessor handles Feedback Loop (FBL) reports. ISPs send these when a
// recipient marks a message as spam. The processor adds the complaining
// recipient to the organization's suppression list.
type FBLProcessor struct {
	suppressionSvc *suppression.Service
	logger         zerolog.Logger
}

// NewFBLProcessor creates a new FBLProcessor with the given dependencies.
func NewFBLProcessor(suppressionSvc *suppression.Service, logger zerolog.Logger) *FBLProcessor {
	return &FBLProcessor{
		suppressionSvc: suppressionSvc,
		logger:         logger.With().Str("component", "fbl_processor").Logger(),
	}
}

// ProcessComplaint handles a single spam complaint by suppressing the
// recipient for the given organization.
func (fp *FBLProcessor) ProcessComplaint(ctx context.Context, orgID uuid.UUID, recipient string) error {
	if err := fp.suppressionSvc.Add(ctx, orgID, recipient, "complaint", nil); err != nil {
		return fmt.Errorf("suppress complaint recipient: %w", err)
	}

	fp.logger.Info().
		Str("org_id", orgID.String()).
		Str("recipient", recipient).
		Msg("recipient suppressed due to spam complaint")

	return nil
}
