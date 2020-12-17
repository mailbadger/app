package templates

import (
	"context"
	"errors"
	"fmt"

	"github.com/cbroglie/mustache"

	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/storage"
	"github.com/mailbadger/app/storage/s3"
)

var (
	ErrParseHTMLPart    = errors.New("failed to parse HTMLPart")
	ErrParseTextPart    = errors.New("failed to parse TextPart")
	ErrParseSubjectPart = errors.New("failed to parse SubjectPart")
)

type Service interface {
	AddTemplate(c context.Context, input *entities.Template) error
	UpdateTemplate(c context.Context, input *entities.Template) error
}

type service struct {
}

func NewTemplateService() Service {
	return &service{}
}

func (s service) AddTemplate(c context.Context, input *entities.Template) error {

	// parse string to validate template params
	_, err := mustache.ParseString(input.HTMLPart)
	if err != nil {
		return ErrParseHTMLPart
	}
	// parse string to validate template params
	_, err = mustache.ParseString(input.TextPart)
	if err != nil {
		return ErrParseTextPart
	}
	// parse string to validate template params
	_, err = mustache.ParseString(input.SubjectPart)
	if err != nil {
		return ErrParseSubjectPart
	}

	err = storage.CreateTemplate(c, input)
	if err != nil {
		return fmt.Errorf("failed to create template error: %w", err)
	}

	err = s3.CreateTemplate(c, input)
	if err != nil {
		return fmt.Errorf("failed to create html template file to s3 error: %w", err)
	}

	return nil
}

func (s service) UpdateTemplate(c context.Context, input *entities.Template) error {

	// parse string to validate template params
	_, err := mustache.ParseString(input.HTMLPart)
	if err != nil {
		return ErrParseHTMLPart
	}
	// parse string to validate template params
	_, err = mustache.ParseString(input.TextPart)
	if err != nil {
		return ErrParseTextPart
	}
	// parse string to validate template params
	_, err = mustache.ParseString(input.SubjectPart)
	if err != nil {
		return ErrParseSubjectPart
	}

	err = storage.UpdateTemplate(c, input)
	if err != nil {
		return fmt.Errorf("failed to create template error: %w", err)
	}

	err = s3.CreateTemplate(c, input)
	if err != nil {
		return fmt.Errorf("failed to create html template file to s3 error: %w", err)
	}

	return nil
}
