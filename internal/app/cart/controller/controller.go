package controller

import (
	"github.com/gofiber/fiber/v3"

	"github.com/rescuebyte/rescuebyte-be/domain/contracts"
	"github.com/rescuebyte/rescuebyte-be/internal/middlewares"
)

type cartController struct {
	cartService contracts.CartService
}

func InitCartController(
	router fiber.Router,
	cartService contracts.CartService,
	middleware *middlewares.Middleware,
) {
	ctr := &cartController{
		cartService: cartService,
	}

	cartGroup := router.Group("/cart", middleware.RequireAuth())

	cartGroup.Get("", ctr.GetCart)
	cartGroup.Get("/count", ctr.GetCartCount)
	cartGroup.Post("/items", ctr.AddToCart)
	cartGroup.Patch("/items/:id", ctr.UpdateCartItem)
	cartGroup.Delete("/items/:id", ctr.DeleteCartItem)
	cartGroup.Delete("/stores/:store_id", ctr.DeleteStoreCartItems)
	cartGroup.Delete("", ctr.ClearCart)
}
