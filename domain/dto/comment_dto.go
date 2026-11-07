package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/rescuebyte/rescuebyte-be/domain/entity"
	"github.com/rescuebyte/rescuebyte-be/domain/helpers/pagination"
)

type CreateCommentRequest struct {
	Content string `json:"content" validate:"required,min=1,max=2000"`
	Rating  int    `json:"rating" validate:"required,gte=1,lte=5"`
}

type CommentAuthorResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	AvatarURL *string   `json:"avatar_url,omitempty"`
	Role      string    `json:"role"`
}

type CommentResponse struct {
	ID        uuid.UUID             `json:"id"`
	UserID    uuid.UUID             `json:"user_id"`
	StoreID   *uuid.UUID            `json:"store_id,omitempty"`
	ProductID *uuid.UUID            `json:"product_id,omitempty"`
	Content   string                `json:"content"`
	Rating    int                   `json:"rating"`
	Author    CommentAuthorResponse `json:"author"`
	CreatedAt string                `json:"created_at"`
	UpdatedAt string                `json:"updated_at"`
}

type GetCommentsQuery struct {
	Limit  uint64 `query:"limit"`
	Page   uint64 `query:"page" validate:"omitempty,gte=1"`
	Offset uint64 `query:"-"`
}

type GetCommentsResponse struct {
	Comments []CommentResponse   `json:"comments"`
	Meta     pagination.Response `json:"meta"`
}

func FormatToCommentResponse(comment entity.Comment) CommentResponse {
	return CommentResponse{
		ID:        comment.ID,
		UserID:    comment.UserID,
		StoreID:   comment.StoreID,
		ProductID: comment.ProductID,
		Content:   comment.Content,
		Rating:    comment.Rating,
		Author: CommentAuthorResponse{
			ID:        comment.UserID,
			Name:      comment.AuthorName,
			AvatarURL: comment.AvatarURL,
			Role:      comment.AuthorRole.String(),
		},
		CreatedAt: comment.CreatedAt.Format(time.RFC3339),
		UpdatedAt: comment.UpdatedAt.Format(time.RFC3339),
	}
}
