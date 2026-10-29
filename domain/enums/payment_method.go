package enums

type PaymentMethod string

const (
	PaymentMethodOnline       PaymentMethod = "ONLINE"
	PaymentMethodCashOnPickup PaymentMethod = "CASH_ON_PICKUP"
)

func (m PaymentMethod) String() string {
	return string(m)
}

func (m PaymentMethod) IsValid() bool {
	switch m {
	case PaymentMethodOnline, PaymentMethodCashOnPickup:
		return true
	default:
		return false
	}
}

type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "PENDING"
	PaymentStatusSettlement PaymentStatus = "SETTLEMENT"
	PaymentStatusExpire     PaymentStatus = "EXPIRE"
	PaymentStatusCancel     PaymentStatus = "CANCEL"
	PaymentStatusDeny       PaymentStatus = "DENY"
	PaymentStatusRefund     PaymentStatus = "REFUND"
)

func (s PaymentStatus) String() string {
	return string(s)
}
