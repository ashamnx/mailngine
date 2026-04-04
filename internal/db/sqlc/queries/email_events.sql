-- name: CreateEmailEvent :one
INSERT INTO email_events (email_id, org_id, event_type, recipient, metadata, ip_address, user_agent)
VALUES ($1, $2, $3, $4, $5, sqlc.arg(ip_address)::inet, $6)
RETURNING *;

-- name: ListEmailEvents :many
SELECT * FROM email_events WHERE email_id = $1 ORDER BY occurred_at ASC;

-- name: ListEmailEventsByOrg :many
SELECT * FROM email_events WHERE org_id = $1 ORDER BY occurred_at DESC LIMIT $2 OFFSET $3;

-- name: CountEmailEventsByOrg :one
SELECT COUNT(*) FROM email_events WHERE org_id = $1;
