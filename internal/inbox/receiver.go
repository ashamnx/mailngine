package inbox

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	sqlcdb "github.com/hellomail/hellomail/internal/db/sqlcdb"
)

// ReceiveMessage processes an incoming email by assigning it to a thread
// and persisting it as an inbox message.
func (s *Service) ReceiveMessage(ctx context.Context, orgID, domainID uuid.UUID, msg *IncomingMessage) (*sqlcdb.InboxMessage, error) {
	// Assign or create a thread for this message.
	threadID, err := s.AssignThread(ctx, orgID, domainID, msg)
	if err != nil {
		return nil, fmt.Errorf("assign thread: %w", err)
	}

	// Prepare nullable fields.
	var threadUUID pgtype.UUID
	if threadID != uuid.Nil {
		threadUUID = pgtype.UUID{Bytes: threadID, Valid: true}
	}

	var messageIDHeader pgtype.Text
	if msg.MessageID != "" {
		messageIDHeader = pgtype.Text{String: msg.MessageID, Valid: true}
	}

	var inReplyTo pgtype.Text
	if msg.InReplyTo != "" {
		inReplyTo = pgtype.Text{String: msg.InReplyTo, Valid: true}
	}

	var referencesHeader pgtype.Text
	if msg.References != "" {
		referencesHeader = pgtype.Text{String: msg.References, Valid: true}
	}

	var fromName pgtype.Text
	if msg.FromName != "" {
		fromName = pgtype.Text{String: msg.FromName, Valid: true}
	}

	toJSON, err := json.Marshal(msg.To)
	if err != nil {
		return nil, fmt.Errorf("marshal to addresses: %w", err)
	}

	var ccJSON []byte
	if len(msg.CC) > 0 {
		ccJSON, err = json.Marshal(msg.CC)
		if err != nil {
			return nil, fmt.Errorf("marshal cc addresses: %w", err)
		}
	}

	var textBodyKey pgtype.Text
	if msg.TextBody != "" {
		textBodyKey = pgtype.Text{String: msg.TextBody, Valid: true}
	}

	var htmlBodyKey pgtype.Text
	if msg.HTMLBody != "" {
		htmlBodyKey = pgtype.Text{String: msg.HTMLBody, Valid: true}
	}

	var snippet pgtype.Text
	if msg.Snippet != "" {
		snippet = pgtype.Text{String: msg.Snippet, Valid: true}
	}

	receivedAt := msg.ReceivedAt
	if receivedAt.IsZero() {
		receivedAt = time.Now().UTC()
	}

	created, err := s.queries.CreateInboxMessage(ctx, sqlcdb.CreateInboxMessageParams{
		OrgID:            orgID,
		DomainID:         domainID,
		ThreadID:         threadUUID,
		MessageIDHeader:  messageIDHeader,
		InReplyTo:        inReplyTo,
		ReferencesHeader: referencesHeader,
		FromAddress:      msg.From,
		FromName:         fromName,
		ToAddresses:      toJSON,
		CcAddresses:      ccJSON,
		Subject:          msg.Subject,
		TextBodyKey:      textBodyKey,
		HtmlBodyKey:      htmlBodyKey,
		Snippet:          snippet,
		ReceivedAt:       receivedAt,
	})
	if err != nil {
		return nil, fmt.Errorf("create inbox message: %w", err)
	}

	return &created, nil
}
