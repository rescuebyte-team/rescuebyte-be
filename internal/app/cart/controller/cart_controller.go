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
