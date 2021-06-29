package entities

import (
	"database/sql/driver"
	"time"
	
	"github.com/segmentio/ksuid"
)

var (
	Job_SubscriberMetrics = "subscriber_metrics"
)

type JobStatus string

func (js JobStatus) Value() (driver.Value, error) {
	return string(js), nil
}

func (js JobStatus) Scan(value interface{}) error {
	value = string(js)
	return nil
}

var (
	JobStatusIdle       JobStatus = "idle"
	JobStatusInProgress JobStatus = "in-progress"
	JobStatusDirty      JobStatus = "dirty"
)

type Job struct {
	ID              int64       `json:"-" gorm:"column:id; primary_key:yes"`
	Name            string      `json:"-"`
	LastProcessedID ksuid.KSUID `json:"-"`
	Status          JobStatus   `json:"-" gorm:"column:status"`
	CreatedAt       time.Time   `json:"-"`
	UpdatedAt       time.Time   `json:"-"`
}
