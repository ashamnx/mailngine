-- name: CreateUser :one
INSERT INTO users (email, name, avatar_url, google_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByGoogleID :one
SELECT * FROM users WHERE google_id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: UpdateUserLastLogin :exec
UPDATE users SET last_login_at = NOW() WHERE id = $1;

-- name: UpsertUser :one
INSERT INTO users (email, name, avatar_url, google_id)
VALUES ($1, $2, $3, $4)
ON CONFLICT (google_id) DO UPDATE
SET name = EXCLUDED.name, avatar_url = EXCLUDED.avatar_url, last_login_at = NOW()
RETURNING *;
