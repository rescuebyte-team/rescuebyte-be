package controller

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"github.com/rescuebyte/rescuebyte-be/domain/dto"
	"github.com/rescuebyte/rescuebyte-be/domain/errx"
	"github.com/rescuebyte/rescuebyte-be/pkg/helpers/http/response"
	"github.com/rescuebyte/rescuebyte-be/pkg/jwt"
)

func (c *orderController) Checkout(ctx fiber.Ctx) error {
	claims, ok := ctx.Locals("claims").(jwt.Claims)
	if !ok {
		return errx.ErrUnauthorized
	}

	var req dto.CheckoutRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return errx.ErrBadRequest.WithError(err)
	}

	res, err := c.orderService.Checkout(ctx.Context(), claims.ID, req)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusCreated, res)
}

func (c *orderController) GetMyOrders(ctx fiber.Ctx) error {
	claims, ok := ctx.Locals("claims").(jwt.Claims)
	if !ok {
		return errx.ErrUnauthorized
	}

	var query dto.GetOrdersQuery
	if err := ctx.Bind().Query(&query); err != nil {
		return errx.ErrBadRequest.WithError(err)
	}

	res, err := c.orderService.GetMyOrders(ctx.Context(), claims.ID, query)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, res)
}

func (c *orderController) GetOrderDetail(ctx fiber.Ctx) error {
	claims, ok := ctx.Locals("claims").(jwt.Claims)
	if !ok {
		return errx.ErrUnauthorized
	}

	idParam := ctx.Params("id")
	orderID, err := uuid.Parse(idParam)
	if err != nil {
		return errx.ErrBadRequest.WithError(err)
	}

	res, err := c.orderService.GetOrderDetail(ctx.Context(), claims.ID, orderID)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, res)
}

func (c *orderController) CancelOrder(ctx fiber.Ctx) error {
	claims, ok := ctx.Locals("claims").(jwt.Claims)
	if !ok {
		return errx.ErrUnauthorized
	}

	idParam := ctx.Params("id")
	orderID, err := uuid.Parse(idParam)
	if err != nil {
		return errx.ErrBadRequest.WithError(err)
	}

	var req dto.CancelOrderRequest
	_ = ctx.Bind().Body(&req)

	if err := c.orderService.CancelOrder(ctx.Context(), claims.ID, orderID, req); err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, map[string]string{
		"message": "Order cancelled successfully",
	})
}

func (c *orderController) UpdateOrderStatus(ctx fiber.Ctx) error {
	claims, ok := ctx.Locals("claims").(jwt.Claims)
	if !ok {
		return errx.ErrUnauthorized
	}

	idParam := ctx.Params("id")
	orderID, err := uuid.Parse(idParam)
	if err != nil {
		return errx.ErrBadRequest.WithError(err)
	}

	var req dto.UpdateOrderStatusRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return errx.ErrBadRequest.WithError(err)
	}

	res, err := c.orderService.UpdateOrderStatusByMerchant(ctx.Context(), claims.ID, orderID, req)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, res)
}
