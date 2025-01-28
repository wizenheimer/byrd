package models

import (
	"time"

	"github.com/google/uuid"
)

type JobErrorEvent struct {
	// JobID is the unique identifier of the job
	JobID uuid.UUID `json:"job_id"`
	// Error is the error message
	Error error `json:"error"`
	// Time is the time of the event
	Time time.Time `json:"time"`
}

// JobError event implementation
func (je *JobErrorEvent) GetEventType() EventType {
	return JobErrorEventType
}

func (je *JobErrorEvent) GetDistinctID() string {
	return je.JobID.String() // You might need to add JobID to JobError or get it from context
}

func (je *JobErrorEvent) GetProperties() map[string]interface{} {
	return map[string]interface{}{
		"error":     je.Error.Error(),
		"timestamp": je.Time.Unix(),
	}
}

func NewJobErrorEvent(jobContext *JobContext, jobError *JobError) *JobErrorEvent {
	return &JobErrorEvent{
		JobID: jobContext.JobID,
		Error: jobError.Error,
		Time:  time.Now(),
	}
}
