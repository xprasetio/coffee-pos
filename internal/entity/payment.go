package entity

import "time"

// Payment methods
const (
	PaymentMethodCash    = "cash"
	PaymentMethodMidtrans = "midtrans"
)

// Payment statuses
const (
	PaymentStatusPending  = "pending"
	PaymentStatusPaid     = "paid"
	PaymentStatusFailed   = "failed"
	PaymentStatusExpired  = "expired"
)

// Payment represents a payment transaction
type Payment struct {
	ID              string     `json:"id"`
	OrderID         string     `json:"order_id"`
	Method          string     `json:"method"`
	Status          string     `json:"status"`
	Amount          int64      `json:"amount"`
	MidtransOrderID *string    `json:"midtrans_order_id,omitempty"`
	MidtransToken   *string    `json:"midtrans_token,omitempty"`
	MidtransURL     *string    `json:"midtrans_url,omitempty"`
	RawNotification *string    `json:"-"`
	PaidAt          *time.Time `json:"paid_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// CheckoutRequest represents a checkout request
type CheckoutRequest struct {
	Method string `json:"method" validate:"required,oneof=cash midtrans"`
}

// CheckoutResponse represents a checkout response
type CheckoutResponse struct {
	Payment       Payment `json:"payment"`
	MidtransToken string  `json:"midtrans_token,omitempty"`
	MidtransURL   string  `json:"midtrans_url,omitempty"`
}
