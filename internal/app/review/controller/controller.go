package controller

import (
	"github.com/gofiber/fiber/v3"

	"github.com/rescuebyte/rescuebyte-be/domain/contracts"
	"github.com/rescuebyte/rescuebyte-be/internal/middlewares"
)

type reviewController struct {
	reviewService contracts.ReviewService
}

func InitReviewController(router fiber.Router, reviewService contracts.ReviewService, middleware *middlewares.Middleware) {
	ctr := &reviewController{reviewService: reviewService}

	storeReviews := router.Group("/stores/:store_id/reviews")
	storeReviews.Get("", ctr.GetStoreReviews)
	storeReviews.Post("", middleware.RequireAuth(), ctr.CreateStoreReview)

	productReviews := router.Group("/products/:product_id/reviews")
	productReviews.Get("", ctr.GetProductReviews)
	productReviews.Post("", middleware.RequireAuth(), ctr.CreateProductReview)

	reviewsGroup := router.Group("/reviews")
	reviewsGroup.Patch("/:id", middleware.RequireAuth(), ctr.UpdateReview)
	reviewsGroup.Delete("/:id", middleware.RequireAuth(), ctr.DeleteReview)
}
