package s3

import (
	"context"

	"github.com/gin-gonic/gin"
)

const key = "s3"

type S3Storage interface {
	DeleteHTMLTemplate(userID int64, templateName string) error
}

// SetToContext sets the s3session to the context
func SetToContext(c *gin.Context, storage S3Storage) {
	c.Set(key, storage)
}

// GetFromContext returns the Storage associated with the context
func GetFromContext(c context.Context) S3Storage {
	return c.Value(key).(S3Storage)
}

// DeleteHTMLTemplate deletes html part of the template saved in s3
func DeleteHTMLTemplate(c context.Context, userID int64, templateName string) error {
	return GetFromContext(c).DeleteHTMLTemplate(userID, templateName)
}
