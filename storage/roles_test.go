package storage

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/mailbadger/app/entities"
)

func TestRoles(t *testing.T) {
	db := openTestDb()

	store := From(db)

	_, err := store.GetRole("foobar")
	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))

	r, err := store.GetRole(entities.AdminRole)

	assert.Nil(t, err)
	assert.Equal(t, r.Name, entities.AdminRole)
}
