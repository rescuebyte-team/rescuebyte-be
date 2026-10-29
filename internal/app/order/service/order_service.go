package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"

	"github.com/rescuebyte/rescuebyte-be/domain/dto"
	"github.com/rescuebyte/rescuebyte-be/domain/entity"
	"github.com/rescuebyte/rescuebyte-be/domain/enums"
	"github.com/rescuebyte/rescuebyte-be/domain/errx"
	"github.com/rescuebyte/rescuebyte-be/domain/helpers/pagination"
)

func generateOrderNumber() string {
	nowStr := time.Now().Format("060102")
	n, err := rand.Int(rand.Reader, big.NewInt(90000))
	randNum := int64(10000)
	if err == nil {
		randNum = 10000 + n.Int64()
	}
	return fmt.Sprintf("RB%s-%05d", nowStr, randNum)
}

func (s *orderService) Checkout(ctx context.Context, userID uuid.UUID, req dto.CheckoutRequest) (*dto.CheckoutResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}

	pickupDate, err := time.Parse("2006-01-02", req.PickupDate)
	if err != nil {
		return nil, errx.ErrBadRequest.WithMessage("Invalid pickup_date format, expected YYYY-MM-DD")
	}

	store, err := s.storeRepo.GetStoreByID(ctx, req.StoreID)
	if err != nil {
		return nil, err
	}
	if store == nil {
		return nil, errx.ErrStoreNotFound
	}

	type preparedItem struct {
		variant *entity.ProductVariant
		product *entity.Product
		qty     int
	}

	preparedItems := make([]preparedItem, 0, len(req.Items))
	var subtotal, totalOriginal int64

	for _, it := range req.Items {
		variant, err := s.productRepo.GetVariantByID(ctx, it.ProductVariantID)
		if err != nil {
			return nil, err
		}
		if !variant.IsActive {
			return nil, errx.ErrVariantInactive
		}

		product, err := s.productRepo.GetProductByID(ctx, variant.ProductID)
		if err != nil {
			return nil, err
		}
		if product.StoreID != req.StoreID {
			return nil, errx.ErrStoreDifferentInOrder
		}

		subtotal += variant.RescuePrice * int64(it.Quantity)
		totalOriginal += variant.NormalPrice * int64(it.Quantity)

		preparedItems = append(preparedItems, preparedItem{
			variant: variant,
			product: product,
			qty:     it.Quantity,
		})
	}

	orderID, err := s.uuid.NewV7()
	if err != nil {
		return nil, err
	}

	orderNumber := generateOrderNumber()
	now := time.Now()

	var serviceFee int64 = 2000
	var packagingFee int64 = 0
	totalAmount := subtotal + serviceFee + packagingFee
	totalDiscount := totalOriginal - subtotal
	if totalDiscount < 0 {
		totalDiscount = 0
	}

	initialStatus := enums.OrderStatusWaitingPayment
	if req.PaymentMethod == string(enums.PaymentMethodCashOnPickup) {
		initialStatus = enums.OrderStatusWaitingConfirmation
	}

	order := entity.Order{
		ID:                 orderID,
		OrderNumber:        orderNumber,
		UserID:             userID,
		StoreID:            req.StoreID,
		CustomerName:       req.CustomerName,
		CustomerPhone:      req.CustomerPhone,
		CustomerEmail:      req.CustomerEmail,
		PickupDate:         pickupDate,
		PickupTimeStart:    req.PickupTimeStart,
		PickupTimeEnd:      req.PickupTimeEnd,
		TotalOriginalPrice: totalOriginal,
		TotalDiscount:      totalDiscount,
		Subtotal:           subtotal,
		ServiceFee:         serviceFee,
		PackagingFee:       packagingFee,
		TotalAmount:        totalAmount,
		Status:             initialStatus,
		PaymentMethod:      req.PaymentMethod,
		Notes:              req.Notes,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	orderItems := make([]entity.OrderItem, 0, len(preparedItems))
	midtransItems := make([]dto.SnapItemDetail, 0, len(preparedItems)+1)

	for _, p := range preparedItems {
		itemID, err := s.uuid.NewV7()
		if err != nil {
			return nil, err
		}
		itemSubtotal := p.variant.RescuePrice * int64(p.qty)

		orderItems = append(orderItems, entity.OrderItem{
			ID:               itemID,
			OrderID:          orderID,
			ProductID:        p.product.ID,
			ProductVariantID: p.variant.ID,
			ProductName:      p.product.Name,
			VariantName:      p.variant.Name,
			ProductImageURL:  p.product.MainImageURL,
			NormalPrice:      p.variant.NormalPrice,
			RescuePrice:      p.variant.RescuePrice,
			Quantity:         p.qty,
			Subtotal:         itemSubtotal,
			CreatedAt:        now,
		})

		midtransItems = append(midtransItems, dto.SnapItemDetail{
			ID:       p.variant.ID.String(),
			Price:    p.variant.RescuePrice,
			Quantity: p.qty,
			Name:     p.product.Name,
		})
	}

	// Service fee for Midtrans
	if serviceFee > 0 {
		midtransItems = append(midtransItems, dto.SnapItemDetail{
			ID:       "fee-service",
			Price:    serviceFee,
			Quantity: 1,
			Name:     "Biaya Layanan RescueBite",
		})
	}

	var snapToken *string
	var snapRedirectURL *string
	paymentID, err := s.uuid.NewV7()
	if err != nil {
		return nil, err
	}

	var payment entity.Payment

	if req.PaymentMethod == string(enums.PaymentMethodOnline) {
		snapReq := dto.SnapTransactionRequest{
			TransactionDetails: dto.SnapTransactionDetails{
				OrderID:     orderNumber,
				GrossAmount: totalAmount,
			},
			CustomerDetails: dto.SnapCustomerDetails{
				FirstName: req.CustomerName,
				Email:     req.CustomerEmail,
				Phone:     req.CustomerPhone,
			},
			ItemDetails: midtransItems,
		}

		snapResp, err := s.midtransSvc.CreateSnapTransaction(ctx, snapReq)
		if err != nil {
			return nil, errx.ErrInternalServer.WithError(err)
		}

		snapToken = &snapResp.Token
		snapRedirectURL = &snapResp.RedirectURL

		payment = entity.Payment{
			ID:              paymentID,
			OrderID:         orderID,
			PaymentType:     string(enums.PaymentMethodOnline),
			Amount:          totalAmount,
			Status:          enums.PaymentStatusPending,
			SnapToken:       snapToken,
			SnapRedirectURL: snapRedirectURL,
			CreatedAt:       now,
			UpdatedAt:       now,
		}
	} else {
		payment = entity.Payment{
			ID:          paymentID,
			OrderID:     orderID,
			PaymentType: string(enums.PaymentMethodCashOnPickup),
			Amount:      totalAmount,
			Status:      enums.PaymentStatusPending,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
	}

	err = s.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		// 1. Deduct stock atomically
		for _, p := range preparedItems {
			if err := s.productRepo.DeductVariantStock(txCtx, p.variant.ID, p.qty); err != nil {
				return err
			}
		}

		// 2. Insert order
		if err := s.orderRepo.CreateOrder(txCtx, order); err != nil {
			return err
		}

		// 3. Insert order items
		if err := s.orderRepo.CreateOrderItems(txCtx, orderItems); err != nil {
			return err
		}

		// 4. Insert payment
		if err := s.paymentRepo.CreatePayment(txCtx, payment); err != nil {
			return err
		}

		// 5. Remove items from cart
		cart, err := s.cartRepo.GetOrCreateCartByUserID(txCtx, userID)
		if err == nil {
			for _, p := range preparedItems {
				item, err := s.cartRepo.GetCartItemByVariantID(txCtx, cart.ID, p.variant.ID)
				if err == nil && item != nil {
					_ = s.cartRepo.DeleteCartItem(txCtx, item.ID)
				}
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	order.Items = orderItems
	res := dto.FormatToOrderDetailResponse(order, *store, &payment)

	return &dto.CheckoutResponse{
		Order:           res,
		SnapToken:       snapToken,
		SnapRedirectURL: snapRedirectURL,
	}, nil
}

func (s *orderService) GetMyOrders(ctx context.Context, userID uuid.UUID, query dto.GetOrdersQuery) (dto.GetOrdersResponse, error) {
	pagination.Normalize(&query.Limit, &query.Page, &query.Offset, "order")

	items, total, err := s.orderRepo.GetCustomerOrders(ctx, userID, query)
	if err != nil {
		return dto.GetOrdersResponse{}, err
	}

	responses := make([]dto.OrderListItemResponse, 0, len(items))
	for _, it := range items {
		pickupStart := dto.FormatTimeStr(it.PickupTimeStart)
		pickupEnd := dto.FormatTimeStr(it.PickupTimeEnd)
		pickupSchedule := fmt.Sprintf("%s | %s – %s WIB", it.PickupDate.Format("02 January 2006"), pickupStart, pickupEnd)

		responses = append(responses, dto.OrderListItemResponse{
			ID:              it.ID,
			OrderNumber:     it.OrderNumber,
			StoreID:         it.StoreID,
			StoreName:       it.StoreName,
			StoreAddress:    it.StoreAddress,
			PickupSchedule:  pickupSchedule,
			TotalAmount:     it.TotalAmount,
			Status:          string(it.Status),
			StatusLabel:     dto.FormatStatusLabel(it.Status),
			PaymentMethod:   it.PaymentMethod,
			TotalItems:      it.TotalItems,
			SummaryItemName: it.SummaryItemName,
			SummaryImageURL: it.SummaryImageURL,
			CreatedAt:       dto.FormatTimeWIB(it.CreatedAt),
		})
	}

	return dto.GetOrdersResponse{
		Orders: responses,
		Meta:   pagination.NewPaginationResponse(total, query.Page, query.Limit),
	}, nil
}

func (s *orderService) GetOrderDetail(ctx context.Context, userID uuid.UUID, orderID uuid.UUID) (*dto.OrderDetailResponse, error) {
	order, err := s.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	store, err := s.storeRepo.GetStoreByID(ctx, order.StoreID)
	if err != nil {
		return nil, err
	}
	if store == nil {
		return nil, errx.ErrStoreNotFound
	}

	// Permission: must be customer or store owner
	if order.UserID != userID && store.UserID != userID {
		return nil, errx.ErrForbidden
	}

	items, err := s.orderRepo.GetOrderItems(ctx, order.ID)
	if err != nil {
		return nil, err
	}
	order.Items = items

	payment, err := s.paymentRepo.GetPaymentByOrderID(ctx, order.ID)
	if err != nil {
		return nil, err
	}

	res := dto.FormatToOrderDetailResponse(*order, *store, payment)
	return &res, nil
}

func (s *orderService) CancelOrder(ctx context.Context, userID uuid.UUID, orderID uuid.UUID, req dto.CancelOrderRequest) error {
	order, err := s.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order.UserID != userID {
		return errx.ErrForbidden
	}

	if order.Status != enums.OrderStatusWaitingPayment && order.Status != enums.OrderStatusWaitingConfirmation {
		return errx.ErrOrderCannotBeCancelled
	}

	items, err := s.orderRepo.GetOrderItems(ctx, order.ID)
	if err != nil {
		return err
	}

	reason := "Dibatalkan oleh pembeli"
	if req.Reason != "" {
		reason = req.Reason
	}

	return s.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		// Restore variant stock
		for _, it := range items {
			if err := s.productRepo.RestoreVariantStock(txCtx, it.ProductVariantID, it.Quantity); err != nil {
				return err
			}
		}

		// Update order status to CANCELLED
		if err := s.orderRepo.UpdateOrderStatus(txCtx, order.ID, enums.OrderStatusCancelled, "cancelled_at", &reason); err != nil {
			return err
		}

		// Update payment status to CANCEL
		_ = s.paymentRepo.UpdatePaymentStatus(txCtx, order.ID, enums.PaymentStatusCancel, nil, nil, nil, nil)

		return nil
	})
}

func (s *orderService) UpdateOrderStatusByMerchant(ctx context.Context, userID uuid.UUID, orderID uuid.UUID, req dto.UpdateOrderStatusRequest) (*dto.OrderDetailResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}

	order, err := s.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	store, err := s.storeRepo.GetStoreByID(ctx, order.StoreID)
	if err != nil {
		return nil, err
	}
	if store == nil || store.UserID != userID {
		return nil, errx.ErrUnauthorizedStoreAccess
	}

	targetStatus := enums.OrderStatus(req.Status)
	timestampField := ""

	switch targetStatus {
	case enums.OrderStatusInProcess:
		timestampField = "confirmed_at"
	case enums.OrderStatusReadyForPickup:
		timestampField = "ready_at"
	case enums.OrderStatusCompleted:
		timestampField = "completed_at"
	case enums.OrderStatusCancelled:
		timestampField = "cancelled_at"
	}

	err = s.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		if targetStatus == enums.OrderStatusCancelled {
			items, err := s.orderRepo.GetOrderItems(txCtx, order.ID)
			if err != nil {
				return err
			}
			for _, it := range items {
				if err := s.productRepo.RestoreVariantStock(txCtx, it.ProductVariantID, it.Quantity); err != nil {
					return err
				}
			}
			_ = s.paymentRepo.UpdatePaymentStatus(txCtx, order.ID, enums.PaymentStatusCancel, nil, nil, nil, nil)
		}

		return s.orderRepo.UpdateOrderStatus(txCtx, order.ID, targetStatus, timestampField, req.CancellationReason)
	})
	if err != nil {
		return nil, err
	}

	return s.GetOrderDetail(ctx, userID, orderID)
}
