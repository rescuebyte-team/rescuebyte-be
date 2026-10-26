package repository

import (
	"github.com/rescuebyte/rescuebyte-be/domain/contracts"
	"github.com/rescuebyte/rescuebyte-be/internal/infra/database"
)

type storeRepository struct {
	dbGetter database.DBGetter
}

func NewStoreRepository(dbGetter database.DBGetter) contracts.StoreRepository {
	return &storeRepository{dbGetter: dbGetter}
}
