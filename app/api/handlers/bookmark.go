package handlers

import (
	"context"
	"net/http"

	"github.com/mromero1591/bookmark-api/business/bookmark"
	"github.com/mromero1591/bookmark-api/business/web/mid"
	"github.com/mromero1591/web-foundation/auth"
	"github.com/mromero1591/web-foundation/web"
	"github.com/pkg/errors"
)

type bookmarkHandler struct {
	bookmarkService bookmark.BookmarkService
}

func SetupBookmarkHandler(app *web.App, a *auth.Auth, bookmarkService bookmark.BookmarkService) error {
	bh := bookmarkHandler{
		bookmarkService: bookmarkService,
	}

	app.Handle(http.MethodPost, "/v1/bookmark", bh.CreateBookmark, mid.Cors("*"), mid.Authenticate(a))
	return nil
}

func (bh bookmarkHandler) CreateBookmark(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var nb bookmark.CreateBookMark
	if err := web.Decode(r, &nb); err != nil {
		return errors.Wrap(err, "unable to decode payload")
	}

	bk, err := bh.bookmarkService.CreateBookmark(ctx, nb)
	if err != nil {
		return errors.Wrapf(err, "Bookmark: %+v", &bk)
	}

	return web.Respond(ctx, w, bk, http.StatusCreated)

}
