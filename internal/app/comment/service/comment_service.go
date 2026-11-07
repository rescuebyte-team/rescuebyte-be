package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/rescuebyte/rescuebyte-be/domain/dto"
	"github.com/rescuebyte/rescuebyte-be/domain/entity"
	"github.com/rescuebyte/rescuebyte-be/domain/errx"
	"github.com/rescuebyte/rescuebyte-be/domain/helpers/pagination"
)

func (s *commentService) CreateStoreComment(ctx context.Context, userID, storeID uuid.UUID, req dto.CreateCommentRequest) (*dto.CommentResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}
	if _, err := s.storeRepo.GetStoreByID(ctx, storeID); err != nil {
		if errors.Is(err, errx.ErrStoreNotFound) {
			return nil, errx.ErrCommentStoreNotFound.WithDetails(map[string]any{"store_id": storeID})
		}
		return nil, err
	}

	comment, err := s.newComment(userID, req)
	if err != nil {
		return nil, err
	}
	comment.StoreID = &storeID
	if err := s.commentRepo.CreateComment(ctx, comment); err != nil {
		return nil, err
	}
	saved, err := s.commentRepo.GetCommentByID(ctx, comment.ID)
	if err != nil {
		return nil, err
	}
	res := dto.FormatToCommentResponse(*saved)
	return &res, nil
}

func (s *commentService) GetStoreComments(ctx context.Context, storeID uuid.UUID, query dto.GetCommentsQuery) (dto.GetCommentsResponse, error) {
	if _, err := s.storeRepo.GetStoreByID(ctx, storeID); err != nil {
		if errors.Is(err, errx.ErrStoreNotFound) {
			return dto.GetCommentsResponse{}, errx.ErrCommentStoreNotFound.WithDetails(map[string]any{"store_id": storeID})
		}
		return dto.GetCommentsResponse{}, err
	}
	return s.list(ctx, true, storeID, query)
}

func (s *commentService) CreateProductComment(ctx context.Context, userID, productID uuid.UUID, req dto.CreateCommentRequest) (*dto.CommentResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}
	if _, err := s.productRepo.GetProductByID(ctx, productID); err != nil {
		if errors.Is(err, errx.ErrProductNotFound) {
			return nil, errx.ErrCommentProductNotFound.WithDetails(map[string]any{"product_id": productID})
		}
		return nil, err
	}

	comment, err := s.newComment(userID, req)
	if err != nil {
		return nil, err
	}
	comment.ProductID = &productID
	if err := s.commentRepo.CreateComment(ctx, comment); err != nil {
		return nil, err
	}
	saved, err := s.commentRepo.GetCommentByID(ctx, comment.ID)
	if err != nil {
		return nil, err
	}
	res := dto.FormatToCommentResponse(*saved)
	return &res, nil
}

func (s *commentService) GetProductComments(ctx context.Context, productID uuid.UUID, query dto.GetCommentsQuery) (dto.GetCommentsResponse, error) {
	if _, err := s.productRepo.GetProductByID(ctx, productID); err != nil {
		if errors.Is(err, errx.ErrProductNotFound) {
			return dto.GetCommentsResponse{}, errx.ErrCommentProductNotFound.WithDetails(map[string]any{"product_id": productID})
		}
		return dto.GetCommentsResponse{}, err
	}
	return s.list(ctx, false, productID, query)
}

func (s *commentService) newComment(userID uuid.UUID, req dto.CreateCommentRequest) (entity.Comment, error) {
	id, err := s.uuid.NewV7()
	if err != nil {
		return entity.Comment{}, err
	}
	now := time.Now()
	return entity.Comment{
		ID:        id,
		UserID:    userID,
		Content:   req.Content,
		Rating:    req.Rating,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (s *commentService) list(ctx context.Context, isStore bool, targetID uuid.UUID, query dto.GetCommentsQuery) (dto.GetCommentsResponse, error) {
	if err := s.validator.Validate(query); err != nil {
		return dto.GetCommentsResponse{}, err
	}
	limit, page, offset := query.Limit, query.Page, query.Offset
	pagination.Normalize(&limit, &page, &offset, "comment")
	query.Limit, query.Page, query.Offset = limit, page, offset

	var comments []entity.Comment
	var total int64
	var err error
	if isStore {
		comments, total, err = s.commentRepo.GetStoreComments(ctx, targetID, query)
	} else {
		comments, total, err = s.commentRepo.GetProductComments(ctx, targetID, query)
	}
	if err != nil {
		return dto.GetCommentsResponse{}, err
	}

	items := make([]dto.CommentResponse, 0, len(comments))
	for _, comment := range comments {
		items = append(items, dto.FormatToCommentResponse(comment))
	}
	return dto.GetCommentsResponse{
		Comments: items,
		Meta:     pagination.NewPaginationResponse(total, page, limit),
	}, nil
}
