package bookmark

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/mromero1591/bookmark-api/database"
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

func (s Store) CreateBookmark(ctx context.Context, b Bookmark) error {
	createBookMarkParams := database.CreateBookmarkParams{
		ID:   b.ID,
		Url:  b.Url,
		Name: b.Name,
		Logo: sql.NullString{
			String: b.Logo,
		},
		CategoryID: b.CategoryID,
		UserID:     b.UserID,
		CreatedAt:  b.CreatedAt,
		UpdatedAt:  b.UpdatedAt,
	}

	_, err := s.db.CreateBookmark(ctx, createBookMarkParams)
	if err != nil {
		return err
	}

	return nil
}

func (s Store) QueryBookMarksByUser(ctx context.Context, userID uuid.UUID) ([]Bookmark, error) {
	fetchedBookmarks, err := s.db.GetBookmarkByUserID(ctx, userID)

	if err != nil {
		return nil, err
	}

	bookmarks := make([]Bookmark, len(fetchedBookmarks))
	for i, b := range fetchedBookmarks {
		bookmarks[i] = Bookmark{
			ID:         b.ID,
			Url:        b.Url,
			Name:       b.Name,
			Logo:       b.Logo.String,
			CategoryID: b.CategoryID,
			UserID:     b.UserID,
			CreatedAt:  b.CreatedAt,
			UpdatedAt:  b.UpdatedAt,
		}

		return bookmarks, nil

	}

	return bookmarks, nil
}
