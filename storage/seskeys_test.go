package storage

import (
	"testing"

	"github.com/news-maily/app/entities"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestSesKeys(t *testing.T) {
	db := openTestDb()
	defer func() {
		err := db.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()
	store := From(db)

	_, err := store.GetSesKeys(1)
	assert.NotNil(t, err)

	keys := &entities.SesKeys{
		UserID:    1,
		AccessKey: "abcd",
		SecretKey: "efgh",
		Region:    "eu-west-1",
	}

	err = store.CreateSesKeys(keys)
	assert.Nil(t, err)

	keys, err = store.GetSesKeys(1)
	assert.Nil(t, err)
	assert.Equal(t, "abcd", keys.AccessKey)
	assert.Equal(t, "efgh", keys.SecretKey)
	assert.Equal(t, "eu-west-1", keys.Region)

	err = store.DeleteSesKeys(1)
	assert.Nil(t, err)

	keys, err = store.GetSesKeys(1)
	assert.NotNil(t, err)
	assert.Nil(t, keys)
}
