package repository

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/xprasetio/coffee-pos/internal/entity"
	"github.com/google/uuid"
)

// StockFilter defines filters for stock movement queries
type StockFilter struct {
	ProductID string
	UserID    string
	Type      string
	Page      int
	Limit     int
}

// StockRepository defines the interface for stock movement data access
type StockRepository interface {
	Create(ctx context.Context, movement *entity.StockMovement) error
	FindByProductID(ctx context.Context, productID string, filter StockFilter) ([]entity.StockMovement, int, error)
	FindAll(ctx context.Context, filter StockFilter) ([]entity.StockMovement, int, error)
	WithTx(tx *sql.Tx) StockRepository
}

// stockRepository implements StockRepository
type stockRepository struct {
	db sqlDB
}

// NewStockRepository creates a new StockRepository instance
func NewStockRepository(db *sql.DB) StockRepository {
	return &stockRepository{db: db}
}

// WithTx returns a new StockRepository instance with the given transaction
func (r *stockRepository) WithTx(tx *sql.Tx) StockRepository {
	return &stockRepository{db: tx}
}

// Create creates a new stock movement record
func (r *stockRepository) Create(ctx context.Context, movement *entity.StockMovement) error {
	movement.ID = uuid.New().String()
	now := time.Now()
	movement.CreatedAt = now
	movement.UpdatedAt = now

	query := `INSERT INTO stock_movements (id, product_id, user_id, type, quantity,
			stock_before, stock_after, notes, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		movement.ID,
		movement.ProductID,
		movement.UserID,
		movement.Type,
		movement.Quantity,
		movement.StockBefore,
		movement.StockAfter,
		movement.Notes,
		movement.CreatedAt,
		movement.UpdatedAt,
	)

	return err
}

// FindByProductID returns stock movements for a specific product with filters and pagination
func (r *stockRepository) FindByProductID(ctx context.Context, productID string, filter StockFilter) ([]entity.StockMovement, int, error) {
	// Set default pagination values
	page := filter.Page
	if page <= 0 {
		page = 1
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	// Build WHERE conditions
	whereConditions := []string{"sm.product_id = ?"}
	args := []interface{}{productID}

	if filter.UserID != "" {
		whereConditions = append(whereConditions, "sm.user_id = ?")
		args = append(args, filter.UserID)
	}

	if filter.Type != "" {
		whereConditions = append(whereConditions, "sm.type = ?")
		args = append(args, filter.Type)
	}

	whereClause := strings.Join(whereConditions, " AND ")

	// Base query with JOINs
	baseQuery := `FROM stock_movements sm
		JOIN products p ON p.id = sm.product_id AND p.deleted_at IS NULL
		JOIN users u ON u.id = sm.user_id AND u.deleted_at IS NULL
		WHERE ` + whereClause

	// Count query for total
	countQuery := `SELECT COUNT(*) ` + baseQuery
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Data query
	dataQuery := `SELECT sm.id, sm.product_id, sm.user_id, sm.type, sm.quantity,
			sm.stock_before, sm.stock_after, sm.notes, sm.created_at, sm.updated_at,
			p.id, p.name,
			u.id, u.name
		` + baseQuery + ` ORDER BY sm.created_at DESC LIMIT ? OFFSET ?`

	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	movements := make([]entity.StockMovement, 0)
	for rows.Next() {
		var movement entity.StockMovement
		var product entity.Product
		var user entity.User

		err := rows.Scan(
			&movement.ID,
			&movement.ProductID,
			&movement.UserID,
			&movement.Type,
			&movement.Quantity,
			&movement.StockBefore,
			&movement.StockAfter,
			&movement.Notes,
			&movement.CreatedAt,
			&movement.UpdatedAt,
			&product.ID,
			&product.Name,
			&user.ID,
			&user.Name,
		)
		if err != nil {
			return nil, 0, err
		}

		movement.Product = &product
		movement.User = &user

		movements = append(movements, movement)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return movements, total, nil
}

// FindAll returns all stock movements with filters and pagination
func (r *stockRepository) FindAll(ctx context.Context, filter StockFilter) ([]entity.StockMovement, int, error) {
	// Set default pagination values
	page := filter.Page
	if page <= 0 {
		page = 1
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	// Build WHERE conditions
	whereConditions := make([]string, 0)
	args := make([]interface{}, 0)

	if filter.ProductID != "" {
		whereConditions = append(whereConditions, "sm.product_id = ?")
		args = append(args, filter.ProductID)
	}

	if filter.UserID != "" {
		whereConditions = append(whereConditions, "sm.user_id = ?")
		args = append(args, filter.UserID)
	}

	if filter.Type != "" {
		whereConditions = append(whereConditions, "sm.type = ?")
		args = append(args, filter.Type)
	}

	whereClause := "1=1"
	if len(whereConditions) > 0 {
		whereClause = strings.Join(whereConditions, " AND ")
	}

	// Base query with JOINs
	baseQuery := `FROM stock_movements sm
		JOIN products p ON p.id = sm.product_id AND p.deleted_at IS NULL
		JOIN users u ON u.id = sm.user_id AND u.deleted_at IS NULL
		WHERE ` + whereClause

	// Count query for total
	countQuery := `SELECT COUNT(*) ` + baseQuery
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Data query
	dataQuery := `SELECT sm.id, sm.product_id, sm.user_id, sm.type, sm.quantity,
			sm.stock_before, sm.stock_after, sm.notes, sm.created_at, sm.updated_at,
			p.id, p.name,
			u.id, u.name
		` + baseQuery + ` ORDER BY sm.created_at DESC LIMIT ? OFFSET ?`

	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	movements := make([]entity.StockMovement, 0)
	for rows.Next() {
		var movement entity.StockMovement
		var product entity.Product
		var user entity.User

		err := rows.Scan(
			&movement.ID,
			&movement.ProductID,
			&movement.UserID,
			&movement.Type,
			&movement.Quantity,
			&movement.StockBefore,
			&movement.StockAfter,
			&movement.Notes,
			&movement.CreatedAt,
			&movement.UpdatedAt,
			&product.ID,
			&product.Name,
			&user.ID,
			&user.Name,
		)
		if err != nil {
			return nil, 0, err
		}

		movement.Product = &product
		movement.User = &user

		movements = append(movements, movement)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return movements, total, nil
}
