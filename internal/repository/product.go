package repository

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/xprasetio/coffee-pos/internal/entity"
	"github.com/google/uuid"
)

// ProductFilter defines filters for product queries
type ProductFilter struct {
	CategoryID string
	IsActive   *bool
	Search     string
	Page       int
	Limit      int
}

// ProductRepository defines the interface for product data access
type ProductRepository interface {
	FindByID(ctx context.Context, id string) (*entity.Product, error)
	FindAll(ctx context.Context, filter ProductFilter) ([]entity.Product, int, error)
	Create(ctx context.Context, product *entity.Product) error
	Update(ctx context.Context, product *entity.Product) error
	Delete(ctx context.Context, id string) error
	UpdateStock(ctx context.Context, id string, stock int) error
	WithTx(tx *sql.Tx) ProductRepository
}

// productRepository implements ProductRepository
type productRepository struct {
	db sqlDB
}

// NewProductRepository creates a new ProductRepository instance
func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db: db}
}

// WithTx returns a new ProductRepository instance with the given transaction
func (r *productRepository) WithTx(tx *sql.Tx) ProductRepository {
	return &productRepository{db: tx}
}

// FindByID finds a product by ID with category joined
func (r *productRepository) FindByID(ctx context.Context, id string) (*entity.Product, error) {
	query := `SELECT p.id, p.sku, p.name, p.description, p.category_id, p.price, p.stock, p.min_stock,
			p.is_active, p.created_at, p.updated_at,
			c.id, c.name
		FROM products p
		LEFT JOIN categories c ON c.id = p.category_id AND c.deleted_at IS NULL
		WHERE p.id = ? AND p.deleted_at IS NULL`

	product := &entity.Product{}
	var categoryID sql.NullString
	var categoryName sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&product.ID,
		&product.SKU,
		&product.Name,
		&product.Description,
		&product.CategoryID,
		&product.Price,
		&product.Stock,
		&product.MinStock,
		&product.IsActive,
		&product.CreatedAt,
		&product.UpdatedAt,
		&categoryID,
		&categoryName,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if categoryID.Valid {
		product.Category = &entity.Category{
			ID:   categoryID.String,
			Name: categoryName.String,
		}
	}

	return product, nil
}

// FindAll returns all products with filters, pagination, and category joined
func (r *productRepository) FindAll(ctx context.Context, filter ProductFilter) ([]entity.Product, int, error) {
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
	whereConditions := []string{"p.deleted_at IS NULL"}
	args := make([]interface{}, 0)

	if filter.CategoryID != "" {
		whereConditions = append(whereConditions, "p.category_id = ?")
		args = append(args, filter.CategoryID)
	}

	if filter.IsActive != nil {
		whereConditions = append(whereConditions, "p.is_active = ?")
		args = append(args, *filter.IsActive)
	}

	if filter.Search != "" {
		whereConditions = append(whereConditions, "p.name LIKE ?")
		args = append(args, "%"+filter.Search+"%")
	}

	whereClause := strings.Join(whereConditions, " AND ")

	// Base query with JOIN
	baseQuery := `FROM products p
		LEFT JOIN categories c ON c.id = p.category_id AND c.deleted_at IS NULL
		WHERE ` + whereClause

	// Count query for total
	countQuery := `SELECT COUNT(*) ` + baseQuery
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Data query
	dataQuery := `SELECT p.id, p.sku, p.name, p.description, p.category_id, p.price, p.stock, p.min_stock,
			p.is_active, p.created_at, p.updated_at,
			c.id, c.name
		` + baseQuery + ` ORDER BY p.created_at DESC LIMIT ? OFFSET ?`

	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	products := make([]entity.Product, 0)
	for rows.Next() {
		var product entity.Product
		var categoryID sql.NullString
		var categoryName sql.NullString

		err := rows.Scan(
			&product.ID,
			&product.SKU,
			&product.Name,
			&product.Description,
			&product.CategoryID,
			&product.Price,
			&product.Stock,
			&product.MinStock,
			&product.IsActive,
			&product.CreatedAt,
			&product.UpdatedAt,
			&categoryID,
			&categoryName,
		)
		if err != nil {
			return nil, 0, err
		}

		if categoryID.Valid {
			product.Category = &entity.Category{
				ID:   categoryID.String,
				Name: categoryName.String,
			}
		}

		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// Create creates a new product
func (r *productRepository) Create(ctx context.Context, product *entity.Product) error {
	product.ID = uuid.New().String()
	now := time.Now()
	product.CreatedAt = now
	product.UpdatedAt = now

	query := `INSERT INTO products (id, sku, name, description, category_id, price, stock, min_stock,
			is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		product.ID,
		product.SKU,
		product.Name,
		product.Description,
		product.CategoryID,
		product.Price,
		product.Stock,
		product.MinStock,
		product.IsActive,
		product.CreatedAt,
		product.UpdatedAt,
	)

	return err
}

// Update updates an existing product
func (r *productRepository) Update(ctx context.Context, product *entity.Product) error {
	product.UpdatedAt = time.Now()

	query := `UPDATE products SET name = ?, description = ?, price = ?,
			is_active = ?, category_id = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL`

	_, err := r.db.ExecContext(ctx, query,
		product.Name,
		product.Description,
		product.Price,
		product.IsActive,
		product.CategoryID,
		product.UpdatedAt,
		product.ID,
	)

	return err
}

// Delete performs a soft delete on a product
func (r *productRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE products SET deleted_at = NOW() WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, id)

	return err
}

// UpdateStock updates the stock of a product
// NOTE: This method should never be called directly from handlers.
// It should always be called together with StockRepository.Create() within
// a single database transaction.
func (r *productRepository) UpdateStock(ctx context.Context, id string, stock int) error {
	query := `UPDATE products SET stock = ?, updated_at = ? WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, stock, time.Now(), id)

	return err
}
