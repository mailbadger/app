package templates

import (
	"context"
	"fmt"

	"github.com/cbroglie/mustache"

	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/storage"
	"github.com/mailbadger/app/storage/s3"
)

type Service interface {
	PostTemplate(c context.Context, input *entities.Template) error
	PutTemplate(c context.Context, input *entities.Template) error
}

type service struct {
}

func NewTemplateService() Service {
	return &service{}
}

func (s service) PostTemplate(c context.Context, input *entities.Template) error {

	// parse string to validate template params
	_, err := mustache.ParseString(input.HTMLPart)
	if err != nil {
		return fmt.Errorf("failed to parse html_part template error: %w", err)
	}
	// parse string to validate template params
	_, err = mustache.ParseString(input.TextPart)
	if err != nil {
		return fmt.Errorf("failed to parse text_part template error: %w", err)
	}

	err = storage.CreateTemplate(c, input)
	if err != nil {
		return fmt.Errorf("failed to create template error: %w", err)
	}

	err = s3.CreateTemplate(c, input)
	if err != nil {
		return fmt.Errorf("failed o create html template file to s3 error: %w", err)
	}

	return nil
}

func (s service) PutTemplate(c context.Context, input *entities.Template) error {

	// parse string to validate template params
	_, err := mustache.ParseString(input.HTMLPart)
	if err != nil {
		return fmt.Errorf("failed to parse html_part template error: %w", err)
	}
	// parse string to validate template params
	_, err = mustache.ParseString(input.TextPart)
	if err != nil {
		return fmt.Errorf("failed to parse text_part template error: %w", err)
	}

	err = storage.UpdateTemplate(c, input)
	if err != nil {
		return fmt.Errorf("failed to create template error: %w", err)
	}

	err = s3.CreateTemplate(c, input)
	if err != nil {
		return fmt.Errorf("failed o create html template file to s3 error: %w", err)
	}

	return nil
}
