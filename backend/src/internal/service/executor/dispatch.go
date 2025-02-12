package executor

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	slackworkspace "github.com/wizenheimer/byrd/src/internal/service/integration/slack"
	"github.com/wizenheimer/byrd/src/internal/service/workspace"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
)

type dispatchExecutor struct {
	ws             workspace.WorkspaceService
	slackWorkspace slackworkspace.SlackWorkspaceService
	logger         *logger.Logger
	runtimeConfig  models.JobExecutorConfig
}

func NewDispatchExecutor(
	ws workspace.WorkspaceService,
	slackWorkspace slackworkspace.SlackWorkspaceService,
	logger *logger.Logger,
	runtimeConfig models.JobExecutorConfig,
) (JobExecutor, error) {
	if ws == nil {
		return nil, errors.New("workspace service is required")
	}
	if logger == nil {
		return nil, errors.New("logger is required")
	}
	d := dispatchExecutor{
		ws:             ws,
		slackWorkspace: slackWorkspace,
		logger:         logger,
		runtimeConfig:  runtimeConfig,
	}
	return &d, nil
}

func (e *dispatchExecutor) Execute(executionContext context.Context, jobState models.JobState) (<-chan models.JobUpdate, <-chan models.JobError) {
	updates := make(chan models.JobUpdate, 1)
	errors := make(chan models.JobError, 1)

	go func() {
		defer close(updates)
		defer close(errors)

		checkpoint := jobState.Checkpoint.BatchID

		workspaceBatchChan, errBatchChan := e.ws.ListActiveWorkspaces(executionContext, e.runtimeConfig.Parallelism, checkpoint)

		batchStartTime := time.Now()

		for {
			select {
			case workspaceBatch, ok := <-workspaceBatchChan:
				if !ok {
					return
				}

				// Get completion channel for the batch
				completionChan := e.processBatch(executionContext, workspaceBatch, errors)

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
				if remainingTime := e.runtimeConfig.LowerBound - elapsedTime; remainingTime > 0 {
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

func (e *dispatchExecutor) processBatch(ctx context.Context, workspaceBatch []uuid.UUID, errors chan models.JobError) <-chan int {
	completions := make(chan int, len(workspaceBatch))

	// Validate timeout
	if e.runtimeConfig.UpperBound <= 0 {
		e.logger.Error("invalid upper bound timeout",
			zap.Duration("upperBound", e.runtimeConfig.UpperBound))
		close(completions)
		return completions
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, e.runtimeConfig.UpperBound)

	var wg sync.WaitGroup

	for index, workspaceID := range workspaceBatch {
		wg.Add(1)
		go func(workspaceIndex int, workspaceID uuid.UUID) {
			defer wg.Done()

			start := time.Now()
			err := e.processWorkspace(timeoutCtx, workspaceID)
			duration := time.Since(start)

			if err != nil {
				e.logger.Error("workspace processing failed",
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

func (e *dispatchExecutor) processWorkspace(ctx context.Context, workspaceID uuid.UUID) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		// Process the workspace
		competitors, _, err := e.ws.ListCompetitorsForWorkspace(ctx, workspaceID, nil, nil)
		if err != nil {
			return err
		}
		errs := make([]error, 0)
		for _, competitor := range competitors {
			if err := e.processCompetitor(ctx, workspaceID, competitor.ID); err != nil {
				e.logger.Error("failed to process competitor", zap.Any("workspaceID", workspaceID), zap.Any("competitorID", competitor.ID), zap.Error(err))
				errs = append(errs, err)
			}
		}
		if len(errs) > 0 {
			err = fmt.Errorf("failed to process some competitors %v", errs)
		}
		return err
	}
}

func (e *dispatchExecutor) processCompetitor(ctx context.Context, workspaceID uuid.UUID, competitorID uuid.UUID) error {
	slackWorkspaceExists, err := e.slackWorkspace.IntegrationExistsForWorkspace(ctx, workspaceID)
	if err != nil {
		e.logger.Error("failed to check if slack workspace exists", zap.Any("workspaceID", workspaceID), zap.Error(err))
		slackWorkspaceExists = false
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		// Send report to slack workspace
		var slackError error
		if slackWorkspaceExists {
			slackError = e.slackWorkspace.DispatchReportToWorkspaceMembers(ctx, workspaceID, competitorID)
		}

		// Process the competitor
		var emailError error
		emailError = e.ws.DispatchReportToWorkspaceMembers(ctx, workspaceID, competitorID)

		if slackError != nil && emailError != nil {
			return fmt.Errorf("failed to send report to slack workspace and email: %v, %v", slackError, emailError)
		} else if slackError != nil {
			return fmt.Errorf("failed to send report to slack workspace: %v", slackError)
		} else if emailError != nil {
			return fmt.Errorf("failed to send report to email: %v", emailError)
		}
		return nil
	}
}

func (e *dispatchExecutor) Terminate(ctx context.Context) error {
	return nil
}
