-- name: CreateDomain :one
INSERT INTO domains (org_id, name, dkim_private_key, dkim_selector)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetDomain :one
SELECT * FROM domains WHERE id = $1 AND org_id = $2;

-- name: GetDomainByName :one
SELECT * FROM domains WHERE name = $1 AND org_id = $2;

-- name: ListDomainsByOrg :many
SELECT * FROM domains WHERE org_id = $1 ORDER BY created_at DESC;

-- name: UpdateDomainStatus :exec
UPDATE domains SET status = @status, verified_at = CASE WHEN @set_verified::bool THEN NOW() ELSE verified_at END WHERE id = @id AND org_id = @org_id;

-- name: UpdateDomainSettings :one
UPDATE domains SET open_tracking = $3, click_tracking = $4 WHERE id = $1 AND org_id = $2 RETURNING *;

-- name: UpdateDomainCloudflare :exec
UPDATE domains SET cloudflare_zone_id = $3, cloudflare_api_token_enc = $4 WHERE id = $1 AND org_id = $2;

-- name: DeleteDomain :exec
DELETE FROM domains WHERE id = $1 AND org_id = $2;

-- name: CreateDNSRecord :one
INSERT INTO dns_records (domain_id, record_type, host, value, purpose) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: ListDNSRecordsByDomain :many
SELECT * FROM dns_records WHERE domain_id = $1 ORDER BY purpose;

-- name: UpdateDNSRecordStatus :exec
UPDATE dns_records SET status = @status, verified_at = CASE WHEN @set_verified::bool THEN NOW() ELSE verified_at END WHERE id = @id AND domain_id = @domain_id;

-- name: GetDomainByNameForInbound :one
SELECT * FROM domains WHERE name = $1 AND status = 'verified';

-- name: DeleteDNSRecordsByDomain :exec
DELETE FROM dns_records WHERE domain_id = $1;
