package storage_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"testing"

	"github.com/google/uuid"

	"github.com/rescuebyte/rescuebyte-be/domain/errx"
	storageService "github.com/rescuebyte/rescuebyte-be/internal/app/storage/service"
)

type mockSupabaseStorage struct {
	uploadedFolder   string
	uploadedFilename string
	uploadedMIME     string
}

func (m *mockSupabaseStorage) Upload(ctx context.Context, folder string, filename string, contentType string, reader io.Reader) (string, error) {
	m.uploadedFolder = folder
	m.uploadedFilename = filename
	m.uploadedMIME = contentType
	return fmt.Sprintf("https://mock.supabase.co/storage/v1/object/public/test/%s/%s", folder, filename), nil
}

func (m *mockSupabaseStorage) Delete(ctx context.Context, fileURL string) error {
	return nil
}

func (m *mockSupabaseStorage) GetPublicURL(folder string, filename string) string {
	return fmt.Sprintf("https://mock.supabase.co/storage/v1/object/public/test/%s/%s", folder, filename)
}

type mockUUIDService struct{}

func (m *mockUUIDService) NewV4() (uuid.UUID, error) {
	return uuid.MustParse("018f0a00-0000-7000-8000-000000000001"), nil
}

func (m *mockUUIDService) NewV7() (uuid.UUID, error) {
	return uuid.MustParse("018f0a00-0000-7000-8000-000000000002"), nil
}

func (m *mockUUIDService) Parse(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

func createTestFileHeader(filename string, content []byte) *multipart.FileHeader {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, filename))
	part, _ := writer.CreatePart(h)
	part.Write(content)
	writer.Close()

	reader := multipart.NewReader(body, writer.Boundary())
	form, _ := reader.ReadForm(int64(len(content) + 1024))
	return form.File["file"][0]
}

func TestStorageService_InvalidFolder(t *testing.T) {
	svc := storageService.NewStorageService(&mockSupabaseStorage{}, &mockUUIDService{})
	fh := createTestFileHeader("test.jpg", []byte("dummy"))

	_, err := svc.UploadFile(context.Background(), fh, "invalid_folder")
	if err == nil {
		t.Fatal("expected error for invalid folder, got nil")
	}

	reqErr, ok := err.(*errx.RequestError)
	if !ok || reqErr.ErrorCode != errx.ErrInvalidFolder.ErrorCode {
		t.Errorf("expected ErrInvalidFolder, got %v", err)
	}
}

func TestStorageService_InvalidFileType(t *testing.T) {
	svc := storageService.NewStorageService(&mockSupabaseStorage{}, &mockUUIDService{})
	// plain text bytes
	fh := createTestFileHeader("test.txt", []byte("This is plain text and not an image"))

	_, err := svc.UploadFile(context.Background(), fh, "products")
	if err == nil {
		t.Fatal("expected error for text/plain file, got nil")
	}

	reqErr, ok := err.(*errx.RequestError)
	if !ok || reqErr.ErrorCode != errx.ErrInvalidFileType.ErrorCode {
		t.Errorf("expected ErrInvalidFileType, got %v", err)
	}
}

func TestStorageService_FileSizeExceeded(t *testing.T) {
	svc := storageService.NewStorageService(&mockSupabaseStorage{}, &mockUUIDService{})
	// Create dummy large file header (> 2MB for avatars)
	largeData := make([]byte, 3*1024*1024)
	fh := createTestFileHeader("large.jpg", largeData)

	_, err := svc.UploadFile(context.Background(), fh, "avatars")
	if err == nil {
		t.Fatal("expected file size exceeded error, got nil")
	}

	reqErr, ok := err.(*errx.RequestError)
	if !ok || reqErr.ErrorCode != errx.ErrFileSizeExceeded.ErrorCode {
		t.Errorf("expected ErrFileSizeExceeded, got %v", err)
	}
}

func TestStorageService_SuccessJPEG(t *testing.T) {
	mockStorage := &mockSupabaseStorage{}
	svc := storageService.NewStorageService(mockStorage, &mockUUIDService{})

	// JPEG magic bytes: FF D8 FF
	jpegData := []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46, 0x00, 0x01}
	fh := createTestFileHeader("photo.jpeg", jpegData)

	res, err := svc.UploadFile(context.Background(), fh, "products")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.ContentType != "image/jpeg" {
		t.Errorf("expected content type image/jpeg, got %s", res.ContentType)
	}

	if mockStorage.uploadedFolder != "products" {
		t.Errorf("expected folder products, got %s", mockStorage.uploadedFolder)
	}

	if res.URL == "" {
		t.Error("expected non-empty URL")
	}
}

func TestStorageService_UploadMultiple(t *testing.T) {
	mockStorage := &mockSupabaseStorage{}
	svc := storageService.NewStorageService(mockStorage, &mockUUIDService{})

	jpegData := []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46, 0x00, 0x01}
	fh1 := createTestFileHeader("photo1.jpeg", jpegData)
	fh2 := createTestFileHeader("photo2.jpeg", jpegData)

	res, err := svc.UploadMultipleFiles(context.Background(), []*multipart.FileHeader{fh1, fh2}, "products")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res) != 2 {
		t.Errorf("expected 2 responses, got %d", len(res))
	}
}
