package controller

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"github.com/rescuebyte/rescuebyte-be/domain/dto"
	"github.com/rescuebyte/rescuebyte-be/domain/errx"
	"github.com/rescuebyte/rescuebyte-be/pkg/helpers/http/response"
	"github.com/rescuebyte/rescuebyte-be/pkg/jwt"
)

// Checkout godoc
// @Summary Checkout order
// @Description Create an order for a store with online payment (Midtrans Snap) or cash on pickup
// @Tags Order
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CheckoutRequest true "Checkout Payload"
// @Success 201 {object} response.Response{payload=dto.CheckoutResponse} "Order created with optional snap token"
// @Failure 400 {object} response.ErrorResponse "Bad request or insufficient stock"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Store or variant not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/orders/checkout [post]
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

// GetMyOrders godoc
// @Summary List customer orders
// @Description Retrieve a list of orders for the customer with tab status filtering and search
// @Tags Order
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param tab query string false "Status tab filter: all, waiting_payment, waiting_confirmation, ready_for_pickup, completed, cancelled"
// @Param search query string false "Search by order number, store name, or product name"
// @Success 200 {object} response.Response{payload=dto.GetOrdersResponse} "List of orders with pagination metadata"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/orders [get]
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

// GetOrderDetail godoc
// @Summary Get order detail
// @Description Retrieve complete order details including items, payment info, and 5-step status tracking
// @Tags Order
// @Security BearerAuth
// @Produce json
// @Param id path string true "Order UUID"
// @Success 200 {object} response.Response{payload=dto.OrderDetailResponse} "Order detail"
// @Failure 400 {object} response.ErrorResponse "Invalid UUID format"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden"
// @Failure 404 {object} response.ErrorResponse "Order not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/orders/{id} [get]
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

// CancelOrder godoc
// @Summary Cancel order
// @Description Cancel an order before it is processed and automatically restore variant stocks
// @Tags Order
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Order UUID"
// @Param request body dto.CancelOrderRequest false "Cancel Order Payload"
// @Success 200 {object} response.Response{payload=map[string]string} "Order cancelled successfully"
// @Failure 400 {object} response.ErrorResponse "Order cannot be cancelled in its current state"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden"
// @Failure 404 {object} response.ErrorResponse "Order not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/orders/{id}/cancel [post]
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

// UpdateOrderStatus godoc
// @Summary Update order status by merchant
// @Description Transition order status (IN_PROCESS, READY_FOR_PICKUP, COMPLETED, CANCELLED)
// @Tags Order
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Order UUID"
// @Param request body dto.UpdateOrderStatusRequest true "Update Status Payload"
// @Success 200 {object} response.Response{payload=dto.OrderDetailResponse} "Updated order detail"
// @Failure 400 {object} response.ErrorResponse "Bad request or invalid status"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden - Store access unauthorized"
// @Failure 404 {object} response.ErrorResponse "Order not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/orders/{id}/status [patch]
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
