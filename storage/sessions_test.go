package storage

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/mailbadger/app/entities"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestSessions(t *testing.T) {
	db := openTestDb()
	defer func() {
		err := db.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()
	store := From(db)

	sess, err := store.GetSession("foobar")
	assert.NotNil(t, err)

	sess = &entities.Session{
		UserID:    1,
		SessionID: "foobar",
	}

	err = store.CreateSession(sess)
	assert.Nil(t, err)

	sess, err = store.GetSession("foobar")
	assert.Nil(t, err)

	assert.Equal(t, sess.SessionID, "foobar")
	assert.Equal(t, sess.User.Username, "admin")
	assert.NotNil(t, sess.User.Boundaries)
	assert.Equal(t, sess.User.Boundaries.Type, "nolimit")

	err = store.DeleteSession("foobar")
	assert.Nil(t, err)

	_, err = store.GetSession("foobar")
	assert.NotNil(t, err)
	assert.True(t, gorm.IsRecordNotFoundError(err))
}
