package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sync"
	"time"
	
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	
	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/storage"
)

func main() {
	driver := os.Getenv("DATABASE_DRIVER")
	config := storage.MakeConfigFromEnv(driver)
	
	s := storage.New(driver, config)
	
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	
	var (
		dateStr      = flag.String("date", yesterday, "Sync payments for current date")
		startDateStr = flag.String("start_date", yesterday, "Start date for syncing payments summaries")
		endDateStr   = flag.String("end_date", yesterday, "End date for syncing payments summaries")
	)
	
	flag.Parse()
	
	startDate, endDate, err := parseDates(yesterday, *dateStr, *startDateStr, *endDateStr)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"date_flag":       dateStr,
			"start_date_flag": startDateStr,
			"end_date_flag":   endDateStr,
		}).WithError(err).Fatal("failed to parse dates")
	}
	
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
		lf.WithFields(logrus.Fields{
			"job_status": job.Status,
		}).WithError(err).
			Fatal("jobs' status is not idle")
	}
	
	job.Status = entities.JobStatusInProgress
	err = s.UpdateJob(job)
	if err != nil {
		lf.WithError(err).
			Fatal("failed to update job to in-progress")
	}
	
	lf.WithField("job_status", job.Status)
	
	defer func() {
		err = s.UpdateJob(job)
		if err != nil {
			lf.WithError(err).
				Fatalf("failed to update job to %s", job.Status)
		}
	}()
	
	job.Status = processEvents(s, startDate, endDate, lf)
}

func parseDates(yesterday, dateStr, startDateStr, endDateStr string) (sd time.Time, ed time.Time, err error) {
	if startDateStr != yesterday && endDateStr != yesterday {
		sd, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("parsing current time: %w", err)
		}
		
		ed, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("parsing current time: %w", err)
		}
		
		return sd, ed.AddDate(0, 0, 1), nil
	}
	
	cd, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("parsing date: %w", err)
	}
	
	return cd, cd.AddDate(0, 0, 1), nil
}

func processEvents(s storage.Storage, startDate, endDate time.Time, lf *logrus.Entry) string {
	events, err := s.GetGroupedSubscriberEvents(startDate, endDate)
	if err != nil {
		lf.WithError(err).
			Error("failed to group events")
		return entities.JobStatusDirty
	}
	
	logrus.Infof("There are %d grouped events from %s, to %s to be proccessed", len(events), startDate, endDate)
	
	if len(events) == 0 {
		lf.Info("there are no events to process")
		return entities.JobStatusIdle
	}
	
	var (
		i                   int
		wg                  sync.WaitGroup
		workers             = 4
		chReducers          = make(chan map[string]*entities.SubscriberMetrics, workers)
		chGroupedEvents     = make(chan *entities.GroupedSubscriberEvents)
		chSubscriberMetrics = make(chan *entities.SubscriberMetrics)
	)
	
	for i = 0; i < workers; i++ {
		wg.Add(1)
		go worker(&wg, chGroupedEvents, chReducers)
	}
	
	for _, event := range events {
		chGroupedEvents <- event
	}
	
	close(chGroupedEvents)
	wg.Wait()
	close(chReducers)
	
	m := make(map[string]*entities.SubscriberMetrics)
	for reducer := range chReducers {
		for k, v := range reducer {
			logrus.Infof("record %+v\n", v)
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
		lf.WithError(err).
			Error("failed to upsert subscriber metrics")
		return  entities.JobStatusIdle
	}
	
	return entities.JobStatusIdle
}

func worker(wg *sync.WaitGroup, events chan *entities.GroupedSubscriberEvents, reducers chan map[string]*entities.SubscriberMetrics) {
	defer wg.Done()
	
	m := make(map[string]*entities.SubscriberMetrics)
	for event := range events {
		k := fmt.Sprintf("%d-%s", event.UserID, event.Date.Format("2006-01-02"))
		if _, ok := m[k]; !ok {
			m[k] = &entities.SubscriberMetrics{
				UserID:       event.UserID,
				Created:      0,
				Deleted:      0,
				Unsubscribed: 0,
				Date:         event.Date, // TODO need to be just date without time
			}
		}
		
		switch entities.EventType(event.EventType) {
		case entities.SubscriberEventTypeCreated:
			m[k].Created += event.Total
		case entities.SubscriberEventTypeDeleted:
			m[k].Deleted += event.Total
		case entities.SubscriberEventTypeUnsubscribed:
			m[k].Unsubscribed += event.Total
		default:
			logrus.WithFields(logrus.Fields{
				"user_id": event.UserID,
				"date":    event.Date,
			}).Error("event type is unsupported")
		}
	}
	
	reducers <- m
}

func upsertSubscriberMetrics(s storage.Storage, metrics chan *entities.SubscriberMetrics) error {
	for metric := range metrics {
		err := s.UpdateSubscriberMetrics(metric)
		if err != nil {
			return err // TODO make our own error to track the date
		}
	}
	
	return nil
}
