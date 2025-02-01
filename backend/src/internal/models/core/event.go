// ./src/internal/models/core/event.go
package models

// EventType is the type of event
type EventType string

const (
	JobErrorEventType          EventType = "job_error"
	JobLifecycleEventType      EventType = "job_lifecycle"
	ScheduleErrorEventType     EventType = "schedule_error"
	ScheduleLifecycleEventType EventType = "schedule_lifecycle"
)

// Event represents the interface that any type can implement to be tracked as an event
type Event interface {
	// GetEventType returns the type of the event
	GetEventType() EventType

	// GetDistinctID returns the unique identifier for the event
	GetDistinctID() string

	// GetProperties returns the event properties as a key-value map
	GetProperties() map[string]interface{}
}
