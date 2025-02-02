// ./src/internal/repository/transaction/manager.go
package transaction

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

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
	ErrShuttingDown    = errors.New("transaction manager is shutting down")
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

var ReadOnlyTxOptions = TxOptions{
	IsoLevel:                    pgx.Serializable,
	Access:                      pgx.ReadOnly,
	RetryOnSerializationFailure: true,
	MaxRetries:                  3,
}

// TxManager handles database transactions
type TxManager struct {
	pool         *pgxpool.Pool
	logger       *logger.Logger
	activeCount  int64
	activeMutex  sync.RWMutex
	isShutdown   bool
	shutdownOnce sync.Once
}

// NewTxManager creates a new transaction manager
func NewTxManager(pool *pgxpool.Pool, logger *logger.Logger) *TxManager {
	tm := &TxManager{
		pool: pool,
		logger: logger.WithFields(map[string]any{
			"module": "transaction_manager",
		}),
	}
	return tm
}

// GetTx retrieves an existing transaction from context
// Returns the transaction and true if found, nil and false otherwise
func GetTx(ctx context.Context) (pgx.Tx, bool) {
	if ctx == nil {
		return nil, false
	}
	tx, ok := ctx.Value(txKey{}).(pgx.Tx)
	return tx, ok
}

// MustGetTx retrieves an existing transaction from context
// Panics if no transaction is found
func MustGetTx(ctx context.Context) (pgx.Tx, error) {
	tx, ok := GetTx(ctx)
	if !ok {
		return nil, ErrTxNotFound
	}
	return tx, nil
}

// GetPool returns the underlying connection pool
func (tm *TxManager) GetPool() *pgxpool.Pool {
	return tm.pool
}

// GetQuerier returns either the transaction from context or falls back to pool
func (tm *TxManager) GetQuerier(ctx context.Context) interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
} {
	if tx, ok := GetTx(ctx); ok {
		return tx
	}
	return tm.pool
}

// trackTransaction safely increments or decrements the active transaction count
func (tm *TxManager) trackTransaction(increment bool) {
	tm.activeMutex.Lock()
	defer tm.activeMutex.Unlock()

	if increment {
		tm.activeCount++
	} else {
		if tm.activeCount > 0 { // Prevent going negative
			tm.activeCount--
		} else {
			tm.logger.Warn("attempted to decrement transaction count below zero")
		}
	}
}

// GetActiveTransactionCount returns the current number of active transactions
func (tm *TxManager) GetActiveTransactionCount() int64 {
	tm.activeMutex.RLock()
	defer tm.activeMutex.RUnlock()
	return tm.activeCount
}

// isShuttingDown checks if the manager is in shutdown mode
func (tm *TxManager) isShuttingDown() bool {
	tm.activeMutex.RLock()
	defer tm.activeMutex.RUnlock()
	return tm.isShutdown
}

// Shutdown gracefully shuts down the transaction manager
func (tm *TxManager) Shutdown(ctx context.Context) error {
	var shutdownErr error

	tm.shutdownOnce.Do(func() {
		tm.activeMutex.Lock()
		tm.isShutdown = true
		tm.activeMutex.Unlock()

		// Log shutdown initiation
		tm.logger.Info("initiating transaction manager shutdown",
			zap.Int64("active_transactions", tm.activeCount))

		// Wait for active transactions to complete with timeout
		deadline := time.After(120 * time.Second)
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				shutdownErr = fmt.Errorf("shutdown context cancelled: %w", ctx.Err())
				return
			case <-deadline:
				// Close the connection pool forcefully, to avoid deadlock
				tm.pool.Close()
				shutdownErr = fmt.Errorf("shutdown timed out with %d active transactions", tm.GetActiveTransactionCount())
				return
			case <-ticker.C:
				if tm.GetActiveTransactionCount() == 0 {
					// Close the connection pool
					tm.pool.Close()
					tm.logger.Info("transaction manager shutdown completed")
					return
				}
			}
		}
	})

	return shutdownErr
}

// RunInTx runs the given function in a transaction
func (tm *TxManager) RunInTx(ctx context.Context, opts *TxOptions, fn func(context.Context) error) error {
	if tm.isShuttingDown() {
		return ErrShuttingDown
	}

	if opts == nil {
		opts = &DefaultTxOptions
	}

	// Check if we already have a transaction
	if _, exists := GetTx(ctx); exists {
		return fn(ctx) // Reuse existing transaction
	}

	// Track this transaction
	tm.trackTransaction(true)
	defer tm.trackTransaction(false) // Ensure we always untrack, even on panic

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

// runSingleTx executes a single transaction attempt with improved panic handling
func (tm *TxManager) runSingleTx(ctx context.Context, opts *TxOptions, fn func(context.Context) error) error {
	tx, err := tm.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   opts.IsoLevel,
		AccessMode: opts.Access,
	})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	ctxWithTx := context.WithValue(ctx, txKey{}, tx)

	// Improved panic handling
	var panicked bool
	var panicValue interface{}

	defer func() {
		if p := recover(); p != nil {
			panicked = true
			panicValue = p

			// Attempt rollback but don't panic if it fails
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				tm.logger.Error("failed to rollback transaction after panic",
					zap.Error(rbErr),
					zap.Any("original_panic", p))
			}
		}
	}()

	// Run the function
	err = fn(ctxWithTx)

	// Handle normal execution path
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			tm.logger.Error("failed to rollback transaction",
				zap.Error(rbErr),
				zap.Error(err))
			return fmt.Errorf("error rolling back transaction: %v (original error: %w)", rbErr, err)
		}
		return err
	}

	// Commit the transaction
	if err := tx.Commit(ctx); err != nil {
		tm.logger.Error("failed to commit transaction", zap.Error(err))
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Re-panic if we recovered from a panic earlier
	if panicked {
		panic(panicValue)
	}

	return nil
}

// RunInTxReadOnly is a convenience method for running read-only transactions
func (tm *TxManager) RunInTxReadOnly(ctx context.Context, fn func(context.Context) error) error {
	opts := ReadOnlyTxOptions
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
func (tm *TxManager) WithTx(ctx context.Context, opts *TxOptions) (context.Context, pgx.Tx, error) {
	if tm.isShuttingDown() {
		return nil, nil, ErrShuttingDown
	}

	if opts == nil {
		opts = &DefaultTxOptions
	}

	// Check if transaction already exists
	if _, exists := GetTx(ctx); exists {
		tm.logger.Error("transaction already exists in context")
		return nil, nil, ErrTxAlreadyExists
	}

	// Track this transaction with deferred cleanup on error
	tm.trackTransaction(true)

	// Start new transaction
	tx, err := tm.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   opts.IsoLevel,
		AccessMode: opts.Access,
	})
	if err != nil {
		tm.trackTransaction(false) // Ensure count is decremented on error
		return nil, nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Add transaction to context
	ctxWithTx := context.WithValue(ctx, txKey{}, tx)
	return ctxWithTx, tx, nil
}
