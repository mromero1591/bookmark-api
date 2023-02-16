// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.0
// source: query.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createUserAccount = `-- name: CreateUserAccount :one
INSERT INTO user_accounts (id, email, name, pwd_hash, created_at, updated_at) 
VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, email, name, pwd_hash, created_at, updated_at
`

type CreateUserAccountParams struct {
	ID        uuid.UUID
	Email     string
	Name      sql.NullString
	PwdHash   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateUserAccountRow struct {
	ID        uuid.UUID
	Email     string
	Name      sql.NullString
	PwdHash   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (q *Queries) CreateUserAccount(ctx context.Context, arg CreateUserAccountParams) (CreateUserAccountRow, error) {
	row := q.db.QueryRowContext(ctx, createUserAccount,
		arg.ID,
		arg.Email,
		arg.Name,
		arg.PwdHash,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i CreateUserAccountRow
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Name,
		&i.PwdHash,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserAccountByEmail = `-- name: GetUserAccountByEmail :one
SELECT id,email, name,pwd_hash, created_at, updated_at FROM user_accounts WHERE email = $1
`

type GetUserAccountByEmailRow struct {
	ID        uuid.UUID
	Email     string
	Name      sql.NullString
	PwdHash   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (q *Queries) GetUserAccountByEmail(ctx context.Context, email string) (GetUserAccountByEmailRow, error) {
	row := q.db.QueryRowContext(ctx, getUserAccountByEmail, email)
	var i GetUserAccountByEmailRow
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Name,
		&i.PwdHash,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}