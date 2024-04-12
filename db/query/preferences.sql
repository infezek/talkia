-- name: CreatePreference :exec
INSERT INTO preferences (chat_id, user_id, preference_key, preference_value)
VALUES (?, ?, ?, ?);

-- name: GetPreferenceByID :one
SELECT id, chat_id, user_id, preference_key, preference_value
FROM preferences
WHERE id = ?;

-- name: UpdatePreference :exec
UPDATE preferences
SET preference_key = ?, preference_value = ?
WHERE id = ?;

-- name: DeletePreference :exec
DELETE FROM preferences
WHERE id = ?;


-- name: GetPreferencesByChatID :many
SELECT id, chat_id, user_id, preference_key, preference_value
FROM preferences
WHERE chat_id = ?;