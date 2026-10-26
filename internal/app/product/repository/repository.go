package repository

import (
	"github.com/rescuebyte/rescuebyte-be/domain/contracts"
	"github.com/rescuebyte/rescuebyte-be/internal/infra/database"
)

type productRepository struct {
	dbGetter database.DBGetter
}

func NewProductRepository(dbGetter database.DBGetter) contracts.ProductRepository {
	return &productRepository{dbGetter: dbGetter}
}
