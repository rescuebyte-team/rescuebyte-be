package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"github.com/rescuebyte/rescuebyte-be/domain/entity"
	"github.com/rescuebyte/rescuebyte-be/domain/enums"
)

var paymentSelectFields = []string{
	"id", "order_id", "payment_type", "payment_channel", "amount", "status",
	"transaction_id", "snap_token", "snap_redirect_url", "raw_response::text", "paid_at", "created_at", "updated_at",
}

func (r *paymentRepository) CreatePayment(ctx context.Context, payment entity.Payment) error {
	query, args, err := sq.
		Insert("payments").
		Columns(
			"id", "order_id", "payment_type", "payment_channel", "amount", "status",
			"transaction_id", "snap_token", "snap_redirect_url", "raw_response", "paid_at", "created_at", "updated_at",
		).
		Values(
			payment.ID, payment.OrderID, payment.PaymentType, payment.PaymentChannel, payment.Amount, payment.Status,
			payment.TransactionID, payment.SnapToken, payment.SnapRedirectURL, payment.RawResponse, payment.PaidAt, payment.CreatedAt, payment.UpdatedAt,
		).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.dbGetter(ctx).ExecContext(ctx, query, args...)
	return err
}

func (r *paymentRepository) GetPaymentByOrderID(ctx context.Context, orderID uuid.UUID) (*entity.Payment, error) {
	query, args, err := sq.
		Select(paymentSelectFields...).
		From("payments").
		Where(sq.Eq{"order_id": orderID}).
		OrderBy("created_at DESC").
		Limit(1).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var payment entity.Payment
	err = r.dbGetter(ctx).GetContext(ctx, &payment, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &payment, nil
}

func (r *paymentRepository) UpdatePaymentStatus(
	ctx context.Context,
	orderID uuid.UUID,
	status enums.PaymentStatus,
	channel *string,
	transactionID *string,
	rawResponse *string,
	paidAt *time.Time,
) error {
	builder := sq.
		Update("payments").
		Set("status", status).
		Set("updated_at", time.Now())

	if channel != nil {
		builder = builder.Set("payment_channel", *channel)
	}
	if transactionID != nil {
		builder = builder.Set("transaction_id", *transactionID)
	}
	if rawResponse != nil {
		builder = builder.Set("raw_response", *rawResponse)
	}
	if paidAt != nil {
		builder = builder.Set("paid_at", *paidAt)
	}

	query, args, err := builder.
		Where(sq.Eq{"order_id": orderID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.dbGetter(ctx).ExecContext(ctx, query, args...)
	return err
}
