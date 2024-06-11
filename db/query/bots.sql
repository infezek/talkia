-- name: CreateBot :exec
INSERT INTO bots (id, user_id, category_id, name, personality, description, avatar_url,background_url, location, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW());

-- name: GetBotByID :one
SELECT b.*, 
       (SELECT COUNT(*) FROM user_like_bot ulb WHERE ulb.bot_id = b.id) as likes 
FROM bots b 
WHERE b.id = ?;


-- name: UpdateBot :exec
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
WHERE id = ?;

-- name: DeleteBot :exec
DELETE FROM bots
WHERE id = ?;

-- name: FindBotsByUserID :many
SELECT b.*, (SELECT COUNT(*) FROM user_like_bot ulb WHERE ulb.bot_id = b.id) as likes FROM bots b WHERE b.user_id = ?;

-- name: FindBotsByName :many
SELECT b.*, (SELECT COUNT(*) FROM user_like_bot ulb WHERE ulb.bot_id = b.id) as likes FROM bots b WHERE name LIKE ? ORDER BY likes DESC LIMIT ? OFFSET ?;