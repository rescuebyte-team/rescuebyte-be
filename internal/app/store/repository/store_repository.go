package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"github.com/rescuebyte/rescuebyte-be/domain/entity"
	"github.com/rescuebyte/rescuebyte-be/domain/errx"
)

var storeSelectFields = []string{
	"id", "user_id", "name", "profile_picture_url", "description",
	"whatsapp_number", "email", "address", "city", "district",
	"status", "operational_start::text", "operational_end::text", "operational_days",
	"latitude", "longitude", "created_at", "updated_at",
}

func (r *storeRepository) CreateStore(ctx context.Context, store entity.Store) error {
	query, args, err := sq.
		Insert("stores").
		Columns(
			"id", "user_id", "name", "profile_picture_url", "description",
			"whatsapp_number", "email", "address", "city", "district",
			"status", "operational_start", "operational_end", "operational_days",
			"latitude", "longitude", "created_at", "updated_at",
		).
		Values(
			store.ID, store.UserID, store.Name, store.ProfilePictureURL, store.Description,
			store.WhatsAppNumber, store.Email, store.Address, store.City, store.District,
			store.Status, store.OperationalStart, store.OperationalEnd, store.OperationalDays,
			store.Latitude, store.Longitude, store.CreatedAt, store.UpdatedAt,
		).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.dbGetter(ctx).ExecContext(ctx, query, args...)
	return err
}

func (r *storeRepository) GetStoreByID(ctx context.Context, id uuid.UUID) (*entity.Store, error) {
	query, args, err := sq.
		Select(storeSelectFields...).
		From("stores").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var store entity.Store
	err = r.dbGetter(ctx).GetContext(ctx, &store, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errx.ErrStoreNotFound
		}
		return nil, err
	}

	return &store, nil
}

func (r *storeRepository) GetStoreByUserID(ctx context.Context, userID uuid.UUID) (*entity.Store, error) {
	query, args, err := sq.
		Select(storeSelectFields...).
		From("stores").
		Where(sq.Eq{"user_id": userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var store entity.Store
	err = r.dbGetter(ctx).GetContext(ctx, &store, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Return nil if store does not exist yet
		}
		return nil, err
	}

	return &store, nil
}

func (r *storeRepository) UpdateStoreProfile(ctx context.Context, store entity.Store) error {
	query, args, err := sq.
		Update("stores").
		Set("name", store.Name).
		Set("profile_picture_url", store.ProfilePictureURL).
		Set("description", store.Description).
		Set("whatsapp_number", store.WhatsAppNumber).
		Set("email", store.Email).
		Set("address", store.Address).
		Set("city", store.City).
		Set("district", store.District).
		Set("latitude", store.Latitude).
		Set("longitude", store.Longitude).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": store.ID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.dbGetter(ctx).ExecContext(ctx, query, args...)
	return err
}

func (r *storeRepository) UpdateStoreOperational(ctx context.Context, store entity.Store) error {
	query, args, err := sq.
		Update("stores").
		Set("status", store.Status).
		Set("operational_start", store.OperationalStart).
		Set("operational_end", store.OperationalEnd).
		Set("operational_days", store.OperationalDays).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": store.ID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.dbGetter(ctx).ExecContext(ctx, query, args...)
	return err
}

func (r *storeRepository) UpdateStoreStatus(ctx context.Context, storeID uuid.UUID, status string) error {
	query, args, err := sq.
		Update("stores").
		Set("status", status).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": storeID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.dbGetter(ctx).ExecContext(ctx, query, args...)
	return err
}
