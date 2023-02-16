package bookmark

import (
	"time"

	"github.com/google/uuid"
)

type (
	Bookmark struct {
		ID         uuid.UUID `json:"id"`
		Url        string    `json:"url"`
		Name       string    `json:"name"`
		Logo       string    `json:"logo"`
		CategoryID uuid.UUID `json:"category_id"`
		UserID     uuid.UUID `json:"user_id"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
	}

	CreateBookMark struct {
		URL        string    `json:"url" validate:"required,url"`
		CategoryID uuid.UUID `json:"category_id" validate:"required"`
		UserID     uuid.UUID `json:"user_id" validate:"required"`
	}
)
