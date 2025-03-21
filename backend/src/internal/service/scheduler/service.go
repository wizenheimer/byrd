// ./src/internal/service/scheduler/service.go
package scheduler

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/robfig/cron/v3"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/recorder"
	"github.com/wizenheimer/byrd/src/internal/repository/schedule"
	"github.com/wizenheimer/byrd/src/internal/scheduler"
	"github.com/wizenheimer/byrd/src/internal/service/workflow"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"github.com/wizenheimer/byrd/src/pkg/utils"
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

	// errorRecord is the error recorder for the scheduler service
	errorRecord *recorder.ErrorRecorder

	// parser is the cron parser
	parser cron.Parser
}

// NewSchedulerService creates a new scheduler service
func NewSchedulerService(
	repository schedule.ScheduleRepository,
	scheduler scheduler.Scheduler,
	workflowService workflow.WorkflowService,
	logger *logger.Logger,
	errorRecord *recorder.ErrorRecorder,
) (SchedulerService, error) {
	parser := utils.NewScheduleParser()

	s := schedulerService{
		repository: repository,
		logger: logger.WithFields(map[string]any{
			"module": "scheduler_service",
		}),
		errorRecord:     errorRecord,
		scheduler:       scheduler,
		workflowService: workflowService,
		parser:          parser,
	}
	return &s, nil
}

// Start starts the scheduler service
func (s *schedulerService) Start(ctx context.Context, recovery bool) error {
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
	err := s.scheduler.Stop()
	if err != nil {
		s.errorRecord.RecordError(ctx, fmt.Errorf("failed to stop scheduler %v", err.Error()))
	}
	return err
}

func (s *schedulerService) triggerWorkflow(ctx context.Context, workflowType models.WorkflowType) func() {
	return func() {
		// Execute the workflow
		_, err := s.workflowService.Submit(ctx, workflowType)
		if err != nil {
			s.errorRecord.RecordError(ctx, fmt.Errorf("failed to submit workflow %v", err.Error()), zap.Any("workflowType", workflowType))
		}
	}
}

func (s *schedulerService) syncWorkflow(ctx context.Context, remoteScheduleID models.ScheduleID) func() {
	return func() {
		// Get the scheduled function
		v, ok := s.scheduledFuncs.Load(remoteScheduleID)
		if !ok {
			s.logger.Error("scheduled function not found", zap.Any("remoteScheduleID", remoteScheduleID))
			return
		}
		sf, ok := v.(*models.ScheduledFunc)
		if !ok {
			s.logger.Error("failed to cast scheduled function", zap.Any("remoteScheduleID", remoteScheduleID))
			return
		}

		// Sync the workflow times
		if err := s.repository.Sync(ctx, remoteScheduleID, sf.LastRun, sf.NextRun); err != nil {
			s.errorRecord.RecordError(ctx, err, zap.Any("remoteScheduleID", remoteScheduleID))
		}
	}
}

// Schedule schedules a new workflow
func (s *schedulerService) Schedule(ctx context.Context, workflowProp models.WorkflowScheduleProps) (models.ScheduleID, error) {
	// Add validation to the workflowProp
	_, err := s.parser.Parse(workflowProp.Spec)
	if err != nil {
		return models.NilScheduleID(), err
	}

	remoteScheduleID := models.NewScheduleID()

	op := &CreateScheduleOperation{
		svc:          s,
		remoteID:     remoteScheduleID,
		workflowProp: workflowProp,
	}

	if err := op.Execute(ctx); err != nil {
		s.errorRecord.RecordError(ctx, fmt.Errorf("failed to schedule workflow, %s", err.Error()), zap.Any("remoteScheduleID", remoteScheduleID))
		return models.NilScheduleID(), err
	}

	return remoteScheduleID, nil
}

// Unschedule unschedules a workflow
func (s *schedulerService) Unschedule(ctx context.Context, remoteScheduleID models.ScheduleID) error {
	op := &DeleteScheduleOperation{
		svc:      s,
		remoteID: remoteScheduleID,
	}

	err := op.Execute(ctx)
	if err != nil {
		s.errorRecord.RecordError(ctx, fmt.Errorf("failed to unschedule workflow, %s", err.Error()), zap.Any("remoteScheduleID", remoteScheduleID))
		return err
	}

	return nil
}

// Reschedule reschedules a workflow
func (s *schedulerService) Reschedule(ctx context.Context, remoteScheduleID models.ScheduleID, workflowProp models.WorkflowScheduleProps) (models.ScheduleID, error) {
	op := &UpdateScheduleOperation{
		svc:          s,
		remoteID:     remoteScheduleID,
		workflowProp: workflowProp,
	}

	if err := op.Execute(ctx); err != nil {
		s.errorRecord.RecordError(ctx, fmt.Errorf("failed to reschedule workflow, %s", err.Error()), zap.Any("remoteScheduleID", remoteScheduleID))
		return models.NilScheduleID(), err
	}

	return remoteScheduleID, nil
}

// Get returns the schedule of a workflow
func (s *schedulerService) Get(ctx context.Context, remoteScheduleID models.ScheduleID) (*models.WorkflowSchedule, error) {
	// Get schedule of a workflow from local state
	v, ok := s.scheduledFuncs.Load(remoteScheduleID)
	if !ok {
		return nil, errors.New("scheduled function not found")
	}
	f, ok := v.(*models.ScheduledFunc)
	if !ok {
		return nil, errors.New("failed to cast scheduled function")
	}

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
	// List the scheduled workflows
	workflowSchedules, err := s.repository.ListScheduledWorkflows(ctx, limit, offset, workflowType)
	if err != nil {
		return nil, err
	}

	return workflowSchedules, nil
}

func (s *schedulerService) Recover(ctx context.Context) {
	// Recover scheduled functions
	workflows, err := s.List(ctx, nil, nil, nil)
	if err != nil {
		s.errorRecord.RecordError(ctx, fmt.Errorf("failed to list scheduled workflows, %s", err.Error()))
		return
	}

	for _, workflow := range workflows {
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
			s.errorRecord.RecordError(ctx, fmt.Errorf("failed to schedule workflow, %s", err.Error()), zap.Any("scheduleID", scheduleID), zap.Any("workflowID", workflow.ID))
			continue
		}

		// Associate the remote schedule ID with scheduled function
		s.scheduledFuncs.Store(scheduleID, f)
	}
}
