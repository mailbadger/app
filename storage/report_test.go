package storage

import (
	"errors"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/mailbadger/app/entities"
)

func TestReport(t *testing.T) {
	db := openTestDb()
	defer func() {
		err := db.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()
	now := time.Now()

	store := From(db)

	reports := []entities.Report{
		{
			UserID:   1,
			Resource: "subscriptions",
			FileName: "subv1",
			Type:     "export",
			Status:   "done",
			Note:     "",
		},
		{
			UserID:   1,
			Resource: "subscriptions",
			FileName: "subv2",
			Type:     "export",
			Status:   "failed",
			Note:     "",
		},
		{
			UserID:   2,
			Resource: "subscriptions",
			FileName: "running",
			Type:     "export",
			Status:   "in_progress",
			Note:     "",
		},
	}
	// test insert report
	for i := range reports {
		err := store.CreateReport(&reports[i])
		assert.Nil(t, err)
	}

	report, err := store.GetReportByFilename("not-found", 1)
	assert.Equal(t, errors.New("record not found"), err)
	assert.Equal(t, new(entities.Report), report)

	report, err = store.GetReportByFilename("subv1", 1)
	assert.Nil(t, err)

	assert.Equal(t, reports[0].FileName, report.FileName)
	assert.Equal(t, reports[0].Resource, report.Resource)

	// test update report
	updatedReport := entities.Report{
		Model: entities.Model{
			ID:        2,
			UpdatedAt: time.Now(),
		},
		UserID:   1,
		FileName: "subv2",
		Status:   "failed",
		Note:     "unable to unmarshal bla",
	}
	err = store.UpdateReport(&updatedReport)
	assert.Nil(t, err)

	// check updated report
	upReport, err := store.GetReportByFilename("subv2", 1)
	assert.Nil(t, err)
	assert.Equal(t, updatedReport.Status, upReport.Status)
	assert.Equal(t, updatedReport.Note, upReport.Note)

	numOfRep, err := store.GetNumberOfReportsForDate(1, now)
	assert.Nil(t, err)
	assert.Equal(t, int64(2), numOfRep)

	runningReport, err := store.GetRunningReportForUser(1)
	assert.Equal(t, errors.New("record not found"), err)
	assert.Equal(t, new(entities.Report), runningReport)

	runningReport, err = store.GetRunningReportForUser(2)
	assert.Nil(t, err)
	assert.Equal(t, reports[2].FileName, runningReport.FileName)
	assert.Equal(t, reports[2].Resource, runningReport.Resource)
	assert.Equal(t, reports[2].Type, runningReport.Type)

	err = store.DeleteAllReportsForUser(1)
	assert.Nil(t, err)

	numOfRep, err = store.GetNumberOfReportsForDate(1, now)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), numOfRep)
}
