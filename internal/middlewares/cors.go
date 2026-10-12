package middlewares

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func Cors() fiber.Handler {
	config := cors.Config{
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"},
		AllowHeaders: []string{
			"Content-Type",
			"Authorization",
			"X-API-Key",
			"Accept",
			"Origin",
			"X-Requested-With",
			"X-Request-Id",
		},
		ExposeHeaders: []string{"Content-Length", "Content-Disposition", "X-Request-Id"},
	}

	return cors.New(config)
}
