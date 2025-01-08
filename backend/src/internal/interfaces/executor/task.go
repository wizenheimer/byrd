package executor

import (
	"context"

	models "github.com/wizenheimer/iris/src/internal/models/core"
)

// TaskExecutor represents the executor for managing tasks
type TaskExecutor interface {
	// Execute executes the task
	Execute(ctx context.Context, t models.Task) (<-chan models.TaskUpdate, <-chan models.TaskError)

	// Terminate terminates the task
	Terminate(ctx context.Context) error
}
