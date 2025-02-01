package models

import (
	"time"
)

type ScheduleErrorEvent struct {
	// ScheduleID is the unique identifier of the schedule
	ScheduleID ScheduleID `json:"schedule_id"`

	// Error is the error message
	Error error `json:"error"`

	// Time is the time of the event
	Time time.Time `json:"time"`
}

// ScheduleError event implementation
func (se *ScheduleErrorEvent) GetEventType() EventType {
	return ScheduleErrorEventType
}

func (se *ScheduleErrorEvent) GetDistinctID() string {
	return se.ScheduleID.String()
}

func (se *ScheduleErrorEvent) GetProperties() map[string]interface{} {
	return map[string]interface{}{
		"error":     se.Error.Error(),
		"timestamp": se.Time.Unix(),
	}
}

func NewScheduleErrorEvent(scheduleID ScheduleID, err error) *ScheduleErrorEvent {
	return &ScheduleErrorEvent{
		ScheduleID: scheduleID,
		Error:      err,
		Time:       time.Now(),
	}
}

type ScheduleStage string

const (
	RecoveryStage ScheduleStage = "recovery"
  ShutdownStage ScheduleStage = "shutdown"
)

type ScheduleEvent struct {
	ScheduleID *ScheduleID   `json:"schedule_id"`
	TimeStamp  time.Time     `json:"time_stamp"`
	Stage      ScheduleStage `json:"stage"`
	Message    string        `json:"message"`
}

func NewScheduleEvent(scheduleID *ScheduleID, stage ScheduleStage, msg string) *ScheduleEvent {
	return &ScheduleEvent{
		ScheduleID: scheduleID,
		TimeStamp:  time.Now(),
		Stage:      stage,
		Message:    msg,
	}
}

func (se *ScheduleEvent) GetEventType() EventType {
	return ScheduleErrorEventType
}

func (se *ScheduleEvent) GetDistinctID() string {
	return se.ScheduleID.String()
}

func (se *ScheduleEvent) GetProperties() map[string]interface{} {
	return map[string]interface{}{
		"message":   se.Message,
		"timestamp": se.TimeStamp.Unix(),
	}
}
