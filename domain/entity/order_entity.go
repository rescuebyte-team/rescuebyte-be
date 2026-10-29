package entity

import (
	"time"

	"github.com/google/uuid"

	"github.com/rescuebyte/rescuebyte-be/domain/enums"
)

type Order struct {
	ID                 uuid.UUID         `db:"id" json:"id"`
	OrderNumber        string            `db:"order_number" json:"order_number"`
	UserID             uuid.UUID         `db:"user_id" json:"user_id"`
	StoreID            uuid.UUID         `db:"store_id" json:"store_id"`
	CustomerName       string            `db:"customer_name" json:"customer_name"`
	CustomerPhone      string            `db:"customer_phone" json:"customer_phone"`
	CustomerEmail      string            `db:"customer_email" json:"customer_email"`
	PickupDate         time.Time         `db:"pickup_date" json:"pickup_date"`
	PickupTimeStart    string            `db:"pickup_time_start" json:"pickup_time_start"`
	PickupTimeEnd      string            `db:"pickup_time_end" json:"pickup_time_end"`
	TotalOriginalPrice int64             `db:"total_original_price" json:"total_original_price"`
	TotalDiscount      int64             `db:"total_discount" json:"total_discount"`
	Subtotal           int64             `db:"subtotal" json:"subtotal"`
	ServiceFee         int64             `db:"service_fee" json:"service_fee"`
	PackagingFee       int64             `db:"packaging_fee" json:"packaging_fee"`
	TotalAmount        int64             `db:"total_amount" json:"total_amount"`
	Status             enums.OrderStatus `db:"status" json:"status"`
	PaymentMethod      string            `db:"payment_method" json:"payment_method"`
	Notes              *string           `db:"notes" json:"notes"`
	CreatedAt          time.Time         `db:"created_at" json:"created_at"`
	PaidAt             *time.Time        `db:"paid_at" json:"paid_at"`
	ConfirmedAt        *time.Time        `db:"confirmed_at" json:"confirmed_at"`
	ReadyAt            *time.Time        `db:"ready_at" json:"ready_at"`
	CompletedAt        *time.Time        `db:"completed_at" json:"completed_at"`
	CancelledAt        *time.Time        `db:"cancelled_at" json:"cancelled_at"`
	CancellationReason *string           `db:"cancellation_reason" json:"cancellation_reason"`
	UpdatedAt          time.Time         `db:"updated_at" json:"updated_at"`

	Items []OrderItem `db:"-" json:"items,omitempty"`
	Store *Store      `db:"-" json:"store,omitempty"`
}

type OrderItem struct {
	ID               uuid.UUID `db:"id" json:"id"`
	OrderID          uuid.UUID `db:"order_id" json:"order_id"`
	ProductID        uuid.UUID `db:"product_id" json:"product_id"`
	ProductVariantID uuid.UUID `db:"product_variant_id" json:"product_variant_id"`
	ProductName      string    `db:"product_name" json:"product_name"`
	VariantName      string    `db:"variant_name" json:"variant_name"`
	ProductImageURL  *string   `db:"product_image_url" json:"product_image_url"`
	NormalPrice      int64     `db:"normal_price" json:"normal_price"`
	RescuePrice      int64     `db:"rescue_price" json:"rescue_price"`
	Quantity         int       `db:"quantity" json:"quantity"`
	Subtotal         int64     `db:"subtotal" json:"subtotal"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
}

type OrderListItem struct {
	ID              uuid.UUID         `db:"id"`
	OrderNumber     string            `db:"order_number"`
	StoreID         uuid.UUID         `db:"store_id"`
	StoreName       string            `db:"store_name"`
	StoreAddress    string            `db:"store_address"`
	StoreCity       string            `db:"store_city"`
	StoreDistrict   string            `db:"store_district"`
	PickupDate      time.Time         `db:"pickup_date"`
	PickupTimeStart string            `db:"pickup_time_start"`
	PickupTimeEnd   string            `db:"pickup_time_end"`
	TotalAmount     int64             `db:"total_amount"`
	Status          enums.OrderStatus `db:"status"`
	PaymentMethod   string            `db:"payment_method"`
	CreatedAt       time.Time         `db:"created_at"`
	TotalItems      int               `db:"total_items"`
	SummaryItemName string            `db:"summary_item_name"`
	SummaryImageURL *string           `db:"summary_image_url"`
}
