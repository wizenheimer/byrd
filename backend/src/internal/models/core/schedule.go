package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
)

type ScheduleID uuid.UUID

func NewScheduleID() ScheduleID {
	return ScheduleID(uuid.New())
}

func (id ScheduleID) String() string {
	return uuid.UUID(id).String()
}

// ScheduledFunc represents a scheduled function
type ScheduledFunc struct {
	// Unique identifier for the scheduled function
	ID ScheduleID

	// Schedule specification
	Spec string

	// Cron entry ID - used to manage the scheduled function
	EntryID cron.EntryID

	// Command to run
	// Command func()

	// Last run time - time in past
	// This is time when the function was last run
	LastRun time.Time

	// Next run time - time in future
	// This is the time when the function is scheduled to run next
	NextRun time.Time

	// Flag to indicate if the function is delayed
	IsDelayed bool

	// Time to delay the function
	DelayUntil time.Time
}
