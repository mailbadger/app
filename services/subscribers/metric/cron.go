package metric

import (
	"context"
	"fmt"
	"time"

	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/logger"
	"github.com/mailbadger/app/storage"
)

// Cron represents the subscribers metrics cronjob.
type Cron struct {
	store storage.Storage
}

// NewCron instantiates a new Cron object.
func NewCron(store storage.Storage) *Cron {
	return &Cron{store}
}

// Start starts executing the job on the given interval.
func (c *Cron) Start(ctx context.Context, d time.Duration) error {
	logger.From(ctx).Debug("cron: starting subscriber metrics cron")
	ticker := time.NewTicker(d)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			err := c.execute(ctx)
			if err != nil {
				logger.From(ctx).WithError(err).Error("cron: execute returned error")
			}
		}
	}
}

func (c *Cron) execute(ctx context.Context) error {
	l := logger.From(ctx)
	l.Debug("cron: starting execute")
	job, err := c.store.GetJobByName(entities.JobSubscriberMetrics)
	if err != nil {
		return fmt.Errorf("cron: failed to get job: %w", err)
	}

	if job.Status != entities.JobStatusIdle {
		return fmt.Errorf("cron: job status is %s when it should be %s", job.Status, entities.JobStatusIdle)
	}

	job.Status = entities.JobStatusInProgress
	err = c.store.UpdateJob(job)
	if err != nil {
		return fmt.Errorf("cron: unable to update job's status: %w", err)
	}

	l.Debug("cron: finished executing")
	return nil
}
