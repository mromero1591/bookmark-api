package users

import (
	"context"
	"database/sql"

	dbSetup "github.com/mromero1591/bookmark-api/business/database"
	"github.com/mromero1591/bookmark-api/database"
)

var (
	sqlcNotFoundErr = "sql: no rows in result set"
)

type Store struct {
	db *database.Queries
}

func NewStore(db *database.Queries) Store {
	if db == nil {
		panic("db is nil")
	}
	return Store{db: db}
}

func (s Store) CreateUserAccount(ctx context.Context, usr User) error {

	createUserParams := database.CreateUserAccountParams{
		ID:        usr.ID,
		Email:     usr.Email,
		Name:      sql.NullString{String: usr.Name, Valid: true},
		PwdHash:   usr.PwdHash,
		CreatedAt: usr.CreatedAt,
		UpdatedAt: usr.UpdatedAt,
	}

	_, err := s.db.CreateUserAccount(ctx, createUserParams)
	if err != nil {
		return err
	}

	return nil
}

func (s Store) GetUserAccountByEmail(ctx context.Context, email string) (User, error) {

	fetchedUser, err := s.db.GetUserAccountByEmail(ctx, email)
	if err != nil {
		if err.Error() == sqlcNotFoundErr {
			return User{}, dbSetup.ErrNotFound
		}
		return User{}, err
	}

	user := User{
		ID:        fetchedUser.ID,
		Email:     fetchedUser.Email,
		PwdHash:   fetchedUser.PwdHash,
		Name:      fetchedUser.Name.String,
		CreatedAt: fetchedUser.CreatedAt,
		UpdatedAt: fetchedUser.UpdatedAt,
	}

	return user, nil
}
