package s3

import (
	"bytes"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/mailbadger/app/entities"
)

func (s *s3storage) CreateHTMLTemplate(html string, bucket string, tmplInput *entities.Template) error {

	input := &s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("AWS_S3_BUCKET")),
		Key:    aws.String(fmt.Sprintf("/PATH_TO_FILE/%d/%s", tmplInput.MerchantID, tmplInput.Name)),
		Body:   bytes.NewReader([]byte(html)),
	}
	_, err := s.s3client.PutObject(input)
	if err != nil {
		return fmt.Errorf("failed to insert html part to s3 error: %w", err)
	}
	return nil
}
