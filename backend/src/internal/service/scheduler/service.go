package scheduler

import (
	"context"
	"errors"
	"sync"

	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/repository/schedule"
	"github.com/wizenheimer/byrd/src/internal/scheduler"
	"github.com/wizenheimer/byrd/src/internal/service/workflow"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
)

type schedulerService struct {
	// repository is the repository for the scheduler service
	repository schedule.ScheduleRepository

	// logger is the logger for the scheduler service
	logger *logger.Logger

	// scheduler is the scheduler for the scheduler service
	scheduler scheduler.Scheduler

	// workflowService is the workflow service for the scheduler service
	workflowService workflow.WorkflowService

	// scheduledFuncs is a map of scheduled functions
	scheduledFuncs sync.Map
}

// NewSchedulerService creates a new scheduler service
func NewSchedulerService(repository schedule.ScheduleRepository, scheduler scheduler.Scheduler, workflowService workflow.WorkflowService, logger *logger.Logger) SchedulerService {
	return &schedulerService{
		repository:      repository,
		logger:          logger.WithFields(map[string]interface{}{"module": "scheduler_service"}),
		scheduler:       scheduler,
		workflowService: workflowService,
	}
}

// Start starts the scheduler service
func (s *schedulerService) Start(ctx context.Context, recovery bool) error {
	s.logger.Info("starting the scheduler service")
	err := s.scheduler.Start()
	if err != nil {
		return err
	}

	if recovery {
		s.Recover(ctx)
	}

	return nil
}

// Gracefully stops the scheduler service gracefully
func (s *schedulerService) Stop(ctx context.Context) error {
	s.logger.Info("stopping the scheduler service")
	return s.scheduler.Stop()
}

func (s *schedulerService) triggerWorkflow(ctx context.Context, workflowType models.WorkflowType) func() {
	return func() {
		s.logger.Debug("triggering workflow", zap.Any("workflow_type", workflowType), zap.Any("workflowService", s.workflowService))
		// Execute the workflow
		jobID, err := s.workflowService.Submit(ctx, workflowType)
		if err != nil {
			s.logger.Error("failed to submit workflow", zap.Error(err))
		} else {
			s.logger.Info("successfully submitted workflow", zap.Any("job_id", jobID))
		}
	}
}

func (s *schedulerService) syncWorkflow(ctx context.Context, remoteScheduleID models.ScheduleID) func() {
	return func() {
		// Get the scheduled function
		v, ok := s.scheduledFuncs.Load(remoteScheduleID)
		if !ok {
			s.logger.Error("scheduled function not found")
			return
		}
		sf := v.(*models.ScheduledFunc)

		// Sync the workflow times
		if err := s.repository.Sync(ctx, remoteScheduleID, sf.LastRun, sf.NextRun); err != nil {
			s.logger.Error("failed to sync workflow times", zap.Error(err))
		}
	}
}

// Schedule schedules a new workflow
func (s *schedulerService) Schedule(ctx context.Context, workflowProp models.WorkflowScheduleProps) (models.ScheduleID, error) {
	s.logger.Info("scheduling a new workflow")

	// Persist the workflow schedule
	remoteScheduleID, err := s.repository.CreateSchedule(ctx, workflowProp)
	if err != nil {
		return models.NilScheduleID(), err
	}

	// Schedule options
	opts := scheduler.ScheduleOptions{
		Hooks:        []func(){s.syncWorkflow(ctx, remoteScheduleID)},
		ScheduleSpec: workflowProp.Spec,
	}
	// Scheduled command
	cmd := s.triggerWorkflow(ctx, workflowProp.WorkflowType)

	// Trigger the scheduler to schedule the workflow
	f, err := s.scheduler.Schedule(cmd, opts)
	if err != nil {
		return models.NilScheduleID(), err
	}

	// Associate the remote schedule ID with scheduled function
	s.scheduledFuncs.Store(remoteScheduleID, f)
	s.logger.Debug("scheduled function", zap.Any("scheduled_func", f))
	return remoteScheduleID, nil
}

