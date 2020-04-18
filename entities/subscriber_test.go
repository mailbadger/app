package entities

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubscriberEntity(t *testing.T) {
	os.Setenv("APP_URL", "https://mailbadger.io")
	os.Setenv("UNSUBSCRIBE_SECRET", "secret")

	sub := &Subscriber{
		Model:    Model{ID: 123},
		MetaJSON: []byte(`{"foo": "bar"}`),
		Email:    "john.doe@example.com",
	}

	err := sub.AppendUnsubscribeURLToMeta()
	assert.Nil(t, err)

	var m map[string]string

	err = json.Unmarshal(sub.MetaJSON, &m)
	assert.Nil(t, err)
	assert.Equal(t, m["foo"], "bar")
	assert.Equal(t, m["unsubscribe_url"], "https://mailbadger.io/unsubscribe?email=john.doe%40example.com&token=77de38e4b50e618a0ebb95db61e2f42697391659d82c064a5f81b9f48d85ccd5")

	tt, err := sub.GenerateUnsubscribeToken(os.Getenv("UNSUBSCRIBE_SECRET"))
	assert.Nil(t, err)
	assert.Equal(t, tt, "77de38e4b50e618a0ebb95db61e2f42697391659d82c064a5f81b9f48d85ccd5")

}
