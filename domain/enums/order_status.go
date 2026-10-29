package enums

type OrderStatus string

const (
	OrderStatusWaitingPayment      OrderStatus = "WAITING_PAYMENT"
	OrderStatusWaitingConfirmation OrderStatus = "WAITING_CONFIRMATION"
	OrderStatusInProcess           OrderStatus = "IN_PROCESS"
	OrderStatusReadyForPickup      OrderStatus = "READY_FOR_PICKUP"
	OrderStatusCompleted           OrderStatus = "COMPLETED"
	OrderStatusCancelled           OrderStatus = "CANCELLED"
)

func (s OrderStatus) String() string {
	return string(s)
}

func (s OrderStatus) IsValid() bool {
	switch s {
	case OrderStatusWaitingPayment,
		OrderStatusWaitingConfirmation,
		OrderStatusInProcess,
		OrderStatusReadyForPickup,
		OrderStatusCompleted,
		OrderStatusCancelled:
		return true
	default:
		return false
	}
}
