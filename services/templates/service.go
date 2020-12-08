package templates

import (
	"context"
	"fmt"

	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/storage"
	"github.com/mailbadger/app/storage/s3"
)

type Service interface {
	GetTemplate(c context.Context, templateID int64, userID int64) (*entities.Template, error)
}

type service struct {
}

func NewTemplateService() Service {
	return &service{}
}

// GetTemplate returns the template with given template id and user id
func (s service) GetTemplate(c context.Context, templateID int64, userID int64) (*entities.Template, error) {
	template, err := storage.GetTemplate(c, templateID, userID)
	if err != nil {
		return nil, fmt.Errorf("get template: %w", err)
	}

	html, err := s3.GetHTMLTemplate(c, userID, template.Name)
	if err != nil {
		return nil, fmt.Errorf("get html template: %w", err)
	}

	template.HTMLPart = html

	return template, err
}
