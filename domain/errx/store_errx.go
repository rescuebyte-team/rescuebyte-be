package errx

import "github.com/gofiber/fiber/v3"

var (
	ErrStoreNotFound = NewError(
		fiber.StatusNotFound,
		"store_not_found",
		"The specified store could not be found.",
	)
	ErrStoreAlreadyExists = NewError(
		fiber.StatusConflict,
		"store_already_exists",
		"You already have a registered store.",
	)
	ErrStoreInactive = NewError(
		fiber.StatusBadRequest,
		"store_inactive",
		"The store is currently inactive.",
	)
)
