package entity

import (
	"time"
)

// Category represents the categories table
type Category struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
}

// CreateCategoryRequest represents create category request body
type CreateCategoryRequest struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
}

// UpdateCategoryRequest represents update category request body
type UpdateCategoryRequest struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
}
