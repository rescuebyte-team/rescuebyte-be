package response

import (
	"errors"

	"github.com/gofiber/fiber/v3"

	"github.com/rescuebyte/rescuebyte-be/domain/errx"
	"github.com/rescuebyte/rescuebyte-be/internal/infra/env"
)

type Response struct {
	RequestID string      `json:"request_id,omitempty"`
	Payload   interface{} `json:"payload"`
}

func SendResponse(ctx fiber.Ctx, code int, payload any) error {
	if code >= 400 {
		payload = formatErrorPayload(payload)
	}

	return ctx.Status(code).JSON(
		Response{
			RequestID: ctx.GetRespHeader("X-Request-Id"),
			Payload:   payload,
		},
	)
}

func formatErrorPayload(payload any) any {
	err, isError := payload.(error)
	if !isError {
		return payload
	}

	var reqErr *errx.RequestError
	if !errors.As(err, &reqErr) {
		return map[string]any{
			"error": map[string]any{
				"message": "unknown error",
				"error":   err.Error(),
			},
		}
	}

	if env.AppEnv.AppEnv != "production" {
		return map[string]any{
			"error": map[string]any{
				"message":    reqErr.Message,
				"error":      errorString(reqErr.Err),
				"error_code": reqErr.ErrorCode,
				"location":   reqErr.Location,
				"details":    reqErr.Details,
			},
		}
	}

	return map[string]any{
		"error": reqErr,
	}
}

func errorString(err error) any {
	if err == nil {
		return nil
	}
	return err.Error()
}
