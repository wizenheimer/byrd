package executor

import (
	"context"
	"sync"
	"time"

	"github.com/wizenheimer/iris/src/internal/domain/interfaces"
	"github.com/wizenheimer/iris/src/internal/domain/models"
	"github.com/wizenheimer/iris/src/pkg/logger"
)

type screenshotTaskExecutor struct {
	config            models.ExecutorConfig
	urlService        interfaces.URLService
	screenshotService interfaces.ScreenshotService
	diffService       interfaces.DiffService
	logger            *logger.Logger

	activeTasks sync.Map // map[string]context.CancelFunc
}

type batchResults struct {
	successful int
	failed     int
}

func NewScreenshotTaskExecutor(
	urlService interfaces.URLService,
	screenshotService interfaces.ScreenshotService,
	diffService interfaces.DiffService,
	logger *logger.Logger,
) (interfaces.TaskExecutor, error) {
	config, err := models.GetExecutorConfig(models.ScreenshotWorkflowType)
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

func (e *screenshotTaskExecutor) Execute(ctx context.Context, task models.Task) (<-chan models.TaskUpdate, <-chan models.TaskError) {
	updates := make(chan models.TaskUpdate, 1)
	errors := make(chan models.TaskError, 1)

	taskCtx, cancel := context.WithCancel(ctx)
	e.activeTasks.Store(task.TaskID, cancel)

	go func() {
		defer e.cleanup(task.TaskID, updates, errors)

		var completed, failed int
		checkpoint := task.Checkpoint

		for {
			select {
			case <-taskCtx.Done():
				return
			default:
				// Get next batch of URLs using checkpoint
				urlBatchChan, errBatchChan := e.urlService.ListURLs(taskCtx, e.config.Parallelism, checkpoint.BatchID)

				for {
					select {
					case <-taskCtx.Done():
						return
					case urlsBatch, ok := <-urlBatchChan:
						if !ok {
                            // Send final update
							updates <- models.TaskUpdate{
								TaskID: task.TaskID,
								Status: models.TaskStatusComplete,
							}
                            // Return
							return
						}
						// Process batch
						urls := urlsBatch.URLs
						batchResults := e.processBatch(taskCtx, urls)
						completed += batchResults.successful
						failed += batchResults.failed
						// Update checkpoint with last URL's ID
						batchID := urls[len(urls)-1].ID
						checkpoint = models.WorkflowCheckpoint{
							BatchID: batchID,
						}
						// Send update
						updates <- models.TaskUpdate{
							TaskID:        task.TaskID,
							Status:        models.TaskStatusRunning,
							Completed:     completed,
							Failed:        failed,
							NewCheckpoint: checkpoint,
						}
					case err, ok := <-errBatchChan:
						if !ok {
							return
						}
						if err != nil {
							errors <- models.TaskError{
								TaskID: task.TaskID,
								Error:  err,
								Time:   time.Now(),
							}
							time.Sleep(e.config.LowerBound)
							continue
						}
					}
				}
			}
		}
	}()

	return updates, errors
}

func (e *screenshotTaskExecutor) processBatch(ctx context.Context, urls []models.URL) batchResults {
	var results batchResults
	var wg sync.WaitGroup
	resultChan := make(chan bool, len(urls)) // true = success, false = failure

	for _, url := range urls {
		wg.Add(1)
		go func(url models.URL) {
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

func (e *screenshotTaskExecutor) processURL(_ context.Context, _ models.URL) error {
	// Create context with timeout for this URL
	// urlCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	// defer cancel()

	// Capture current screenshot

	// Get previous screenshot for comparison

	// Create diff

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

func (e *screenshotTaskExecutor) cleanup(taskID string, updates chan<- models.TaskUpdate, errors chan<- models.TaskError) {
	e.activeTasks.Delete(taskID)
	close(updates)
	close(errors)
}
