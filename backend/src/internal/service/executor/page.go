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

	checkpoint := jobState.Checkpoint.BatchID
	pageBatchChan, errBatchChan := pe.pageService.ListActivePages(executionContext, pe.runtimeConfig.Parallelism, checkpoint)

	go func() {
		defer close(updates)
		defer close(errors)

		batchStartTime := time.Now()

		for {
			select {
			// Handle the page batch
			case pageBatch, ok := <-pageBatchChan:
				if !ok {
					goto COMPLETION
				}

				// Process the batch
				pe.processBatch(executionContext, pageBatch, updates, errors)

				// Calculate remaining time to wait
				elapsedTime := time.Since(batchStartTime)
				if remainingTime := pe.runtimeConfig.LowerBound - elapsedTime; remainingTime > 0 {
					select {
					case <-time.After(remainingTime):
					case <-executionContext.Done():
						return
					}
				}

				// Reset timer for next batch
				batchStartTime = time.Now()

				// Handle the errors, reserialize them and send them to the errors channel
			case err, ok := <-errBatchChan:
				if !ok {
					return
				}

				// Serialize the batch error into job error
				// and send it to the errors channel
				errors <- models.JobError{
					Error: err,
					Time:  time.Now(),
				}

				// Handle the context cancellation
			case <-executionContext.Done():
				return
			}

		}

	COMPLETION:
		updates <- models.JobUpdate{
			Time:      time.Now(),
			Completed: 0,
			Failed:    0,
			NewCheckpoint: models.JobCheckpoint{
				BatchID: nil,
			},
		}
	}()

	return updates, errors
}

func (pe *pageExecutor) processBatch(ctx context.Context, pageBatch []models.Page, updates chan models.JobUpdate, errors chan models.JobError) {
	// TODO: add lower bound and upper bound for the batch
	var wg sync.WaitGroup

	timeoutCtx, cancel := context.WithTimeout(ctx, pe.runtimeConfig.UpperBound)
	defer cancel()

	// Holds the index of the last page processed
	processedIndexChan := make(chan int, len(pageBatch))
	maxIndex := 0
	var mu sync.Mutex // mutex to protect maxIndex

	// Iterate over the page batch, and spawn a worker for each page
	for index, page := range pageBatch {
		// Track the number of workers
		wg.Add(1)

		// Spawn a worker to refresh the page
		// Respecting the context timeout
		go func(pageIndex int, pageID uuid.UUID) {
			defer wg.Done()
			err := pe.pageService.RefreshPage(timeoutCtx, page.ID)
			if err != nil {
				errors <- models.JobError{
					Error: err,
					Time:  time.Now(),
				}
			} else {
				processedIndexChan <- pageIndex
			}
		}(index, page.ID)
	}

	// Close processedIndexChan channel after all workers finish
	go func() {
		wg.Wait()
		close(processedIndexChan)
	}()

	// Process results until timeout or completion
	for {
		select {
		// Handle the timeout
		case <-timeoutCtx.Done():
			goto FINISH
		// Pick the maximum index processed thus far
		case index, ok := <-processedIndexChan:
			if !ok {
				goto FINISH
			}
			mu.Lock()
			if index > maxIndex {
				maxIndex = index
			}
			mu.Unlock()
		}

	}

	// Handle the completion of the batch
FINISH:
	updates <- models.JobUpdate{
		Time:      time.Now(),
		Completed: int64(maxIndex),
		Failed:    int64(len(pageBatch) - maxIndex),
		NewCheckpoint: models.JobCheckpoint{
			BatchID: &pageBatch[maxIndex].ID,
		},
	}
}

func (pe *pageExecutor) Terminate(ctx context.Context) error {
	// No cleanup required for the page executor
	// as it does not hold any shared resources
	return nil
}
