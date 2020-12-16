package s3

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/gin-gonic/gin"

	"github.com/mailbadger/app/entities"
)

const key = "s3"

// SetToContext sets the s3 client interface to the context
func SetToContext(c *gin.Context, s3client s3iface.S3API) {
	c.Set(key, s3client)
}

// GetFromContext returns the s3 client interface associated with the context
func GetFromContext(c context.Context) s3iface.S3API {
	return c.Value(key).(s3iface.S3API)
}

// PutS3Object uploads a object to s3.
func PutS3Object(c context.Context, input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	return GetFromContext(c).PutObject(input)
}

// CreateTemplate uploads html file to s3.
func CreateTemplate(c context.Context, tmplInput *entities.Template) error {

	input := &s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("AWS_S3_BUCKET")),
		Key:    aws.String(fmt.Sprintf("/PATH_TO_FILE/%d/%s", tmplInput.UserID, tmplInput.Name)),
		Body:   bytes.NewReader([]byte(tmplInput.HTMLPart)),
	}

	_, err := PutS3Object(c, input)
	if err != nil {
		return fmt.Errorf("failed to insert html part to s3 error: %w", err)
	}
	return nil
}
