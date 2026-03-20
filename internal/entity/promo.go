package entity

import "time"

// Promo types
const (
	PromoTypePercentage = "percentage"
	PromoTypeFixed      = "fixed"
)

// Promo represents a promotional offer in the coffee shop POS system
type Promo struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Code      string     `json:"code"`
	Type      string     `json:"type"`
	Value     int64      `json:"value"` // percentage (e.g., 10 = 10%) or nominal in cents
	MinOrder  int64      `json:"min_order"`
	MaxDiscount *int64   `json:"max_discount,omitempty"` // cap for percentage, nullable
	UsageLimit  *int     `json:"usage_limit,omitempty"`  // NULL = unlimited
	UsedCount   int      `json:"used_count"`
	StartedAt   time.Time `json:"started_at"`
	EndedAt     time.Time `json:"ended_at"`
	IsActive    bool     `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"-"`
}

// IsValid checks if the promo is currently valid
func (p *Promo) IsValid() bool {
	if !p.IsActive {
		return false
	}

	now := time.Now()
	if now.Before(p.StartedAt) || now.After(p.EndedAt) {
		return false
	}

	if p.UsageLimit != nil && p.UsedCount >= *p.UsageLimit {
		return false
	}

	return true
}

// CreatePromoRequest represents the request to create a new promo
type CreatePromoRequest struct {
	Name       string     `json:"name" validate:"required,min=2,max=200"`
	Code       string     `json:"code" validate:"required,min=3,max=50"`
	Type       string     `json:"type" validate:"required,oneof=percentage fixed"`
	Value      int64      `json:"value" validate:"required,min=1"`
	MinOrder   int64      `json:"min_order" validate:"min=0"`
	MaxDiscount *int64    `json:"max_discount" validate:"omitempty,min=1"`
	UsageLimit  *int      `json:"usage_limit" validate:"omitempty,min=1"`
	StartedAt   time.Time `json:"started_at" validate:"required"`
	EndedAt     time.Time `json:"ended_at" validate:"required"`
}

// UpdatePromoRequest represents the request to update an existing promo
// All fields are optional
type UpdatePromoRequest struct {
	Name      string     `json:"name" validate:"omitempty,min=2,max=200"`
	IsActive  *bool      `json:"is_active"`
	StartedAt *time.Time `json:"started_at"`
	EndedAt   *time.Time `json:"ended_at"`
}
