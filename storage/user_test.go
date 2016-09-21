package storage

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	db := openTestDb()
	defer db.Close()

	store := From(db)

	//Test get admin user
	user, err := store.GetUser(1)
	assert.Nil(t, err)
	assert.Equal(t, user.Username, "admin")

	//Update user test
	user.Username = "foo"
	err = store.UpdateUser(&user)
	assert.Nil(t, err)

	user, err = store.GetUser(1)
	assert.Equal(t, user.Username, "foo")

	//ApiKey length test
	key, _ := base64.URLEncoding.DecodeString(user.ApiKey)
	assert.Nil(t, err)
	assert.Len(t, key, 32)
}
