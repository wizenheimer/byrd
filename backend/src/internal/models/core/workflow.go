// ./src/internal/models/core/workflow.go
package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// WorkflowType is an enum for the type of workflow
type WorkflowType string

const (
	// ScreenshotWorkflowType is a workflow that takes a screenshot
	ScreenshotWorkflowType WorkflowType = "screenshot"
	// ReportWorkflowType is a workflow that generates a report
	ReportWorkflowType WorkflowType = "report"
	// DispatchWorkflowType is a workflow that dispatches a generated report
	DispatchWorkflowType WorkflowType = "dispatch"
)

// Parses a string into a WorkflowType
func ParseWorkflowType(s string) (WorkflowType, error) {
	switch WorkflowType(s) {
	case ScreenshotWorkflowType, ReportWorkflowType, DispatchWorkflowType:
		return WorkflowType(s), nil
	default:
		return "", fmt.Errorf("invalid workflow type: %s", s)
	}
}

// WorkflowSchedule represents a scheduled workflow
type WorkflowSchedule struct {
	// ID is the unique identifier for the scheduled workflow
	ID ScheduleID `json:"id"`

	// WorkflowType is the type of the workflow
	WorkflowType WorkflowType `json:"workflow_type"`

	// About is the description of the workflow
	About string `json:"about"`

	// Spec is the cron specification for the workflow
	Spec string `json:"spec"`

	// LastRun is the time when the workflow was last run
	// Using sql.NullTime for NULL handling
	LastRun sql.NullTime `json:"last_run,omitempty"`

	// NextRun is the time when the workflow is scheduled to run next
	// Using sql.NullTime for NULL handling
	NextRun sql.NullTime `json:"next_run,omitempty"`

	// CreatedAt is the time when the workflow was created
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt is the time when the workflow was last updated
	UpdatedAt time.Time `json:"updated_at"`
}

type WorkflowScheduleProps struct {
	// WorkflowType is the type of the workflow
	WorkflowType WorkflowType `json:"workflow_type" default:"screenshot" validate:"required,oneof=screenshot report"`

	// About is the description of the workflow
	About string `json:"about" default:""`

	// Spec is the cron specification for the workflow
	Spec string `json:"spec" required:"true" validate:"required"`
}

// WorkflowScheduleJSON is an internal type for JSON marshaling/unmarshaling
type workflowScheduleJSON struct {
	ID           string       `json:"id"`
	WorkflowType WorkflowType `json:"workflow_type"`
	About        string       `json:"about"`
	Spec         string       `json:"spec"`
	LastRun      *time.Time   `json:"last_run,omitempty"`
	NextRun      *time.Time   `json:"next_run,omitempty"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

// MarshalJSON implements custom JSON marshaling for WorkflowSchedule
func (w WorkflowSchedule) MarshalJSON() ([]byte, error) {
	j := workflowScheduleJSON{
		ID:           w.ID.String(),
		WorkflowType: w.WorkflowType,
		About:        w.About,
		Spec:         w.Spec,
		CreatedAt:    w.CreatedAt,
		UpdatedAt:    w.UpdatedAt,
	}

	if w.LastRun.Valid {
		j.LastRun = &w.LastRun.Time
	}

	if w.NextRun.Valid {
		j.NextRun = &w.NextRun.Time
	}

	return json.Marshal(j)
}

// UnmarshalJSON implements custom JSON unmarshaling for WorkflowSchedule
func (w *WorkflowSchedule) UnmarshalJSON(data []byte) error {
	var j workflowScheduleJSON
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}

	// Parse UUID from string
	id, err := uuid.Parse(j.ID)
	if err != nil {
		return fmt.Errorf("invalid schedule ID: %w", err)
	}
	w.ID = ScheduleID(id)

	w.WorkflowType = j.WorkflowType
	w.About = j.About
	w.Spec = j.Spec
	w.CreatedAt = j.CreatedAt
	w.UpdatedAt = j.UpdatedAt

	// Handle nullable times
	if j.LastRun != nil {
		w.LastRun = sql.NullTime{
			Time:  *j.LastRun,
			Valid: true,
		}
	} else {
		w.LastRun = sql.NullTime{Valid: false}
	}

	if j.NextRun != nil {
		w.NextRun = sql.NullTime{
			Time:  *j.NextRun,
			Valid: true,
		}
	} else {
		w.NextRun = sql.NullTime{Valid: false}
	}

	return nil
}
