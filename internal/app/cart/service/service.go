package service

import (
	"github.com/rescuebyte/rescuebyte-be/domain/contracts"
	"github.com/rescuebyte/rescuebyte-be/internal/infra/database"
	"github.com/rescuebyte/rescuebyte-be/pkg/uuid"
	"github.com/rescuebyte/rescuebyte-be/pkg/validator"
)

type cartService struct {
	transactor  database.Transactor
	cartRepo    contracts.CartRepository
	productRepo contracts.ProductRepository
	validator   validator.Service
	uuid        uuid.Service
}

func NewCartService(
	transactor database.Transactor,
	cartRepo contracts.CartRepository,
	productRepo contracts.ProductRepository,
	validator validator.Service,
	uuid uuid.Service,
) contracts.CartService {
	return &cartService{
		transactor:  transactor,
		cartRepo:    cartRepo,
		productRepo: productRepo,
		validator:   validator,
		uuid:        uuid,
	}
}
