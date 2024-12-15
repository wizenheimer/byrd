package models

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

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

// String returns the string representation of WorkflowStatus
func (s WorkflowStatus) String() string {
	switch s {
	case WorkflowStatusPending:
		return "pending"
	case WorkflowStatusRunning:
		return "running"
	case WorkflowStatusCompleted:
		return "completed"
	case WorkflowStatusFailed:
		return "failed"
	case WorkflowStatusAborted:
		return "aborted"
	case WorkflowStatusExpired:
		return "expired"
	default:
		return "unknown"
	}
}

// ParseWorkflowStatus converts a string to WorkflowStatus
func ParseWorkflowStatus(s string) (WorkflowStatus, error) {
	switch strings.ToLower(s) {
	case "running":
		return WorkflowStatusRunning, nil
	case "expired":
		return WorkflowStatusExpired, nil
	case "pending":
		return WorkflowStatusPending, nil
	case "completed":
		return WorkflowStatusCompleted, nil
	case "failed":
		return WorkflowStatusFailed, nil
	case "aborted":
		return WorkflowStatusAborted, nil
	default:
		return WorkflowStatusExpired, fmt.Errorf("invalid workflow status: %s", s)
	}
}

// Checkpoint represents a checkpoint for a workflow
type Checkpoint struct {
	// BatchID is the checkpointed batch id for the workflow
	BatchID *string `json:"batch_id"`
	// nil if not checkpointed
	// Stage is the stage of the workflow
	Stage *int `json:"stage"`
}

// WorkflowIdentifier uniquely identifies a workflow
type WorkflowIdentifier struct {
	// Type is the type of the workflow
	Type WorkflowType `json:"type"`
	// Year is the year of the workflow
	Year int `json:"year"`
	// WeekNumber is the week number of the workflow
	WeekNumber int `json:"week_number"`
	// BucketNumber is the bucket number of the workflow
	BucketNumber int `json:"bucket_number"`
}

// formatWorkflowID creates a consistent string representation
func (id *WorkflowIdentifier) Serialize(prefix string, status WorkflowStatus) string {
	return fmt.Sprintf("%s-%s-%s-%d-%d-%d",
		prefix,
		status.String(),
		id.Type.String(),
		id.Year,
		id.WeekNumber,
		id.BucketNumber)
}

// parseWorkflowID parses the string representation into a WorkflowIdentifier
func ParseWorkflowID(serialized string) (*WorkflowIdentifier, string, WorkflowStatus, error) {
	parts := strings.Split(serialized, "-")
	if len(parts) != 6 {
		return nil, "workflow", WorkflowStatusExpired,
			fmt.Errorf("invalid workflow id format: expected 6 parts, got %d", len(parts))
	}

	id := &WorkflowIdentifier{}
	var err error

	prefix := parts[0]
	status, err := ParseWorkflowStatus(parts[1]) // You'll need this function
	if err != nil {
		return nil, "workflow", WorkflowStatusExpired, fmt.Errorf("invalid status: %w", err)
	}

	id.Type, err = ParseWorkflowType(parts[2]) // You'll need this function
	if err != nil {
		return nil, "workflow", WorkflowStatusExpired, fmt.Errorf("invalid type: %w", err)
	}
	id.Year, err = strconv.Atoi(parts[3])
	if err != nil {
		return nil, "workflow", WorkflowStatusExpired, fmt.Errorf("invalid year: %w", err)
	}

	id.WeekNumber, err = strconv.Atoi(parts[4])
	if err != nil {
		return nil, "workflow", WorkflowStatusExpired, fmt.Errorf("invalid week number: %w", err)
	}

	id.BucketNumber, err = strconv.Atoi(parts[5])
	if err != nil {
		return nil, "workflow", WorkflowStatusExpired, fmt.Errorf("invalid bucket number: %w", err)
	}

	return id, prefix, status, nil
}

// WorkflowType represents the type of a workflow
// This can be either a screenshot or report workflow
type WorkflowType string

const (
	// ScreenshotWorkflowType is the type of workflow that takes screenshots
	ScreenshotWorkflowType WorkflowType = "screenshot"
	// ReportWorkflowType is the type of workflow that generates reports
	ReportWorkflowType WorkflowType = "report"
)

const (
	ScreenshotWorkflowRepositoryPrefix string = "sswf"
	ReportWorkflowRepositoryPrefix     string = "rpwf"
)

func (wt WorkflowType) Prefix() string {
	switch wt {
	case ScreenshotWorkflowType:
		return ScreenshotWorkflowRepositoryPrefix
	case ReportWorkflowType:
		return ReportWorkflowRepositoryPrefix
	default:
		return ScreenshotWorkflowRepositoryPrefix
	}
}

func (wt WorkflowType) String() string {
	switch wt {
	case ScreenshotWorkflowType:
		return "screenshot"
	case ReportWorkflowType:
		return "report"
	default:
		return "unknown"
	}
}

func ParseWorkflowType(s string) (WorkflowType, error) {
	switch strings.ToLower(s) {
	case "screenshot":
		return ScreenshotWorkflowType, nil
	case "report":
		return ReportWorkflowType, nil
	default:
		return ScreenshotWorkflowType, fmt.Errorf("invalid workflow type: %s", s)
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
	// Cancel is the context cancel function
	Cancel context.CancelFunc
	// ExecutorID is the executor id
	ExecutorID uuid.UUID
	// Status is the status of the workflow
	Status WorkflowStatus
	Mutex  sync.RWMutex
}
