package main

import (
	"context"
	"errors"
	"os"
	"sync"
	
	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/storage"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

var (
	limit int64 = 1000
)

func main() {
	driver := os.Getenv("DATABASE_DRIVER")
	config := storage.MakeConfigFromEnv(driver)
	
	s := storage.New(driver, config)
	
	job, err := s.GetJobByName(entities.Job_SubscriberMetrics)
	if err != nil {
		logrus.WithField("job_name", entities.Job_SubscriberMetrics).
			WithError(err).
			Fatal("failed to fetch job")
	}
	
	if job.Status != entities.JobStatusIdle {
		logrus.WithField("job_name", entities.Job_SubscriberMetrics).
			WithError(err).
			Fatal("jobs' status is not idle")
	}
	
	job.Status = entities.JobStatusInProgress
	err = s.UpdateJob(job)
	if err != nil {
		logrus.WithField("job_name", entities.Job_SubscriberMetrics).
			WithError(err).
			Fatal("failed to update job to in-progress")
	}
	
	defer func() {
		err = s.UpdateJob(job)
		if err != nil {
			logrus.WithField("job_name" ,entities.Job_SubscriberMetrics).
				WithError(err).
				Fatalf("failed to update job to %s", job.Status)
		}
	}()
	
	events, err := s.GetEventsAfterID(job.LastProcessedID, limit)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"job_name":          entities.Job_SubscriberMetrics,
			"last_processed_id": job.LastProcessedID,
		}).WithError(err).
			Fatal("failed to fetch events")
	}
	
	if len(events) == 0 {
		logrus.Info("there are no events to process")
		return
	}
	
	workers := 4
	var mux sync.Mutex
	chEvents := make(chan *entities.SubscriberEvent)
	
	g, _ := errgroup.WithContext(context.Background())
	
	for i := 0; i < workers; i++ {
		g.Go(func() error {
			for _, event := range events {
				sm := &entities.SubscribersMetrics{
					UserID:       event.UserID,
					Created:      0,
					Deleted:      0,
					Unsubscribed: 0,
					Date:         event.CreatedAt, // TODO need to be just date without time
				}
				
				switch event.EventType {
				case entities.SubscriberEventTypeCreated:
					sm.Created += 1
				case entities.SubscriberEventTypeDeleted:
					sm.Deleted += 1
				case entities.SubscriberEventTypeUnsubscribed:
					sm.Unsubscribed += 1
				default:
					logrus.WithField("event_id", event.ID).Error("event type is unsupported")
					return errors.New("unsupported event type")
				}
				
				mux.Lock()
				err = s.UpdateSubscriberMetrics(sm)
				if err != nil {
					return err
				}
				
				mux.Unlock()
			}
			
			return nil
		})
	}
	
	lastProcessedID := events[len(events)-1].ID
	for _, event := range events {
		chEvents <- event
	}
	
	close(chEvents)
	err = g.Wait()
	if err != nil {
		job.Status = entities.JobStatusDirty
	}
	
	job.LastProcessedID = lastProcessedID
	job.Status = entities.JobStatusIdle
}
