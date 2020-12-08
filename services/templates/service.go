package templates

import (
	"context"
	"fmt"

	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/storage"
	s3storage "github.com/mailbadger/app/storage/s3"
)

// TemplateService contains all methods for operating with templates
type TemplateService interface {
	DeleteTemplate(c context.Context, template *entities.Template) error
}

// service implements the Service interface
type service struct {
}

// NewTemplateService returns a new Service
func NewTemplateService() TemplateService {
	return &service{}
}

// DeleteTemplate deletes the given template
func (s service) DeleteTemplate(c context.Context, template *entities.Template) error {
	err := s3storage.DeleteHTMLTemplate(c, template.UserID, template.Name)
	if err != nil {
		return fmt.Errorf("delete html template: %w", err)
	}

	err = storage.DeleteTemplate(c, template.ID, template.UserID)
	if err != nil {
		return fmt.Errorf("delete template: %w", err)
	}

	return nil
}
