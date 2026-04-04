-- name: CreateAuditLog :one
INSERT INTO audit_logs (org_id, user_id, api_key_id, action, resource_type, resource_id, metadata, ip_address, user_agent)
VALUES ($1, $2, $3, $4, $5, $6, $7, sqlc.arg(ip_address)::inet, $8)
RETURNING *;

-- name: ListAuditLogs :many
SELECT * FROM audit_logs WHERE org_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: CountAuditLogs :one
SELECT COUNT(*) FROM audit_logs WHERE org_id = $1;

-- name: GetAuditLog :one
SELECT * FROM audit_logs WHERE id = $1 AND org_id = $2;
