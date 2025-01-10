// ./src/internal/models/api/workflow.go
package models

import (
	"fmt"
	"time"

	models "github.com/wizenheimer/byrd/src/internal/models/core"
)

// WorkflowRequest is a request to start a workflow
type WorkflowRequest struct {
	// Type is the type of the workflow
	Type *models.WorkflowType `json:"workflow_type"`
	// Year is the year of the workflow
	Year *int `json:"year"`
	// WeekNumber is the week number of the workflow
	WeekNumber *int `json:"week_number"`
	// WeekDay is the week day of the workflow
	WeekDay *int `json:"week_day"`
}

// Validate checks if the WorkflowRequest is valid
// If safe is true, it will set default values for missing or invalid fields
func (wr WorkflowRequest) Validate(safe bool) error {
	if wr.Type == nil {
		return fmt.Errorf("missing workflow type, got %v", wr.Type)
	} else if *wr.Type != models.ScreenshotWorkflowType && *wr.Type != models.ReportWorkflowType {
		if safe {
			*wr.Type = models.ScreenshotWorkflowType
		} else {
			return fmt.Errorf("invalid workflow type: %v", *wr.Type)
		}
	}

	if wr.Year == nil {
		return fmt.Errorf("missing year, got %v", wr.Year)
	} else if *wr.Year < 2000 || *wr.Year > 2100 {
		if safe {
			currentYear := time.Now().Year()
			*wr.Year = currentYear
		} else {
			return fmt.Errorf("invalid year: %v", *wr.Year)
		}
	}

	if wr.WeekNumber == nil {
		return fmt.Errorf("missing week number, got %v", wr.WeekNumber)
	} else if *wr.WeekNumber < 1 || *wr.WeekNumber > 53 {
		if safe {
			_, currentWeek := time.Now().ISOWeek()
			*wr.WeekNumber = currentWeek
		} else {
			return fmt.Errorf("invalid week number: %v", *wr.WeekNumber)
		}
	}

	if wr.WeekDay == nil {
		return fmt.Errorf("missing week day, got %v", wr.WeekDay)
	} else if *wr.WeekDay < 1 || *wr.WeekDay > 7 {
		if safe {
			_, _, currentWeekDay := time.Now().Date()
			*wr.WeekDay = int(currentWeekDay)
		} else {
			return fmt.Errorf("invalid week day: %v", *wr.WeekDay)
		}
	}

	return nil
}

// WorkflowResponse is a response to a workflow request
type WorkflowResponse struct {
	// WorkflowID is the identifier of the workflow
	WorkflowID models.WorkflowIdentifier `json:"workflow_id"`
	// WorkflowState is the state of the workflow
	WorkflowState WorkflowState `json:"workflow_state"`
}

// WorkflowState captures the state of the workflow
type WorkflowState struct {
	// Status is the current status of the workflow
	Status models.WorkflowStatus `json:"status"`
	// Checkpoint is the current checkpoint of the workflow
	Checkpoint models.WorkflowCheckpoint `json:"checkpoint"`
}
