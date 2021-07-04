package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"
	
	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/storage"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func main() {
	driver := os.Getenv("DATABASE_DRIVER")
	config := storage.MakeConfigFromEnv(driver)
	
	s := storage.New(driver, config)
	
	startDate := time.Now().AddDate(0, 0, -1)
	endDate := time.Now()
	
	lf := logrus.WithFields(logrus.Fields{
		"job_name":   entities.Job_SubscriberMetrics,
		"start_date": startDate,
		"end_date":   endDate,
	})
	
	job, err := s.GetJobByName(entities.Job_SubscriberMetrics)
	if err != nil {
		lf.WithError(err).
			Fatal("failed to fetch job")
	}
	
	if job.Status != entities.JobStatusIdle {
		lf.WithError(err).
			Fatal("jobs' status is not idle")
	}
	
	job.Status = entities.JobStatusInProgress
	err = s.UpdateJob(job)
	if err != nil {
		lf.WithError(err).
			Fatal("failed to update job to in-progress")
	}
	
	defer func() {
		err = s.UpdateJob(job)
		if err != nil {
			lf.WithError(err).
				Fatalf("failed to update job to %s", job.Status)
		}
	}()
	
	events, err := s.GetGroupedSubscriberEvents(startDate, endDate)
	if err != nil {
		job.Status = entities.JobStatusDirty
		lf.WithError(err).
			Error("failed to group events")
		return
	}
	
	if len(events) == 0 {
		job.Status = entities.JobStatusIdle
		lf.Info("there are no events to process")
		return
	}
	
	var (
		i                   int
		wg                  sync.WaitGroup
		workers             = 4
		chReducers          = make(chan map[string]*entities.SubscribersMetrics, workers)
		chGroupedEvents     = make(chan *entities.GroupedSubscriberEvents)
		chSubscriberMetrics = make(chan *entities.SubscribersMetrics)
	)
	
	for i = 0; i < workers; i++ {
		wg.Add(1)
		go processEvent(&wg, chGroupedEvents, chReducers)
	}
	
	for _, event := range events {
		chGroupedEvents <- event
	}
	
	close(chGroupedEvents)
	wg.Wait()
	close(chGroupedEvents)
	
	m := <-chReducers
	for reducer := range chReducers {
		for k, v := range reducer {
			if _, ok := m[k]; !ok {
				m[k] = v
				continue
			}
			
			// it is safe to add the values since all types are counted daily
			m[k].Created += v.Created
			m[k].Deleted += v.Deleted
			m[k].Unsubscribed += v.Unsubscribed
		}
	}
	
	errGroup, _ := errgroup.WithContext(context.Background())
	
	for i = 0; i < workers; i++ {
		errGroup.Go(func() error {
			return upsertSubscriberMetrics(s, chSubscriberMetrics)
		})
	}
	
	for _, v := range m {
		chSubscriberMetrics <- v
	}
	
	close(chSubscriberMetrics)
	err = errGroup.Wait()
	if err != nil {
		job.Status = entities.JobStatusIdle
		lf.WithError(err).
			Error("failed to upsert subscriber metrics")
		return
	}
	
	job.Status = entities.JobStatusIdle
}

func processEvent(wg *sync.WaitGroup, events chan *entities.GroupedSubscriberEvents, reducers chan map[string]*entities.SubscribersMetrics) {
	defer wg.Done()
	
	m := make(map[string]*entities.SubscribersMetrics)
	for event := range events {
		k := fmt.Sprintf("%d-%s", event.UserID, event.Date.Format("2006-01-02"))
		if _, ok := m[k]; !ok {
			m[k] = &entities.SubscribersMetrics{
				UserID:       event.UserID,
				Created:      0,
				Deleted:      0,
				Unsubscribed: 0,
				Date:         event.Date, // TODO need to be just date without time
			}
			
			switch event.EventType {
			case entities.SubscriberEventTypeCreated:
				m[k].Created = event.Total
			case entities.SubscriberEventTypeDeleted:
				m[k].Deleted = event.Total
			case entities.SubscriberEventTypeUnsubscribed:
				m[k].Unsubscribed = event.Total
			default:
				logrus.WithFields(logrus.Fields{
					"user_id": event.UserID,
					"date":    event.Date,
				}).Error("event type is unsupported")
			}
		}
	}
	
	reducers <- m
}

func upsertSubscriberMetrics(s storage.Storage, metrics chan *entities.SubscribersMetrics) error {
	for metric := range metrics {
		err := s.UpdateSubscriberMetrics(metric)
		if err != nil {
			return err // TODO make our own error to track the date
		}
	}
	
	return nil
}
