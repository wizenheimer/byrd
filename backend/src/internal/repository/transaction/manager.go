// ./src/internal/repository/transaction/manager.go
package transaction

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

// Common errors
var (
	ErrTxNotFound      = errors.New("transaction not found in context")
	ErrTxAlreadyExists = errors.New("transaction already exists in context")
	ErrInvalidTxOpts   = errors.New("invalid transaction options")
)

// Key for storing transaction in context
type txKey struct{}

// TxOptions wraps sql.TxOptions with additional configuration
type TxOptions struct {
	// Isolation level for the transaction
	Isolation sql.IsolationLevel
	// ReadOnly transaction
	ReadOnly bool
	// RetryOnSerializationFailure determines if the transaction should be retried on serialization failures
	RetryOnSerializationFailure bool
	// MaxRetries is the maximum number of retries for serialization failures
	MaxRetries int
}

// DefaultTxOptions provides sensible defaults for transaction options
var DefaultTxOptions = TxOptions{
	Isolation:                   sql.LevelSerializable,
	ReadOnly:                    false,
	RetryOnSerializationFailure: true,
	MaxRetries:                  3,
}

// TxManager handles database transactions
type TxManager struct {
	db *sql.DB
}

// NewTxManager creates a new transaction manager
func NewTxManager(db *sql.DB) *TxManager {
	return &TxManager{db: db}
}

// GetTx retrieves an existing transaction from context
// Returns the transaction and true if found, nil and false otherwise
func GetTx(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(*sql.Tx)
	return tx, ok
}

// MustGetTx retrieves an existing transaction from context
// Panics if no transaction is found
func MustGetTx(ctx context.Context) *sql.Tx {
	tx, ok := GetTx(ctx)
	if !ok {
		panic(ErrTxNotFound)
	}
	return tx
}

// GetDB returns the underlying database connection
func (tm *TxManager) GetDB() *sql.DB {
	return tm.db
}

// GetRunner returns either the transaction from context or falls back to db
func (tm *TxManager) GetRunner(ctx context.Context) interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
} {
	if tx, ok := GetTx(ctx); ok {
		return tx
	}
	return tm.db
}

// RunInTx runs the given function in a transaction
// If a transaction already exists in the context, it will be reused
func (tm *TxManager) RunInTx(ctx context.Context, opts *TxOptions, fn func(context.Context) error) error {
	// Use default options if none provided
	if opts == nil {
		opts = &DefaultTxOptions
	}

	// Check if we already have a transaction
	if _, exists := GetTx(ctx); exists {
		// Reuse existing transaction
		return fn(ctx)
	}

	// Initialize retry counter
	retries := 0

	for {
		err := tm.runSingleTx(ctx, opts, fn)
		if err == nil {
			return nil
		}

		// Check if we should retry on serialization failure
		if opts.RetryOnSerializationFailure && isSerializationFailure(err) && retries < opts.MaxRetries {
			retries++
			continue
		}

		return err
	}
}

// runSingleTx executes a single transaction attempt
func (tm *TxManager) runSingleTx(ctx context.Context, opts *TxOptions, fn func(context.Context) error) error {
	// Start new transaction
	tx, err := tm.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: opts.Isolation,
		ReadOnly:  opts.ReadOnly,
	})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Add transaction to context
	ctxWithTx := context.WithValue(ctx, txKey{}, tx)

	// Handle panic and rollback
	defer func() {
		if p := recover(); p != nil {
			err := tx.Rollback()
			if err != nil {
				panic(fmt.Errorf("error rolling back transaction: %v (original panic: %v)", err, p))
			}
			panic(p) // re-throw panic after rollback
		}
	}()

	// Run the function
	err = fn(ctxWithTx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("error rolling back transaction: %v (original error: %w)", rbErr, err)
		}
		return err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// RunInTxReadOnly is a convenience method for running read-only transactions
func (tm *TxManager) RunInTxReadOnly(ctx context.Context, fn func(context.Context) error) error {
	opts := DefaultTxOptions
	opts.ReadOnly = true
	return tm.RunInTx(ctx, &opts, fn)
}

// isSerializationFailure checks if the error is a serialization failure that should be retried
func isSerializationFailure(err error) bool {
	// PostgreSQL serialization failure error code: 40001
	// You might need to add more error codes or adapt this for different databases
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "40001"
	}
	return false
}

// WithTx creates a new transaction and adds it to the context
// This is useful when you need more control over the transaction lifecycle
func (tm *TxManager) WithTx(ctx context.Context, opts *TxOptions) (context.Context, *sql.Tx, error) {
	if opts == nil {
		opts = &DefaultTxOptions
	}

	// Check if transaction already exists
	if _, exists := GetTx(ctx); exists {
		return nil, nil, ErrTxAlreadyExists
	}

	// Start new transaction
	tx, err := tm.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: opts.Isolation,
		ReadOnly:  opts.ReadOnly,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Add transaction to context
	ctxWithTx := context.WithValue(ctx, txKey{}, tx)
	return ctxWithTx, tx, nil
}
