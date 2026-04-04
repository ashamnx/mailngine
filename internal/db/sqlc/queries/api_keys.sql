-- name: CreateAPIKey :one
INSERT INTO api_keys (org_id, domain_id, name, prefix, key_hash, permission, expires_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetAPIKeyByHash :one
SELECT * FROM api_keys
WHERE key_hash = $1 AND revoked_at IS NULL AND (expires_at IS NULL OR expires_at > NOW());

-- name: GetAPIKey :one
SELECT * FROM api_keys WHERE id = $1 AND org_id = $2;

-- name: ListAPIKeysByOrg :many
SELECT * FROM api_keys
WHERE org_id = $1 AND revoked_at IS NULL
ORDER BY created_at DESC;

-- name: RevokeAPIKey :exec
UPDATE api_keys SET revoked_at = NOW() WHERE id = $1 AND org_id = $2;

-- name: UpdateAPIKeyLastUsed :exec
UPDATE api_keys SET last_used_at = NOW() WHERE id = $1;
