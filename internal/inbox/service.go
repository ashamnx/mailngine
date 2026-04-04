// Package inbox implements the inbound email inbox for Hello Mail.
// It provides thread management, message CRUD, label operations,
// and email threading for received messages.
package inbox

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	sqlcdb "github.com/hellomail/hellomail/internal/db/sqlcdb"
)

// Sentinel errors for inbox operations.
var (
	ErrThreadNotFound  = errors.New("thread not found")
	ErrMessageNotFound = errors.New("message not found")
	ErrLabelNotFound   = errors.New("label not found")
	ErrLabelExists     = errors.New("label already exists")
)

// Service provides inbox operations for threads, messages, and labels.
type Service struct {
	db      *pgxpool.Pool
	queries *sqlcdb.Queries
	logger  zerolog.Logger
}

// NewService creates a new inbox Service with the given dependencies.
func NewService(db *pgxpool.Pool, logger zerolog.Logger) *Service {
	return &Service{
		db:      db,
		queries: sqlcdb.New(db),
		logger:  logger.With().Str("component", "inbox_service").Logger(),
	}
}

// ListThreads returns a paginated list of threads for an organization.
func (s *Service) ListThreads(ctx context.Context, orgID uuid.UUID, page, perPage int) ([]sqlcdb.InboxThread, int64, error) {
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

	threads, err := s.queries.ListThreads(ctx, sqlcdb.ListThreadsParams{
		OrgID:  orgID,
		Limit:  int32(perPage),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("list threads: %w", err)
	}

	total, err := s.queries.CountThreads(ctx, orgID)
	if err != nil {
		return nil, 0, fmt.Errorf("count threads: %w", err)
	}

	return threads, total, nil
}

// GetThread retrieves a thread and all its messages, scoped to an organization.
func (s *Service) GetThread(ctx context.Context, orgID, threadID uuid.UUID) (*sqlcdb.InboxThread, []sqlcdb.InboxMessage, error) {
	thread, err := s.queries.GetThread(ctx, sqlcdb.GetThreadParams{
		ID:    threadID,
		OrgID: orgID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil, ErrThreadNotFound
		}
		return nil, nil, fmt.Errorf("get thread: %w", err)
	}

	messages, err := s.queries.ListMessagesByThread(ctx, sqlcdb.ListMessagesByThreadParams{
		ThreadID: pgtype.UUID{Bytes: threadID, Valid: true},
		OrgID:    orgID,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("list messages by thread: %w", err)
	}

	return &thread, messages, nil
}

// DeleteThread removes a thread and its associated data, scoped to an organization.
func (s *Service) DeleteThread(ctx context.Context, orgID, threadID uuid.UUID) error {
	// Verify the thread exists and belongs to the org.
	_, err := s.queries.GetThread(ctx, sqlcdb.GetThreadParams{
		ID:    threadID,
		OrgID: orgID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrThreadNotFound
		}
		return fmt.Errorf("get thread: %w", err)
	}

	if err := s.queries.DeleteThread(ctx, sqlcdb.DeleteThreadParams{
		ID:    threadID,
		OrgID: orgID,
	}); err != nil {
		return fmt.Errorf("delete thread: %w", err)
	}

	return nil
}

// GetMessage retrieves a single inbox message, scoped to an organization.
func (s *Service) GetMessage(ctx context.Context, orgID, messageID uuid.UUID) (*sqlcdb.InboxMessage, error) {
	msg, err := s.queries.GetInboxMessage(ctx, sqlcdb.GetInboxMessageParams{
		ID:    messageID,
		OrgID: orgID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrMessageNotFound
		}
		return nil, fmt.Errorf("get message: %w", err)
	}
	return &msg, nil
}

// UpdateMessageFlags performs a partial update on message boolean flags.
// Only fields where the pointer is non-nil are updated.
func (s *Service) UpdateMessageFlags(ctx context.Context, orgID, messageID uuid.UUID, isRead, isStarred, isArchived, isTrashed *bool) error {
	// Fetch the current message to merge flags.
	msg, err := s.queries.GetInboxMessage(ctx, sqlcdb.GetInboxMessageParams{
		ID:    messageID,
		OrgID: orgID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrMessageNotFound
		}
		return fmt.Errorf("get message for flag update: %w", err)
	}

	// Merge: use new value if provided, otherwise keep existing.
	readVal := msg.IsRead
	if isRead != nil {
		readVal = *isRead
	}
	starredVal := msg.IsStarred
	if isStarred != nil {
		starredVal = *isStarred
	}
	archivedVal := msg.IsArchived
	if isArchived != nil {
		archivedVal = *isArchived
	}
	trashedVal := msg.IsTrashed
	if isTrashed != nil {
		trashedVal = *isTrashed
	}

	if err := s.queries.UpdateMessageFlags(ctx, sqlcdb.UpdateMessageFlagsParams{
		ID:         messageID,
		OrgID:      orgID,
		IsRead:     readVal,
		IsStarred:  starredVal,
		IsArchived: archivedVal,
		IsTrashed:  trashedVal,
	}); err != nil {
		return fmt.Errorf("update message flags: %w", err)
	}

	return nil
}

// DeleteMessage removes an inbox message, scoped to an organization.
func (s *Service) DeleteMessage(ctx context.Context, orgID, messageID uuid.UUID) error {
	_, err := s.queries.GetInboxMessage(ctx, sqlcdb.GetInboxMessageParams{
		ID:    messageID,
		OrgID: orgID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrMessageNotFound
		}
		return fmt.Errorf("get message: %w", err)
	}

	if err := s.queries.DeleteInboxMessage(ctx, sqlcdb.DeleteInboxMessageParams{
		ID:    messageID,
		OrgID: orgID,
	}); err != nil {
		return fmt.Errorf("delete message: %w", err)
	}

	return nil
}

// SearchMessages performs a case-insensitive search on message subjects.
func (s *Service) SearchMessages(ctx context.Context, orgID uuid.UUID, query string, page, perPage int) ([]sqlcdb.InboxMessage, error) {
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

	messages, err := s.queries.SearchMessages(ctx, sqlcdb.SearchMessagesParams{
		OrgID:   orgID,
		Column2: pgtype.Text{String: query, Valid: true},
		Limit:   int32(perPage),
		Offset:  int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("search messages: %w", err)
	}

	return messages, nil
}

// CreateLabel creates a new label for an organization.
func (s *Service) CreateLabel(ctx context.Context, orgID uuid.UUID, name, color string) (*sqlcdb.InboxLabel, error) {
	var colorField pgtype.Text
	if color != "" {
		colorField = pgtype.Text{String: color, Valid: true}
	}

	label, err := s.queries.CreateLabel(ctx, sqlcdb.CreateLabelParams{
		OrgID: orgID,
		Name:  name,
		Color: colorField,
	})
	if err != nil {
		// Check for unique constraint violation (org_id, name).
		if isDuplicateKeyError(err) {
			return nil, ErrLabelExists
		}
		return nil, fmt.Errorf("create label: %w", err)
	}

	return &label, nil
}

// ListLabels returns all labels for an organization, ordered by name.
func (s *Service) ListLabels(ctx context.Context, orgID uuid.UUID) ([]sqlcdb.InboxLabel, error) {
	labels, err := s.queries.ListLabelsByOrg(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("list labels: %w", err)
	}
	return labels, nil
}

// DeleteLabel removes a label, scoped to an organization.
func (s *Service) DeleteLabel(ctx context.Context, orgID, labelID uuid.UUID) error {
	if err := s.queries.DeleteLabel(ctx, sqlcdb.DeleteLabelParams{
		ID:    labelID,
		OrgID: orgID,
	}); err != nil {
		return fmt.Errorf("delete label: %w", err)
	}
	return nil
}

// AddMessageLabel associates a label with a message.
func (s *Service) AddMessageLabel(ctx context.Context, messageID, labelID uuid.UUID) error {
	if err := s.queries.AddMessageLabel(ctx, sqlcdb.AddMessageLabelParams{
		MessageID: messageID,
		LabelID:   labelID,
	}); err != nil {
		return fmt.Errorf("add message label: %w", err)
	}
	return nil
}

// RemoveMessageLabel dissociates a label from a message.
func (s *Service) RemoveMessageLabel(ctx context.Context, messageID, labelID uuid.UUID) error {
	if err := s.queries.RemoveMessageLabel(ctx, sqlcdb.RemoveMessageLabelParams{
		MessageID: messageID,
		LabelID:   labelID,
	}); err != nil {
		return fmt.Errorf("remove message label: %w", err)
	}
	return nil
}

// isDuplicateKeyError checks if the error is a PostgreSQL unique constraint violation (23505).
func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	// pgx wraps PostgreSQL errors as *pgconn.PgError.
	var pgErr interface{ SQLState() string }
	if errors.As(err, &pgErr) {
		return pgErr.SQLState() == "23505"
	}
	return false
}
