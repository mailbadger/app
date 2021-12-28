package metric

import (
	"context"
	"time"

	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/logger"
	"github.com/mailbadger/app/storage"
)

type Cron struct {
	store storage.Storage
}

func NewCron(store storage.Storage) *Cron {
	return &Cron{store}
}

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
	logger.From(ctx).Debug("cron: starting execute")
	//todo do the metrics
	c.store.GetJobByName(entities.JobSubscriberMetrics)
	logger.From(ctx).Debug("cron: finished executing")
	return nil
}
