package order_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/rescuebyte/rescuebyte-be/domain/dto"
	"github.com/rescuebyte/rescuebyte-be/domain/entity"
	"github.com/rescuebyte/rescuebyte-be/domain/enums"
)

func TestBuildOrderSteps(t *testing.T) {
	now := time.Now()
	paidAt := now.Add(2 * time.Minute)
	confirmedAt := now.Add(5 * time.Minute)
	readyAt := now.Add(45 * time.Minute)
	completedAt := now.Add(90 * time.Minute)

	order := entity.Order{
		ID:            uuid.New(),
		OrderNumber:   "RB251026-00123",
		Status:        enums.OrderStatusCompleted,
		PaymentMethod: "ONLINE",
		CreatedAt:     now,
		PaidAt:        &paidAt,
		ConfirmedAt:   &confirmedAt,
		ReadyAt:       &readyAt,
		CompletedAt:   &completedAt,
	}

	steps := dto.BuildOrderSteps(order)

	assert.Equal(t, 5, len(steps))

	// All 5 steps should be completed
	assert.True(t, steps[0].IsCompleted)
	assert.Equal(t, "Pesanan Dibuat", steps[0].StepName)

	assert.True(t, steps[1].IsCompleted)
	assert.Equal(t, "Pembayaran", steps[1].StepName)
	assert.Equal(t, "Lunas", steps[1].StepDescription)

	assert.True(t, steps[2].IsCompleted)
	assert.Equal(t, "Dikonfirmasi", steps[2].StepName)

	assert.True(t, steps[3].IsCompleted)
	assert.Equal(t, "Siap Diambil", steps[3].StepName)

	assert.True(t, steps[4].IsCompleted)
	assert.Equal(t, "Selesai", steps[4].StepName)
}

func TestBuildOrderStepsCashOnPickup(t *testing.T) {
	now := time.Now()
	order := entity.Order{
		ID:            uuid.New(),
		OrderNumber:   "RB251026-00124",
		Status:        enums.OrderStatusWaitingConfirmation,
		PaymentMethod: "CASH_ON_PICKUP",
		CreatedAt:     now,
	}

	steps := dto.BuildOrderSteps(order)

	assert.Equal(t, 5, len(steps))
	assert.True(t, steps[0].IsCompleted)
	assert.True(t, steps[1].IsCompleted) // Cash on pickup is automatically ready for confirmation
	assert.Equal(t, "Bayar di Tempat (Kasir)", steps[1].StepDescription)
	assert.False(t, steps[2].IsCompleted)
	assert.False(t, steps[3].IsCompleted)
	assert.False(t, steps[4].IsCompleted)
}

func TestFormatStatusLabel(t *testing.T) {
	assert.Equal(t, "Menunggu Pembayaran", dto.FormatStatusLabel(enums.OrderStatusWaitingPayment))
	assert.Equal(t, "Menunggu Konfirmasi", dto.FormatStatusLabel(enums.OrderStatusWaitingConfirmation))
	assert.Equal(t, "Diproses", dto.FormatStatusLabel(enums.OrderStatusInProcess))
	assert.Equal(t, "Siap Diambil", dto.FormatStatusLabel(enums.OrderStatusReadyForPickup))
	assert.Equal(t, "Selesai", dto.FormatStatusLabel(enums.OrderStatusCompleted))
	assert.Equal(t, "Dibatalkan", dto.FormatStatusLabel(enums.OrderStatusCancelled))
}
