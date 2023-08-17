-- name: GetMessagesByUserID :many
SELECT m.id, m.user_id, m.role, m.content, m.parent_id, m.agent_id, u.user_name as username
FROM (
  SELECT id, user_id, role, content, parent_id, agent_id, created_at
    FROM messages
    WHERE user_id = $1
) m
JOIN users u
ON m.user_id = u.id
ORDER BY m.created_at DESC LIMIT $2;

-- name: CreateMessage :one
INSERT INTO messages (
  user_id, role, content, parent_id, agent_id
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;


-- name: GetUserByJID :one
SELECT * FROM users
WHERE jid = $1 LIMIT 1;


-- name: GetUserByTelegramID :one
SELECT * FROM users
WHERE telegram_id = $1 LIMIT 1;

-- name: GetAgentByName :one
SELECT * FROM agents
WHERE name = $1 LIMIT 1;


-- name: GetAllAgents :many
SELECT * FROM agents;
