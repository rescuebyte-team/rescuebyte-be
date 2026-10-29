package contracts

import (
	"context"

	"github.com/google/uuid"

	"github.com/rescuebyte/rescuebyte-be/domain/dto"
	"github.com/rescuebyte/rescuebyte-be/domain/entity"
)

type CartRepository interface {
	GetOrCreateCartByUserID(ctx context.Context, userID uuid.UUID) (*entity.Cart, error)
	GetCartItemByID(ctx context.Context, id uuid.UUID) (*entity.CartItem, error)
	GetCartItemByVariantID(ctx context.Context, cartID, variantID uuid.UUID) (*entity.CartItem, error)
	GetCartItemDetails(ctx context.Context, cartID uuid.UUID) ([]entity.CartItemDetail, error)
	GetCartItemsCount(ctx context.Context, cartID uuid.UUID) (int, error)
	UpsertCartItem(ctx context.Context, item entity.CartItem) error
	UpdateCartItemQuantity(ctx context.Context, id uuid.UUID, quantity int) error
	DeleteCartItem(ctx context.Context, id uuid.UUID) error
	DeleteCartItemsByStoreID(ctx context.Context, cartID, storeID uuid.UUID) error
	ClearCart(ctx context.Context, cartID uuid.UUID) error
}

type CartService interface {
	GetCart(ctx context.Context, userID uuid.UUID, userLat, userLon *float64) (*dto.CartResponse, error)
	GetCartCount(ctx context.Context, userID uuid.UUID) (*dto.CartCountResponse, error)
	AddToCart(ctx context.Context, userID uuid.UUID, req dto.AddToCartRequest) error
	UpdateCartItem(ctx context.Context, userID, itemID uuid.UUID, req dto.UpdateCartItemRequest) error
	DeleteCartItem(ctx context.Context, userID, itemID uuid.UUID) error
	DeleteStoreCartItems(ctx context.Context, userID, storeID uuid.UUID) error
	ClearCart(ctx context.Context, userID uuid.UUID) error
}
