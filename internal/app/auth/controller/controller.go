package controller

import (
	"github.com/gofiber/fiber/v3"

	"github.com/rescuebyte/rescuebyte-be/domain/contracts"
	"github.com/rescuebyte/rescuebyte-be/internal/middlewares"
)

type authController struct {
	authService contracts.AuthService
}

func InitAuthController(
	router fiber.Router,
	authService contracts.AuthService,
	middleware *middlewares.Middleware,
) {
	authCtr := &authController{
		authService: authService,
	}

	authGroup := router.Group("/auth")
	authGroup.Post("/register", authCtr.Register)
	authGroup.Post("/login", authCtr.Login)
	authGroup.Post("/refresh", authCtr.RefreshToken)
	authGroup.Post("/logout", authCtr.Logout)
	authGroup.Get("/google", authCtr.GetGoogleAuthURL)
	authGroup.Get("/google/callback", authCtr.GoogleCallback)
	authGroup.Get("/me", middleware.RequireAuth(), authCtr.GetMe)
}
