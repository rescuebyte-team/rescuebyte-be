package middlewares

import (
	"context"

	"github.com/gofiber/fiber/v3"

	"github.com/rescuebyte/rescuebyte-be/pkg/log"
)

func RequestIDContext() fiber.Handler {
	return func(c fiber.Ctx) error {
		reqID := c.GetRespHeader("X-Request-Id")
		if reqID != "" {
			ctx := context.WithValue(c.Context(), log.RequestIDKey, reqID)
			c.SetContext(ctx)
		}
		return c.Next()
	}
}
