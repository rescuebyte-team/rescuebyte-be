package service

import (
	"github.com/rescuebyte/rescuebyte-be/domain/contracts"
	"github.com/rescuebyte/rescuebyte-be/internal/infra/database"
	"github.com/rescuebyte/rescuebyte-be/pkg/midtrans"
)

type paymentService struct {
	transactor  database.Transactor
	paymentRepo contracts.PaymentRepository
	orderRepo   contracts.OrderRepository
	productRepo contracts.ProductRepository
	midtransSvc midtrans.Service
}

func NewPaymentService(
	transactor database.Transactor,
	paymentRepo contracts.PaymentRepository,
	orderRepo contracts.OrderRepository,
	productRepo contracts.ProductRepository,
	midtransSvc midtrans.Service,
) contracts.PaymentService {
	return &paymentService{
		transactor:  transactor,
		paymentRepo: paymentRepo,
		orderRepo:   orderRepo,
		productRepo: productRepo,
		midtransSvc: midtransSvc,
	}
}
