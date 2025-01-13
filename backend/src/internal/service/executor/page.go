package executor

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/service/page"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
)

type pageExecutor struct {
	// pageService represents the page service for the workflow
	pageService page.PageService

	// logger represents the logger for the workflow
	logger *logger.Logger

	// runtimeConfig represents the runtime configuration for the workflow
	runtimeConfig models.JobExecutorConfig
}

func NewPageExecutor(pageService page.PageService, runtimeConfig models.JobExecutorConfig, logger *logger.Logger) (JobExecutor, error) {
	logger.Debug("initializing page executor", zap.Any("runtimeConfig", runtimeConfig))
	if logger == nil {
		return nil, errors.New("logger is required")
	}

	pe := pageExecutor{
		pageService:   pageService,
		runtimeConfig: runtimeConfig,
		logger:        logger.WithFields(map[string]interface{}{"module": "page_executor"}),
	}

	return &pe, nil
}

func (pe *pageExecutor) Execute(executionContext context.Context, jobState models.JobState) (<-chan models.JobUpdate, <-chan models.JobError) {
	pe.logger.Debug("executing page executor", zap.Any("jobState", jobState))
	updates := make(chan models.JobUpdate, 1)
	errors := make(chan models.JobError, 1)

	go func() {
		defer close(updates)
		defer close(errors)

		checkpoint := jobState.Checkpoint.BatchID
		pageBatchChan, errBatchChan := pe.pageService.ListActivePages(executionContext, pe.runtimeConfig.Parallelism, checkpoint)

		batchStartTime := time.Now()

		for {
			select {
			case pageBatch, ok := <-pageBatchChan:
				if !ok {
					return
				}

				// Get the completion channel for this batch
				completionChan := pe.processBatch(executionContext, pageBatch, errors)

				// Track max completion for this batch
				maxIndex := -1
				for index := range completionChan {
					if index > maxIndex {
						maxIndex = index
					}
				}

				// Send update if we processed anything
				if maxIndex >= 0 && maxIndex < len(pageBatch) {
					select {
					case updates <- models.JobUpdate{
						Time:      time.Now(),
						Completed: int64(maxIndex + 1),
						Failed:    int64(len(pageBatch) - (maxIndex + 1)),
						NewCheckpoint: models.JobCheckpoint{
							BatchID: &pageBatch[maxIndex],
						},
					}:
					case <-executionContext.Done():
						return
					}
				}

				// Handle rate limiting
				elapsedTime := time.Since(batchStartTime)
				if remainingTime := pe.runtimeConfig.LowerBound - elapsedTime; remainingTime > 0 {
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
				case errors <- models.JobError{Error: err, Time: time.Now()}:
				case <-executionContext.Done():
					return
				}

			case <-executionContext.Done():
				return
			}
		}
	}()

	return updates, errors
}

func (pe *pageExecutor) processBatch(ctx context.Context, pageBatch []uuid.UUID, errors chan models.JobError) <-chan int {
	pe.logger.Debug("processing page batch",
		zap.Any("batch", pageBatch),
		zap.Duration("upperBound", pe.runtimeConfig.UpperBound))

	completions := make(chan int, len(pageBatch))

	// Validate timeout
	if pe.runtimeConfig.UpperBound <= 0 {
		pe.logger.Error("invalid upper bound timeout",
			zap.Duration("upperBound", pe.runtimeConfig.UpperBound))
		close(completions)
		return completions
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, pe.runtimeConfig.UpperBound)

	var wg sync.WaitGroup

	for index, pageID := range pageBatch {
		wg.Add(1)
		go func(pageIndex int, pageID uuid.UUID) {
			defer wg.Done()

			start := time.Now()
			err := pe.processPage(timeoutCtx, pageID)
			duration := time.Since(start)

			if err != nil {
				pe.logger.Debug("page processing failed",
					zap.Any("pageID", pageID),
					zap.Duration("duration", duration),
					zap.Error(err))

				select {
				case errors <- models.JobError{Error: err, Time: time.Now()}:
				case <-timeoutCtx.Done():
				}
				return
			}

			pe.logger.Debug("page processing succeeded",
				zap.Any("pageID", pageID),
				zap.Duration("duration", duration))

			select {
			case completions <- pageIndex:
			case <-timeoutCtx.Done():
				pe.logger.Debug("completion send timed out",
					zap.Any("pageID", pageID))
			}
		}(index, pageID)
	}

	// Close completion channel when all work is done
	go func() {
		wg.Wait()
		pe.logger.Debug("all workers completed")
		close(completions)
		cancel() // Clean up timeout context
	}()

	return completions
}

func (pe *pageExecutor) processPage(ctx context.Context, pageID uuid.UUID) error {
	pe.logger.Debug("processing page", zap.Any("pageID", pageID))
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return pe.pageService.RefreshPage(ctx, pageID)
	}
}

func (pe *pageExecutor) Terminate(ctx context.Context) error {
	// No cleanup required for the page executor
	// as it does not hold any shared resources
	return nil
}
