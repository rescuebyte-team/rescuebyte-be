package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/rescuebyte/rescuebyte-be/domain/dto"
	"github.com/rescuebyte/rescuebyte-be/domain/errx"
)

var allowedMIMETypes = map[string]string{
	"image/jpeg":  ".jpg",
	"image/pjpeg": ".jpg",
	"image/png":   ".png",
	"image/webp":  ".webp",
}

func (s *storageService) UploadFile(ctx context.Context, file *multipart.FileHeader, folder string) (*dto.UploadResponse, error) {
	if file == nil {
		return nil, errx.ErrFileRequired
	}

	cleanFolder := strings.ToLower(strings.TrimSpace(folder))
	if cleanFolder != "products" && cleanFolder != "stores" && cleanFolder != "avatars" {
		return nil, errx.ErrInvalidFolder
	}

	maxSize := int64(5 * 1024 * 1024) // 5 MB for products
	if cleanFolder == "stores" || cleanFolder == "avatars" {
		maxSize = 2 * 1024 * 1024 // 2 MB for stores and avatars
	}

	if file.Size > maxSize {
		return nil, errx.ErrFileSizeExceeded
	}

	src, err := file.Open()
	if err != nil {
		return nil, errx.ErrBadRequest.WithError(err)
	}
	defer src.Close()

	buffer := make([]byte, 512)
	n, err := src.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, errx.ErrBadRequest.WithError(err)
	}

	detectedMIME := http.DetectContentType(buffer[:n])
	ext, ok := allowedMIMETypes[detectedMIME]
	if !ok {
		return nil, errx.ErrInvalidFileType
	}

	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return nil, errx.ErrInternalServer.WithError(err)
	}

	newUUID, err := s.uuid.NewV7()
	if err != nil {
		return nil, errx.ErrInternalServer.WithError(err)
	}

	filename := fmt.Sprintf("%s%s", newUUID.String(), ext)

	fileURL, err := s.supabaseStorage.Upload(ctx, cleanFolder, filename, detectedMIME, src)
	if err != nil {
		return nil, errx.ErrStorageUploadFailed.WithError(err)
	}

	return &dto.UploadResponse{
		URL:         fileURL,
		Filename:    filename,
		ContentType: detectedMIME,
		Size:        file.Size,
	}, nil
}

func (s *storageService) UploadMultipleFiles(ctx context.Context, files []*multipart.FileHeader, folder string) ([]dto.UploadResponse, error) {
	if len(files) == 0 {
		return nil, errx.ErrFileRequired
	}

	if len(files) > 5 {
		return nil, errx.ErrBadRequest.WithMessage("Maximum 5 files can be uploaded at once")
	}

	responses := make([]dto.UploadResponse, 0, len(files))
	for _, file := range files {
		res, err := s.UploadFile(ctx, file, folder)
		if err != nil {
			return nil, err
		}
		responses = append(responses, *res)
	}

	return responses, nil
}

func (s *storageService) DeleteFile(ctx context.Context, fileURL string) error {
	if strings.TrimSpace(fileURL) == "" {
		return nil
	}
	return s.supabaseStorage.Delete(ctx, fileURL)
}
