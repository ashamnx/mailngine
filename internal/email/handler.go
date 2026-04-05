package email

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	sqlcdb "github.com/mailngine/mailngine/internal/db/sqlcdb"
	"github.com/mailngine/mailngine/internal/queue"
)

// SendHandler processes email:send asynq tasks.
// It loads the email from the database, builds the MIME message,
// signs it with DKIM, and delivers it via SMTP.
type SendHandler struct {
	db      *pgxpool.Pool
	queries *sqlcdb.Queries
	smtp    *SMTPSender
	logger  zerolog.Logger
}

// NewSendHandler creates a new handler for email:send tasks.
func NewSendHandler(db *pgxpool.Pool, smtp *SMTPSender, logger zerolog.Logger) *SendHandler {
	return &SendHandler{
		db:      db,
		queries: sqlcdb.New(db),
		smtp:    smtp,
		logger:  logger.With().Str("component", "email_send_handler").Logger(),
	}
}

// Register adds the handler to the asynq mux.
func (h *SendHandler) Register(mux *asynq.ServeMux) {
	mux.HandleFunc(queue.TaskSendEmail, h.ProcessTask)
}

// ProcessTask handles a single email:send task.
func (h *SendHandler) ProcessTask(ctx context.Context, task *asynq.Task) error {
	var payload SendTaskPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("unmarshal payload: %w", err)
	}

	logger := h.logger.With().Str("email_id", payload.EmailID.String()).Logger()
	logger.Info().Msg("processing email:send task")

	// 1. Look up org_id for this email.
	row, err := h.queries.GetEmailOrgID(ctx, payload.EmailID)
	if err != nil {
		return fmt.Errorf("get email org_id: %w", err)
	}

	// 2. Load the full email record.
	email, err := h.queries.GetEmail(ctx, sqlcdb.GetEmailParams{
		ID:    payload.EmailID,
		OrgID: row.OrgID,
	})
	if err != nil {
		return fmt.Errorf("get email: %w", err)
	}

	// Skip if already sent (idempotency on retry).
	if email.Status != "queued" {
		logger.Info().Str("status", email.Status).Msg("email already processed, skipping")
		return nil
	}

	// 3. Extract domain from the from address and load domain record (for DKIM).
	parts := strings.SplitN(email.FromAddress, "@", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid from address: %s", email.FromAddress)
	}
	domainName := parts[1]

	domain, err := h.queries.GetVerifiedDomainByName(ctx, sqlcdb.GetVerifiedDomainByNameParams{
		Name:  domainName,
		OrgID: row.OrgID,
	})
	if err != nil {
		return fmt.Errorf("get verified domain %s: %w", domainName, err)
	}

	// 4. Parse recipient lists from JSON.
	var toAddrs, ccAddrs, bccAddrs []string
	if err := json.Unmarshal(email.ToAddresses, &toAddrs); err != nil {
		return fmt.Errorf("unmarshal to addresses: %w", err)
	}
	if len(email.CcAddresses) > 0 {
		_ = json.Unmarshal(email.CcAddresses, &ccAddrs)
	}
	if len(email.BccAddresses) > 0 {
		_ = json.Unmarshal(email.BccAddresses, &bccAddrs)
	}

	// 5. Parse custom headers.
	var headers map[string]string
	if len(email.Headers) > 0 {
		_ = json.Unmarshal(email.Headers, &headers)
	}

	// 6. Build the from display value.
	from := email.FromAddress
	if email.FromName.Valid && email.FromName.String != "" {
		from = fmt.Sprintf("%s <%s>", email.FromName.String, email.FromAddress)
	}

	// 7. Read body content.
	var textBody, htmlBody string
	if email.TextBodyKey.Valid {
		textBody = email.TextBodyKey.String
	}
	if email.HtmlBodyKey.Valid {
		htmlBody = email.HtmlBodyKey.String
	}

	// 8. Build the MIME message.
	messageID := ""
	if email.MessageID.Valid {
		messageID = email.MessageID.String
	}

	mimeMsg, err := BuildMIMEMessage(from, email.Subject, textBody, htmlBody, toAddrs, ccAddrs, bccAddrs,
		email.ReplyTo.String, headers, messageID)
	if err != nil {
		return fmt.Errorf("build MIME message: %w", err)
	}

	// 9. Sign with DKIM.
	if domain.DkimPrivateKey.Valid {
		mimeMsg, err = SignMessage(mimeMsg, domain.Name, domain.DkimSelector, domain.DkimPrivateKey.String)
		if err != nil {
			logger.Warn().Err(err).Msg("DKIM signing failed, sending unsigned")
		}
	}

	// 10. Collect all envelope recipients (To + CC + BCC).
	allRecipients := make([]string, 0, len(toAddrs)+len(ccAddrs)+len(bccAddrs))
	allRecipients = append(allRecipients, toAddrs...)
	allRecipients = append(allRecipients, ccAddrs...)
	allRecipients = append(allRecipients, bccAddrs...)

	// 11. Send via SMTP.
	if err := h.smtp.Send(ctx, email.FromAddress, allRecipients, mimeMsg); err != nil {
		logger.Error().Err(err).Msg("SMTP delivery failed")
		return fmt.Errorf("smtp send: %w", err)
	}

	// 12. Update status to "sent".
	if err := h.queries.UpdateEmailStatus(ctx, sqlcdb.UpdateEmailStatusParams{
		ID:           payload.EmailID,
		OrgID:        row.OrgID,
		Status:       "sent",
		SetSent:      true,
		SetDelivered: false,
	}); err != nil {
		logger.Error().Err(err).Msg("failed to update email status")
		return fmt.Errorf("update email status: %w", err)
	}

	// 13. Record sent event for each recipient.
	for _, recipient := range toAddrs {
		if _, err := h.queries.CreateEmailEvent(ctx, sqlcdb.CreateEmailEventParams{
			EmailID:   payload.EmailID,
			OrgID:     row.OrgID,
			EventType: "sent",
			Recipient: recipient,
			IpAddress: "0.0.0.0",
		}); err != nil {
			logger.Warn().Err(err).Str("recipient", recipient).Msg("failed to record sent event")
		}
	}

	logger.Info().Int("recipients", len(allRecipients)).Msg("email sent successfully")
	return nil
}
