package errx

import (
	"github.com/gofiber/fiber/v3"
)

var (
	ErrEmailAlreadyExists = NewError(
		fiber.StatusConflict,
		"email_already_exists",
		"An account with this email address already exists.",
	)
	ErrInvalidCredentials = NewError(
		fiber.StatusUnauthorized,
		"invalid_credentials",
		"Invalid email or password.",
	)
	ErrUserNotFound = NewError(
		fiber.StatusNotFound,
		"user_not_found",
		"The specified user could not be found.",
	)
	ErrPasswordNotSet = NewError(
		fiber.StatusBadRequest,
		"password_not_set",
		"This account was registered via Google OAuth. Please sign in with Google.",
	)
	ErrInvalidRefreshToken = NewError(
		fiber.StatusUnauthorized,
		"invalid_refresh_token",
		"The provided refresh token is invalid or has expired.",
	)
	ErrRefreshTokenRevoked = NewError(
		fiber.StatusUnauthorized,
		"refresh_token_revoked",
		"The provided refresh token has been revoked.",
	)
	ErrOAuthExchangeFailed = NewError(
		fiber.StatusBadRequest,
		"oauth_exchange_failed",
		"Failed to exchange code with OAuth provider.",
	)
	ErrOAuthGetUserInfoFailed = NewError(
		fiber.StatusBadRequest,
		"oauth_get_user_info_failed",
		"Failed to retrieve user information from OAuth provider.",
	)
)
