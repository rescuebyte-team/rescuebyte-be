package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/rescuebyte/rescuebyte-be/domain/entity"
	"github.com/rescuebyte/rescuebyte-be/domain/enums"
	"github.com/rescuebyte/rescuebyte-be/domain/helpers/pagination"
)

type CheckoutItemRequest struct {
	ProductVariantID uuid.UUID `json:"product_variant_id" validate:"required"`
	Quantity         int       `json:"quantity" validate:"required,min=1"`
}

type CheckoutRequest struct {
	StoreID         uuid.UUID             `json:"store_id" validate:"required"`
	Items           []CheckoutItemRequest `json:"items" validate:"required,min=1,dive"`
	CustomerName    string                `json:"customer_name" validate:"required,min=2,max=255"`
	CustomerPhone   string                `json:"customer_phone" validate:"required,min=8,max=20"`
	CustomerEmail   string                `json:"customer_email" validate:"required,email"`
	PickupDate      string                `json:"pickup_date" validate:"required"` // YYYY-MM-DD
	PickupTimeStart string                `json:"pickup_time_start" validate:"required"`
	PickupTimeEnd   string                `json:"pickup_time_end" validate:"required"`
	PaymentMethod   string                `json:"payment_method" validate:"required,oneof=ONLINE CASH_ON_PICKUP"`
	Notes           *string               `json:"notes" validate:"omitempty"`
}

type OrderItemResponse struct {
	ID               uuid.UUID `json:"id"`
	ProductID        uuid.UUID `json:"product_id"`
	ProductVariantID uuid.UUID `json:"product_variant_id"`
	ProductName      string    `json:"product_name"`
	VariantName      string    `json:"variant_name"`
	ProductImageURL  *string   `json:"product_image_url"`
	NormalPrice      int64     `json:"normal_price"`
	RescuePrice      int64     `json:"rescue_price"`
	Quantity         int       `json:"quantity"`
	Subtotal         int64     `json:"subtotal"`
}

type OrderStepTrackingResponse struct {
	StepNumber      int     `json:"step_number"`
	StepName        string  `json:"step_name"`
	StepDescription string  `json:"step_description"`
	IsCompleted     bool    `json:"is_completed"`
	CompletedAt     *string `json:"completed_at,omitempty"`
}

type OrderDetailResponse struct {
	ID                 uuid.UUID                   `json:"id"`
	OrderNumber        string                      `json:"order_number"`
	StoreID            uuid.UUID                   `json:"store_id"`
	StoreName          string                      `json:"store_name"`
	StoreAddress       string                      `json:"store_address"`
	StoreCity          string                      `json:"store_city"`
	StoreDistrict      string                      `json:"store_district"`
	CustomerName       string                      `json:"customer_name"`
	CustomerPhone      string                      `json:"customer_phone"`
	CustomerEmail      string                      `json:"customer_email"`
	PickupDate         string                      `json:"pickup_date"`
	PickupTime         string                      `json:"pickup_time"`
	TotalOriginalPrice int64                       `json:"total_original_price"`
	TotalDiscount      int64                       `json:"total_discount"`
	Subtotal           int64                       `json:"subtotal"`
	ServiceFee         int64                       `json:"service_fee"`
	PackagingFee       int64                       `json:"packaging_fee"`
	TotalAmount        int64                       `json:"total_amount"`
	Status             string                      `json:"status"`
	StatusLabel        string                      `json:"status_label"`
	PaymentMethod      string                      `json:"payment_method"`
	Notes              *string                     `json:"notes"`
	CreatedAt          string                      `json:"created_at"`
	Items              []OrderItemResponse         `json:"items"`
	Steps              []OrderStepTrackingResponse `json:"steps"`
	Payment            *PaymentResponse            `json:"payment,omitempty"`
}

type CheckoutResponse struct {
	Order           OrderDetailResponse `json:"order"`
	SnapToken       *string             `json:"snap_token,omitempty"`
	SnapRedirectURL *string             `json:"snap_redirect_url,omitempty"`
}

type OrderListItemResponse struct {
	ID              uuid.UUID `json:"id"`
	OrderNumber     string    `json:"order_number"`
	StoreID         uuid.UUID `json:"store_id"`
	StoreName       string    `json:"store_name"`
	StoreAddress    string    `json:"store_address"`
	PickupSchedule  string    `json:"pickup_schedule"`
	TotalAmount     int64     `json:"total_amount"`
	Status          string    `json:"status"`
	StatusLabel     string    `json:"status_label"`
	PaymentMethod   string    `json:"payment_method"`
	TotalItems      int       `json:"total_items"`
	SummaryItemName string    `json:"summary_item_name"`
	SummaryImageURL *string   `json:"summary_image_url"`
	CreatedAt       string    `json:"created_at"`
}

