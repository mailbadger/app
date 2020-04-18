package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/news-maily/app/logger"
	"github.com/sirupsen/logrus"
)

// SetLoggerEntry sets a request logger entry to the context, along with fields 'request_id' and 'user_id'
// which will propagate in each logged message.
func SetLoggerEntry() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqLogger := logrus.StandardLogger()
		fields := logrus.Fields{}

		reqID := GetReqID(c)
		if reqID != "" {
			fields["request_id"] = reqID
		}

		u := GetUser(c)
		if u != nil {
			fields["user_id"] = u.ID
		}

		logger.SetToContext(c, reqLogger.WithFields(fields))

		c.Next()
	}
}
