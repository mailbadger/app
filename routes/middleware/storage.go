package middleware

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/news-maily/api/storage"
)

// Storage is a middleware that inits the Storage and attaches it to the context.
func Storage() gin.HandlerFunc {
	driver := os.Getenv("DATABASE_DRIVER")
	config := makeConfigFromEnv(driver)
	s := storage.New(driver, config)

	return func(c *gin.Context) {
		storage.SetToContext(c, s)
		c.Next()
	}
}

func makeConfigFromEnv(driver string) string {
	switch driver {
	case "sqlite3":
		return os.Getenv("SQLITE3_FILE")
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true",
			os.Getenv("MYSQL_USER"),
			os.Getenv("MYSQL_PASS"),
			os.Getenv("MYSQL_HOST"),
			os.Getenv("MYSQL_DATABASE"),
		)
	default:
		return ""
	}
}
