package models

import (
	"fmt"
	"time"

	"github.com/wizenheimer/iris/src/pkg/utils/ptr"
)

type WorkflowRequest struct {
	// Type is the type of the workflow
	Type string `json:"type"`
	// Year is the year of the workflow
	Year *int `json:"year"`
	// WeekNumber is the week number of the workflow
	WeekNumber *int `json:"week_number"`
	// BucketNumber is the bucket number of the workflow
	BucketNumber *int `json:"bucket_number"`
}

func (wr *WorkflowRequest) Validate(safe bool) error {
	now := time.Now()

	// If the year is not set, set it based on the current year
	if wr.Year == nil {
		if safe {
			// Set current year
			wr.Year = ptr.To(now.Year())
		} else {
			return fmt.Errorf("year is required")
		}
	}

	// If the week number is not set, set it based on the current week
	if wr.WeekNumber == nil {
		if safe {
			// Set current week number
			_, week := now.ISOWeek()
			wr.WeekNumber = ptr.To(week)
		} else {
			return fmt.Errorf("week number is required")
		}
	}

	// If the bucket number is not set, set it based on the current day of the week
	if wr.BucketNumber == nil {
		if safe {
			// Set current bucket number
			currentWeekday := now.Weekday()
			if currentWeekday == 0 {
				currentWeekday = 7 // Sunday is 0, but we want it to be 7
			}

			// Calculate the bucket number
			if currentWeekday <= 3 { // Monday, Tuesday, Wednesday
				wr.BucketNumber = ptr.To(1)
			} else { // Thursday, Friday, Saturday, Sunday
				wr.BucketNumber = ptr.To(2)
			}
		} else {
			return fmt.Errorf("bucket number is required")
		}
	}
	return nil
}

type WorkflowStatus string

const (
	// WorkflowStatusPending is the status of a workflow that has not yet started
	// This is set to happen in the future
	WorkflowStatusPending WorkflowStatus = "pending"
	// WorkflowStatusRunning is the status of a workflow that is currently running
	WorkflowStatusRunning WorkflowStatus = "running"
	// WorkflowStatusCompleted is the status of a workflow that has completed
	WorkflowStatusCompleted WorkflowStatus = "completed"
	// WorkflowStatusFailed is the status of a workflow that has failed to start
	WorkflowStatusFailed WorkflowStatus = "failed"
	// WorkflowStatusAborted is the status of a workflow that has been aborted
	WorkflowStatusAborted WorkflowStatus = "aborted"
	// WorkflowStatusExpired is the status of a workflow that is unknown
	// This happens when Redis TTL expires
	// Currently TTL is set to 4 days (1/2 week)
	WorkflowStatusExpired WorkflowStatus = "expired"
)

// The WorkflowResponse is sourced from Redis
// Redis in turn gets populed by the workflow service
// The workflow service is responsible for synchronizing
// the local state with Redis
type WorkflowResponse struct {
	// Status is the status of the workflow
	Status WorkflowStatus `json:"status"`
	// Type is the type of the workflow
	Type string `json:"type"`
	// Year is the year of the workflow
	Year int `json:"year"`
	// WeekNumber is the week number of the workflow
	WeekNumber int `json:"week_number"`
	// BucketNumber is the bucket number of the workflow
	BucketNumber int `json:"bucket_number"`
	// BatchID is the checkpointed batch id for the workflow
	BatchID *string `json:"batch_id"`
	// nil if not checkpointed
}

type WorkflowListResponse struct {
	// Workflows is the list of workflows
	Workflows []WorkflowResponse `json:"workflows"`
	// Total is the total number of workflows
	Total int `json:"total"`
}

// WorkflowIdentifier uniquely identifies a workflow
type WorkflowIdentifier struct {
	// Type is the type of the workflow
	Type string
	// Year is the year of the workflow
	Year int
	// WeekNumber is the week number of the workflow
	WeekNumber int
	// BucketNumber is the bucket number of the workflow
	BucketNumber int
}

// formatWorkflowID creates a consistent string representation
func (id *WorkflowIdentifier) Serialize() string {
	return fmt.Sprintf("%s-%d-%d-%d", id.Type, id.Year, id.WeekNumber, id.BucketNumber)
}

// parseWorkflowID parses the string representation into a WorkflowIdentifier
func ParseWorkflowID(serialized string) (*WorkflowIdentifier, error) {
	var id WorkflowIdentifier
	_, err := fmt.Sscanf(serialized, "%s-%d-%d-%d", &id.Type, &id.Year, &id.WeekNumber, &id.BucketNumber)
	if err != nil {
		return nil, err
	}
	return &id, nil
}