// Unschedule unschedules a workflow
func (s *schedulerService) Unschedule(ctx context.Context, remoteScheduleID models.ScheduleID) error {
	s.logger.Info("unscheduling a workflow")

	// Delete the workflow schedule from the repository
	if err := s.repository.DeleteSchedule(ctx, remoteScheduleID); err != nil {
		return err
	}

	value, ok := s.scheduledFuncs.Load(remoteScheduleID)
	if !ok {
		return errors.New("scheduled function not found")
	}

	f := value.(*models.ScheduledFunc)

	// Remove the scheduled function from the scheduler
	if err := s.scheduler.Delete(f.ID); err != nil {
		return err
	}

	// Remove the scheduled function from the local state
	s.scheduledFuncs.Delete(remoteScheduleID)
	return nil
}

// Reschedule reschedules a workflow
func (s *schedulerService) Reschedule(ctx context.Context, remoteScheduleID models.ScheduleID, workflowProp models.WorkflowScheduleProps) (models.ScheduleID, error) {
	s.logger.Info("rescheduling a workflow")
	// Update the workflow schedule
	if err := s.repository.UpdateSchedule(ctx, remoteScheduleID, workflowProp); err != nil {
		return models.NilScheduleID(), err
	}

	v, ok := s.scheduledFuncs.Load(remoteScheduleID)
	if !ok {
		return models.NilScheduleID(), errors.New("scheduled function not found")
	}
	f := v.(*models.ScheduledFunc)

	// Remove the scheduled function from the local state
	s.scheduledFuncs.Delete(remoteScheduleID)

	// Schedule options
	opts := scheduler.ScheduleOptions{
		Hooks:        []func(){s.syncWorkflow(ctx, remoteScheduleID)},
		ScheduleSpec: workflowProp.Spec,
	}
	// Scheduled command
	cmd := s.triggerWorkflow(ctx, workflowProp.WorkflowType)

	// Update the scheduled function
	f, err := s.scheduler.Update(f.ID, cmd, opts)
	if err != nil {
		return models.NilScheduleID(), err
	}

	// Store the scheduled function
	s.scheduledFuncs.Store(remoteScheduleID, f)
	return remoteScheduleID, nil
}

// Get returns the schedule of a workflow
func (s *schedulerService) Get(ctx context.Context, remoteScheduleID models.ScheduleID) (*models.WorkflowSchedule, error) {
	s.logger.Info("getting the schedule of a workflow")
	// Get schedule of a workflow from local state
	v, ok := s.scheduledFuncs.Load(remoteScheduleID)
	if !ok {
		return nil, errors.New("scheduled function not found")
	}
	f := v.(*models.ScheduledFunc)

	_, err := s.scheduler.Get(f.ID)
	if err != nil {
		return nil, err
	}

	// Get the workflow schedule
	workflowSchedule, err := s.repository.GetSchedule(ctx, remoteScheduleID)
	if err != nil {
		return nil, err
	}

	return &workflowSchedule, nil
}

// List returns the list of scheduled workflows
func (s *schedulerService) List(ctx context.Context, limit, offset *int, workflowType *models.WorkflowType) ([]models.WorkflowSchedule, error) {
	s.logger.Info("listing the scheduled workflows")
	// List the scheduled workflows
	workflowSchedules, err := s.repository.ListScheduledWorkflows(ctx, limit, offset, workflowType)
	if err != nil {
		return nil, err
	}

	return workflowSchedules, nil
}

func (s *schedulerService) Recover(ctx context.Context) {
	s.logger.Info("recovering scheduled functions")
	// Recover scheduled functions
	workflows, err := s.List(ctx, nil, nil, nil)
	if err != nil {
		s.logger.Error("failed to list scheduled workflows", zap.Error(err))
		return
	}

	for _, workflow := range workflows {
		s.logger.Info("recovering scheduled function", zap.Any("workflow", workflow))
		scheduleID := workflow.ID
		// Schedule options
		opts := scheduler.ScheduleOptions{
			Hooks:        []func(){s.syncWorkflow(ctx, scheduleID)},
			ScheduleSpec: workflow.Spec,
		}
		// Scheduled command
		cmd := s.triggerWorkflow(ctx, workflow.WorkflowType)

		// Trigger the scheduler to schedule the workflow
		f, err := s.scheduler.Schedule(cmd, opts)
		if err != nil {
			s.logger.Error("failed to schedule workflow", zap.Error(err))
			continue
		}

		// Associate the remote schedule ID with scheduled function
		s.scheduledFuncs.Store(scheduleID, f)
	}
}
