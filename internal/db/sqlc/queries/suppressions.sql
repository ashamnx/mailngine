-- name: CreateSuppression :one
INSERT INTO suppressions (org_id, email, reason, metadata)
VALUES ($1, $2, $3, $4)
ON CONFLICT (org_id, email) DO NOTHING
RETURNING *;

-- name: GetSuppression :one
SELECT * FROM suppressions WHERE org_id = $1 AND email = $2;

-- name: ListSuppressionsByOrg :many
SELECT * FROM suppressions WHERE org_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: CountSuppressionsByOrg :one
SELECT COUNT(*) FROM suppressions WHERE org_id = $1;

-- name: DeleteSuppression :exec
DELETE FROM suppressions WHERE id = $1 AND org_id = $2;

-- name: CheckSuppressed :one
SELECT EXISTS(SELECT 1 FROM suppressions WHERE org_id = $1 AND email = $2) AS is_suppressed;
