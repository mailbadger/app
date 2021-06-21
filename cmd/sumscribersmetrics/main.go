package main

import (
	"os"
	
	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/storage"
	"github.com/sirupsen/logrus"
)

func main() {
	driver := os.Getenv("DATABASE_DRIVER")
	config := storage.MakeConfigFromEnv(driver)
	
	s := storage.New(driver, config)
	
	_, err := s.GetJobByName(entities.Job_SubscriberMetrics)
	if err != nil {
		logrus.WithField("job_name", entities.Job_SubscriberMetrics).
			WithError(err).
			Fatal("failed to fetch job")
	}
	
	// ---------------------
	/*
		Fetch latest events
	*/
	// ---------------------
	
	
	// ---------------------
	/*
		For each summary record fetch latest date record and add the sum if it is the same date update it,
		otherwise insert new one with the date from the summary record (keep in mind that this could be the first one)
	*/
	// ---------------------
	
	// ---------------------
	/*
		Add last processed id
	*/
	// ---------------------
}
