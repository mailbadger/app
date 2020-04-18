package logger

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const key = "logger"

// From returns a logger entry from the given context.
func From(ctx context.Context) *logrus.Entry {
	return ctx.Value(key).(*logrus.Entry)
}

// SetToContext sets the given entry in the context.
func SetToContext(ctx *gin.Context, entry *logrus.Entry) {
	ctx.Set(key, entry)
}
