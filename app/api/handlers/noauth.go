package handlers

import (
	"context"
	"net/http"

	dbSetup "github.com/mromero1591/bookmark-api/business/database"
	"github.com/mromero1591/bookmark-api/business/sys/validate"
	"github.com/mromero1591/bookmark-api/business/users"
	"github.com/mromero1591/bookmark-api/business/web/mid"
	"github.com/mromero1591/web-foundation/auth"
	"github.com/mromero1591/web-foundation/web"
	"github.com/pkg/errors"
)

type noauthHandler struct {
	auth        *auth.Auth
	userService users.UserService
}

// SetupNoAuthHandler: creates the api endpoints for handling all non auth endpoints
func SetupNoAuthHandler(app *web.App, auth *auth.Auth, userService users.UserService) error {
	//noauth group
	ng := noauthHandler{
		auth:        auth,
		userService: userService,
	}

	//health
	app.Handle(http.MethodGet, "/health", ng.health, mid.Cors("*"))
	app.Handle(http.MethodPost, "/v1/signup", ng.signUp, mid.Cors("*"))

	return nil
}

func (ng noauthHandler) health(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return web.Respond(ctx, w, "ok", http.StatusOK)
}

// UserAdd adds new users into the database.
func (ng noauthHandler) signUp(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	var nu users.NewUser
	if err := web.Decode(r, &nu); err != nil {
		return errors.Wrap(err, "unable to decode payload")
	}

	usr, err := ng.userService.CreateUserAccount(ctx, nu)
	if err != nil {
		return errors.Wrapf(err, "User: %+v", &usr)
	}

	//once user is created than we will authenticated.
	claims, err := ng.userService.Authenticate(ctx, nu.Email, nu.Password)
	if err != nil {
		switch errors.Cause(err) {
		case dbSetup.ErrNotFound:
			return validate.NewRequestError(err, http.StatusNotFound)
		case dbSetup.ErrAuthenticationFailure:
			return validate.NewRequestError(err, http.StatusUnauthorized)
		default:
			return errors.Wrap(err, "authenticating")
		}
	}

	var tkn struct {
		Token string `json:"token"`
	}
	tkn.Token, err = ng.auth.GenerateToken(claims)
	if err != nil {
		return errors.Wrap(err, "generating token")
	}

	response := struct {
		Token string     `json:"token"`
		User  users.User `json:"user"`
	}{
		Token: tkn.Token,
		User:  usr,
	}

	return web.Respond(ctx, w, response, http.StatusCreated)
}
