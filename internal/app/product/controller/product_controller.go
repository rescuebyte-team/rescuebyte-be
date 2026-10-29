package controller

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"github.com/rescuebyte/rescuebyte-be/domain/dto"
	"github.com/rescuebyte/rescuebyte-be/domain/errx"
	"github.com/rescuebyte/rescuebyte-be/pkg/helpers/http/response"
	"github.com/rescuebyte/rescuebyte-be/pkg/jwt"
)

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new food rescue product with variants and optional gallery images
// @Tags Product
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateProductRequest true "Create Product Payload"
// @Success 201 {object} response.Response{payload=dto.ProductResponse} "Created product details"
// @Failure 400 {object} response.ErrorResponse "Bad request or validation error"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Store not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/products [post]
func (c *productController) CreateProduct(ctx fiber.Ctx) error {
	claims, ok := ctx.Locals("claims").(jwt.Claims)
	if !ok {
		return errx.ErrUnauthorized
	}

	var req dto.CreateProductRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return errx.ErrBadRequest.WithError(err)
	}

	res, err := c.productService.CreateProduct(ctx.Context(), claims.ID, req)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusCreated, res)
}

// GetMyProducts godoc
// @Summary List merchant products
// @Description Retrieve a list of products belonging to the merchant's store with filtering and search
// @Tags Product
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param search query string false "Search by product name"
// @Param tab query string false "Filter tab: all, active, low_stock, out_of_stock"
// @Param sort query string false "Sort order: newest, oldest"
// @Success 200 {object} response.Response{payload=[]dto.ProductResponse} "List of products"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Store not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/products/me/list [get]
func (c *productController) GetMyProducts(ctx fiber.Ctx) error {
	claims, ok := ctx.Locals("claims").(jwt.Claims)
	if !ok {
		return errx.ErrUnauthorized
	}

	var query dto.GetProductsQuery
	if err := ctx.Bind().Query(&query); err != nil {
		return errx.ErrBadRequest.WithError(err)
	}

	products, total, err := c.productService.GetMyProducts(ctx.Context(), claims.ID, query)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, map[string]any{
		"items": products,
		"total": total,
		"page":  query.Page,
		"limit": query.Limit,
	})
}

// GetAllProducts godoc
// @Summary Get all products
// @Description Retrieve a list of active products for marketplace card display with optional filtering
// @Tags Product
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param store_id query string false "Filter by store UUID"
// @Param nama query string false "Filter by product name"
// @Param name query string false "Alias filter by product name"
// @Param search query string false "Search by product or store name"
// @Param city query string false "Filter by store city"
// @Param latitude query number false "User latitude for distance calculation"
// @Param longitude query number false "User longitude for distance calculation"
// @Success 200 {object} response.Response{payload=dto.GetProductsResponse} "List of products with pagination metadata"
// @Failure 400 {object} response.ErrorResponse "Bad request or validation error"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/products [get]
func (c *productController) GetAllProducts(ctx fiber.Ctx) error {
	var query dto.GetAllProductsQuery
	if err := ctx.Bind().Query(&query); err != nil {
		return errx.ErrBadRequest.WithError(err)
	}

	res, err := c.productService.GetAllProducts(ctx.Context(), query)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, res)
}

// GetProductDetail godoc
// @Summary Get product detail
// @Description Retrieve product detail by its UUID
// @Tags Product
// @Produce json
// @Param id path string true "Product UUID"
// @Success 200 {object} response.Response{payload=dto.ProductResponse} "Product details"
// @Failure 400 {object} response.ErrorResponse "Invalid UUID format"
// @Failure 404 {object} response.ErrorResponse "Product not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/products/{id} [get]
func (c *productController) GetProductDetail(ctx fiber.Ctx) error {
	idParam := ctx.Params("id")
	productID, err := uuid.Parse(idParam)
	if err != nil {
		return errx.ErrBadRequest.WithError(err)
	}

	res, err := c.productService.GetProductDetail(ctx.Context(), productID)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, res)
}

// UpdateProduct godoc
// @Summary Update product
// @Description Update product information
// @Tags Product
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Product UUID"
// @Param request body dto.UpdateProductRequest true "Update Product Payload"
// @Success 200 {object} response.Response{payload=dto.ProductResponse} "Updated product details"
// @Failure 400 {object} response.ErrorResponse "Bad request or validation error"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden - Store access unauthorized"
// @Failure 404 {object} response.ErrorResponse "Product not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/products/{id} [put]
func (c *productController) UpdateProduct(ctx fiber.Ctx) error {
	claims, ok := ctx.Locals("claims").(jwt.Claims)
	if !ok {
		return errx.ErrUnauthorized
	}

	idParam := ctx.Params("id")
	productID, err := uuid.Parse(idParam)
	if err != nil {
		return errx.ErrBadRequest.WithError(err)
	}

	var req dto.UpdateProductRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return errx.ErrBadRequest.WithError(err)
	}

	res, err := c.productService.UpdateProduct(ctx.Context(), claims.ID, productID, req)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, res)
}

// ToggleProductStatus godoc
// @Summary Toggle product status
// @Description Toggle product active/inactive state
// @Tags Product
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Product UUID"
// @Param request body dto.ToggleProductStatusRequest true "Toggle Status Payload"
// @Success 200 {object} response.Response{payload=dto.ProductResponse} "Updated product details"
// @Failure 400 {object} response.ErrorResponse "Invalid UUID format"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden - Store access unauthorized"
// @Failure 404 {object} response.ErrorResponse "Product not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/products/{id}/status [patch]
func (c *productController) ToggleProductStatus(ctx fiber.Ctx) error {
	claims, ok := ctx.Locals("claims").(jwt.Claims)
	if !ok {
		return errx.ErrUnauthorized
	}

	idParam := ctx.Params("id")
	productID, err := uuid.Parse(idParam)
	if err != nil {
		return errx.ErrBadRequest.WithError(err)
	}

	var req dto.ToggleProductStatusRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return errx.ErrBadRequest.WithError(err)
	}

	res, err := c.productService.ToggleProductStatus(ctx.Context(), claims.ID, productID, req)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, res)
}

// DeleteProduct godoc
// @Summary Delete product
// @Description Delete a product and its associated variants and images
// @Tags Product
// @Security BearerAuth
// @Produce json
// @Param id path string true "Product UUID"
// @Success 200 {object} response.Response{payload=map[string]string} "Product successfully deleted"
// @Failure 400 {object} response.ErrorResponse "Invalid UUID format"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden - Store access unauthorized"
// @Failure 404 {object} response.ErrorResponse "Product not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/products/{id} [delete]
func (c *productController) DeleteProduct(ctx fiber.Ctx) error {
	claims, ok := ctx.Locals("claims").(jwt.Claims)
	if !ok {
		return errx.ErrUnauthorized
	}

	idParam := ctx.Params("id")
	productID, err := uuid.Parse(idParam)
	if err != nil {
		return errx.ErrBadRequest.WithError(err)
	}

	if err := c.productService.DeleteProduct(ctx.Context(), claims.ID, productID); err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, map[string]string{
		"message": "Product successfully deleted",
	})
}
