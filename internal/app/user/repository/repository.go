package repository

import (
	"github.com/rescuebyte/rescuebyte-be/domain/contracts"
	"github.com/rescuebyte/rescuebyte-be/internal/infra/database"
)

type userRepository struct {
	dbGetter database.DBGetter
}

func NewUserRepository(dbGetter database.DBGetter) contracts.UserRepository {
	return &userRepository{dbGetter: dbGetter}
}
