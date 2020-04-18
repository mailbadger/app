package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const reqIDKey = "reqid"

// RequestID is a middleware to set a random request id to the context and the response header in a form of:
// X-Request-Id: <uuid4>
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.Request.Header.Get("X-Request-Id")

		if requestID == "" {
			uuid4, err := uuid.NewRandom()
			if err != nil {
				logrus.WithError(err).Error("RequestID: Unable to generate uuid")
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": "We are unable to process the request at the moment. Please try again.",
				})
				return
			}
			requestID = uuid4.String()
		}

		// Expose it for use in the application
		c.Set(reqIDKey, requestID)

		// Set X-Request-Id header
		c.Writer.Header().Set("X-Request-Id", requestID)
		c.Next()
	}
}

// GetReqID fetches the request ID from the given context.
func GetReqID(ctx context.Context) string {
	return ctx.Value(reqIDKey).(string)
}
