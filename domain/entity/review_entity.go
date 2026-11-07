package entity

import (
	"time"

	"github.com/google/uuid"

	"github.com/rescuebyte/rescuebyte-be/domain/enums"
)

type Review struct {
	ID         uuid.UUID         `db:"id" json:"id"`
	UserID     uuid.UUID         `db:"user_id" json:"user_id"`
	StoreID    *uuid.UUID        `db:"store_id" json:"store_id,omitempty"`
	ProductID  *uuid.UUID        `db:"product_id" json:"product_id,omitempty"`
	Content    string            `db:"content" json:"content"`
	Rating     int               `db:"rating" json:"rating"`
	CreatedAt  time.Time         `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time         `db:"updated_at" json:"updated_at"`
	AuthorName string            `db:"author_name" json:"author_name"`
	AvatarURL  *string           `db:"avatar_url" json:"avatar_url"`
	AuthorRole enums.UserRoleIdx `db:"author_role" json:"author_role"`
}
