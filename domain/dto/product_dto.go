package dto

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/rescuebyte/rescuebyte-be/domain/entity"
	"github.com/rescuebyte/rescuebyte-be/domain/helpers/pagination"
)

type CreateProductVariantRequest struct {
	Name         string `json:"name" validate:"required,max=100"`
	Sequence     int    `json:"sequence" validate:"required,min=1"`
	NormalPrice  int64  `json:"normal_price" validate:"required,gte=0"`
	RescuePrice  int64  `json:"rescue_price" validate:"required,gte=0"`
	InitialStock int    `json:"initial_stock" validate:"required,gte=0"`
	CurrentStock int    `json:"current_stock" validate:"required,gte=0"`
}

type CreateProductRequest struct {
	Name         string                        `json:"name" validate:"required,min=2,max=255"`
	Description  *string                       `json:"description" validate:"omitempty"`
	MainImageURL *string                       `json:"main_image_url" validate:"omitempty,url"`
	Variants     []CreateProductVariantRequest `json:"variants" validate:"required,min=1,dive"`
	ImageURLs    []string                      `json:"image_urls" validate:"omitempty,dive,url"`
}

type UpdateProductRequest struct {
	Name         string  `json:"name" validate:"required,min=2,max=255"`
	Description  *string `json:"description" validate:"omitempty"`
	MainImageURL *string `json:"main_image_url" validate:"omitempty,url"`
	IsActive     *bool   `json:"is_active" validate:"omitempty"`
}

type UpdateProductVariantRequest struct {
	ID           *uuid.UUID `json:"id" validate:"omitempty"`
	Name         string     `json:"name" validate:"required,max=100"`
	Sequence     int        `json:"sequence" validate:"required,min=1"`
	NormalPrice  int64      `json:"normal_price" validate:"required,gte=0"`
	RescuePrice  int64      `json:"rescue_price" validate:"required,gte=0"`
	InitialStock int        `json:"initial_stock" validate:"required,gte=0"`
	CurrentStock int        `json:"current_stock" validate:"required,gte=0"`
	IsActive     bool       `json:"is_active"`
}

type ToggleProductStatusRequest struct {
	IsActive bool `json:"is_active"`
}

type GetProductsQuery struct {
	Page   int    `query:"page"`
	Limit  int    `query:"limit"`
	Search string `query:"search"`
	Tab    string `query:"tab"`  // "all", "active", "low_stock", "out_of_stock"
	Sort   string `query:"sort"` // "newest", "oldest"
}

type GetAllProductsQuery struct {
	Limit     uint64 `query:"limit" validate:"omitempty,number"`
	Page      uint64 `query:"page" validate:"omitempty,number,gte=1"`
	Offset    uint64
	StoreID   string   `query:"store_id" validate:"omitempty,uuid"`
	Nama      string   `query:"nama" validate:"omitempty"`
	Name      string   `query:"name" validate:"omitempty"`
	Search    string   `query:"search" validate:"omitempty"`
	City      string   `query:"city" validate:"omitempty"`
	Latitude  *float64 `query:"latitude" validate:"omitempty,latitude"`
	Longitude *float64 `query:"longitude" validate:"omitempty,longitude"`
}

type ProductCardResponse struct {
	ID                 uuid.UUID `json:"id"`
	StoreID            uuid.UUID `json:"store_id"`
	StoreName          string    `json:"store_name"`
	Name               string    `json:"name"`
	Description        *string   `json:"description"`
	MainImageURL       *string   `json:"main_image_url"`
	Stock              int       `json:"stock"`
	NormalPrice        int64     `json:"normal_price"`
	RescuePrice        int64     `json:"rescue_price"`
	DiscountPercentage int       `json:"discount_percentage"`
	OperationalStart   string    `json:"operational_start"`
	OperationalEnd     string    `json:"operational_end"`
	PickupTime         string    `json:"pickup_time"`
	DistanceKm         *float64  `json:"distance_km,omitempty"`
}

