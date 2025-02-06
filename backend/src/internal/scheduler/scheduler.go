// ./src/internal/scheduler/scheduler.go
package scheduler

import (
	"fmt"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"github.com/wizenheimer/byrd/src/pkg/utils"
	"go.uber.org/zap"
)

type scheduler struct {
	cron      *cron.Cron
	schedules sync.Map // map[models.ScheduleID]*models.ScheduledFunc
	running   bool
	mu        sync.Mutex
	parser    cron.Parser
	logger    *logger.Logger
}

// NewScheduler creates a new instance of the scheduler
func NewScheduler(logger *logger.Logger) Scheduler {
	cronLogger := NewCronLogger(logger)
	parser := utils.NewScheduleParser()

	return &scheduler{
		// Initialize the cron instance with custom parser and logger
		cron: cron.New(
			cron.WithSeconds(),
			cron.WithParser(parser),
			cron.WithLogger(cronLogger),
		),

		// Add the logger
		logger: logger,

		// Add the parser
		parser: parser,
	}
}

// Start the scheduler instance
func (s *scheduler) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return nil
	}

	s.cron.Start()
	s.running = true

	scheduleCount := 0
	s.schedules.Range(func(_, _ interface{}) bool {
		scheduleCount++
		return true
	})

	return nil
}

func (s *scheduler) Schedule(cmd func(), opts ScheduleOptions) (*models.ScheduledFunc, error) {
	// Validate scheduleSpec and cmd
	if cmd == nil {
		return nil, fmt.Errorf("command function cannot be nil")
	}
	if _, err := s.parse(opts.ScheduleSpec); err != nil {
		return nil, fmt.Errorf("failed to parse schedule: %w", err)
	}

	// Create a new ID for the scheduled function
	id := models.NewScheduleID()

	// Wrap the command with the ID and hooks
	wrappedCmd := s.wrapCommand(id, cmd, opts.Hooks...)

	var scheduledFunc models.ScheduledFunc

	// Check if the function should be delayed
	if opts.Delay > 0 {
		// Schedule the function with a delay
		scheduledFunc = models.ScheduledFunc{
			ID:         id,
			Spec:       opts.ScheduleSpec,
			DelayUntil: time.Now().Add(opts.Delay),
			LastRun:    time.Time{}, // Initialize with zero time
			State:      models.DelayedFuncState,
		}

		// Schedule the function after the delay
		go func() {
			// Delay the function
			timer := time.NewTimer(opts.Delay)
			defer timer.Stop()

			<-timer.C

			s.mu.Lock()
			if !s.running {
				s.mu.Unlock()
				return
			}

			s.mu.Unlock()

			// Schedule the function
			wrappedCmd = s.wrapCommand(id, cmd, opts.Hooks...)
			entryID, err := s.cron.AddFunc(opts.ScheduleSpec, wrappedCmd)
			if err != nil {
				return
			}

			entry := s.cron.Entry(entryID)
			if existing, ok := s.schedules.Load(id); ok {
				sf, ok := existing.(*models.ScheduledFunc)
				if !ok {
					s.logger.Error("failed to load scheduled function due to type mismatch", zap.Any("id", id))
					return
				}
				sf.EntryID = entryID
				sf.NextRun = entry.Next
				sf.State = models.ActiveFuncState
				s.schedules.Store(id, sf)
			}
		}()
	} else {
		// Schedule the function immediately
		entryID, err := s.cron.AddFunc(opts.ScheduleSpec, wrappedCmd)
		if err != nil {
			return nil, fmt.Errorf("failed to schedule function: %w", err)
		}

		entry := s.cron.Entry(entryID)
		scheduledFunc = models.ScheduledFunc{
			ID:      id,
			Spec:    opts.ScheduleSpec,
			State:   models.ActiveFuncState,
			EntryID: entryID,
			LastRun: time.Time{}, // Initialize with zero time
			NextRun: entry.Next,
		}
	}

	s.schedules.Store(id, &scheduledFunc)
	return &scheduledFunc, nil
}

