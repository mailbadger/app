package storage

import (
	"errors"
	"testing"

	"github.com/mailbadger/app/entities"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestSessions(t *testing.T) {
	db := openTestDb()
	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			logrus.Error(err)
		}
		sqlDB.Close()
	}()
	store := From(db)

	sess, err := store.GetSession("foobar")
	assert.NotNil(t, err)
	assert.Nil(t, sess)

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
	assert.Equal(t, sess.User.Boundaries.Type, entities.BoundaryTypeNoLimit)

	err = store.DeleteSession("foobar")
	assert.Nil(t, err)

	_, err = store.GetSession("foobar")
	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))

	// test delete all sessions for user
	sess = &entities.Session{
		UserID:    1,
		SessionID: "delete-session",
	}

	err = store.CreateSession(sess)
	assert.Nil(t, err)

	err = store.DeleteAllSessionsForUser(1)
	assert.Nil(t, err)

	_, err = store.GetSession("delete-session")
	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}
