package contracts

import (
	"context"

	"github.com/google/uuid"

	"github.com/rescuebyte/rescuebyte-be/domain/dto"
	"github.com/rescuebyte/rescuebyte-be/domain/entity"
	"github.com/rescuebyte/rescuebyte-be/domain/enums"
)

type ReviewRepository interface {
	CreateReview(ctx context.Context, review entity.Review) error
	GetReviewByID(ctx context.Context, reviewID uuid.UUID) (*entity.Review, error)
	UpdateReview(ctx context.Context, review entity.Review, req dto.UpdateReviewRequest) error
	DeleteReview(ctx context.Context, reviewID uuid.UUID) error
	GetStoreReviews(ctx context.Context, storeID uuid.UUID, query dto.GetReviewsQuery) ([]entity.Review, int64, error)
	GetProductReviews(ctx context.Context, productID uuid.UUID, query dto.GetReviewsQuery) ([]entity.Review, int64, error)
}

type ReviewService interface {
	CreateStoreReview(ctx context.Context, userID, storeID uuid.UUID, req dto.CreateReviewRequest) (*dto.ReviewResponse, error)
	GetStoreReviews(ctx context.Context, storeID uuid.UUID, query dto.GetReviewsQuery) (dto.GetReviewsResponse, error)
	CreateProductReview(ctx context.Context, userID, productID uuid.UUID, req dto.CreateReviewRequest) (*dto.ReviewResponse, error)
	GetProductReviews(ctx context.Context, productID uuid.UUID, query dto.GetReviewsQuery) (dto.GetReviewsResponse, error)
	UpdateReview(ctx context.Context, userID uuid.UUID, roleID enums.UserRoleIdx, reviewID uuid.UUID, req dto.UpdateReviewRequest) (*dto.ReviewResponse, error)
	DeleteReview(ctx context.Context, userID uuid.UUID, roleID enums.UserRoleIdx, reviewID uuid.UUID) error
}
