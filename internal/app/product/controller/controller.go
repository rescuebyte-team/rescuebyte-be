package controller

import (
	"github.com/gofiber/fiber/v3"

	"github.com/rescuebyte/rescuebyte-be/domain/contracts"
	"github.com/rescuebyte/rescuebyte-be/internal/middlewares"
)

type productController struct {
	productService contracts.ProductService
}

func InitProductController(
	router fiber.Router,
	productService contracts.ProductService,
	middleware *middlewares.Middleware,
) {
	ctr := &productController{
		productService: productService,
	}

	productGroup := router.Group("/products")

	// Public product detail
	productGroup.Get("/:id", ctr.GetProductDetail)

	// Protected merchant product operations
	protected := productGroup.Group("", middleware.RequireAuth())
	protected.Post("/", ctr.CreateProduct)
	protected.Get("/me/list", ctr.GetMyProducts)
	protected.Put("/:id", ctr.UpdateProduct)
	protected.Patch("/:id/status", ctr.ToggleProductStatus)
	protected.Delete("/:id", ctr.DeleteProduct)
}
