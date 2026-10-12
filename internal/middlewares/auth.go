package middlewares

import (
	"strings"

	"github.com/gofiber/fiber/v3"

	"github.com/rescuebyte/rescuebyte-be/domain/enums"
	"github.com/rescuebyte/rescuebyte-be/domain/errx"
	"github.com/rescuebyte/rescuebyte-be/pkg/jwt"
)

func (m *Middleware) extractTokenClaims(ctx fiber.Ctx) (*jwt.Claims, error) {
	header := ctx.Get("Authorization")
	if header == "" {
		return nil, errx.ErrNoBearerToken
	}

	headerSlice := strings.Split(header, " ")
	if len(headerSlice) != 2 || strings.ToLower(headerSlice[0]) != "bearer" {
		return nil, errx.ErrInvalidBearerToken
	}

	token := headerSlice[1]
	var claims jwt.Claims
	err := m.jwt.Decode(token, &claims)
	if err != nil {
		return nil, errx.ErrInvalidBearerToken.WithError(err)
	}

	if claims.TokenType != "access" {
		return nil, errx.ErrInvalidBearerToken
	}

	return &claims, nil
}

func (m *Middleware) RequireAuth() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		claims, err := m.extractTokenClaims(ctx)
		if err != nil {
			return err
		}

		ctx.Locals("claims", *claims)
		return ctx.Next()
	}
}

func (m *Middleware) OptionalAuth() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		claims, err := m.extractTokenClaims(ctx)
		if err == nil && claims != nil {
			ctx.Locals("claims", *claims)
		}

		return ctx.Next()
	}
}

// RequireRoles grants access only to users who have one of the accepted roles
func (m *Middleware) RequireRoles(acceptedRoles ...enums.UserRole) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		claims, ok := ctx.Locals("claims").(jwt.Claims)
		if !ok {
			return errx.ErrUnauthorized
		}

		if len(acceptedRoles) == 0 {
			return ctx.Next()
		}

		for _, role := range acceptedRoles {
			if claims.Role == role {
				return ctx.Next()
			}
		}

		return errx.ErrForbidden.WithMessage("You do not have the required role to access this resource.")
	}
}
