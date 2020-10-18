package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/mailbadger/app/storage"
)

// Storage is a middleware that inits the Storage and attaches it to the context.
func Storage(s storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		storage.SetToContext(c, s)
		c.Next()
	}
}
