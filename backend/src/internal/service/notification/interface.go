// ./src/internal/service/notification/interface.go
package notification

import (
	"context"
	"errors"
	"time"

	models "github.com/wizenheimer/byrd/src/internal/models/core"
)

const (
	DefaultBufferSize int           = 25
	MinimumPriority   int           = 0
	MaximumPriority   int           = 10
	DefaultPriority   int           = MinimumPriority
	ProcessingTimeout time.Duration = 5 * time.Second
	MaxChannels                     = 100
)

var (
	ErrServiceBusy       = errors.New("notification service is busy")
	ErrInvalidPriority   = errors.New("invalid priority level")
	ErrInvalidBufferSize = errors.New("invalid buffer size")
	ErrTooManyChannels   = errors.New("maximum number of channels reached")
)

type NotificationService interface {
	// GetAlertChannel returns a channel for alerts or error if service is busy/timeout
	GetAlertChannel(ctx context.Context, priority int, bufferSize int) (<-chan models.Alert, error)

	// GetEventChannel returns a channel for events or error if service is busy/timeout
	GetEventChannel(ctx context.Context, priority int, bufferSize int) (<-chan models.Event, error)

	// GetEmailChannel returns a channel for emails or error if service is busy/timeout
	GetEmailChannel(ctx context.Context, priority int, bufferSize int) (<-chan models.Email, error)

	// Returns the shared log channel for all writers
	// Log channel is used as a fallback when the service is busy
	GetLogChannel() chan<- any

	// Start starts the notification service
	Start()

	// Gracefully stops the notification service
	Stop() error

	// Close closes the notification service channels
	Close() error
}
