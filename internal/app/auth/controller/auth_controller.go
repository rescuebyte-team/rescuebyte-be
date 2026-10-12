package controller

import (
	"github.com/gofiber/fiber/v3"

	"github.com/rescuebyte/rescuebyte-be/domain/dto"
	"github.com/rescuebyte/rescuebyte-be/domain/errx"
	"github.com/rescuebyte/rescuebyte-be/pkg/helpers/http/response"
	"github.com/rescuebyte/rescuebyte-be/pkg/jwt"
)

func (c *authController) Register(ctx fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return errx.ErrBadRequest.WithError(err)
	}

	res, err := c.authService.Register(ctx.Context(), req)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusCreated, res)
}

func (c *authController) Login(ctx fiber.Ctx) error {
	var req dto.LoginRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return errx.ErrBadRequest.WithError(err)
	}

	res, err := c.authService.Login(ctx.Context(), req)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, res)
}

func (c *authController) RefreshToken(ctx fiber.Ctx) error {
	var req dto.RefreshTokenRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return errx.ErrBadRequest.WithError(err)
	}

	res, err := c.authService.RefreshToken(ctx.Context(), req)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, res)
}

func (c *authController) Logout(ctx fiber.Ctx) error {
	var req dto.RefreshTokenRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return errx.ErrBadRequest.WithError(err)
	}

	if err := c.authService.Logout(ctx.Context(), req.RefreshToken); err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, map[string]string{
		"message": "Successfully logged out",
	})
}

func (c *authController) GetGoogleAuthURL(ctx fiber.Ctx) error {
	res, err := c.authService.GetGoogleAuthURL(ctx.Context())
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, res)
}

func (c *authController) GoogleCallback(ctx fiber.Ctx) error {
	var query dto.GoogleCallbackQuery
	if err := ctx.Bind().Query(&query); err != nil {
		return errx.ErrBadRequest.WithError(err)
	}

	res, err := c.authService.HandleGoogleCallback(ctx.Context(), query)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, res)
}

func (c *authController) GetMe(ctx fiber.Ctx) error {
	claims, ok := ctx.Locals("claims").(jwt.Claims)
	if !ok {
		return errx.ErrUnauthorized
	}

	res, err := c.authService.GetMe(ctx.Context(), claims.ID)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, res)
}
