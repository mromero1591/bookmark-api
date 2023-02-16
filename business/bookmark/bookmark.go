package bookmark

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mromero1591/bookmark-api/business/metadata"
	"github.com/mromero1591/bookmark-api/business/sys/validate"
	"github.com/pkg/errors"
)

type StoreAPI interface {
	CreateBookmark(ctx context.Context, b Bookmark) error
	QueryBookMarksByUser(ctx context.Context, userID uuid.UUID) ([]Bookmark, error)
}

type BookmarkService struct {
	store StoreAPI
}

func NewBookmarkService(store StoreAPI) BookmarkService {
	return BookmarkService{store: store}
}

func (s BookmarkService) CreateBookmark(ctx context.Context, nb CreateBookMark) (Bookmark, error) {

	if err := validate.Check(nb); err != nil {

	}

	resp, err := http.Get(nb.URL)

	if err != nil {
		return Bookmark{}, errors.Wrap(err, "failed to get url")
	}

	defer resp.Body.Close()

	meta := metadata.Extract(resp.Body)

	bookmark := Bookmark{
		ID:         uuid.New(),
		Url:        nb.URL,
		Name:       meta.Title,
		Logo:       meta.Image,
		CategoryID: nb.CategoryID,
		UserID:     nb.UserID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	return bookmark, s.store.CreateBookmark(ctx, bookmark)
}

func (s BookmarkService) QueryBookMarksByUser(ctx context.Context, userID uuid.UUID) ([]Bookmark, error) {
	return s.store.QueryBookMarksByUser(ctx, userID)
}
