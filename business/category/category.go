package category

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mromero1591/bookmark-api/business/sys/validate"
	"github.com/pkg/errors"
)

type StoreAPI interface {
	CreateCategory(ctx context.Context, cat Category) error
	QueryCategoryByUserID(ctx context.Context, userID uuid.UUID) ([]Category, error)
	DeleteCategory(ctx context.Context, id uuid.UUID) error
}

type CategoryService struct {
	store StoreAPI
}

func NewCategoryService(store StoreAPI) CategoryService {
	return CategoryService{store: store}
}

func (s CategoryService) CreateCategory(ctx context.Context, new CreateCategory) (Category, error) {
	if err := validate.Check(new); err != nil {
		return Category{}, errors.Wrap(err, "failed to validate category")
	}

	cat := Category{
		ID:        uuid.New(),
		Name:      strings.TrimSpace(strings.ToLower(new.Name)),
		UserID:    new.UserID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.store.CreateCategory(ctx, cat); err != nil {
		return Category{}, errors.Wrap(err, "failed to create category")
	}

	return cat, nil
}

func (s CategoryService) QueryCategoryByUserID(ctx context.Context, userID uuid.UUID) ([]Category, error) {
	categories, err := s.store.QueryCategoryByUserID(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query categories")
	}

	return categories, nil
}

func (s CategoryService) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	if err := s.store.DeleteCategory(ctx, id); err != nil {
		return errors.Wrap(err, "failed to delete category")
	}

	return nil
}
