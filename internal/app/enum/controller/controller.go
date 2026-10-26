package controller

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rescuebyte/rescuebyte-be/domain/contracts"
	baseCtr "github.com/rescuebyte/rescuebyte-be/internal/app/base/controller"
)

type enumController struct {
	baseCtr.BaseController
	enumService contracts.EnumService
}

func InitEnumController(router fiber.Router, enumService contracts.EnumService) {
	ctr := &enumController{
		enumService: enumService,
	}
	router.Get("/enums", ctr.GetEnums)
}
