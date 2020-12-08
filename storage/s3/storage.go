package s3

import (
	"context"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"

	"github.com/mailbadger/app/entities"
)

const key = "s3_session"

type S3 interface {
	GetHTMLTemplate(userID int64, templateName string) (*s3.GetObjectOutput, error)
}

type storage struct {
}

// SetToContext sets the s3session to the context
func SetToContext(c *gin.Context, sess *session.Session) {
	c.Set(key, sess)
}

// GetFromContext returns the Storage associated with the context
func GetFromContext(c context.Context) S3 {
	return c.Value(key).(S3)
}

func GetHTMLTemplate(c context.Context, userID int64, templateName string) (*entities.Template, error) {
	return nil, nil
}
