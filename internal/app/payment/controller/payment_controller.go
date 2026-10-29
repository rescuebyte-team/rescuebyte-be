package controller

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"github.com/rescuebyte/rescuebyte-be/domain/dto"
	"github.com/rescuebyte/rescuebyte-be/domain/errx"
	"github.com/rescuebyte/rescuebyte-be/pkg/helpers/http/response"
	"github.com/rescuebyte/rescuebyte-be/pkg/jwt"
)

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
