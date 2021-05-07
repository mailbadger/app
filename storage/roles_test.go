package storage

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/mailbadger/app/entities"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestRoles(t *testing.T) {
	db := openTestDb()
	defer func() {
		err := db.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	store := From(db)

	_, err := store.GetRole("foobar")
	assert.NotNil(t, err)
	assert.EqualError(t, err, gorm.ErrRecordNotFound.Error())

	r, err := store.GetRole(entities.AdminRole)

	assert.Nil(t, err)
	assert.Equal(t, r.Name, entities.AdminRole)
}
