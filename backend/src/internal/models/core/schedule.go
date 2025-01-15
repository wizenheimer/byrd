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

func NilScheduleID() ScheduleID {
    return ScheduleID(uuid.Nil)
}

type ScheduleFuncState string

const (
    // delayed - function is delayed, will be run at a later time
    DelayedFuncState ScheduleFuncState = "delayed"

    // active - function is active, will be run at the next scheduled time
    ActiveFuncState ScheduleFuncState = "active"

    // stale - function is stale, will not be run at the next scheduled time
    StaleFuncState ScheduleFuncState = "stale"
)

// ScheduledFunc represents a scheduled function
type ScheduledFunc struct {
	// Unique identifier for the scheduled function
	ID ScheduleID

	// Schedule specification
	Spec string

	// Cron entry ID - used to manage the scheduled function
	EntryID cron.EntryID

    // State of the scheduled function
    State ScheduleFuncState

	// Last run time - time in past
	// This is time when the function was last run
	LastRun time.Time

	// Next run time - time in future
	// This is the time when the function is scheduled to run next
	NextRun time.Time

	// Time to delay the function
	DelayUntil time.Time
}
