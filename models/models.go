// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Agent struct {
	ID           uuid.UUID
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Name         string
	Constitution string
}

type Connector struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
}

type Group struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        sql.NullString
	Description sql.NullString
	RealID      sql.NullString
	ConnectorID uuid.NullUUID
}

type Message struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.NullUUID
	GroupID   uuid.NullUUID
	Role      sql.NullString
	Content   sql.NullString
	ParentID  uuid.NullUUID
	AgentID   uuid.NullUUID
}

type User struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	PhoneNumber string
	Jid         sql.NullString
	TelegramID  sql.NullString
	Context     sql.NullString
	UserName    sql.NullString
}

type UsersGroup struct {
	UserID    uuid.UUID
	GroupID   uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}
