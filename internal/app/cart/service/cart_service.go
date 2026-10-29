package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/rescuebyte/rescuebyte-be/domain/dto"
	"github.com/rescuebyte/rescuebyte-be/domain/entity"
	"github.com/rescuebyte/rescuebyte-be/domain/errx"
)

func (s *cartService) GetCart(ctx context.Context, userID uuid.UUID, userLat, userLon *float64) (*dto.CartResponse, error) {
	cart, err := s.cartRepo.GetOrCreateCartByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	details, err := s.cartRepo.GetCartItemDetails(ctx, cart.ID)
	if err != nil {
		return nil, err
	}

	res := dto.FormatToCartResponse(details, userLat, userLon)
	return &res, nil
}

func (s *cartService) GetCartCount(ctx context.Context, userID uuid.UUID) (*dto.CartCountResponse, error) {
	cart, err := s.cartRepo.GetOrCreateCartByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	count, err := s.cartRepo.GetCartItemsCount(ctx, cart.ID)
	if err != nil {
		return nil, err
	}

	return &dto.CartCountResponse{Count: count}, nil
}

func (s *cartService) AddToCart(ctx context.Context, userID uuid.UUID, req dto.AddToCartRequest) error {
	if err := s.validator.Validate(req); err != nil {
		return err
	}

	variant, err := s.productRepo.GetVariantByID(ctx, req.ProductVariantID)
	if err != nil {
		return err
	}
	if !variant.IsActive {
		return errx.ErrVariantInactive
	}

	cart, err := s.cartRepo.GetOrCreateCartByUserID(ctx, userID)
	if err != nil {
		return err
	}

	existingItem, err := s.cartRepo.GetCartItemByVariantID(ctx, cart.ID, req.ProductVariantID)
	if err != nil {
		return err
	}

	totalDesiredQty := req.Quantity
	if existingItem != nil {
		totalDesiredQty += existingItem.Quantity
	}

	if totalDesiredQty > variant.CurrentStock {
		return errx.ErrInsufficientStock
	}

	itemID, err := s.uuid.NewV7()
	if err != nil {
		return err
	}
	now := time.Now()

	cartItem := entity.CartItem{
		ID:               itemID,
		CartID:           cart.ID,
		ProductVariantID: req.ProductVariantID,
		Quantity:         req.Quantity,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	return s.cartRepo.UpsertCartItem(ctx, cartItem)
}

func (s *cartService) UpdateCartItem(ctx context.Context, userID, itemID uuid.UUID, req dto.UpdateCartItemRequest) error {
	if err := s.validator.Validate(req); err != nil {
		return err
	}

	cart, err := s.cartRepo.GetOrCreateCartByUserID(ctx, userID)
	if err != nil {
		return err
	}

	item, err := s.cartRepo.GetCartItemByID(ctx, itemID)
	if err != nil {
		return err
	}
	if item == nil || item.CartID != cart.ID {
		return errx.ErrCartItemNotFound
	}

	variant, err := s.productRepo.GetVariantByID(ctx, item.ProductVariantID)
	if err != nil {
		return err
	}
	if !variant.IsActive {
		return errx.ErrVariantInactive
	}

	if req.Quantity > variant.CurrentStock {
		return errx.ErrInsufficientStock
	}

	return s.cartRepo.UpdateCartItemQuantity(ctx, itemID, req.Quantity)
}

func (s *cartService) DeleteCartItem(ctx context.Context, userID, itemID uuid.UUID) error {
	cart, err := s.cartRepo.GetOrCreateCartByUserID(ctx, userID)
	if err != nil {
		return err
	}

	item, err := s.cartRepo.GetCartItemByID(ctx, itemID)
	if err != nil {
		return err
	}
	if item == nil || item.CartID != cart.ID {
		return errx.ErrCartItemNotFound
	}

	return s.cartRepo.DeleteCartItem(ctx, itemID)
}

func (s *cartService) DeleteStoreCartItems(ctx context.Context, userID, storeID uuid.UUID) error {
	cart, err := s.cartRepo.GetOrCreateCartByUserID(ctx, userID)
	if err != nil {
		return err
	}

	return s.cartRepo.DeleteCartItemsByStoreID(ctx, cart.ID, storeID)
}

func (s *cartService) ClearCart(ctx context.Context, userID uuid.UUID) error {
	cart, err := s.cartRepo.GetOrCreateCartByUserID(ctx, userID)
	if err != nil {
		return err
	}

	return s.cartRepo.ClearCart(ctx, cart.ID)
}