// Recover recovers scheduled functions that got pre-empted due to a restart
// It executes the command immediately if the next run time is in the past
func (s *scheduler) Recover(scheduleSpec string, cmd func(), lastRun *time.Time, nextRun *time.Time) (*models.ScheduledFunc, error) {
	currentTime := time.Now()

	// Validate schedule
	schedule, err := s.parse(scheduleSpec)
	if err != nil {
		return nil, fmt.Errorf("failed to parse schedule: %w", err)
	}

	// Set default times if not provided
	// Cold start scenario
	if lastRun == nil {
		lastRun = &currentTime
	}

	// Calculate next run if not provided
	// Cold start scenario
	if nextRun == nil {
		next := schedule.Next(*lastRun)
		nextRun = &next
	}

	// If the next run time is in the past, we need to:
	// Execute the command immediately
	// Sync the state with the schedule
	if nextRun.Before(currentTime) {
		// Execute the command
		s.safeExecute(cmd)

		// Update last run time to current time
		lastRun = &currentTime

		// Calculate next run time based on the schedule
		r := schedule.Next(*lastRun)
		nextRun = &r
	}

	// Create a new ID for the recovered function
	id := models.NewScheduleID()
	wrappedCmd := s.wrapCommand(id, cmd)

	// Schedule the function
	entryID, err := s.cron.AddFunc(scheduleSpec, wrappedCmd)
	if err != nil {
		return nil, fmt.Errorf("failed to schedule recovered function: %w", err)
	}

	// Create the scheduled function with provided times
	f := &models.ScheduledFunc{
		ID:      id,
		Spec:    scheduleSpec,
		EntryID: entryID,
		LastRun: *lastRun,
		NextRun: *nextRun,
	}

	// Sync the next run time with the cron entry
	if entry := s.cron.Entry(entryID); !entry.Next.Equal(f.NextRun) {
		f.NextRun = entry.Next
	}

	// Store it in our map
	s.schedules.Store(id, f)

	return f, nil
}

// Update a scheduled function with a new schedule specification and command
// This doesn't stop a running function
// It merely updates the schedule specification and command for the next run
func (s *scheduler) Update(id models.ScheduleID, cmd func(), opts ScheduleOptions) (*models.ScheduledFunc, error) {
	// Get the existing scheduled function
	sf, err := s.Get(id)
	if err != nil {
		return nil, err
	}

	// Remove existing function
	s.cron.Remove(sf.EntryID)
	sf.State = models.StaleFuncState

	// Schedule the new function
	return s.Schedule(cmd, opts)
}

// Delete a scheduled function
// This stops the function from running in the future
func (s *scheduler) Delete(id models.ScheduleID) error {
	// Get the scheduled function
	sf, err := s.Get(id)
	if err != nil {
		return err
	}

	// If the function is delayed, we don't need to remove it from the cron because it hasn't been scheduled yet
	if sf.State != models.DelayedFuncState {
		s.cron.Remove(sf.EntryID)
		sf.State = models.StaleFuncState
	}

	// Delete the scheduled function
	s.schedules.Delete(id)
	return nil
}

// Get a scheduled function by ID
// This returns the scheduled function with the last run and next run times
func (s *scheduler) Get(id models.ScheduleID) (*models.ScheduledFunc, error) {
	if value, ok := s.schedules.Load(id); ok {
		sf, ok := value.(*models.ScheduledFunc)
		if !ok {
			return nil, fmt.Errorf("failed to load scheduled function due to type mismatch")
		}
		// If the function is not delayed,
		// and has been scheduled
		// update the last run and next run times
		if sf.State != models.DelayedFuncState && sf.EntryID > 0 {
			entry := s.cron.Entry(sf.EntryID)
			sf.LastRun = entry.Prev
			sf.NextRun = entry.Next
		}
		return sf, nil
	}
	return nil, fmt.Errorf("scheduled function not found")
}

// List all scheduled functions
// This returns all the scheduled functions with the last run and next run times
func (s *scheduler) List() []*models.ScheduledFunc {
	var funcs []*models.ScheduledFunc
	s.schedules.Range(func(key, value interface{}) bool {
		sf, ok := value.(*models.ScheduledFunc)
		if !ok {
			s.logger.Error("failed to load scheduled function due to type mismatch", zap.Any("id", key))
			return true
		}
		// If the function is not delayed,
		// and has been scheduled
		// update the last run and next run times
		if sf.State != models.DelayedFuncState && sf.EntryID > 0 {
			entry := s.cron.Entry(sf.EntryID)
			sf.LastRun = entry.Prev
			sf.NextRun = entry.Next
		}
		funcs = append(funcs, sf)
		return true
	})
	return funcs
}

