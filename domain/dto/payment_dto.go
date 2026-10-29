package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/rescuebyte/rescuebyte-be/domain/entity"
)

type PaymentResponse struct {
	ID              uuid.UUID `json:"id"`
	OrderID         uuid.UUID `json:"order_id"`
	PaymentType     string    `json:"payment_type"`
	PaymentChannel  *string   `json:"payment_channel,omitempty"`
	Amount          int64     `json:"amount"`
	Status          string    `json:"status"`
	SnapToken       *string   `json:"snap_token,omitempty"`
	SnapRedirectURL *string   `json:"snap_redirect_url,omitempty"`
	PaidAt          *string   `json:"paid_at,omitempty"`
}

type MidtransNotificationPayload struct {
	OrderID           string `json:"order_id"`
	TransactionStatus string `json:"transaction_status"`
	FraudStatus       string `json:"fraud_status,omitempty"`
	StatusCode        string `json:"status_code"`
	GrossAmount       string `json:"gross_amount"`
	PaymentType       string `json:"payment_type"`
	TransactionTime   string `json:"transaction_time"`
	TransactionID     string `json:"transaction_id"`
	SignatureKey      string `json:"signature_key"`
}

func FormatToPaymentResponse(p entity.Payment) PaymentResponse {
	var paidAtStr *string
	if p.PaidAt != nil {
		s := p.PaidAt.Format(time.RFC3339)
		paidAtStr = &s
	}

	return PaymentResponse{
		ID:              p.ID,
		OrderID:         p.OrderID,
		PaymentType:     p.PaymentType,
		PaymentChannel:  p.PaymentChannel,
		Amount:          p.Amount,
		Status:          string(p.Status),
		SnapToken:       p.SnapToken,
		SnapRedirectURL: p.SnapRedirectURL,
		PaidAt:          paidAtStr,
	}
}
