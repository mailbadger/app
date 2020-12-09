package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/mailbadger/app/storage/s3"
)

// S3Storage is a middleware that inits the S3Storage and attaches it to the context.
func S3Storage(storage s3.S3Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		s3.SetToContext(c, storage)
		c.Next()
	}
}
