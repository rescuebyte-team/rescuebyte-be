package controller

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"github.com/rescuebyte/rescuebyte-be/domain/dto"
	"github.com/rescuebyte/rescuebyte-be/domain/errx"
	"github.com/rescuebyte/rescuebyte-be/pkg/helpers/http/response"
	"github.com/rescuebyte/rescuebyte-be/pkg/jwt"
)

// HandleMidtransWebhook godoc
// @Summary Midtrans payment webhook
// @Description Handle Midtrans payment notification with SHA-512 signature validation and automatic status transition
// @Tags Payment
// @Accept json
// @Produce json
// @Param payload body dto.MidtransNotificationPayload true "Midtrans Notification Payload"
// @Success 200 {object} response.Response{payload=map[string]string} "Payment notification processed successfully"
// @Failure 400 {object} response.ErrorResponse "Bad request"
// @Failure 401 {object} response.ErrorResponse "Invalid signature"
// @Failure 404 {object} response.ErrorResponse "Order not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/payments/midtrans/webhook [post]
func (c *paymentController) HandleMidtransWebhook(ctx fiber.Ctx) error {
	var payload dto.MidtransNotificationPayload
	if err := ctx.Bind().Body(&payload); err != nil {
		return errx.ErrBadRequest.WithError(err)
	}

	if err := c.paymentService.HandleMidtransNotification(ctx.Context(), payload); err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, map[string]string{
		"message": "Payment notification processed successfully",
	})
}

// GetPaymentStatus godoc
// @Summary Check payment status
// @Description Check payment status for a specific order
// @Tags Payment
// @Security BearerAuth
// @Produce json
// @Param order_id path string true "Order UUID"
// @Success 200 {object} response.Response{payload=dto.PaymentResponse} "Payment status details"
// @Failure 400 {object} response.ErrorResponse "Invalid UUID format"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden"
// @Failure 404 {object} response.ErrorResponse "Payment not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/payments/order/{order_id} [get]
func (c *paymentController) GetPaymentStatus(ctx fiber.Ctx) error {
	claims, ok := ctx.Locals("claims").(jwt.Claims)
	if !ok {
		return errx.ErrUnauthorized
	}

	orderIDParam := ctx.Params("order_id")
	orderID, err := uuid.Parse(orderIDParam)
	if err != nil {
		return errx.ErrBadRequest.WithError(err)
	}

	res, err := c.paymentService.GetPaymentStatus(ctx.Context(), claims.ID, orderID)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, res)
}
