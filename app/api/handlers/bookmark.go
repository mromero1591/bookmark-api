package handlers

import (
	"context"
	"net/http"

	"github.com/google/uuid"
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
	app.Handle(http.MethodGet, "/v1/bookmark", bh.QueryBookmarks, mid.Cors("*"), mid.Authenticate(a))
	app.Handle(http.MethodDelete, "/v1/bookmark/:id", bh.DeleteBookmark, mid.Cors("*"), mid.Authenticate(a))
	return nil
}

func (bh bookmarkHandler) CreateBookmark(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var nb bookmark.CreateBookMark
	if err := web.Decode(r, &nb); err != nil {
		return errors.Wrap(err, "unable to decode payload")
	}

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return errors.Wrapf(err, "Unable to parse user_id: %s", claims.Subject)
	}

	nb.UserID = userID

	bk, err := bh.bookmarkService.CreateBookmark(ctx, nb)
	if err != nil {
		return errors.Wrapf(err, "Bookmark: %+v", &bk)
	}

	return web.Respond(ctx, w, bk, http.StatusCreated)

}

func (bh bookmarkHandler) QueryBookmarks(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return errors.Wrapf(err, "Unable to parse user_id: %s", claims.Subject)
	}

	bks, err := bh.bookmarkService.QueryBookMarksByUser(ctx, userID)
	if err != nil {
		return errors.Wrapf(err, "Bookmarks: %+v", &bks)
	}

	return web.Respond(ctx, w, bks, http.StatusOK)
}

func (bh bookmarkHandler) DeleteBookmark(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	id := web.Param(r, "id")

	deleteID, err := uuid.Parse(id)
	if err != nil {
		return errors.Wrapf(err, "Unable to parse id: %s", id)
	}

	if err := bh.bookmarkService.DeleteBookmark(ctx, deleteID); err != nil {
		return errors.Wrapf(err, "Bookmark: %s", deleteID)
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}
