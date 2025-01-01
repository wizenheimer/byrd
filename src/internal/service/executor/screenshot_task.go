package executor

import (
	"context"
	"sync"
	"time"

	exc "github.com/wizenheimer/iris/src/internal/interfaces/executor"
	svc "github.com/wizenheimer/iris/src/internal/interfaces/service"
	core_models "github.com/wizenheimer/iris/src/internal/models/core"

	"github.com/wizenheimer/iris/src/pkg/logger"
	"go.uber.org/zap"
)

type screenshotTaskExecutor struct {
	config            core_models.ExecutorConfig
	urlService        svc.URLService
	screenshotService svc.ScreenshotService
	diffService       svc.DiffService
	logger            *logger.Logger

	activeTasks sync.Map // map[string]context.CancelFunc
}

type batchResults struct {
	successful int
	failed     int
}

func NewScreenshotTaskExecutor(
	urlService svc.URLService,
	screenshotService svc.ScreenshotService,
	diffService svc.DiffService,
	logger *logger.Logger,
) (exc.TaskExecutor, error) {
	config, err := core_models.GetExecutorConfig(core_models.ScreenshotWorkflowType)
	if err != nil {
		return nil, err
	}

	taskExecutor := &screenshotTaskExecutor{
		config:            config,
		urlService:        urlService,
		screenshotService: screenshotService,
		diffService:       diffService,
		logger:            logger,
	}

	return taskExecutor, nil
}

func (e *screenshotTaskExecutor) Execute(ctx context.Context, task core_models.Task) (<-chan core_models.TaskUpdate, <-chan core_models.TaskError) {
	updates := make(chan core_models.TaskUpdate, 1)
	errors := make(chan core_models.TaskError, 1)

	taskCtx, cancel := context.WithCancel(ctx)
	e.activeTasks.Store(task.TaskID, cancel)

	go func() {
		defer e.cleanup(task.TaskID, updates, errors)

		var completed, failed int
		checkpoint := task.Checkpoint

		// Get URL batches stream
		urlBatchChan, errBatchChan := e.urlService.ListURLs(taskCtx, e.config.Parallelism, checkpoint.BatchID)

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
				errors <- core_models.TaskError{
					TaskID: task.TaskID,
					Error:  err,
					Time:   time.Now(),
				}
				// Back off on error
				time.Sleep(e.config.LowerBound)

			case batch, ok := <-urlBatchChan:
				if !ok {
					e.logger.Debug("url batch channel closed, task complete",
						zap.Any("task_id", task.TaskID),
						zap.Any("completed", completed),
						zap.Any("failed", failed))
					// Send final update
					updates <- core_models.TaskUpdate{
						TaskID:        task.TaskID,
						Status:        core_models.TaskStatusComplete,
						Completed:     completed,
						Failed:        failed,
						NewCheckpoint: checkpoint,
					}
					return
				}

				if len(batch.URLs) == 0 {
					e.logger.Debug("empty batch received",
						zap.Any("task_id", task.TaskID))
					continue
				}

				// Process batch
				results := e.processBatch(taskCtx, batch.URLs)
				completed += results.successful
				failed += results.failed

				// Update checkpoint with last URL's ID
				if lastURL := batch.URLs[len(batch.URLs)-1]; lastURL.ID != nil {
					checkpoint = core_models.WorkflowCheckpoint{
						BatchID: lastURL.ID,
					}
				}

				// Send progress update
				updates <- core_models.TaskUpdate{
					TaskID:        task.TaskID,
					Status:        core_models.TaskStatusRunning,
					Completed:     completed,
					Failed:        failed,
					NewCheckpoint: checkpoint,
				}

			}
		}
	}()

	return updates, errors
}

func (e *screenshotTaskExecutor) processBatch(ctx context.Context, urls []core_models.URL) batchResults {
	var results batchResults
	var wg sync.WaitGroup
	resultChan := make(chan bool, len(urls)) // true = success, false = failure

	for _, url := range urls {
		wg.Add(1)
		go func(url core_models.URL) {
			defer wg.Done()

			err := e.processURL(ctx, url)
			resultChan <- err == nil
		}(url)
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

func (e *screenshotTaskExecutor) processURL(ctx context.Context, url core_models.URL) error {
	// Create context with timeout for this URL
	urlContext, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Capture current screenshot
	_, currentHtmlContentResp, err := e.screenshotService.Refresh(urlContext, url.URL, core_models.ScreenshotRequestOptions{
		URL: url.URL,
	})
	if err != nil {
		return err
	}

	// Get previous screenshot for comparison
	_, previousHtmlContentResp, err := e.screenshotService.Retrieve(urlContext, url.URL)
	if err != nil {
		return err
	}

	// Compare screenshots
	diffHTMLResp, err := e.diffService.Compare(urlContext, currentHtmlContentResp, previousHtmlContentResp, "competitor_analysis", true)
	if err != nil {
		return err
	}

	e.logger.Debug("diff analysis completed", zap.Any("url", url.URL), zap.Any("diff", diffHTMLResp))
	return nil
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

func (e *screenshotTaskExecutor) cleanup(taskID string, updates chan<- core_models.TaskUpdate, errors chan<- core_models.TaskError) {
	e.activeTasks.Delete(taskID)
	close(updates)
	close(errors)
}
