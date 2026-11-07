package errx

import "github.com/gofiber/fiber/v3"

var (
	ErrCommentStoreNotFound = NewError(
		fiber.StatusNotFound,
		"comment_store_not_found",
		"The store for this comment could not be found.",
	)
	ErrCommentProductNotFound = NewError(
		fiber.StatusNotFound,
		"comment_product_not_found",
		"The product for this comment could not be found.",
	)
	ErrCommentNotFound = NewError(
		fiber.StatusNotFound,
		"comment_not_found",
		"The specified comment could not be found.",
	)
	ErrInvalidCommentRating = NewError(
		fiber.StatusBadRequest,
		"invalid_comment_rating",
		"Comment rating must be between 1 and 5.",
	)
)
