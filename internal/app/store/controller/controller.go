package controller

import (
	"github.com/gofiber/fiber/v3"

	"github.com/rescuebyte/rescuebyte-be/domain/contracts"
	"github.com/rescuebyte/rescuebyte-be/internal/middlewares"
)

type storeController struct {
	storeService contracts.StoreService
}

func InitStoreController(
	router fiber.Router,
	storeService contracts.StoreService,
	middleware *middlewares.Middleware,
) {
	ctr := &storeController{
		storeService: storeService,
	}

	storeGroup := router.Group("/stores")
	storeGroup.Get("", ctr.GetAllStores)
	storeGroup.Get("/:id", ctr.GetStoreByID)

	storeGroup.Post("/", middleware.RequireAuth(), ctr.CreateStore)
	storeGroup.Get("/me/profile", middleware.RequireAuth(), ctr.GetMyStore)
	storeGroup.Put("/me/profile", middleware.RequireAuth(), ctr.UpdateStoreProfile)
	storeGroup.Put("/me/operational", middleware.RequireAuth(), ctr.UpdateStoreOperational)
	storeGroup.Patch("/me/status", middleware.RequireAuth(), ctr.ToggleStoreStatus)
}
