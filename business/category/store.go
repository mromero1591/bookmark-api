package category

import (
	"context"

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

func (s Store) CreateCategory(ctx context.Context, cat Category) error {

	createCategoryParams := database.CreateCategoryParams{
		ID:        cat.ID,
		Name:      cat.Name,
		UserID:    cat.UserID,
		CreatedAt: cat.CreatedAt,
		UpdatedAt: cat.UpdatedAt,
	}

	_, err := s.db.CreateCategory(ctx, createCategoryParams)
	if err != nil {
		return err
	}

	return nil
}

func (s Store) QueryCategoryByUserID(ctx context.Context, userID uuid.UUID) ([]Category, error) {

	fetchedCategories, err := s.db.GetCategoryByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	categories := make([]Category, len(fetchedCategories))
	for i, c := range fetchedCategories {
		categories[i] = Category{
			ID:        c.ID,
			Name:      c.Name,
			UserID:    c.UserID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		}
	}

	return categories, nil
}

func (s Store) DeleteCategory(ctx context.Context, id uuid.UUID) error {

	if err := s.db.DeleteCategory(ctx, id); err != nil {
		return err
	}

	return nil
}
