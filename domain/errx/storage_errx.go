package errx

import "github.com/gofiber/fiber/v3"

var (
	ErrFileRequired = NewError(
		fiber.StatusBadRequest,
		"file_required",
		"No file was provided in the request.",
	)
	ErrInvalidFileType = NewError(
		fiber.StatusBadRequest,
		"invalid_file_type",
		"Invalid file type. Only JPEG, PNG, and WebP images are allowed.",
	)
	ErrFileSizeExceeded = NewError(
		fiber.StatusBadRequest,
		"file_size_exceeded",
		"The uploaded file exceeds the maximum allowed size.",
	)
	ErrInvalidFolder = NewError(
		fiber.StatusBadRequest,
		"invalid_folder",
		"The specified upload destination folder is invalid.",
	)
	ErrStorageUploadFailed = NewError(
		fiber.StatusInternalServerError,
		"storage_upload_failed",
		"Failed to upload the file to storage.",
	)
)
