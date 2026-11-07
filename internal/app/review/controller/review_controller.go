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

// CreateStoreReview godoc
// @Summary Create a store review
// @Description Create a rating and review for a store as the authenticated user
// @Tags Review
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param store_id path string true "Store UUID"
// @Param request body dto.CreateReviewRequest true "Review and rating"
// @Success 201 {object} response.Response{payload=dto.ReviewResponse}
// @Failure 400 {object} response.ErrorResponse "Invalid ID or request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Store not found"
// @Router /api/v1/stores/{store_id}/reviews [post]
func (c *reviewController) CreateStoreReview(ctx fiber.Ctx) error {
	claims, ok := ctx.Locals("claims").(jwt.Claims)
	if !ok {
		return errx.ErrUnauthorized
	}
	storeID, err := parseTargetID(ctx, "store_id")
	if err != nil {
		return err
	}
	var req dto.CreateReviewRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return errx.ErrBadRequest.WithError(err)
	}
	result, err := c.reviewService.CreateStoreReview(ctx.Context(), claims.ID, storeID, req)
	if err != nil {
		return err
	}
	return response.SendResponse(ctx, fiber.StatusCreated, result)
}

// GetStoreReviews godoc
// @Summary List store reviews
// @Description List newest reviews for a store, including author details
// @Tags Review
// @Produce json
// @Param store_id path string true "Store UUID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} response.Response{payload=dto.GetReviewsResponse}
// @Failure 400 {object} response.ErrorResponse "Invalid ID or query"
// @Failure 404 {object} response.ErrorResponse "Store not found"
// @Router /api/v1/stores/{store_id}/reviews [get]
func (c *reviewController) GetStoreReviews(ctx fiber.Ctx) error {
	storeID, err := parseTargetID(ctx, "store_id")
	if err != nil {
		return err
	}
	var query dto.GetReviewsQuery
	if err := ctx.Bind().Query(&query); err != nil {
		return errx.ErrBadRequest.WithError(err)
	}
	result, err := c.reviewService.GetStoreReviews(ctx.Context(), storeID, query)
	if err != nil {
		return err
	}
	return response.SendResponse(ctx, fiber.StatusOK, result)
}

// CreateProductReview godoc
// @Summary Create a product review
// @Description Create a rating and review for a product as the authenticated user
// @Tags Review
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param product_id path string true "Product UUID"
// @Param request body dto.CreateReviewRequest true "Review and rating"
// @Success 201 {object} response.Response{payload=dto.ReviewResponse}
// @Failure 400 {object} response.ErrorResponse "Invalid ID or request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Product not found"
// @Router /api/v1/products/{product_id}/reviews [post]
func (c *reviewController) CreateProductReview(ctx fiber.Ctx) error {
	claims, ok := ctx.Locals("claims").(jwt.Claims)
	if !ok {
		return errx.ErrUnauthorized
	}
	productID, err := parseTargetID(ctx, "product_id")
	if err != nil {
		return err
	}
	var req dto.CreateReviewRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return errx.ErrBadRequest.WithError(err)
	}
	result, err := c.reviewService.CreateProductReview(ctx.Context(), claims.ID, productID, req)
	if err != nil {
		return err
	}
	return response.SendResponse(ctx, fiber.StatusCreated, result)
}

// GetProductReviews godoc
// @Summary List product reviews
// @Description List newest reviews for a product, including author details
// @Tags Review
// @Produce json
// @Param product_id path string true "Product UUID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} response.Response{payload=dto.GetReviewsResponse}
// @Failure 400 {object} response.ErrorResponse "Invalid ID or query"
// @Failure 404 {object} response.ErrorResponse "Product not found"
// @Router /api/v1/products/{product_id}/reviews [get]
func (c *reviewController) GetProductReviews(ctx fiber.Ctx) error {
	productID, err := parseTargetID(ctx, "product_id")
	if err != nil {
		return err
	}
	var query dto.GetReviewsQuery
	if err := ctx.Bind().Query(&query); err != nil {
		return errx.ErrBadRequest.WithError(err)
	}
	result, err := c.reviewService.GetProductReviews(ctx.Context(), productID, query)
	if err != nil {
		return err
	}
	return response.SendResponse(ctx, fiber.StatusOK, result)
}

// UpdateReview godoc
// @Summary Edit a review
// @Description Edit a review's content and rating. Only its author or an admin can edit it.
// @Tags Review
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Review UUID"
// @Param request body dto.UpdateReviewRequest true "Updated review and rating"
// @Success 200 {object} response.Response{payload=dto.ReviewResponse}
// @Failure 400 {object} response.ErrorResponse "Invalid ID or request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Only the review author or an admin can edit this review"
// @Failure 404 {object} response.ErrorResponse "Review not found"
// @Router /api/v1/reviews/{id} [patch]
func (c *reviewController) UpdateReview(ctx fiber.Ctx) error {
	claims, ok := ctx.Locals("claims").(jwt.Claims)
	if !ok {
		return errx.ErrUnauthorized
	}
	reviewID, err := parseTargetID(ctx, "id")
	if err != nil {
		return err
	}
	var req dto.UpdateReviewRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return errx.ErrBadRequest.WithError(err)
	}
	result, err := c.reviewService.UpdateReview(ctx.Context(), claims.ID, claims.Role, reviewID, req)
	if err != nil {
		return err
	}
	return response.SendResponse(ctx, fiber.StatusOK, result)
}

// DeleteReview godoc
// @Summary Delete a review
// @Description Delete a review. Only its author or an admin can delete it.
// @Tags Review
// @Security BearerAuth
// @Produce json
// @Param id path string true "Review UUID"
// @Success 200 {object} response.Response{payload=map[string]string}
// @Failure 400 {object} response.ErrorResponse "Invalid ID"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Only the review author or an admin can delete this review"
// @Failure 404 {object} response.ErrorResponse "Review not found"
// @Router /api/v1/reviews/{id} [delete]
func (c *reviewController) DeleteReview(ctx fiber.Ctx) error {
	claims, ok := ctx.Locals("claims").(jwt.Claims)
	if !ok {
		return errx.ErrUnauthorized
	}
	reviewID, err := parseTargetID(ctx, "id")
	if err != nil {
		return err
	}
	if err := c.reviewService.DeleteReview(ctx.Context(), claims.ID, claims.Role, reviewID); err != nil {
		return err
	}
	return response.SendResponse(ctx, fiber.StatusOK, map[string]string{"message": "Review deleted successfully"})
}
