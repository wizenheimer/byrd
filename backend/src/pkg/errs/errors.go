// ./src/pkg/errs/errors.go
package errs

import (
	"encoding/json"

	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
)

// ErrorContext is an alias for map[string]any to store error parameters
type ErrorContext map[string]any

// Error is a map that stores errors with their list of contexts
type Error map[error][]ErrorContext

// MarshalJSON implements json.Marshaler
func (e Error) MarshalJSON() ([]byte, error) {
	// Create a map with string keys for JSON marshaling
	m := make(map[string][]ErrorContext)
	for k, v := range e {
		m[k.Error()] = v
	}
	return json.Marshal(m)
}

// New creates a new Error instance
func New() Error {
	return make(Error)
}

// Add appends a new context to an error
func (e Error) Add(err error, ctx map[string]any) {
	errCtx := ErrorContext(ctx)
	e[err] = append(e[err], errCtx)
}

// AddContext appends a new context to an error
func (e Error) AddContexts(err error, ctx []ErrorContext) {
	e[err] = append(e[err], ctx...)
}

// Contains checks if an error exists in the Error map
func (e Error) Contains(err error) bool {
	_, exists := e[err]
	return exists
}

// HasFatalErrors checks if there are any errors other than the provided non-fatal errors
func (e Error) HasFatalErrors(nonFatalErrors ...error) bool {
	nonFatal := make(map[error]bool)
	for _, err := range nonFatalErrors {
		nonFatal[err] = true
	}

	for err := range e {
		if !nonFatal[err] {
			return true
		}
	}
	return false
}

// HasErrors checks if there are any errors
func (e Error) HasErrors() bool {
	return len(e) > 0
}

func (e Error) Log(msg string, logger *logger.Logger, fields ...zap.Field) {
	if len(e) == 0 {
		return
	}

	fields = append(fields, zap.Any("errors", e))

	logger.Error(msg, fields...)
}

func (e Error) Merge(others ...Error) {
	for _, other := range others {
		for err, contexts := range other {
			e.AddContexts(err, contexts)
		}
	}
}

// FilterFatalErrors returns a new Error instance containing only fatal errors
func (e Error) FilterFatalErrors(nonFatalErrors ...error) Error {
	nonFatal := make(map[error]bool)
	for _, err := range nonFatalErrors {
		nonFatal[err] = true
	}

	fatalErrors := New()
	for err, contexts := range e {
		if !nonFatal[err] {
			fatalErrors[err] = contexts
		}
	}
	return fatalErrors
}

// Propagate consolidates errors to a higher level with a new error type
func (e Error) Propagate(levelError error, nonFatalErrors ...error) Error {
	// Filter out non-fatal errors
	fatalErrors := e.FilterFatalErrors(nonFatalErrors...)

	if len(fatalErrors) == 0 {
		return nil
	}

	// Create new error store for the higher level
	higherLevel := New()

	// Convert errors to array of ErrorDetail
	errContexts := []ErrorContext{}

	// Collect all errors and their contexts
	for err, contexts := range fatalErrors {
		for _, ctx := range contexts {
			errContexts = append(errContexts, ErrorContext{
				"error_type":   err.Error(),
				"error_detail": ctx,
			})
		}
	}

	// Add consolidated error to higher level with the array structure
	higherLevel.AddContexts(levelError, errContexts)

	return higherLevel
}

// Error implements the error interface and returns a string representation of all errors
func (e Error) Error() string {
	if len(e) == 0 {
		return "no errors"
	}

	// Convert to a map[string][]ErrorContext for consistent JSON marshaling
	m := make(map[string][]ErrorContext)
	for k, v := range e {
		m[k.Error()] = v
	}

	// Marshal to JSON
	bytes, err := json.MarshalIndent(m, "", "    ")
	if err != nil {
		return "error marshaling errors: " + err.Error()
	}

	return string(bytes)
}
