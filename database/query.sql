
-- name: GetMessagesByUserID :many
SELECT * FROM messages
WHERE user_id = $1 LIMIT 20;

-- name: CreateMessage :one
INSERT INTO messages (
  id, user_id, role, content, parent_id
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;


-- name: GetUserByJID :one
SELECT * FROM users
WHERE jid = $1 LIMIT 1;


-- name: GetUserByTelegramID :one
SELECT * FROM users
WHERE telegram_id = $1 LIMIT 1;