-- name: CreateMessage :exec
INSERT INTO messages (id, chat_id, user_id, who, message, created_at)
VALUES (?, ?, ?, ?, ?, NOW());

-- name: GetMessageByID :one
SELECT id, chat_id, user_id, who, message
FROM messages
WHERE id = ?;

-- name: DeleteMessage :exec
DELETE FROM messages WHERE id = ?;


-- name: GetMessagesByChatID :many
SELECT id, chat_id, user_id, who, message FROM messages WHERE chat_id = ? ORDER BY id DESC LIMIT 10;