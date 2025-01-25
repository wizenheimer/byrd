// ./src/internal/repository/schedule/interface.go
package schedule

import (
	"context"
	"time"

	models "github.com/wizenheimer/byrd/src/internal/models/core"
)

type ScheduleRepository interface {
	// CreateSchedule schedules a new workflow in the repository
	CreateSchedule(ctx context.Context, workflowProps models.WorkflowScheduleProps) (models.ScheduleID, error)

	// CreateScheduleWithID schedules a new workflow in the repository with a specific ID
	CreateScheduleWithID(ctx context.Context, scheduleID models.ScheduleID, workflowProps models.WorkflowScheduleProps) (models.ScheduleID, error)

	// GetSchedule returns the schedule of a workflow
	GetSchedule(ctx context.Context, scheduleID models.ScheduleID) (models.WorkflowSchedule, error)

	// UpdateSchedule updates the schedule of a workflow in the repository
	UpdateSchedule(ctx context.Context, scheduleID models.ScheduleID, workflowProps models.WorkflowScheduleProps) error

	// DeleteSchedule deletes the scheduled workflows
	DeleteSchedule(ctx context.Context, scheduleID models.ScheduleID) error

	// ListScheduledWorkflows returns the list of scheduled workflows
	ListScheduledWorkflows(ctx context.Context, limit, offset *int, workflowType *models.WorkflowType) ([]models.WorkflowSchedule, error)

	// Sync
	Sync(ctx context.Context, scheduleID models.ScheduleID, lastRun, nextRun time.Time) error
}
