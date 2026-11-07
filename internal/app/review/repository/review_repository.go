package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/rescuebyte/rescuebyte-be/domain/contracts"
	"github.com/rescuebyte/rescuebyte-be/domain/dto"
	"github.com/rescuebyte/rescuebyte-be/domain/entity"
	"github.com/rescuebyte/rescuebyte-be/domain/errx"
	"github.com/rescuebyte/rescuebyte-be/pkg/helpers"
	"github.com/rescuebyte/rescuebyte-be/pkg/helpers/pgerror"
)

func (r *reviewRepository) CreateReview(ctx context.Context, review entity.Review) error {
	query, args, err := sq.Insert("reviews").
		Columns("id", "user_id", "store_id", "product_id", "content", "rating", "created_at", "updated_at").
		Values(review.ID, review.UserID, review.StoreID, review.ProductID, review.Content, review.Rating, review.CreatedAt, review.UpdatedAt).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.dbGetter(ctx).ExecContext(ctx, query, args...)
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		customErr := r.handleTargetForeignKeyError(err, pgErr, review, "reviewRepository.CreateReview")
		if customErr != nil {
			return customErr
		}
	}
	return err
}

func (r *reviewRepository) GetStoreReviews(ctx context.Context, storeID uuid.UUID, query dto.GetReviewsQuery) ([]entity.Review, int64, error) {
	return r.getReviews(ctx, sq.Eq{"c.store_id": storeID}, query)
}

func (r *reviewRepository) GetProductReviews(ctx context.Context, productID uuid.UUID, query dto.GetReviewsQuery) ([]entity.Review, int64, error) {
	return r.getReviews(ctx, sq.Eq{"c.product_id": productID}, query)
}

func (r *reviewRepository) GetReviewByID(ctx context.Context, reviewID uuid.UUID) (*entity.Review, error) {
	query, args, err := reviewSelect().
		Where(sq.Eq{"c.id": reviewID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}
	var review entity.Review
	err = r.dbGetter(ctx).GetContext(ctx, &review, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errx.ErrReviewNotFound.WithDetails(map[string]any{"id": reviewID}).WithLocation("reviewRepository.GetReviewByID")
		}
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepository) UpdateReview(ctx context.Context, review entity.Review, req dto.UpdateReviewRequest) error {
	query, args, err := sq.Update("reviews").
		Set("content", req.Content).
		Set("rating", req.Rating).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": review.ID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	result, err := r.dbGetter(ctx).ExecContext(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23503" {
				return r.handleTargetForeignKeyError(err, pgErr, review, "reviewRepository.UpdateReview")
			}
			if pgErr.Code == "23514" && pgErr.ConstraintName == "reviews_rating_check" {
				return errx.ErrInvalidReviewRating.WithDetails(map[string]any{"rating": req.Rating}).WithError(err).WithLocation("reviewRepository.UpdateReview")
			}
		}
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	return helpers.CheckRowsAffected(rowsAffected, errx.ErrReviewNotFound.
		WithDetails(map[string]any{"id": review.ID}).
		WithLocation("reviewRepository.UpdateReview"))
}

func (r *reviewRepository) DeleteReview(ctx context.Context, reviewID uuid.UUID) error {
	query, args, err := sq.Delete("reviews").
		Where(sq.Eq{"id": reviewID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	result, err := r.dbGetter(ctx).ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	return helpers.CheckRowsAffected(rowsAffected, errx.ErrReviewNotFound.
		WithDetails(map[string]any{"id": reviewID}).
		WithLocation("reviewRepository.DeleteReview"))
}

func (r *reviewRepository) handleTargetForeignKeyError(
	err error,
	pgErr *pgconn.PgError,
	review entity.Review,
	location string,
) error {
	customErr := pgerror.HandlePgError(*pgErr, []pgerror.PgError{
		{
			Code:           pgerror.ForeignKey,
			ConstraintName: "reviews_store_id_fkey",
			Err: errx.ErrReviewStoreNotFound.
				WithDetails(map[string]any{
					"store_id": review.StoreID,
				}).
				WithError(err).
				WithLocation(location),
		},
		{
			Code:           pgerror.ForeignKey,
			ConstraintName: "reviews_product_id_fkey",
			Err: errx.ErrReviewProductNotFound.
				WithDetails(map[string]any{
					"product_id": review.ProductID,
				}).
				WithError(err).
				WithLocation(location),
		},
		{
			Code:           pgerror.ForeignKey,
			ConstraintName: "reviews_user_id_fkey",
			Err: errx.ErrUserNotFound.
				WithDetails(map[string]any{
					"user_id": review.UserID,
				}).
				WithError(err).
				WithLocation(location),
		},
	})
	if customErr != nil {
		return customErr
	}
	return err
}

func reviewSelect() sq.SelectBuilder {
	return sq.Select(
		"c.id", "c.user_id", "c.store_id", "c.product_id", "c.content", "c.rating", "c.created_at", "c.updated_at",
		"u.name AS author_name", "u.avatar_url", "u.role_id AS author_role",
	).From("reviews c").Join("users u ON u.id = c.user_id")
}

func (r *reviewRepository) getReviews(ctx context.Context, target sq.Eq, query dto.GetReviewsQuery) ([]entity.Review, int64, error) {
	countSQL, countArgs, err := sq.Select("COUNT(*)").
		From("reviews c").
		Where(target).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, 0, err
	}

	var total int64
	if err := r.dbGetter(ctx).GetContext(ctx, &total, countSQL, countArgs...); err != nil {
		return nil, 0, err
	}

	builder := reviewSelect().
		Where(target).
		OrderBy("c.created_at DESC", "c.id DESC")
	if query.Limit > 0 {
		builder = builder.Limit(query.Limit).Offset(query.Offset)
	}

	listSQL, listArgs, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, 0, err
	}
	reviews := make([]entity.Review, 0)
	if err := r.dbGetter(ctx).SelectContext(ctx, &reviews, listSQL, listArgs...); err != nil {
		return nil, 0, err
	}
	return reviews, total, nil
}

var _ contracts.ReviewRepository = (*reviewRepository)(nil)