type GetOrdersQuery struct {
	Limit  uint64 `query:"limit" validate:"omitempty,number"`
	Page   uint64 `query:"page" validate:"omitempty,number,gte=1"`
	Offset uint64
	Tab    string `query:"tab" validate:"omitempty"` // all, waiting_payment, waiting_confirmation, ready_for_pickup, completed, cancelled
	Search string `query:"search" validate:"omitempty"`
}

type GetOrdersResponse struct {
	Orders []OrderListItemResponse `json:"orders"`
	Meta   pagination.Response     `json:"meta"`
}

type UpdateOrderStatusRequest struct {
	Status             string  `json:"status" validate:"required,oneof=IN_PROCESS READY_FOR_PICKUP COMPLETED CANCELLED"`
	CancellationReason *string `json:"cancellation_reason" validate:"omitempty"`
}

type CancelOrderRequest struct {
	Reason string `json:"reason" validate:"omitempty,max=255"`
}

func FormatStatusLabel(status enums.OrderStatus) string {
	switch status {
	case enums.OrderStatusWaitingPayment:
		return "Menunggu Pembayaran"
	case enums.OrderStatusWaitingConfirmation:
		return "Menunggu Konfirmasi"
	case enums.OrderStatusInProcess:
		return "Diproses"
	case enums.OrderStatusReadyForPickup:
		return "Siap Diambil"
	case enums.OrderStatusCompleted:
		return "Selesai"
	case enums.OrderStatusCancelled:
		return "Dibatalkan"
	default:
		return string(status)
	}
}

func FormatTimeWIB(t time.Time) string {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		loc = time.FixedZone("WIB", 7*3600)
	}
	tWIB := t.In(loc)
	months := []string{
		"", "Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember",
	}
	return fmt.Sprintf("%d %s %d, %02d:%02d WIB",
		tWIB.Day(), months[tWIB.Month()], tWIB.Year(), tWIB.Hour(), tWIB.Minute(),
	)
}

func FormatTimeShortWIB(t time.Time) string {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		loc = time.FixedZone("WIB", 7*3600)
	}
	tWIB := t.In(loc)
	months := []string{
		"", "Okt", "Nov", "Des", "Jan", "Feb", "Mar",
		"Apr", "Mei", "Jun", "Jul", "Agt", "Sep",
	}
	mIdx := int(tWIB.Month())
	monthStr := ""
	if mIdx >= 1 && mIdx <= 12 {
		monthNames := []string{"", "Jan", "Feb", "Mar", "Apr", "Mei", "Jun", "Jul", "Agt", "Sep", "Okt", "Nov", "Des"}
		monthStr = monthNames[mIdx]
	} else {
		monthStr = months[0]
	}
	return fmt.Sprintf("%d %s • %02d:%02d WIB", tWIB.Day(), monthStr, tWIB.Hour(), tWIB.Minute())
}

