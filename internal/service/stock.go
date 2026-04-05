package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/xprasetio/coffee-pos/internal/entity"
	"github.com/xprasetio/coffee-pos/internal/repository"
	"github.com/xprasetio/coffee-pos/pkg/txmanager"
)

// StockService defines the interface for stock business logic
type StockService interface {
	GetStock(ctx context.Context, productID string) (*entity.Product, error)
	Adjust(ctx context.Context, productID string, userID string, req entity.StockAdjustmentRequest) error
}

// stockService implements StockService
type stockService struct {
	stockRepo   repository.StockRepository
	productRepo repository.ProductRepository
	txManager   *txmanager.TxManager
}

// NewStockService creates a new StockService instance
func NewStockService(
	stockRepo repository.StockRepository,
	productRepo repository.ProductRepository,
	txManager *txmanager.TxManager,
) StockService {
	return &stockService{
		stockRepo:   stockRepo,
		productRepo: productRepo,
		txManager:   txManager,
	}
}

// GetStock returns the current stock of a product
func (s *stockService) GetStock(ctx context.Context, productID string) (*entity.Product, error) {
	product, err := s.productRepo.FindByID(ctx, productID)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("produk tidak ditemukan")
	}
	return product, nil
}

// Adjust adjusts the stock of a product
func (s *stockService) Adjust(ctx context.Context, productID string, userID string, req entity.StockAdjustmentRequest) error {
	product, err := s.productRepo.FindByID(ctx, productID)
	if err != nil {
		return err
	}
	if product == nil {
		return errors.New("produk tidak ditemukan")
	}

	// Calculate signed quantity and new stock based on type
	var signedQty int
	var newStock int

	switch req.Type {
	case entity.MovementTypeIn, entity.MovementTypeAdjustment:
		signedQty = req.Quantity
		newStock = product.Stock + req.Quantity
	case entity.MovementTypeOut:
		signedQty = req.Quantity * -1
		newStock = product.Stock - req.Quantity
	}

	if newStock < 0 {
		return errors.New("stok tidak cukup")
	}

	return s.txManager.WithTx(ctx, func(tx *sql.Tx) error {
		txStockRepo := s.stockRepo.WithTx(tx)
		txProductRepo := s.productRepo.WithTx(tx)

		movement := &entity.StockMovement{
			ProductID:   productID,
			UserID:      userID,
			Type:        req.Type,
			Quantity:    signedQty,
			StockBefore: product.Stock,
			StockAfter:  newStock,
			Notes:       req.Notes,
		}

		if err := txStockRepo.Create(ctx, movement); err != nil {
			return err
		}

		return txProductRepo.UpdateStock(ctx, productID, newStock)
	})
}
