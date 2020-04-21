package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/news-maily/app/queue"
	"github.com/sirupsen/logrus"
)

// Producer is a middleware that inits the Producer and attaches it to the context.
func Producer() gin.HandlerFunc {
	return func(c *gin.Context) {
		p, err := queue.NewProducer(os.Getenv("NSQD_HOST"), os.Getenv("NSQD_PORT"))
		if err != nil {
			logrus.WithError(err).Error("unable to instantiate queue producer")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "We are unable to process the request. Please try again.",
			})
			return
		}

		queue.SetToContext(c, p)
		c.Next()
	}
}
