-- name: CreateUser :exec
INSERT INTO users 
(id, username, email, password, platform, location, language, created_at)
VALUES (?, ?, ?, ?, ?, ?, ?, NOW());

-- name: GetUserByID :one
SELECT * FROM users WHERE id = ?;

-- name: UpdateUser :exec
UPDATE users SET 
username = ?,
language = ?
WHERE id = ?;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = ?;

-- name: FindChatsByID :many
SELECT id, user_id, bot_id FROM chats WHERE user_id = ?;

-- name: FindUserByEmail :one
SELECT * FROM users WHERE email = ?;

-- name: FindChatsByUserID :many
SELECT * FROM chats WHERE user_id = ?;

-- name: DeleteCategoriesUserID :exec
DELETE FROM users_categories WHERE user_id = ?;

-- name: AddCategoriesByUserID :exec
INSERT INTO users_categories (id, user_id, category_id) VALUES (?, ?, ?);

-- name: PublishBot :exec
UPDATE bots SET published = TRUE WHERE id = ? AND user_id = ?;

-- name: LikeToBot :exec
INSERT INTO user_like_bot (id, user_id, bot_id, created_at) VALUES (?, ?, ?, NOW());

-- name: FindLikeBotByUserID :many
SELECT * FROM user_like_bot ulb JOIN bots b ON user_like_bot.bot_id = bots.id WHERE ulb.user_id = ?;

-- name: FindLikeBotByBotID :many
SELECT * FROM user_like_bot ulb JOIN users u ON user_like_bot.user_id = users.id WHERE ulb.bot_id = ?;

-- name: FindLikeBotByUserAndBot :one
SELECT * FROM user_like_bot WHERE user_id = ? AND bot_id = ?;


-- name: UpdateAvatarURL :exec
UPDATE users SET avatar_url = ? WHERE id = ?;

-- name: ListBotsByUserID :many
SELECT b.*,
       (SELECT COUNT(*) FROM user_like_bot WHERE bot_id = b.id) AS likes
FROM bots AS b
WHERE b.user_id = ?
ORDER BY b.created_at DESC
LIMIT ? OFFSET ?;

-- name: ListBotsByUserIDCount :one
SELECT count(*) FROM bots WHERE user_id = ?;

-- name: ListCategoriesByUserID :many
SELECT c.* FROM users_categories uc JOIN categories c ON uc.category_id = c.id WHERE uc.user_id = ?;