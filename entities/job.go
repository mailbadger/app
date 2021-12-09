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
	ID                int64     `json:"-" gorm:"column:id; primary_key:yes"`
	Name              string    `json:"-"`
	LastProcessedDate time.Time `json:"-"`
	Status            string    `json:"-"`
	CreatedAt         time.Time `json:"-"`
	UpdatedAt         time.Time `json:"-"`
}
