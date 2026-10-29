package entity

import (
	"time"

	"github.com/google/uuid"
)

type Cart struct {
	ID        uuid.UUID  `db:"id" json:"id"`
	UserID    uuid.UUID  `db:"user_id" json:"user_id"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	Items     []CartItem `db:"-" json:"items,omitempty"`
}

type CartItem struct {
	ID               uuid.UUID `db:"id" json:"id"`
	CartID           uuid.UUID `db:"cart_id" json:"cart_id"`
	ProductVariantID uuid.UUID `db:"product_variant_id" json:"product_variant_id"`
	Quantity         int       `db:"quantity" json:"quantity"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}

type CartItemDetail struct {
	ID               uuid.UUID `db:"id"`
	CartID           uuid.UUID `db:"cart_id"`
	ProductVariantID uuid.UUID `db:"product_variant_id"`
	Quantity         int       `db:"quantity"`
	ProductID        uuid.UUID `db:"product_id"`
	ProductName      string    `db:"product_name"`
	ProductDesc      *string   `db:"product_description"`
	ProductImageURL  *string   `db:"product_image_url"`
	VariantName      string    `db:"variant_name"`
	NormalPrice      int64     `db:"normal_price"`
	RescuePrice      int64     `db:"rescue_price"`
	CurrentStock     int       `db:"current_stock"`
	VariantActive    bool      `db:"variant_active"`
	ProductActive    bool      `db:"product_active"`
	StoreID          uuid.UUID `db:"store_id"`
	StoreName        string    `db:"store_name"`
	StoreAddress     string    `db:"store_address"`
	StoreCity        string    `db:"store_city"`
	StoreDistrict    string    `db:"store_district"`
	StoreLatitude    float64   `db:"store_latitude"`
	StoreLongitude   float64   `db:"store_longitude"`
	OperationalStart string    `db:"operational_start"`
	OperationalEnd   string    `db:"operational_end"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}