type ProductVariantResponse struct {
	ID           uuid.UUID `json:"id"`
	ProductID    uuid.UUID `json:"product_id"`
	Name         string    `json:"name"`
	Sequence     int       `json:"sequence"`
	NormalPrice  int64     `json:"normal_price"`
	RescuePrice  int64     `json:"rescue_price"`
	InitialStock int       `json:"initial_stock"`
	CurrentStock int       `json:"current_stock"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ProductImageResponse struct {
	ID           uuid.UUID `json:"id"`
	ImageURL     string    `json:"image_url"`
	DisplayOrder int       `json:"display_order"`
}

type ProductResponse struct {
	ID           uuid.UUID                `json:"id"`
	StoreID      uuid.UUID                `json:"store_id"`
	Name         string                   `json:"name"`
	Description  *string                  `json:"description"`
	MainImageURL *string                  `json:"main_image_url"`
	IsActive     bool                     `json:"is_active"`
	TotalInitial int                      `json:"total_initial_stock"`
	TotalCurrent int                      `json:"total_current_stock"`
	MinRescue    int64                    `json:"min_rescue_price"`
	MinNormal    int64                    `json:"min_normal_price"`
	Variants     []ProductVariantResponse `json:"variants,omitempty"`
	Images       []ProductImageResponse   `json:"images,omitempty"`
	CreatedAt    time.Time                `json:"created_at"`
	UpdatedAt    time.Time                `json:"updated_at"`
}

type GetProductsResponse struct {
	Products []ProductCardResponse `json:"products"`
	Meta     pagination.Response   `json:"meta"`
}

func FormatTimeStr(t string) string {
	t = strings.TrimSpace(t)
	if len(t) >= 5 {
		return t[:5]
	}
	return t
}

func CalculateDistanceKm(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371.0 // kilometers
	dLat := (lat2 - lat1) * (math.Pi / 180.0)
	dLon := (lon2 - lon1) * (math.Pi / 180.0)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*(math.Pi/180.0))*math.Cos(lat2*(math.Pi/180.0))*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return math.Round(earthRadius*c*10) / 10
}

func FormatToProductCardResponse(item entity.ProductCardItem, userLat, userLon *float64) ProductCardResponse {
	var discount int
	if item.MinNormalPrice > 0 && item.MinRescuePrice < item.MinNormalPrice {
		discount = int(math.Round(float64(item.MinNormalPrice-item.MinRescuePrice) / float64(item.MinNormalPrice) * 100))
	}

	opStart := FormatTimeStr(item.OperationalStart)
	opEnd := FormatTimeStr(item.OperationalEnd)
	pickupTime := ""
	if opStart != "" && opEnd != "" {
		pickupTime = fmt.Sprintf("Ambil: %s – %s WIB", opStart, opEnd)
	}

	var distanceKm *float64
	if userLat != nil && userLon != nil {
		dist := CalculateDistanceKm(*userLat, *userLon, item.StoreLatitude, item.StoreLongitude)
		distanceKm = &dist
	}

	return ProductCardResponse{
		ID:                 item.ID,
		StoreID:            item.StoreID,
		StoreName:          item.StoreName,
		Name:               item.Name,
		Description:        item.Description,
		MainImageURL:       item.MainImageURL,
		Stock:              item.TotalCurrentStock,
		NormalPrice:        item.MinNormalPrice,
		RescuePrice:        item.MinRescuePrice,
		DiscountPercentage: discount,
		OperationalStart:   opStart,
		OperationalEnd:     opEnd,
		PickupTime:         pickupTime,
		DistanceKm:         distanceKm,
	}
}

func FormatToProductVariantResponse(v entity.ProductVariant) ProductVariantResponse {
	return ProductVariantResponse{
		ID:           v.ID,
		ProductID:    v.ProductID,
		Name:         v.Name,
		Sequence:     v.Sequence,
		NormalPrice:  v.NormalPrice,
		RescuePrice:  v.RescuePrice,
		InitialStock: v.InitialStock,
		CurrentStock: v.CurrentStock,
		IsActive:     v.IsActive,
		CreatedAt:    v.CreatedAt,
		UpdatedAt:    v.UpdatedAt,
	}
}

func FormatToProductResponse(p entity.Product) ProductResponse {
	var totalInit, totalCurr int
	var minRescue, minNormal int64 = -1, -1

	variants := make([]ProductVariantResponse, 0, len(p.Variants))
	for _, v := range p.Variants {
		variants = append(variants, FormatToProductVariantResponse(v))
		totalInit += v.InitialStock
		totalCurr += v.CurrentStock

		if minRescue == -1 || v.RescuePrice < minRescue {
			minRescue = v.RescuePrice
		}
		if minNormal == -1 || v.NormalPrice < minNormal {
			minNormal = v.NormalPrice
		}
	}

	images := make([]ProductImageResponse, 0, len(p.Images))
	for _, img := range p.Images {
		images = append(images, ProductImageResponse{
			ID:           img.ID,
			ImageURL:     img.ImageURL,
			DisplayOrder: img.DisplayOrder,
		})
	}

	if minRescue == -1 {
		minRescue = 0
	}
	if minNormal == -1 {
		minNormal = 0
	}

	return ProductResponse{
		ID:           p.ID,
		StoreID:      p.StoreID,
		Name:         p.Name,
		Description:  p.Description,
		MainImageURL: p.MainImageURL,
		IsActive:     p.IsActive,
		TotalInitial: totalInit,
		TotalCurrent: totalCurr,
		MinRescue:    minRescue,
		MinNormal:    minNormal,
		Variants:     variants,
		Images:       images,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
	}
}
