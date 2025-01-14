package scheduler

import (
	"time"

	models "github.com/wizenheimer/byrd/src/internal/models/core"
)

// Scheduler is an interface for scheduling functions
type Scheduler interface {
	// Start the scheduler
	Start() error

	// Stop the scheduler
	Stop() error

	// Recover recovers scheduled functions that got pre-empted due to a restart
	Recover(scheduleSpec string, cmd func(), lastRun *time.Time, nextRun *time.Time) (*models.ScheduledFunc, error)

	// Schedule a function to run based on the schedule specification
	Schedule(scheduleSpec string, cmd func()) (*models.ScheduledFunc, error)

	// Schedule a function to run based on the schedule specification with a delay
	ScheduleWithDelay(scheduleSpec string, delay time.Duration, cmd func()) (*models.ScheduledFunc, error)

	// Update a scheduled function with a new schedule specification and command
	Update(id models.ScheduleID, scheduleSpec string, cmd func()) error

	// Delete a scheduled function
	Delete(id models.ScheduleID) error

	// Get a scheduled function by ID
	Get(id models.ScheduleID) (*models.ScheduledFunc, error)

	// List all scheduled functions
	List() []*models.ScheduledFunc
}
