package entity

import (
	"time"
)

// Product represents the products table
type Product struct {
	ID          string     `json:"id"`
	SKU         string     `json:"sku"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	CategoryID  string     `json:"category_id"`
	Category    *Category  `json:"category,omitempty"`
	Price       int64      `json:"price"`
	Stock       int        `json:"stock"`
	MinStock    int        `json:"min_stock"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"-"`
}

// CreateProductRequest represents create product request body
type CreateProductRequest struct {
	SKU         string `json:"sku" validate:"required,max=20"`
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Description string `json:"description" validate:"max=1000"`
	CategoryID  string `json:"category_id" validate:"required,uuid4"`
	Price       int64  `json:"price" validate:"required,min=1"`
	Stock       int    `json:"stock" validate:"min=0"`
	MinStock    int    `json:"min_stock" validate:"min=1"`
	IsActive    bool   `json:"is_active"`
}

// UpdateProductRequest represents update product request body
type UpdateProductRequest struct {
	SKU         string `json:"sku" validate:"max=20"`
	Name        string `json:"name" validate:"min=2,max=100"`
	Description string `json:"description" validate:"max=1000"`
	CategoryID  string `json:"category_id" validate:"uuid4"`
	Price       int64  `json:"price" validate:"min=1"`
	Stock       int    `json:"stock" validate:"min=0"`
	MinStock    int    `json:"min_stock" validate:"min=1"`
	IsActive    bool   `json:"is_active"`
}

// ProductResponse represents product data sent to client
type ProductResponse struct {
	ID          string       `json:"id"`
	SKU         string       `json:"sku"`
	Name        string       `json:"name"`
	Description string       `json:"description,omitempty"`
	CategoryID  string       `json:"category_id"`
	Category    *Category    `json:"category,omitempty"`
	Price       int64        `json:"price"`
	Stock       int          `json:"stock"`
	MinStock    int          `json:"min_stock"`
	IsActive    bool         `json:"is_active"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// ToResponse converts Product to ProductResponse
func (p *Product) ToResponse() ProductResponse {
	return ProductResponse{
		ID:          p.ID,
		SKU:         p.SKU,
		Name:        p.Name,
		Description: p.Description,
		CategoryID:  p.CategoryID,
		Category:    p.Category,
		Price:       p.Price,
		Stock:       p.Stock,
		MinStock:    p.MinStock,
		IsActive:    p.IsActive,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
