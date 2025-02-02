package executor

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/service/workspace"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
)

type reportExecutor struct {
	workspaceService workspace.WorkspaceService

	logger *logger.Logger

	runtimeConfig models.JobExecutorConfig
}

func NewReportExecutor(workspaceService workspace.WorkspaceService, logger *logger.Logger, runtimConfig models.JobExecutorConfig) (JobExecutor, error) {
	if workspaceService == nil {
		return nil, errors.New("workspaceService is required")
	}

	if logger == nil {
		return nil, errors.New("logger is required")
	}

	re := reportExecutor{
		workspaceService: workspaceService,
		logger: logger.WithFields(map[string]any{
			"module": "report_executor",
		}),
		runtimeConfig: runtimConfig,
	}

	return &re, nil
}

func (re *reportExecutor) Execute(executionContext context.Context, jobState models.JobState) (<-chan models.JobUpdate, <-chan models.JobError) {
	re.logger.Info("Executing report job", zap.Any("job_state", jobState))

	updates := make(chan models.JobUpdate, 1)
	errors := make(chan models.JobError, 1)

	go func() {
		defer close(updates)
		defer close(errors)

		checkpoint := jobState.Checkpoint.BatchID

		workspaceBatchChan, errBatchChan := re.workspaceService.ListActiveWorkspaces(executionContext, re.runtimeConfig.Parallelism, checkpoint)

		batchStartTime := time.Now()

		for {
			select {
			case workspaceBatch, ok := <-workspaceBatchChan:
				if !ok {
					return
				}

				// Get completion channel for the batch
				completionChan := re.processBatch(executionContext, workspaceBatch, errors)

				// Find the max index of the batch
				maxIndex := -1
				for index := range completionChan {
					if index > maxIndex {
						maxIndex = index
					}
				}

				// Send the update
				if maxIndex >= 0 && maxIndex < len(workspaceBatch) {
					select {
					case updates <- models.JobUpdate{
						Time:      time.Now(),
						Completed: int64(maxIndex + 1),
						Failed:    int64(len(workspaceBatch) - (maxIndex + 1)),
						NewCheckpoint: models.JobCheckpoint{
							BatchID: &workspaceBatch[maxIndex],
						},
					}:
					case <-executionContext.Done():
						return
					}
				}

				// Handle rate limiting
				elapsedTime := time.Since(batchStartTime)
				if remainingTime := re.runtimeConfig.LowerBound - elapsedTime; remainingTime > 0 {
					select {
					case <-time.After(remainingTime):
					case <-executionContext.Done():
						return
					}
				}
				batchStartTime = time.Now()

			case err, ok := <-errBatchChan:
				if !ok {
					return
				}
				select {
				case errors <- models.JobError{
					Error: err,
					Time:  time.Now(),
				}:
				case <-executionContext.Done():
					// Return if the context is done
					return
				}

			case <-executionContext.Done():
				return
			}
		}

	}()

	return updates, errors
}

func (re *reportExecutor) processBatch(ctx context.Context, workspaceBatch []uuid.UUID, errors chan models.JobError) <-chan int {
	completions := make(chan int, len(workspaceBatch))

	// Validate timeout
	if re.runtimeConfig.UpperBound <= 0 {
		re.logger.Error("invalid upper bound timeout",
			zap.Duration("upperBound", re.runtimeConfig.UpperBound))
		close(completions)
		return completions
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, re.runtimeConfig.UpperBound)

	var wg sync.WaitGroup

	for index, workspaceID := range workspaceBatch {
		wg.Add(1)
		go func(workspaceIndex int, workspaceID uuid.UUID) {
			defer wg.Done()

			start := time.Now()
			err := re.processWorkspace(timeoutCtx, workspaceID)
			duration := time.Since(start)

			if err != nil {
				re.logger.Error("workspace processing failed",
					zap.Any("workspaceID", workspaceID),
					zap.Duration("duration", duration),
					zap.Error(err))

				select {
				case errors <- models.JobError{Error: err, Time: time.Now()}:
				case <-timeoutCtx.Done():
				}
				return
			}

			select {
			case completions <- workspaceIndex:
			case <-timeoutCtx.Done():
				re.logger.Error("timeout exceeded for report executor",
					zap.Duration("upperBound", re.runtimeConfig.UpperBound),
					zap.Duration("duration", time.Since(start)),
					zap.Any("workspaceID", workspaceID),
					zap.Error(err))
			}
		}(index, workspaceID)
	}

	// Close completion channel when all work is done
	go func() {
		wg.Wait()
		close(completions)
		cancel() // Clean up timeout context
	}()

	return completions
}

func (re *reportExecutor) processWorkspace(ctx context.Context, workspaceID uuid.UUID) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		// Process the workspace
		competitors, _, err := re.workspaceService.ListCompetitorsForWorkspace(ctx, workspaceID, nil, nil)
		if err != nil {
			return err
		}
		errs := make([]error, 0)
		for _, competitor := range competitors {
			if err := re.processCompetitor(ctx, workspaceID, competitor.ID); err != nil {
				re.logger.Error("failed to process competitor", zap.Any("workspaceID", workspaceID), zap.Any("competitorID", competitor.ID), zap.Error(err))
				errs = append(errs, err)
			}
		}
		if len(errs) > 0 {
			err = fmt.Errorf("failed to process some competitors, %v", errs)
		}
		return err
	}
}

func (re *reportExecutor) processCompetitor(ctx context.Context, workspaceID uuid.UUID, competitorID uuid.UUID) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		// Process the competitor
		_, err := re.workspaceService.CreateReport(ctx, workspaceID, competitorID)
		if err != nil {
			return err
		}
		return nil
	}
}

func (pe *reportExecutor) Terminate(ctx context.Context) error {
	// No cleanup required for the page executor
	// as it does not hold any shared resources
	return nil
}
