package service

import (
	"context"
	"errors"

	"github.com/xprasetio/coffee-pos/internal/entity"
	"github.com/xprasetio/coffee-pos/internal/repository"
)

// CategoryService defines the interface for category business logic
type CategoryService interface {
	FindAll(ctx context.Context) ([]entity.Category, error)
	FindByID(ctx context.Context, id string) (*entity.Category, error)
	Create(ctx context.Context, req entity.CreateCategoryRequest) (*entity.Category, error)
	Update(ctx context.Context, id string, req entity.UpdateCategoryRequest) (*entity.Category, error)
	Delete(ctx context.Context, id string) error
}

// categoryService implements CategoryService
type categoryService struct {
	categoryRepo repository.CategoryRepository
}

// NewCategoryService creates a new CategoryService instance
func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{categoryRepo: categoryRepo}
}

// FindAll returns all categories
func (s *categoryService) FindAll(ctx context.Context) ([]entity.Category, error) {
	return s.categoryRepo.FindAll(ctx)
}

// FindByID finds a category by ID
func (s *categoryService) FindByID(ctx context.Context, id string) (*entity.Category, error) {
	result, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, errors.New("kategori tidak ditemukan")
	}
	return result, nil
}

// Create creates a new category
func (s *categoryService) Create(ctx context.Context, req entity.CreateCategoryRequest) (*entity.Category, error) {
	existing, err := s.categoryRepo.FindByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("nama kategori sudah digunakan")
	}

	category := &entity.Category{
		Name: req.Name,
	}

	if err := s.categoryRepo.Create(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

// Update updates an existing category
func (s *categoryService) Update(ctx context.Context, id string, req entity.UpdateCategoryRequest) (*entity.Category, error) {
	category, err := s.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	existing, err := s.categoryRepo.FindByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	if existing != nil && existing.ID != id {
		return nil, errors.New("nama kategori sudah digunakan")
	}

	category.Name = req.Name

	if err := s.categoryRepo.Update(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

// Delete deletes a category by ID
func (s *categoryService) Delete(ctx context.Context, id string) error {
	if _, err := s.FindByID(ctx, id); err != nil {
		return err
	}

	return s.categoryRepo.Delete(ctx, id)
}
