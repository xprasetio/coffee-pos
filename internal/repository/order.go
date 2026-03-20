package repository

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/xprasetio/coffee-pos/internal/entity"
	"github.com/google/uuid"
)

// OrderFilter defines filters for order queries
type OrderFilter struct {
	ShiftID   string
	CashierID string
	Status    string
	Page      int
	Limit     int
}

// OrderRepository defines the interface for order data access
type OrderRepository interface {
	FindByID(ctx context.Context, id string) (*entity.Order, error)
	FindAll(ctx context.Context, filter OrderFilter) ([]entity.Order, int, error)
	Create(ctx context.Context, order *entity.Order) error
	Update(ctx context.Context, order *entity.Order) error
	AddItem(ctx context.Context, item *entity.OrderItem) error
	UpdateItem(ctx context.Context, item *entity.OrderItem) error
	DeleteItem(ctx context.Context, itemID string) error
	FindItemByID(ctx context.Context, itemID string) (*entity.OrderItem, error)
	WithTx(tx *sql.Tx) OrderRepository
}

// orderRepository implements OrderRepository
type orderRepository struct {
	db sqlDB
}

// NewOrderRepository creates a new OrderRepository instance
func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepository{db: db}
}

// WithTx returns a new OrderRepository instance with the given transaction
func (r *orderRepository) WithTx(tx *sql.Tx) OrderRepository {
	return &orderRepository{db: tx}
}

// FindByID finds an order by ID with items loaded via separate query
func (r *orderRepository) FindByID(ctx context.Context, id string) (*entity.Order, error) {
	// Query order with cashier and optional table/promo
	query := `SELECT o.id, o.shift_id, o.cashier_id, o.table_id, o.promo_id,
			o.status, o.subtotal, o.discount_amount, o.total, o.notes,
			o.created_at, o.updated_at,
			u.id, u.name,
			t.id, t.name,
			pm.id, pm.name
		FROM orders o
		LEFT JOIN users u ON u.id = o.cashier_id AND u.deleted_at IS NULL
		LEFT JOIN tables t ON t.id = o.table_id AND t.deleted_at IS NULL
		LEFT JOIN promos pm ON pm.id = o.promo_id AND pm.deleted_at IS NULL
		WHERE o.id = ? AND o.deleted_at IS NULL`

	order := &entity.Order{}
	var tableID, tableName sql.NullString
	var promoID, promoName sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&order.ID,
		&order.ShiftID,
		&order.CashierID,
		&order.TableID,
		&order.PromoID,
		&order.Status,
		&order.Subtotal,
		&order.DiscountAmount,
		&order.Total,
		&order.Notes,
		&order.CreatedAt,
		&order.UpdatedAt,
		&order.Cashier.ID,
		&order.Cashier.Name,
		&tableID,
		&tableName,
		&promoID,
		&promoName,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Build table info if exists
	if tableID.Valid {
		order.Table = &entity.Table{
			ID:   tableID.String,
			Name: tableName.String,
		}
	}

	// Build promo info if exists
	if promoID.Valid {
		order.Promo = &entity.Promo{
			ID:   promoID.String,
			Name: promoName.String,
		}
	}

	// Query order items separately
	itemsQuery := `SELECT oi.id, oi.order_id, oi.product_id, oi.quantity, oi.price, oi.subtotal,
			oi.notes, oi.created_at, oi.updated_at,
			p.id, p.name, p.price
		FROM order_items oi
		JOIN products p ON p.id = oi.product_id AND p.deleted_at IS NULL
		WHERE oi.order_id = ? AND oi.deleted_at IS NULL`

	rows, err := r.db.QueryContext(ctx, itemsQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]entity.OrderItem, 0)
	for rows.Next() {
		var item entity.OrderItem
		var product entity.Product

		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.Quantity,
			&item.Price,
			&item.Subtotal,
			&item.Notes,
			&item.CreatedAt,
			&item.UpdatedAt,
			&product.ID,
			&product.Name,
			&product.Price,
		)
		if err != nil {
			return nil, err
		}

		item.Product = &product
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	order.Items = items

	return order, nil
}

