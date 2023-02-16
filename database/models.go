// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.0

package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type UserAccount struct {
	ID        uuid.UUID
	Email     string
	PwdHash   string
	Name      sql.NullString
	CreatedAt time.Time
	UpdatedAt time.Time
}
