package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/rescuebyte/rescuebyte-be/domain/dto"
	"github.com/rescuebyte/rescuebyte-be/domain/entity"
	"github.com/rescuebyte/rescuebyte-be/domain/errx"
	"github.com/rescuebyte/rescuebyte-be/domain/helpers/pagination"
)

func (s *productService) CreateProduct(ctx context.Context, userID uuid.UUID, req dto.CreateProductRequest) (*dto.ProductResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}

	store, err := s.storeRepo.GetStoreByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if store == nil {
		return nil, errx.ErrStoreNotFound
	}

	if len(req.Variants) == 0 {
		return nil, errx.ErrProductMustHaveVariant
	}

	now := time.Now()
	productID, err := s.uuid.NewV7()
	if err != nil {
		return nil, err
	}

	product := entity.Product{
		ID:           productID,
		StoreID:      store.ID,
		Name:         req.Name,
		Description:  req.Description,
		MainImageURL: req.MainImageURL,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	variants := make([]entity.ProductVariant, 0, len(req.Variants))
	for _, v := range req.Variants {
		vID, err := s.uuid.NewV7()
		if err != nil {
			return nil, err
		}
		variants = append(variants, entity.ProductVariant{
			ID:           vID,
			ProductID:    productID,
			Name:         v.Name,
			Sequence:     v.Sequence,
			NormalPrice:  v.NormalPrice,
			RescuePrice:  v.RescuePrice,
			InitialStock: v.InitialStock,
			CurrentStock: v.CurrentStock,
			IsActive:     true,
			CreatedAt:    now,
			UpdatedAt:    now,
		})
	}

	images := make([]entity.ProductImage, 0, len(req.ImageURLs))
	for idx, imgURL := range req.ImageURLs {
		imgID, err := s.uuid.NewV7()
		if err != nil {
			return nil, err
		}
		images = append(images, entity.ProductImage{
			ID:           imgID,
			ProductID:    productID,
			ImageURL:     imgURL,
			DisplayOrder: idx + 1,
			CreatedAt:    now,
		})
	}

	err = s.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := s.productRepo.CreateProduct(txCtx, product); err != nil {
			return err
		}
		if err := s.productRepo.CreateProductVariants(txCtx, variants); err != nil {
			return err
		}
		if len(images) > 0 {
			if err := s.productRepo.CreateProductImages(txCtx, images); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	product.Variants = variants
	product.Images = images

	res := dto.FormatToProductResponse(product)
	return &res, nil
}

func (s *productService) GetAllProducts(ctx context.Context, query dto.GetAllProductsQuery) (dto.GetProductsResponse, error) {
	if err := s.validator.Validate(query); err != nil {
		return dto.GetProductsResponse{}, err
	}

	pagination.Normalize(&query.Limit, &query.Page, &query.Offset, "product")

	items, total, err := s.productRepo.GetAllProducts(ctx, query)
	if err != nil {
		return dto.GetProductsResponse{}, err
	}

	if len(items) == 0 {
		return dto.GetProductsResponse{
			Products: []dto.ProductCardResponse{},
			Meta:     pagination.NewPaginationResponse(total, query.Page, query.Limit),
		}, nil
	}

	responses := make([]dto.ProductCardResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, dto.FormatToProductCardResponse(item, query.Latitude, query.Longitude))
	}

	return dto.GetProductsResponse{
		Products: responses,
		Meta:     pagination.NewPaginationResponse(total, query.Page, query.Limit),
	}, nil
}

func (s *productService) GetMyProducts(ctx context.Context, userID uuid.UUID, query dto.GetProductsQuery) ([]dto.ProductResponse, int64, error) {
	store, err := s.storeRepo.GetStoreByUserID(ctx, userID)
	if err != nil {
		return nil, 0, err
	}
	if store == nil {
		return nil, 0, errx.ErrStoreNotFound
	}

	products, total, err := s.productRepo.GetProductsByStoreID(ctx, store.ID, query)
	if err != nil {
		return nil, 0, err
	}

	if len(products) == 0 {
		return []dto.ProductResponse{}, total, nil
	}

	productIDs := make([]uuid.UUID, 0, len(products))
	for _, p := range products {
		productIDs = append(productIDs, p.ID)
	}

	variantMap, err := s.productRepo.GetVariantsByProductIDs(ctx, productIDs)
	if err != nil {
		return nil, 0, err
	}

	imageMap, err := s.productRepo.GetImagesByProductIDs(ctx, productIDs)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.ProductResponse, 0, len(products))
	for _, p := range products {
		p.Variants = variantMap[p.ID]
		p.Images = imageMap[p.ID]
		responses = append(responses, dto.FormatToProductResponse(p))
	}

	return responses, total, nil
}

func (s *productService) GetProductDetail(ctx context.Context, id uuid.UUID) (*dto.ProductResponse, error) {
	product, err := s.productRepo.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}

	variantMap, err := s.productRepo.GetVariantsByProductIDs(ctx, []uuid.UUID{product.ID})
	if err != nil {
		return nil, err
	}

	imageMap, err := s.productRepo.GetImagesByProductIDs(ctx, []uuid.UUID{product.ID})
	if err != nil {
		return nil, err
	}

	product.Variants = variantMap[product.ID]
	product.Images = imageMap[product.ID]

	res := dto.FormatToProductResponse(*product)
	return &res, nil
}

func (s *productService) UpdateProduct(ctx context.Context, userID uuid.UUID, productID uuid.UUID, req dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}

	store, err := s.storeRepo.GetStoreByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if store == nil {
		return nil, errx.ErrStoreNotFound
	}

	product, err := s.productRepo.GetProductByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	if product.StoreID != store.ID {
		return nil, errx.ErrUnauthorizedStoreAccess
	}

	product.Name = req.Name
	product.Description = req.Description
	product.MainImageURL = req.MainImageURL
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}
	product.UpdatedAt = time.Now()

	err = s.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := s.productRepo.UpdateProduct(txCtx, *product); err != nil {
			return err
		}
		if req.Stock != nil {
			if err := s.productRepo.UpdateProductStock(txCtx, productID, *req.Stock); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.GetProductDetail(ctx, productID)
}

func (s *productService) ToggleProductStatus(ctx context.Context, userID uuid.UUID, productID uuid.UUID, req dto.ToggleProductStatusRequest) (*dto.ProductResponse, error) {
	store, err := s.storeRepo.GetStoreByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if store == nil {
		return nil, errx.ErrStoreNotFound
	}

	product, err := s.productRepo.GetProductByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	if product.StoreID != store.ID {
		return nil, errx.ErrUnauthorizedStoreAccess
	}

	if err := s.productRepo.UpdateProductStatus(ctx, productID, req.IsActive); err != nil {
		return nil, err
	}

	product.IsActive = req.IsActive
	return s.GetProductDetail(ctx, productID)
}

func (s *productService) DeleteProduct(ctx context.Context, userID uuid.UUID, productID uuid.UUID) error {
	store, err := s.storeRepo.GetStoreByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if store == nil {
		return errx.ErrStoreNotFound
	}

	product, err := s.productRepo.GetProductByID(ctx, productID)
	if err != nil {
		return err
	}

	if product.StoreID != store.ID {
		return errx.ErrUnauthorizedStoreAccess
	}

	return s.productRepo.DeleteProduct(ctx, productID)
}
