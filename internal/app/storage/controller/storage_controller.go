package controller

import (
	"github.com/gofiber/fiber/v3"

	_ "github.com/rescuebyte/rescuebyte-be/domain/dto"
	"github.com/rescuebyte/rescuebyte-be/domain/errx"
	"github.com/rescuebyte/rescuebyte-be/pkg/helpers/http/response"
)

// Upload godoc
// @Summary Upload an image
// @Description Upload a single image file (JPEG, PNG, WebP) to Supabase Storage
// @Tags Storage
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param folder formData string false "Target folder: products, stores, avatars (default: products)"
// @Param file formData file true "Image file to upload (max 5MB for products, 2MB for stores/avatars)"
// @Success 201 {object} response.Response{payload=dto.UploadResponse} "Uploaded image details"
// @Failure 400 {object} response.ErrorResponse "Bad request or validation error"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Storage upload failed"
// @Router /api/v1/storage/upload [post]
func (c *storageController) Upload(ctx fiber.Ctx) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		return errx.ErrFileRequired.WithError(err)
	}

	folder := ctx.FormValue("folder")
	if folder == "" {
		folder = ctx.Query("folder")
	}
	if folder == "" {
		folder = "products"
	}

	res, err := c.storageService.UploadFile(ctx.Context(), file, folder)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusCreated, res)
}

// UploadMultiple godoc
// @Summary Upload multiple images
// @Description Upload multiple image files (JPEG, PNG, WebP) to Supabase Storage (max 5 files)
// @Tags Storage
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param folder formData string false "Target folder: products, stores, avatars (default: products)"
// @Param files formData []file true "Image files to upload"
// @Success 201 {object} response.Response{payload=[]dto.UploadResponse} "List of uploaded images details"
// @Failure 400 {object} response.ErrorResponse "Bad request or validation error"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Storage upload failed"
// @Router /api/v1/storage/upload-multiple [post]
func (c *storageController) UploadMultiple(ctx fiber.Ctx) error {
	form, err := ctx.MultipartForm()
	if err != nil {
		return errx.ErrFileRequired.WithError(err)
	}

	files := form.File["files"]
	if len(files) == 0 {
		return errx.ErrFileRequired
	}

	folder := ctx.FormValue("folder")
	if folder == "" {
		folder = ctx.Query("folder")
	}
	if folder == "" {
		folder = "products"
	}

	res, err := c.storageService.UploadMultipleFiles(ctx.Context(), files, folder)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusCreated, res)
}
