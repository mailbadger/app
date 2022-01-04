package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mailbadger/app/logger"
	"github.com/sirupsen/logrus"
)

// Logger provides logging middleware.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := GetReqID(c)
		entry := logger.From(c).WithField("request_id", reqID)
		logger.SetToContext(c, entry)

		start := time.Now()

		path := c.Request.URL.Path
		c.Next()

		end := time.Now().UTC()
		latency := end.Sub(start)

		entry = logger.From(c).WithFields(logrus.Fields{
			"status":     c.Writer.Status(),
			"method":     c.Request.Method,
			"path":       path,
			"ip":         c.ClientIP(),
			"latency":    latency,
			"user-agent": c.Request.UserAgent(),
			"time":       end.Format(time.RFC3339),
		})

		if len(c.Errors) > 0 {
			// Append error field if this is an erroneous request.
			entry.Error(c.Errors.String())
		} else {
			entry.Info()
		}
	}
}
