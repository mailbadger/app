package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const reqIDKey = "reqid"

// RequestID is a middleware to set a random request id to the context and the response header in a form of:
// X-Request-Id: <uuid4>
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.Request.Header.Get("X-Request-Id")

		if requestID == "" {
			uuid4 := uuid.New()
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
