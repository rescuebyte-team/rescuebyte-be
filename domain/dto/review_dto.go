package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/rescuebyte/rescuebyte-be/domain/entity"
	"github.com/rescuebyte/rescuebyte-be/domain/helpers/pagination"
)

type CreateReviewRequest struct {
	Content string `json:"content" validate:"required,min=1,max=2000"`
	Rating  int    `json:"rating" validate:"required,gte=1,lte=5"`
}

type UpdateReviewRequest struct {
	Content string `json:"content" validate:"required,min=1,max=2000"`
	Rating  int    `json:"rating" validate:"required,gte=1,lte=5"`
}

type ReviewAuthorResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	AvatarURL *string   `json:"avatar_url,omitempty"`
	Role      string    `json:"role"`
}

type ReviewResponse struct {
	ID        uuid.UUID            `json:"id"`
	UserID    uuid.UUID            `json:"user_id"`
	StoreID   *uuid.UUID           `json:"store_id,omitempty"`
	ProductID *uuid.UUID           `json:"product_id,omitempty"`
	Content   string               `json:"content"`
	Rating    int                  `json:"rating"`
	Author    ReviewAuthorResponse `json:"author"`
	CreatedAt string               `json:"created_at"`
	UpdatedAt string               `json:"updated_at"`
}

type GetReviewsQuery struct {
	Limit  uint64 `query:"limit"`
	Page   uint64 `query:"page" validate:"omitempty,gte=1"`
	Offset uint64 `query:"-"`
}

type GetReviewsResponse struct {
	Reviews []ReviewResponse    `json:"reviews"`
	Meta    pagination.Response `json:"meta"`
}

func FormatToReviewResponse(review entity.Review) ReviewResponse {
	return ReviewResponse{
		ID:        review.ID,
		UserID:    review.UserID,
		StoreID:   review.StoreID,
		ProductID: review.ProductID,
		Content:   review.Content,
		Rating:    review.Rating,
		Author: ReviewAuthorResponse{
			ID:        review.UserID,
			Name:      review.AuthorName,
			AvatarURL: review.AvatarURL,
			Role:      review.AuthorRole.String(),
		},
		CreatedAt: review.CreatedAt.Format(time.RFC3339),
		UpdatedAt: review.UpdatedAt.Format(time.RFC3339),
	}
}
