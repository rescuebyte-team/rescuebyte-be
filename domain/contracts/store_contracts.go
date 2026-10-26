package contracts

import (
	"context"

	"github.com/google/uuid"

	"github.com/rescuebyte/rescuebyte-be/domain/dto"
	"github.com/rescuebyte/rescuebyte-be/domain/entity"
)

type StoreRepository interface {
	CreateStore(ctx context.Context, store entity.Store) error
	GetStoreByID(ctx context.Context, id uuid.UUID) (*entity.Store, error)
	GetStoreByUserID(ctx context.Context, userID uuid.UUID) (*entity.Store, error)
	UpdateStoreProfile(ctx context.Context, store entity.Store) error
	UpdateStoreOperational(ctx context.Context, store entity.Store) error
	UpdateStoreStatus(ctx context.Context, storeID uuid.UUID, status string) error
}

type StoreService interface {
	CreateStore(ctx context.Context, userID uuid.UUID, req dto.CreateStoreRequest) (*dto.StoreResponse, error)
	GetMyStore(ctx context.Context, userID uuid.UUID) (*dto.StoreResponse, error)
	GetStoreByID(ctx context.Context, id uuid.UUID) (*dto.StoreResponse, error)
	UpdateStoreProfile(ctx context.Context, userID uuid.UUID, req dto.UpdateStoreProfileRequest) (*dto.StoreResponse, error)
	UpdateStoreOperational(ctx context.Context, userID uuid.UUID, req dto.UpdateStoreOperationalRequest) (*dto.StoreResponse, error)
	ToggleStoreStatus(ctx context.Context, userID uuid.UUID) (*dto.StoreResponse, error)
}
