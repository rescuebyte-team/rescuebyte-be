package cart_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/rescuebyte/rescuebyte-be/domain/dto"
	"github.com/rescuebyte/rescuebyte-be/domain/entity"
)

func TestFormatToCartResponse(t *testing.T) {
	store1ID := uuid.New()
	store2ID := uuid.New()

	v1ID := uuid.New()
	v2ID := uuid.New()
	v3ID := uuid.New()

	details := []entity.CartItemDetail{
		{
			ID:               uuid.New(),
			CartID:           uuid.New(),
			ProductVariantID: v1ID,
			Quantity:         2,
			ProductID:        uuid.New(),
			ProductName:      "Mystery Bakery Box",
			VariantName:      "Reguler",
			NormalPrice:      40000,
			RescuePrice:      20000,
			CurrentStock:     5,
			VariantActive:    true,
			ProductActive:    true,
			StoreID:          store1ID,
			StoreName:        "TerasRasa Bakery",
			StoreAddress:     "Jl. Soekarno Hatta",
			StoreCity:        "Malang",
			StoreDistrict:    "Lowokwaru",
			StoreLatitude:    -7.9400,
			StoreLongitude:   112.6100,
			OperationalStart: "15:00:00",
			OperationalEnd:   "18:00:00",
		},
		{
			ID:               uuid.New(),
			CartID:           uuid.New(),
			ProductVariantID: v2ID,
			Quantity:         1,
			ProductID:        uuid.New(),
			ProductName:      "Mystery Pastry Box",
			VariantName:      "Reguler",
			NormalPrice:      36000,
			RescuePrice:      18000,
			CurrentStock:     3,
			VariantActive:    true,
			ProductActive:    true,
			StoreID:          store1ID,
			StoreName:        "TerasRasa Bakery",
			StoreAddress:     "Jl. Soekarno Hatta",
			StoreCity:        "Malang",
			StoreDistrict:    "Lowokwaru",
			StoreLatitude:    -7.9400,
			StoreLongitude:   112.6100,
			OperationalStart: "15:00:00",
			OperationalEnd:   "18:00:00",
		},
		{
			ID:               uuid.New(),
			CartID:           uuid.New(),
			ProductVariantID: v3ID,
			Quantity:         1,
			ProductID:        uuid.New(),
			ProductName:      "Mystery Meal Box",
			VariantName:      "Default",
			NormalPrice:      50000,
			RescuePrice:      25000,
			CurrentStock:     2,
			VariantActive:    true,
			ProductActive:    true,
			StoreID:          store2ID,
			StoreName:        "GreenBowl",
			StoreAddress:     "Jl. Ijen",
			StoreCity:        "Malang",
			StoreDistrict:    "Klojen",
			StoreLatitude:    -7.9700,
			StoreLongitude:   112.6200,
			OperationalStart: "12:00:00",
			OperationalEnd:   "15:00:00",
		},
	}

	userLat := -7.9500
	userLon := 112.6150

	res := dto.FormatToCartResponse(details, &userLat, &userLon)

	assert.Equal(t, 3, res.TotalItems)
	assert.Equal(t, 4, res.TotalBoxes) // 2 + 1 + 1 = 4 boxes
	// Subtotal: 2*20000 + 1*18000 + 1*25000 = 40000 + 18000 + 25000 = 83000
	assert.Equal(t, int64(83000), res.Subtotal)
	// Original: 2*40000 + 1*36000 + 1*50000 = 80000 + 36000 + 50000 = 166000
	assert.Equal(t, int64(166000), res.TotalOriginalPrice)
	// Discount: 166000 - 83000 = 83000
	assert.Equal(t, int64(83000), res.TotalDiscount)
	// Service fee: 2 stores * 2000 = 4000
	assert.Equal(t, int64(4000), res.ServiceFee)
	// Total Payment: 83000 + 4000 = 87000
	assert.Equal(t, int64(87000), res.TotalPayment)

	assert.Equal(t, 2, len(res.Stores))
	assert.Equal(t, "TerasRasa Bakery", res.Stores[0].StoreName)
	assert.Equal(t, 2, len(res.Stores[0].Items))
	assert.Equal(t, int64(58000), res.Stores[0].StoreSubtotal) // 40000 + 18000
	assert.Equal(t, "GreenBowl", res.Stores[1].StoreName)
	assert.Equal(t, 1, len(res.Stores[1].Items))
}
