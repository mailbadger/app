package templates

import (
	"context"
	"fmt"

	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/storage"
	"github.com/mailbadger/app/storage/s3"
)

var (
	bucket = "bucketo za templates"
)

type Service interface {
	PostTemplate(c context.Context, input *entities.Template) error
}

type service struct {
}

func NewTemplateService() Service {
	return &service{}
}

func (s service) PostTemplate(c context.Context, input *entities.Template) error {
	err := storage.CreateTemplate(c, input)
	if err != nil {
		return fmt.Errorf("failed to create template error: %w", err)
	}

	err = s3.CreateHTMLTemplate(c, input.HTMLPart, bucket, input)
	if err != nil {
		return fmt.Errorf("failed o create html template file to s3 error: %w", err)
	}

	return nil
}
