package controller

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rescuebyte/rescuebyte-be/domain/dto"
	"github.com/rescuebyte/rescuebyte-be/pkg/helpers/http/response"
)

// GetEnums godoc
// @Summary Get system enums
// @Description Retrieve system enums and lookup values (optional filter by comma-separated types)
// @Tags Enums
// @Accept json
// @Produce json
// @Param types query string false "Comma-separated enum types filter (e.g. user_role,auth_provider)"
// @Success 200 {object} response.Response{payload=map[string]any} "Enums retrieved successfully"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/enums [get]
func (c *enumController) GetEnums(ctx fiber.Ctx) error {
	var query dto.GetEnumsQuery
	if err := ctx.Bind().Query(&query); err != nil {
		return err
	}

	result, err := c.enumService.GetEnums(query.Types)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, result)
}
