package scheduler

import (
	"fmt"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/pkg/logger"
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
	parser := cron.NewParser(
		cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)

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
		s.logger.Info("scheduler already running", zap.String("status", "running"))
		return nil
	}

	s.cron.Start()
	s.running = true

	s.logger.Info("scheduler started", zap.String("status", "running"))
	return nil
}

// Schedule a function to run based on the schedule specification
// This doesn't trigger the function immediately
// It schedules the function to run based on the schedule specification
func (s *scheduler) Schedule(scheduleSpec string, cmd func()) (*models.ScheduledFunc, error) {
	s.logger.Info("scheduling a new function", zap.String("scheduleSpec", scheduleSpec))
	// Add validation
	if cmd == nil {
		return nil, fmt.Errorf("command function cannot be nil")
	}

	if _, err := s.parse(scheduleSpec); err != nil {
		return nil, fmt.Errorf("failed to parse schedule: %w", err)
	}

	id := models.NewScheduleID()
	wrappedCmd := s.wrapCommand(id, cmd)

	entryID, err := s.cron.AddFunc(scheduleSpec, wrappedCmd)
	if err != nil {
		return nil, fmt.Errorf("failed to schedule function: %w", err)
	}

	entry := s.cron.Entry(entryID)
	scheduledFunc := models.ScheduledFunc{
		ID:      id,
		Spec:    scheduleSpec,
		EntryID: entryID,
		LastRun: time.Time{}, // Initialize with zero time
		NextRun: entry.Next,
	}

	s.schedules.Store(id, &scheduledFunc)
	return &scheduledFunc, nil
}

// Schedule a function to run based on the schedule specification with a delay
// This schedules the function to run after the specified delay
// It doesn't trigger the function immediately
func (s *scheduler) ScheduleWithDelay(scheduleSpec string, delay time.Duration, cmd func()) (*models.ScheduledFunc, error) {
	s.logger.Info("scheduling a new function with delay", zap.String("scheduleSpec", scheduleSpec), zap.Duration("delay", delay))
	// Add validation
	if cmd == nil {
		return nil, fmt.Errorf("command function cannot be nil")
	}
	if scheduleSpec == "" {
		return nil, fmt.Errorf("schedule specification cannot be empty")
	}

	// Validate scheduleSpec
	_, err := s.parse(scheduleSpec)
	if err != nil {
		return nil, fmt.Errorf("failed to parse schedule: %w", err)
	}

	id := models.NewScheduleID()
	delayUntil := time.Now().Add(delay)

	scheduledFunc := models.ScheduledFunc{
		ID:         id,
		Spec:       scheduleSpec,
		IsDelayed:  true,
		DelayUntil: delayUntil,
		LastRun:    time.Time{}, // Initialize with zero time
	}
	s.schedules.Store(id, &scheduledFunc)

	go func() {
		timer := time.NewTimer(delay)
		defer timer.Stop()

		<-timer.C

		s.mu.Lock()
		if !s.running {
			s.mu.Unlock()
			return
		}
		s.mu.Unlock()

		wrappedCmd := s.wrapCommand(id, cmd)
		entryID, err := s.cron.AddFunc(scheduleSpec, wrappedCmd)
		if err != nil {
			return
		}

		entry := s.cron.Entry(entryID)
		if existing, ok := s.schedules.Load(id); ok {
			sf := existing.(*models.ScheduledFunc)
			sf.EntryID = entryID
			sf.NextRun = entry.Next
			sf.IsDelayed = false
			s.schedules.Store(id, sf)
		}
	}()

	return &scheduledFunc, nil
}

// Recover recovers scheduled functions that got pre-empted due to a restart
// It executes the command immediately if the next run time is in the past
func (s *scheduler) Recover(scheduleSpec string, cmd func(), lastRun *time.Time, nextRun *time.Time) (*models.ScheduledFunc, error) {
	s.logger.Info("recovering a scheduled function", zap.String("scheduleSpec", scheduleSpec), zap.Any("lastRun", lastRun), zap.Any("nextRun", nextRun))

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
func (s *scheduler) Update(id models.ScheduleID, scheduleSpec string, cmd func()) error {
	s.logger.Info("updating a scheduled function", zap.String("scheduleSpec", scheduleSpec), zap.Any("id", id))

	sf, err := s.Get(id)
	if err != nil {
		return err
	}

	// Remove existing function
	s.cron.Remove(sf.EntryID)

	// Schedule new function with wrapped command
	wrappedCmd := s.wrapCommand(id, cmd)
	entryID, err := s.cron.AddFunc(scheduleSpec, wrappedCmd)
	if err != nil {
		return fmt.Errorf("failed to update function: %w", err)
	}

	entry := s.cron.Entry(entryID)
	sf.Spec = scheduleSpec
	sf.EntryID = entryID
	// sf.Command = cmd
	sf.NextRun = entry.Next

	s.schedules.Store(id, sf)
	return nil
}

// Delete a scheduled function
// This stops the function from running in the future
func (s *scheduler) Delete(id models.ScheduleID) error {
	s.logger.Info("deleting a scheduled function", zap.Any("id", id))

	sf, err := s.Get(id)
	if err != nil {
		return err
	}

	if !sf.IsDelayed {
		s.cron.Remove(sf.EntryID)
	}
	s.schedules.Delete(id)
	return nil
}

// Get a scheduled function by ID
// This returns the scheduled function with the last run and next run times
func (s *scheduler) Get(id models.ScheduleID) (*models.ScheduledFunc, error) {
	s.logger.Info("getting a scheduled function", zap.Any("id", id))

	if value, ok := s.schedules.Load(id); ok {
		sf := value.(*models.ScheduledFunc)
		if !sf.IsDelayed && sf.EntryID > 0 {
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
	s.logger.Info("listing all scheduled functions")

	var funcs []*models.ScheduledFunc
	s.schedules.Range(func(key, value interface{}) bool {
		sf := value.(*models.ScheduledFunc)
		if !sf.IsDelayed && sf.EntryID > 0 {
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
	s.logger.Info("stopping the scheduler")

	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		s.logger.Info("scheduler already stopped", zap.String("status", "stopped"))
		return nil
	}

	s.cron.Stop()
	s.running = false
	s.logger.Info("scheduler stopped", zap.String("status", "stopped"))
	return nil
}

// This wrapper function is used to execute the command safely
// It also updates the LastRun and NextRun times for the scheduled function once the command is executed in the future
func (s *scheduler) wrapCommand(id models.ScheduleID, cmd func()) func() {
	return func() {
		if value, ok := s.schedules.Load(id); ok {
			sf := value.(*models.ScheduledFunc)

			// Execute the command safely
            s.logger.Info("executing scheduled function", zap.Any("id", id))
			s.safeExecute(cmd)

			// Update LastRun and NextRun times
            s.logger.Info("updating last run and next run times", zap.Any("id", id))
			entry := s.cron.Entry(sf.EntryID)
			sf.LastRun = entry.Prev
			sf.NextRun = entry.Next

			// Store the updated scheduled function
            s.logger.Info("storing updated scheduled function", zap.Any("id", id))
			s.schedules.Store(id, sf)
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

	f := value.(*models.ScheduledFunc)

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

	f := value.(*models.ScheduledFunc)

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
