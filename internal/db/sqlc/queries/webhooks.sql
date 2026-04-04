-- name: CreateWebhook :one
INSERT INTO webhooks (org_id, url, events, secret, is_active)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetWebhook :one
SELECT * FROM webhooks WHERE id = $1 AND org_id = $2;

-- name: ListWebhooksByOrg :many
SELECT * FROM webhooks WHERE org_id = $1 ORDER BY created_at DESC;

-- name: UpdateWebhook :one
UPDATE webhooks SET url = $3, events = $4, is_active = $5
WHERE id = $1 AND org_id = $2
RETURNING *;

-- name: DeleteWebhook :exec
DELETE FROM webhooks WHERE id = $1 AND org_id = $2;

-- name: ListWebhooksByOrgAndEvent :many
SELECT * FROM webhooks WHERE org_id = $1 AND is_active = true AND events @> $2;

-- name: CreateWebhookDelivery :one
INSERT INTO webhook_deliveries (webhook_id, event_type, payload, status)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateWebhookDelivery :exec
UPDATE webhook_deliveries
SET response_status = $2,
    response_body = $3,
    attempt = $4,
    status = $5,
    delivered_at = CASE WHEN $5 = 'delivered' THEN NOW() ELSE delivered_at END,
    next_retry_at = $6
WHERE id = $1;

-- name: ListWebhookDeliveries :many
SELECT * FROM webhook_deliveries WHERE webhook_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: GetWebhookDelivery :one
SELECT * FROM webhook_deliveries WHERE id = $1;

-- name: GetWebhookByDeliveryID :one
SELECT w.* FROM webhooks w
JOIN webhook_deliveries wd ON wd.webhook_id = w.id
WHERE wd.id = $1;
