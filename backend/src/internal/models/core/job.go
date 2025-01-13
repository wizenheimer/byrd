package models

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Job represents a running workflow
// Each Job has a unique identifier and a state
type Job struct {
	// JobID is the unique identifier of the job
	JobID uuid.UUID `json:"job_id"`

	// JobState is the state of the job
	JobState
}

func NewJob() *Job {
	return &Job{
		JobID: uuid.New(),
		JobState: JobState{
			Status: JobStatusRunning,
			Checkpoint: JobCheckpoint{
				BatchID: nil,
			},
		},
	}
}

// JobState captures the state of the job
type JobState struct {
	// Status is the current status of the job
	Status JobStatus `json:"status" validate:"required" default:"running"`

	// Checkpoint is the current checkpoint of the job
	Checkpoint JobCheckpoint `json:"checkpoint"`

	// Completed is the number of iterations completed
	Completed int64 `json:"completed"`

	// Failed is the number of iterations failed
	Failed int64 `json:"failed"`
}

// JobCheckpoint captures the current checkpoint of the workflow
type JobCheckpoint struct {
	// BatchID is the batch ID of the current checkpoint
	BatchID *uuid.UUID `json:"batch_id"`
}

// WorkflowStatus is an enum for the status of a workflow
type JobStatus string

const (
	JobStatusRunning   JobStatus = "running"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "failed"
	JobStatusAborted   JobStatus = "aborted"
	JobStatusUnknown   JobStatus = "unknown"
)

func ParseJobStatus(s string) (JobStatus, error) {
	switch JobStatus(s) {
	case JobStatusRunning, JobStatusCompleted,
		JobStatusFailed, JobStatusAborted, JobStatusUnknown:
		return JobStatus(s), nil
	default:
		return "", fmt.Errorf("invalid workflow status: %s", s)
	}
}

// JobError represents an error that occurred while executing a task
type JobError struct {
	Error error     `json:"error"`
	Time  time.Time `json:"time"`
}

// JobUpdate represents an update to a task
type JobUpdate struct {
	Time          time.Time     `json:"time"`
	Completed     int64         `json:"completed"`
	Failed        int64         `json:"failed"`
	NewCheckpoint JobCheckpoint `json:"new_checkpoint"`
}

type JobContext struct {
	// Job embeds the Job struct
	Job

	// cancel is the cancel function for the job
	cancel context.CancelFunc

	// mutex is the mutex for the job context
	mutex sync.Mutex
}

func NewJobContextForJob(job *Job) (*JobContext, context.Context) {
	ctx, cancel := context.WithCancel(context.Background())
	jobContext := JobContext{
		Job:    *job,
		cancel: cancel,
	}
	return &jobContext, ctx
}

func (jc *JobContext) IncrementCompleted(completed int64) {
	jc.mutex.Lock()
	defer jc.mutex.Unlock()

	jc.Completed += completed
}

func (jc *JobContext) IncrementFailed(failed int64) {
	jc.mutex.Lock()
	defer jc.mutex.Unlock()

	jc.Failed += failed
}

func (jc *JobContext) UpdateStatus(status JobStatus) {
	jc.mutex.Lock()
	defer jc.mutex.Unlock()

	jc.Status = status
}

func (jc *JobContext) GetStatus() JobStatus {
	jc.mutex.Lock()
	defer jc.mutex.Unlock()

	return jc.Status
}

func (jc *JobContext) UpdateCheckpoint(checkpoint JobCheckpoint) {
	jc.mutex.Lock()
	defer jc.mutex.Unlock()

	jc.Checkpoint = checkpoint
}

func (jc *JobContext) GetCheckpoint() JobCheckpoint {
	jc.mutex.Lock()
	defer jc.mutex.Unlock()

	return jc.Checkpoint
}

func (jc *JobContext) UpdateState(status JobStatus, checkpoint JobCheckpoint) {
	jc.mutex.Lock()
	defer jc.mutex.Unlock()

	jc.Status = status
	jc.Checkpoint = checkpoint
}

func (jc *JobContext) GetState() (JobStatus, JobCheckpoint) {
	jc.mutex.Lock()
	defer jc.mutex.Unlock()

	return jc.Status, jc.Checkpoint
}

func (jc *JobContext) HandleUpdate(jobUpdate *JobUpdate) {
	jc.mutex.Lock()
	defer jc.mutex.Unlock()

	jc.Status = JobStatusRunning
	jc.Completed += jobUpdate.Completed
	jc.Failed += jobUpdate.Failed
	jc.Checkpoint = jobUpdate.NewCheckpoint
}

func (jc *JobContext) HandleError(jobError *JobError) {
	jc.mutex.Lock()
	defer jc.mutex.Unlock()

	jc.Failed += 1
}

func (jc *JobContext) HandleCompletion() {
	jc.mutex.Lock()
	defer jc.mutex.Unlock()

	jc.Status = JobStatusCompleted
}

func (jc *JobContext) HandleCancellation() {
	// TODO: add pre and post hooks
	jc.Status = JobStatusAborted
	jc.cancel()
}

// ExecutorConfig represents the configuration for an executor
type JobExecutorConfig struct {
	// Number of tasks to execute in parallel
	// This is used to limit the number of tasks that can be executed concurrently
	Parallelism int `json:"batch_size"`
	// Lower bound for the time to wait before executing the next batch
	// This is used to prevent the executor from executing tasks too frequently
	LowerBound time.Duration `json:"lower_bound"`
	// Upper bound for the time to wait before executing the next batch
	// This is used to prevent the executor from getting stuck with the same batch
	UpperBound time.Duration `json:"upper_bound"`
}
