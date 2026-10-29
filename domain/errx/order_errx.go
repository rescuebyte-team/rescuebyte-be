package errx

import "github.com/gofiber/fiber/v3"

var (
	ErrOrderNotFound = NewError(
		fiber.StatusNotFound,
		"order_not_found",
		"The specified order could not be found.",
	)
	ErrEmptyOrderItems = NewError(
		fiber.StatusBadRequest,
		"empty_order_items",
		"Checkout must include at least one item.",
	)
	ErrStoreDifferentInOrder = NewError(
		fiber.StatusBadRequest,
		"store_different_in_order",
		"All items in a single checkout must belong to the same store.",
	)
	ErrInvalidOrderStatusTransition = NewError(
		fiber.StatusBadRequest,
		"invalid_order_status_transition",
		"The order cannot transition to the requested status.",
	)
	ErrOrderCannotBeCancelled = NewError(
		fiber.StatusBadRequest,
		"order_cannot_be_cancelled",
		"The order can no longer be cancelled in its current state.",
	)
	ErrInvalidPaymentSignature = NewError(
		fiber.StatusUnauthorized,
		"invalid_payment_signature",
		"Payment notification signature verification failed.",
	)
)
