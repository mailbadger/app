package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/news-maily/api/queue"
	"github.com/sirupsen/logrus"
)

// Producer is a middleware that inits the Producer and attaches it to the context.
func Producer() gin.HandlerFunc {
	p, err := queue.NewProducer(os.Getenv("NSQD_HOST"), os.Getenv("NSQD_PORT"))
	if err != nil {
		logrus.Errorln(err)
	}

	return func(c *gin.Context) {
		queue.SetToContext(c, p)
		c.Next()
	}
}
