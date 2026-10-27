package contracts

import (
	"context"
	"mime/multipart"

	"github.com/rescuebyte/rescuebyte-be/domain/dto"
)

type StorageService interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader, folder string) (*dto.UploadResponse, error)
	UploadMultipleFiles(ctx context.Context, files []*multipart.FileHeader, folder string) ([]dto.UploadResponse, error)
	DeleteFile(ctx context.Context, fileURL string) error
}
