-- name: ListTrends :many
SELECT b.*, COUNT(ulb.id) AS likes
FROM bots b LEFT JOIN user_like_bot ulb ON b.id = ulb.bot_id
WHERE  ulb.created_at >= ? AND ulb.created_at <= ?
GROUP BY  b.id, b.name
ORDER BY  likes DESC LIMIT 10;



-- name: ListCategoriesTrends :many
SELECT c.*, SUM(sub.likes) AS total_likes
FROM categories AS c
LEFT JOIN (
    SELECT b.category_id, b.id AS bot_id, COUNT(ulb.id) AS likes
    FROM bots AS b
    LEFT JOIN user_like_bot AS ulb ON b.id = ulb.bot_id
    WHERE  ulb.created_at >= ? AND ulb.created_at <= ?
    GROUP BY b.category_id, b.id
) AS sub ON c.id = sub.category_id
GROUP BY c.id, c.name
ORDER BY c.name;


-- name: ListBotsOfCategoryTrends :many
SELECT b.*, COUNT(ulb.id) AS likes
FROM bots b LEFT JOIN user_like_bot ulb ON b.id = ulb.bot_id
WHERE  ulb.created_at >= ? AND ulb.created_at <= ? AND b.category_id = ?
GROUP BY  b.id, b.name
ORDER BY  likes DESC LIMIT 10;