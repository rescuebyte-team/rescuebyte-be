package repository

import (
	"github.com/rescuebyte/rescuebyte-be/domain/contracts"
	"github.com/rescuebyte/rescuebyte-be/internal/infra/database"
)

type paymentRepository struct {
	dbGetter database.DBGetter
}

func NewPaymentRepository(dbGetter database.DBGetter) contracts.PaymentRepository {
	return &paymentRepository{dbGetter: dbGetter}
}
