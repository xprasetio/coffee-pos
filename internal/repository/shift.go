package repository

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/xprasetio/coffee-pos/internal/entity"
	"github.com/google/uuid"
)

// ShiftRepository defines the interface for shift data access
type ShiftRepository interface {
	FindByID(ctx context.Context, id string) (*entity.Shift, error)
	FindOpenByCashierID(ctx context.Context, cashierID string) (*entity.Shift, error)
	FindAll(ctx context.Context, cashierID string, page, limit int) ([]entity.Shift, int, error)
	Create(ctx context.Context, shift *entity.Shift) error
	Close(ctx context.Context, id string, closingCash int64, notes string) error
	WithTx(tx *sql.Tx) ShiftRepository
}

// shiftRepository implements ShiftRepository
type shiftRepository struct {
	db sqlDB
}

// NewShiftRepository creates a new ShiftRepository instance
func NewShiftRepository(db *sql.DB) ShiftRepository {
	return &shiftRepository{db: db}
}

// WithTx returns a new ShiftRepository instance with the given transaction
func (r *shiftRepository) WithTx(tx *sql.Tx) ShiftRepository {
	return &shiftRepository{db: tx}
}

// FindByID finds a shift by ID with cashier joined
func (r *shiftRepository) FindByID(ctx context.Context, id string) (*entity.Shift, error) {
	query := `SELECT s.id, s.cashier_id, s.opened_at, s.closed_at, s.opening_cash,
			s.closing_cash, s.total_transactions, s.status, s.notes,
			s.created_at, s.updated_at,
			u.id, u.name
		FROM shifts s
		LEFT JOIN users u ON u.id = s.cashier_id AND u.deleted_at IS NULL
		WHERE s.id = ? AND s.deleted_at IS NULL`

	shift := &entity.Shift{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&shift.ID,
		&shift.CashierID,
		&shift.OpenedAt,
		&shift.ClosedAt,
		&shift.OpeningCash,
		&shift.ClosingCash,
		&shift.TotalTransactions,
		&shift.Status,
		&shift.Notes,
		&shift.CreatedAt,
		&shift.UpdatedAt,
		&shift.Cashier.ID,
		&shift.Cashier.Name,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return shift, nil
}

// FindOpenByCashierID finds an open shift for a specific cashier
func (r *shiftRepository) FindOpenByCashierID(ctx context.Context, cashierID string) (*entity.Shift, error) {
	query := `SELECT s.id, s.cashier_id, s.opened_at, s.closed_at, s.opening_cash,
			s.closing_cash, s.total_transactions, s.status, s.notes,
			s.created_at, s.updated_at,
			u.id, u.name
		FROM shifts s
		LEFT JOIN users u ON u.id = s.cashier_id AND u.deleted_at IS NULL
		WHERE s.cashier_id = ? AND s.status = ? AND s.deleted_at IS NULL`

	shift := &entity.Shift{}

	err := r.db.QueryRowContext(ctx, query, cashierID, entity.ShiftStatusOpen).Scan(
		&shift.ID,
		&shift.CashierID,
		&shift.OpenedAt,
		&shift.ClosedAt,
		&shift.OpeningCash,
		&shift.ClosingCash,
		&shift.TotalTransactions,
		&shift.Status,
		&shift.Notes,
		&shift.CreatedAt,
		&shift.UpdatedAt,
		&shift.Cashier.ID,
		&shift.Cashier.Name,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return shift, nil
}

// FindAll returns all shifts with optional cashier filter and pagination
func (r *shiftRepository) FindAll(ctx context.Context, cashierID string, page, limit int) ([]entity.Shift, int, error) {
	// Set default pagination values
	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	// Build WHERE conditions
	whereConditions := []string{"s.deleted_at IS NULL"}
	args := make([]interface{}, 0)

	if cashierID != "" {
		whereConditions = append(whereConditions, "s.cashier_id = ?")
		args = append(args, cashierID)
	}

	whereClause := strings.Join(whereConditions, " AND ")

	// Base query with JOIN
	baseQuery := `FROM shifts s
		LEFT JOIN users u ON u.id = s.cashier_id AND u.deleted_at IS NULL
		WHERE ` + whereClause

	// Count query for total
	countQuery := `SELECT COUNT(*) ` + baseQuery
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Data query
	dataQuery := `SELECT s.id, s.cashier_id, s.opened_at, s.closed_at, s.opening_cash,
			s.closing_cash, s.total_transactions, s.status, s.notes,
			s.created_at, s.updated_at,
			u.id, u.name
		` + baseQuery + ` ORDER BY s.opened_at DESC LIMIT ? OFFSET ?`

	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	shifts := make([]entity.Shift, 0)
	for rows.Next() {
		var shift entity.Shift

		err := rows.Scan(
			&shift.ID,
			&shift.CashierID,
			&shift.OpenedAt,
			&shift.ClosedAt,
			&shift.OpeningCash,
			&shift.ClosingCash,
			&shift.TotalTransactions,
			&shift.Status,
			&shift.Notes,
			&shift.CreatedAt,
			&shift.UpdatedAt,
			&shift.Cashier.ID,
			&shift.Cashier.Name,
		)
		if err != nil {
			return nil, 0, err
		}

		shifts = append(shifts, shift)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return shifts, total, nil
}

// Create creates a new shift
func (r *shiftRepository) Create(ctx context.Context, shift *entity.Shift) error {
	shift.ID = uuid.New().String()
	now := time.Now()
	shift.CreatedAt = now
	shift.UpdatedAt = now
	shift.OpenedAt = now
	shift.Status = entity.ShiftStatusOpen

	query := `INSERT INTO shifts (id, cashier_id, opening_cash, status, opened_at,
			created_at, updated_at, total_transactions, notes)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		shift.ID,
		shift.CashierID,
		shift.OpeningCash,
		shift.Status,
		shift.OpenedAt,
		shift.CreatedAt,
		shift.UpdatedAt,
		shift.TotalTransactions,
		shift.Notes,
	)

	return err
}

// Close closes a shift
func (r *shiftRepository) Close(ctx context.Context, id string, closingCash int64, notes string) error {
	query := `UPDATE shifts SET 
			status = ?,
			closed_at = NOW(),
			closing_cash = ?,
			notes = ?,
			updated_at = NOW()
		WHERE id = ? AND status = ?`

	_, err := r.db.ExecContext(ctx, query,
		entity.ShiftStatusClosed,
		closingCash,
		notes,
		id,
		entity.ShiftStatusOpen,
	)

	return err
}
