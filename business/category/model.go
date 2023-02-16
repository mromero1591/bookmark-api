package category

import (
	"time"

	"github.com/google/uuid"
)

type (
	Category struct {
		ID        uuid.UUID `json:"id"`
		Name      string    `json:"name"`
		UserID    uuid.UUID `json:"user_id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	CreateCategory struct {
		Name   string    `json:"name"`
		UserID uuid.UUID `json:"user_id" validate:"required"`
	}
)
