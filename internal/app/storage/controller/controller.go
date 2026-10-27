package controller

import (
	"github.com/gofiber/fiber/v3"

	"github.com/rescuebyte/rescuebyte-be/domain/contracts"
	"github.com/rescuebyte/rescuebyte-be/internal/middlewares"
)

type storageController struct {
	storageService contracts.StorageService
}

func InitStorageController(
	router fiber.Router,
	storageService contracts.StorageService,
	middleware *middlewares.Middleware,
) {
	ctr := &storageController{
		storageService: storageService,
	}

	storageGroup := router.Group("/storage", middleware.RequireAuth())
	storageGroup.Post("/upload", ctr.Upload)
	storageGroup.Post("/upload-multiple", ctr.UploadMultiple)
}
