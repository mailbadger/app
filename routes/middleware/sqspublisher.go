package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/mailbadger/app/sqs"
)

// SQSPublisher is a middleware that adds the SQS publisher to the context.
func SQSPublisher(pub sqs.Publisher) gin.HandlerFunc {
	return func(c *gin.Context) {
		sqs.SetPublisherToContext(c, pub)
		c.Next()
	}
}
