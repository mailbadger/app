package storage

import (
	"testing"

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
			Status:   "inprogress",
			Note:     "",
		},
	}
	// test insert report
	for _, r := range reports {
		err := store.CreateReport(&r)
		assert.Nil(t, err)
	}

	report, err := store.GetReportByFilename("subv1", 1)
	assert.Nil(t, err)

	assert.Equal(t, reports[0].FileName, report.FileName)
	assert.Equal(t, reports[0].Resource, report.Resource)

	// test update report
	updatedReport := entities.Report{
		UserID:   1,
		FileName: "subv2",
		Status:   "failed",
		Note:     "unable to unmarshal bla",
	}
	err = store.UpdateReport(&updatedReport)

	// check updated report
	upReport, err := store.GetReportByFilename("subv2", 1)
	assert.Nil(t, err)

	assert.Equal(t, updatedReport.Status, upReport.Status)
	assert.Equal(t, updatedReport.Note, upReport.Note)

}
