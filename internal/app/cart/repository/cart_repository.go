package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"github.com/rescuebyte/rescuebyte-be/domain/entity"
)

func (r *cartRepository) GetOrCreateCartByUserID(ctx context.Context, userID uuid.UUID) (*entity.Cart, error) {
	query, args, err := sq.
		Select("id", "user_id", "created_at", "updated_at").
		From("carts").
		Where(sq.Eq{"user_id": userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var cart entity.Cart
	err = r.dbGetter(ctx).GetContext(ctx, &cart, query, args...)
	if err == nil {
		return &cart, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	// Insert new cart
	newCartID, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	now := time.Now()

	insertQuery, insertArgs, err := sq.
		Insert("carts").
		Columns("id", "user_id", "created_at", "updated_at").
		Values(newCartID, userID, now, now).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	_, err = r.dbGetter(ctx).ExecContext(ctx, insertQuery, insertArgs...)
	if err != nil {
		// If another request inserted concurrently, try selecting again
		err2 := r.dbGetter(ctx).GetContext(ctx, &cart, query, args...)
		if err2 == nil {
			return &cart, nil
		}
		return nil, err
	}

	return &entity.Cart{
		ID:        newCartID,
		UserID:    userID,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (r *cartRepository) GetCartItemByID(ctx context.Context, id uuid.UUID) (*entity.CartItem, error) {
	query, args, err := sq.
		Select("id", "cart_id", "product_variant_id", "quantity", "created_at", "updated_at").
		From("cart_items").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var item entity.CartItem
	err = r.dbGetter(ctx).GetContext(ctx, &item, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (r *cartRepository) GetCartItemByVariantID(ctx context.Context, cartID, variantID uuid.UUID) (*entity.CartItem, error) {
	query, args, err := sq.
		Select("id", "cart_id", "product_variant_id", "quantity", "created_at", "updated_at").
		From("cart_items").
		Where(sq.Eq{
			"cart_id":            cartID,
			"product_variant_id": variantID,
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var item entity.CartItem
	err = r.dbGetter(ctx).GetContext(ctx, &item, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (r *cartRepository) GetCartItemDetails(ctx context.Context, cartID uuid.UUID) ([]entity.CartItemDetail, error) {
	query, args, err := sq.
		Select(
			"ci.id",
			"ci.cart_id",
			"ci.product_variant_id",
			"ci.quantity",
			"p.id AS product_id",
			"p.name AS product_name",
			"p.description AS product_description",
			"p.main_image_url AS product_image_url",
			"pv.name AS variant_name",
			"pv.normal_price",
			"pv.rescue_price",
			"pv.current_stock",
			"pv.is_active AS variant_active",
			"p.is_active AS product_active",
			"s.id AS store_id",
			"s.name AS store_name",
			"s.address AS store_address",
			"s.city AS store_city",
			"s.district AS store_district",
			"s.latitude AS store_latitude",
			"s.longitude AS store_longitude",
			"COALESCE(s.operational_start::text, '') AS operational_start",
			"COALESCE(s.operational_end::text, '') AS operational_end",
			"ci.created_at",
			"ci.updated_at",
		).
		From("cart_items ci").
		Join("product_variants pv ON pv.id = ci.product_variant_id").
		Join("products p ON p.id = pv.product_id").
		Join("stores s ON s.id = p.store_id").
		Where(sq.Eq{"ci.cart_id": cartID}).
		OrderBy("s.name ASC", "ci.created_at DESC").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var items []entity.CartItemDetail
	err = r.dbGetter(ctx).SelectContext(ctx, &items, query, args...)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *cartRepository) GetCartItemsCount(ctx context.Context, cartID uuid.UUID) (int, error) {
	query, args, err := sq.
		Select("COALESCE(SUM(quantity), 0)").
		From("cart_items").
		Where(sq.Eq{"cart_id": cartID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return 0, err
	}

	var count int
	err = r.dbGetter(ctx).GetContext(ctx, &count, query, args...)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *cartRepository) UpsertCartItem(ctx context.Context, item entity.CartItem) error {
	rawQuery := `
		INSERT INTO cart_items (id, cart_id, product_variant_id, quantity, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (cart_id, product_variant_id)
		DO UPDATE SET
			quantity = cart_items.quantity + EXCLUDED.quantity,
			updated_at = EXCLUDED.updated_at
	`
	_, err := r.dbGetter(ctx).ExecContext(ctx, rawQuery, item.ID, item.CartID, item.ProductVariantID, item.Quantity, item.CreatedAt, item.UpdatedAt)
	return err
}

func (r *cartRepository) UpdateCartItemQuantity(ctx context.Context, id uuid.UUID, quantity int) error {
	query, args, err := sq.
		Update("cart_items").
		Set("quantity", quantity).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.dbGetter(ctx).ExecContext(ctx, query, args...)
	return err
}

func (r *cartRepository) DeleteCartItem(ctx context.Context, id uuid.UUID) error {
	query, args, err := sq.
		Delete("cart_items").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.dbGetter(ctx).ExecContext(ctx, query, args...)
	return err
}

func (r *cartRepository) DeleteCartItemsByStoreID(ctx context.Context, cartID, storeID uuid.UUID) error {
	rawQuery := `
		DELETE FROM cart_items ci
		USING product_variants pv, products p
		WHERE ci.product_variant_id = pv.id
		  AND pv.product_id = p.id
		  AND ci.cart_id = $1
		  AND p.store_id = $2
	`
	_, err := r.dbGetter(ctx).ExecContext(ctx, rawQuery, cartID, storeID)
	return err
}

func (r *cartRepository) ClearCart(ctx context.Context, cartID uuid.UUID) error {
	query, args, err := sq.
		Delete("cart_items").
		Where(sq.Eq{"cart_id": cartID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.dbGetter(ctx).ExecContext(ctx, query, args...)
	return err
}
