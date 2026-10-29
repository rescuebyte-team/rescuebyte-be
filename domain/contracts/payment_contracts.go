package contracts

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/rescuebyte/rescuebyte-be/domain/dto"
	"github.com/rescuebyte/rescuebyte-be/domain/entity"
	"github.com/rescuebyte/rescuebyte-be/domain/enums"
)

type PaymentRepository interface {
	CreatePayment(ctx context.Context, payment entity.Payment) error
	GetPaymentByOrderID(ctx context.Context, orderID uuid.UUID) (*entity.Payment, error)
	UpdatePaymentStatus(ctx context.Context, orderID uuid.UUID, status enums.PaymentStatus, channel *string, transactionID *string, rawResponse *string, paidAt *time.Time) error
}

type PaymentService interface {
	HandleMidtransNotification(ctx context.Context, payload dto.MidtransNotificationPayload) error
	GetPaymentStatus(ctx context.Context, userID uuid.UUID, orderID uuid.UUID) (*dto.PaymentResponse, error)
}
