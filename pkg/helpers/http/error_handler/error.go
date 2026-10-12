package errorhandler

import (
	"errors"

	"github.com/gofiber/fiber/v3"

	"github.com/rescuebyte/rescuebyte-be/domain/errx"
	"github.com/rescuebyte/rescuebyte-be/pkg/helpers/http/response"
	"github.com/rescuebyte/rescuebyte-be/pkg/log"
	"github.com/rescuebyte/rescuebyte-be/pkg/validator"
)

func ErrorHandler(c fiber.Ctx, err error) error {
	var valErr validator.MultiError
	if errors.As(err, &valErr) {
		return response.SendResponse(c, fiber.StatusUnprocessableEntity, map[string]any{
			"message":    "Validation error",
			"error":      valErr,
			"error_code": "VALIDATION_ERROR",
		})
	}

	var reqErr *errx.RequestError
	if errors.As(err, &reqErr) {
		log.Error(log.Fields{
			"request_id": c.GetRespHeader("X-Request-Id"),
			"error_code": reqErr.ErrorCode,
			"location":   reqErr.Location,
			"details":    reqErr.Details,
			"error":      reqErr.Err,
		}, "[ErrorHandler] Request error")

		return response.SendResponse(c, reqErr.StatusCode, reqErr)
	}

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		log.Error(log.Fields{
			"request_id": c.GetRespHeader("X-Request-Id"),
			"error":      fiberErr.Error(),
		}, "[ErrorHandler] Fiber error")

		return response.SendResponse(c, fiberErr.Code, fiberErr)
	}

	log.Error(log.Fields{
		"request_id": c.GetRespHeader("X-Request-Id"),
		"error":      err.Error(),
	}, "[ErrorHandler] Internal server error")

	return response.SendResponse(c, fiber.StatusInternalServerError, errx.ErrInternalServer.WithError(err))
}
