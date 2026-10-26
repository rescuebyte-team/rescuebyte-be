package service

import "github.com/rescuebyte/rescuebyte-be/domain/contracts"

type enumService struct{}

func NewEnumService() contracts.EnumService {
	return &enumService{}
}
