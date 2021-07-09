package storage

import (
	"testing"
	
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	
	"github.com/mailbadger/app/entities"
)

func TestJobs(t *testing.T) {
	db := openTestDb()
	defer func() {
		err := db.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()
	
	store := From(db)
	
	// This jobs are created with the migrations
	
	job, err := store.GetJobByName("test")
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.Equal(t, &entities.Job{}, job)
	
	job, err = store.GetJobByName(entities.Job_SubscriberMetrics)
	assert.Nil(t, err)
	assert.Equal(t, entities.Job_SubscriberMetrics, job.Name)
	
	job.Status = entities.StatusInProgress
	err = store.UpdateJob(job)
	assert.Nil(t, err)
	
	job, err = store.GetJobByName(entities.Job_SubscriberMetrics)
	assert.Nil(t, err)
	assert.Equal(t, entities.Job_SubscriberMetrics, job.Name)
	assert.Equal(t, entities.JobStatusInProgress, job.Status)
}
