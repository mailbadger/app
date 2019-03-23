package storage

import (
	"testing"

	"github.com/news-maily/api/entities"
	"github.com/stretchr/testify/assert"
)

func TestSesKeys(t *testing.T) {
	db := openTestDb()
	defer db.Close()

	store := From(db)

	keys, err := store.GetSesKeys(1)
	assert.NotNil(t, err)

	keys = &entities.SesKeys{
		UserId:    1,
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

	err = store.DeleteSesKeys(keys.Id, 1)
	assert.Nil(t, err)

	keys, err = store.GetSesKeys(1)
	assert.NotNil(t, err)
	assert.Nil(t, keys)
}
