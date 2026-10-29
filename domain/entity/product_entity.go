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

type ProductCardItem struct {
	ID                uuid.UUID `db:"id" json:"id"`
	StoreID           uuid.UUID `db:"store_id" json:"store_id"`
	StoreName         string    `db:"store_name" json:"store_name"`
	Name              string    `db:"name" json:"name"`
	Description       *string   `db:"description" json:"description"`
	MainImageURL      *string   `db:"main_image_url" json:"main_image_url"`
	OperationalStart  string    `db:"operational_start" json:"operational_start"`
	OperationalEnd    string    `db:"operational_end" json:"operational_end"`
	StoreCity         string    `db:"store_city" json:"store_city"`
	StoreDistrict     string    `db:"store_district" json:"store_district"`
	StoreLatitude     float64   `db:"store_latitude" json:"store_latitude"`
	StoreLongitude    float64   `db:"store_longitude" json:"store_longitude"`
	TotalCurrentStock int       `db:"total_current_stock" json:"total_current_stock"`
	MinNormalPrice    int64     `db:"min_normal_price" json:"min_normal_price"`
	MinRescuePrice    int64     `db:"min_rescue_price" json:"min_rescue_price"`
}
