package entities

import (
	"errors"
	"fmt"
	"time"

	"github.com/cbroglie/mustache"
)

var (
	ErrMissingDefaultData = errors.New("missing default data")
)

type Template struct {
	Model
	UserID      int64  `json:"user_id"`
	Name        string `json:"name"`
	HTMLPart    string `json:"html_part" gorm:"-"`
	TextPart    string `json:"text_part"`
	SubjectPart string `json:"subject_part"`
}

func (t Template) ValidateData(data map[string]string) error {
	err := validateData(t.SubjectPart, data)
	if err != nil {
		return fmt.Errorf("validate subject part: %w", err)
	}

	err = validateData(t.TextPart, data)
	if err != nil {
		return fmt.Errorf("validate text part: %w", err)
	}

	err = validateData(t.HTMLPart, data)
	if err != nil {
		return fmt.Errorf("validate html part: %w", err)
	}

	return nil
}

func validateData(templateString string, data map[string]string) error {
	template, err := mustache.ParseString(templateString)
	if err != nil {
		return fmt.Errorf("parse string: %w", err)
	}

	for _, tag := range template.Tags() {
		_, exist := data[tag.Name()]
		if !exist {
			return fmt.Errorf("%s tag: %w", tag.Name(), ErrMissingDefaultData)
		}
	}

	return nil
}

type TemplateCollection struct {
	NextToken  string         `json:"next_token"`
	Collection []TemplateMeta `json:"collection"`
}

type TemplateMeta struct {
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
}

type TemplatesCollectionItem struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	SubjectPart string    `json:"subject_part"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (c TemplatesCollectionItem) GetID() int64 {
	return c.ID
}
