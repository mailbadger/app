package entities

import "time"

var (
	Job_SubscriberMetrics = "subscriber_metrics"
)

type Job struct {
	ID              int64     `json:"-" gorm:"column:id; primary_key:yes"`
	Name            string    `json:"-"`
	LastProcessedID int64     `json:"-"`
	CreatedAt       time.Time `json:"-"`
	UpdatedAt       time.Time `json:"-"`
}
