-- name: CreateEmail :one
INSERT INTO emails (org_id, domain_id, api_key_id, idempotency_key, from_address, from_name, to_addresses, cc_addresses, bcc_addresses, reply_to, subject, text_body_key, html_body_key, headers, tags, template_id, template_data, status, scheduled_at, message_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
RETURNING *;

-- name: GetEmail :one
SELECT * FROM emails WHERE id = $1 AND org_id = $2;

-- name: GetEmailByIdempotencyKey :one
SELECT * FROM emails WHERE org_id = $1 AND idempotency_key = $2;

-- name: ListEmailsByOrg :many
SELECT * FROM emails WHERE org_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: UpdateEmailStatus :exec
UPDATE emails SET status = @status, sent_at = CASE WHEN @set_sent::bool THEN NOW() ELSE sent_at END, delivered_at = CASE WHEN @set_delivered::bool THEN NOW() ELSE delivered_at END WHERE id = @id AND org_id = @org_id;

-- name: CountEmailsByOrg :one
SELECT COUNT(*) FROM emails WHERE org_id = $1;

-- name: GetEmailOrgID :one
SELECT id, org_id FROM emails WHERE id = $1;

-- name: GetVerifiedDomainByName :one
SELECT * FROM domains WHERE name = $1 AND org_id = $2 AND status = 'verified';
