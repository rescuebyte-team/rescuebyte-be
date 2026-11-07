package errx

import "github.com/gofiber/fiber/v3"

var (
	ErrReviewStoreNotFound = NewError(
		fiber.StatusNotFound,
		"review_store_not_found",
		"The store for this review could not be found.",
	)
	ErrReviewProductNotFound = NewError(
		fiber.StatusNotFound,
		"review_product_not_found",
		"The product for this review could not be found.",
	)
	ErrReviewNotFound = NewError(
		fiber.StatusNotFound,
		"review_not_found",
		"The specified review could not be found.",
	)
	ErrInvalidReviewRating = NewError(
		fiber.StatusBadRequest,
		"invalid_review_rating",
		"Review rating must be between 1 and 5.",
	)
)
