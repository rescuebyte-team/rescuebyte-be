package controller

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"github.com/rescuebyte/rescuebyte-be/domain/dto"
	"github.com/rescuebyte/rescuebyte-be/domain/errx"
	"github.com/rescuebyte/rescuebyte-be/pkg/helpers/http/response"
	"github.com/rescuebyte/rescuebyte-be/pkg/jwt"
)

// CreateStore godoc
// @Summary Create a new store
// @Description Register a new store for the authenticated user
// @Tags Store
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateStoreRequest true "Create Store Request Payload"
// @Success 201 {object} response.Response{payload=dto.StoreResponse} "Created Store details"
// @Failure 400 {object} response.ErrorResponse "Bad request or validation error"
// @Failure 401 {object} response.ErrorResponse "Unauthorized or missing token"
// @Failure 409 {object} response.ErrorResponse "User already has a store"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/stores [post]
func (c *storeController) CreateStore(ctx fiber.Ctx) error {
	claims, ok := ctx.Locals("claims").(jwt.Claims)
	if !ok {
		return errx.ErrUnauthorized
	}

	var req dto.CreateStoreRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return errx.ErrBadRequest.WithError(err)
	}

	res, err := c.storeService.CreateStore(ctx.Context(), claims.ID, req)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusCreated, res)
}

// GetMyStore godoc
// @Summary Get authenticated user's store
// @Description Retrieve the store profile of the currently logged-in user
// @Tags Store
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.Response{payload=dto.StoreResponse} "Store profile details"
// @Failure 401 {object} response.ErrorResponse "Unauthorized or missing token"
// @Failure 404 {object} response.ErrorResponse "Store not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/stores/me/profile [get]
func (c *storeController) GetMyStore(ctx fiber.Ctx) error {
	claims, ok := ctx.Locals("claims").(jwt.Claims)
	if !ok {
		return errx.ErrUnauthorized
	}

	res, err := c.storeService.GetMyStore(ctx.Context(), claims.ID)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, res)
}

// GetAllStores godoc
// @Summary Get all stores
// @Description Retrieve a list of active stores with optional filtering by name
// @Tags Store
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param name query string false "Filter by store name"
// @Param nama query string false "Alias filter by store name"
// @Success 200 {object} response.Response{payload=dto.GetStoresResponse} "List of stores with pagination metadata"
// @Failure 400 {object} response.ErrorResponse "Bad request or validation error"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/stores [get]
func (c *storeController) GetAllStores(ctx fiber.Ctx) error {
	var query dto.GetAllStoresQuery
	if err := ctx.Bind().Query(&query); err != nil {
		return errx.ErrBadRequest.WithError(err)
	}

	res, err := c.storeService.GetAllStores(ctx.Context(), query)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, res)
}

// GetStoreByID godoc
// @Summary Get store by ID
// @Description Retrieve public store profile by its UUID
// @Tags Store
// @Produce json
// @Param id path string true "Store UUID"
// @Success 200 {object} response.Response{payload=dto.StoreResponse} "Store profile details"
// @Failure 400 {object} response.ErrorResponse "Invalid UUID format"
// @Failure 404 {object} response.ErrorResponse "Store not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/stores/{id} [get]
func (c *storeController) GetStoreByID(ctx fiber.Ctx) error {
	idParam := ctx.Params("id")
	storeID, err := uuid.Parse(idParam)
	if err != nil {
		return errx.ErrBadRequest.WithError(err)
	}

	res, err := c.storeService.GetStoreByID(ctx.Context(), storeID)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, res)
}

// UpdateStoreProfile godoc
// @Summary Update store profile
// @Description Update basic store details (name, description, contact, location)
// @Tags Store
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.UpdateStoreProfileRequest true "Update Store Profile Payload"
// @Success 200 {object} response.Response{payload=dto.StoreResponse} "Updated store details"
// @Failure 400 {object} response.ErrorResponse "Bad request or validation error"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Store not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/stores/me/profile [put]
func (c *storeController) UpdateStoreProfile(ctx fiber.Ctx) error {
	claims, ok := ctx.Locals("claims").(jwt.Claims)
	if !ok {
		return errx.ErrUnauthorized
	}

	var req dto.UpdateStoreProfileRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return errx.ErrBadRequest.WithError(err)
	}

	res, err := c.storeService.UpdateStoreProfile(ctx.Context(), claims.ID, req)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, res)
}

// UpdateStoreOperational godoc
// @Summary Update store operational info
// @Description Update operational hours, operating days, and status
// @Tags Store
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.UpdateStoreOperationalRequest true "Update Operational Payload"
// @Success 200 {object} response.Response{payload=dto.StoreResponse} "Updated store details"
// @Failure 400 {object} response.ErrorResponse "Bad request or validation error"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Store not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/stores/me/operational [put]
func (c *storeController) UpdateStoreOperational(ctx fiber.Ctx) error {
	claims, ok := ctx.Locals("claims").(jwt.Claims)
	if !ok {
		return errx.ErrUnauthorized
	}

	var req dto.UpdateStoreOperationalRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return errx.ErrBadRequest.WithError(err)
	}

	res, err := c.storeService.UpdateStoreOperational(ctx.Context(), claims.ID, req)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, res)
}

// ToggleStoreStatus godoc
// @Summary Toggle store active status
// @Description Switch store status between ACTIVE and INACTIVE
// @Tags Store
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.Response{payload=dto.StoreResponse} "Updated store with toggled status"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Store not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/stores/me/status [patch]
func (c *storeController) ToggleStoreStatus(ctx fiber.Ctx) error {
	claims, ok := ctx.Locals("claims").(jwt.Claims)
	if !ok {
		return errx.ErrUnauthorized
	}

	res, err := c.storeService.ToggleStoreStatus(ctx.Context(), claims.ID)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, res)
}
