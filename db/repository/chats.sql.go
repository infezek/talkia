// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: chats.sql

package repository

import (
	"context"
)

const createChat = `-- name: CreateChat :exec
INSERT INTO chats (id, user_id, bot_id, created_at)
VALUES (?, ?, ?, NOW())
`

type CreateChatParams struct {
	ID     string
	UserID string
	BotID  string
}

func (q *Queries) CreateChat(ctx context.Context, arg CreateChatParams) error {
	_, err := q.db.ExecContext(ctx, createChat, arg.ID, arg.UserID, arg.BotID)
	return err
}

const deleteChat = `-- name: DeleteChat :exec
DELETE FROM chats
WHERE id = ?
`

func (q *Queries) DeleteChat(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, deleteChat, id)
	return err
}

const getChatByID = `-- name: GetChatByID :one
SELECT id, user_id, bot_id
FROM chats
WHERE id = ?
`

type GetChatByIDRow struct {
	ID     string
	UserID string
	BotID  string
}

func (q *Queries) GetChatByID(ctx context.Context, id string) (GetChatByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getChatByID, id)
	var i GetChatByIDRow
	err := row.Scan(&i.ID, &i.UserID, &i.BotID)
	return i, err
}

const listChatsByUserID = `-- name: ListChatsByUserID :many
SELECT id, user_id, bot_id, created_at, updated_at FROM chats WHERE user_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?
`

type ListChatsByUserIDParams struct {
	UserID string
	Limit  int32
	Offset int32
}

func (q *Queries) ListChatsByUserID(ctx context.Context, arg ListChatsByUserIDParams) ([]Chat, error) {
	rows, err := q.db.QueryContext(ctx, listChatsByUserID, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Chat
	for rows.Next() {
		var i Chat
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.BotID,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const listChatsByUserIDCount = `-- name: ListChatsByUserIDCount :one
SELECT count(*) FROM chats WHERE user_id = ?
`

func (q *Queries) ListChatsByUserIDCount(ctx context.Context, userID string) (int64, error) {
	row := q.db.QueryRowContext(ctx, listChatsByUserIDCount, userID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const updateChat = `-- name: UpdateChat :exec
UPDATE chats
SET user_id = ?, bot_id = ?
WHERE id = ?
`

type UpdateChatParams struct {
	UserID string
	BotID  string
	ID     string
}

func (q *Queries) UpdateChat(ctx context.Context, arg UpdateChatParams) error {
	_, err := q.db.ExecContext(ctx, updateChat, arg.UserID, arg.BotID, arg.ID)
	return err
}
