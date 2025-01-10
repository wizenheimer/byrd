// ./src/internal/models/core/task.go
package models

import "time"

type TaskStatus string

const (
	TaskStatusPending  TaskStatus = "pending"
	TaskStatusRunning  TaskStatus = "running"
	TaskStatusComplete TaskStatus = "complete"
	TaskStatusFailed   TaskStatus = "failed"
)

// Task represents a task that needs to be executed
type Task struct {
	TaskID     string             `json:"task_id"`
	WorkflowID WorkflowIdentifier `json:"workflow_id"`
	Checkpoint WorkflowCheckpoint `json:"checkpoint"`
}

// TaskError represents an error that occurred while executing a task
type TaskError struct {
	TaskID string    `json:"task_id"`
	Error  error     `json:"error"`
	Time   time.Time `json:"time"`
}

// TaskUpdate represents an update to a task
type TaskUpdate struct {
	TaskID        string             `json:"task_id"`
	Status        TaskStatus         `json:"status"`
	Completed     int                `json:"completed"`
	Failed        int                `json:"failed"`
	NewCheckpoint WorkflowCheckpoint `json:"new_checkpoint"`
}
