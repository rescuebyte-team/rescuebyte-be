package entity

import (
	"time"

	"github.com/google/uuid"
)

type ProductImage struct {
	ID           uuid.UUID `db:"id" json:"id"`
	ProductID    uuid.UUID `db:"product_id" json:"product_id"`
	ImageURL     string    `db:"image_url" json:"image_url"`
	DisplayOrder int       `db:"display_order" json:"display_order"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}
