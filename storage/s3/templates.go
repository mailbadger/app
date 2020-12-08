package s3

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

const bucket = "BUCKET HERE"

// DeleteHTMLTemplate deletes html part of the template saved in s3
func (s *s3storage) DeleteHTMLTemplate(userID int64, templateName string) error {
	obj, err := s.s3client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fmt.Sprintf("/PATH_TO_FILE/%d/%s", userID, templateName)),
	})
	if err != nil {
		return fmt.Errorf("delete object: %w", err)
	}

	if !aws.BoolValue(obj.DeleteMarker) {
		return errors.New("failed to delete")
	}

	return nil
}
