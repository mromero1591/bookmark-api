package users

import (
	"time"

	"github.com/google/uuid"
)

type (
	User struct {
		ID        uuid.UUID `json:"id"`
		Email     string    `json:"email"`
		PwdHash   string    `json:"-"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	NewUser struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required`
		Name     string `json:"name"`
	}
)
