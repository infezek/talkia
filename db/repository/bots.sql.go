// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: bots.sql

package repository

import (
	"context"
	"time"
)

const createBot = `-- name: CreateBot :exec
INSERT INTO bots (id, user_id, category_id, name, personality, description, avatar_url,background_url, location, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
`

type CreateBotParams struct {
	ID            string
	UserID        string
	CategoryID    string
	Name          string
	Personality   string
	Description   string
	AvatarUrl     string
	BackgroundUrl string
	Location      string
}

func (q *Queries) CreateBot(ctx context.Context, arg CreateBotParams) error {
	_, err := q.db.ExecContext(ctx, createBot,
		arg.ID,
		arg.UserID,
		arg.CategoryID,
		arg.Name,
		arg.Personality,
		arg.Description,
		arg.AvatarUrl,
		arg.BackgroundUrl,
		arg.Location,
	)
	return err
}

const deleteBot = `-- name: DeleteBot :exec
DELETE FROM bots
WHERE id = ?
`

func (q *Queries) DeleteBot(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, deleteBot, id)
	return err
}

const findBotsByName = `-- name: FindBotsByName :many
SELECT b.id, b.user_id, b.category_id, b.name, b.personality, b.description, b.avatar_url, b.background_url, b.location, b.published, b.active, b.created_at, b.updated_at, (SELECT COUNT(*) FROM user_like_bot ulb WHERE ulb.bot_id = b.id) as likes FROM bots b WHERE name LIKE ? ORDER BY likes DESC LIMIT ? OFFSET ?
`

type FindBotsByNameParams struct {
	Name   string
	Limit  int32
	Offset int32
}

type FindBotsByNameRow struct {
	ID            string
	UserID        string
	CategoryID    string
	Name          string
	Personality   string
	Description   string
	AvatarUrl     string
	BackgroundUrl string
	Location      string
	Published     bool
	Active        bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Likes         int64
}

func (q *Queries) FindBotsByName(ctx context.Context, arg FindBotsByNameParams) ([]FindBotsByNameRow, error) {
	rows, err := q.db.QueryContext(ctx, findBotsByName, arg.Name, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindBotsByNameRow
	for rows.Next() {
		var i FindBotsByNameRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.CategoryID,
			&i.Name,
			&i.Personality,
			&i.Description,
			&i.AvatarUrl,
			&i.BackgroundUrl,
			&i.Location,
			&i.Published,
			&i.Active,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Likes,
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

const findBotsByUserID = `-- name: FindBotsByUserID :many
SELECT b.id, b.user_id, b.category_id, b.name, b.personality, b.description, b.avatar_url, b.background_url, b.location, b.published, b.active, b.created_at, b.updated_at, (SELECT COUNT(*) FROM user_like_bot ulb WHERE ulb.bot_id = b.id) as likes FROM bots b WHERE b.user_id = ?
`

type FindBotsByUserIDRow struct {
	ID            string
	UserID        string
	CategoryID    string
	Name          string
	Personality   string
	Description   string
	AvatarUrl     string
	BackgroundUrl string
	Location      string
	Published     bool
	Active        bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Likes         int64
}

func (q *Queries) FindBotsByUserID(ctx context.Context, userID string) ([]FindBotsByUserIDRow, error) {
	rows, err := q.db.QueryContext(ctx, findBotsByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindBotsByUserIDRow
	for rows.Next() {
		var i FindBotsByUserIDRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.CategoryID,
			&i.Name,
			&i.Personality,
			&i.Description,
			&i.AvatarUrl,
			&i.BackgroundUrl,
			&i.Location,
			&i.Published,
			&i.Active,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Likes,
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

const getBotByID = `-- name: GetBotByID :one
SELECT b.id, b.user_id, b.category_id, b.name, b.personality, b.description, b.avatar_url, b.background_url, b.location, b.published, b.active, b.created_at, b.updated_at, 
       (SELECT COUNT(*) FROM user_like_bot ulb WHERE ulb.bot_id = b.id) as likes 
FROM bots b 
WHERE b.id = ?
`

type GetBotByIDRow struct {
	ID            string
	UserID        string
	CategoryID    string
	Name          string
	Personality   string
	Description   string
	AvatarUrl     string
	BackgroundUrl string
	Location      string
	Published     bool
	Active        bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Likes         int64
}

func (q *Queries) GetBotByID(ctx context.Context, id string) (GetBotByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getBotByID, id)
	var i GetBotByIDRow
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.CategoryID,
		&i.Name,
		&i.Personality,
		&i.Description,
		&i.AvatarUrl,
		&i.BackgroundUrl,
		&i.Location,
		&i.Published,
		&i.Active,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Likes,
	)
	return i, err
}

const updateBot = `-- name: UpdateBot :exec
UPDATE bots
SET name = ?, 
personality = ?, 
description = ?,
background_url = ?,
avatar_url = ?, 
location = ?,
category_id = ?,
active = ?,
updated_at = NOW()
WHERE id = ?
`

type UpdateBotParams struct {
	Name          string
	Personality   string
	Description   string
	BackgroundUrl string
	AvatarUrl     string
	Location      string
	CategoryID    string
	Active        bool
	ID            string
}

func (q *Queries) UpdateBot(ctx context.Context, arg UpdateBotParams) error {
	_, err := q.db.ExecContext(ctx, updateBot,
		arg.Name,
		arg.Personality,
		arg.Description,
		arg.BackgroundUrl,
		arg.AvatarUrl,
		arg.Location,
		arg.CategoryID,
		arg.Active,
		arg.ID,
	)
	return err
}
