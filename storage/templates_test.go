package storage

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/mailbadger/app/entities"
)

func TestTemplate(t *testing.T) {
	db := openTestDb()
	defer func() {
		err := db.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()
	store := From(db)

	//templates for insert
	templates := []entities.Template{
		{
			UserID:   1,
			Name:     "template1",
			TextPart: "asd {{.name}}",
			Subject:  "subject",
		},
		{
			UserID:   1,
			Name:     "template2",
			TextPart: "asd {{.name}}",
			Subject:  "subject2",
		},
	}

	// test insert templates
	for _, te := range templates {
		err := store.CreateTemplate(&te)
		assert.Nil(t, err)
	}

	templates[1] = entities.Template{
		UserID:   1,
		Name:     "template2",
		TextPart: "asd {{.name}} and {{.surname}}",
		Subject:  "subject2",
	}

	err := store.UpdateTemplate(&templates[1])
	assert.Nil(t, err)

}
