package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v3"

	"github.com/rescuebyte/rescuebyte-be/pkg/log"
)

func LoggerConfig() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		latency := time.Since(start)

		status := c.Response().StatusCode()
		fields := log.Fields{
			"status":     status,
			"method":     c.Method(),
			"path":       c.Path(),
			"latency":    latency.String(),
			"ip":         c.IP(),
			"request_id": c.GetRespHeader("X-Request-Id"),
		}

		if err != nil {
			fields["error"] = err.Error()
		}

		if status >= 500 {
			log.Error(fields, "[HTTP] Server error")
		} else if status >= 400 {
			log.Warn(fields, "[HTTP] Client error")
		} else {
			log.Info(fields, "[HTTP] Request handled")
		}

		return err
	}
}
