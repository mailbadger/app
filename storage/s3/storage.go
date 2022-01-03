package s3

import (
	"context"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"

	"github.com/gin-gonic/gin"
)

func NewClient() (*s3.S3, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	return s3.New(sess), nil
}

const key = "s3"

// SetToContext sets the s3 client interface to the context
func SetToContext(c *gin.Context, s3client s3iface.S3API) {
	c.Set(key, s3client)
}

// GetFromContext returns the s3 client interface associated with the context
func GetFromContext(c context.Context) s3iface.S3API {
	return c.Value(key).(s3iface.S3API)
}
