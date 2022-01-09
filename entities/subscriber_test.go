package entities

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSubscriberEntity(t *testing.T) {
	var (
		subID int64 = 123
		now         = time.Now()
	)

	sub := &Subscriber{
		Model: Model{
			ID:        subID,
			CreatedAt: now,
			UpdatedAt: now,
		},
		MetaJSON: []byte(`{"foo": "bar"}`),
		Email:    "john.doe@example.com",
	}

	url, err := sub.GetUnsubscribeURL("foobar", "secret", "http://example.com")
	assert.Nil(t, err)

	m, err := sub.GetMetadata()
	assert.Nil(t, err)

	assert.Equal(t, m["foo"], "bar")
	assert.Equal(t, url, "http://example.com/unsubscribe.html?email=john.doe%40example.com&t=77de38e4b50e618a0ebb95db61e2f42697391659d82c064a5f81b9f48d85ccd5&uuid=foobar")

	tt, err := sub.GenerateUnsubscribeToken("secret")
	assert.Nil(t, err)
	assert.Equal(t, tt, "77de38e4b50e618a0ebb95db61e2f42697391659d82c064a5f81b9f48d85ccd5")

	id := sub.GetID()
	assert.Equal(t, subID, id)

	createdAt := sub.GetCreatedAt()
	assert.Equal(t, now, createdAt)

	updatedAt := sub.GetUpdatedAt()
	assert.Equal(t, now, updatedAt)
}
