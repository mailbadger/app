package s3

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/mailbadger/app/entities"
)

const key = "s3"

type S3Storage interface {
	CreateHTMLTemplate(html string, bucket string, tmplInput *entities.Template) error
}

// SetToContext sets the s3session to the context
func SetToContext(c *gin.Context, storage S3Storage) {
	c.Set(key, storage)
}

// GetFromContext returns the Storage associated with the context
func GetFromContext(c context.Context) S3Storage {
	return c.Value(key).(S3Storage)
}

// CreateHTMLTemplate uploads html file to s3.
func CreateHTMLTemplate(c context.Context, html string, bucket string, tmplInput *entities.Template) error {
	return GetFromContext(c).CreateHTMLTemplate(html, bucket, tmplInput)
}
