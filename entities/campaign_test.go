package entities

import (
	"testing"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
)

func TestSetEventID(t *testing.T) {
	now:=time.Now()

	// We dont need full campaign for this test only the campaign's schedule
	c := &Campaign{
		UserID:     1,
		Name:       "test campaign",
		TemplateID: 1,
		Model:Model{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	c.SetEventID()
	assert.NotNil(t, c.EventID)

	id := c.GetID()
	assert.NotEqual(t, 0, id)

	ca := c.GetCreatedAt()
	assert.Equal(t, now, ca)

	ua := c.GetUpdatedAt()
	assert.Equal(t, now, ua)

	uid := ksuid.New()
	// Again for this test we don't need full campaign schedule only the id
	c.Schedule = &CampaignSchedule{
		ID:         uid,
		CampaignID: 1,
	}

	c.SetEventID()
	assert.NotNil(t, c.EventID)
	assert.Equal(t, uid, *c.EventID)
}