func BuildOrderSteps(order entity.Order) []OrderStepTrackingResponse {
	// 5-step tracker:
	// 1. Pesanan Dibuat
	// 2. Pembayaran
	// 3. Dikonfirmasi
	// 4. Siap Diambil
	// 5. Selesai

	step1Time := FormatTimeShortWIB(order.CreatedAt)
	steps := []OrderStepTrackingResponse{
		{
			StepNumber:      1,
			StepName:        "Pesanan Dibuat",
			StepDescription: "Pesanan berhasil dibuat",
			IsCompleted:     true,
			CompletedAt:     &step1Time,
		},
	}

	// Step 2: Pembayaran
	var step2Time *string
	isStep2Completed := false
	desc2 := "Menunggu pembayaran"
	if order.PaymentMethod == "CASH_ON_PICKUP" {
		isStep2Completed = true
		desc2 = "Bayar di Tempat (Kasir)"
		step2Time = &step1Time
	} else if order.PaidAt != nil {
		isStep2Completed = true
		desc2 = "Lunas"
		tStr := FormatTimeShortWIB(*order.PaidAt)
		step2Time = &tStr
	}
	steps = append(steps, OrderStepTrackingResponse{
		StepNumber:      2,
		StepName:        "Pembayaran",
		StepDescription: desc2,
		IsCompleted:     isStep2Completed,
		CompletedAt:     step2Time,
	})

	// Step 3: Dikonfirmasi
	var step3Time *string
	isStep3Completed := false
	desc3 := "Menunggu konfirmasi toko"
	if order.ConfirmedAt != nil {
		isStep3Completed = true
		desc3 = "Diterima Toko"
		tStr := FormatTimeShortWIB(*order.ConfirmedAt)
		step3Time = &tStr
	} else if order.Status == enums.OrderStatusInProcess || order.Status == enums.OrderStatusReadyForPickup || order.Status == enums.OrderStatusCompleted {
		isStep3Completed = true
		desc3 = "Diterima Toko"
	}
	steps = append(steps, OrderStepTrackingResponse{
		StepNumber:      3,
		StepName:        "Dikonfirmasi",
		StepDescription: desc3,
		IsCompleted:     isStep3Completed,
		CompletedAt:     step3Time,
	})

	// Step 4: Siap Diambil
	var step4Time *string
	isStep4Completed := false
	desc4 := "Sedang dipersiapkan"
	if order.ReadyAt != nil {
		isStep4Completed = true
		desc4 = "Di Rak Khusus"
		tStr := FormatTimeShortWIB(*order.ReadyAt)
		step4Time = &tStr
	} else if order.Status == enums.OrderStatusReadyForPickup || order.Status == enums.OrderStatusCompleted {
		isStep4Completed = true
		desc4 = "Di Rak Khusus"
	}
	steps = append(steps, OrderStepTrackingResponse{
		StepNumber:      4,
		StepName:        "Siap Diambil",
		StepDescription: desc4,
		IsCompleted:     isStep4Completed,
		CompletedAt:     step4Time,
	})

	// Step 5: Selesai
	var step5Time *string
	isStep5Completed := false
	desc5 := "Menunggu pengambilan"
	if order.CompletedAt != nil {
		isStep5Completed = true
		desc5 = "Diambil Pembeli"
		tStr := FormatTimeShortWIB(*order.CompletedAt)
		step5Time = &tStr
	} else if order.Status == enums.OrderStatusCompleted {
		isStep5Completed = true
		desc5 = "Diambil Pembeli"
	}
	steps = append(steps, OrderStepTrackingResponse{
		StepNumber:      5,
		StepName:        "Selesai",
		StepDescription: desc5,
		IsCompleted:     isStep5Completed,
		CompletedAt:     step5Time,
	})

	return steps
}

func FormatToOrderDetailResponse(order entity.Order, store entity.Store, payment *entity.Payment) OrderDetailResponse {
	itemResponses := make([]OrderItemResponse, 0, len(order.Items))
	for _, it := range order.Items {
		itemResponses = append(itemResponses, OrderItemResponse{
			ID:               it.ID,
			ProductID:        it.ProductID,
			ProductVariantID: it.ProductVariantID,
			ProductName:      it.ProductName,
			VariantName:      it.VariantName,
			ProductImageURL:  it.ProductImageURL,
			NormalPrice:      it.NormalPrice,
			RescuePrice:      it.RescuePrice,
			Quantity:         it.Quantity,
			Subtotal:         it.Subtotal,
		})
	}

	pickupStart := FormatTimeStr(order.PickupTimeStart)
	pickupEnd := FormatTimeStr(order.PickupTimeEnd)
	pickupTimeStr := fmt.Sprintf("%s – %s WIB", pickupStart, pickupEnd)

	var paymentResp *PaymentResponse
	if payment != nil {
		p := FormatToPaymentResponse(*payment)
		paymentResp = &p
	}

	return OrderDetailResponse{
		ID:                 order.ID,
		OrderNumber:        order.OrderNumber,
		StoreID:            order.StoreID,
		StoreName:          store.Name,
		StoreAddress:       store.Address,
		StoreCity:          store.City,
		StoreDistrict:      store.District,
		CustomerName:       order.CustomerName,
		CustomerPhone:      order.CustomerPhone,
		CustomerEmail:      order.CustomerEmail,
		PickupDate:         order.PickupDate.Format("2006-01-02"),
		PickupTime:         pickupTimeStr,
		TotalOriginalPrice: order.TotalOriginalPrice,
		TotalDiscount:      order.TotalDiscount,
		Subtotal:           order.Subtotal,
		ServiceFee:         order.ServiceFee,
		PackagingFee:       order.PackagingFee,
		TotalAmount:        order.TotalAmount,
		Status:             string(order.Status),
		StatusLabel:        FormatStatusLabel(order.Status),
		PaymentMethod:      order.PaymentMethod,
		Notes:              order.Notes,
		CreatedAt:          FormatTimeWIB(order.CreatedAt),
		Items:              itemResponses,
		Steps:              BuildOrderSteps(order),
		Payment:            paymentResp,
	}
}
