package service

import (
	"github.com/rescuebyte/rescuebyte-be/domain/contracts"
	"github.com/rescuebyte/rescuebyte-be/internal/infra/database"
	"github.com/rescuebyte/rescuebyte-be/pkg/midtrans"
	"github.com/rescuebyte/rescuebyte-be/pkg/uuid"
	"github.com/rescuebyte/rescuebyte-be/pkg/validator"
)

type orderService struct {
	transactor  database.Transactor
	orderRepo   contracts.OrderRepository
	cartRepo    contracts.CartRepository
	productRepo contracts.ProductRepository
	storeRepo   contracts.StoreRepository
	paymentRepo contracts.PaymentRepository
	midtransSvc midtrans.Service
	validator   validator.Service
	uuid        uuid.Service
}

func NewOrderService(
	transactor database.Transactor,
	orderRepo contracts.OrderRepository,
	cartRepo contracts.CartRepository,
	productRepo contracts.ProductRepository,
	storeRepo contracts.StoreRepository,
	paymentRepo contracts.PaymentRepository,
	midtransSvc midtrans.Service,
	validator validator.Service,
	uuid uuid.Service,
) contracts.OrderService {
	return &orderService{
		transactor:  transactor,
		orderRepo:   orderRepo,
		cartRepo:    cartRepo,
		productRepo: productRepo,
		storeRepo:   storeRepo,
		paymentRepo: paymentRepo,
		midtransSvc: midtransSvc,
		validator:   validator,
		uuid:        uuid,
	}
}
