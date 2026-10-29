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

type SnapTransactionDetails struct {
	OrderID     string `json:"order_id"`
	GrossAmount int64  `json:"gross_amount"`
}

type SnapCustomerDetails struct {
	FirstName string `json:"first_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

type SnapItemDetail struct {
	ID       string `json:"id"`
	Price    int64  `json:"price"`
	Quantity int    `json:"quantity"`
	Name     string `json:"name"`
}

type SnapTransactionRequest struct {
	TransactionDetails SnapTransactionDetails `json:"transaction_details"`
	CustomerDetails    SnapCustomerDetails    `json:"customer_details"`
	ItemDetails        []SnapItemDetail       `json:"item_details"`
}

type SnapTransactionResponse struct {
	Token         string   `json:"token"`
	RedirectURL   string   `json:"redirect_url"`
	ErrorMessages []string `json:"error_messages,omitempty"`
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
