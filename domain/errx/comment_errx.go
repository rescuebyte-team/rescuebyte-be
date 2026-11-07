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
)
