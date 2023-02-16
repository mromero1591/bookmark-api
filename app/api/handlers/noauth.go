package handlers

import (
	"context"
	"net/http"

	"github.com/mromero1591/bookmark-api/business/web/mid"
	"github.com/mromero1591/web-foundation/auth"
	"github.com/mromero1591/web-foundation/web"
)

type noauthHandler struct {
	auth *auth.Auth
}

// SetupNoAuthHandler: creates the api endpoints for handling all non auth endpoints
func SetupNoAuthHandler(app *web.App, auth *auth.Auth) error {
	//noauth group
	ng := noauthHandler{
		auth: auth,
	}

	//health
	app.Handle(http.MethodGet, "/health", ng.health, mid.Cors("*"))

	return nil
}

func (ng noauthHandler) health(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return web.Respond(ctx, w, "ok", http.StatusOK)
}
