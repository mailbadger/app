package middleware

import (
	"os"

	"github.com/FilipNikolovski/news-maily/storage"
	"github.com/gin-gonic/gin"
)

// Storage is a middleware that inits the Storage and attaches it to the context.
func Storage() gin.HandlerFunc {
	s := storage.New(os.Getenv("DATABASE_DRIVER"), os.Getenv("DATABASE_CONFIG"))

	return func(c *gin.Context) {
		storage.SetToContext(c, s)
		c.Next()
	}
}
