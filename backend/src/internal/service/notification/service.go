package notification

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/wizenheimer/byrd/src/internal/alert"
	"github.com/wizenheimer/byrd/src/internal/event"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
)

const (
	DefaultBufferSize int           = 25
	MinimumPriority   int           = 0
	MaximumPriority   int           = 10
	DefaultPriority   int           = MinimumPriority
	ProcessingTimeout time.Duration = 5 * time.Second
)

var (
	ErrServiceBusy       = errors.New("notification service is busy")
	ErrInvalidPriority   = errors.New("invalid priority level")
	ErrInvalidBufferSize = errors.New("invalid buffer size")
)

type NotificationService interface {
	// GetAlertChannel returns a channel for alerts or error if service is busy/timeout
	GetAlertChannel(ctx context.Context, priority int, bufferSize int) (<-chan models.Alert, error)

	// GetEventChannel returns a channel for events or error if service is busy/timeout
	GetEventChannel(ctx context.Context, priority int, bufferSize int) (<-chan models.Event, error)

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

type notificationService struct {
	alertChannels []chan models.Alert
	eventChannels []chan models.Event
	logChannel    chan any
	alertClient   alert.AlertClient
	eventClient   event.EventClient
	logger        *logger.Logger
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
	mu            sync.RWMutex

	// Fan-in channels
	alertFanIn chan models.Alert
	eventFanIn chan models.Event
}

func NewNotificationService(alertClient alert.AlertClient, eventClient event.EventClient, logger *logger.Logger) NotificationService {
	ctx, cancel := context.WithCancel(context.Background())
	return &notificationService{
		alertClient:   alertClient,
		eventClient:   eventClient,
		alertChannels: make([]chan models.Alert, 0),
		eventChannels: make([]chan models.Event, 0),
		logChannel:    make(chan any, DefaultBufferSize),
		logger:        logger.WithFields(map[string]interface{}{"module": "notification_service"}),
		ctx:           ctx,
		cancel:        cancel,
		alertFanIn:    make(chan models.Alert, DefaultBufferSize),
		eventFanIn:    make(chan models.Event, DefaultBufferSize),
	}
}

// GetAlertChannel returns a channel for alerts or error if service is busy/timeout
func (s *notificationService) GetAlertChannel(ctx context.Context, priority int, bufferSize int) (<-chan models.Alert, error) {
	if bufferSize <= 0 {
		return nil, ErrInvalidBufferSize
	}
	if priority < MinimumPriority || priority > MaximumPriority {
		return nil, ErrInvalidPriority
	}

	// Try to acquire lock with timeout
	// Timeout is there to prevent deadlocks
	lockChan := make(chan struct{})
	go func() {
		s.mu.Lock()
		close(lockChan)
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-lockChan:
		// Lock acquired, defer unlock
		defer s.mu.Unlock()

		ch := make(chan models.Alert, bufferSize)
		s.alertChannels = append(s.alertChannels, ch)
		go s.forwardAlerts(ch)
		return ch, nil
	case <-time.After(ProcessingTimeout):
		return nil, ErrServiceBusy
	}
}

func (s *notificationService) forwardAlerts(ch <-chan models.Alert) {
	for {
		select {
		case <-s.ctx.Done():
			return
		case alert, ok := <-ch:
			if !ok {
				return
			}
			select {
			case s.alertFanIn <- alert:
			case <-time.After(ProcessingTimeout):
				s.logger.Error("Alert forwarding timed out")
			}
		}
	}
}

func (s *notificationService) GetEventChannel(ctx context.Context, priority int, bufferSize int) (<-chan models.Event, error) {
	if bufferSize <= 0 {
		return nil, ErrInvalidBufferSize
	}
	if priority < MinimumPriority || priority > MaximumPriority {
		return nil, ErrInvalidPriority
	}

	// Try to acquire lock with timeout
	// Timeout is there to prevent deadlocks
	lockChan := make(chan struct{})
	go func() {
		s.mu.Lock()
		close(lockChan)
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-lockChan:
		// Lock acquired, defer unlock
		defer s.mu.Unlock()

		ch := make(chan models.Event, bufferSize)
		s.eventChannels = append(s.eventChannels, ch)
		go s.forwardEvents(ch)
		return ch, nil
	case <-time.After(ProcessingTimeout):
		return nil, ErrServiceBusy
	}
}

func (s *notificationService) forwardEvents(ch <-chan models.Event) {
	for {
		select {
		case <-s.ctx.Done():
			return
		case event, ok := <-ch:
			if !ok {
				return
			}
			select {
			case s.eventFanIn <- event:
			case <-time.After(ProcessingTimeout):
				s.logger.Error("Event forwarding timed out")
			}
		}
	}
}

func (s *notificationService) Start() {
	s.logger.Info("Starting notification service")
	s.wg.Add(3)
	go s.startAlertWorker()
	go s.startEventWorker()
	go s.startLogWorker()
}

func (s *notificationService) Stop() error {
	s.logger.Info("Stopping notification service")
	s.cancel()
	s.wg.Wait()
	return nil
}

func (s *notificationService) Close() error {
	s.logger.Info("Closing notification service")

	s.mu.Lock()
	defer s.mu.Unlock()

	// Close fan-in channels
	close(s.alertFanIn)
	close(s.eventFanIn)
  close(s.logChannel)

	// Close all alert channels
	for _, ch := range s.alertChannels {
		close(ch)
	}
	s.alertChannels = nil

	// Close all event channels
	for _, ch := range s.eventChannels {
		close(ch)
	}
	s.eventChannels = nil

	return nil
}

func (s *notificationService) startAlertWorker() {
	defer s.wg.Done()

	for {
		select {
		case <-s.ctx.Done():
			s.logger.Info("Alert worker stopping")
			return
		case alert, ok := <-s.alertFanIn:
			if !ok {
				return
			}
			if err := s.alertClient.SendAlert(s.ctx, alert); err != nil {
				s.logger.Error("Failed to send alert", zap.Error(err))
        s.logChannel <- alert
			}
		}
	}
}

func (s *notificationService) startEventWorker() {
	defer s.wg.Done()

	for {
		select {
		case <-s.ctx.Done():
			s.logger.Info("Event worker stopping")
			return
		case event, ok := <-s.eventFanIn:
			if !ok {
				return
			}
			if err := s.eventClient.SendEvent(s.ctx, event); err != nil {
				s.logger.Error("Failed to send event", zap.Error(err))
        s.logChannel <- event
			}
		}
	}
}

// Simple getter for the shared log channel
func (s *notificationService) GetLogChannel() chan<- any {
	return s.logChannel
}

func (s *notificationService) startLogWorker() {
	defer s.wg.Done()

	for {
		select {
		case <-s.ctx.Done():
			s.logger.Info("Log worker stopping")
			return
		case msg := <-s.logChannel:
			// Log the message with its type
			s.logger.Info("Log message",
				zap.Any("value", msg),
			)
		}
	}
}
