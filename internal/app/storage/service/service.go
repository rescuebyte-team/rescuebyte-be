package service

import (
	"github.com/rescuebyte/rescuebyte-be/domain/contracts"
	"github.com/rescuebyte/rescuebyte-be/pkg/supabase"
	"github.com/rescuebyte/rescuebyte-be/pkg/uuid"
)

type storageService struct {
	supabaseStorage supabase.StorageService
	uuid            uuid.Service
}

func NewStorageService(
	supabaseStorage supabase.StorageService,
	uuid uuid.Service,
) contracts.StorageService {
	return &storageService{
		supabaseStorage: supabaseStorage,
		uuid:            uuid,
	}
}
