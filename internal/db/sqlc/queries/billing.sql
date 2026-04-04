-- name: GetCurrentMonthUsage :one
SELECT COALESCE(SUM(emails_sent), 0)::int as total_sent,
       COALESCE(SUM(emails_received), 0)::int as total_received,
       COALESCE(SUM(api_calls), 0)::int as total_api_calls
FROM usage_daily
WHERE org_id = $1
  AND date >= date_trunc('month', CURRENT_DATE)::date
  AND date <= CURRENT_DATE;

-- name: GetMonthlyUsageHistory :many
SELECT * FROM usage_monthly WHERE org_id = $1 ORDER BY month DESC LIMIT $2;

-- name: UpsertUsageMonthly :exec
INSERT INTO usage_monthly (org_id, month, total_sent, total_received, total_api_calls, plan_limit, overage_count)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (org_id, month) DO UPDATE SET
  total_sent = EXCLUDED.total_sent,
  total_received = EXCLUDED.total_received,
  total_api_calls = EXCLUDED.total_api_calls,
  overage_count = EXCLUDED.overage_count;