// Stop the scheduler
// This stops the scheduler from running
func (s *scheduler) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return nil
	}

	s.schedules.Range(func(key, value interface{}) bool {
		scheduleKey, ok := key.(models.ScheduleID)
		if !ok {
			s.logger.Error("failed to load scheduled function due to type mismatch", zap.Any("id", key))
			return true
		}
		if err := s.Delete(scheduleKey); err != nil {
			s.logger.Error("failed to delete scheduled function", zap.Any("id", key), zap.Error(err))
		}
		return true
	})

	s.cron.Stop()
	s.running = false
	return nil
}

// This wrapper function is used to execute the command safely
// It also updates the LastRun and NextRun times for the scheduled function once the command is executed in the future
func (s *scheduler) wrapCommand(id models.ScheduleID, cmd func(), hooks ...func()) func() {
	return func() {
		if value, ok := s.schedules.Load(id); ok {
			sf, ok := value.(*models.ScheduledFunc)
			if !ok {
				s.logger.Error("failed to load scheduled function due to type mismatch", zap.Any("id", id))
				return
			}

			// Execute the command safely
			s.safeExecute(cmd)

			// Update LastRun and NextRun times
			entry := s.cron.Entry(sf.EntryID)
			sf.LastRun = entry.Prev
			sf.NextRun = entry.Next

			// Store the updated scheduled function
			s.schedules.Store(id, sf)

			// Execute hooks
			for _, hook := range hooks {
				s.safeExecute(hook)
			}
		}
	}
}

// This function is used to execute the command safely
// It recovers from panics and logs the error
func (s *scheduler) safeExecute(cmd func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				// Log the panic or handle it appropriately
				s.logger.Error("panic occurred while executing scheduled function", zap.Any("recovered", r))
			}
		}()
		cmd()
	}()
}

// Get the next run time for a scheduled function based on the schedule specification
func (s *scheduler) NextRun(scheduleID models.ScheduleID, runtime *time.Time) (time.Time, error) {
	// Get the scheduled function
	value, ok := s.schedules.Load(scheduleID)
	if !ok {
		return time.Time{}, fmt.Errorf("scheduled function not found")
	}

	f, ok := value.(*models.ScheduledFunc)
	if !ok {
		return time.Time{}, fmt.Errorf("failed to load scheduled function due to type mismatch")
	}

	var currentTime time.Time
	if runtime == nil {
		currentTime = time.Now()
	}
	schedule, err := s.parse(f.Spec)
	if err != nil {
		return currentTime, err
	}

	return schedule.Next(currentTime), nil
}

// Get the previous run time for a scheduled function based on the schedule specification
func (s *scheduler) PrevRun(scheduleID models.ScheduleID, runtime *time.Time) (time.Time, error) {
	// Get the scheduled function
	value, ok := s.schedules.Load(scheduleID)
	if !ok {
		return time.Time{}, fmt.Errorf("scheduled function not found")
	}

	f, ok := value.(*models.ScheduledFunc)
	if !ok {
		return time.Time{}, fmt.Errorf("failed to load scheduled function due to type mismatch")
	}

	var currentTime time.Time
	if runtime == nil {
		currentTime = time.Now()
	}
	schedule, err := s.parse(f.Spec)
	if err != nil {
		return currentTime, fmt.Errorf("failed to parse schedule: %w", err)
	}

	// Start from LastRun
	prevTime := f.LastRun
	nextTime := schedule.Next(prevTime)

	// Keep moving forward until we find the last time before current time
	for nextTime.Before(currentTime) {
		prevTime = nextTime
		nextTime = schedule.Next(prevTime)
	}

	return prevTime, nil
}

// Parse a schedule specification
func (s *scheduler) parse(scheduleSpec string) (cron.Schedule, error) {
	return s.parser.Parse(scheduleSpec)
}
