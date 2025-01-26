// ./src/internal/service/notification/service.go
package notification

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/wizenheimer/byrd/src/internal/alert"
	"github.com/wizenheimer/byrd/src/internal/email"
	"github.com/wizenheimer/byrd/src/internal/event"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
)

type notificationService struct {
	alertChannels []chan models.Alert
	alertClient   alert.AlertClient

	eventClient   event.EventClient
	eventChannels []chan models.Event

	emailClient   email.EmailClient
	emailChannels []chan models.Email

	logChannel chan any
	logger     *logger.Logger
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
	mu         sync.RWMutex

	// Fan-in channels
	alertFanIn       chan models.Alert
	eventFanIn       chan models.Event
	emailClientFanIn chan models.Email

	// Channel counters
	alertChannelCount atomic.Int32
	eventChannelCount atomic.Int32
	emailChannelCount atomic.Int32
}

func NewNotificationService(alertClient alert.AlertClient, eventClient event.EventClient, emailClient email.EmailClient, logger *logger.Logger) NotificationService {
	ctx, cancel := context.WithCancel(context.Background())

	service := notificationService{
		alertClient:   alertClient,
		alertChannels: make([]chan models.Alert, 0),

		eventClient:   eventClient,
		eventChannels: make([]chan models.Event, 0),

		emailClient:   emailClient,
		emailChannels: make([]chan models.Email, 0),

		logChannel: make(chan any, DefaultBufferSize),
		logger:     logger.WithFields(map[string]interface{}{"module": "notification_service"}),

		ctx:              ctx,
		cancel:           cancel,
		alertFanIn:       make(chan models.Alert, DefaultBufferSize),
		eventFanIn:       make(chan models.Event, DefaultBufferSize),
		emailClientFanIn: make(chan models.Email, DefaultBufferSize*2),
	}

	// Initialize counters
	service.alertChannelCount.Store(0)
	service.eventChannelCount.Store(0)
	service.emailChannelCount.Store(0)

	return &service
}

// GetAlertChannel returns a channel for alerts or error if service is busy/timeout
func (s *notificationService) GetAlertChannel(ctx context.Context, priority int, bufferSize int) (<-chan models.Alert, error) {
	if bufferSize <= 0 {
		return nil, ErrInvalidBufferSize
	}
	if priority < MinimumPriority || priority > MaximumPriority {
		return nil, ErrInvalidPriority
	}

	// Check channel count before acquiring lock
	if s.alertChannelCount.Load() >= MaxChannels {
		return nil, ErrTooManyChannels
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

		// Double check after acquiring lock
		if s.alertChannelCount.Load() >= MaxChannels {
			return nil, ErrTooManyChannels
		}

		ch := make(chan models.Alert, bufferSize)
		s.alertChannels = append(s.alertChannels, ch)
		s.alertChannelCount.Add(1)
		go s.forwardAlerts(ch)
		return ch, nil
	case <-time.After(ProcessingTimeout):
		return nil, ErrServiceBusy
	}
}

func (s *notificationService) forwardAlerts(ch <-chan models.Alert) {
	defer func() {
		s.mu.Lock()
		// Remove channel from slice
		for i, c := range s.alertChannels {
			if c == ch {
				s.alertChannels = append(s.alertChannels[:i], s.alertChannels[i+1:]...)
				break
			}
		}
		s.mu.Unlock()
	}()

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

	// Check channel count before acquiring lock
	if s.eventChannelCount.Load() >= MaxChannels {
		return nil, ErrTooManyChannels
	}

	lockChan := make(chan struct{})
	go func() {
		s.mu.Lock()
		close(lockChan)
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-lockChan:
		defer s.mu.Unlock()

		// Double check after acquiring lock
		if s.eventChannelCount.Load() >= MaxChannels {
			return nil, ErrTooManyChannels
		}

		ch := make(chan models.Event, bufferSize)
		s.eventChannels = append(s.eventChannels, ch)
		s.eventChannelCount.Add(1)

		go func() {
			s.forwardEvents(ch)
			// Decrement counter when the forwarding goroutine exits
			s.eventChannelCount.Add(-1)
		}()

		return ch, nil
	case <-time.After(ProcessingTimeout):
		return nil, ErrServiceBusy
	}
}

func (s *notificationService) forwardEvents(ch <-chan models.Event) {
	defer func() {
		s.mu.Lock()
		// Remove channel from slice
		for i, c := range s.eventChannels {
			if c == ch {
				s.eventChannels = append(s.eventChannels[:i], s.eventChannels[i+1:]...)
				break
			}
		}
		s.mu.Unlock()
	}()

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

func (s *notificationService) GetEmailChannel(ctx context.Context, priority int, bufferSize int) (<-chan models.Email, error) {
	if bufferSize <= 0 {
		return nil, ErrInvalidBufferSize
	}
	if priority < MinimumPriority || priority > MaximumPriority {
		return nil, ErrInvalidPriority
	}

	// Check channel count before acquiring lock
	if s.emailChannelCount.Load() >= MaxChannels {
		return nil, ErrTooManyChannels
	}

	lockChan := make(chan struct{})
	go func() {
		s.mu.Lock()
		close(lockChan)
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-lockChan:
		defer s.mu.Unlock()

		// Double check after acquiring lock
		if s.emailChannelCount.Load() >= MaxChannels {
			return nil, ErrTooManyChannels
		}

		ch := make(chan models.Email, bufferSize)
		s.emailChannels = append(s.emailChannels, ch)
		s.emailChannelCount.Add(1)

		go func() {
			s.forwardEmails(ch)
			// Decrement counter when the forwarding goroutine exits
			s.emailChannelCount.Add(-1)
		}()

		return ch, nil
	case <-time.After(ProcessingTimeout):
		return nil, ErrServiceBusy
	}
}

func (s *notificationService) forwardEmails(ch <-chan models.Email) {
	defer func() {
		s.mu.Lock()
		// Remove channel from slice
		for i, c := range s.emailChannels {
			if c == ch {
				s.emailChannels = append(s.emailChannels[:i], s.emailChannels[i+1:]...)
				break
			}
		}
		s.mu.Unlock()
	}()

	for {
		select {
		case <-s.ctx.Done():
			return
		case email, ok := <-ch:
			if !ok {
				return
			}
			select {
			case s.emailClientFanIn <- email:
			case <-time.After(ProcessingTimeout):
				s.logger.Error("Email forwarding timed out")
			}
		}
	}
}

func (s *notificationService) Start() {
	s.logger.Info("Starting notification service")
	s.wg.Add(4)
	go s.startAlertWorker()
	go s.startEventWorker()
	go s.startLogWorker()
	go s.startEmailWorker()
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

	// Missing close for emailChannels
	for _, ch := range s.emailChannels {
		close(ch)
	}
	s.emailChannels = nil

	// Missing close for emailClientFanIn
	close(s.emailClientFanIn)

	return nil
}

func (s *notificationService) startEmailWorker() {
	defer s.wg.Done()

	for {
		select {
		case <-s.ctx.Done():
			s.logger.Info("Email worker stopping")
			return
		case email, ok := <-s.emailClientFanIn:
			if !ok {
				return
			}
			if err := s.emailClient.Send(s.ctx, email); err != nil {
				s.logger.Error("Failed to send email", zap.Error(err))
				s.logChannel <- email
			}
		}
	}
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
