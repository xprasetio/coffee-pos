package service

import (
	"context"
	"errors"

	"github.com/xprasetio/coffee-pos/internal/entity"
	"github.com/xprasetio/coffee-pos/internal/repository"
)

// TableService handles table business logic
type TableService struct {
	tableRepo repository.TableRepository
}

// NewTableService creates a new TableService instance
func NewTableService(tableRepo repository.TableRepository) *TableService {
	return &TableService{tableRepo: tableRepo}
}

// FindAll returns all tables
func (s *TableService) FindAll(ctx context.Context) ([]entity.Table, error) {
	return s.tableRepo.FindAll(ctx)
}

// FindByID finds a table by ID
func (s *TableService) FindByID(ctx context.Context, id string) (*entity.Table, error) {
	table, err := s.tableRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if table == nil {
		return nil, errors.New("meja tidak ditemukan")
	}
	return table, nil
}

// Create creates a new table
func (s *TableService) Create(ctx context.Context, req entity.CreateTableRequest) (*entity.Table, error) {
	existing, err := s.tableRepo.FindByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("nama meja sudah digunakan")
	}

	table := &entity.Table{
		Name:     req.Name,
		Capacity: req.Capacity,
		Status:   entity.TableStatusAvailable,
	}

	if err := s.tableRepo.Create(ctx, table); err != nil {
		return nil, err
	}

	return table, nil
}

// Update updates an existing table
func (s *TableService) Update(ctx context.Context, id string, req entity.UpdateTableRequest) (*entity.Table, error) {
	table, err := s.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	needUpdate := false

	if req.Name != "" {
		existing, err := s.tableRepo.FindByName(ctx, req.Name)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.ID != id {
			return nil, errors.New("nama meja sudah digunakan")
		}
		table.Name = req.Name
		needUpdate = true
	}

	if req.Capacity > 0 {
		table.Capacity = req.Capacity
		needUpdate = true
	}

	if needUpdate {
		if err := s.tableRepo.Update(ctx, table); err != nil {
			return nil, err
		}
	}

	if req.Status != "" {
		if err := s.tableRepo.UpdateStatus(ctx, id, req.Status); err != nil {
			return nil, err
		}
		table.Status = req.Status
	}

	return table, nil
}

// Delete deletes a table by ID
func (s *TableService) Delete(ctx context.Context, id string) error {
	table, err := s.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if table == nil {
		return errors.New("meja tidak ditemukan")
	}

	return s.tableRepo.Delete(ctx, id)
}
