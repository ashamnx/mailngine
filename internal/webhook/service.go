package webhook

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/url"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	"github.com/mailngine/mailngine/internal/db/sqlcdb"
	"github.com/mailngine/mailngine/internal/queue"
)

// ErrUnsafeURL is returned when a webhook URL targets a private/internal IP address.
var ErrUnsafeURL = fmt.Errorf("webhook URL must not target private or internal IP addresses")

// validateWebhookURL checks that the URL uses https or http and does not resolve
// to a private, loopback, or link-local IP address (SSRF protection).
func validateWebhookURL(rawURL string) error {
	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("webhook URL must use http or https scheme")
	}

	host := u.Hostname()
	ips, err := net.LookupIP(host)
	if err != nil {
		return fmt.Errorf("cannot resolve webhook host %q: %w", host, err)
	}

	for _, ip := range ips {
		if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsUnspecified() {
			return ErrUnsafeURL
		}
	}

	return nil
}

type Service struct {
	db      *pgxpool.Pool
	queries *sqlcdb.Queries
	queue   *asynq.Client
	logger  zerolog.Logger
}

func NewService(db *pgxpool.Pool, queue *asynq.Client, logger zerolog.Logger) *Service {
	return &Service{
		db:      db,
		queries: sqlcdb.New(db),
		queue:   queue,
		logger:  logger,
	}
}

func (s *Service) Create(ctx context.Context, orgID uuid.UUID, webhookURL string, events []string) (*sqlcdb.Webhook, error) {
	if err := validateWebhookURL(webhookURL); err != nil {
		return nil, fmt.Errorf("validate webhook URL: %w", err)
	}

	secret, err := generateSecret()
	if err != nil {
		return nil, fmt.Errorf("generate webhook secret: %w", err)
	}

	eventsJSON, err := json.Marshal(events)
	if err != nil {
		return nil, fmt.Errorf("marshal events: %w", err)
	}

	webhook, err := s.queries.CreateWebhook(ctx, sqlcdb.CreateWebhookParams{
		OrgID:    orgID,
		Url:      webhookURL,
		Events:   eventsJSON,
		Secret:   secret,
		IsActive: true,
	})
	if err != nil {
		return nil, fmt.Errorf("create webhook: %w", err)
	}
	return &webhook, nil
}

func (s *Service) List(ctx context.Context, orgID uuid.UUID) ([]sqlcdb.Webhook, error) {
	return s.queries.ListWebhooksByOrg(ctx, orgID)
}

func (s *Service) Get(ctx context.Context, orgID, webhookID uuid.UUID) (*sqlcdb.Webhook, error) {
	wh, err := s.queries.GetWebhook(ctx, sqlcdb.GetWebhookParams{ID: webhookID, OrgID: orgID})
	if err != nil {
		return nil, fmt.Errorf("get webhook: %w", err)
	}
	return &wh, nil
}

func (s *Service) Update(ctx context.Context, orgID, webhookID uuid.UUID, webhookURL string, events []string, isActive bool) (*sqlcdb.Webhook, error) {
	if err := validateWebhookURL(webhookURL); err != nil {
		return nil, fmt.Errorf("validate webhook URL: %w", err)
	}

	eventsJSON, err := json.Marshal(events)
	if err != nil {
		return nil, fmt.Errorf("marshal events: %w", err)
	}

	wh, err := s.queries.UpdateWebhook(ctx, sqlcdb.UpdateWebhookParams{
		ID:       webhookID,
		OrgID:    orgID,
		Url:      webhookURL,
		Events:   eventsJSON,
		IsActive: isActive,
	})
	if err != nil {
		return nil, fmt.Errorf("update webhook: %w", err)
	}
	return &wh, nil
}

func (s *Service) Delete(ctx context.Context, orgID, webhookID uuid.UUID) error {
	return s.queries.DeleteWebhook(ctx, sqlcdb.DeleteWebhookParams{ID: webhookID, OrgID: orgID})
}

func (s *Service) ListDeliveries(ctx context.Context, webhookID uuid.UUID, page, perPage int) ([]sqlcdb.WebhookDelivery, error) {
	offset := (page - 1) * perPage
	return s.queries.ListWebhookDeliveries(ctx, sqlcdb.ListWebhookDeliveriesParams{
		WebhookID: webhookID,
		Limit:     int32(perPage),
		Offset:    int32(offset),
	})
}

func (s *Service) Dispatch(ctx context.Context, orgID uuid.UUID, eventType string, payload any) error {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal webhook payload: %w", err)
	}

	eventFilter, _ := json.Marshal([]string{eventType})
	webhooks, err := s.queries.ListWebhooksByOrgAndEvent(ctx, sqlcdb.ListWebhooksByOrgAndEventParams{
		OrgID:  orgID,
		Events: eventFilter,
	})
	if err != nil {
		return fmt.Errorf("list webhooks for event: %w", err)
	}

	for _, wh := range webhooks {
		delivery, err := s.queries.CreateWebhookDelivery(ctx, sqlcdb.CreateWebhookDeliveryParams{
			WebhookID: wh.ID,
			EventType: eventType,
			Payload:   payloadJSON,
			Status:    "pending",
		})
		if err != nil {
			s.logger.Error().Err(err).Str("webhook_id", wh.ID.String()).Msg("failed to create webhook delivery")
			continue
		}

		taskPayload, _ := json.Marshal(map[string]string{"delivery_id": delivery.ID.String()})
		task := asynq.NewTask(queue.TaskDispatchWebhook, taskPayload)
		if _, err := s.queue.Enqueue(task); err != nil {
			s.logger.Error().Err(err).Str("delivery_id", delivery.ID.String()).Msg("failed to enqueue webhook delivery")
		}
	}

	return nil
}

func generateSecret() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
