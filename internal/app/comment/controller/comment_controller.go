package controller

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"github.com/rescuebyte/rescuebyte-be/domain/dto"
	"github.com/rescuebyte/rescuebyte-be/domain/errx"
	"github.com/rescuebyte/rescuebyte-be/pkg/helpers/http/response"
	"github.com/rescuebyte/rescuebyte-be/pkg/jwt"
)

func parseTargetID(ctx fiber.Ctx, param string) (uuid.UUID, error) {
	id, err := uuid.Parse(ctx.Params(param))
	if err != nil {
		return uuid.Nil, errx.ErrBadRequest.WithMessage("Invalid target ID.").WithError(err).WithLocation("path." + param)
	}
	return id, nil
}

// CreateStoreComment godoc
// @Summary Create a store comment
// @Description Create a rating and comment for a store as the authenticated user
// @Tags Comment
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param store_id path string true "Store UUID"
// @Param request body dto.CreateCommentRequest true "Comment and rating"
// @Success 201 {object} response.Response{payload=dto.CommentResponse}
// @Failure 400 {object} response.ErrorResponse "Invalid ID or request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Store not found"
// @Router /api/v1/stores/{store_id}/comments [post]
func (c *commentController) CreateStoreComment(ctx fiber.Ctx) error {
	claims, ok := ctx.Locals("claims").(jwt.Claims)
	if !ok {
		return errx.ErrUnauthorized
	}
	storeID, err := parseTargetID(ctx, "store_id")
	if err != nil {
		return err
	}
	var req dto.CreateCommentRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return errx.ErrBadRequest.WithError(err)
	}
	result, err := c.commentService.CreateStoreComment(ctx.Context(), claims.ID, storeID, req)
	if err != nil {
		return err
	}
	return response.SendResponse(ctx, fiber.StatusCreated, result)
}

// GetStoreComments godoc
// @Summary List store comments
// @Description List newest comments for a store, including author details
// @Tags Comment
// @Produce json
// @Param store_id path string true "Store UUID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} response.Response{payload=dto.GetCommentsResponse}
// @Failure 400 {object} response.ErrorResponse "Invalid ID or query"
// @Failure 404 {object} response.ErrorResponse "Store not found"
// @Router /api/v1/stores/{store_id}/comments [get]
func (c *commentController) GetStoreComments(ctx fiber.Ctx) error {
	storeID, err := parseTargetID(ctx, "store_id")
	if err != nil {
		return err
	}
	var query dto.GetCommentsQuery
	if err := ctx.Bind().Query(&query); err != nil {
		return errx.ErrBadRequest.WithError(err)
	}
	result, err := c.commentService.GetStoreComments(ctx.Context(), storeID, query)
	if err != nil {
		return err
	}
	return response.SendResponse(ctx, fiber.StatusOK, result)
}

// CreateProductComment godoc
// @Summary Create a product comment
// @Description Create a rating and comment for a product as the authenticated user
// @Tags Comment
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param product_id path string true "Product UUID"
// @Param request body dto.CreateCommentRequest true "Comment and rating"
// @Success 201 {object} response.Response{payload=dto.CommentResponse}
// @Failure 400 {object} response.ErrorResponse "Invalid ID or request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Product not found"
// @Router /api/v1/products/{product_id}/comments [post]
func (c *commentController) CreateProductComment(ctx fiber.Ctx) error {
	claims, ok := ctx.Locals("claims").(jwt.Claims)
	if !ok {
		return errx.ErrUnauthorized
	}
	productID, err := parseTargetID(ctx, "product_id")
	if err != nil {
		return err
	}
	var req dto.CreateCommentRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return errx.ErrBadRequest.WithError(err)
	}
	result, err := c.commentService.CreateProductComment(ctx.Context(), claims.ID, productID, req)
	if err != nil {
		return err
	}
	return response.SendResponse(ctx, fiber.StatusCreated, result)
}

// GetProductComments godoc
// @Summary List product comments
// @Description List newest comments for a product, including author details
// @Tags Comment
// @Produce json
// @Param product_id path string true "Product UUID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} response.Response{payload=dto.GetCommentsResponse}
// @Failure 400 {object} response.ErrorResponse "Invalid ID or query"
// @Failure 404 {object} response.ErrorResponse "Product not found"
// @Router /api/v1/products/{product_id}/comments [get]
func (c *commentController) GetProductComments(ctx fiber.Ctx) error {
	productID, err := parseTargetID(ctx, "product_id")
	if err != nil {
		return err
	}
	var query dto.GetCommentsQuery
	if err := ctx.Bind().Query(&query); err != nil {
		return errx.ErrBadRequest.WithError(err)
	}
	result, err := c.commentService.GetProductComments(ctx.Context(), productID, query)
	if err != nil {
		return err
	}
	return response.SendResponse(ctx, fiber.StatusOK, result)
}
