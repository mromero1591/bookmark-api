package handlers

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/mromero1591/bookmark-api/business/category"
	"github.com/mromero1591/bookmark-api/business/web/mid"
	"github.com/mromero1591/web-foundation/auth"
	"github.com/mromero1591/web-foundation/web"
	"github.com/pkg/errors"
)

type categoryHandler struct {
	categoryService category.CategoryService
}

func SetupCategoryHandler(app *web.App, a *auth.Auth, categoryService category.CategoryService) error {
	ch := categoryHandler{
		categoryService: categoryService,
	}

	app.Handle(http.MethodPost, "/v1/category", ch.createCategory, mid.Cors("*"), mid.Authenticate(a))
	app.Handle(http.MethodGet, "/v1/category", ch.queryCategoriesByUserID, mid.Cors("*"), mid.Authenticate(a))
	return nil
}

func (ch categoryHandler) createCategory(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var nc category.CreateCategory
	if err := web.Decode(r, &nc); err != nil {
		return errors.Wrap(err, "unable to decode payload")
	}

	cat, err := ch.categoryService.CreateCategory(ctx, nc)
	if err != nil {
		return errors.Wrapf(err, "Category: %+v", &cat)
	}

	return web.Respond(ctx, w, cat, http.StatusCreated)
}

func (ch categoryHandler) queryCategoriesByUserID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return errors.Wrapf(err, "Unable to parse user_id: %s", claims.Subject)
	}

	categories, err := ch.categoryService.QueryCategoryByUserID(ctx, userID)
	if err != nil {
		return errors.Wrap(err, "unable to decode payload")
	}

	return web.Respond(ctx, w, categories, http.StatusOK)
}
