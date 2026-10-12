package service

import (
	"github.com/rescuebyte/rescuebyte-be/domain/contracts"
	"github.com/rescuebyte/rescuebyte-be/internal/infra/database"
	"github.com/rescuebyte/rescuebyte-be/pkg/validator"
)

type userService struct {
	transactor database.Transactor
	userRepo   contracts.UserRepository
	validator  validator.Service
}

func NewUserService(
	transactor database.Transactor,
	userRepo contracts.UserRepository,
	validator validator.Service,
) contracts.UserService {
	return &userService{
		transactor: transactor,
		userRepo:   userRepo,
		validator:  validator,
	}
}
