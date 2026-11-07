package service

import (
	"github.com/rescuebyte/rescuebyte-be/domain/contracts"
	"github.com/rescuebyte/rescuebyte-be/pkg/uuid"
	"github.com/rescuebyte/rescuebyte-be/pkg/validator"
)

type reviewService struct {
	reviewRepo  contracts.ReviewRepository
	storeRepo   contracts.StoreRepository
	productRepo contracts.ProductRepository
	validator   validator.Service
	uuid        uuid.Service
}

func NewReviewService(
	reviewRepo contracts.ReviewRepository,
	storeRepo contracts.StoreRepository,
	productRepo contracts.ProductRepository,
	validator validator.Service,
	uuid uuid.Service,
) contracts.ReviewService {
	return &reviewService{
		reviewRepo:  reviewRepo,
		storeRepo:   storeRepo,
		productRepo: productRepo,
		validator:   validator,
		uuid:        uuid,
	}
}
