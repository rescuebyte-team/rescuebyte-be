package controller

import (
	"github.com/gofiber/fiber/v3"

	"github.com/rescuebyte/rescuebyte-be/domain/dto"
	"github.com/rescuebyte/rescuebyte-be/domain/errx"
	"github.com/rescuebyte/rescuebyte-be/pkg/helpers/http/response"
	"github.com/rescuebyte/rescuebyte-be/pkg/jwt"
)

// Register godoc
// @Summary Register new user account
// @Description Register a new user with name, email, and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "User Registration Details"
// @Success 201 {object} response.Response{payload=dto.AuthResponse} "User successfully registered"
// @Failure 400 {object} response.ErrorResponse "Bad request or validation error"
// @Failure 409 {object} response.ErrorResponse "Email already exists"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/auth/register [post]
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

// Login godoc
// @Summary Login user
// @Description Authenticate user credentials and return JWT tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login Credentials"
// @Success 200 {object} response.Response{payload=dto.AuthResponse} "Login successful"
// @Failure 400 {object} response.ErrorResponse "Bad request or validation error"
// @Failure 401 {object} response.ErrorResponse "Invalid email or password"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/auth/login [post]
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

// RefreshToken godoc
// @Summary Refresh access token
// @Description Generate a new access and refresh token pair using an existing valid refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh Token Payload"
// @Success 200 {object} response.Response{payload=dto.AuthResponse} "Token successfully refreshed"
// @Failure 400 {object} response.ErrorResponse "Bad request or missing token"
// @Failure 401 {object} response.ErrorResponse "Invalid or expired refresh token"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/auth/refresh [post]
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

// Logout godoc
// @Summary Logout user
// @Description Revoke the refresh token and invalidate the user session
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh Token Payload"
// @Success 200 {object} response.Response{payload=map[string]string} "Successfully logged out"
// @Failure 400 {object} response.ErrorResponse "Bad request or missing token"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/auth/logout [post]
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

// GetGoogleAuthURL godoc
// @Summary Get Google OAuth2 URL
// @Description Retrieve Google OAuth2 consent URL and CSRF state token
// @Tags Auth
// @Produce json
// @Success 200 {object} response.Response{payload=dto.GoogleAuthURLResponse} "Google OAuth URL generated"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/auth/google [get]
func (c *authController) GetGoogleAuthURL(ctx fiber.Ctx) error {
	res, err := c.authService.GetGoogleAuthURL(ctx.Context())
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, res)
}

// GoogleCallback godoc
// @Summary Handle Google OAuth2 callback
// @Description Authenticate or register a user through Google OAuth2 callback query parameters
// @Tags Auth
// @Produce json
// @Param code query string true "OAuth Authorization Code"
// @Param state query string true "OAuth State"
// @Success 200 {object} response.Response{payload=dto.AuthResponse} "Authentication successful"
// @Failure 400 {object} response.ErrorResponse "Bad request or state mismatch"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/auth/google/callback [get]
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

// GetMe godoc
// @Summary Get current authenticated user
// @Description Retrieve the profile of currently logged-in user from token claims
// @Tags Auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.Response{payload=dto.UserResponse} "Current user profile"
// @Failure 401 {object} response.ErrorResponse "Unauthorized or missing token"
// @Failure 404 {object} response.ErrorResponse "User not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/auth/me [get]
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
