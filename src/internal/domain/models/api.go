package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
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

// The WorkflowResponse is sourced from Redis
// Redis in turn gets populed by the workflow service
// The workflow service is responsible for synchronizing
// the local state with Redis
type WorkflowResponse struct {
	// Status is the status of the workflow
	Status WorkflowStatus `json:"status"`
	// Type is the type of the workflow
	Type WorkflowType `json:"type"`
	// Year is the year of the workflow
	Year int `json:"year"`
	// WeekNumber is the week number of the workflow
	WeekNumber int `json:"week_number"`
	// BucketNumber is the bucket number of the workflow
	BucketNumber int `json:"bucket_number"`
	// BatchID is the checkpointed batch id for the workflow
	BatchID *uuid.UUID `json:"batch_id"`
	// nil if not checkpointed
	// Stage is the stage of the workflow
	Stage *int `json:"stage"`
}

// WorkflowListResponse is a list of workflows
type WorkflowListResponse struct {
	// WorkflowStatus is the status of the workflow
	WorkflowStatus WorkflowStatus `json:"workflow_status"`
	// Total is the total number of workflows
	Total int `json:"total"`
	// Workflows is the list of workflows
	Workflows []WorkflowResponse `json:"workflows"`
}
