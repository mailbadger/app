package entities

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValidateData(t *testing.T) {
	template := Template{
		BaseTemplate: BaseTemplate{
			Name:        "test-template",
			SubjectPart: "Hello {{name}}",
		},
		HTMLPart: "<h1>My favourite animal is {{fave_animal}}<h1>",
		TextPart: "My favourite animal is {{fave_animal}}",
	}

	err := template.ValidateData(map[string]string{
		"name": "Djale",
	})
	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, ErrMissingDefaultData))

	err = template.ValidateData(map[string]string{
		"name":        "Djale",
		"fave_animal": "Dog",
	})
	assert.Nil(t, err)

	// test failed to parse subject part
	template.SubjectPart = "Hello {{{name}}"

	err = template.ValidateData(map[string]string{
		"name":        "Djale",
		"fave_animal": "Dog",
	})
	assert.NotNil(t, err)

	// test failed to parse text part
	template.TextPart = "My favourite animal is {{{fave_animal}}"

	err = template.ValidateData(map[string]string{
		"name":        "Djale",
		"fave_animal": "Dog",
	})
	assert.NotNil(t, err)

	// test failed to parse html part
	template.HTMLPart = "<h1>My favourite animal is {{{fave_animal}}<h1>"

	err = template.ValidateData(map[string]string{
		"name":        "Djale",
		"fave_animal": "Dog",
	})
	assert.NotNil(t, err)
}

func TestGetters(t *testing.T) {
	now := time.Now()

	template := Template{
		BaseTemplate: BaseTemplate{
			Model: Model{
				ID:        212,
				CreatedAt: now,
				UpdatedAt: now,
			},
			UserID:      2,
			Name:        "test-template",
			SubjectPart: "Hello {{name}}",
		},
		HTMLPart:     "<h1>My favourite animal is {{fave_animal}}<h1>",
		TextPart:     "My favourite animal is {{fave_animal}}",
	}

	b := template.GetBase()
	assert.NotNil(t, b)
	assert.Equal(t, template.BaseTemplate, *b)

	id := template.GetID()
	assert.Equal(t, template.ID, id)

	tableName := template.BaseTemplate.TableName()
	assert.Equal(t, "templates", tableName)
}
