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
	// Public store endpoints
	storeGroup.Get("", ctr.GetAllStores)
	storeGroup.Get("/:id", ctr.GetStoreByID)

	// Protected merchant store endpoints
	protected := storeGroup.Group("", middleware.RequireAuth())
	protected.Post("/", ctr.CreateStore)
	protected.Get("/me/profile", ctr.GetMyStore)
	protected.Put("/me/profile", ctr.UpdateStoreProfile)
	protected.Put("/me/operational", ctr.UpdateStoreOperational)
	protected.Patch("/me/status", ctr.ToggleStoreStatus)
}
