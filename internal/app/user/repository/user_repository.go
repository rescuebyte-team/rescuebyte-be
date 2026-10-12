package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"github.com/rescuebyte/rescuebyte-be/domain/contracts"
	"github.com/rescuebyte/rescuebyte-be/domain/entity"
	"github.com/rescuebyte/rescuebyte-be/domain/enums"
	"github.com/rescuebyte/rescuebyte-be/domain/errx"
	"github.com/rescuebyte/rescuebyte-be/internal/infra/database"
)

type userRepository struct {
	dbGetter database.DBGetter
}

func NewUserRepository(dbGetter database.DBGetter) contracts.UserRepository {
	return &userRepository{dbGetter: dbGetter}
}

var userSelectFields = []string{
	"id", "email", "password_hash", "name", "avatar_url",
	"role", "provider", "provider_id", "created_at", "updated_at",
}

func (r *userRepository) CreateUser(ctx context.Context, user entity.User) error {
	query, args, err := sq.
		Insert("users").
		Columns("id", "email", "password_hash", "name", "avatar_url", "role", "provider", "provider_id", "created_at", "updated_at").
		Values(user.ID, user.Email, user.PasswordHash, user.Name, user.AvatarURL, user.Role, user.Provider, user.ProviderID, user.CreatedAt, user.UpdatedAt).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.dbGetter(ctx).ExecContext(ctx, query, args...)
	return err
}

func (r *userRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	query, args, err := sq.
		Select(userSelectFields...).
		From("users").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var user entity.User
	err = r.dbGetter(ctx).GetContext(ctx, &user, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errx.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	query, args, err := sq.
		Select(userSelectFields...).
		From("users").
		Where(sq.Eq{"email": email}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var user entity.User
	err = r.dbGetter(ctx).GetContext(ctx, &user, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // not found, return nil without error
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetUserByProviderID(ctx context.Context, provider enums.AuthProvider, providerID string) (*entity.User, error) {
	query, args, err := sq.
		Select(userSelectFields...).
		From("users").
		Where(sq.Eq{"provider": provider, "provider_id": providerID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var user entity.User
	err = r.dbGetter(ctx).GetContext(ctx, &user, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user entity.User) error {
	query, args, err := sq.
		Update("users").
		Set("name", user.Name).
		Set("avatar_url", user.AvatarURL).
		Set("role", user.Role).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": user.ID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.dbGetter(ctx).ExecContext(ctx, query, args...)
	return err
}

func (r *userRepository) CreateRefreshToken(ctx context.Context, token entity.RefreshToken) error {
	query, args, err := sq.
		Insert("refresh_tokens").
		Columns("id", "user_id", "token", "expires_at", "revoked_at", "created_at").
		Values(token.ID, token.UserID, token.Token, token.ExpiresAt, token.RevokedAt, token.CreatedAt).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.dbGetter(ctx).ExecContext(ctx, query, args...)
	return err
}

func (r *userRepository) GetRefreshToken(ctx context.Context, token string) (*entity.RefreshToken, error) {
	query, args, err := sq.
		Select("id", "user_id", "token", "expires_at", "revoked_at", "created_at").
		From("refresh_tokens").
		Where(sq.Eq{"token": token}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var rt entity.RefreshToken
	err = r.dbGetter(ctx).GetContext(ctx, &rt, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errx.ErrInvalidRefreshToken
		}
		return nil, err
	}

	return &rt, nil
}

func (r *userRepository) RevokeRefreshToken(ctx context.Context, token string) error {
	now := time.Now()
	query, args, err := sq.
		Update("refresh_tokens").
		Set("revoked_at", now).
		Where(sq.Eq{"token": token}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.dbGetter(ctx).ExecContext(ctx, query, args...)
	return err
}

func (r *userRepository) RevokeAllUserRefreshTokens(ctx context.Context, userID uuid.UUID) error {
	now := time.Now()
	query, args, err := sq.
		Update("refresh_tokens").
		Set("revoked_at", now).
		Where(sq.Eq{"user_id": userID, "revoked_at": nil}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.dbGetter(ctx).ExecContext(ctx, query, args...)
	return err
}
