package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/xprasetio/coffee-pos/internal/entity"
	"github.com/google/uuid"
)

// PromoRepository defines the interface for promo data access
type PromoRepository interface {
	FindByID(ctx context.Context, id string) (*entity.Promo, error)
	FindByCode(ctx context.Context, code string) (*entity.Promo, error)
	FindAll(ctx context.Context, page, limit int) ([]entity.Promo, int, error)
	Create(ctx context.Context, promo *entity.Promo) error
	Update(ctx context.Context, promo *entity.Promo) error
	Delete(ctx context.Context, id string) error
	IncrementUsedCount(ctx context.Context, id string) error
	WithTx(tx *sql.Tx) PromoRepository
}

// promoRepository implements PromoRepository
type promoRepository struct {
	db sqlDB
}

// NewPromoRepository creates a new PromoRepository instance
func NewPromoRepository(db *sql.DB) PromoRepository {
	return &promoRepository{db: db}
}

// WithTx returns a new PromoRepository instance with the given transaction
func (r *promoRepository) WithTx(tx *sql.Tx) PromoRepository {
	return &promoRepository{db: tx}
}

// FindByID finds a promo by ID
func (r *promoRepository) FindByID(ctx context.Context, id string) (*entity.Promo, error) {
	query := `SELECT id, name, code, type, value, min_order, max_discount, usage_limit,
			used_count, started_at, ended_at, is_active, created_at, updated_at, deleted_at
		FROM promos WHERE id = ? AND deleted_at IS NULL`

	promo := &entity.Promo{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&promo.ID,
		&promo.Name,
		&promo.Code,
		&promo.Type,
		&promo.Value,
		&promo.MinOrder,
		&promo.MaxDiscount,
		&promo.UsageLimit,
		&promo.UsedCount,
		&promo.StartedAt,
		&promo.EndedAt,
		&promo.IsActive,
		&promo.CreatedAt,
		&promo.UpdatedAt,
		&promo.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return promo, nil
}

// FindByCode finds a promo by code
func (r *promoRepository) FindByCode(ctx context.Context, code string) (*entity.Promo, error) {
	query := `SELECT id, name, code, type, value, min_order, max_discount, usage_limit,
			used_count, started_at, ended_at, is_active, created_at, updated_at, deleted_at
		FROM promos WHERE code = ? AND deleted_at IS NULL`

	promo := &entity.Promo{}

	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&promo.ID,
		&promo.Name,
		&promo.Code,
		&promo.Type,
		&promo.Value,
		&promo.MinOrder,
		&promo.MaxDiscount,
		&promo.UsageLimit,
		&promo.UsedCount,
		&promo.StartedAt,
		&promo.EndedAt,
		&promo.IsActive,
		&promo.CreatedAt,
		&promo.UpdatedAt,
		&promo.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return promo, nil
}

// FindAll returns all promos with pagination
func (r *promoRepository) FindAll(ctx context.Context, page, limit int) ([]entity.Promo, int, error) {
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

	// Count query for total
	countQuery := `SELECT COUNT(*) FROM promos WHERE deleted_at IS NULL`
	var total int
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Data query
	dataQuery := `SELECT id, name, code, type, value, min_order, max_discount, usage_limit,
			used_count, started_at, ended_at, is_active, created_at, updated_at, deleted_at
		FROM promos WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, dataQuery, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	promos := make([]entity.Promo, 0)
	for rows.Next() {
		var promo entity.Promo

		err := rows.Scan(
			&promo.ID,
			&promo.Name,
			&promo.Code,
			&promo.Type,
			&promo.Value,
			&promo.MinOrder,
			&promo.MaxDiscount,
			&promo.UsageLimit,
			&promo.UsedCount,
			&promo.StartedAt,
			&promo.EndedAt,
			&promo.IsActive,
			&promo.CreatedAt,
			&promo.UpdatedAt,
			&promo.DeletedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		promos = append(promos, promo)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return promos, total, nil
}

// Create creates a new promo
func (r *promoRepository) Create(ctx context.Context, promo *entity.Promo) error {
	promo.ID = uuid.New().String()
	now := time.Now()
	promo.CreatedAt = now
	promo.UpdatedAt = now

	query := `INSERT INTO promos (id, name, code, type, value, min_order, max_discount,
			usage_limit, used_count, started_at, ended_at, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		promo.ID,
		promo.Name,
		promo.Code,
		promo.Type,
		promo.Value,
		promo.MinOrder,
		promo.MaxDiscount,
		promo.UsageLimit,
		promo.UsedCount,
		promo.StartedAt,
		promo.EndedAt,
		promo.IsActive,
		promo.CreatedAt,
		promo.UpdatedAt,
	)

	return err
}

// Update updates an existing promo
func (r *promoRepository) Update(ctx context.Context, promo *entity.Promo) error {
	promo.UpdatedAt = time.Now()

	query := `UPDATE promos SET name = ?, type = ?, value = ?, min_order = ?,
			max_discount = ?, usage_limit = ?, started_at = ?, ended_at = ?,
			is_active = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL`

	_, err := r.db.ExecContext(ctx, query,
		promo.Name,
		promo.Type,
		promo.Value,
		promo.MinOrder,
		promo.MaxDiscount,
		promo.UsageLimit,
		promo.StartedAt,
		promo.EndedAt,
		promo.IsActive,
		promo.UpdatedAt,
		promo.ID,
	)

	return err
}

// Delete performs a soft delete on a promo
func (r *promoRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE promos SET deleted_at = NOW() WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, id)

	return err
}

// IncrementUsedCount increments the used_count of a promo by 1
func (r *promoRepository) IncrementUsedCount(ctx context.Context, id string) error {
	query := `UPDATE promos SET used_count = used_count + 1, updated_at = NOW() WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, id)

	return err
}
