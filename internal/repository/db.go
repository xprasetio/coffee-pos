package repository

import (
	"context"
	"database/sql"
)

// sqlDB defines the interface for database operations
// This allows repositories to work with both *sql.DB and *sql.Tx
type sqlDB interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}
