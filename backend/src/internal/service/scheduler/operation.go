// ./src/internal/service/scheduler/operation.go
package scheduler

import (
	"context"
	"errors"
	"fmt"

	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/scheduler"
	"go.uber.org/zap"
)

// RollbackOperation represents a reversible operation
type RollbackOperation interface {
	// Execute performs the operation
	Execute(ctx context.Context) error
	// Rollback reverts the operation
	Rollback(ctx context.Context) error
}

// CompoundOperation represents a series of operations that should be executed atomically
type CompoundOperation struct {
	operations []RollbackOperation
	executed   []RollbackOperation
}

// NewCompoundOperation creates a new compound operation
func NewCompoundOperation() *CompoundOperation {
	return &CompoundOperation{
		operations: make([]RollbackOperation, 0),
		executed:   make([]RollbackOperation, 0),
	}
}

// AddOperation adds an operation to the compound operation
func (co *CompoundOperation) AddOperation(op RollbackOperation) {
	co.operations = append(co.operations, op)
}

// Execute executes all operations in order, rolling back on failure
func (co *CompoundOperation) Execute(ctx context.Context) error {
	for _, op := range co.operations {
		if err := op.Execute(ctx); err != nil {
			// Rollback all executed operations in reverse order
			for i := len(co.executed) - 1; i >= 0; i-- {
				if rbErr := co.executed[i].Rollback(ctx); rbErr != nil {
					return fmt.Errorf("rollback failed: %v (original error: %v)", rbErr, err)
				}
			}
			return err
		}
		co.executed = append(co.executed, op)
	}
	return nil
}

// CreateScheduleOperation represents the operation of creating a schedule
type CreateScheduleOperation struct {
	svc           *schedulerService
	remoteID      models.ScheduleID
	workflowProp  models.WorkflowScheduleProps
	scheduledFunc *models.ScheduledFunc
}

func (op *CreateScheduleOperation) Execute(ctx context.Context) error {
	// Schedule in repository first
	_, err := op.svc.repository.CreateScheduleWithID(ctx, op.remoteID, op.workflowProp)
	if err != nil {
		return err
	}

	// Schedule in scheduler
	opts := scheduler.ScheduleOptions{
		Hooks:        []func(){op.svc.syncWorkflow(ctx, op.remoteID)},
		ScheduleSpec: op.workflowProp.Spec,
	}
	cmd := op.svc.triggerWorkflow(ctx, op.workflowProp.WorkflowType)
	f, err := op.svc.scheduler.Schedule(cmd, opts)
	if err != nil {
		return err
	}

	op.scheduledFunc = f
	op.svc.scheduledFuncs.Store(op.remoteID, f)
	return nil
}

func (op *CreateScheduleOperation) Rollback(ctx context.Context) error {
	logger := op.svc.logger
	// Delete from scheduler if needed
	if op.scheduledFunc != nil {
		if err := op.svc.scheduler.Delete(op.scheduledFunc.ID); err != nil {
			logger.Fatal("failed to rollback create operation", zap.Error(err))
			return err
		}
		op.svc.scheduledFuncs.Delete(op.remoteID)
	}

	// Delete from repository
	if err := op.svc.repository.DeleteSchedule(ctx, op.remoteID); err != nil {
		logger.Fatal("failed to rollback create operation", zap.Error(err))
		return err
	}

	return nil
}

// UpdateScheduleOperation represents the operation of updating a schedule
type UpdateScheduleOperation struct {
	svc             *schedulerService
	remoteID        models.ScheduleID
	workflowProp    models.WorkflowScheduleProps
	oldWorkflowProp models.WorkflowScheduleProps
	oldFunc         *models.ScheduledFunc
}

func (op *UpdateScheduleOperation) Execute(ctx context.Context) error {
	// Store old state for potential rollback
	schedule, err := op.svc.repository.GetSchedule(ctx, op.remoteID)
	if err != nil {
		return err
	}
	op.oldWorkflowProp = models.WorkflowScheduleProps{
		WorkflowType: schedule.WorkflowType,
		Spec:         schedule.Spec,
	}

	v, ok := op.svc.scheduledFuncs.Load(op.remoteID)
	if !ok {
		return fmt.Errorf("scheduled function not found")
	}
	op.oldFunc, ok = v.(*models.ScheduledFunc)
	if !ok {
		return fmt.Errorf("failed to cast scheduled function")
	}

	// Update repository
	if err := op.svc.repository.UpdateSchedule(ctx, op.remoteID, op.workflowProp); err != nil {
		return err
	}

	// Update scheduler
	opts := scheduler.ScheduleOptions{
		Hooks:        []func(){op.svc.syncWorkflow(ctx, op.remoteID)},
		ScheduleSpec: op.workflowProp.Spec,
	}
	f, err := op.svc.scheduler.Update(op.oldFunc.ID, op.svc.triggerWorkflow(ctx, op.workflowProp.WorkflowType), opts)
	if err != nil {
		return err
	}

	op.svc.scheduledFuncs.Store(op.remoteID, f)
	return nil
}

