package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/news-maily/app/storage"
)

// Storage is a middleware that inits the Storage and attaches it to the context.
func Storage() gin.HandlerFunc {
	driver := os.Getenv("DATABASE_DRIVER")
	config := storage.MakeConfigFromEnv(driver)
	s := storage.New(driver, config)

	return func(c *gin.Context) {
		storage.SetToContext(c, s)
		c.Next()
	}
}
