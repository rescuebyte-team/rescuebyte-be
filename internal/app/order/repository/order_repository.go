package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"github.com/rescuebyte/rescuebyte-be/domain/dto"
	"github.com/rescuebyte/rescuebyte-be/domain/entity"
	"github.com/rescuebyte/rescuebyte-be/domain/enums"
	"github.com/rescuebyte/rescuebyte-be/domain/errx"
)

var orderSelectFields = []string{
	"id", "order_number", "user_id", "store_id", "customer_name", "customer_phone", "customer_email",
	"pickup_date", "pickup_time_start::text", "pickup_time_end::text", "total_original_price", "total_discount",
	"subtotal", "service_fee", "packaging_fee", "total_amount", "status", "payment_method", "notes",
	"created_at", "paid_at", "confirmed_at", "ready_at", "completed_at", "cancelled_at", "cancellation_reason", "updated_at",
}

func (r *orderRepository) CreateOrder(ctx context.Context, order entity.Order) error {
	query, args, err := sq.
		Insert("orders").
		Columns(
			"id", "order_number", "user_id", "store_id", "customer_name", "customer_phone", "customer_email",
			"pickup_date", "pickup_time_start", "pickup_time_end", "total_original_price", "total_discount",
			"subtotal", "service_fee", "packaging_fee", "total_amount", "status", "payment_method", "notes",
			"created_at", "paid_at", "confirmed_at", "ready_at", "completed_at", "cancelled_at", "cancellation_reason", "updated_at",
		).
		Values(
			order.ID, order.OrderNumber, order.UserID, order.StoreID, order.CustomerName, order.CustomerPhone, order.CustomerEmail,
			order.PickupDate, order.PickupTimeStart, order.PickupTimeEnd, order.TotalOriginalPrice, order.TotalDiscount,
			order.Subtotal, order.ServiceFee, order.PackagingFee, order.TotalAmount, order.Status, order.PaymentMethod, order.Notes,
			order.CreatedAt, order.PaidAt, order.ConfirmedAt, order.ReadyAt, order.CompletedAt, order.CancelledAt, order.CancellationReason, order.UpdatedAt,
		).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.dbGetter(ctx).ExecContext(ctx, query, args...)
	return err
}

func (r *orderRepository) CreateOrderItems(ctx context.Context, items []entity.OrderItem) error {
	if len(items) == 0 {
		return nil
	}

	builder := sq.
		Insert("order_items").
		Columns(
			"id", "order_id", "product_id", "product_variant_id", "product_name",
			"variant_name", "product_image_url", "normal_price", "rescue_price", "quantity", "subtotal", "created_at",
		)

	for _, it := range items {
		builder = builder.Values(
			it.ID, it.OrderID, it.ProductID, it.ProductVariantID, it.ProductName,
			it.VariantName, it.ProductImageURL, it.NormalPrice, it.RescuePrice, it.Quantity, it.Subtotal, it.CreatedAt,
		)
	}

	query, args, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}

	_, err = r.dbGetter(ctx).ExecContext(ctx, query, args...)
	return err
}

func (r *orderRepository) GetOrderByID(ctx context.Context, id uuid.UUID) (*entity.Order, error) {
	query, args, err := sq.
		Select(orderSelectFields...).
		From("orders").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var order entity.Order
	err = r.dbGetter(ctx).GetContext(ctx, &order, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errx.ErrOrderNotFound
		}
		return nil, err
	}

	return &order, nil
}

func (r *orderRepository) GetOrderByNumber(ctx context.Context, orderNumber string) (*entity.Order, error) {
	query, args, err := sq.
		Select(orderSelectFields...).
		From("orders").
		Where(sq.Eq{"order_number": orderNumber}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var order entity.Order
	err = r.dbGetter(ctx).GetContext(ctx, &order, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errx.ErrOrderNotFound
		}
		return nil, err
	}

	return &order, nil
}

