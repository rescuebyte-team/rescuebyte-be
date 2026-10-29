package contracts

import (
	"context"

	"github.com/google/uuid"

	"github.com/rescuebyte/rescuebyte-be/domain/dto"
	"github.com/rescuebyte/rescuebyte-be/domain/entity"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, product entity.Product) error
	CreateProductVariants(ctx context.Context, variants []entity.ProductVariant) error
	CreateProductImages(ctx context.Context, images []entity.ProductImage) error
	GetProductByID(ctx context.Context, id uuid.UUID) (*entity.Product, error)
	GetAllProducts(ctx context.Context, query dto.GetAllProductsQuery) ([]entity.ProductCardItem, int64, error)
	GetProductsByStoreID(ctx context.Context, storeID uuid.UUID, query dto.GetProductsQuery) ([]entity.Product, int64, error)
	GetVariantsByProductIDs(ctx context.Context, productIDs []uuid.UUID) (map[uuid.UUID][]entity.ProductVariant, error)
	GetImagesByProductIDs(ctx context.Context, productIDs []uuid.UUID) (map[uuid.UUID][]entity.ProductImage, error)
	UpdateProduct(ctx context.Context, product entity.Product) error
	UpdateProductStatus(ctx context.Context, productID uuid.UUID, isActive bool) error
	DeleteProduct(ctx context.Context, id uuid.UUID) error
}

type ProductService interface {
	CreateProduct(ctx context.Context, userID uuid.UUID, req dto.CreateProductRequest) (*dto.ProductResponse, error)
	GetAllProducts(ctx context.Context, query dto.GetAllProductsQuery) (dto.GetProductsResponse, error)
	GetMyProducts(ctx context.Context, userID uuid.UUID, query dto.GetProductsQuery) ([]dto.ProductResponse, int64, error)
	GetProductDetail(ctx context.Context, id uuid.UUID) (*dto.ProductResponse, error)
	UpdateProduct(ctx context.Context, userID uuid.UUID, productID uuid.UUID, req dto.UpdateProductRequest) (*dto.ProductResponse, error)
	ToggleProductStatus(ctx context.Context, userID uuid.UUID, productID uuid.UUID, req dto.ToggleProductStatusRequest) (*dto.ProductResponse, error)
	DeleteProduct(ctx context.Context, userID uuid.UUID, productID uuid.UUID) error
}
