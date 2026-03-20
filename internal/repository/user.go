package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/xprasetio/coffee-pos/internal/entity"
	"github.com/google/uuid"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	FindByID(ctx context.Context, id string) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	Create(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, user *entity.User) error
	FindAll(ctx context.Context) ([]entity.User, error)
	WithTx(tx *sql.Tx) UserRepository
}

// userRepository implements UserRepository
type userRepository struct {
	db sqlDB
}

// NewUserRepository creates a new UserRepository instance
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

// WithTx returns a new UserRepository instance with the given transaction
func (r *userRepository) WithTx(tx *sql.Tx) UserRepository {
	return &userRepository{db: tx}
}

// FindByID finds a user by ID
func (r *userRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	query := `SELECT id, name, email, password, role, is_active, created_at, updated_at, deleted_at
		FROM users WHERE id = ? AND deleted_at IS NULL`

	user := &entity.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

// FindByEmail finds a user by email
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `SELECT id, name, email, password, role, is_active, created_at, updated_at, deleted_at
		FROM users WHERE email = ? AND deleted_at IS NULL`

	user := &entity.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Create creates a new user
func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	user.ID = uuid.New().String()
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	query := `INSERT INTO users (id, name, email, password, role, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Name,
		user.Email,
		user.Password,
		user.Role,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
	)

	return err
}

// Update updates an existing user
func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	user.UpdatedAt = time.Now()

	query := `UPDATE users SET name=?, email=?, role=?, is_active=?, updated_at=?
		WHERE id=? AND deleted_at IS NULL`

	_, err := r.db.ExecContext(ctx, query,
		user.Name,
		user.Email,
		user.Role,
		user.IsActive,
		user.UpdatedAt,
		user.ID,
	)

	return err
}

// FindAll returns all users
func (r *userRepository) FindAll(ctx context.Context) ([]entity.User, error) {
	query := `SELECT id, name, email, password, role, is_active, created_at, updated_at, deleted_at
		FROM users WHERE deleted_at IS NULL ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]entity.User, 0)
	for rows.Next() {
		var user entity.User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.Role,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
