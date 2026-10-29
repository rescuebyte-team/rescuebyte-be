package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/rescuebyte/rescuebyte-be/domain/dto"
	"github.com/rescuebyte/rescuebyte-be/domain/enums"
	"github.com/rescuebyte/rescuebyte-be/domain/errx"
	"github.com/rescuebyte/rescuebyte-be/pkg/log"
)

func (s *paymentService) HandleMidtransNotification(ctx context.Context, payload dto.MidtransNotificationPayload) error {
	// 1. Verify SHA-512 Signature
	isValid := s.midtransSvc.VerifySignature(payload.OrderID, payload.StatusCode, payload.GrossAmount, payload.SignatureKey)
	if !isValid {
		log.Warn(log.Fields{
			"order_id": payload.OrderID,
		}, "[Payment] Invalid midtrans signature verification")
		return errx.ErrInvalidPaymentSignature
	}

	// 2. Fetch order
	order, err := s.orderRepo.GetOrderByNumber(ctx, payload.OrderID)
	if err != nil {
		return err
	}

	now := time.Now()
	rawBytes, _ := json.Marshal(payload)
	rawStr := string(rawBytes)

	var targetPaymentStatus enums.PaymentStatus
	var targetOrderStatus *enums.OrderStatus
	var paidAt *time.Time
	shouldRestoreStock := false

	switch payload.TransactionStatus {
	case "capture":
		if payload.FraudStatus == "challenge" {
			targetPaymentStatus = enums.PaymentStatusPending
		} else {
			targetPaymentStatus = enums.PaymentStatusSettlement
			status := enums.OrderStatusWaitingConfirmation
			targetOrderStatus = &status
			paidAt = &now
		}
	case "settlement":
		targetPaymentStatus = enums.PaymentStatusSettlement
		status := enums.OrderStatusWaitingConfirmation
		targetOrderStatus = &status
		paidAt = &now
	case "pending":
		targetPaymentStatus = enums.PaymentStatusPending
	case "deny", "cancel", "expire":
		if payload.TransactionStatus == "expire" {
			targetPaymentStatus = enums.PaymentStatusExpire
		} else {
			targetPaymentStatus = enums.PaymentStatusCancel
		}
		status := enums.OrderStatusCancelled
		targetOrderStatus = &status
		shouldRestoreStock = true
	case "refund", "partial_refund":
		targetPaymentStatus = enums.PaymentStatusRefund
	default:
		log.Warn(log.Fields{
			"status": payload.TransactionStatus,
		}, "[Payment] Unknown transaction status from midtrans")
		return nil
	}

	channel := payload.PaymentType

	return s.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		// Update payment
		if err := s.paymentRepo.UpdatePaymentStatus(
			txCtx,
			order.ID,
			targetPaymentStatus,
			&channel,
			&payload.TransactionID,
			&rawStr,
			paidAt,
		); err != nil {
			return err
		}

		// Update order if status changed
		if targetOrderStatus != nil && order.Status != *targetOrderStatus {
			timestampField := ""
			var cancelReason *string

			if *targetOrderStatus == enums.OrderStatusWaitingConfirmation {
				timestampField = "paid_at"
			} else if *targetOrderStatus == enums.OrderStatusCancelled {
				timestampField = "cancelled_at"
				r := "Pembayaran gagal atau kadaluarsa (" + payload.TransactionStatus + ")"
				cancelReason = &r
			}

			if err := s.orderRepo.UpdateOrderStatus(txCtx, order.ID, *targetOrderStatus, timestampField, cancelReason); err != nil {
				return err
			}

			if shouldRestoreStock {
				items, err := s.orderRepo.GetOrderItems(txCtx, order.ID)
				if err == nil {
					for _, it := range items {
						_ = s.productRepo.RestoreVariantStock(txCtx, it.ProductVariantID, it.Quantity)
					}
				}
			}
		}

		return nil
	})
}

func (s *paymentService) GetPaymentStatus(ctx context.Context, userID uuid.UUID, orderID uuid.UUID) (*dto.PaymentResponse, error) {
	order, err := s.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if order.UserID != userID {
		return nil, errx.ErrForbidden
	}

	payment, err := s.paymentRepo.GetPaymentByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if payment == nil {
		return nil, errx.ErrNotFound
	}

	res := dto.FormatToPaymentResponse(*payment)
	return &res, nil
}
