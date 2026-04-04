-- name: CreateTemplate :one
INSERT INTO templates (org_id, name, subject, html_body, text_body, variables)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetTemplate :one
SELECT * FROM templates WHERE id = $1 AND org_id = $2;

-- name: GetTemplateByName :one
SELECT * FROM templates WHERE org_id = $1 AND name = $2;

-- name: ListTemplatesByOrg :many
SELECT * FROM templates WHERE org_id = $1 ORDER BY name;

-- name: UpdateTemplate :one
UPDATE templates
SET name = $3, subject = $4, html_body = $5, text_body = $6, variables = $7
WHERE id = $1 AND org_id = $2
RETURNING *;

-- name: DeleteTemplate :exec
DELETE FROM templates WHERE id = $1 AND org_id = $2;
