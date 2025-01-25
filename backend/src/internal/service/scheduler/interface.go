// ./src/internal/service/scheduler/interface.go
package scheduler

import (
	"context"

	models "github.com/wizenheimer/byrd/src/internal/models/core"
	// "github.com/wizenheimer/byrd/src/internal/scheduler"
)

type SchedulerService interface {
	// Start starts the scheduler service
	Start(ctx context.Context, recovery bool) error

	// Gracefully stops the scheduler service gracefully
	Stop(ctx context.Context) error

	// Schedule schedules a new workflow
	Schedule(ctx context.Context, workflowProp models.WorkflowScheduleProps) (models.ScheduleID, error)

	// Unschedule unschedules a workflow
	Unschedule(ctx context.Context, remoteScheduleID models.ScheduleID) error

	// Reschedule reschedules a workflow
	Reschedule(ctx context.Context, remoteScheduleID models.ScheduleID, workflowProp models.WorkflowScheduleProps) (models.ScheduleID, error)

	// Get returns the schedule of a workflow
	Get(ctx context.Context, remoteScheduleID models.ScheduleID) (*models.WorkflowSchedule, error)

	// List returns the list of scheduled workflows
	List(ctx context.Context, limit, offset *int, workflowType *models.WorkflowType) ([]models.WorkflowSchedule, error)
}
