package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/xprasetio/coffee-pos/internal/entity"
	"github.com/google/uuid"
)

// CategoryRepository defines the interface for category data access
type CategoryRepository interface {
	FindByID(ctx context.Context, id string) (*entity.Category, error)
	FindByName(ctx context.Context, name string) (*entity.Category, error)
	FindAll(ctx context.Context) ([]entity.Category, error)
	Create(ctx context.Context, category *entity.Category) error
	Update(ctx context.Context, category *entity.Category) error
	Delete(ctx context.Context, id string) error
	WithTx(tx *sql.Tx) CategoryRepository
}

// categoryRepository implements CategoryRepository
type categoryRepository struct {
	db sqlDB
}

// NewCategoryRepository creates a new CategoryRepository instance
func NewCategoryRepository(db *sql.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

// WithTx returns a new CategoryRepository instance with the given transaction
func (r *categoryRepository) WithTx(tx *sql.Tx) CategoryRepository {
	return &categoryRepository{db: tx}
}

// FindByID finds a category by ID
func (r *categoryRepository) FindByID(ctx context.Context, id string) (*entity.Category, error) {
	query := `SELECT id, name, created_at, updated_at, deleted_at
		FROM categories WHERE id = ? AND deleted_at IS NULL`

	category := &entity.Category{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&category.ID,
		&category.Name,
		&category.CreatedAt,
		&category.UpdatedAt,
		&category.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return category, nil
}

// FindByName finds a category by name
func (r *categoryRepository) FindByName(ctx context.Context, name string) (*entity.Category, error) {
	query := `SELECT id, name, created_at, updated_at, deleted_at
		FROM categories WHERE name = ? AND deleted_at IS NULL`

	category := &entity.Category{}
	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&category.ID,
		&category.Name,
		&category.CreatedAt,
		&category.UpdatedAt,
		&category.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return category, nil
}

// FindAll returns all categories ordered by name ascending
func (r *categoryRepository) FindAll(ctx context.Context) ([]entity.Category, error) {
	query := `SELECT id, name, created_at, updated_at, deleted_at
		FROM categories WHERE deleted_at IS NULL ORDER BY name ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]entity.Category, 0)
	for rows.Next() {
		var category entity.Category
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.CreatedAt,
			&category.UpdatedAt,
			&category.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

// Create creates a new category
func (r *categoryRepository) Create(ctx context.Context, category *entity.Category) error {
	category.ID = uuid.New().String()
	now := time.Now()
	category.CreatedAt = now
	category.UpdatedAt = now

	query := `INSERT INTO categories (id, name, created_at, updated_at) VALUES (?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		category.ID,
		category.Name,
		category.CreatedAt,
		category.UpdatedAt,
	)

	return err
}

// Update updates an existing category (only Name field)
func (r *categoryRepository) Update(ctx context.Context, category *entity.Category) error {
	category.UpdatedAt = time.Now()

	query := `UPDATE categories SET name=?, updated_at=?
		WHERE id=? AND deleted_at IS NULL`

	_, err := r.db.ExecContext(ctx, query,
		category.Name,
		category.UpdatedAt,
		category.ID,
	)

	return err
}

// Delete performs a soft delete on a category
func (r *categoryRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE categories SET deleted_at = NOW() WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, id)

	return err
}
