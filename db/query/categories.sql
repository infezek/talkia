-- name: CreateCategory :exec
INSERT INTO categories 
(id, name, active, created_at)
VALUES (?, ?, ?, ?);

-- name: ListCategories :many
SELECT c.* FROM categories c WHERE c.active = True  ORDER BY created_at DESC LIMIT ? OFFSET ?;

-- name: ListCategoriesCount :one
SELECT count(*) FROM categories c WHERE c.active = True;


-- name: FindCategoriesByUserID :many
SELECT c.* FROM categories c
JOIN users_categories uc ON c.id = uc.category_id
WHERE uc.user_id = ?;


-- name: FindCategoryByID :one
SELECT * FROM categories WHERE id = ?;

-- name: FindCategoryByName :one
SELECT * FROM categories WHERE name = ?;

-- name: UpdateCategory :exec
UPDATE categories
SET name = ?, active = ?
WHERE id = ?;
