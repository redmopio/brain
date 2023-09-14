-- name: GetMessagesByGroupID :many
SELECT m.id, m.user_id, m.role, m.content, m.parent_id, m.agent_id, u.user_name as username
FROM (
  SELECT id, user_id, role, content, parent_id, agent_id, created_at
    FROM messages
    WHERE group_id = $1
) m
JOIN users u
ON m.user_id = u.id
ORDER BY m.created_at DESC LIMIT $2;

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
  group_id, user_id, role, content, parent_id, agent_id
) VALUES (
  $1, $2, $3, $4, $5, $6
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

-- name: CreateConnector :one
INSERT INTO connectors (
  name
) VALUES (
  $1
) RETURNING *;

-- name: CreateGroup :one
INSERT INTO groups (
  real_id, name, description, connector_id
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetGroupsFromUser :many
SELECT g.id , g.name, c.name as connector, g.description
FROM groups g
JOIN users_groups ug  
  ON g.id = ug.group_id
JOIN connectors c
  ON g.connector_id = c.id
WHERE ug.user_id = $1;

-- name: AddUserToGroup :one
INSERT INTO users_groups (
  user_id, group_id
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetGroupByID :one
SELECT * FROM groups
WHERE id = $1 LIMIT 1;

-- name: GetGroupByRealID :one
SELECT * FROM groups
WHERE real_id = $1 LIMIT 1;


-- name: GetConnectorByName :one
SELECT * FROM connectors
WHERE name = $1 LIMIT 1;

-- name: GetUsersFromGroup :many
SELECT u.id, u.user_name, u.phone_number, u.jid, u.telegram_id, u.context
FROM users u
JOIN users_groups ug
  ON u.id = ug.user_id
WHERE ug.group_id = $1;
