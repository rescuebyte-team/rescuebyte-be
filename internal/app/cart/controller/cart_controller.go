package controller

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"github.com/rescuebyte/rescuebyte-be/domain/dto"
	"github.com/rescuebyte/rescuebyte-be/domain/errx"
	"github.com/rescuebyte/rescuebyte-be/pkg/helpers/http/response"
	"github.com/rescuebyte/rescuebyte-be/pkg/jwt"
)

type getCartQuery struct {
	Latitude  *float64 `query:"latitude"`
	Longitude *float64 `query:"longitude"`
}

// GetCart godoc
// @Summary Get customer cart
// @Description Retrieve cart items grouped by store with financial calculations
// @Tags Cart
// @Security BearerAuth
// @Produce json
// @Param latitude query number false "User latitude"
// @Param longitude query number false "User longitude"
// @Success 200 {object} response.Response{payload=dto.CartResponse} "Cart details"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/cart [get]
func (c *cartController) GetCart(ctx fiber.Ctx) error {
	claims, ok := ctx.Locals("claims").(jwt.Claims)
	if !ok {
		return errx.ErrUnauthorized
	}

	var query getCartQuery
	_ = ctx.Bind().Query(&query)

	res, err := c.cartService.GetCart(ctx.Context(), claims.ID, query.Latitude, query.Longitude)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, res)
}

// GetCartCount godoc
// @Summary Get cart items count
// @Description Get total count of items/boxes in cart for the floating cart badge
// @Tags Cart
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.Response{payload=dto.CartCountResponse} "Cart count"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/cart/count [get]
func (c *cartController) GetCartCount(ctx fiber.Ctx) error {
	claims, ok := ctx.Locals("claims").(jwt.Claims)
	if !ok {
		return errx.ErrUnauthorized
	}

	res, err := c.cartService.GetCartCount(ctx.Context(), claims.ID)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, res)
}

// AddToCart godoc
// @Summary Add item to cart
// @Description Add a product variant to the customer's cart
// @Tags Cart
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.AddToCartRequest true "Add to Cart Payload"
// @Success 201 {object} response.Response{payload=map[string]string} "Item added successfully"
// @Failure 400 {object} response.ErrorResponse "Bad request or insufficient stock"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Variant not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/cart/items [post]
func (c *cartController) AddToCart(ctx fiber.Ctx) error {
	claims, ok := ctx.Locals("claims").(jwt.Claims)
	if !ok {
		return errx.ErrUnauthorized
	}

	var req dto.AddToCartRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return errx.ErrBadRequest.WithError(err)
	}

	if err := c.cartService.AddToCart(ctx.Context(), claims.ID, req); err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusCreated, map[string]string{
		"message": "Item added to cart successfully",
	})
}

// UpdateCartItem godoc
// @Summary Update cart item quantity
// @Description Change the quantity of an item in the cart
// @Tags Cart
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Cart Item UUID"
// @Param request body dto.UpdateCartItemRequest true "Update Quantity Payload"
// @Success 200 {object} response.Response{payload=map[string]string} "Item updated successfully"
// @Failure 400 {object} response.ErrorResponse "Bad request or insufficient stock"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Cart item not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/cart/items/{id} [patch]
func (c *cartController) UpdateCartItem(ctx fiber.Ctx) error {
	claims, ok := ctx.Locals("claims").(jwt.Claims)
	if !ok {
		return errx.ErrUnauthorized
	}

	idParam := ctx.Params("id")
	itemID, err := uuid.Parse(idParam)
	if err != nil {
		return errx.ErrBadRequest.WithError(err)
	}

	var req dto.UpdateCartItemRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return errx.ErrBadRequest.WithError(err)
	}

	if err := c.cartService.UpdateCartItem(ctx.Context(), claims.ID, itemID, req); err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, map[string]string{
		"message": "Cart item updated successfully",
	})
}

// DeleteCartItem godoc
// @Summary Remove item from cart
// @Description Delete a specific item from the cart
// @Tags Cart
// @Security BearerAuth
// @Produce json
// @Param id path string true "Cart Item UUID"
// @Success 200 {object} response.Response{payload=map[string]string} "Item removed successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid UUID format"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Cart item not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/cart/items/{id} [delete]
func (c *cartController) DeleteCartItem(ctx fiber.Ctx) error {
	claims, ok := ctx.Locals("claims").(jwt.Claims)
	if !ok {
		return errx.ErrUnauthorized
	}

	idParam := ctx.Params("id")
	itemID, err := uuid.Parse(idParam)
	if err != nil {
		return errx.ErrBadRequest.WithError(err)
	}

	if err := c.cartService.DeleteCartItem(ctx.Context(), claims.ID, itemID); err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, map[string]string{
		"message": "Cart item removed successfully",
	})
}

// DeleteStoreCartItems godoc
// @Summary Remove all items from a store
// @Description Delete all cart items belonging to a specific store
// @Tags Cart
// @Security BearerAuth
// @Produce json
// @Param store_id path string true "Store UUID"
// @Success 200 {object} response.Response{payload=map[string]string} "Store items removed successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid UUID format"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/cart/stores/{store_id} [delete]
func (c *cartController) DeleteStoreCartItems(ctx fiber.Ctx) error {
	claims, ok := ctx.Locals("claims").(jwt.Claims)
	if !ok {
		return errx.ErrUnauthorized
	}

	storeIDParam := ctx.Params("store_id")
	storeID, err := uuid.Parse(storeIDParam)
	if err != nil {
		return errx.ErrBadRequest.WithError(err)
	}

	if err := c.cartService.DeleteStoreCartItems(ctx.Context(), claims.ID, storeID); err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, map[string]string{
		"message": "Store items removed from cart successfully",
	})
}

// ClearCart godoc
// @Summary Clear cart
// @Description Remove all items from the customer's cart
// @Tags Cart
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.Response{payload=map[string]string} "Cart cleared successfully"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/cart [delete]
func (c *cartController) ClearCart(ctx fiber.Ctx) error {
	claims, ok := ctx.Locals("claims").(jwt.Claims)
	if !ok {
		return errx.ErrUnauthorized
	}

	if err := c.cartService.ClearCart(ctx.Context(), claims.ID); err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, map[string]string{
		"message": "Cart cleared successfully",
	})
}
