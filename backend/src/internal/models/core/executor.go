// ./src/internal/models/core/executor.go
package models

import (
	"fmt"
	"time"
)

// ExecutorConfig represents the configuration for an executor
type ExecutorConfig struct {
	// Number of tasks to execute in parallel
	// This is used to limit the number of tasks that can be executed concurrently
	Parallelism int `json:"batch_size"`
	// Lower bound for the time to wait before executing the next batch
	// This is used to prevent the executor from executing tasks too frequently
	LowerBound time.Duration `json:"lower_bound"`
	// Upper bound for the time to wait before executing the next batch
	// This is used to prevent the executor from getting stuck with the same batch
	UpperBound time.Duration `json:"upper_bound"`
}

// GetExecutorConfig returns the configuration for the given workflow type
// This is used to determine the configuration for the executor based on the workflow type
func GetExecutorConfig(workflowType WorkflowType) (ExecutorConfig, error) {
	// Default executor configuration
	exConfig := ExecutorConfig{
		Parallelism: 10,
		LowerBound:  1 * time.Minute,
		UpperBound:  30 * time.Minute,
	}

	switch workflowType {
	case ScreenshotWorkflowType:
		return exConfig, nil
	case ReportWorkflowType:
		return exConfig, nil
	}

	return exConfig, fmt.Errorf("invalid workflow type: %s", workflowType)
}
