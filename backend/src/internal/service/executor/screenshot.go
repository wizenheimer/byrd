package executor

import (
	"context"
	"sync"
	"time"

	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/service/page"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"github.com/wizenheimer/byrd/src/pkg/utils"
	"go.uber.org/zap"
)

type screenshotTaskExecutor struct {
	config      models.ExecutorConfig
	pageService page.PageService
	logger      *logger.Logger

	activeTasks sync.Map // map[string]context.CancelFunc
}

type batchResults struct {
	successful int
	failed     int
}

func NewScreenshotTaskExecutor(
	pageService page.PageService,
	logger *logger.Logger,
) (TaskExecutor, error) {
	config, err := models.GetExecutorConfig(models.ScreenshotWorkflowType)
	if err != nil {
		return nil, err
	}

	taskExecutor := &screenshotTaskExecutor{
		config:      config,
		pageService: pageService,
		logger:      logger,
	}

	return taskExecutor, nil
}

func (e *screenshotTaskExecutor) Execute(ctx context.Context, task models.Task) (<-chan models.TaskUpdate, <-chan models.TaskError) {
	updates := make(chan models.TaskUpdate, 1)
	errors := make(chan models.TaskError, 1)

	taskCtx, cancel := context.WithCancel(ctx)
	e.activeTasks.Store(task.TaskID, cancel)

	go func() {
		defer e.cleanup(task.TaskID, updates, errors)

		var completed, failed int
		checkpoint := task.Checkpoint

		// Get URL batches stream
		pageBatchChan, errBatchChan := e.pageService.ListActivePages(taskCtx, e.config.Parallelism, checkpoint.BatchID)

		for {
			select {
			case <-taskCtx.Done():
				e.logger.Debug("task context cancelled",
					zap.Any("task_id", task.TaskID),
					zap.Any("completed", completed),
					zap.Any("failed", failed))
				return

			case err, ok := <-errBatchChan:
				if !ok {
					e.logger.Debug("error channel closed",
						zap.Any("task_id", task.TaskID))
					return
				}
				errors <- models.TaskError{
					TaskID: task.TaskID,
					Error:  err,
					Time:   time.Now(),
				}
				// Back off on error
				time.Sleep(e.config.LowerBound)

			case pages, ok := <-pageBatchChan:
				if !ok {
					e.logger.Debug("url batch channel closed, task complete",
						zap.Any("task_id", task.TaskID),
						zap.Any("completed", completed),
						zap.Any("failed", failed))
					// Send final update
					updates <- models.TaskUpdate{
						TaskID:        task.TaskID,
						Status:        models.TaskStatusComplete,
						Completed:     completed,
						Failed:        failed,
						NewCheckpoint: checkpoint,
					}
					return
				}

				if len(pages) == 0 {
					e.logger.Debug("empty batch received",
						zap.Any("task_id", task.TaskID))
					continue
				}

				// Process batch
				results := e.processBatch(taskCtx, pages)
				completed += results.successful
				failed += results.failed

				// Update checkpoint with last URL's ID
				lastPage := pages[max(len(pages)-1, 0)]
				checkpoint = models.WorkflowCheckpoint{
					BatchID: utils.ToPtr(lastPage.ID),
				}

				// Send progress update
				updates <- models.TaskUpdate{
					TaskID:        task.TaskID,
					Status:        models.TaskStatusRunning,
					Completed:     completed,
					Failed:        failed,
					NewCheckpoint: checkpoint,
				}

			}
		}
	}()

	return updates, errors
}

func (e *screenshotTaskExecutor) processBatch(ctx context.Context, pages []models.Page) batchResults {
	var results batchResults
	var wg sync.WaitGroup
	resultChan := make(chan bool, len(pages)) // true = success, false = failure

	for _, p := range pages {
		wg.Add(1)
		go func(p models.Page) {
			defer wg.Done()

			err := e.pageService.RefreshPage(ctx, p.ID)
			resultChan <- err == nil
		}(p)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(resultChan)

	// Count results
	for result := range resultChan {
		if result {
			results.successful++
		} else {
			results.failed++
		}
	}

	return results
}

func (e *screenshotTaskExecutor) Terminate(ctx context.Context) error {
	var wg sync.WaitGroup

	e.activeTasks.Range(func(key, value interface{}) bool {
		wg.Add(1)
		go func(cancel context.CancelFunc) {
			defer wg.Done()
			cancel()
		}(value.(context.CancelFunc))
		return true
	})

	// Wait for all tasks to clean up
	wg.Wait()
	return nil
}

func (e *screenshotTaskExecutor) cleanup(taskID string, updates chan<- models.TaskUpdate, errors chan<- models.TaskError) {
	e.activeTasks.Delete(taskID)
	close(updates)
	close(errors)
}
