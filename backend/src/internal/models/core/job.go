// ./src/internal/models/core/job.go
package models

import (
	"context"
	"database/sql"
	"encoding/json"
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

func NewJobState() *JobState {
	return &JobState{
		Status: JobStatusRunning,
		Checkpoint: JobCheckpoint{
			BatchID: nil,
		},
	}
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

// WorkflowRecord represents a historical record of a workflow
type JobRecord struct {
	// ID is the unique identifier for the job record
	ID uuid.UUID `json:"id"`
	// WorkflowType is the type of the workflow
	WorkflowType WorkflowType `json:"workflow_type"`
	// JobID is the unique identifier of the job
	JobID uuid.UUID `json:"job_id"`
	// StartTime is the time when the job started
	StartTime sql.NullTime `json:"start_time,omitempty"`
	// EndTime is the time when the job ended
	EndTime sql.NullTime `json:"end_time,omitempty"`
	// CancelTime is the time when the job was cancelled
	CancelTime sql.NullTime `json:"cancel_time,omitempty"`
	// Preemptions is the number of times the job was pre-empted
	Preemptions int `json:"preemptions"`
}

// jobRecordJSON is an internal type for JSON marshaling/unmarshaling
type jobRecordJSON struct {
	ID           string       `json:"id"`
	WorkflowType WorkflowType `json:"workflow_type"`
	JobID        string       `json:"job_id"`
	StartTime    *time.Time   `json:"start_time,omitempty"`
	EndTime      *time.Time   `json:"end_time,omitempty"`
	CancelTime   *time.Time   `json:"cancel_time,omitempty"`
	Preemptions  int          `json:"preemptions"`
}

// MarshalJSON implements custom JSON marshaling for JobRecord
func (j JobRecord) MarshalJSON() ([]byte, error) {
	record := jobRecordJSON{
		ID:           j.ID.String(),
		WorkflowType: j.WorkflowType,
		JobID:        j.JobID.String(),
		Preemptions:  j.Preemptions,
	}

	if j.StartTime.Valid {
		record.StartTime = &j.StartTime.Time
	}

	if j.EndTime.Valid {
		record.EndTime = &j.EndTime.Time
	}

	if j.CancelTime.Valid {
		record.CancelTime = &j.CancelTime.Time
	}

	return json.Marshal(record)
}

// UnmarshalJSON implements custom JSON unmarshaling for JobRecord
func (j *JobRecord) UnmarshalJSON(data []byte) error {
	var record jobRecordJSON
	if err := json.Unmarshal(data, &record); err != nil {
		return fmt.Errorf("failed to unmarshal job record: %w", err)
	}

	// Parse ID UUID
	id, err := uuid.Parse(record.ID)
	if err != nil {
		return fmt.Errorf("invalid record ID: %w", err)
	}
	j.ID = id

	// Parse JobID UUID
	jobID, err := uuid.Parse(record.JobID)
	if err != nil {
		return fmt.Errorf("invalid job ID: %w", err)
	}
	j.JobID = jobID

	j.WorkflowType = record.WorkflowType
	j.Preemptions = record.Preemptions

	// Handle nullable times
	if record.StartTime != nil {
		j.StartTime = sql.NullTime{
			Time:  *record.StartTime,
			Valid: true,
		}
	} else {
		j.StartTime = sql.NullTime{Valid: false}
	}

	if record.EndTime != nil {
		j.EndTime = sql.NullTime{
			Time:  *record.EndTime,
			Valid: true,
		}
	} else {
		j.EndTime = sql.NullTime{Valid: false}
	}

	if record.CancelTime != nil {
		j.CancelTime = sql.NullTime{
			Time:  *record.CancelTime,
			Valid: true,
		}
	} else {
		j.CancelTime = sql.NullTime{Valid: false}
	}

	return nil
}
