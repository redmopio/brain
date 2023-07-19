-- name: GetConversation :one
SELECT * FROM conversations
WHERE id = $1 LIMIT 1;

-- name: ListConversations :many
SELECT * FROM conversations
ORDER BY updated_at;

-- name: CreateConversation :one
INSERT INTO conversations (
  phone_number, jid, context, conversation_buffer, conversation_summary, user_name, tools
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: DeleteConversation :exec
DELETE FROM conversations
WHERE id = $1;


-- name: GetConversationByJid :one
SELECT * FROM conversations
WHERE jid = $1 LIMIT 1;

-- name: GetConversationByPhoneNumber :one
SELECT * FROM conversations
WHERE phone_number = $1 LIMIT 1;

-- name: UpdateConversationContext :one
UPDATE conversations
SET context = $1
WHERE id = $2
RETURNING *;

-- name: UpdateConversationBuffer :one
UPDATE conversations
SET conversation_buffer = $1
WHERE id = $2
RETURNING *;

-- name: UpdateConversationSummary :one
UPDATE conversations
SET conversation_summary = $1
WHERE id = $2
RETURNING *;
