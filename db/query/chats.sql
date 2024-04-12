-- name: CreateChat :exec
INSERT INTO chats (id, user_id, bot_id, created_at)
VALUES (?, ?, ?, NOW());

-- name: GetChatByID :one
SELECT id, user_id, bot_id
FROM chats
WHERE id = ?;

-- name: UpdateChat :exec
UPDATE chats
SET user_id = ?, bot_id = ?
WHERE id = ?;

-- name: DeleteChat :exec
DELETE FROM chats
WHERE id = ?;

-- name: ListChatsByUserID :many
SELECT * FROM chats WHERE user_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?;

-- name: ListChatsByUserIDCount :one
SELECT count(*) FROM chats WHERE user_id = ?;