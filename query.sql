-- name: GetUserAccountByEmail :one
SELECT id,email, name,pwd_hash, created_at, updated_at FROM user_accounts WHERE email = $1;

-- name: CreateUserAccount :one
INSERT INTO user_accounts (id, email, name, pwd_hash, created_at, updated_at) 
VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, email, name, pwd_hash, created_at, updated_at;

-- name: CreateCategory :one
INSERT INTO category (id, name, user_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5) RETURNING id, name, user_id, created_at, updated_at;

-- name: CreateBookmark :one
INSERT INTO bookmark (id, url, name, logo, category_id, user_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, url, name, logo, category_id, user_id, created_at, updated_at;

-- name: GetCategoryByUserID :many
SELECT id, name, user_id, created_at, updated_at FROM category WHERE user_id = $1;

-- name: GetBookmarkByUserID :many
SELECT id, url, name, logo, category_id, user_id, created_at, updated_at FROM bookmark WHERE user_id = $1;

-- name: DeleteBookmark :exec
DELETE FROM bookmark WHERE id = $1;

-- name: DeleteCategory :exec
DELETE FROM category WHERE id = $1;