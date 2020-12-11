package s3

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

const bucket = "BUCKET HERE"

// GetHTMLTemplate returns html part of the template saved in s3
func GetHTMLTemplate(c context.Context, userID int64, templateName string) (string, error) {
	obj, err := GetFromContext(c).GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fmt.Sprintf("/PATH_TO_FILE/%d/%s", userID, templateName)),
	})
	if err != nil {
		return "", fmt.Errorf("get object: %w", err)
	}

	var body []byte
	_, err = obj.Body.Read(body)
	if err != nil {
		return "", fmt.Errorf("read: %w", err)
	}

	return string(body), nil
}

