package dto

import (
	"fmt"
	"math"

	"github.com/google/uuid"

	"github.com/rescuebyte/rescuebyte-be/domain/entity"
)

type AddToCartRequest struct {
	ProductVariantID uuid.UUID `json:"product_variant_id" validate:"required"`
	Quantity         int       `json:"quantity" validate:"required,min=1"`
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" validate:"required,min=1"`
}

type CartItemResponse struct {
	ID                 uuid.UUID `json:"id"`
	ProductVariantID   uuid.UUID `json:"product_variant_id"`
	ProductID          uuid.UUID `json:"product_id"`
	ProductName        string    `json:"product_name"`
	VariantName        string    `json:"variant_name"`
	ProductImageURL    *string   `json:"product_image_url"`
	NormalPrice        int64     `json:"normal_price"`
	RescuePrice        int64     `json:"rescue_price"`
	DiscountPercentage int       `json:"discount_percentage"`
	Quantity           int       `json:"quantity"`
	Subtotal           int64     `json:"subtotal"`
	CurrentStock       int       `json:"current_stock"`
	IsAvailable        bool      `json:"is_available"`
}

type StoreCartGroup struct {
	StoreID          uuid.UUID          `json:"store_id"`
	StoreName        string             `json:"store_name"`
	StoreAddress     string             `json:"store_address"`
	StoreCity        string             `json:"store_city"`
	StoreDistrict    string             `json:"store_district"`
	OperationalStart string             `json:"operational_start"`
	OperationalEnd   string             `json:"operational_end"`
	PickupTime       string             `json:"pickup_time"`
	DistanceKm       *float64           `json:"distance_km,omitempty"`
	Items            []CartItemResponse `json:"items"`
	StoreSubtotal    int64              `json:"store_subtotal"`
}

type CartResponse struct {
	TotalItems         int              `json:"total_items"`
	TotalBoxes         int              `json:"total_boxes"`
	TotalOriginalPrice int64            `json:"total_original_price"`
	TotalDiscount      int64            `json:"total_discount"`
	Subtotal           int64            `json:"subtotal"`
	ServiceFee         int64            `json:"service_fee"`
	TotalPayment       int64            `json:"total_payment"`
	Stores             []StoreCartGroup `json:"stores"`
}

type CartCountResponse struct {
	Count int `json:"count"`
}

func FormatToCartResponse(details []entity.CartItemDetail, userLat, userLon *float64) CartResponse {
	storeMap := make(map[uuid.UUID]*StoreCartGroup)
	var storeOrder []uuid.UUID

	var totalOriginal, subtotal int64
	var totalBoxes int

	for _, item := range details {
		discount := 0
		if item.NormalPrice > 0 && item.RescuePrice < item.NormalPrice {
			discount = int(math.Round(float64(item.NormalPrice-item.RescuePrice) / float64(item.NormalPrice) * 100))
		}

		itemSubtotal := item.RescuePrice * int64(item.Quantity)
		isAvailable := item.CurrentStock >= item.Quantity && item.VariantActive && item.ProductActive

		itemResp := CartItemResponse{
			ID:                 item.ID,
			ProductVariantID:   item.ProductVariantID,
			ProductID:          item.ProductID,
			ProductName:        item.ProductName,
			VariantName:        item.VariantName,
			ProductImageURL:    item.ProductImageURL,
			NormalPrice:        item.NormalPrice,
			RescuePrice:        item.RescuePrice,
			DiscountPercentage: discount,
			Quantity:           item.Quantity,
			Subtotal:           itemSubtotal,
			CurrentStock:       item.CurrentStock,
			IsAvailable:        isAvailable,
		}

		subtotal += itemSubtotal
		totalOriginal += item.NormalPrice * int64(item.Quantity)
		totalBoxes += item.Quantity

		if group, exists := storeMap[item.StoreID]; exists {
			group.Items = append(group.Items, itemResp)
			group.StoreSubtotal += itemSubtotal
		} else {
			opStart := FormatTimeStr(item.OperationalStart)
			opEnd := FormatTimeStr(item.OperationalEnd)
			pickupTime := ""
			if opStart != "" && opEnd != "" {
				pickupTime = fmt.Sprintf("Waktu Pickup %s – %s WIB", opStart, opEnd)
			}

			var distanceKm *float64
			if userLat != nil && userLon != nil {
				d := CalculateDistanceKm(*userLat, *userLon, item.StoreLatitude, item.StoreLongitude)
				distanceKm = &d
			}

			newGroup := &StoreCartGroup{
				StoreID:          item.StoreID,
				StoreName:        item.StoreName,
				StoreAddress:     item.StoreAddress,
				StoreCity:        item.StoreCity,
				StoreDistrict:    item.StoreDistrict,
				OperationalStart: opStart,
				OperationalEnd:   opEnd,
				PickupTime:       pickupTime,
				DistanceKm:       distanceKm,
				Items:            []CartItemResponse{itemResp},
				StoreSubtotal:    itemSubtotal,
			}
			storeMap[item.StoreID] = newGroup
			storeOrder = append(storeOrder, item.StoreID)
		}
	}

	stores := make([]StoreCartGroup, 0, len(storeOrder))
	for _, storeID := range storeOrder {
		stores = append(stores, *storeMap[storeID])
	}

	var serviceFee int64 = 0
	if len(stores) > 0 {
		serviceFee = 2000 * int64(len(stores)) // Rp 2.000 per store
	}

	totalDiscount := totalOriginal - subtotal
	if totalDiscount < 0 {
		totalDiscount = 0
	}

	totalPayment := subtotal + serviceFee

	return CartResponse{
		TotalItems:         len(details),
		TotalBoxes:         totalBoxes,
		TotalOriginalPrice: totalOriginal,
		TotalDiscount:      totalDiscount,
		Subtotal:           subtotal,
		ServiceFee:         serviceFee,
		TotalPayment:       totalPayment,
		Stores:             stores,
	}
}
