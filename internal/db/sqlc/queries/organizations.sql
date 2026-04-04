-- name: CreateOrganization :one
INSERT INTO organizations (name, slug, plan, monthly_limit)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetOrganization :one
SELECT * FROM organizations WHERE id = $1;

-- name: GetOrganizationBySlug :one
SELECT * FROM organizations WHERE slug = $1;

-- name: UpdateOrganization :one
UPDATE organizations
SET name = $2, plan = $3, monthly_limit = $4, overage_enabled = $5
WHERE id = $1
RETURNING *;

-- name: UpdateOrganizationName :one
UPDATE organizations SET name = $2 WHERE id = $1 RETURNING *;
