package repository

import (
	"github.com/rescuebyte/rescuebyte-be/domain/contracts"
	"github.com/rescuebyte/rescuebyte-be/internal/infra/database"
)

type orderRepository struct {
	dbGetter database.DBGetter
}

func NewOrderRepository(dbGetter database.DBGetter) contracts.OrderRepository {
	return &orderRepository{dbGetter: dbGetter}
}
