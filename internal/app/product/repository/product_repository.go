package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"github.com/rescuebyte/rescuebyte-be/domain/dto"
	"github.com/rescuebyte/rescuebyte-be/domain/entity"
	"github.com/rescuebyte/rescuebyte-be/domain/errx"
)

var productSelectFields = []string{
	"id", "store_id", "name", "description", "main_image_url", "is_active", "created_at", "updated_at",
}

var variantSelectFields = []string{
	"id", "product_id", "name", "sequence", "normal_price", "rescue_price",
	"initial_stock", "current_stock", "is_active", "created_at", "updated_at",
}

var imageSelectFields = []string{
	"id", "product_id", "image_url", "display_order", "created_at",
}

func (r *productRepository) CreateProduct(ctx context.Context, product entity.Product) error {
	query, args, err := sq.
		Insert("products").
		Columns("id", "store_id", "name", "description", "main_image_url", "is_active", "created_at", "updated_at").
		Values(product.ID, product.StoreID, product.Name, product.Description, product.MainImageURL, product.IsActive, product.CreatedAt, product.UpdatedAt).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.dbGetter(ctx).ExecContext(ctx, query, args...)
	return err
}

func (r *productRepository) CreateProductVariants(ctx context.Context, variants []entity.ProductVariant) error {
	if len(variants) == 0 {
		return nil
	}

	builder := sq.
		Insert("product_variants").
		Columns("id", "product_id", "name", "sequence", "normal_price", "rescue_price", "initial_stock", "current_stock", "is_active", "created_at", "updated_at")

	for _, v := range variants {
		builder = builder.Values(v.ID, v.ProductID, v.Name, v.Sequence, v.NormalPrice, v.RescuePrice, v.InitialStock, v.CurrentStock, v.IsActive, v.CreatedAt, v.UpdatedAt)
	}

	query, args, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}

	_, err = r.dbGetter(ctx).ExecContext(ctx, query, args...)
	return err
}

func (r *productRepository) CreateProductImages(ctx context.Context, images []entity.ProductImage) error {
	if len(images) == 0 {
		return nil
	}

	builder := sq.
		Insert("product_images").
		Columns("id", "product_id", "image_url", "display_order", "created_at")

	for _, img := range images {
		builder = builder.Values(img.ID, img.ProductID, img.ImageURL, img.DisplayOrder, img.CreatedAt)
	}

	query, args, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}

	_, err = r.dbGetter(ctx).ExecContext(ctx, query, args...)
	return err
}

