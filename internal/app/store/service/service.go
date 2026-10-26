package service

import (
	"github.com/rescuebyte/rescuebyte-be/domain/contracts"
	"github.com/rescuebyte/rescuebyte-be/internal/infra/database"
	"github.com/rescuebyte/rescuebyte-be/pkg/uuid"
	"github.com/rescuebyte/rescuebyte-be/pkg/validator"
)

type storeService struct {
	transactor database.Transactor
	storeRepo  contracts.StoreRepository
	userRepo   contracts.UserRepository
	validator  validator.Service
	uuid       uuid.Service
}

func NewStoreService(
	transactor database.Transactor,
	storeRepo contracts.StoreRepository,
	userRepo contracts.UserRepository,
	validator validator.Service,
	uuid uuid.Service,
) contracts.StoreService {
	return &storeService{
		transactor: transactor,
		storeRepo:  storeRepo,
		userRepo:   userRepo,
		validator:  validator,
		uuid:       uuid,
	}
}
