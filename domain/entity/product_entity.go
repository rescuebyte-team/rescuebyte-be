package entity

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID           uuid.UUID        `db:"id" json:"id"`
	StoreID      uuid.UUID        `db:"store_id" json:"store_id"`
	Name         string           `db:"name" json:"name"`
	Description  *string          `db:"description" json:"description"`
	MainImageURL *string          `db:"main_image_url" json:"main_image_url"`
	IsActive     bool             `db:"is_active" json:"is_active"`
	CreatedAt    time.Time        `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time        `db:"updated_at" json:"updated_at"`
	Variants     []ProductVariant `db:"-" json:"variants,omitempty"`
	Images       []ProductImage   `db:"-" json:"images,omitempty"`
}
