package service

import (
	"context"
	"errors"

	"github.com/xprasetio/coffee-pos/internal/entity"
	"github.com/xprasetio/coffee-pos/internal/repository"
)

// ProductService defines the interface for product business logic
type ProductService interface {
	FindAll(ctx context.Context, filter repository.ProductFilter) ([]entity.Product, int, error)
	FindByID(ctx context.Context, id string) (*entity.Product, error)
	Create(ctx context.Context, req entity.CreateProductRequest) (*entity.Product, error)
	Update(ctx context.Context, id string, req entity.UpdateProductRequest) (*entity.Product, error)
	Delete(ctx context.Context, id string) error
}

// productService implements ProductService
type productService struct {
	productRepo  repository.ProductRepository
	categoryRepo repository.CategoryRepository
}

// NewProductService creates a new ProductService instance
func NewProductService(
	productRepo repository.ProductRepository,
	categoryRepo repository.CategoryRepository,
) ProductService {
	return &productService{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
	}
}

// FindAll returns all products with filtering and pagination
func (s *productService) FindAll(ctx context.Context, filter repository.ProductFilter) ([]entity.Product, int, error) {
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}

	return s.productRepo.FindAll(ctx, filter)
}

// FindByID finds a product by ID
func (s *productService) FindByID(ctx context.Context, id string) (*entity.Product, error) {
	result, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, errors.New("produk tidak ditemukan")
	}
	return result, nil
}

// Create creates a new product
func (s *productService) Create(ctx context.Context, req entity.CreateProductRequest) (*entity.Product, error) {
	category, err := s.categoryRepo.FindByID(ctx, req.CategoryID)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New("kategori tidak ditemukan")
	}

	product := &entity.Product{
		SKU:         req.SKU,
		Name:        req.Name,
		Description: req.Description,
		CategoryID:  req.CategoryID,
		Price:       req.Price,
		Stock:       req.Stock,
		MinStock:    req.MinStock,
		IsActive:    true,
	}

	if err := s.productRepo.Create(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

// Update updates an existing product
func (s *productService) Update(ctx context.Context, id string, req entity.UpdateProductRequest) (*entity.Product, error) {
	product, err := s.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.CategoryID != "" {
		category, err := s.categoryRepo.FindByID(ctx, req.CategoryID)
		if err != nil {
			return nil, err
		}
		if category == nil {
			return nil, errors.New("kategori tidak ditemukan")
		}
		product.CategoryID = req.CategoryID
	}

	if req.Name != "" {
		product.Name = req.Name
	}

	if req.Description != "" {
		product.Description = req.Description
	}

	if req.Price > 0 {
		product.Price = req.Price
	}

	if req.IsActive {
		product.IsActive = req.IsActive
	}

	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

// Delete deletes a product by ID
func (s *productService) Delete(ctx context.Context, id string) error {
	if _, err := s.FindByID(ctx, id); err != nil {
		return err
	}

	return s.productRepo.Delete(ctx, id)
}
