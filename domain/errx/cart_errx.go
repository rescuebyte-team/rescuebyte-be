package errx

import "github.com/gofiber/fiber/v3"

var (
	ErrCartItemNotFound = NewError(
		fiber.StatusNotFound,
		"cart_item_not_found",
		"The specified cart item was not found.",
	)
	ErrInsufficientStock = NewError(
		fiber.StatusBadRequest,
		"insufficient_stock",
		"The requested quantity exceeds available stock.",
	)
	ErrVariantInactive = NewError(
		fiber.StatusBadRequest,
		"variant_inactive",
		"This product variant is currently not available for purchase.",
	)
)
