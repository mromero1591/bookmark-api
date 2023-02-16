-- name: GetUserAccountByEmail :one
SELECT id,email, name,pwd_hash, created_at, updated_at FROM user_accounts WHERE email = $1;

-- name: CreateUserAccount :one
INSERT INTO user_accounts (id, email, name, pwd_hash, created_at, updated_at) 
VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, email, name, pwd_hash, created_at, updated_at;