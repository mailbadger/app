package templates

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/cbroglie/mustache"

	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/storage"
)

var (
	templatesBucket = os.Getenv("TEMPLATE_BUCKET")

	ErrParseHTMLPart    = errors.New("failed to parse HTMLPart")
	ErrParseTextPart    = errors.New("failed to parse TextPart")
	ErrParseSubjectPart = errors.New("failed to parse SubjectPart")
)

type Service interface {
	AddTemplate(c context.Context, input *entities.Template) error
	UpdateTemplate(c context.Context, input *entities.Template) error
	DeleteTemplate(c context.Context, templateID, userID int64) error
}

// service implements the Service interface
type service struct {
	db storage.Storage
	s3 s3iface.S3API
}

func NewTemplateService(db storage.Storage, s3 s3iface.S3API) Service {
	return &service{
		db: db,
		s3: s3,
	}
}

func (s service) AddTemplate(c context.Context, template *entities.Template) error {

	// parse string to validate template params
	_, err := mustache.ParseString(template.HTMLPart)
	if err != nil {
		return ErrParseHTMLPart
	}
	// parse string to validate template params
	_, err = mustache.ParseString(template.TextPart)
	if err != nil {
		return ErrParseTextPart
	}
	// parse string to validate template params
	_, err = mustache.ParseString(template.SubjectPart)
	if err != nil {
		return ErrParseSubjectPart
	}

	s3Input := &s3.PutObjectInput{
		Bucket: aws.String(templatesBucket),
		Key:    aws.String(fmt.Sprintf("%d/%d", template.UserID, template.ID)),
		Body:   bytes.NewReader([]byte(template.HTMLPart)),
	}

	_, err = s.s3.PutObject(s3Input)
	if err != nil {
		return fmt.Errorf("upload template: put s3 object: %w", err)
	}

	err = s.db.CreateTemplate(template)
	if err != nil {
		return fmt.Errorf("create template: %w", err)
	}

	return nil
}

func (s service) UpdateTemplate(c context.Context, template *entities.Template) error {

	// parse string to validate template params
	_, err := mustache.ParseString(template.HTMLPart)
	if err != nil {
		return ErrParseHTMLPart
	}
	// parse string to validate template params
	_, err = mustache.ParseString(template.TextPart)
	if err != nil {
		return ErrParseTextPart
	}
	// parse string to validate template params
	_, err = mustache.ParseString(template.SubjectPart)
	if err != nil {
		return ErrParseSubjectPart
	}

	s3Input := &s3.PutObjectInput{
		Bucket: aws.String(templatesBucket),
		Key:    aws.String(fmt.Sprintf("%d/%d", template.UserID, template.ID)),
		Body:   bytes.NewReader([]byte(template.HTMLPart)),
	}

	_, err = s.s3.PutObject(s3Input)
	if err != nil {
		return fmt.Errorf("upload template: put s3 object: %w", err)
	}

	err = s.db.UpdateTemplate(template)
	if err != nil {
		return fmt.Errorf("update template: %w", err)
	}

	return nil
}

// DeleteTemplate deletes the given template
func (s *service) DeleteTemplate(c context.Context, templateID, userID int64) error {
	_, err := s.db.GetTemplate(templateID, userID)
	if err != nil {
		return fmt.Errorf("get template: %w", err)
	}

	_, err = s.s3.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(templatesBucket),
		Key:    aws.String(fmt.Sprintf("%d/%d", userID, templateID)),
	})

	err = s.db.DeleteTemplate(templateID, userID)
	if err != nil {
		return fmt.Errorf("delete template: %w", err)
	}

	return nil
}