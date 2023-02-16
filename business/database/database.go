package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/mromero1591/bookmark-api/app/api/config"
	"github.com/pkg/errors"
)

var (
	ErrNotFound              = errors.New("not found")
	ErrInvalidID             = errors.New("ID is not in its proper form")
	ErrAuthenticationFailure = errors.New("authentication failed")
	ErrForbidden             = errors.New("attempted action is not allowed")
)

func Open(cfg config.Config) (*sql.DB, error) {
	fmt.Printf("\nDB | %+v", cfg.DB)
	db, err := sql.Open("postgres", fmt.Sprintf("user=postgres password=%s host=%s port=5432 dbname=postgres", cfg.DB.Password, cfg.DB.Host))
	if err != nil {
		return nil, err
	}

	return db, nil
}