func (r *productRepository) GetProductByID(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	query, args, err := sq.
		Select(productSelectFields...).
		From("products").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var product entity.Product
	err = r.dbGetter(ctx).GetContext(ctx, &product, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errx.ErrProductNotFound
		}
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) GetAllProducts(ctx context.Context, query dto.GetAllProductsQuery) ([]entity.ProductCardItem, int64, error) {
	builder := sq.
		Select(
			"p.id",
			"p.store_id",
			"s.name AS store_name",
			"p.name",
			"p.description",
			"p.main_image_url",
			"COALESCE(s.operational_start::text, '') AS operational_start",
			"COALESCE(s.operational_end::text, '') AS operational_end",
			"s.city AS store_city",
			"s.district AS store_district",
			"s.latitude AS store_latitude",
			"s.longitude AS store_longitude",
			"COALESCE(SUM(pv.current_stock), 0)::int AS total_current_stock",
			"COALESCE(MIN(pv.normal_price), 0)::bigint AS min_normal_price",
			"COALESCE(MIN(pv.rescue_price), 0)::bigint AS min_rescue_price",
		).
		From("products p").
		Join("stores s ON s.id = p.store_id").
		LeftJoin("product_variants pv ON pv.product_id = p.id AND pv.is_active = true").
		Where(sq.Eq{
			"p.is_active": true,
			"s.status":    "ACTIVE",
		})

	countBuilder := sq.
		Select("COUNT(DISTINCT p.id)").
		From("products p").
		Join("stores s ON s.id = p.store_id").
		Where(sq.Eq{
			"p.is_active": true,
			"s.status":    "ACTIVE",
		})

	if query.StoreID != "" {
		if storeUUID, err := uuid.Parse(query.StoreID); err == nil {
			builder = builder.Where(sq.Eq{"p.store_id": storeUUID})
			countBuilder = countBuilder.Where(sq.Eq{"p.store_id": storeUUID})
		}
	}

	if query.City != "" {
		cityPattern := fmt.Sprintf("%%%s%%", query.City)
		builder = builder.Where(sq.ILike{"s.city": cityPattern})
		countBuilder = countBuilder.Where(sq.ILike{"s.city": cityPattern})
	}

	searchKeyword := query.Search
	if searchKeyword == "" {
		searchKeyword = query.Nama
	}
	if searchKeyword == "" {
		searchKeyword = query.Name
	}
	if searchKeyword != "" {
		searchPattern := fmt.Sprintf("%%%s%%", searchKeyword)
		orClause := sq.Or{
			sq.ILike{"p.name": searchPattern},
			sq.ILike{"s.name": searchPattern},
		}
		builder = builder.Where(orClause)
		countBuilder = countBuilder.Where(orClause)
	}

	// Count total
	countSql, countArgs, err := countBuilder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, 0, err
	}
	var total int64
	if err := r.dbGetter(ctx).GetContext(ctx, &total, countSql, countArgs...); err != nil {
		return nil, 0, err
	}

	// Grouping
	builder = builder.GroupBy("p.id", "s.id")

	// Sorting
	builder = builder.OrderBy("p.created_at DESC")

	// Pagination
	if query.Limit > 0 {
		builder = builder.Limit(query.Limit).Offset(query.Offset)
	}

	listSql, listArgs, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, 0, err
	}

	var items []entity.ProductCardItem
	if err := r.dbGetter(ctx).SelectContext(ctx, &items, listSql, listArgs...); err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *productRepository) GetProductsByStoreID(ctx context.Context, storeID uuid.UUID, query dto.GetProductsQuery) ([]entity.Product, int64, error) {
	builder := sq.
		Select("p.id", "p.store_id", "p.name", "p.description", "p.main_image_url", "p.is_active", "p.created_at", "p.updated_at").
		From("products p").
		Where(sq.Eq{"p.store_id": storeID})

	countBuilder := sq.
		Select("COUNT(DISTINCT p.id)").
		From("products p").
		Where(sq.Eq{"p.store_id": storeID})

	if query.Search != "" {
		searchPattern := fmt.Sprintf("%%%s%%", query.Search)
		builder = builder.Where(sq.ILike{"p.name": searchPattern})
		countBuilder = countBuilder.Where(sq.ILike{"p.name": searchPattern})
	}

	switch query.Tab {
	case "active":
		builder = builder.Where(sq.Eq{"p.is_active": true})
		countBuilder = countBuilder.Where(sq.Eq{"p.is_active": true})
	case "low_stock":
		// Has at least one variant with remaining stock between 1 and 2
		subQuery := "EXISTS (SELECT 1 FROM product_variants pv WHERE pv.product_id = p.id AND pv.current_stock > 0 AND pv.current_stock <= 2)"
		builder = builder.Where(subQuery)
		countBuilder = countBuilder.Where(subQuery)
	case "out_of_stock":
		// All variants have 0 stock
		subQuery := "NOT EXISTS (SELECT 1 FROM product_variants pv WHERE pv.product_id = p.id AND pv.current_stock > 0)"
		builder = builder.Where(subQuery)
		countBuilder = countBuilder.Where(subQuery)
	}

	// Count total
	countSql, countArgs, err := countBuilder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, 0, err
	}
	var total int64
	if err := r.dbGetter(ctx).GetContext(ctx, &total, countSql, countArgs...); err != nil {
		return nil, 0, err
	}

	// Sorting
	if query.Sort == "oldest" {
		builder = builder.OrderBy("p.created_at ASC")
	} else {
		builder = builder.OrderBy("p.created_at DESC")
	}

	// Pagination
	limit := uint64(10)
	if query.Limit > 0 {
		limit = uint64(query.Limit)
	}
	page := uint64(1)
	if query.Page > 0 {
		page = uint64(query.Page)
	}
	offset := (page - 1) * limit

	builder = builder.Limit(limit).Offset(offset)

	listSql, listArgs, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, 0, err
	}

	var products []entity.Product
	if err := r.dbGetter(ctx).SelectContext(ctx, &products, listSql, listArgs...); err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) GetVariantsByProductIDs(ctx context.Context, productIDs []uuid.UUID) (map[uuid.UUID][]entity.ProductVariant, error) {
	result := make(map[uuid.UUID][]entity.ProductVariant)
	if len(productIDs) == 0 {
		return result, nil
	}

	query, args, err := sq.
		Select(variantSelectFields...).
		From("product_variants").
		Where(sq.Eq{"product_id": productIDs}).
		OrderBy("sequence ASC").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var variants []entity.ProductVariant
	if err := r.dbGetter(ctx).SelectContext(ctx, &variants, query, args...); err != nil {
		return nil, err
	}

	for _, v := range variants {
		result[v.ProductID] = append(result[v.ProductID], v)
	}

	return result, nil
}

