// Package email implements the email sending pipeline for Mailngine.
// It handles validation, queuing, and orchestration of outbound email delivery.
package email

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	sqlcdb "github.com/mailngine/mailngine/internal/db/sqlcdb"
	"github.com/mailngine/mailngine/internal/queue"
)

// SendRequest represents the payload for sending an email via the API.
type SendRequest struct {
	From           string            `json:"from" validate:"required,email"`
	To             []string          `json:"to" validate:"required,min=1,dive,email"`
	CC             []string          `json:"cc,omitempty" validate:"omitempty,dive,email"`
	BCC            []string          `json:"bcc,omitempty" validate:"omitempty,dive,email"`
	ReplyTo        string            `json:"reply_to,omitempty" validate:"omitempty,email"`
	Subject        string            `json:"subject" validate:"required"`
	HTML           string            `json:"html,omitempty"`
	Text           string            `json:"text,omitempty"`
	Headers        map[string]string `json:"headers,omitempty"`
	Tags           []string          `json:"tags,omitempty"`
	IdempotencyKey string            `json:"idempotency_key,omitempty"`
	ScheduledAt    *time.Time        `json:"scheduled_at,omitempty"`
}

// EmailResponse is the API response representation of an email record.
type EmailResponse struct {
	ID          uuid.UUID  `json:"id"`
	From        string     `json:"from"`
	To          []string   `json:"to"`
	CC          []string   `json:"cc,omitempty"`
	BCC         []string   `json:"bcc,omitempty"`
	Subject     string     `json:"subject"`
	Status      string     `json:"status"`
	MessageID   string     `json:"message_id,omitempty"`
	ScheduledAt *time.Time `json:"scheduled_at,omitempty"`
	SentAt      *time.Time `json:"sent_at,omitempty"`
	DeliveredAt *time.Time `json:"delivered_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

// SendTaskPayload is the JSON payload enqueued for the email:send asynq task.
type SendTaskPayload struct {
	EmailID uuid.UUID `json:"email_id"`
}

// ErrDomainNotVerified is returned when the sender's domain is not verified.
var ErrDomainNotVerified = errors.New("domain is not verified for this organization")

// ErrInvalidFromAddress is returned when the from address cannot be parsed.
var ErrInvalidFromAddress = errors.New("invalid from address")

// ErrEmailNotFound is returned when an email record does not exist.
var ErrEmailNotFound = errors.New("email not found")

// ErrRecipientSuppressed is returned when one or more recipients are on the suppression list.
var ErrRecipientSuppressed = errors.New("recipient is suppressed")

// SuppressionChecker checks if an email address is suppressed for an organization.
type SuppressionChecker interface {
	Check(ctx context.Context, orgID uuid.UUID, email string) (bool, error)
}

// Service provides email sending operations.
type Service struct {
	db             *pgxpool.Pool
	queries        *sqlcdb.Queries
	queue          *asynq.Client
	suppressionSvc SuppressionChecker
	logger         zerolog.Logger
}

// NewService creates a new email Service with the given dependencies.
func NewService(db *pgxpool.Pool, queue *asynq.Client, suppressionSvc SuppressionChecker, logger zerolog.Logger) *Service {
	return &Service{
		db:             db,
		queries:        sqlcdb.New(db),
		queue:          queue,
		suppressionSvc: suppressionSvc,
		logger:         logger.With().Str("component", "email_service").Logger(),
	}
}

// SendEmail validates the request, persists the email record, and enqueues it for delivery.
func (s *Service) SendEmail(ctx context.Context, orgID uuid.UUID, apiKeyID *uuid.UUID, req *SendRequest) (*sqlcdb.Email, error) {
	// Extract domain from the "from" address.
	domainName, err := extractDomain(req.From)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidFromAddress, err)
	}

	// Verify the domain is registered and verified for this org.
	domain, err := s.queries.GetVerifiedDomainByName(ctx, sqlcdb.GetVerifiedDomainByNameParams{
		Name:  domainName,
		OrgID: orgID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: %s", ErrDomainNotVerified, domainName)
		}
		return nil, fmt.Errorf("lookup domain: %w", err)
	}

	// Idempotency check: if a key is provided and an email already exists, return it.
	if req.IdempotencyKey != "" {
		existing, err := s.queries.GetEmailByIdempotencyKey(ctx, sqlcdb.GetEmailByIdempotencyKeyParams{
			OrgID:          orgID,
			IdempotencyKey: pgtype.Text{String: req.IdempotencyKey, Valid: true},
		})
		if err == nil {
			return &existing, nil
		}
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("check idempotency key: %w", err)
		}
	}

	// Check all recipients against the suppression list.
	allRecipients := make([]string, 0, len(req.To)+len(req.CC)+len(req.BCC))
	allRecipients = append(allRecipients, req.To...)
	allRecipients = append(allRecipients, req.CC...)
	allRecipients = append(allRecipients, req.BCC...)

	for _, rcpt := range allRecipients {
		suppressed, err := s.suppressionSvc.Check(ctx, orgID, rcpt)
		if err != nil {
			return nil, fmt.Errorf("check suppression for %s: %w", rcpt, err)
		}
		if suppressed {
			return nil, fmt.Errorf("%w: %s", ErrRecipientSuppressed, rcpt)
		}
	}

	// Generate a unique Message-ID per RFC 5322.
	messageID := fmt.Sprintf("<%s@mailngine.com>", uuid.New().String())

	// Marshal JSON fields.
	toJSON, err := json.Marshal(req.To)
	if err != nil {
		return nil, fmt.Errorf("marshal to addresses: %w", err)
	}

	var ccJSON, bccJSON, headersJSON, tagsJSON []byte
	if len(req.CC) > 0 {
		ccJSON, err = json.Marshal(req.CC)
		if err != nil {
			return nil, fmt.Errorf("marshal cc addresses: %w", err)
		}
	}
	if len(req.BCC) > 0 {
		bccJSON, err = json.Marshal(req.BCC)
		if err != nil {
			return nil, fmt.Errorf("marshal bcc addresses: %w", err)
		}
	}
	if len(req.Headers) > 0 {
		headersJSON, err = json.Marshal(req.Headers)
		if err != nil {
			return nil, fmt.Errorf("marshal headers: %w", err)
		}
	}
	if len(req.Tags) > 0 {
		tagsJSON, err = json.Marshal(req.Tags)
		if err != nil {
			return nil, fmt.Errorf("marshal tags: %w", err)
		}
	}

	// Build optional pgtype fields.
	var apiKeyUUID pgtype.UUID
	if apiKeyID != nil {
		apiKeyUUID = pgtype.UUID{Bytes: *apiKeyID, Valid: true}
	}

	var idempotencyKey pgtype.Text
	if req.IdempotencyKey != "" {
		idempotencyKey = pgtype.Text{String: req.IdempotencyKey, Valid: true}
	}

	var replyTo pgtype.Text
	if req.ReplyTo != "" {
		replyTo = pgtype.Text{String: req.ReplyTo, Valid: true}
	}

	var textBodyKey pgtype.Text
	if req.Text != "" {
		textBodyKey = pgtype.Text{String: req.Text, Valid: true}
	}

	var htmlBodyKey pgtype.Text
	if req.HTML != "" {
		htmlBodyKey = pgtype.Text{String: req.HTML, Valid: true}
	}

	var scheduledAt pgtype.Timestamptz
	if req.ScheduledAt != nil {
		scheduledAt = pgtype.Timestamptz{Time: *req.ScheduledAt, Valid: true}
	}

	// Extract from_name if present (e.g. "John Doe <john@example.com>").
	fromName, fromAddr := parseFromAddress(req.From)
	var fromNameField pgtype.Text
	if fromName != "" {
		fromNameField = pgtype.Text{String: fromName, Valid: true}
	}

	// Insert the email record.
	email, err := s.queries.CreateEmail(ctx, sqlcdb.CreateEmailParams{
		OrgID:          orgID,
		DomainID:       domain.ID,
		ApiKeyID:       apiKeyUUID,
		IdempotencyKey: idempotencyKey,
		FromAddress:    fromAddr,
		FromName:       fromNameField,
		ToAddresses:    toJSON,
		CcAddresses:    ccJSON,
		BccAddresses:   bccJSON,
		ReplyTo:        replyTo,
		Subject:        req.Subject,
		TextBodyKey:    textBodyKey,
		HtmlBodyKey:    htmlBodyKey,
		Headers:        headersJSON,
		Tags:           tagsJSON,
		TemplateID:     pgtype.UUID{},
		TemplateData:   nil,
		Status:         "queued",
		ScheduledAt:    scheduledAt,
		MessageID:      pgtype.Text{String: messageID, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("create email record: %w", err)
	}

	// Enqueue the send task.
	payload, err := json.Marshal(SendTaskPayload{EmailID: email.ID})
	if err != nil {
		return nil, fmt.Errorf("marshal task payload: %w", err)
	}

	opts := []asynq.Option{
		asynq.MaxRetry(3),
		asynq.Queue("default"),
	}
	if req.ScheduledAt != nil {
		opts = append(opts, asynq.ProcessAt(*req.ScheduledAt))
	}

	task := asynq.NewTask(queue.TaskSendEmail, payload, opts...)
	if _, err := s.queue.EnqueueContext(ctx, task); err != nil {
		s.logger.Error().Err(err).Str("email_id", email.ID.String()).Msg("failed to enqueue email task")
		// The email is already persisted as "queued"; a retry mechanism can pick it up.
		// We do not return an error to the caller since the record was created.
	}

	return &email, nil
}

// GetEmail retrieves a single email by ID scoped to an organization.
func (s *Service) GetEmail(ctx context.Context, orgID, emailID uuid.UUID) (*sqlcdb.Email, error) {
	email, err := s.queries.GetEmail(ctx, sqlcdb.GetEmailParams{
		ID:    emailID,
		OrgID: orgID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrEmailNotFound
		}
		return nil, fmt.Errorf("get email: %w", err)
	}
	return &email, nil
}

// ListEmails returns a paginated list of emails for an organization.
func (s *Service) ListEmails(ctx context.Context, orgID uuid.UUID, page, perPage int) ([]sqlcdb.Email, int64, error) {
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

	emails, err := s.queries.ListEmailsByOrg(ctx, sqlcdb.ListEmailsByOrgParams{
		OrgID:  orgID,
		Limit:  int32(perPage),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("list emails: %w", err)
	}

	total, err := s.queries.CountEmailsByOrg(ctx, orgID)
	if err != nil {
		return nil, 0, fmt.Errorf("count emails: %w", err)
	}

	return emails, total, nil
}

// ToResponse converts a database Email record into an API response.
func ToResponse(e *sqlcdb.Email) EmailResponse {
	resp := EmailResponse{
		ID:        e.ID,
		From:      e.FromAddress,
		Subject:   e.Subject,
		Status:    e.Status,
		CreatedAt: e.CreatedAt,
	}

	// Unmarshal JSONB fields.
	if e.ToAddresses != nil {
		_ = json.Unmarshal(e.ToAddresses, &resp.To)
	}
	if e.CcAddresses != nil {
		_ = json.Unmarshal(e.CcAddresses, &resp.CC)
	}
	if e.BccAddresses != nil {
		_ = json.Unmarshal(e.BccAddresses, &resp.BCC)
	}
	if e.MessageID.Valid {
		resp.MessageID = e.MessageID.String
	}
	if e.ScheduledAt.Valid {
		t := e.ScheduledAt.Time
		resp.ScheduledAt = &t
	}
	if e.SentAt.Valid {
		t := e.SentAt.Time
		resp.SentAt = &t
	}
	if e.DeliveredAt.Valid {
		t := e.DeliveredAt.Time
		resp.DeliveredAt = &t
	}

	return resp
}

// extractDomain returns the domain part of an email address.
func extractDomain(addr string) (string, error) {
	// Handle "Name <email>" format.
	if idx := strings.LastIndex(addr, "<"); idx != -1 {
		end := strings.LastIndex(addr, ">")
		if end == -1 || end <= idx {
			return "", fmt.Errorf("malformed address: %s", addr)
		}
		addr = addr[idx+1 : end]
	}

	parts := strings.SplitN(addr, "@", 2)
	if len(parts) != 2 || parts[1] == "" {
		return "", fmt.Errorf("invalid email address: %s", addr)
	}
	return strings.ToLower(parts[1]), nil
}

// parseFromAddress splits a "Name <email>" formatted address into name and bare email.
// If no angle brackets are present, it returns an empty name and the original address.
func parseFromAddress(addr string) (name, email string) {
	if idx := strings.LastIndex(addr, "<"); idx != -1 {
		end := strings.LastIndex(addr, ">")
		if end > idx {
			name = strings.TrimSpace(addr[:idx])
			// Remove surrounding quotes from the display name if present.
			name = strings.Trim(name, "\"")
			email = addr[idx+1 : end]
			return name, email
		}
	}
	return "", addr
}
