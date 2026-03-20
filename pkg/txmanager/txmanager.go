package txmanager

import (
	"context"
	"database/sql"
)

// TxManager manages database transactions
type TxManager struct {
	db *sql.DB
}

// New creates a new TxManager instance
func New(db *sql.DB) *TxManager {
	return &TxManager{db: db}
}

// WithTx executes a function within a database transaction
// If the function returns an error, the transaction is rolled back
// If the function succeeds, the transaction is committed
func (tm *TxManager) WithTx(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := tm.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit()
}
