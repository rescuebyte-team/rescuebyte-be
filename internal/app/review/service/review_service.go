package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/rescuebyte/rescuebyte-be/domain/dto"
	"github.com/rescuebyte/rescuebyte-be/domain/entity"
	"github.com/rescuebyte/rescuebyte-be/domain/enums"
	"github.com/rescuebyte/rescuebyte-be/domain/errx"
	"github.com/rescuebyte/rescuebyte-be/domain/helpers/pagination"
)

func (s *reviewService) CreateStoreReview(ctx context.Context, userID, storeID uuid.UUID, req dto.CreateReviewRequest) (*dto.ReviewResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}
	if _, err := s.storeRepo.GetStoreByID(ctx, storeID); err != nil {
		if errors.Is(err, errx.ErrStoreNotFound) {
			return nil, errx.ErrReviewStoreNotFound.WithDetails(map[string]any{"store_id": storeID})
		}
		return nil, err
	}

	review, err := s.newReview(userID, req)
	if err != nil {
		return nil, err
	}
	review.StoreID = &storeID
	if err := s.reviewRepo.CreateReview(ctx, review); err != nil {
		return nil, err
	}
	saved, err := s.reviewRepo.GetReviewByID(ctx, review.ID)
	if err != nil {
		return nil, err
	}
	res := dto.FormatToReviewResponse(*saved)
	return &res, nil
}

func (s *reviewService) GetStoreReviews(ctx context.Context, storeID uuid.UUID, query dto.GetReviewsQuery) (dto.GetReviewsResponse, error) {
	if _, err := s.storeRepo.GetStoreByID(ctx, storeID); err != nil {
		if errors.Is(err, errx.ErrStoreNotFound) {
			return dto.GetReviewsResponse{}, errx.ErrReviewStoreNotFound.WithDetails(map[string]any{"store_id": storeID})
		}
		return dto.GetReviewsResponse{}, err
	}
	return s.list(ctx, true, storeID, query)
}

func (s *reviewService) CreateProductReview(ctx context.Context, userID, productID uuid.UUID, req dto.CreateReviewRequest) (*dto.ReviewResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}
	if _, err := s.productRepo.GetProductByID(ctx, productID); err != nil {
		if errors.Is(err, errx.ErrProductNotFound) {
			return nil, errx.ErrReviewProductNotFound.WithDetails(map[string]any{"product_id": productID})
		}
		return nil, err
	}

	review, err := s.newReview(userID, req)
	if err != nil {
		return nil, err
	}
	review.ProductID = &productID
	if err := s.reviewRepo.CreateReview(ctx, review); err != nil {
		return nil, err
	}
	saved, err := s.reviewRepo.GetReviewByID(ctx, review.ID)
	if err != nil {
		return nil, err
	}
	res := dto.FormatToReviewResponse(*saved)
	return &res, nil
}

func (s *reviewService) GetProductReviews(ctx context.Context, productID uuid.UUID, query dto.GetReviewsQuery) (dto.GetReviewsResponse, error) {
	if _, err := s.productRepo.GetProductByID(ctx, productID); err != nil {
		if errors.Is(err, errx.ErrProductNotFound) {
			return dto.GetReviewsResponse{}, errx.ErrReviewProductNotFound.WithDetails(map[string]any{"product_id": productID})
		}
		return dto.GetReviewsResponse{}, err
	}
	return s.list(ctx, false, productID, query)
}

func (s *reviewService) UpdateReview(ctx context.Context, userID uuid.UUID, roleID enums.UserRoleIdx, reviewID uuid.UUID, req dto.UpdateReviewRequest) (*dto.ReviewResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}
	review, err := s.reviewRepo.GetReviewByID(ctx, reviewID)
	if err != nil {
		return nil, err
	}
	if roleID != enums.RoleAdminIdx && review.UserID != userID {
		return nil, errx.ErrForbidden
	}

	if err := s.reviewRepo.UpdateReview(ctx, *review, req); err != nil {
		return nil, err
	}
	updated, err := s.reviewRepo.GetReviewByID(ctx, reviewID)
	if err != nil {
		return nil, err
	}
	res := dto.FormatToReviewResponse(*updated)
	return &res, nil
}

func (s *reviewService) DeleteReview(ctx context.Context, userID uuid.UUID, roleID enums.UserRoleIdx, reviewID uuid.UUID) error {
	review, err := s.reviewRepo.GetReviewByID(ctx, reviewID)
	if err != nil {
		return err
	}
	if roleID != enums.RoleAdminIdx && review.UserID != userID {
		return errx.ErrForbidden
	}
	return s.reviewRepo.DeleteReview(ctx, reviewID)
}

func (s *reviewService) newReview(userID uuid.UUID, req dto.CreateReviewRequest) (entity.Review, error) {
	id, err := s.uuid.NewV7()
	if err != nil {
		return entity.Review{}, err
	}
	now := time.Now()
	return entity.Review{
		ID:        id,
		UserID:    userID,
		Content:   req.Content,
		Rating:    req.Rating,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (s *reviewService) list(ctx context.Context, isStore bool, targetID uuid.UUID, query dto.GetReviewsQuery) (dto.GetReviewsResponse, error) {
	if err := s.validator.Validate(query); err != nil {
		return dto.GetReviewsResponse{}, err
	}
	limit, page, offset := query.Limit, query.Page, query.Offset
	pagination.Normalize(&limit, &page, &offset, "review")
	query.Limit, query.Page, query.Offset = limit, page, offset

	var reviews []entity.Review
	var total int64
	var err error
	if isStore {
		reviews, total, err = s.reviewRepo.GetStoreReviews(ctx, targetID, query)
	} else {
		reviews, total, err = s.reviewRepo.GetProductReviews(ctx, targetID, query)
	}
	if err != nil {
		return dto.GetReviewsResponse{}, err
	}

	items := make([]dto.ReviewResponse, 0, len(reviews))
	for _, review := range reviews {
		items = append(items, dto.FormatToReviewResponse(review))
	}
	return dto.GetReviewsResponse{
		Reviews: items,
		Meta:    pagination.NewPaginationResponse(total, page, limit),
	}, nil
}