// FindAll returns all orders with filters and pagination (without items)
func (r *orderRepository) FindAll(ctx context.Context, filter OrderFilter) ([]entity.Order, int, error) {
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
	whereConditions := []string{"o.deleted_at IS NULL"}
	args := make([]interface{}, 0)

	if filter.ShiftID != "" {
		whereConditions = append(whereConditions, "o.shift_id = ?")
		args = append(args, filter.ShiftID)
	}

	if filter.CashierID != "" {
		whereConditions = append(whereConditions, "o.cashier_id = ?")
		args = append(args, filter.CashierID)
	}

	if filter.Status != "" {
		whereConditions = append(whereConditions, "o.status = ?")
		args = append(args, filter.Status)
	}

	whereClause := strings.Join(whereConditions, " AND ")

	// Base query with JOINs
	baseQuery := `FROM orders o
		LEFT JOIN users u ON u.id = o.cashier_id AND u.deleted_at IS NULL
		WHERE ` + whereClause

	// Count query for total
	countQuery := `SELECT COUNT(*) ` + baseQuery
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Data query (without items)
	dataQuery := `SELECT o.id, o.shift_id, o.cashier_id, o.table_id, o.promo_id,
			o.status, o.subtotal, o.discount_amount, o.total, o.notes,
			o.created_at, o.updated_at,
			u.id, u.name
		` + baseQuery + ` ORDER BY o.created_at DESC LIMIT ? OFFSET ?`

	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	orders := make([]entity.Order, 0)
	for rows.Next() {
		var order entity.Order

		err := rows.Scan(
			&order.ID,
			&order.ShiftID,
			&order.CashierID,
			&order.TableID,
			&order.PromoID,
			&order.Status,
			&order.Subtotal,
			&order.DiscountAmount,
			&order.Total,
			&order.Notes,
			&order.CreatedAt,
			&order.UpdatedAt,
			&order.Cashier.ID,
			&order.Cashier.Name,
		)
		if err != nil {
			return nil, 0, err
		}

		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// Create creates a new order
func (r *orderRepository) Create(ctx context.Context, order *entity.Order) error {
	order.ID = uuid.New().String()
	now := time.Now()
	order.CreatedAt = now
	order.UpdatedAt = now

	query := `INSERT INTO orders (id, shift_id, cashier_id, table_id, promo_id,
			status, subtotal, discount_amount, total, notes, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		order.ID,
		order.ShiftID,
		order.CashierID,
		order.TableID,
		order.PromoID,
		order.Status,
		order.Subtotal,
		order.DiscountAmount,
		order.Total,
		order.Notes,
		order.CreatedAt,
		order.UpdatedAt,
	)

	return err
}

// Update updates an existing order
func (r *orderRepository) Update(ctx context.Context, order *entity.Order) error {
	order.UpdatedAt = time.Now()

	query := `UPDATE orders SET table_id = ?, promo_id = ?, status = ?,
			subtotal = ?, discount_amount = ?, total = ?, notes = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL`

	_, err := r.db.ExecContext(ctx, query,
		order.TableID,
		order.PromoID,
		order.Status,
		order.Subtotal,
		order.DiscountAmount,
		order.Total,
		order.Notes,
		order.UpdatedAt,
		order.ID,
	)

	return err
}

// AddItem adds a new item to an order
func (r *orderRepository) AddItem(ctx context.Context, item *entity.OrderItem) error {
	item.ID = uuid.New().String()
	now := time.Now()
	item.CreatedAt = now
	item.UpdatedAt = now

	// Calculate subtotal in repository
	item.Subtotal = item.Price * int64(item.Quantity)

	query := `INSERT INTO order_items (id, order_id, product_id, quantity, price,
			subtotal, notes, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		item.ID,
		item.OrderID,
		item.ProductID,
		item.Quantity,
		item.Price,
		item.Subtotal,
		item.Notes,
		item.CreatedAt,
		item.UpdatedAt,
	)

	return err
}

// UpdateItem updates an existing order item
func (r *orderRepository) UpdateItem(ctx context.Context, item *entity.OrderItem) error {
	item.UpdatedAt = time.Now()

	query := `UPDATE order_items SET quantity = ?, notes = ?, subtotal = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL`

	_, err := r.db.ExecContext(ctx, query,
		item.Quantity,
		item.Notes,
		item.Subtotal,
		item.UpdatedAt,
		item.ID,
	)

	return err
}

// DeleteItem hard deletes an order item
func (r *orderRepository) DeleteItem(ctx context.Context, itemID string) error {
	query := `DELETE FROM order_items WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, itemID)

	return err
}

// FindItemByID finds an order item by ID (without product JOIN)
func (r *orderRepository) FindItemByID(ctx context.Context, itemID string) (*entity.OrderItem, error) {
	query := `SELECT id, order_id, product_id, quantity, price, subtotal,
			notes, created_at, updated_at
		FROM order_items WHERE id = ? AND deleted_at IS NULL`

	item := &entity.OrderItem{}

	err := r.db.QueryRowContext(ctx, query, itemID).Scan(
		&item.ID,
		&item.OrderID,
		&item.ProductID,
		&item.Quantity,
		&item.Price,
		&item.Subtotal,
		&item.Notes,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return item, nil
}