func (op *UpdateScheduleOperation) Rollback(ctx context.Context) error {
	logger := op.svc.logger
	// Restore old state in repository
	if err := op.svc.repository.UpdateSchedule(ctx, op.remoteID, op.oldWorkflowProp); err != nil {
		logger.Fatal("failed to rollback update operation", zap.Error(err))
		return err
	}

	// Restore old state in scheduler
	opts := scheduler.ScheduleOptions{
		Hooks:        []func(){op.svc.syncWorkflow(ctx, op.remoteID)},
		ScheduleSpec: op.oldWorkflowProp.Spec,
	}
	cmd := op.svc.triggerWorkflow(ctx, op.oldWorkflowProp.WorkflowType)
	f, err := op.svc.scheduler.Update(op.oldFunc.ID, cmd, opts)
	if err != nil {
		logger.Fatal("failed to rollback update operation", zap.Error(err))
		return err
	}

	op.svc.scheduledFuncs.Store(op.remoteID, f)
	return nil
}

// DeleteScheduleOperation represents the operation of deleting a schedule
type DeleteScheduleOperation struct {
	svc             *schedulerService
	remoteID        models.ScheduleID
	oldWorkflowProp models.WorkflowScheduleProps
	oldFunc         *models.ScheduledFunc
	wasDeleted      bool
}

func (op *DeleteScheduleOperation) Execute(ctx context.Context) error {
	logger := op.svc.logger

	// Store old state for potential rollback
	schedule, err := op.svc.repository.GetSchedule(ctx, op.remoteID)
	if err != nil {
		return err
	}
	op.oldWorkflowProp = models.WorkflowScheduleProps{
		WorkflowType: schedule.WorkflowType,
		Spec:         schedule.Spec,
	}

	value, ok := op.svc.scheduledFuncs.Load(op.remoteID)
	if !ok {
		return errors.New("scheduled function not found")
	}
	op.oldFunc, ok = value.(*models.ScheduledFunc)
	if !ok {
		return errors.New("failed to cast scheduled function")
	}

	// Delete from repository first
	if err := op.svc.repository.DeleteSchedule(ctx, op.remoteID); err != nil {
		return err
	}

	// Delete from scheduler
	if err := op.svc.scheduler.Delete(op.oldFunc.ID); err != nil {
		// Rollback repository deletion
		if rbErr := op.Rollback(ctx); rbErr != nil {
			logger.Fatal("failed to rollback repository deletion", zap.Error(rbErr))
			return fmt.Errorf("failed to rollback repository deletion: %v (original error: %v)", rbErr, err)
		}
		return err
	}

	// Remove from local state
	op.svc.scheduledFuncs.Delete(op.remoteID)
	op.wasDeleted = true

	return nil
}

func (op *DeleteScheduleOperation) Rollback(ctx context.Context) error {
	logger := op.svc.logger

	if !op.wasDeleted {
		// Nothing was deleted, no need to rollback
		logger.Debug("nothing to rollback for deletion operation, skipping rollback operation", zap.Any("scheduleID", op.remoteID))
		return nil
	}

	// Restore in repository
	_, err := op.svc.repository.CreateScheduleWithID(ctx, op.remoteID, op.oldWorkflowProp)
	if err != nil {
		logger.Fatal("failed to restore repository state", zap.Error(err))
		return fmt.Errorf("failed to restore repository state: %v", err)
	}

	// Restore in scheduler
	opts := scheduler.ScheduleOptions{
		Hooks:        []func(){op.svc.syncWorkflow(ctx, op.remoteID)},
		ScheduleSpec: op.oldWorkflowProp.Spec,
	}
	f, err := op.svc.scheduler.Schedule(op.svc.triggerWorkflow(ctx, op.oldWorkflowProp.WorkflowType), opts)
	if err != nil {
		logger.Fatal("failed to restore scheduler state", zap.Error(err))
		return fmt.Errorf("failed to restore scheduler state: %v", err)
	}

	// Restore local state
	op.svc.scheduledFuncs.Store(op.remoteID, f)
	op.wasDeleted = false
	return nil
}
