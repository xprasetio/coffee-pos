package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/xprasetio/coffee-pos/internal/entity"
	"github.com/google/uuid"
)

// PaymentRepository defines the interface for payment data access
type PaymentRepository interface {
	FindByID(ctx context.Context, id string) (*entity.Payment, error)
	FindByOrderID(ctx context.Context, orderID string) (*entity.Payment, error)
	FindByMidtransOrderID(ctx context.Context, midtransOrderID string) (*entity.Payment, error)
	Create(ctx context.Context, payment *entity.Payment) error
	UpdateStatus(ctx context.Context, id string, status string, paidAt *time.Time) error
	UpdateMidtransData(ctx context.Context, id string, token string, url string, midtransOrderID string) error
	SaveRawNotification(ctx context.Context, id string, raw string) error
	WithTx(tx *sql.Tx) PaymentRepository
}

// paymentRepository implements PaymentRepository
type paymentRepository struct {
	db sqlDB
}

// NewPaymentRepository creates a new PaymentRepository instance
func NewPaymentRepository(db *sql.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

// WithTx returns a new PaymentRepository instance with the given transaction
func (r *paymentRepository) WithTx(tx *sql.Tx) PaymentRepository {
	return &paymentRepository{db: tx}
}

// FindByID finds a payment by ID
func (r *paymentRepository) FindByID(ctx context.Context, id string) (*entity.Payment, error) {
	query := `SELECT id, order_id, method, status, amount, midtrans_order_id,
			midtrans_token, midtrans_url, raw_notification, paid_at,
			created_at, updated_at
		FROM payments WHERE id = ? AND deleted_at IS NULL`

	payment := &entity.Payment{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.Method,
		&payment.Status,
		&payment.Amount,
		&payment.MidtransOrderID,
		&payment.MidtransToken,
		&payment.MidtransURL,
		&payment.RawNotification,
		&payment.PaidAt,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return payment, nil
}

// FindByOrderID finds a payment by order ID
func (r *paymentRepository) FindByOrderID(ctx context.Context, orderID string) (*entity.Payment, error) {
	query := `SELECT id, order_id, method, status, amount, midtrans_order_id,
			midtrans_token, midtrans_url, raw_notification, paid_at,
			created_at, updated_at
		FROM payments WHERE order_id = ? AND deleted_at IS NULL`

	payment := &entity.Payment{}

	err := r.db.QueryRowContext(ctx, query, orderID).Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.Method,
		&payment.Status,
		&payment.Amount,
		&payment.MidtransOrderID,
		&payment.MidtransToken,
		&payment.MidtransURL,
		&payment.RawNotification,
		&payment.PaidAt,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return payment, nil
}

// FindByMidtransOrderID finds a payment by midtrans order ID
func (r *paymentRepository) FindByMidtransOrderID(ctx context.Context, midtransOrderID string) (*entity.Payment, error) {
	query := `SELECT id, order_id, method, status, amount, midtrans_order_id,
			midtrans_token, midtrans_url, raw_notification, paid_at,
			created_at, updated_at
		FROM payments WHERE midtrans_order_id = ? AND deleted_at IS NULL`

	payment := &entity.Payment{}

	err := r.db.QueryRowContext(ctx, query, midtransOrderID).Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.Method,
		&payment.Status,
		&payment.Amount,
		&payment.MidtransOrderID,
		&payment.MidtransToken,
		&payment.MidtransURL,
		&payment.RawNotification,
		&payment.PaidAt,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return payment, nil
}

// Create creates a new payment record
// NOTE: The order_id column in payments table is UNIQUE — one order
// can only have one payment record. If Create is called twice for the
// same order, MySQL will throw a duplicate key error. Service layer
// must check FindByOrderID before calling Create.
func (r *paymentRepository) Create(ctx context.Context, payment *entity.Payment) error {
	payment.ID = uuid.New().String()
	now := time.Now()
	payment.CreatedAt = now
	payment.UpdatedAt = now

	if payment.Status == "" {
		payment.Status = entity.PaymentStatusPending
	}

	query := `INSERT INTO payments (id, order_id, method, status, amount,
			midtrans_order_id, midtrans_token, midtrans_url,
			raw_notification, paid_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		payment.ID,
		payment.OrderID,
		payment.Method,
		payment.Status,
		payment.Amount,
		payment.MidtransOrderID,
		payment.MidtransToken,
		payment.MidtransURL,
		payment.RawNotification,
		payment.PaidAt,
		payment.CreatedAt,
		payment.UpdatedAt,
	)

	return err
}

// UpdateStatus updates the status and paid_at of a payment
func (r *paymentRepository) UpdateStatus(ctx context.Context, id string, status string, paidAt *time.Time) error {
	query := `UPDATE payments SET status = ?, paid_at = ?, updated_at = NOW() WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, status, paidAt, id)

	return err
}

// UpdateMidtransData updates the midtrans-related fields of a payment
func (r *paymentRepository) UpdateMidtransData(ctx context.Context, id string, token string, url string, midtransOrderID string) error {
	query := `UPDATE payments SET midtrans_token = ?, midtrans_url = ?,
			midtrans_order_id = ?, updated_at = NOW() WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, token, url, midtransOrderID, id)

	return err
}

// SaveRawNotification saves the raw webhook notification payload
func (r *paymentRepository) SaveRawNotification(ctx context.Context, id string, raw string) error {
	query := `UPDATE payments SET raw_notification = ?, updated_at = NOW() WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, raw, id)

	return err
}
