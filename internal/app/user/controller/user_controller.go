package controller

import (
	"github.com/gofiber/fiber/v3"

	"github.com/rescuebyte/rescuebyte-be/domain/dto"
	"github.com/rescuebyte/rescuebyte-be/domain/errx"
	"github.com/rescuebyte/rescuebyte-be/pkg/helpers/http/response"
	"github.com/rescuebyte/rescuebyte-be/pkg/jwt"
)

// GetProfile godoc
// @Summary Get user profile
// @Description Retrieve the profile of the currently authenticated user
// @Tags User
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.Response{payload=dto.UserResponse} "User profile details"
// @Failure 401 {object} response.ErrorResponse "Unauthorized or missing token"
// @Failure 404 {object} response.ErrorResponse "User not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/users/profile [get]
func (c *userController) GetProfile(ctx fiber.Ctx) error {
	claims, ok := ctx.Locals("claims").(jwt.Claims)
	if !ok {
		return errx.ErrUnauthorized
	}

	res, err := c.userService.GetProfile(ctx.Context(), claims.ID)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, res)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update the authenticated user's name and/or avatar URL
// @Tags User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.UpdateProfileRequest true "Profile Update Payload"
// @Success 200 {object} response.Response{payload=dto.UserResponse} "Updated user profile"
// @Failure 400 {object} response.ErrorResponse "Bad request or validation error"
// @Failure 401 {object} response.ErrorResponse "Unauthorized or missing token"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/users/profile [put]
func (c *userController) UpdateProfile(ctx fiber.Ctx) error {
	claims, ok := ctx.Locals("claims").(jwt.Claims)
	if !ok {
		return errx.ErrUnauthorized
	}

	var req dto.UpdateProfileRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return errx.ErrBadRequest.WithError(err)
	}

	res, err := c.userService.UpdateProfile(ctx.Context(), claims.ID, req)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, res)
}

// UpdateUserRole godoc
// @Summary Update user role
// @Description Modify a target user's role (Admin role required)
// @Tags User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.UpdateUserRoleRequest true "User Role Update Payload"
// @Success 200 {object} response.Response{payload=dto.UserResponse} "Updated user with new role"
// @Failure 400 {object} response.ErrorResponse "Bad request or validation error"
// @Failure 401 {object} response.ErrorResponse "Unauthorized or missing token"
// @Failure 403 {object} response.ErrorResponse "Forbidden - Admin access required"
// @Failure 404 {object} response.ErrorResponse "User not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/users/role [patch]
func (c *userController) UpdateUserRole(ctx fiber.Ctx) error {
	var req dto.UpdateUserRoleRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return errx.ErrBadRequest.WithError(err)
	}

	res, err := c.userService.UpdateUserRole(ctx.Context(), req)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, res)
}
