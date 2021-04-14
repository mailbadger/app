package storage

import (
	"testing"
	"time"

	"github.com/mailbadger/app/utils"
	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/mailbadger/app/entities"
)

func TestScheduledCampaign(t *testing.T) {
	db := openTestDb()
	defer func() {
		err := db.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	var now = time.Now()

	store := From(db)

	segmentIDS := []int64{1, 2, 3, 4, 5, 6}
	segmentIDSstring := utils.SliceIntToString(segmentIDS, ",")

	//Test create scheduled campaign
	c := &entities.CampaignSchedule{
		ID:                  ksuid.New(),
		UserID:              1,
		CampaignID:          1,
		ScheduledAt:         now,
		Source:              "bla@email.com",
		FromName:            "from name",
		SegmentIDs:          segmentIDSstring,
		DefaultTemplateData: []byte(`{"foo":"bar"}`),
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	err := store.CreateCampaignSchedule(c)
	assert.Nil(t, err)

	// Test delete scheduled campaign
	err = store.DeleteCampaignSchedule(c.CampaignID)
	assert.Nil(t, err)
}
