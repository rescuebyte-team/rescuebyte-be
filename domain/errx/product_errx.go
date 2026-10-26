package errx

import "github.com/gofiber/fiber/v3"

var (
	ErrProductNotFound = NewError(
		fiber.StatusNotFound,
		"product_not_found",
		"The specified product could not be found.",
	)
	ErrVariantNotFound = NewError(
		fiber.StatusNotFound,
		"variant_not_found",
		"The specified product variant could not be found.",
	)
	ErrUnauthorizedStoreAccess = NewError(
		fiber.StatusForbidden,
		"unauthorized_store_access",
		"You do not have permission to access or modify this store's products.",
	)
	ErrProductMustHaveVariant = NewError(
		fiber.StatusBadRequest,
		"product_must_have_variant",
		"A product must contain at least one variant.",
	)
)
