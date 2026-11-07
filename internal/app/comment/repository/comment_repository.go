package repository

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/rescuebyte/rescuebyte-be/domain/contracts"
	"github.com/rescuebyte/rescuebyte-be/domain/dto"
	"github.com/rescuebyte/rescuebyte-be/domain/entity"
	"github.com/rescuebyte/rescuebyte-be/domain/errx"
)

func (r *commentRepository) CreateComment(ctx context.Context, comment entity.Comment) error {
	query, args, err := sq.Insert("comments").
		Columns("id", "user_id", "store_id", "product_id", "content", "rating", "created_at", "updated_at").
		Values(comment.ID, comment.UserID, comment.StoreID, comment.ProductID, comment.Content, comment.Rating, comment.CreatedAt, comment.UpdatedAt).
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
	if errors.As(err, &pgErr) && pgErr.Code == "23503" {
		switch pgErr.ConstraintName {
		case "comments_store_id_fkey":
			return errx.ErrCommentStoreNotFound.WithDetails(map[string]any{"store_id": comment.StoreID}).WithError(err).WithLocation("commentRepository.CreateComment")
		case "comments_product_id_fkey":
			return errx.ErrCommentProductNotFound.WithDetails(map[string]any{"product_id": comment.ProductID}).WithError(err).WithLocation("commentRepository.CreateComment")
		case "comments_user_id_fkey":
			return errx.ErrUserNotFound.WithDetails(map[string]any{"user_id": comment.UserID}).WithError(err).WithLocation("commentRepository.CreateComment")
		}
	}
	return err
}

func (r *commentRepository) GetStoreComments(ctx context.Context, storeID uuid.UUID, query dto.GetCommentsQuery) ([]entity.Comment, int64, error) {
	return r.getComments(ctx, sq.Eq{"c.store_id": storeID}, query)
}

func (r *commentRepository) GetProductComments(ctx context.Context, productID uuid.UUID, query dto.GetCommentsQuery) ([]entity.Comment, int64, error) {
	return r.getComments(ctx, sq.Eq{"c.product_id": productID}, query)
}

func (r *commentRepository) GetCommentByID(ctx context.Context, commentID uuid.UUID) (*entity.Comment, error) {
	query, args, err := commentSelect().
		Where(sq.Eq{"c.id": commentID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}
	var comment entity.Comment
	err = r.dbGetter(ctx).GetContext(ctx, &comment, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errx.ErrNotFound
		}
		return nil, err
	}
	return &comment, nil
}

func commentSelect() sq.SelectBuilder {
	return sq.Select(
		"c.id", "c.user_id", "c.store_id", "c.product_id", "c.content", "c.rating", "c.created_at", "c.updated_at",
		"u.name AS author_name", "u.avatar_url", "u.role_id AS author_role",
	).From("comments c").Join("users u ON u.id = c.user_id")
}

func (r *commentRepository) getComments(ctx context.Context, target sq.Eq, query dto.GetCommentsQuery) ([]entity.Comment, int64, error) {
	countSQL, countArgs, err := sq.Select("COUNT(*)").
		From("comments c").
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

	builder := commentSelect().
		Where(target).
		OrderBy("c.created_at DESC", "c.id DESC")
	if query.Limit > 0 {
		builder = builder.Limit(query.Limit).Offset(query.Offset)
	}

	listSQL, listArgs, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, 0, err
	}
	comments := make([]entity.Comment, 0)
	if err := r.dbGetter(ctx).SelectContext(ctx, &comments, listSQL, listArgs...); err != nil {
		return nil, 0, err
	}
	return comments, total, nil
}

var _ contracts.CommentRepository = (*commentRepository)(nil)
