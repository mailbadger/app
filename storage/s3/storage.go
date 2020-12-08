package s3

import (
	"context"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"

	"github.com/mailbadger/app/entities"
)

const key = "s3"

type S3Storage interface {
	GetHTMLTemplate(userID int64, templateName string) (*s3.GetObjectOutput, error)
}

// SetToContext sets the s3session to the context
func SetToContext(c *gin.Context, storage S3Storage) {
	c.Set(key, storage)
}

// GetFromContext returns the Storage associated with the context
func GetFromContext(c context.Context) S3Storage {
	return c.Value(key).(S3Storage)
}

func GetHTMLTemplate(c context.Context, userID int64, templateName string) (*entities.Template, error) {
	return nil, nil
}
