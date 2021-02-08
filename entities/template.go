package entities

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cbroglie/mustache"
	"golang.org/x/sync/errgroup"
)

var (
	ErrMissingDefaultData = errors.New("missing default data")
)

const (
	TagName = "name"
	TagUnsubscribeUrl = "unsubscribe_url"
)

// BaseTemplate represents the base params of each template
type BaseTemplate struct {
	Model
	UserID      int64  `json:"user_id"`
	Name        string `json:"name"`
	SubjectPart string `json:"subject_part"`
}

// GetID returns the id of the template
func (c BaseTemplate) GetID() int64 {
	return c.ID
}

// TableName overrides the table name used by BaseTemplate to `templates`
func (BaseTemplate) TableName() string {
	return "templates"
}

// Template represents the email body template
type Template struct {
	BaseTemplate
	HTMLPart string `json:"html_part" gorm:"-"`
	TextPart string `json:"text_part"`
}

// GetBase returns the base of the template
func (t Template) GetBase() *BaseTemplate {
	return &BaseTemplate{
		Model: Model{
			ID:        t.ID,
			CreatedAt: t.CreatedAt,
			UpdatedAt: t.UpdatedAt,
		},
		UserID:      t.UserID,
		Name:        t.Name,
		SubjectPart: t.SubjectPart,
	}
}

// ValidateData checks if all template tags are covered with provided data
func (t Template) ValidateData(data map[string]string) error {
	g, _ := errgroup.WithContext(context.Background())

	g.Go(func() error {
		err := validateData(t.SubjectPart, data)
		if err != nil {
			return fmt.Errorf("validate subject part: %w", err)
		}

		return nil
	})

	g.Go(func() error {
		err := validateData(t.TextPart, data)
		if err != nil {
			return fmt.Errorf("validate text part: %w", err)
		}

		return nil
	})

	g.Go(func() error {
		err := validateData(t.HTMLPart, data)
		if err != nil {
			return fmt.Errorf("validate html part: %w", err)
		}

		return nil
	})

	return g.Wait()
}

func validateData(templateString string, data map[string]string) error {
	template, err := mustache.ParseString(templateString)
	if err != nil {
		return fmt.Errorf("parse string: %w", err)
	}

	for _, tag := range template.Tags() {
		if tag.Name() == TagName || tag.Name() == TagUnsubscribeUrl {
			continue
		}
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
