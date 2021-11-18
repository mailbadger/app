package storage

import (
	"testing"

	"github.com/mailbadger/app/entities"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestRoles(t *testing.T) {
	db := openTestDb()
	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			logrus.Error(err)
		}
		sqlDB.Close()
	}()

	store := From(db)

	_, err := store.GetRole("foobar")
	assert.NotNil(t, err)
	assert.EqualError(t, err, gorm.ErrRecordNotFound.Error())

	r, err := store.GetRole(entities.AdminRole)

	assert.Nil(t, err)
	assert.Equal(t, r.Name, entities.AdminRole)
}
