package storage

import (
	"errors"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/mailbadger/app/entities"
)

func TestTemplates(t *testing.T) {
	db := openTestDb()
	defer func() {
		err := db.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()
	store := From(db)

	template, err :=store.GetTemplate(0, 0)
	assert.Equal(t, &entities.Template{}, template)
	assert.Equal(t, errors.New("record not found"), err)
}
