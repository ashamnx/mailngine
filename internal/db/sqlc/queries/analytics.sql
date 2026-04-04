-- name: UpsertUsageDaily :exec
INSERT INTO usage_daily (org_id, date, emails_sent, emails_delivered, emails_bounced, emails_received, api_calls)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (org_id, date) DO UPDATE SET
  emails_sent = usage_daily.emails_sent + EXCLUDED.emails_sent,
  emails_delivered = usage_daily.emails_delivered + EXCLUDED.emails_delivered,
  emails_bounced = usage_daily.emails_bounced + EXCLUDED.emails_bounced,
  emails_received = usage_daily.emails_received + EXCLUDED.emails_received,
  api_calls = usage_daily.api_calls + EXCLUDED.api_calls;

-- name: GetUsageDaily :many
SELECT * FROM usage_daily WHERE org_id = $1 AND date >= $2 AND date <= $3 ORDER BY date;

-- name: GetUsageSummary :one
SELECT
  COALESCE(SUM(emails_sent), 0)::int as total_sent,
  COALESCE(SUM(emails_delivered), 0)::int as total_delivered,
  COALESCE(SUM(emails_bounced), 0)::int as total_bounced,
  COALESCE(SUM(emails_received), 0)::int as total_received
FROM usage_daily WHERE org_id = $1 AND date >= $2 AND date <= $3;

-- name: CountEmailEventsByType :many
SELECT event_type, COUNT(*) as count FROM email_events WHERE org_id = $1 AND occurred_at >= $2 AND occurred_at <= $3 GROUP BY event_type;
