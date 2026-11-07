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

	productGroup.Get("", ctr.GetAllProducts)
	productGroup.Get("/:id", ctr.GetProductDetail)

	productGroup.Post("/", middleware.RequireAuth(), ctr.CreateProduct)
	productGroup.Get("/me/list", middleware.RequireAuth(), ctr.GetMyProducts)
	productGroup.Put("/:id", middleware.RequireAuth(), ctr.UpdateProduct)
	productGroup.Patch("/:id/status", middleware.RequireAuth(), ctr.ToggleProductStatus)
	productGroup.Delete("/:id", middleware.RequireAuth(), ctr.DeleteProduct)
}
