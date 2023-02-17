package users

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	dbSetup "github.com/mromero1591/bookmark-api/business/database"
	"github.com/mromero1591/bookmark-api/business/sys/validate"
	"github.com/mromero1591/web-foundation/auth"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type StoreAPI interface {
	CreateUserAccount(ctx context.Context, usr User) error
	GetUserAccountByEmail(ctx context.Context, email string) (User, error)
}

type UserService struct {
	store StoreAPI
}

func NewUserService(store StoreAPI) UserService {
	return UserService{store: store}
}

func (s UserService) CreateUserAccount(ctx context.Context, new NewUser) (User, error) {
	if err := validate.Check(new); err != nil {
		return User{}, errors.Wrap(err, "validating data")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(new.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, errors.Wrap(err, "generating password hash")
	}

	usr := User{
		ID:        uuid.New(),
		Email:     new.Email,
		PwdHash:   string(hash),
		Name:      new.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.store.CreateUserAccount(ctx, usr); err != nil {
		return User{}, errors.Wrap(err, "storing user")
	}

	return usr, nil

}

func (s UserService) GetUserAccountByEmail(ctx context.Context, email string) (User, error) {
	return s.store.GetUserAccountByEmail(ctx, email)
}

// Authenticate finds a user by their email and verifies their password. On
// success it returns a Claims User representing this user. The claims can be
// used to generate a token for future authentication.
func (s UserService) Authenticate(ctx context.Context, email, password string) (auth.Claims, error) {
	usr, err := s.store.GetUserAccountByEmail(ctx, email)
	if err != nil {
		return auth.Claims{}, dbSetup.ErrNotFound
	}

	// Compare the provided password with the saved hash. Use the bcrypt
	// comparison function so it is cryptographically secure.
	if err := bcrypt.CompareHashAndPassword([]byte(usr.PwdHash), []byte(password)); err != nil {
		return auth.Claims{}, dbSetup.ErrAuthenticationFailure
	}

	// If we are this far the request is valid. Create some claims for the user
	// and generate their token.
	claims := auth.Claims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    "bookmark-api",
			Subject:   usr.ID.String(),
			ExpiresAt: time.Now().Add(time.Hour * 8760).Unix(), //one year.
			IssuedAt:  time.Now().UTC().Unix(),
		},
		Roles: []string{"USER"},
	}

	return claims, nil
}
