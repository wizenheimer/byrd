package models

import "time"

// EventType is the type of event
type EventType string

// EventSource is the source of the event
type EventSource string

// Event is a generic event model
type Event struct {
	// Data is the event data
	Data any
	// Timestamp is the time the event occurred
	Time time.Time
	// Type is the event type
	Type EventType
	// Source is the event source
	Source EventSource
}

func NewEvent(data any, eventType EventType, eventSource EventSource) Event {
	return Event{
		Data:   data,
		Time:   time.Now(),
		Type:   eventType,
		Source: eventSource,
	}
}
