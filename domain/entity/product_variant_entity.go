package entity

import (
	"time"

	"github.com/google/uuid"
)

type ProductVariant struct {
	ID           uuid.UUID `db:"id" json:"id"`
	ProductID    uuid.UUID `db:"product_id" json:"product_id"`
	Name         string    `db:"name" json:"name"`
	Sequence     int       `db:"sequence" json:"sequence"`
	NormalPrice  int64     `db:"normal_price" json:"normal_price"`
	RescuePrice  int64     `db:"rescue_price" json:"rescue_price"`
	InitialStock int       `db:"initial_stock" json:"initial_stock"`
	CurrentStock int       `db:"current_stock" json:"current_stock"`
	IsActive     bool      `db:"is_active" json:"is_active"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}
