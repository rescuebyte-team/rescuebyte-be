package controller

import (
	"github.com/gofiber/fiber/v3"

	"github.com/rescuebyte/rescuebyte-be/domain/contracts"
	"github.com/rescuebyte/rescuebyte-be/domain/enums"
	"github.com/rescuebyte/rescuebyte-be/internal/middlewares"
)

type userController struct {
	userService contracts.UserService
}

func InitUserController(
	router fiber.Router,
	userService contracts.UserService,
	middleware *middlewares.Middleware,
) {
	userCtr := &userController{
		userService: userService,
	}

	userGroup := router.Group("/users", middleware.RequireAuth())
	userGroup.Get("/profile", userCtr.GetProfile)
	userGroup.Put("/profile", userCtr.UpdateProfile)
	userGroup.Patch("/role", middleware.RequireRoles(enums.RoleAdminIdx), userCtr.UpdateUserRole)
}
