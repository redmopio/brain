// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: query.sql

package models

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createMessage = `-- name: CreateMessage :one
INSERT INTO messages (
  user_id, role, content, parent_id, agent_id
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING id, created_at, updated_at, user_id, role, content, parent_id, agent_id
`

type CreateMessageParams struct {
	UserID   uuid.NullUUID
	Role     sql.NullString
	Content  sql.NullString
	ParentID uuid.NullUUID
	AgentID  uuid.NullUUID
}

func (q *Queries) CreateMessage(ctx context.Context, arg CreateMessageParams) (Message, error) {
	row := q.db.QueryRowContext(ctx, createMessage,
		arg.UserID,
		arg.Role,
		arg.Content,
		arg.ParentID,
		arg.AgentID,
	)
	var i Message
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.Role,
		&i.Content,
		&i.ParentID,
		&i.AgentID,
	)
	return i, err
}

const getAgentByName = `-- name: GetAgentByName :one
SELECT id, created_at, updated_at, name, constitution FROM agents
WHERE name = $1 LIMIT 1
`

func (q *Queries) GetAgentByName(ctx context.Context, name string) (Agent, error) {
	row := q.db.QueryRowContext(ctx, getAgentByName, name)
	var i Agent
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Constitution,
	)
	return i, err
}

const getMessagesByUserID = `-- name: GetMessagesByUserID :many
SELECT m.id, m.user_id, m.role, m.content, m.parent_id, m.agent_id, u.user_name as username
FROM (
  SELECT id, user_id, role, content, parent_id, agent_id, created_at
    FROM messages
    WHERE user_id = $1
) m
JOIN users u
ON m.user_id = u.id
ORDER BY m.created_at DESC LIMIT $2
`

type GetMessagesByUserIDParams struct {
	UserID uuid.NullUUID
	Limit  int32
}

type GetMessagesByUserIDRow struct {
	ID       uuid.UUID
	UserID   uuid.NullUUID
	Role     sql.NullString
	Content  sql.NullString
	ParentID uuid.NullUUID
	AgentID  uuid.NullUUID
	Username sql.NullString
}

func (q *Queries) GetMessagesByUserID(ctx context.Context, arg GetMessagesByUserIDParams) ([]GetMessagesByUserIDRow, error) {
	rows, err := q.db.QueryContext(ctx, getMessagesByUserID, arg.UserID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetMessagesByUserIDRow
	for rows.Next() {
		var i GetMessagesByUserIDRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Role,
			&i.Content,
			&i.ParentID,
			&i.AgentID,
			&i.Username,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserByJID = `-- name: GetUserByJID :one
SELECT id, created_at, updated_at, phone_number, jid, telegram_id, context, user_name FROM users
WHERE jid = $1 LIMIT 1
`

func (q *Queries) GetUserByJID(ctx context.Context, jid sql.NullString) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByJID, jid)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.PhoneNumber,
		&i.Jid,
		&i.TelegramID,
		&i.Context,
		&i.UserName,
	)
	return i, err
}

const getUserByTelegramID = `-- name: GetUserByTelegramID :one
SELECT id, created_at, updated_at, phone_number, jid, telegram_id, context, user_name FROM users
WHERE telegram_id = $1 LIMIT 1
`

func (q *Queries) GetUserByTelegramID(ctx context.Context, telegramID sql.NullString) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByTelegramID, telegramID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.PhoneNumber,
		&i.Jid,
		&i.TelegramID,
		&i.Context,
		&i.UserName,
	)
	return i, err
}
