package contracts

import (
	"context"

	"github.com/google/uuid"

	"github.com/rescuebyte/rescuebyte-be/domain/dto"
	"github.com/rescuebyte/rescuebyte-be/domain/entity"
	"github.com/rescuebyte/rescuebyte-be/domain/enums"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order entity.Order) error
	CreateOrderItems(ctx context.Context, items []entity.OrderItem) error
	GetOrderByID(ctx context.Context, id uuid.UUID) (*entity.Order, error)
	GetOrderByNumber(ctx context.Context, orderNumber string) (*entity.Order, error)
	GetOrderItems(ctx context.Context, orderID uuid.UUID) ([]entity.OrderItem, error)
	GetCustomerOrders(ctx context.Context, userID uuid.UUID, query dto.GetOrdersQuery) ([]entity.OrderListItem, int64, error)
	GetStoreOrders(ctx context.Context, storeID uuid.UUID, query dto.GetOrdersQuery) ([]entity.OrderListItem, int64, error)
	UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status enums.OrderStatus, timestampField string, cancelReason *string) error
}

type OrderService interface {
	Checkout(ctx context.Context, userID uuid.UUID, req dto.CheckoutRequest) (*dto.CheckoutResponse, error)
	GetMyOrders(ctx context.Context, userID uuid.UUID, query dto.GetOrdersQuery) (dto.GetOrdersResponse, error)
	GetSellerOrders(ctx context.Context, userID uuid.UUID, query dto.GetOrdersQuery) (dto.GetOrdersResponse, error)
	GetOrderDetail(ctx context.Context, userID uuid.UUID, orderID uuid.UUID) (*dto.OrderDetailResponse, error)
	CancelOrder(ctx context.Context, userID uuid.UUID, orderID uuid.UUID, req dto.CancelOrderRequest) error
	UpdateOrderStatusByMerchant(ctx context.Context, userID uuid.UUID, orderID uuid.UUID, req dto.UpdateOrderStatusRequest) (*dto.OrderDetailResponse, error)
}
