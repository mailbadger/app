package templates

import (
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/news-maily/api/emails"
)

type Storage interface {
	GetTemplate(input *ses.GetTemplateInput) (*ses.GetTemplateOutput, error)
	ListTemplates(input *ses.ListTemplatesInput) (*ses.ListTemplatesOutput, error)
	DeleteTemplate(input *ses.DeleteTemplateInput) (*ses.DeleteTemplateOutput, error)
	CreateTemplate(input *ses.CreateTemplateInput) (*ses.CreateTemplateOutput, error)
	UpdateTemplate(input *ses.UpdateTemplateInput) (*ses.UpdateTemplateOutput, error)
}

func NewSesTemplateStore(key, secret, region string) (Storage, error) {
	return emails.NewSESClient(key, secret, region)
}
