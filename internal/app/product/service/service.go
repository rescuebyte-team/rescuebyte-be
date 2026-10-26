package service

import (
	"github.com/rescuebyte/rescuebyte-be/domain/contracts"
	"github.com/rescuebyte/rescuebyte-be/internal/infra/database"
	"github.com/rescuebyte/rescuebyte-be/pkg/uuid"
	"github.com/rescuebyte/rescuebyte-be/pkg/validator"
)

type productService struct {
	transactor  database.Transactor
	productRepo contracts.ProductRepository
	storeRepo   contracts.StoreRepository
	validator   validator.Service
	uuid        uuid.Service
}

func NewProductService(
	transactor database.Transactor,
	productRepo contracts.ProductRepository,
	storeRepo contracts.StoreRepository,
	validator validator.Service,
	uuid uuid.Service,
) contracts.ProductService {
	return &productService{
		transactor:  transactor,
		productRepo: productRepo,
		storeRepo:   storeRepo,
		validator:   validator,
		uuid:        uuid,
	}
}
