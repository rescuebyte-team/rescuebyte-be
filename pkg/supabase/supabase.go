package supabase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/rescuebyte/rescuebyte-be/internal/infra/env"
	"github.com/rescuebyte/rescuebyte-be/pkg/log"
)

type StorageService interface {
	Upload(ctx context.Context, folder string, filename string, contentType string, reader io.Reader) (string, error)
	Delete(ctx context.Context, fileURL string) error
	GetPublicURL(folder string, filename string) string
}

type supabaseStorageImpl struct {
	baseURL string
	apiKey  string
	bucket  string
	client  *http.Client
}

var SupabaseStorage StorageService = NewSupabaseStorage()

func NewSupabaseStorage() StorageService {
	return &supabaseStorageImpl{
		baseURL: strings.TrimRight(env.AppEnv.SupabaseURL, "/"),
		apiKey:  env.AppEnv.SupabaseKey,
		bucket:  env.AppEnv.SupabaseStorageBucket,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *supabaseStorageImpl) Upload(ctx context.Context, folder string, filename string, contentType string, reader io.Reader) (string, error) {
	if s.baseURL == "" || s.apiKey == "" || s.bucket == "" {
		return "", fmt.Errorf("supabase storage is not properly configured: SUPABASE_URL, SUPABASE_KEY, or SUPABASE_STORAGE_BUCKET is empty")
	}

	cleanFolder := strings.Trim(folder, "/")
	cleanFilename := strings.Trim(filename, "/")
	objectPath := fmt.Sprintf("%s/%s", cleanFolder, cleanFilename)
	uploadURL := fmt.Sprintf("%s/storage/v1/object/%s/%s", s.baseURL, s.bucket, objectPath)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uploadURL, reader)
	if err != nil {
		return "", fmt.Errorf("failed to create upload request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))
	req.Header.Set("apikey", s.apiKey)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("x-upsert", "true")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("supabase storage request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Error(log.Fields{
			"status_code": resp.StatusCode,
			"response":    string(bodyBytes),
			"url":         uploadURL,
		}, "[SupabaseStorage] Upload failed")
		return "", fmt.Errorf("supabase storage upload failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return s.GetPublicURL(cleanFolder, cleanFilename), nil
}

func (s *supabaseStorageImpl) GetPublicURL(folder string, filename string) string {
	cleanFolder := strings.Trim(folder, "/")
	cleanFilename := strings.Trim(filename, "/")
	return fmt.Sprintf("%s/storage/v1/object/public/%s/%s/%s", s.baseURL, s.bucket, cleanFolder, cleanFilename)
}

func (s *supabaseStorageImpl) Delete(ctx context.Context, fileURL string) error {
	if s.baseURL == "" || s.apiKey == "" || s.bucket == "" {
		return fmt.Errorf("supabase storage is not properly configured: SUPABASE_URL, SUPABASE_KEY, or SUPABASE_STORAGE_BUCKET is empty")
	}

	path := s.extractPath(fileURL)
	if path == "" {
		return fmt.Errorf("invalid file path to delete: %s", fileURL)
	}

	deleteURL := fmt.Sprintf("%s/storage/v1/object/%s", s.baseURL, s.bucket)
	payload, err := json.Marshal(map[string]any{
		"prefixes": []string{path},
	})
	if err != nil {
		return fmt.Errorf("failed to marshal delete payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, deleteURL, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("failed to create delete request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))
	req.Header.Set("apikey", s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("supabase storage delete request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Error(log.Fields{
			"status_code": resp.StatusCode,
			"response":    string(bodyBytes),
			"path":        path,
		}, "[SupabaseStorage] Delete failed")
		return fmt.Errorf("supabase storage delete failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

func (s *supabaseStorageImpl) extractPath(fileURL string) string {
	marker := fmt.Sprintf("/storage/v1/object/public/%s/", s.bucket)
	idx := strings.Index(fileURL, marker)
	if idx != -1 {
		return fileURL[idx+len(marker):]
	}
	return strings.TrimPrefix(fileURL, "/")
}
