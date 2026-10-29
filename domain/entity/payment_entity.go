package entity

import (
	"time"

	"github.com/google/uuid"

	"github.com/rescuebyte/rescuebyte-be/domain/enums"
)

type Payment struct {
	ID              uuid.UUID           `db:"id" json:"id"`
	OrderID         uuid.UUID           `db:"order_id" json:"order_id"`
	PaymentType     string              `db:"payment_type" json:"payment_type"`
	PaymentChannel  *string             `db:"payment_channel" json:"payment_channel"`
	Amount          int64               `db:"amount" json:"amount"`
	Status          enums.PaymentStatus `db:"status" json:"status"`
	TransactionID   *string             `db:"transaction_id" json:"transaction_id"`
	SnapToken       *string             `db:"snap_token" json:"snap_token"`
	SnapRedirectURL *string             `db:"snap_redirect_url" json:"snap_redirect_url"`
	RawResponse     *string             `db:"raw_response" json:"raw_response,omitempty"`
	PaidAt          *time.Time          `db:"paid_at" json:"paid_at"`
	CreatedAt       time.Time           `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time           `db:"updated_at" json:"updated_at"`
}
