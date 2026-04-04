-- Threads

-- name: CreateThread :one
INSERT INTO inbox_threads (org_id, domain_id, subject, participant_addresses, last_message_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetThread :one
SELECT * FROM inbox_threads WHERE id = $1 AND org_id = $2;

-- name: ListThreads :many
SELECT * FROM inbox_threads WHERE org_id = $1 ORDER BY last_message_at DESC LIMIT $2 OFFSET $3;

-- name: CountThreads :one
SELECT COUNT(*) FROM inbox_threads WHERE org_id = $1;

-- name: UpdateThreadLastMessage :exec
UPDATE inbox_threads SET last_message_at = $2, message_count = message_count + 1, participant_addresses = $3 WHERE id = $1;

-- name: DeleteThread :exec
DELETE FROM inbox_threads WHERE id = $1 AND org_id = $2;

-- name: FindThreadBySubjectAndDomain :one
SELECT * FROM inbox_threads
WHERE org_id = $1 AND domain_id = $2 AND subject = $3 AND last_message_at > $4
ORDER BY last_message_at DESC
LIMIT 1;

-- Messages

-- name: CreateInboxMessage :one
INSERT INTO inbox_messages (org_id, domain_id, thread_id, message_id_header, in_reply_to, references_header, from_address, from_name, to_addresses, cc_addresses, subject, text_body_key, html_body_key, snippet, received_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
RETURNING *;

-- name: GetInboxMessage :one
SELECT * FROM inbox_messages WHERE id = $1 AND org_id = $2;

-- name: ListMessagesByThread :many
SELECT * FROM inbox_messages WHERE thread_id = $1 AND org_id = $2 ORDER BY received_at ASC;

-- name: UpdateMessageFlags :exec
UPDATE inbox_messages SET is_read = $3, is_starred = $4, is_archived = $5, is_trashed = $6 WHERE id = $1 AND org_id = $2;

-- name: MarkMessageRead :exec
UPDATE inbox_messages SET is_read = true WHERE id = $1 AND org_id = $2;

-- name: DeleteInboxMessage :exec
DELETE FROM inbox_messages WHERE id = $1 AND org_id = $2;

-- name: SearchMessages :many
SELECT * FROM inbox_messages WHERE org_id = $1 AND subject ILIKE '%' || $2 || '%' ORDER BY received_at DESC LIMIT $3 OFFSET $4;

-- name: FindMessageByMessageID :one
SELECT * FROM inbox_messages WHERE message_id_header = $1 AND org_id = $2;

-- Labels

-- name: CreateLabel :one
INSERT INTO inbox_labels (org_id, name, color) VALUES ($1, $2, $3) RETURNING *;

-- name: ListLabelsByOrg :many
SELECT * FROM inbox_labels WHERE org_id = $1 ORDER BY name;

-- name: UpdateLabel :one
UPDATE inbox_labels SET name = $3, color = $4 WHERE id = $1 AND org_id = $2 RETURNING *;

-- name: DeleteLabel :exec
DELETE FROM inbox_labels WHERE id = $1 AND org_id = $2;

-- name: AddMessageLabel :exec
INSERT INTO inbox_message_labels (message_id, label_id) VALUES ($1, $2) ON CONFLICT DO NOTHING;

-- name: RemoveMessageLabel :exec
DELETE FROM inbox_message_labels WHERE message_id = $1 AND label_id = $2;

-- name: ListMessageLabels :many
SELECT il.* FROM inbox_labels il JOIN inbox_message_labels iml ON il.id = iml.label_id WHERE iml.message_id = $1;
