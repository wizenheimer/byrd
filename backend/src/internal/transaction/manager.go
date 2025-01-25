// ./src/internal/transaction/manager.go
// ./src/internal/repository/transaction/manager.go
package transaction

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
)

// Common errors
var (
	ErrTxNotFound      = errors.New("transaction not found in context")
	ErrTxAlreadyExists = errors.New("transaction already exists in context")
	ErrInvalidTxOpts   = errors.New("invalid transaction options")
)

// Key for storing transaction in context
type txKey struct{}

// TxOptions wraps pgx.TxOptions with additional configuration
type TxOptions struct {
	// IsoLevel is the isolation level for the transaction
	IsoLevel pgx.TxIsoLevel
	// Access is the transaction access mode (ReadWrite or ReadOnly)
	Access pgx.TxAccessMode
	// RetryOnSerializationFailure determines if the transaction should be retried on serialization failures
	RetryOnSerializationFailure bool
	// MaxRetries is the maximum number of retries for serialization failures
	MaxRetries int
}

// DefaultTxOptions provides sensible defaults for transaction options
var DefaultTxOptions = TxOptions{
	IsoLevel:                    pgx.Serializable,
	Access:                      pgx.ReadWrite,
	RetryOnSerializationFailure: true,
	MaxRetries:                  3,
}

// TxManager handles database transactions
type TxManager struct {
	pool   *pgxpool.Pool
	logger *logger.Logger
}

// NewTxManager creates a new transaction manager
func NewTxManager(pool *pgxpool.Pool, logger *logger.Logger) *TxManager {
	tm := TxManager{
		pool:   pool,
		logger: logger.WithFields(map[string]interface{}{"module": "transaction_manager"}),
	}
	tm.logger.Debug("created a new transaction manager")
	return &tm
}

// GetTx retrieves an existing transaction from context
// Returns the transaction and true if found, nil and false otherwise
func GetTx(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(pgx.Tx)
	return tx, ok
}

// MustGetTx retrieves an existing transaction from context
// Panics if no transaction is found
func MustGetTx(ctx context.Context) pgx.Tx {
	tx, ok := GetTx(ctx)
	if !ok {
		panic(ErrTxNotFound)
	}
	return tx
}

// GetPool returns the underlying connection pool
func (tm *TxManager) GetPool() *pgxpool.Pool {
	tm.logger.Debug("getting pool for making query")
	return tm.pool
}

// GetQuerier returns either the transaction from context or falls back to pool
func (tm *TxManager) GetQuerier(ctx context.Context) interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
} {
	if tx, ok := GetTx(ctx); ok {
		tm.logger.Debug("transaction querier found in context, using it for making query")
		return tx
	}
	tm.logger.Debug("using pool for making query")
	return tm.pool
}

// RunInTx runs the given function in a transaction
// If a transaction already exists in the context, it will be reused
func (tm *TxManager) RunInTx(ctx context.Context, opts *TxOptions, fn func(context.Context) error) error {
	// Use default options if none provided
	if opts == nil {
		tm.logger.Debug("using default transaction options")
		opts = &DefaultTxOptions
	}

	// Check if we already have a transaction
	if _, exists := GetTx(ctx); exists {
		// Reuse existing transaction
		tm.logger.Debug("reusing existing transaction")
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
			tm.logger.Debug("retrying transaction due to serialization failure", zap.Error(err))
			retries++
			continue
		}

		tm.logger.Debug("failed to run transaction", zap.Error(err))
		return err
	}
}

// runSingleTx executes a single transaction attempt
func (tm *TxManager) runSingleTx(ctx context.Context, opts *TxOptions, fn func(context.Context) error) error {
	// Start new transaction
	tx, err := tm.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   opts.IsoLevel,
		AccessMode: opts.Access,
	})
	if err != nil {
		tm.logger.Debug("failed to begin transaction", zap.Error(err))
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Add transaction to context
	ctxWithTx := context.WithValue(ctx, txKey{}, tx)

	// Handle panic and rollback
	defer func() {
		if p := recover(); p != nil {
			tm.logger.Error("panic in transaction", zap.Any("error", p))
			err := tx.Rollback(ctx)
			if err != nil {
				tm.logger.Error("failed to rollback transaction", zap.Error(err))
				panic(fmt.Errorf("error rolling back transaction: %v (original panic: %v)", err, p))
			}
			panic(p) // re-throw panic after rollback
		}
	}()

	// Run the function
	err = fn(ctxWithTx)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			tm.logger.Error("failed to rollback transaction", zap.Error(rbErr))
			return fmt.Errorf("error rolling back transaction: %v (original error: %w)", rbErr, err)
		}
		return err
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		tm.logger.Error("failed to commit transaction", zap.Error(err))
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// RunInTxReadOnly is a convenience method for running read-only transactions
func (tm *TxManager) RunInTxReadOnly(ctx context.Context, fn func(context.Context) error) error {
	tm.logger.Debug("running read-only transaction")
	opts := DefaultTxOptions
	opts.Access = pgx.ReadOnly
	return tm.RunInTx(ctx, &opts, fn)
}

// isSerializationFailure checks if the error is a serialization failure that should be retried
func isSerializationFailure(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		// PostgreSQL serialization failure error code: 40001
		return pgErr.Code == "40001"
	}
	return false
}

// WithTx creates a new transaction and adds it to the context
// This is useful when you need more control over the transaction lifecycle
func (tm *TxManager) WithTx(ctx context.Context, opts *TxOptions) (context.Context, pgx.Tx, error) {
	if opts == nil {
		tm.logger.Debug("using default transaction options")
		opts = &DefaultTxOptions
	}

	// Check if transaction already exists
	if _, exists := GetTx(ctx); exists {
		tm.logger.Error("transaction already exists in context")
		return nil, nil, ErrTxAlreadyExists
	}

	// Start new transaction
	tx, err := tm.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   opts.IsoLevel,
		AccessMode: opts.Access,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Add transaction to context
	ctxWithTx := context.WithValue(ctx, txKey{}, tx)
	tm.logger.Debug("transaction created and added to context")
	return ctxWithTx, tx, nil
}
