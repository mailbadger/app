package s3

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

const bucket = "BUCKET HERE"

var (
	ErrDeleteFailed = errors.New("failed to delete")
)

// DeleteHTMLTemplate deletes html part of the template saved in s3
func DeleteHTMLTemplate(c context.Context, userID int64, templateName string) error {
	obj, err := GetFromContext(c).DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fmt.Sprintf("/PATH_TO_FILE/%d/%s", userID, templateName)),
	})
	if err != nil {
		return fmt.Errorf("delete object: %w", err)
	}

	if !aws.BoolValue(obj.DeleteMarker) {
		return ErrDeleteFailed
	}

	return nil
}
