package models

import (
	"context"
	"fmt"
	"sync"
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

// WorkflowStatus represents the status of a workflow
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

type Checkpoint struct {
	// BatchID is the checkpointed batch id for the workflow
	BatchID *string `json:"batch_id"`
	// nil if not checkpointed
	// Stage is the stage of the workflow
	Stage *int `json:"stage"`
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
	BatchID *string `json:"batch_id"`
	// nil if not checkpointed
	// Stage is the stage of the workflow
	Stage *int `json:"stage"`
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
	Type WorkflowType
	// Year is the year of the workflow
	Year int
	// WeekNumber is the week number of the workflow
	WeekNumber int
	// BucketNumber is the bucket number of the workflow
	BucketNumber int
}

// formatWorkflowID creates a consistent string representation
func (id *WorkflowIdentifier) Serialize(prefix string, status WorkflowStatus) string {
	return fmt.Sprintf("%s-%v-%s-%d-%d-%d", prefix, status, id.Type, id.Year, id.WeekNumber, id.BucketNumber)
}

// parseWorkflowID parses the string representation into a WorkflowIdentifier
func ParseWorkflowID(serialized string) (*WorkflowIdentifier, string, WorkflowStatus, error) {
	var id WorkflowIdentifier
	var status WorkflowStatus
	var prefix string
	_, err := fmt.Sscanf(serialized, "%s-%v-%v-%d-%d-%d", &prefix, &status, &id.Type, &id.Year, &id.WeekNumber, &id.BucketNumber)
	if err != nil {
		return nil, "workflow", WorkflowStatusExpired, fmt.Errorf("failed to parse workflow id: %w", err)
	}
	return &id, prefix, status, nil
}

type WorkflowType string

const (
	ScreenshotWorkflowType WorkflowType = "screenshot"
	ReportWorkflowType     WorkflowType = "report"
)

const (
	ScreenshotWorkflowRepositoryPrefix string = "workflow-screenshot"
	ReportWorkflowRepositoryPrefix     string = "workflow-report"
)

func GetWorkflowPrefixFromWorkflowType(wfType WorkflowType) string {
	switch wfType {
	case ScreenshotWorkflowType:
		return ScreenshotWorkflowRepositoryPrefix
	case ReportWorkflowType:
		return ReportWorkflowRepositoryPrefix
	default:
		return ScreenshotWorkflowRepositoryPrefix
	}
}

func GetWorkflowTypeFromWorkflowPrefix(prefix string) WorkflowType {
	switch prefix {
	case ScreenshotWorkflowRepositoryPrefix:
		return ScreenshotWorkflowType
	case ReportWorkflowRepositoryPrefix:
		return ReportWorkflowType
	default:
		return ScreenshotWorkflowType
	}
}

// WorkflowUpdate represents a heartbeat or status update from the executor
type WorkflowUpdate struct {
	// ID is the identifier of the workflow
	ID *WorkflowIdentifier `json:"id"`
	// Checkpoint is the checkpointed batch id for the workflow
	Checkpoint *Checkpoint `json:"checkpoint"`
	// Timestamp is the time of the update
	Timestamp time.Time `json:"timestamp"`
	// Status is the status of the workflow
	Status WorkflowStatus `json:"status"`
}

// WorkflowError represents an error from the executor
type WorkflowError struct {
	// ID is the identifier of the workflow
	ID *WorkflowIdentifier `json:"id"`
	// Error is the error received
	Error error `json:"error"`
	// Timestamp is the time of the error
	Timestamp time.Time `json:"timestamp"`
}

// WorkflowState to track running workflows
type WorkflowState struct {
	Cancel     context.CancelFunc
	ExecutorID uuid.UUID
	Status     WorkflowStatus
	Mutex      sync.RWMutex
}
