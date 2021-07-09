package entities

import (
	"database/sql/driver"
	"errors"
	"time"
)

var (
	Job_SubscriberMetrics = "subscriber_metrics"
)

type JobStatus string

func (j *JobStatus) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("failed to scan JobStatus")
	}
	
	*j = JobStatus(str)
	
	return nil
}

func (j *JobStatus) Value() (driver.Value, error) {
	return j, nil
}

var (
	JobStatusIdle       = "idle"
	JobStatusInProgress = "in-progress"
	JobStatusDirty      = "dirty"
)

type Job struct {
	Model
	Name              string    `json:"name"`
	LastProcessedDate time.Time `json:"last_processed_date"`
	Status            string    `json:"status"`
}
