package main

import (
	"fmt"
	"os"
	"sync"
	
	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/storage"
	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"
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
	
	// Fetch last N events
	var events []*entities.SubscriberEvent
	
	
	var wg sync.WaitGroup
	workers := 5
	chReducers := make(chan map[string]*entities.SubscribersMetrics, workers)
	chEvents := make(chan *entities.SubscriberEvent)
	
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go processEvent(chEvents, chReducers)
	}
	
	lastProcessedID := ksuid.KSUID{}
	for _, event := range events {
		chEvents <- event
		lastProcessedID = event.ID
	}
	
	close(chEvents)
	wg.Wait()
	
	m := <-chReducers
	for reducer := range chReducers {
		for k, v := range reducer {
			if _, ok := m[k]; !ok {
				m[k] = v
				continue
			}
			
			m[k].Created += v.Created
			m[k].Deleted += v.Deleted
			m[k].Unsubscribed += v.Unsubscribed
		}
	}
	
	var sm []*entities.SubscribersMetrics
	for _, v := range m {
		sm = append(sm, v)
	}
	
	job.LastProcessedID = lastProcessedID
	// Make that big transaction
	s.CreateSubscriberMetrics(sm, job)
}

func processEvent(events chan *entities.SubscriberEvent, reducers chan map[string]*entities.SubscribersMetrics) {
	m := make(map[string]*entities.SubscribersMetrics)
	for event := range events {
		k := fmt.Sprintf("%d-%s", event.UserID, event.CreatedAt.Format("2006-01-02"))
		if _, ok := m[k]; !ok {
			m[k] = &entities.SubscribersMetrics{
				UserID:       event.UserID,
				Created:      0,
				Deleted:      0,
				Unsubscribed: 0,
				Date:         event.CreatedAt, // TODO need to be just date without time
			}
			
			switch event.EventType {
			case entities.SubscriberEventTypeCreated:
				m[k].Created += 1
			case entities.SubscriberEventTypeDeleted:
				m[k].Deleted += 1
			case entities.SubscriberEventTypeUnsubscribed:
				m[k].Unsubscribed += 1
			default:
				logrus.WithField("event_id", event.ID).Error("event type is unsupported")
			}
		}
	}
	
	reducers <- m
}
