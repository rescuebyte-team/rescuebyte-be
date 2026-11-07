package supabase_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rescuebyte/rescuebyte-be/internal/infra/env"
	"github.com/rescuebyte/rescuebyte-be/pkg/supabase"
)

func configureStorage(t *testing.T, baseURL string) {
	t.Helper()
	original := *env.AppEnv
	t.Cleanup(func() { *env.AppEnv = original })
	env.AppEnv.SupabaseURL = baseURL
	env.AppEnv.SupabaseKey = "dummy-key"
	env.AppEnv.SupabaseStorageBucket = "rescuebyte-public"
}

func TestGetPublicURL(t *testing.T) {
	configureStorage(t, "https://xyz.supabase.co/")
	storage := supabase.NewSupabaseStorage()

	url := storage.GetPublicURL("/products/", "/sample.webp")
	expected := "https://xyz.supabase.co/storage/v1/object/public/rescuebyte-public/products/sample.webp"
	if url != expected {
		t.Fatalf("expected %s, got %s", expected, url)
	}
}

func TestDeleteExtractsObjectPath(t *testing.T) {
	for _, input := range []string{
		"https://xyz.supabase.co/storage/v1/object/public/rescuebyte-public/products/test.webp",
		"/products/test.webp",
	} {
		t.Run(input, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodDelete {
					t.Errorf("expected DELETE, got %s", r.Method)
				}
				if r.URL.Path != "/storage/v1/object/rescuebyte-public" {
					t.Errorf("unexpected delete endpoint: %s", r.URL.Path)
				}
				var payload struct {
					Prefixes []string `json:"prefixes"`
				}
				if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
					t.Errorf("decode delete payload: %v", err)
					return
				}
				if len(payload.Prefixes) != 1 || payload.Prefixes[0] != "products/test.webp" {
					t.Errorf("unexpected object path: %#v", payload.Prefixes)
				}
				w.WriteHeader(http.StatusNoContent)
			}))
			defer server.Close()
			configureStorage(t, server.URL)
			storage := supabase.NewSupabaseStorage()
			if err := storage.Delete(context.Background(), input); err != nil {
				t.Fatalf("Delete() error = %v", err)
			}
		})
	}
}

func TestUnconfiguredStorage(t *testing.T) {
	configureStorage(t, "")
	env.AppEnv.SupabaseKey = ""
	storage := supabase.NewSupabaseStorage()

	_, err := storage.Upload(context.Background(), "products", "test.jpg", "image/jpeg", strings.NewReader("dummy"))
	if err == nil {
		t.Fatal("expected error for unconfigured storage, got nil")
	}

	if err := storage.Delete(context.Background(), "products/test.jpg"); err == nil {
		t.Fatal("expected error for unconfigured storage, got nil")
	}
}