func (r *orderRepository) GetOrderItems(ctx context.Context, orderID uuid.UUID) ([]entity.OrderItem, error) {
	query, args, err := sq.
		Select(
			"id", "order_id", "product_id", "product_variant_id", "product_name",
			"variant_name", "product_image_url", "normal_price", "rescue_price", "quantity", "subtotal", "created_at",
		).
		From("order_items").
		Where(sq.Eq{"order_id": orderID}).
		OrderBy("created_at ASC").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var items []entity.OrderItem
	err = r.dbGetter(ctx).SelectContext(ctx, &items, query, args...)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *orderRepository) GetCustomerOrders(ctx context.Context, userID uuid.UUID, query dto.GetOrdersQuery) ([]entity.OrderListItem, int64, error) {
	builder := sq.
		Select(
			"o.id",
			"o.order_number",
			"o.store_id",
			"s.name AS store_name",
			"s.address AS store_address",
			"s.city AS store_city",
			"s.district AS store_district",
			"o.pickup_date",
			"o.pickup_time_start::text AS pickup_time_start",
			"o.pickup_time_end::text AS pickup_time_end",
			"o.total_amount",
			"o.status",
			"o.payment_method",
			"o.created_at",
			"COALESCE(COUNT(oi.id), 0)::int AS total_items",
			"COALESCE(MIN(oi.product_name), '') AS summary_item_name",
			"MIN(oi.product_image_url) AS summary_image_url",
		).
		From("orders o").
		Join("stores s ON s.id = o.store_id").
		LeftJoin("order_items oi ON oi.order_id = o.id").
		Where(sq.Eq{"o.user_id": userID})

	countBuilder := sq.
		Select("COUNT(DISTINCT o.id)").
		From("orders o").
		Join("stores s ON s.id = o.store_id").
		LeftJoin("order_items oi ON oi.order_id = o.id").
		Where(sq.Eq{"o.user_id": userID})

	// Tab filtering
	switch query.Tab {
	case "waiting_payment":
		builder = builder.Where(sq.Eq{"o.status": enums.OrderStatusWaitingPayment})
		countBuilder = countBuilder.Where(sq.Eq{"o.status": enums.OrderStatusWaitingPayment})
	case "waiting_confirmation":
		builder = builder.Where(sq.Eq{"o.status": enums.OrderStatusWaitingConfirmation})
		countBuilder = countBuilder.Where(sq.Eq{"o.status": enums.OrderStatusWaitingConfirmation})
	case "ready_for_pickup":
		builder = builder.Where(sq.Eq{"o.status": enums.OrderStatusReadyForPickup})
		countBuilder = countBuilder.Where(sq.Eq{"o.status": enums.OrderStatusReadyForPickup})
	case "completed":
		builder = builder.Where(sq.Eq{"o.status": enums.OrderStatusCompleted})
		countBuilder = countBuilder.Where(sq.Eq{"o.status": enums.OrderStatusCompleted})
	case "cancelled":
		builder = builder.Where(sq.Eq{"o.status": enums.OrderStatusCancelled})
		countBuilder = countBuilder.Where(sq.Eq{"o.status": enums.OrderStatusCancelled})
	}

	// Search filter
	if query.Search != "" {
		pattern := fmt.Sprintf("%%%s%%", query.Search)
		searchClause := sq.Or{
			sq.ILike{"o.order_number": pattern},
			sq.ILike{"s.name": pattern},
			sq.ILike{"oi.product_name": pattern},
		}
		builder = builder.Where(searchClause)
		countBuilder = countBuilder.Where(searchClause)
	}

	// Count total
	countSql, countArgs, err := countBuilder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, 0, err
	}
	var total int64
	if err := r.dbGetter(ctx).GetContext(ctx, &total, countSql, countArgs...); err != nil {
		return nil, 0, err
	}

	// Grouping & Sorting
	builder = builder.GroupBy("o.id", "s.id").OrderBy("o.created_at DESC")

	// Pagination
	if query.Limit > 0 {
		builder = builder.Limit(query.Limit).Offset(query.Offset)
	}

	listSql, listArgs, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, 0, err
	}

	var items []entity.OrderListItem
	if err := r.dbGetter(ctx).SelectContext(ctx, &items, listSql, listArgs...); err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *orderRepository) GetStoreOrders(ctx context.Context, storeID uuid.UUID, query dto.GetOrdersQuery) ([]entity.OrderListItem, int64, error) {
	builder := sq.
		Select(
			"o.id",
			"o.order_number",
			"o.store_id",
			"s.name AS store_name",
			"s.address AS store_address",
			"s.city AS store_city",
			"s.district AS store_district",
			"o.pickup_date",
			"o.pickup_time_start::text AS pickup_time_start",
			"o.pickup_time_end::text AS pickup_time_end",
			"o.total_amount",
			"o.status",
			"o.payment_method",
			"o.created_at",
			"COALESCE(COUNT(oi.id), 0)::int AS total_items",
			"COALESCE(MIN(oi.product_name), '') AS summary_item_name",
			"MIN(oi.product_image_url) AS summary_image_url",
		).
		From("orders o").
		Join("stores s ON s.id = o.store_id").
		LeftJoin("order_items oi ON oi.order_id = o.id").
		Where(sq.Eq{"o.store_id": storeID})

	countBuilder := sq.
		Select("COUNT(DISTINCT o.id)").
		From("orders o").
		Join("stores s ON s.id = o.store_id").
		LeftJoin("order_items oi ON oi.order_id = o.id").
		Where(sq.Eq{"o.store_id": storeID})

	if query.Tab != "" && query.Tab != "all" {
		builder = builder.Where(sq.Eq{"o.status": query.Tab})
		countBuilder = countBuilder.Where(sq.Eq{"o.status": query.Tab})
	}

	if query.Search != "" {
		pattern := fmt.Sprintf("%%%s%%", query.Search)
		searchClause := sq.Or{
			sq.ILike{"o.order_number": pattern},
			sq.ILike{"o.customer_name": pattern},
			sq.ILike{"oi.product_name": pattern},
		}
		builder = builder.Where(searchClause)
		countBuilder = countBuilder.Where(searchClause)
	}

	countSql, countArgs, err := countBuilder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, 0, err
	}
	var total int64
	if err := r.dbGetter(ctx).GetContext(ctx, &total, countSql, countArgs...); err != nil {
		return nil, 0, err
	}

	builder = builder.GroupBy("o.id", "s.id").OrderBy("o.created_at DESC")

	if query.Limit > 0 {
		builder = builder.Limit(query.Limit).Offset(query.Offset)
	}

	listSql, listArgs, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, 0, err
	}

	var items []entity.OrderListItem
	if err := r.dbGetter(ctx).SelectContext(ctx, &items, listSql, listArgs...); err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *orderRepository) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status enums.OrderStatus, timestampField string, cancelReason *string) error {
	builder := sq.
		Update("orders").
		Set("status", status).
		Set("updated_at", time.Now())

	if timestampField != "" {
		builder = builder.Set(timestampField, time.Now())
	}
	if cancelReason != nil {
		builder = builder.Set("cancellation_reason", *cancelReason)
	}

	query, args, err := builder.
		Where(sq.Eq{"id": orderID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.dbGetter(ctx).ExecContext(ctx, query, args...)
	return err
}
