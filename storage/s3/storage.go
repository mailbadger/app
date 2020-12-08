package s3

import (
	"context"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gin-gonic/gin"
)

const key = "s3_session"

type S3 interface {
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
