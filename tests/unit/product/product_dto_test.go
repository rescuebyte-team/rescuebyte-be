package product_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/rescuebyte/rescuebyte-be/domain/dto"
	"github.com/rescuebyte/rescuebyte-be/domain/entity"
)

func TestFormatTimeStr(t *testing.T) {
	assert.Equal(t, "19:45", dto.FormatTimeStr("19:45:00"))
	assert.Equal(t, "19:45", dto.FormatTimeStr("19:45"))
	assert.Equal(t, "08:00", dto.FormatTimeStr("08:00:15"))
	assert.Equal(t, "", dto.FormatTimeStr(""))
}

func TestCalculateDistanceKm(t *testing.T) {
	// Jakarta to Bandung is roughly ~120-150 km
	jktLat, jktLon := -6.2088, 106.8456
	bdgLat, bdgLon := -6.9175, 107.6191

	dist := dto.CalculateDistanceKm(jktLat, jktLon, bdgLat, bdgLon)
	assert.True(t, dist > 110 && dist < 160, "Distance should be ~120-150 km")

	// Same point should be 0.0
	sameDist := dto.CalculateDistanceKm(jktLat, jktLon, jktLat, jktLon)
	assert.Equal(t, 0.0, sameDist)
}

func TestFormatToProductCardResponse(t *testing.T) {
	prodID := uuid.New()
	storeID := uuid.New()
	desc := "Artisan French Viennoiserie Surprise Box"
	img := "https://example.com/pastry.jpg"

	item := entity.ProductCardItem{
		ID:                prodID,
		StoreID:           storeID,
		StoreName:         "Monsieur Spoon Senopati",
		Name:              "PASTRY SURPRISE",
		Description:       &desc,
		MainImageURL:      &img,
		OperationalStart:  "19:45:00",
		OperationalEnd:    "21:00:00",
		StoreCity:         "Jakarta Selatan",
		StoreDistrict:     "Kebayoran Baru",
		StoreLatitude:     -6.2300,
		StoreLongitude:    106.8100,
		TotalCurrentStock: 2,
		MinNormalPrice:    110000,
		MinRescuePrice:    35000,
	}

	userLat := -6.2350
	userLon := 106.8120

	res := dto.FormatToProductCardResponse(item, &userLat, &userLon)

	assert.Equal(t, prodID, res.ID)
	assert.Equal(t, storeID, res.StoreID)
	assert.Equal(t, "Monsieur Spoon Senopati", res.StoreName)
	assert.Equal(t, "PASTRY SURPRISE", res.Name)
	assert.Equal(t, 2, res.Stock)
	assert.Equal(t, int64(110000), res.NormalPrice)
	assert.Equal(t, int64(35000), res.RescuePrice)
	assert.Equal(t, 68, res.DiscountPercentage) // (110000-35000)/110000 = 68.18% -> 68%
	assert.Equal(t, "19:45", res.OperationalStart)
	assert.Equal(t, "21:00", res.OperationalEnd)
	assert.Equal(t, "Ambil: 19:45 – 21:00 WIB", res.PickupTime)
	assert.NotNil(t, res.DistanceKm)
	assert.True(t, *res.DistanceKm > 0.0)

	// Test without user coordinates
	resWithoutCoords := dto.FormatToProductCardResponse(item, nil, nil)
	assert.Nil(t, resWithoutCoords.DistanceKm)
}

func TestUpdateProductRequest_Validation(t *testing.T) {
	// Valid with stock
	stock := 10
	req := dto.UpdateProductRequest{
		Name:  "Updated Product",
		Stock: &stock,
	}
	assert.Equal(t, "Updated Product", req.Name)
	assert.NotNil(t, req.Stock)
	assert.Equal(t, 10, *req.Stock)

	// Valid with zero stock
	zeroStock := 0
	reqZero := dto.UpdateProductRequest{
		Name:  "Zero Stock Product",
		Stock: &zeroStock,
	}
	assert.NotNil(t, reqZero.Stock)
	assert.Equal(t, 0, *reqZero.Stock)

	// Valid without stock (nil)
	reqNil := dto.UpdateProductRequest{
		Name: "Product No Stock Update",
	}
	assert.Nil(t, reqNil.Stock)
}
