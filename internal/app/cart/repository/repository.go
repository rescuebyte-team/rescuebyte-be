package repository

import (
	"github.com/rescuebyte/rescuebyte-be/domain/contracts"
	"github.com/rescuebyte/rescuebyte-be/internal/infra/database"
)

type cartRepository struct {
	dbGetter database.DBGetter
}

func NewCartRepository(dbGetter database.DBGetter) contracts.CartRepository {
	return &cartRepository{dbGetter: dbGetter}
}
