package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/xprasetio/coffee-pos/internal/entity"
	"github.com/google/uuid"
)

// TableRepository defines the interface for table data access
type TableRepository interface {
	FindByID(ctx context.Context, id string) (*entity.Table, error)
	FindByName(ctx context.Context, name string) (*entity.Table, error)
	FindAll(ctx context.Context) ([]entity.Table, error)
	Create(ctx context.Context, table *entity.Table) error
	Update(ctx context.Context, table *entity.Table) error
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status string) error
	WithTx(tx *sql.Tx) TableRepository
}

// tableRepository implements TableRepository
type tableRepository struct {
	db sqlDB
}

// NewTableRepository creates a new TableRepository instance
func NewTableRepository(db *sql.DB) TableRepository {
	return &tableRepository{db: db}
}

// WithTx returns a new TableRepository instance with the given transaction
func (r *tableRepository) WithTx(tx *sql.Tx) TableRepository {
	return &tableRepository{db: tx}
}

// FindByID finds a table by ID
func (r *tableRepository) FindByID(ctx context.Context, id string) (*entity.Table, error) {
	query := `SELECT id, name, capacity, status, created_at, updated_at, deleted_at
		FROM tables WHERE id = ? AND deleted_at IS NULL`

	table := &entity.Table{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&table.ID,
		&table.Name,
		&table.Capacity,
		&table.Status,
		&table.CreatedAt,
		&table.UpdatedAt,
		&table.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return table, nil
}

// FindByName finds a table by name
func (r *tableRepository) FindByName(ctx context.Context, name string) (*entity.Table, error) {
	query := `SELECT id, name, capacity, status, created_at, updated_at, deleted_at
		FROM tables WHERE name = ? AND deleted_at IS NULL`

	table := &entity.Table{}

	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&table.ID,
		&table.Name,
		&table.Capacity,
		&table.Status,
		&table.CreatedAt,
		&table.UpdatedAt,
		&table.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return table, nil
}

// FindAll returns all tables ordered by name ascending
func (r *tableRepository) FindAll(ctx context.Context) ([]entity.Table, error) {
	query := `SELECT id, name, capacity, status, created_at, updated_at, deleted_at
		FROM tables WHERE deleted_at IS NULL ORDER BY name ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tables := make([]entity.Table, 0)
	for rows.Next() {
		var table entity.Table
		err := rows.Scan(
			&table.ID,
			&table.Name,
			&table.Capacity,
			&table.Status,
			&table.CreatedAt,
			&table.UpdatedAt,
			&table.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tables, nil
}

// Create creates a new table
func (r *tableRepository) Create(ctx context.Context, table *entity.Table) error {
	table.ID = uuid.New().String()
	now := time.Now()
	table.CreatedAt = now
	table.UpdatedAt = now

	query := `INSERT INTO tables (id, name, capacity, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		table.ID,
		table.Name,
		table.Capacity,
		table.Status,
		table.CreatedAt,
		table.UpdatedAt,
	)

	return err
}

// Update updates an existing table
func (r *tableRepository) Update(ctx context.Context, table *entity.Table) error {
	table.UpdatedAt = time.Now()

	query := `UPDATE tables SET name = ?, capacity = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL`

	_, err := r.db.ExecContext(ctx, query,
		table.Name,
		table.Capacity,
		table.UpdatedAt,
		table.ID,
	)

	return err
}

// Delete performs a soft delete on a table
func (r *tableRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE tables SET deleted_at = NOW() WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, id)

	return err
}

// UpdateStatus updates the status of a table
func (r *tableRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	query := `UPDATE tables SET status = ?, updated_at = ? WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, status, time.Now(), id)

	return err
}
