package controller

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/rescuebyte/rescuebyte-be/domain/dto"
	"github.com/rescuebyte/rescuebyte-be/domain/enums"
	"github.com/rescuebyte/rescuebyte-be/domain/errx"
	"github.com/rescuebyte/rescuebyte-be/pkg/jwt"
	"github.com/rescuebyte/rescuebyte-be/pkg/log"
)

type BaseController struct{}

func (b *BaseController) BuildLogActor(ctx fiber.Ctx) (*dto.LogActor, error) {
	claims, ok := ctx.Locals("claims").(jwt.Claims)
	if !ok {
		log.Ctx(ctx.Context()).Error(log.Fields{
			"error": "failed to get claims from context",
		}, "[BaseController][BuildLogActor] failed to get claims from context")

		return nil, errx.ErrClaimsNotFound
	}

	return &dto.LogActor{
		ID:    claims.ID,
		Name:  claims.Name,
		Email: claims.Email,
	}, nil
}

func (b *BaseController) ExtractActorRole(ctx fiber.Ctx) enums.UserRoleIdx {
	claims, ok := ctx.Locals("claims").(jwt.Claims)
	if !ok {
		return 0
	}
	return claims.Role
}

func (b *BaseController) ExtractActorID(ctx fiber.Ctx) uuid.UUID {
	claims, ok := ctx.Locals("claims").(jwt.Claims)
	if !ok {
		return uuid.Nil
	}
	return claims.ID
}

func (b *BaseController) IsAuthenticated(ctx fiber.Ctx) bool {
	_, ok := ctx.Locals("claims").(jwt.Claims)
	return ok
}
