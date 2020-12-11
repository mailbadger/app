package middleware

import (
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/gin-gonic/gin"

	"github.com/mailbadger/app/storage/s3"
)

// S3Client is a middleware that inits the S3 client interface and attaches it to the context.
func S3Client(client s3iface.S3API) gin.HandlerFunc {
	return func(c *gin.Context) {
		s3.SetToContext(c, client)
		c.Next()
	}
}