func (r *productRepository) GetVariantByID(ctx context.Context, id uuid.UUID) (*entity.ProductVariant, error) {
	query, args, err := sq.
		Select(variantSelectFields...).
		From("product_variants").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var variant entity.ProductVariant
	err = r.dbGetter(ctx).GetContext(ctx, &variant, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errx.ErrVariantNotFound
		}
		return nil, err
	}

	return &variant, nil
}

func (r *productRepository) DeductVariantStock(ctx context.Context, id uuid.UUID, quantity int) error {
	rawQuery := `
		UPDATE product_variants
		SET current_stock = current_stock - $1,
		    updated_at = NOW()
		WHERE id = $2 AND current_stock >= $1
	`
	res, err := r.dbGetter(ctx).ExecContext(ctx, rawQuery, quantity, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errx.ErrInsufficientStock
	}
	return nil
}

func (r *productRepository) RestoreVariantStock(ctx context.Context, id uuid.UUID, quantity int) error {
	rawQuery := `
		UPDATE product_variants
		SET current_stock = current_stock + $1,
		    updated_at = NOW()
		WHERE id = $2
	`
	_, err := r.dbGetter(ctx).ExecContext(ctx, rawQuery, quantity, id)
	return err
}

func (r *productRepository) GetImagesByProductIDs(ctx context.Context, productIDs []uuid.UUID) (map[uuid.UUID][]entity.ProductImage, error) {
	result := make(map[uuid.UUID][]entity.ProductImage)
	if len(productIDs) == 0 {
		return result, nil
	}

	query, args, err := sq.
		Select(imageSelectFields...).
		From("product_images").
		Where(sq.Eq{"product_id": productIDs}).
		OrderBy("display_order ASC").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var images []entity.ProductImage
	if err := r.dbGetter(ctx).SelectContext(ctx, &images, query, args...); err != nil {
		return nil, err
	}

	for _, img := range images {
		result[img.ProductID] = append(result[img.ProductID], img)
	}

	return result, nil
}

func (r *productRepository) UpdateProduct(ctx context.Context, product entity.Product) error {
	query, args, err := sq.
		Update("products").
		Set("name", product.Name).
		Set("description", product.Description).
		Set("main_image_url", product.MainImageURL).
		Set("is_active", product.IsActive).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": product.ID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.dbGetter(ctx).ExecContext(ctx, query, args...)
	return err
}

func (r *productRepository) UpdateProductStatus(ctx context.Context, productID uuid.UUID, isActive bool) error {
	query, args, err := sq.
		Update("products").
		Set("is_active", isActive).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": productID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.dbGetter(ctx).ExecContext(ctx, query, args...)
	return err
}

func (r *productRepository) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	query, args, err := sq.
		Delete("products").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.dbGetter(ctx).ExecContext(ctx, query, args...)
	return err
}
