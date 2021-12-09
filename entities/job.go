package entities

import (
	"time"
	
	"github.com/segmentio/ksuid"
)

var (
	Job_SubscriberMetrics = "subscriber_metrics"
)

type Job struct {
	ID              int64       `json:"-" gorm:"column:id; primary_key:yes"`
	Name            string      `json:"-"`
	LastProcessedID ksuid.KSUID `json:"-"`
	CreatedAt       time.Time   `json:"-"`
	UpdatedAt       time.Time   `json:"-"`
}
