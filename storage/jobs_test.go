package storage

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/mailbadger/app/entities"
)

func TestJobs(t *testing.T) {
	db := openTestDb()

	store := From(db)

	job, err := store.GetJobByName("foo")
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
	assert.Equal(t, &entities.Job{}, job)

	job, err = store.GetJobByName(entities.JobSubscriberMetrics)
	assert.Nil(t, err)
	assert.Equal(t, entities.JobSubscriberMetrics, job.Name)

	job.Status = entities.JobStatusInProgress
	fmt.Println(job)
	err = store.UpdateJob(job)
	assert.Nil(t, err)

	job, err = store.GetJobByName(entities.JobSubscriberMetrics)
	assert.Nil(t, err)
	assert.Equal(t, entities.JobSubscriberMetrics, job.Name)
	assert.Equal(t, entities.JobStatusInProgress, job.Status)
}
