package controller

import (
	"github.com/gofiber/fiber/v3"

	"github.com/rescuebyte/rescuebyte-be/domain/contracts"
	"github.com/rescuebyte/rescuebyte-be/internal/middlewares"
)

type orderController struct {
	orderService contracts.OrderService
}

func InitOrderController(
	router fiber.Router,
	orderService contracts.OrderService,
	middleware *middlewares.Middleware,
) {
	ctr := &orderController{
		orderService: orderService,
	}

	orderGroup := router.Group("/orders", middleware.RequireAuth())

	orderGroup.Post("/checkout", ctr.Checkout)
	orderGroup.Get("", ctr.GetMyOrders)
	orderGroup.Get("/:id", ctr.GetOrderDetail)
	orderGroup.Post("/:id/cancel", ctr.CancelOrder)
	orderGroup.Patch("/:id/status", ctr.UpdateOrderStatus)
}
