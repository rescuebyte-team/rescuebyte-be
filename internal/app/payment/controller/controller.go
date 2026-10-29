package controller

import (
	"github.com/gofiber/fiber/v3"

	"github.com/rescuebyte/rescuebyte-be/domain/contracts"
	"github.com/rescuebyte/rescuebyte-be/internal/middlewares"
)

type paymentController struct {
	paymentService contracts.PaymentService
}

func InitPaymentController(
	router fiber.Router,
	paymentService contracts.PaymentService,
	middleware *middlewares.Middleware,
) {
	ctr := &paymentController{
		paymentService: paymentService,
	}

	paymentGroup := router.Group("/payments")
	paymentGroup.Post("/midtrans/webhook", ctr.HandleMidtransWebhook)

	protected := paymentGroup.Group("", middleware.RequireAuth())
	protected.Get("/order/:order_id", ctr.GetPaymentStatus)
}
