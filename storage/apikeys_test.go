package storage

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/mailbadger/app/entities"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestAPIKeys(t *testing.T) {
	db := openTestDb()
	defer func() {
		err := db.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()
	store := From(db)

	_, err := store.GetAPIKey("foobar")
	assert.NotNil(t, err)

	keys, err := store.GetAPIKeys(1)
	assert.Nil(t, err)
	assert.Empty(t, keys)

	k := &entities.APIKey{
		UserID:    1,
		Active:    true,
		SecretKey: "foobar",
	}

	err = store.CreateAPIKey(k)
	assert.Nil(t, err)

	k, err = store.GetAPIKey("foobar")
	assert.Nil(t, err)
	assert.Equal(t, k.SecretKey, "foobar")
	assert.True(t, k.Active)
	assert.Equal(t, k.User.Username, "admin")
	assert.NotNil(t, k.User.Boundaries)
	assert.Equal(t, k.User.Boundaries.Type, entities.BoundaryTypeNoLimit)

	k.Active = false
	err = store.UpdateAPIKey(k)
	assert.Nil(t, err)

	keys, err = store.GetAPIKeys(1)
	assert.Nil(t, err)
	assert.NotEmpty(t, keys)

	_, err = store.GetAPIKey("foobar")
	assert.NotNil(t, err)
	assert.True(t, gorm.IsRecordNotFoundError(err))

	err = store.DeleteAPIKey(k.ID, 1)
	assert.Nil(t, err)

	_, err = store.GetAPIKey("foobar")
	assert.NotNil(t, err)
	assert.True(t, gorm.IsRecordNotFoundError(err))

}
