package service

import (
	"github.com/rescuebyte/rescuebyte-be/domain/contracts"
	"github.com/rescuebyte/rescuebyte-be/pkg/uuid"
	"github.com/rescuebyte/rescuebyte-be/pkg/validator"
)

type commentService struct {
	commentRepo contracts.CommentRepository
	storeRepo   contracts.StoreRepository
	productRepo contracts.ProductRepository
	validator   validator.Service
	uuid        uuid.Service
}

func NewCommentService(
	commentRepo contracts.CommentRepository,
	storeRepo contracts.StoreRepository,
	productRepo contracts.ProductRepository,
	validator validator.Service,
	uuid uuid.Service,
) contracts.CommentService {
	return &commentService{
		commentRepo: commentRepo,
		storeRepo:   storeRepo,
		productRepo: productRepo,
		validator:   validator,
		uuid:        uuid,
	}
}
