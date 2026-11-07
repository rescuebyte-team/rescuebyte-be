package contracts

import (
	"context"

	"github.com/google/uuid"

	"github.com/rescuebyte/rescuebyte-be/domain/dto"
	"github.com/rescuebyte/rescuebyte-be/domain/entity"
	"github.com/rescuebyte/rescuebyte-be/domain/enums"
)

type CommentRepository interface {
	CreateComment(ctx context.Context, comment entity.Comment) error
	GetCommentByID(ctx context.Context, commentID uuid.UUID) (*entity.Comment, error)
	UpdateComment(ctx context.Context, comment entity.Comment, req dto.UpdateCommentRequest) error
	DeleteComment(ctx context.Context, commentID uuid.UUID) error
	GetStoreComments(ctx context.Context, storeID uuid.UUID, query dto.GetCommentsQuery) ([]entity.Comment, int64, error)
	GetProductComments(ctx context.Context, productID uuid.UUID, query dto.GetCommentsQuery) ([]entity.Comment, int64, error)
}

type CommentService interface {
	CreateStoreComment(ctx context.Context, userID, storeID uuid.UUID, req dto.CreateCommentRequest) (*dto.CommentResponse, error)
	GetStoreComments(ctx context.Context, storeID uuid.UUID, query dto.GetCommentsQuery) (dto.GetCommentsResponse, error)
	CreateProductComment(ctx context.Context, userID, productID uuid.UUID, req dto.CreateCommentRequest) (*dto.CommentResponse, error)
	GetProductComments(ctx context.Context, productID uuid.UUID, query dto.GetCommentsQuery) (dto.GetCommentsResponse, error)
	UpdateComment(ctx context.Context, userID uuid.UUID, roleID enums.UserRoleIdx, commentID uuid.UUID, req dto.UpdateCommentRequest) (*dto.CommentResponse, error)
	DeleteComment(ctx context.Context, userID uuid.UUID, roleID enums.UserRoleIdx, commentID uuid.UUID) error
}
